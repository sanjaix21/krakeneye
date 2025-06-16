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
			"https://www2.rarbggo.to/",
			"https://www.rarbgproxy.to/",
			"https://www.proxyrarbg.to/",
		},
	},
}
