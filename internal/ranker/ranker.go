// TODO: add 1440P/2K
// include download as a metric for ranking
// refine ranking systems
package ranker

import (
	"math"
	"sanjaix21/krakeneye/internal/parser"
	"strings"
)

/*
COMPREHENSIVE TORRENT RANKING SYSTEM - 100 POINTS TOTAL

Priority Distribution (based on user requirements):
1. Seeders (30 points) - Speed of download is crucial
2. Size (25 points) - Optimal file size for quality
3. Resolution (20 points) - Video quality matters
4. Source (15 points) - Source quality (IMAX, Blu-ray, etc.)
5. Codecs (7 points) - Video/Audio codec quality
6. Uploader Trust (3 points) - Trusted uploaders bonus

Total: 100 points
*/

type RankTorrent struct {
	// Storing this mostly for debug purpose
	SizeScore       float64
	SeedScore       float64
	ResolutionScore float64
	SourceScore     float64
	CodecsScore     float64
	UploaderScore   float64
	TorrentScore    float64
}

func (rt *RankTorrent) RankTorrentFile(torrent parser.TorrentFile) float64 {
	torrentScore := 0.0

	torrentScore += rt.RankSize(torrent)
	torrentScore += rt.RankSeeds(torrent)
	torrentScore += rt.RankResolution(torrent)
	torrentScore += rt.RankSource(torrent)
	torrentScore += rt.RankCodecs(torrent)
	torrentScore += rt.RankUploader(torrent)

	rt.TorrentScore = torrentScore
	torrent.Score = torrentScore
	return torrentScore
}

// Seeds Ranking (30 Points Max)
func (rt *RankTorrent) RankSeeds(torrent parser.TorrentFile) float64 {
	seeders := torrent.Seeders
	leechers := torrent.Leechers

	var seedScore float64

	switch {
	case seeders >= 100:
		// excellent availabilty - very fast download speed
		seedScore = 30.0

	case seeders >= 50:
		// very good availabilty - fast download
		progress := float64(seeders-50) / 50.0
		seedScore = 25.0 + (progress * 5.0) // 25-20 points speed

	case seeders >= 20:
		// Good availabilty - decent download
		progress := float64(seeders-20) / 20.0
		seedScore = 20.0 + (progress * 5.0) // 20 - 25 points speed

	case seeders >= 10:
		// moderate availabilty - slower download speed but not too bad
		progress := float64(seeders-10) / 10.0
		seedScore = 12.0 + (progress * 8.0) // 12 - 20 points

	case seeders >= 5:
		// Low availabilty - very slow download speed
		progress := float64(seeders-5) / 5.0
		seedScore = 4.0 + (progress * 8.0) // 4 - 12 points

	case seeders > 0:
		// very low availabilty - very slow download speed
		progress := float64(seeders-5) / 5.0
		seedScore = 1.0 + (progress * 3.0) // 1 - 4 points

	default:
		return 0.0 // Dead Torrent
	}

	// seeders - leechers ratio
	// better ratio = better speed

	if seeders > 0 {
		seedRatio := float64(seeders) / math.Max(float64(leechers), 1.0)

		switch {
		case seedRatio >= 3.0:
			// Excellent Ratio
			seedScore *= 1.15

		case seedRatio >= 2.0:
			// Great Ratio
			seedScore *= 1.10

		case seedRatio >= 1.0:
			// Good Ratio
			seedScore *= 1.05

		case seedRatio >= 0.5:
			// Balanced ratio
			// No change
			seedScore *= 1.0

		case seedRatio >= 0.2:
			// Poor Ratio
			// congested traffic
			seedScore *= 0.95

		default:
			// very poor ratio - heavly congested traffic
			seedScore *= 0.90
		}

	}
	rt.SeedScore = seedScore
	return math.Min(seedScore, 30.0)
}

func (rt *RankTorrent) RankSize(torrent parser.TorrentFile) float64 {
	categoryLower := torrent.Category

	sizeScore := 0.0
	switch {
	case strings.Contains(categoryLower, "movies"):
		sizeScore = rt.rankMovieSize(torrent.Size, torrent.Resolution, torrent.Source)

	case strings.Contains(categoryLower, "tv"):
		sizeScore = rt.rankTvSize(torrent.Size, torrent.Resolution, torrent.Source)

	default:
		sizeScore = rt.rankGenericSize(torrent.Size, torrent.Resolution)

	}

	rt.SizeScore = sizeScore
	return sizeScore
}

// Resolution Ranking (20 Points Max)
func (rt *RankTorrent) RankResolution(torrent parser.TorrentFile) float64 {
	resolution := strings.ToUpper(torrent.Resolution)

	var resolutionScore float64
	switch resolution {
	case "2160P":
		resolutionScore = 15.0
	case "1080P":
		resolutionScore = 20.0 // 1080P is what most people prefer so it has more points
	case "720P":
		resolutionScore = 10.0
	case "480P":
		resolutionScore = 5.0
	default:
		resolutionScore = 8.0 // Unknown Resolution
	}

	rt.ResolutionScore = resolutionScore
	return resolutionScore
}

// Source Ranking (15 Points Max)
func (rt *RankTorrent) RankSource(torrent parser.TorrentFile) float64 {
	var sourceScore float64
	sourceUpper := strings.ToUpper(torrent.Source)

	switch sourceUpper {
	case "IMAX":
		sourceScore = 12.0
	case "BLURAY":
		sourceScore = 15.0 // most prefer bluray so more score
	case "WEB":
		sourceScore = 10.0
	case "CAM":
		sourceScore = 5.0
	default:
		sourceScore = 10.0 // for unknown mostly be web
	}

	rt.SourceScore = sourceScore
	return sourceScore
}

// Codec Ranking (7 POINTS MAX)
func (rt *RankTorrent) RankCodecs(torrent parser.TorrentFile) float64 {
	videoCodec := torrent.VideoCodec
	audioCodec := torrent.AudioCodec
	bitDepth := torrent.BitDepth

	// Video codec scoring (5 points max)
	var videoScore float64
	switch videoCodec {
	case "AV1":
		// AV1 - newest, most efficient codec, best compression
		videoScore = 5.0

	case "HEVC/x265":
		// HEVC/x265 - modern, efficient codec, great quality
		videoScore = 4.5

	case "x264":
		// x264 - widely compatible, good quality, mature codec
		videoScore = 3.5

	case "Unknown":
		// Unknown codec - could be anything, neutral score
		videoScore = 2.0

	default:
		// Fallback for any other detected codecs
		videoScore = 2.5
	}

	// Audio codec scoring (1.5 points max)
	var audioScore float64
	switch audioCodec {
	case "Dolby Atmos":
		// Dolby Atmos - premium object-based audio, best quality
		audioScore = 1.5

	case "DTS-HD/TrueHD":
		// DTS-HD/TrueHD - lossless audio codecs, excellent quality
		audioScore = 1.4

	case "DTS":
		// DTS - good quality lossy codec, better than AC3
		audioScore = 1.2

	case "EAC3":
		// Enhanced AC3 (DD+) - improved version of AC3
		audioScore = 1.1

	case "AC3":
		// AC3/Dolby Digital - standard surround sound
		audioScore = 1.0

	case "AAC":
		// AAC - modern, efficient stereo/multi-channel codec
		audioScore = 0.9

	case "OPUS":
		// OPUS - very efficient modern codec, excellent for streaming
		audioScore = 0.8

	case "MP3":
		// MP3 - widely compatible but older, limited to stereo
		audioScore = 0.6

	case "Unknown":
		// Unknown audio codec - could be anything
		audioScore = 0.5

	default:
		// Fallback for any other detected audio codecs
		audioScore = 0.7
	}

	// Bit depth bonus (0.5 points max)
	var bitDepthScore float64
	switch bitDepth {
	case "10-bit":
		// 10-bit - HDR capable, better color gradients, future-proof
		bitDepthScore = 0.5

	case "8-bit":
		// 8-bit - standard bit depth, widely compatible
		bitDepthScore = 0.2

	default:
		// Unknown bit depth - assume standard
		bitDepthScore = 0.2
	}

	return videoScore + audioScore + bitDepthScore
}

func (rt *RankTorrent) RankUploader(torrent parser.TorrentFile) float64 {
	isTrusted := torrent.Trusted
	var uploaderScore float64
	if isTrusted {
		uploaderScore = 3.0
	} else {
		uploaderScore = 0.0
	}

	rt.UploaderScore = uploaderScore
	return uploaderScore
}

func (rt *RankTorrent) rankMovieSize(size float64, resolution string, source string) float64 {
	if size <= 0 {
		return 0
	}

	var sweetSpot, tolerance float64

	switch resolution {
	case "480P":
		sweetSpot = rt.getMovieSweetSpot480p(source)
		tolerance = 0.5 // +- 0.5GB

	case "720P":
		sweetSpot = rt.getMovieSweetSpot720p(source)
		tolerance = 1.0 // +- 1.0GB

	case "1080P":
		sweetSpot = rt.getMovieSweetSpot1080p(source)
		tolerance = 2.0 // +- 2.0GB

	case "2160P": // 4k
		sweetSpot = rt.getMovieSweetSpot2160p(source)
		tolerance = 6.0 // +- 6.0GB

	default:
		return rt.rankGenericSize(size, resolution)
	}

	return rt.calculateSizeScore(size, sweetSpot, tolerance, 20.0)
}

func (rt *RankTorrent) rankTvSize(size float64, resolution string, source string) float64 {
	if size <= 0 {
		return 0
	}

	var sweetSpot, tolerance float64

	switch resolution {
	case "480P":
		sweetSpot = rt.getTvSweetSpot480p(source)
		tolerance = 2.0

	case "720P":
		sweetSpot = rt.getTvSweetSpot720p(source)
		tolerance = 3.5

	case "1080P":
		sweetSpot = rt.getTvSweetSpot1080p(source)
		tolerance = 5.0

	case "2160P":
		sweetSpot = rt.getTvSweetSpot2160p(source)
		tolerance = 10.0

	default:
		return rt.rankGenericSize(size, resolution)
	}

	return rt.calculateSizeScore(size, sweetSpot, tolerance, 20.0)
}

func (rt *RankTorrent) calculateSizeScore(
	size float64,
	sweetSpot float64,
	tolerance float64,
	maxScore float64,
) float64 {
	if sweetSpot <= 0 {
		return 0
	}

	distance := math.Abs(size - sweetSpot)

	// within tolerance
	if distance <= tolerance {
		score := maxScore * (1 - (distance / tolerance * 0.1)) // 90-100% score when within tolerance
		return math.Min(score, maxScore)
	}

	// Outside tolerance
	excessDistance := distance - tolerance
	decayFactor := math.Exp(-excessDistance / sweetSpot) // Normalize by sweetspot

	return maxScore * 0.9 * decayFactor // upto 90% score when outside tolerance
}

func (rt *RankTorrent) getMovieSweetSpot480p(source string) float64 {
	switch {
	case strings.Contains(source, "BLURAY"):
		return 2.0
	case strings.Contains(source, "WEB"):
		return 1.0
	case strings.Contains(source, "CAM"):
		return 0.7
	default:
		return 1.0
	}
}

func (rt *RankTorrent) getMovieSweetSpot720p(source string) float64 {
	switch {
	case strings.Contains(source, "IMAX"):
		return 8.0
	case strings.Contains(source, "BLURAY"):
		return 6.0
	case strings.Contains(source, "WEB"):
		return 3.5
	case strings.Contains(source, "CAM"):
		return 2.0
	default:
		return 4.5
	}
}

func (rt *RankTorrent) getMovieSweetSpot1080p(source string) float64 {
	switch {
	case strings.Contains(source, "IMAX"):
		return 12.0
	case strings.Contains(source, "BLURAY"):
		return 10.0
	case strings.Contains(source, "WEB"):
		return 7.0
	case strings.Contains(source, "CAM"):
		return 3.0
	default:
		return 4.0
	}
}

func (rt *RankTorrent) getMovieSweetSpot2160p(source string) float64 {
	switch {
	case strings.Contains(source, "IMAX"):
		return 60.0
	case strings.Contains(source, "BLURAY"):
		return 50.0
	case strings.Contains(source, "WEB"):
		return 25.0
	default:
		return 35.0
	}
}

func (rt *RankTorrent) getTvSweetSpot480p(source string) float64 {
	switch {
	case strings.Contains(source, "BLURAY"):
		return 5.0
	case strings.Contains(source, "WEB"):
		return 3.5
	case strings.Contains(source, "CAM"):
		return 2.0
	default:
		return 2.5
	}
}

func (rt *RankTorrent) getTvSweetSpot720p(source string) float64 {
	switch {
	case strings.Contains(source, "IMAX"):
		return 12.0
	case strings.Contains(source, "BLURAY"):
		return 10.0
	case strings.Contains(source, "WEB"):
		return 8.0
	case strings.Contains(source, "CAM"):
		return 3.5
	default:
		return 4.0
	}
}

func (rt *RankTorrent) getTvSweetSpot1080p(source string) float64 {
	switch {
	case strings.Contains(source, "IMAX"):
		return 55.0
	case strings.Contains(source, "BLURAY"):
		return 30.0
	case strings.Contains(source, "WEB"):
		return 20.0
	case strings.Contains(source, "CAM"):
		return 7.5
	default:
		return 10.0
	}
}

func (rt *RankTorrent) getTvSweetSpot2160p(source string) float64 {
	switch {
	case strings.Contains(source, "IMAX"):
		return 100.0
	case strings.Contains(source, "BLURAY"):
		return 60.0
	case strings.Contains(source, "WEB"):
		return 30.0
	default:
		return 30.0
	}
}

func (rt *RankTorrent) rankGenericSize(size float64, resolution string) float64 {
	if size <= 0 {
		return 0
	}

	var idealSize, tolerance float64
	switch resolution {
	case "480P":
		idealSize = 1.0
	case "720P":
		idealSize = 3.5
	case "1080P":
		idealSize = 10.0
	case "2160P":
		idealSize = 35.0
	default:
		return 10.0 // for unknown/NONE resolution
	}
	tolerance = idealSize * 0.5
	return rt.calculateSizeScore(size, idealSize, tolerance, 20)
}
