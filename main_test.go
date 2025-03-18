package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArtistHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/artist/1", nil)
	w := httptest.NewRecorder()

	artistHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestSearchHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/search?q=queen", nil)
	w := httptest.NewRecorder()

	searchHandler(w, req)
	var results []Artist
	json.Unmarshal(w.Body.Bytes(), &results)

	if len(results) < 1 {
		t.Error("Expected at least 1 search result")
	}
}
