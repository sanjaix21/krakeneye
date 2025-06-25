// KrakenEye Retro App
class KrakenEye {
  constructor() {
    this.searchInput = document.getElementById('searchInput');
    this.statusSection = document.getElementById('statusSection');
    this.statusText = document.getElementById('statusText');
    this.progressBar = document.querySelector('.progress-fill');
    this.results = document.getElementById('results');
    
    this.initializeEventListeners();
    this.createMatrixEffect();
  }

  initializeEventListeners() {
    // Enter key support
    this.searchInput.addEventListener('keydown', (e) => {
      if (e.key === 'Enter') {
        this.searchTorrents();
      }
    });

    // Input focus effects
    this.searchInput.addEventListener('focus', () => {
      this.searchInput.parentElement.classList.add('focused');
    });

    this.searchInput.addEventListener('blur', () => {
      this.searchInput.parentElement.classList.remove('focused');
    });
  }

  createMatrixEffect() {
    const matrixBg = document.getElementById('matrix-bg');
    const chars = '01アイウエオカキクケコサシスセソタチツテトナニヌネノハヒフヘホマミムメモヤユヨラリルレロワヲン';
    
    // Create matrix columns
    for (let i = 0; i < 50; i++) {
      const column = document.createElement('div');
      column.style.position = 'absolute';
      column.style.left = `${Math.random() * 100}%`;
      column.style.fontSize = '14px';
      column.style.color = 'rgba(74, 222, 128, 0.3)';
      column.style.fontFamily = 'monospace';
      column.style.whiteSpace = 'pre';
      column.style.lineHeight = '1.2';
      column.style.animation = `matrix ${10 + Math.random() * 20}s linear infinite`;
      column.style.animationDelay = `${Math.random() * 5}s`;
      
      // Generate random characters for this column
      let text = '';
      for (let j = 0; j < 20; j++) {
        text += chars[Math.floor(Math.random() * chars.length)] + '\n';
      }
      column.textContent = text;
      
      matrixBg.appendChild(column);
    }
  }

  async searchTorrents() {
    const query = this.searchInput.value.trim();
    if (!query) {
      this.showError('PLEASE ENTER A SEARCH QUERY');
      return;
    }

    this.showStatus('INITIALIZING KRAKEN PROTOCOLS...');
    this.clearResults();
    
    try {
      // Simulate realistic loading stages
      await this.simulateLoadingStages();
      
      const response = await fetch(`/search?q=${encodeURIComponent(query)}`);
      
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }
      
      this.showStatus('PARSING TORRENT DATA...');
      await this.delay(500);
      
      const data = await response.json();
      
      this.showStatus('RANKING RESULTS...');
      await this.delay(300);
      
      this.hideStatus();
      this.displayResults(data);
      
    } catch (error) {
      console.error('Search error:', error);
      this.hideStatus();
      this.showError('KRAKEN SCAN FAILED - CHECK CONNECTION');
    }
  }

  async simulateLoadingStages() {
    const stages = [
      { text: 'SCANNING PIRACY NETWORKS...', duration: 800 },
      { text: 'CONNECTING TO RARBG MIRROR...', duration: 600 },
      { text: 'BYPASSING SECURITY PROTOCOLS...', duration: 700 },
      { text: 'EXTRACTING TORRENT METADATA...', duration: 900 },
    ];

    for (let i = 0; i < stages.length; i++) {
      this.showStatus(stages[i].text);
      this.updateProgress((i + 1) / stages.length * 80); // 80% for loading stages
      await this.delay(stages[i].duration);
    }
  }

  showStatus(message) {
    this.statusText.textContent = message;
    this.statusText.classList.add('loading-dots');
    this.statusSection.classList.remove('hidden');
    this.statusSection.classList.add('fade-in');
  }

  hideStatus() {
    this.statusSection.classList.add('hidden');
    this.statusText.classList.remove('loading-dots');
    this.updateProgress(0);
  }

  updateProgress(percentage) {
    this.progressBar.style.width = `${percentage}%`;
  }

  clearResults() {
    this.results.innerHTML = '';
  }

  displayResults(torrents) {
    if (!torrents || torrents.length === 0) {
      this.results.innerHTML = `
        <div class="no-results">
          ⚠️ NO TORRENTS FOUND IN THE KRAKEN'S NET
          <br><br>
          TRY DIFFERENT SEARCH TERMS
        </div>
      `;
      return;
    }

    this.results.innerHTML = torrents.map((torrent, index) => `
      <div class="torrent-card slide-up" style="animation-delay: ${index * 0.1}s">
        <div class="torrent-title">
          ${this.escapeHtml(torrent.Name || 'UNKNOWN TORRENT')}
        </div>
        
        <div class="torrent-info">
          <div class="info-item">
            <span class="info-label">📦 SIZE:</span>
            <span class="info-value">${this.formatSize(torrent.Size)}</span>
          </div>
          <div class="info-item">
            <span class="info-label">📺 RES:</span>
            <span class="info-value">${torrent.Resolution || 'UNKNOWN'}</span>
          </div>
          <div class="info-item">
            <span class="info-label">🌱 SEEDS:</span>
            <span class="info-value">${torrent.Seeders || '?'}</span>
          </div>
          <div class="info-item">
            <span class="info-label">🩸 LEECH:</span>
            <span class="info-value">${torrent.Leechers || '?'}</span>
          </div>
          <div class="info-item">
            <span class="info-label">🎞️ SOURCE:</span>
            <span class="info-value">${torrent.Source || 'UNKNOWN'}</span>
          </div>
          <div class="info-item">
            <span class="info-label">🏴‍☠️ SITE:</span>
            <span class="info-value">${torrent.SiteName || 'UNKNOWN'}</span>
          </div>
        </div>
        
        <div class="torrent-actions">
          <button onclick="krakenEye.copyMagnet('${this.escapeHtml(torrent.MagnetLink || '')}')" 
                  class="magnet-btn">
            🧲 COPY MAGNET
          </button>
          <div class="kraken-score">
            🐉 KRAKEN SCORE: <span class="score-value">${(torrent.Score || 0).toFixed(1)}</span>
          </div>
        </div>
      </div>
    `).join('');
  }

  showError(message) {
    this.results.innerHTML = `
      <div class="error-message">
        ❌ ${message}
      </div>
    `;
  }

  async copyMagnet(magnetLink) {
    if (!magnetLink) {
      this.showNotification('❌ NO MAGNET LINK AVAILABLE', 'error');
      return;
    }

    try {
      await navigator.clipboard.writeText(magnetLink);
      this.showNotification('🧲 MAGNET LINK COPIED TO CLIPBOARD!', 'success');
    } catch (error) {
      console.error('Copy failed:', error);
      this.showNotification('❌ FAILED TO COPY MAGNET LINK', 'error');
    }
  }

  showNotification(message, type = 'info') {
    // Remove existing notifications
    const existing = document.querySelector('.notification');
    if (existing) {
      existing.remove();
    }

    const notification = document.createElement('div');
    notification.className = `notification fixed top-4 right-4 z-50 px-6 py-3 rounded-lg font-orbitron font-bold text-sm transform translate-x-full transition-transform duration-300`;
    
    if (type === 'success') {
      notification.classList.add('bg-green-500', 'text-black');
    } else if (type === 'error') {
      notification.classList.add('bg-red-500', 'text-white');
    } else {
      notification.classList.add('bg-cyan-500', 'text-black');
    }
    
    notification.textContent = message;
    document.body.appendChild(notification);

    // Animate in
    setTimeout(() => {
      notification.style.transform = 'translateX(0)';
    }, 100);

    // Animate out and remove
    setTimeout(() => {
      notification.style.transform = 'translateX(100%)';
      setTimeout(() => {
        if (notification.parentNode) {
          notification.remove();
        }
      }, 300);
    }, 3000);
  }

  formatSize(sizeGB) {
    if (!sizeGB || sizeGB === 0) return '? GB';
    if (sizeGB < 1) return `${(sizeGB * 1024).toFixed(0)} MB`;
    if (sizeGB >= 1024) return `${(sizeGB / 1024).toFixed(1)} TB`;
    return `${sizeGB.toFixed(1)} GB`;
  }

  escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
  }

  delay(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
  }
}

// Initialize the app
const krakenEye = new KrakenEye();

// Global function for button clicks
window.searchTorrents = () => krakenEye.searchTorrents();