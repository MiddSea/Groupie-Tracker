document.getElementById('searchInput').addEventListener('input', function(e) {
    const query = e.target.value.trim().toLowerCase();
    const suggestions = document.getElementById('suggestions');
    
    if (query.length < 1) {
        suggestions.style.display = 'none';
        return;
    }

    fetch(`/search?q=${encodeURIComponent(query)}`)
        .then(response => response.json())
        .then(artists => {
            suggestions.innerHTML = '';
            artists.forEach(artist => {
                if (artist.name.toLowerCase().startsWith(query)) {
                    const div = document.createElement('div');
                    div.className = 'suggestion-item';
                    div.textContent = artist.name;
                    div.onclick = () => window.location.href = `/artist/${artist.id}`;
                    suggestions.appendChild(div);
                }
            });
            suggestions.style.display = artists.length ? 'block' : 'none';
        });
});

// Add loading feedback
searchInput.addEventListener('input', function() {
    suggestions.innerHTML = '<div class="loading">Searching...</div>';
  });
  
  // Add error feedback
  function showSearchError() {
    suggestions.innerHTML = '<div class="error">Search unavailable</div>';
  }