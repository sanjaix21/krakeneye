package webui

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/sanjaix21/krakeneye/internal/parser"
	"github.com/sanjaix21/krakeneye/internal/ranker"
	"github.com/sanjaix21/krakeneye/internal/sites"
)

func StartServer(port int) {
	fmt.Printf("🕸️  Launching KrakenEye Retro WebUI on http://localhost:%d\n", port)

	result, err := sites.FindFirstWorkingMirror()
	if err != nil {
		log.Fatalf("❌ No working mirror found. Error: %v", err)
	}

	fmt.Printf("✅ Connected to: %s (%s)\n", result.SiteName, result.Mirror)

	torrentParser, err := parser.NewParser(result.SiteName, result.Mirror)
	if err != nil {
		log.Fatalf("❌ Could not create parser: %v", err)
	}

	// Serve static files with proper MIME types
	fs := http.FileServer(http.Dir("internal/webui/static"))
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set proper MIME types for static files
		switch {
		case r.URL.Path == "/" || r.URL.Path == "/index.html":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		case r.URL.Path == "/styles.css":
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		case r.URL.Path == "/app.js":
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		}
		
		// Add security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		
		fs.ServeHTTP(w, r)
	}))

	// Search endpoint with enhanced error handling
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for development
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		query := r.URL.Query().Get("q")
		if query == "" {
			http.Error(w, `{"error": "Search query is required"}`, http.StatusBadRequest)
			return
		}

		fmt.Printf("🔍 Searching for: %s\n", query)

		// Search torrents
		torrents, err := torrentParser.Search(query)
		if err != nil {
			fmt.Printf("❌ Search failed: %v\n", err)
			if err.Error() == "none" {
				// Return empty array instead of error for no results
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode([]parser.TorrentFile{})
				return
			}
			http.Error(w, `{"error": "Search failed"}`, http.StatusInternalServerError)
			return
		}

		fmt.Printf("📦 Found %d torrents, enriching data...\n", len(torrents))

		// Enrich torrents with detailed information
		enrichedTorrents := torrentParser.EnrichTorrents(torrents)
		
		// Rank and score torrents
		rankerFunc := &ranker.RankTorrent{}
		var enrichedPtrs []*parser.TorrentFile
		
		for i := range enrichedTorrents {
			enrichedTorrents[i].Score = rankerFunc.RankTorrentFile(enrichedTorrents[i])
			enrichedTorrents[i].SiteName = result.SiteName
			enrichedPtrs = append(enrichedPtrs, &enrichedTorrents[i])
		}

		// Sort by score (highest first)
		sort.Slice(enrichedPtrs, func(i, j int) bool {
			return enrichedPtrs[i].Score > enrichedPtrs[j].Score
		})

		fmt.Printf("🏆 Ranked %d torrents, sending response\n", len(enrichedPtrs))

		// Return JSON response
		if err := json.NewEncoder(w).Encode(enrichedPtrs); err != nil {
			fmt.Printf("❌ JSON encoding failed: %v\n", err)
			http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
			return
		}
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "healthy",
			"site":   result.SiteName,
			"mirror": result.Mirror,
		})
	})

	fmt.Printf("🚀 Server starting on port %d...\n", port)
	fmt.Printf("🌐 Open your browser to: http://localhost:%d\n", port)
	fmt.Printf("🏴‍☠️ Ready to hunt torrents!\n\n")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}