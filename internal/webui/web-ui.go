package webui

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sanjaix21/krakeneye/internal/parser"
	"github.com/sanjaix21/krakeneye/internal/ranker"
	"github.com/sanjaix21/krakeneye/internal/sites"
)

func StartServer(port int) {
	fmt.Printf("🕸️  Launching KrakenEye WebUI on http://localhost:%d\n", port)

	result, err := sites.FindFirstWorkingMirror()
	if err != nil {
		log.Fatalf("No working mirror found. Error: %v", err)
	}

	torrentParser, err := parser.NewParser(result.SiteName, result.Mirror)
	if err != nil {
		log.Fatalf("Could not create parser: %v", err)
	}

	// Serve static HTML + JS
	http.Handle("/", http.FileServer(http.Dir("internal/webui/static")))

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		torrents, err := torrentParser.Search(query)
		if err != nil {
			http.Error(w, "Search failed", http.StatusInternalServerError)
			return
		}

		enriched := torrentParser.EnrichTorrents(torrents)
		rankerFunc := &ranker.RankTorrent{}
		var enrichedPtrs []*parser.TorrentFile
		for i := range enriched {
			enriched[i].Score = rankerFunc.RankTorrentFile(enriched[i])
			enriched[i].SiteName = result.SiteName
			enrichedPtrs = append(enrichedPtrs, &enriched[i])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(enrichedPtrs)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
