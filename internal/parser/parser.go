package parser

import (
	"regexp"
	"strconv"
	"strings"
)

type TorrentFile struct {
	Name       string
	Href       string
	Size       float64
	SizeRaw    string
	Seeders    int
	Leechers   int
	Uploader   string
	MagnetLink string
	Language   string
	MetaInfo   string
	Source     string
	UploadDate string
	Category   string
	Resolution string
	Trusted    bool
	VideoCodec string
	AudioCodec string
	Container  string
	BitDepth   string
}

type TorrentParser interface {
	Search(query string) ([]TorrentFile, error)
	EnrichTorrents(torrents []TorrentFile) []TorrentFile
}

func ParseSizeToGB(sizeStr string) float64 {
	re := regexp.MustCompile(`([0-9.]+)\s*([A-Za-z]+)`)
	matches := re.FindStringSubmatch(strings.TrimSpace(sizeStr))

	if len(matches) < 3 {
		return 0
	}

	sizeFloat, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0
	}

	unit := strings.ToUpper(matches[2])

	switch unit {
	case "KB":
		return sizeFloat / (1024 * 1024)
	case "MB":
		return sizeFloat / 1024
	case "GB":
		return sizeFloat
	case "TB":
		return sizeFloat * 1024
	default:
		return 0
	}
}
