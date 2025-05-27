package parser

import (
	"fmt"
	"strings"
)

func NewParser(siteName string, mirrorURL string) (TorrentParser, error) {
	switch strings.ToLower(siteName) {
	case "rarbg":
		return NewRarbgParser(mirrorURL), nil

	default:
		return nil, fmt.Errorf("usupported site: %s", siteName)
	}
}
