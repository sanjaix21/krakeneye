package display

import (
	"fmt"
	"sort"

	"github.com/sanjaix21/krakeneye/internal/parser"
	"github.com/sanjaix21/krakeneye/internal/ranker"
)

type DebugDisplay struct {
	ranker *ranker.RankTorrent
}

type DisplayManager struct {
	torrents []*parser.TorrentFile
}

func NewDebugDisplay() *DebugDisplay {
	return &DebugDisplay{
		ranker: &ranker.RankTorrent{},
	}
}

func NewDisplayManager(torrents []*parser.TorrentFile) *DisplayManager {
	return &DisplayManager{
		torrents: torrents,
	}
}

func (dd *DebugDisplay) PrintTorrentDebug(torrent parser.TorrentFile, index int) {
	fmt.Printf("🧿 Torrent #%d\n", index)
	fmt.Println("────────────────────────────")
	fmt.Printf("🎬 Name       : %s\n", torrent.Name)
	fmt.Printf("🔗 Link       : %s\n", torrent.Href)
	fmt.Printf("📁 Size       : %s (%.2f GB)\n", torrent.SizeRaw, torrent.Size)
	fmt.Printf("🧲 Magnet     : %s\n", torrent.MagnetLink)
	fmt.Printf("📊 Category   : %s\n", torrent.Category)
	fmt.Printf("📅 Uploaded   : %s\n", torrent.UploadDate)
	fmt.Printf("🚀 Seeders    : %d\n", torrent.Seeders)
	fmt.Printf("🩸 Leechers   : %d\n", torrent.Leechers)
	fmt.Printf("📤 Uploader   : %s (Trusted: %t)\n", torrent.Uploader, torrent.Trusted)
	fmt.Printf("🌐 Language   : %s\n", torrent.Language)
	fmt.Printf("⏬ Downloads  : %d\n", torrent.Downloads)
	fmt.Printf("🎞️ Source     : %s\n", torrent.Source)
	fmt.Printf("🖥️ Resolution : %s\n", torrent.Resolution)
	fmt.Printf("🎧 Audio      : %s\n", torrent.AudioCodec)
	fmt.Printf("📼 Video      : %s\n", torrent.VideoCodec)
	fmt.Printf("📦 Container  : %s\n", torrent.Container)
	fmt.Printf("🌈 Bit Depth  : %s\n", torrent.BitDepth)
	fmt.Printf("📃 MetaInfo   : %s\n", torrent.MetaInfo)
	fmt.Println()
}

func (dd *DebugDisplay) PrintSizeScoreDebug(torrent parser.TorrentFile) {
	score := dd.ranker.RankSize(torrent)
	fmt.Printf("📦 Torrent Debug Report\n")
	fmt.Printf("🔤 Name:        %s\n", torrent.Name)
	fmt.Printf("📂 Category:    %s\n", torrent.Category)
	fmt.Printf("🖥 Resolution:  %s\n", torrent.Resolution)
	fmt.Printf("🎞 Source:      %s\n", torrent.Source)
	fmt.Printf("💾 Size:        %.2f GB\n", torrent.Size)
	fmt.Printf("🏅 Size Score:  %.2f / 20\n", score)
	fmt.Println("⚓------------------------------")
	fmt.Println()
}

func (dd *DebugDisplay) PrintSeedScoreDebug(torrent parser.TorrentFile) {
	score := dd.ranker.RankSeeds(torrent)
	fmt.Printf("📦 Torrent Debug Report\n")
	fmt.Printf("🔤 Name:        %s\n", torrent.Name)
	fmt.Printf("🚀 Seeders    : %d\n", torrent.Seeders)
	fmt.Printf("🩸 Leechers   : %d\n", torrent.Leechers)
	fmt.Printf("🏅 Seed Score:  %.2f / 20\n", score)
	fmt.Println("⚓------------------------------")
	fmt.Println()
}

func (dd *DebugDisplay) PrintTorrentScoreDebug(torrent parser.TorrentFile) {
	torrentScore := dd.ranker.RankTorrentFile(torrent)

	fmt.Printf("🏴‍☠️ KrakenEye Torrent Score Debug\n")
	fmt.Printf("🔤 Name:             %s\n", torrent.Name)
	fmt.Printf("───────────────────────────────────────\n")
	fmt.Printf("📏 Size Score        : %.2f / 25\n", dd.ranker.SizeScore)
	fmt.Printf("🌱 Seeder Score      : %.2f / 30\n", dd.ranker.SeedScore)
	fmt.Printf("🖥 Resolution Score  : %.2f / 20\n", dd.ranker.ResolutionScore)
	fmt.Printf("🎞 Source Score      : %.2f / 15\n", dd.ranker.SourceScore)
	fmt.Printf("🎧 Codecs Score      : %.2f / 7\n", dd.ranker.CodecsScore)
	fmt.Printf("🧑‍🚀 Uploader Score   : %.2f / 3\n", dd.ranker.UploaderScore)
	fmt.Println("───────────────────────────────────────")
	fmt.Printf("🏁 TOTAL SCORE       : %.2f / 100\n", torrentScore)
	fmt.Println("⚓---------------------------------------")
	fmt.Println()
}

func (dm *DisplayManager) ListTorrents() {
	sort.Slice(dm.torrents, func(i, j int) bool {
		return dm.torrents[i].Score > dm.torrents[j].Score
	})

	fmt.Println("Ranked Torrent List:")
	fmt.Println(
		"---------------------------------------------------------------------------------------------",
	)
	fmt.Printf(
		"%-3s %-50s %-7s %-8s %-6s %-6s\n",
		"#",
		"Name",
		"Size(GB)",
		"Seeders",
		"Res",
		"Score",
	)
	fmt.Println(
		"---------------------------------------------------------------------------------------------",
	)

	// Display Each Torrent
	for i, torrent := range dm.torrents {
		fmt.Printf("%-3d %-50s %-7.2f %-8d %-6s %-6.2f\n",
			i+1,
			truncateString(torrent.Name, 50),
			torrent.Size, // Convert bytes to GB
			torrent.Seeders,
			torrent.Resolution,
			torrent.Score,
		)
	}
}

func truncateString(str string, maxLen int) string {
	if len(str) <= maxLen {
		return str
	}
	return str[:maxLen-3] + "..."
}
