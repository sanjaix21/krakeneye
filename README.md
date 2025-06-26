# 🧿 KrakenEye

## 🌐 Live Demo

Try out the live version of the Web UI here:  
🔗 [https://krakeneye-web.onrender.com/](https://krakeneye-web.onrender.com/)

---

**KrakenEye** is a command-line and web-based torrent aggregator tool. It scrapes torrent data from public trackers, beginning with **RARBG**, and is designed for performance, reliability, and future extensibility.

> ⚠️ **Disclaimer:** This project is intended for educational and research purposes only. Use responsibly and ensure compliance with all applicable laws in your jurisdiction.

---

## ⚙️ Version

**Current:** `v0.1.0-beta`  
This is a beta release with both CLI and Web UI support, RARBG parsing, and magnet link extraction.

---

## 🎯 Features

- Connects to a working RARBG mirror
- Accepts search queries (e.g., `interstellar 2014`)
- Parses and displays:
  - Torrent title
  - Size
  - Category
  - Seeders / Leechers
  - Uploader
  - Upload date
  - Magnet link

---

## 📥 How to Use

### CLI Mode

1. Build the project:
   ```bash
   go build -o krakeneye
   ```

2. Run the binary:
   ```bash
   ./krakeneye
   ```

3. Enter your search term when prompted:
   ```bash
   🔍 Enter search query (e.g. interstellar 2014): the matrix 1999
   ```

### Web UI Mode

1. Build the project:
   ```bash
   go build -o krakeneye
   ```

2. Run with the `--web` flag:
   ```bash
   ./krakeneye --web
   ```

3. Open your browser and go to: [http://localhost:8787](http://localhost:8787)

---

## 🤝 Contributions

Contributions are welcome. If you'd like to contribute, please open an issue or submit a pull request. Ensure your code is well-documented and tested. For major changes, please open a discussion first.

---

## ⚓ License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.
