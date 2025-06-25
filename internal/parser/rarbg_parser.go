// TODO:
package parser

import (
	"errors"
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
	req, _ := http.NewRequest("GET", searchURL, nil)
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/114.0.0.0 Safari/537.36",
	)
	resp, err := http.DefaultClient.Do(req)
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

	if len(torrents) <= 0 {
		return torrents, errors.New("none")
	}
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
		parts := strings.Split(category, "/")
		topCategory := parts[0]
		torrent.Category = topCategory
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
		if seedersInt, err := strconv.Atoi(seedersText); err == nil {
			torrent.Seeders = seedersInt
		}
	}
}

// Leechers
func (r *RarbgParser) extractLeechers(s *goquery.Selection, torrent *TorrentFile) {
	// Leechers
	if leechersText := strings.TrimSpace(s.Find("td.lista").Eq(6).Text()); leechersText != "" {
		if leechersInt, err := strconv.Atoi(leechersText); err == nil {
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

	req, _ := http.NewRequest("GET", detailURL, nil)
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/114.0.0.0 Safari/537.36",
	)

	resp, err := http.DefaultClient.Do(req)
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

	doc.Find("table.lista tr").Each(func(i int, s *goquery.Selection) {
		header := strings.TrimSpace(s.Find("td.header2").Text())
		value := strings.TrimSpace(s.Find("td.lista").Text())

		switch header {
		case "Description:":
			torrent.MetaInfo = value
			r.parseVideoSpecs(value, torrent)

		case "Language:":
			torrent.Language = value

		case "Downloads:":
			if downloads, err := strconv.Atoi(value); err == nil {
				torrent.Downloads = downloads
			}
		}
	})

	r.parseFilenameMetaData(torrent.Name, torrent)

	torrent.Trusted = r.isTrustedUploader(torrent.Uploader)

	return nil
}

func (r *RarbgParser) parseVideoSpecs(description string, torrent *TorrentFile) {
	descLower := strings.ToLower(description)
	metaInfoLower := strings.ToLower(torrent.MetaInfo)

	// Video Codec
	if strings.Contains(descLower, "hevc") || strings.Contains(descLower, "x265") ||
		strings.Contains(descLower, "h265") || strings.Contains(descLower, "avc") {
		torrent.VideoCodec = "HEVC/x265"
	} else if strings.Contains(descLower, "x264") || strings.Contains(descLower, "h264") || strings.Contains(descLower, "h.264") {
		torrent.VideoCodec = "x264"
	} else if strings.Contains(descLower, "av1") {
		torrent.VideoCodec = "AV1"
	} else {
		torrent.VideoCodec = "Unknown"
	}

	// Audio Codec
	if strings.Contains(descLower, "atmos") {
		torrent.AudioCodec = "Dolby Atmos"
	} else if strings.Contains(descLower, "dts-hd") || strings.Contains(descLower, "truehd") {
		torrent.AudioCodec = "DTS-HD/TrueHD"
	} else if strings.Contains(descLower, "dts") {
		torrent.AudioCodec = "DTS"
	} else if strings.Contains(descLower, "aac") {
		torrent.AudioCodec = "AAC"
	} else if strings.Contains(descLower, "opus") {
		torrent.AudioCodec = "OPUS"
	} else if strings.Contains(descLower, "mp3") {
		torrent.AudioCodec = "MP3"
	} else if strings.Contains(descLower, "eac3") || strings.Contains(descLower, "ddp") {
		torrent.AudioCodec = "EAC3"
	} else if strings.Contains(descLower, "ac3") || strings.Contains(descLower, "dd5.1") {
		torrent.AudioCodec = "AC3"
	} else {
		torrent.AudioCodec = "Unknown"
	}

	// Containers
	if strings.Contains(descLower, "matroska") || strings.Contains(torrent.Name, "mkv") {
		torrent.Container = "Matroska/MKV"
	} else if strings.Contains(torrent.Name, ".mp4") {
		torrent.Container = "MP4"
	} else {
		torrent.Container = "Unknown"
	}

	// Bit depth
	bitDepth := "8-bit"
	if strings.Contains(descLower, "10bit") || strings.Contains(descLower, "10-bit") ||
		strings.Contains(metaInfoLower, "10bit") || strings.Contains(metaInfoLower, "10-bit") ||
		strings.Contains(descLower, "hdr") || strings.Contains(metaInfoLower, "hdr") ||
		strings.Contains(
			descLower,
			"dolby vision",
		) || strings.Contains(metaInfoLower, "dolby vision") {
		bitDepth = "10-bit"
	}
	torrent.BitDepth = bitDepth
}

// Resolution
func (r *RarbgParser) parseFilenameMetaData(filename string, torrent *TorrentFile) {
	name := strings.ToUpper(filename)

	switch {
	case strings.Contains(name, "2160P") || strings.Contains(name, "4K") || strings.Contains(name, "UHD"):
		torrent.Resolution = "2160P"

	case strings.Contains(name, "1080P") || strings.Contains(name, "FHD"):
		torrent.Resolution = "1080P"

	case strings.Contains(name, "720P"):
		torrent.Resolution = "720P"

	case strings.Contains(name, "480P") || strings.Contains(name, "SD"):
		torrent.Resolution = "480P"

	default:
		torrent.Resolution = "Unknown"
	}

	// Source
	if strings.Contains(name, "IMAX") {
		torrent.Source = "IMAX"
	} else if strings.Contains(name, "BLURAY") || strings.Contains(name, "BLU-RAY") || strings.Contains(name, "BDREMUX") || strings.Contains(name, "BDRIP") {
		torrent.Source = "BLURAY"
	} else if strings.Contains(name, "WEBRIP") || strings.Contains(name, "WEB-DL") || strings.Contains(name, "WEB") || strings.Contains(name, "AMZN") || strings.Contains(name, "NF") || strings.Contains(name, "HMAX") {
		torrent.Source = "WEB"
	} else if strings.Contains(name, "CAM") || strings.Contains(name, "CAMRIP") {
		torrent.Source = "CAM"
	} else {
		torrent.Source = "Unknown"
	}
}

// Trusted Uploaders

func (r *RarbgParser) isTrustedUploader(uploader string) bool {
	trustedUploaders := []string{
		"RARBG", "YTS", "ETRG", "Prof", "PMEDIA", "Wrath", "FGT",
		"SPARKS", "UTR", "PSA", "DON", "GalaxyRG", "GalaxyTV",
		"QxR", "Tigole", "CtrlHD", "NTb", "TBS", "RMTeam", "Judas",
		"SUSPENSE", "EBP", "icecracked", "DataDiva", "Accid",
		"1DNCreW", "bone111", "NikaNika", "Maxoverpower", "IONICBOII",
		"Petehollow", "Telly", "mkvCinemas", "TAoE", "prudence25",
	}

	uploaderLower := strings.ToLower(uploader)
	for _, trusted := range trustedUploaders {
		if strings.Contains(uploaderLower, strings.ToLower(trusted)) {
			return true
		}
	}

	return false
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
