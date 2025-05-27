package parser

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type RarbgParser struct {
	BaseURL string
}

func NewRarbgParser(mirrorURL string) *RarbgParser {
	return &RarbgParser{
		BaseURL: mirrorURL,
	}
}

func (r *RarbgParser) Search(query string) ([]TorrentFile, error) {
	searchURL := fmt.Sprintf("%ssearch/?search=%s", r.BaseURL, strings.ReplaceAll(query, " ", "+"))
	fmt.Println("🌍 Search URL:", searchURL)
	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("received non-OK status code from RARBG: %s", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("⚠️ Warning: failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fail to parse HTML: %w", err)
	}

	var torrents []TorrentFile

	doc.Find("table.lista2t tr.lista2").Each(func(i int, s *goquery.Selection) {
		if torrent := r.ParseTableRow(s); torrent != nil {
			torrents = append(torrents, *torrent)
		}
	})

	fmt.Printf("📦 Found %d torrents\n", len(torrents))
	return torrents, nil
}

func (r *RarbgParser) ParseTableRow(s *goquery.Selection) *TorrentFile {
	linkElement := s.Find("td.lista").Eq(1).Find("a").First()

	href, exists := linkElement.Attr("href")
	if !exists {
		return nil // skip invalid rows
	}

	torrent := &TorrentFile{
		Name: r.extractName(linkElement),
		Href: href,
	}

	// Extract all basic metadata
	r.extractCategory(s, torrent)
	r.extractUploadDate(s, torrent)
	r.extractSize(s, torrent)
	r.extractSeeders(s, torrent)
	r.extractLeechers(s, torrent)
	r.extractUploaders(s, torrent)

	return torrent
}

// Name
func (r *RarbgParser) extractName(linkElement *goquery.Selection) string {
	return strings.TrimSpace(linkElement.Text())
}

// Category
func (r *RarbgParser) extractCategory(s *goquery.Selection, torrent *TorrentFile) {
	if category := strings.TrimSpace(s.Find("td.lista").Eq(2).Text()); category != "" {
		torrent.Category = category
	}
}

// UploadDate
func (r *RarbgParser) extractUploadDate(s *goquery.Selection, torrent *TorrentFile) {
	if uploadDate := strings.TrimSpace(s.Find("td.lista").Eq(3).Text()); uploadDate != "" {
		torrent.UploadDate = uploadDate
	}
}

// Size
func (r *RarbgParser) extractSize(s *goquery.Selection, torrent *TorrentFile) {
	if sizeText := strings.TrimSpace(s.Find("td.lista").Eq(4).Text()); sizeText != "" {
		torrent.SizeRaw = sizeText

		sizeFloat := ParseSizeToGB(sizeText)
		torrent.Size = sizeFloat
	}
}

// Seeders
func (r *RarbgParser) extractSeeders(s *goquery.Selection, torrent *TorrentFile) {
	if seedersText := strings.TrimSpace(s.Find("td.lista").Eq(5).Text()); seedersText != "" {
		if seedersInt, err := strconv.Atoi(seedersText); err != nil {
			torrent.Seeders = seedersInt
		}
	}
}

// Leechers
func (r *RarbgParser) extractLeechers(s *goquery.Selection, torrent *TorrentFile) {
	// Leechers
	if leechersText := strings.TrimSpace(s.Find("td.lista").Eq(6).Text()); leechersText != "" {
		if leechersInt, err := strconv.Atoi(leechersText); err != nil {
			torrent.Leechers = leechersInt
		}
	}
}

func (r *RarbgParser) extractUploaders(s *goquery.Selection, torrent *TorrentFile) {
	// Uploader
	if uploader := strings.TrimSpace(s.Find("td.lista").Eq(7).Text()); uploader != "" {
		torrent.Uploader = uploader
	}
}

func (r *RarbgParser) FetchTorrentDetails(torrent *TorrentFile) error {
	detailURL := torrent.Href
	if strings.HasPrefix(torrent.Href, "/") {
		detailURL = r.BaseURL + strings.TrimPrefix(torrent.Href, "/")
	}

	// TODO:
	// fmt.Printf("🔍 Fetching details for: %s\n", torrent.Name)

	resp, err := http.Get(detailURL)
	if err != nil {
		return fmt.Errorf("failed to fetch torrent details: %w", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("⚠️ Warning: failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to parse detail HTML: %w", err)
	}

	// getting magnetlink
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			if strings.Contains(href, "magnet:") && torrent.MagnetLink == "" {
				torrent.MagnetLink = href
			}
		}
	})

	return nil
}

func (r *RarbgParser) EnrichTorrents(torrents []TorrentFile) []TorrentFile {
	var enrichedTorrents []TorrentFile

	for _, torrent := range torrents {
		torrentCopy := torrent

		if err := r.FetchTorrentDetails(&torrentCopy); err != nil {
			fmt.Printf("⚠️  Failed to fetch details for %s: %v\n", torrent.Name, err)
		}

		enrichedTorrents = append(enrichedTorrents, torrentCopy)
	}

	return enrichedTorrents
}
