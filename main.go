package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sanjaix21/krakeneye/internal/display"
	"github.com/sanjaix21/krakeneye/internal/parser"
	"github.com/sanjaix21/krakeneye/internal/sites"
)

func _getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("🔍 Enter search query (e.g. interstellar 2014): ")
	query, _ := reader.ReadString('\n')
	return strings.TrimSpace(query)
}

func main() {
	fmt.Println("🏴‍☠️ Scanning for a working piracy site mirror...")

	result, err := sites.FindFirstWorkingMirror()
	if err != nil {
		log.Fatalf("❌ No working mirror found. Error: %v", err)
	}

	fmt.Println("✅ Working Mirror Found!")
	fmt.Printf("🔸 Site    : %s\n", result.SiteName)
	fmt.Printf("🔗 Mirror  : %s\n", result.Mirror)

	// tempQuery := getUserInput()
	// tempQuery := "interstellar"
	tempQuery := "brooklyn nine nine s01"

	torrentParser, err := parser.NewParser(result.SiteName, result.Mirror)
	if err != nil {
		log.Fatalf("Could not create parser: %v", err)
	}

	torrents, err := torrentParser.Search(tempQuery)
	if err != nil {
		log.Fatal(err)
	}

	enrichedTorrents := torrentParser.EnrichTorrents(torrents)
	for idx, torrent := range enrichedTorrents {
		display.PrintTorrentDebug(torrent, idx+1)
	}
}
