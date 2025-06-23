function searchTorrents() {
  const query = document.getElementById("searchInput").value.trim();
  const loading = document.getElementById("loading");
  const results = document.getElementById("results");

  if (!query) return;

  results.innerHTML = "";
  loading.classList.remove("hidden");

  fetch(`/search?q=${encodeURIComponent(query)}`)
    .then(res => res.json())
    .then(data => {
      loading.classList.add("hidden");

      if (!data.length) {
        results.innerHTML = "<p class='text-center text-red-500'>⚠️ No results found</p>";
        return;
      }

      results.innerHTML = data.map(t => `
        <div class="bg-gradient-to-br from-gray-900 to-red-950 p-4 rounded-2xl shadow-lg border border-red-700 transition-transform hover:scale-105 duration-200 overflow-hidden">
        <h2 class="text-xl font-bold text-yellow-300 break-words mb-2">${t.Name}</h2>
        <div class="text-sm text-gray-300 space-y-1">
          <p>🎬 <span class="text-white">Size:</span> ${t.Size || "?"}</p>
          <p>📺 <span class="text-white">Resolution:</span> ${t.Resolution || "Unknown"}</p>
          <p>🌱 <span class="text-white">Seeders:</span> ${t.Seeders || "?"}</p>
          <p>🧭 <span class="text-white">Source:</span> ${t.SiteName || "Unknown"}</p>
          <p>🧲 <button onclick='copyMagnet("${t.MagnetLink}")' class="mt-1 bg-red-600 hover:bg-red-500 px-3 py-1 rounded-full text-white font-bold">Magnet Link</button></p>
          <p class="text-right text-xs text-red-400 italic">🐉 KrakenEye Score: ${t.Score?.toFixed(2)}</p>
        </div>
      </div>
  `).join("");
    })
    .catch(() => {
      loading.classList.add("hidden");
      results.innerHTML = "<p class='text-center text-red-500'>⚠️ Error fetching results</p>";
    });
}

function copyMagnet(link) {
  navigator.clipboard.writeText(link);
  alert("🧲 Magnet link copied!");
}

// Add enter key support
document.getElementById("searchInput").addEventListener("keydown", function (e) {
  if (e.key === "Enter") {
    searchTorrents();
  }
});

