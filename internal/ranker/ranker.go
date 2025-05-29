package ranker

import "github.com/sanjaix21/krakeneye/internal/parser"

type RankTorrent struct{}

func (rt *RankTorrent) RankTorrentFile(torrent parser.TorrentFile) float64 {
	score := 0.0

	score += rt.RankSize(torrent)
	return score
}

func (rt *RankTorrent) RankSize(torrent parser.TorrentFile) float64 {
	category := torrent.Category
	sizeScore := 0.0
	switch category {
	case "Movies":
		sizeScore = rankMovies(torrent.Size, torrent.Resolution, torrent.Source)

	case "TV":
		sizeScore = rankTV(torrent.Size, torrent.Resolution, torrent.Source)

	default:
		sizeScore = 0.0
	}

	return sizeScore
}

func rankMovies(size float64, resolution string, source string) float64 {
	if resolution == "480P" {
	}
}
