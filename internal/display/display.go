package display

import (
	"fmt"

	"github.com/sanjaix21/krakeneye/internal/parser"
)

func PrintTorrentDebug(torrent parser.TorrentFile, index int) {
	fmt.Printf("🧿 Torrent #%d\n", index)
	fmt.Println("────────────────────────────")
	fmt.Printf("🎬 Name       : %s\n", torrent.Name)
	fmt.Printf("📁 Size       : %s (%.2f GB)\n", torrent.SizeRaw, torrent.Size)
	fmt.Printf("🔗 Magnet     : %s\n", torrent.MagnetLink)
	fmt.Printf("📂 Category   : %s\n", torrent.Category)
	fmt.Printf("📅 Uploaded   : %s\n", torrent.UploadDate)
	fmt.Printf("🚀 Seeders    : %d\n", torrent.Seeders)
	fmt.Printf("🧲 Leechers   : %d\n", torrent.Leechers)
	fmt.Printf("📤 Uploader   : %s (Trusted: %t)\n", torrent.Uploader, torrent.Trusted)
	fmt.Printf("🌐 Language   : %s\n", torrent.Language)
	fmt.Printf("📉 Downloads  : %d\n", torrent.Downloads)
	fmt.Printf("🎞️ Source     : %s\n", torrent.Source)
	fmt.Printf("🖥️ Resolution : %s\n", torrent.Resolution)
	fmt.Printf("🎧 Audio      : %s\n", torrent.AudioCodec)
	fmt.Printf("📼 Video      : %s\n", torrent.VideoCodec)
	fmt.Printf("📦 Container  : %s\n", torrent.Container)
	fmt.Printf("🌈 Bit Depth  : %s\n", torrent.BitDepth)
	fmt.Printf("📃 MetaInfo   : %s\n", torrent.MetaInfo)
	fmt.Println()
}
