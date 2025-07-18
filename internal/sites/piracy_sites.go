package sites

type Site struct {
	Name     string
	Category string // like torrents/streaming (useful for future updates)
	Primary  string
	Weight   int // higher number = higher prority
	Mirrors  []string
}

// sites are given in prority order, add new sites in the correct place based on your priority
var PiracySites = []Site{
	{
		Name:     "rarbg",
		Category: "torrents",
		Primary:  "rarbg.to",
		Weight:   3,
		Mirrors: []string{
			"https://en.rarbg.gg/",
			"https://rargb.to/",
			"https://www.rarbgproxy.to/",
			"https://www.proxyrarbg.to/",
			"https://www2.rarbggo.to/",
		},
	},

	// {
	// 	Name:     "https://kickasstorrents.cr/",
	// 	Category: "torrents",
	// 	Primary:  "https://kickasstorrents.cr/",
	// 	Weight:   2,
	// 	Mirrors: []string{
	// 		"https://kickasstorrents.cr/",
	// 		"https://kickasstorrents.cr/",
	// 		"https://kickass.sx/",
	// 		"https://katcr.to/",
	// 		"https://katcr.to/",
	// 		"https://kickasstorrent.cr/",
	// 	},
	// },
	// {
	// 	Name:     "1337x",
	// 	Category: "torrents",
	// 	Primary:  "https://1337x.to/",
	// 	Weight:   1,
	// 	Mirrors: []string{
	// 		"https://1337x.to/",
	// 		"https://1337x.st/",
	// 		"https://x1337x.cc/",
	// 		"https://x1337x.ws/",
	// 		"https://x1337x.eu/",
	// 		"https://x1337x.se/",
	// 	},
	// },
}
