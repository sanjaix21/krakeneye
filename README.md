# 🧿 KrakenEye

**KrakenEye** is a command-line and web-based torrent discovery tool built for pirates, by pirates. It scrapes torrent data from public trackers, currently starting with **RARBG**, with plans to expand into **1337x**, **KickassTorrent**, and more.

> ⚠️ **Disclaimer:** This project is for educational and research purposes only. Use responsibly and comply with all applicable laws in your region.

---

## ⚙️ Version

**Current:** `v0.1.0-alpha`

- CLI + WebUI functional  
- Telegram bot in development  
- Early parser, ranker, and scoring in place

---

## 🎯 What It Does

- 🌍 Connects to a working RARBG mirror
- 🔎 Accepts search queries from user or Web UI
- 📄 Parses and ranks torrent results
- 🧠 Displays:
  - Title
  - Size
  - Resolution
  - Source
  - Seeders
  - KrakenEye Score™
  - Magnet link (copyable via UI)

---

## 🚀 How To Use

### 🖥️ CLI Mode

```bash
go run main.go
```

Then follow prompts:

```bash
🔍 Enter search query (e.g. interstellar 2014): the matrix 1999
```

---

### 🌐 Web UI Mode

```bash
go run main.go --web
```

- Default port: `8787` (auto-increments if busy)
- Open browser at: [http://localhost:8787](http://localhost:8787)

A pirate-themed interface with search bar, elegant results, copy magnet buttons, and KrakenEye scoring.

---

## 📁 Folder Structure

```plaintext
krakeneye/
├── main.go
├── go.mod
├── Dockerfile
├── internal/
│   ├── parser/              # Handles different site parsers (RARBG, future: 1337x, etc.)
│   ├── ranker/              # Torrent ranking system (file size, seeders, source, etc.)
│   ├── display/             # CLI output logic
│   ├── sites/               # Mirror detection & piracy site definitions
│   └── webui/               # Web UI server & static frontend
│       ├── static/          # HTML, CSS, JS, images
│       └── web-ui.go        # Starts the WebUI server
└── internal/telegrambot/    # (WIP) Telegram bot interface
```

---

## 🗺️ Roadmap

### 🔧 Immediate Tasks

- 🛠 Improve scoring logic (codecs, uploader trust)
- 🕵️ Add fallback to 1337x when RARBG fails
- 📱 Polish Telegram bot commands (`/search`, `/get`)

### 🏴‍☠️ Planned Features

- ✅ Web UI (done)
- 🧠 Intelligent ranking engine
- 🧭 Proxy + region unlock
- 💾 .torrent file download support
- 🤖 ML-powered recommendation engine (stretch goal)
- 🐳 Docker & DevOps CI/CD pipelines

---

## 🤝 Contributions

All hands on deck! This ship is still being built — contributions, feature requests, and suggestions are welcome.

---

## ⚓ License

MIT License (See [LICENSE](./LICENSE) file)

---

**Raise the Kraken.** 🏴‍☠️  
Built with Go and the spirit of piracy by [sanjaix21](https://github.com/sanjaix21)
