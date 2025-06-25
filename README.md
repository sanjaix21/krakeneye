# 🧿 KrakenEye

**KrakenEye** is a command-line torrent aggregator tool built for pirates, by pirates. It scrapes torrent data from public trackers, starting with **RARBG**, with plans to expand into **1337x** and **KickassTorrent**.

> ⚠️ **Disclaimer:** This project is for educational and research purposes only. Use responsibly and comply with all applicable laws in your region.

---

## ⚙️ Version

**Current:** `v0.0.1-alpha`

This is an early prototype release meant to test the parser and search engine.

---

## 🎯 What It Does

- 🌐 Connects to a working RARBG mirror
- 🔍 Accepts a search query from the user (e.g., `interstellar 2014`)
- 📄 Parses and displays:
  - Torrent title
  - Size
  - Category
  - Seeders / Leechers
  - Uploader
  - Upload date
  - Magnet link (✅ early feature)

---

## 📥 How To Use

```bash
go run main.go
```

Then enter your search term when prompted:

```bash
🔍 Enter search query (e.g. interstellar 2014): interstellar 2014
```

The tool will:
1. Find a working RARBG mirror
2. Perform a search using your query
3. Parse and display magnet links and torrent info

---

## 🗺️ Roadmap

### 🔧 Immediate Tasks
- Finish complete RARBG parser (multiple pages, mirror failover)
- Implement torrent ranking logic based on seeders/quality
- Add CLI help menu and usage flags

### 🏴‍☠️ Planned Features
- Support for **1337x**, **KickassTorrent**, and more
- Intelligent ranking system
- Proxy support for geo-restricted users
- UI dashboard for web-based access
- .torrent file support
- ML-based ranking and recommendation engine (stretch goal)

---

## 🧠 Example

```
🔍 Enter search query (e.g. interstellar 2014): the matrix 1999
```

Result:
```
🎬 Title: The.Matrix.1999.1080p.BluRay.x264
📁 Size: 1.9 GB
🔢 Seeders: 1245
🧷 Magnet: magnet:?xt=urn:btih:...
```

---

## 📁 File Structure

```plaintext
krakeneye/
├── main.go
├── go.mod
├── go.sum
└── internal/
    ├── parser/
    │   ├── factory.go          # Determines which parser to use
    │   ├── parser.go           # Common parser interface
    │   └── rarbg_parser.go     # Parses RARBG listings
    ├── ranker/
    │   └── ranker.go           # Torrent ranking logic (TODO)
    └── sites/
        ├── get_working_mirror.go  # Finds a working mirror
        └── piracy_sites.go        # Defines supported sites
```

---

## 🤝 Contributions

All hands on deck! This ship is still being built — contributions, feature requests, and suggestions are welcome.

---

## ⚓ License

MIT License (See [LICENSE](./LICENSE) file)

---

**Raise the Kraken.** 🏴‍☠️  
Built with Go and the spirit of piracy.
