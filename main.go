package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/sanjaix21/krakeneye/internal/display"
	"github.com/sanjaix21/krakeneye/internal/parser"
	"github.com/sanjaix21/krakeneye/internal/ranker"
	"github.com/sanjaix21/krakeneye/internal/sites"
)

func getUserInput(query string) string {
	reader := bufio.NewReader(os.Stdin)

	switch {
	case strings.Contains(query, "search"):
		fmt.Printf("🔍 Enter search query (e.g. interstellar 2014): ")
		query, _ := reader.ReadString('\n')
		return strings.TrimSpace(query)
	case strings.Contains(query, "option"):
		fmt.Printf("➡️ Enter id to get magnet link (e.g. 2): ")
		option, _ := reader.ReadString('\n')
		return strings.TrimSpace(option)
	case strings.Contains(query, "new"):
		fmt.Printf("❓ Want to make a new search (y/n): ")
		option, _ := reader.ReadString('\n')
		return strings.TrimSpace(option)
	default:
		return ""
	}
}

func searchMedia(torrentParser parser.TorrentParser) ([]parser.TorrentFile, error) {
	tempQuery := getUserInput("search")

	torrents, err := torrentParser.Search(tempQuery)
	if err != nil {
		if err.Error() == "none" {
			fmt.Printf("No torrents found. Try checking name of the movie/tv\n")
		}
		log.Fatal(err)
	}

	return torrents, nil
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

	torrentParser, err := parser.NewParser(result.SiteName, result.Mirror)
	if err != nil {
		log.Fatalf("Could not create parser: %v", err)
	}

	for {

		torrents, err := searchMedia(torrentParser)
		if err != nil {
			log.Fatalf("Failed to search for media")
		}

		enrichedTorrents := torrentParser.EnrichTorrents(torrents)
		rankerFunc := &ranker.RankTorrent{}

		var torrentPointers []*parser.TorrentFile
		for i := range enrichedTorrents {
			enrichedTorrents[i].Score = rankerFunc.RankTorrentFile(enrichedTorrents[i])
			torrentPointers = append(torrentPointers, &enrichedTorrents[i])
		}

		displayOutput := display.NewDisplayManager(torrentPointers)
		displayOutput.ListTorrents()
		option, err := strconv.Atoi(getUserInput("option"))
		fmt.Printf("%T\n", option)
		if err != nil {
			fmt.Printf("Unable to convert string to int\n")
			return
		}
		fmt.Println("⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘")
		fmt.Println(torrentPointers[option].MagnetLink)
		fmt.Println("⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘")

		newSearch := getUserInput("new")
		if newSearch != "y" {
			fmt.Println("Thanks for using KrakenEye 🐉. May the Force be with you 🌠")
			break
		}
	}

	enrichedTorrents := torrentParser.EnrichTorrents(torrents)
	rankerFunc := &ranker.RankTorrent{}
	var torrentPointers []*parser.TorrentFile
	for i := range enrichedTorrents {
		enrichedTorrents[i].Score = rankerFunc.RankTorrentFile(enrichedTorrents[i])
		torrentPointers = append(torrentPointers, &enrichedTorrents[i])
	}

	displayOutput := display.NewDisplayManager(torrentPointers)
	displayOutput.ListTorrents()
	option, err := strconv.Atoi(getUserInput("option"))
	if err != nil {
		fmt.Printf("Unable to convert string to int\n")
		return
	}
	fmt.Println("⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘")
	fmt.Println(torrentPointers[option].MagnetLink)
	fmt.Println("⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘⫘")
}
