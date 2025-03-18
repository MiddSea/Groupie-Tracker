package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

// Data structures
type Artist struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	ImageURL     string   `json:"image"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type Date struct {
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type ArtistDetails struct {
	Artist
	Locations []string
	Dates     []string
	Concerts  map[string][]string
}

type RelationsResponse struct {
	Index []Relation `json:"index"`
}

var (
	artistsDetails []ArtistDetails
	dataMutex      sync.RWMutex
	templates      *template.Template
	apiCache       = cache.New(5*time.Minute, 10*time.Minute)
)

func init() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
}

func fetchWithCache(url string, v interface{}) error {
	if cached, found := apiCache.Get(url); found {
		return json.Unmarshal(cached.([]byte), v)
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	apiCache.Set(url, body, cache.DefaultExpiration)
	return json.Unmarshal(body, v)
}

func fetchAllData() ([]ArtistDetails, error) {
	var wg sync.WaitGroup
	var artists []Artist
	var locations []Location
	var dates []Date
	var relations []Relation
	var errs [4]error

	wg.Add(4)

	go func() {
		defer wg.Done()
		artists, errs[0] = fetchArtists()
	}()

	go func() {
		defer wg.Done()
		locations, errs[1] = fetchLocations()
	}()

	go func() {
		defer wg.Done()
		dates, errs[2] = fetchDates()
	}()

	go func() {
		defer wg.Done()
		relations, errs[3] = fetchRelations()
	}()

	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return nil, fmt.Errorf("data fetch error: %v", err)
		}
	}

	return combineData(artists, locations, dates, relations), nil
}

func combineData(artists []Artist, locations []Location, dates []Date, relations []Relation) []ArtistDetails {
	details := make([]ArtistDetails, len(artists))

	for i, artist := range artists {
		detail := ArtistDetails{
			Artist:   artist,
			Concerts: make(map[string][]string),
		}

		for _, loc := range locations {
			if loc.ID == artist.ID {
				detail.Locations = loc.Locations
				break
			}
		}

		for _, d := range dates {
			if d.ID == artist.ID {
				detail.Dates = d.Dates
				break
			}
		}

		for _, rel := range relations {
			if rel.ID == artist.ID {
				detail.Concerts = rel.DatesLocations
				break
			}
		}

		details[i] = detail
	}
	return details
}

func fetchArtists() ([]Artist, error) {
	var artists []Artist
	err := fetchWithCache("https://groupietrackers.herokuapp.com/api/artists", &artists)
	return artists, err
}

func fetchLocations() ([]Location, error) {
	var wrapper struct{ Index []Location }
	err := fetchWithCache("https://groupietrackers.herokuapp.com/api/locations", &wrapper)
	return wrapper.Index, err
}

func fetchDates() ([]Date, error) {
	var wrapper struct{ Index []Date }
	err := fetchWithCache("https://groupietrackers.herokuapp.com/api/dates", &wrapper)
	return wrapper.Index, err
}

func fetchRelations() ([]Relation, error) {
	var response RelationsResponse
	err := fetchWithCache("https://groupietrackers.herokuapp.com/api/relation", &response)
	return response.Index, err
}

func fetchDataPeriodically() {
	ticker := time.NewTicker(1 * time.Minute)
	fetchData()
	for range ticker.C {
		fetchData()
	}
}

func fetchData() error {
	combined, err := fetchAllData()
	if err != nil {
		log.Println("Data fetch error:", err)
		return err
	}

	dataMutex.Lock()
	artistsDetails = combined
	dataMutex.Unlock()
	return nil
}

func main() {
	go fetchDataPeriodically()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", middleware(homeHandler))
	http.HandleFunc("/artist/", middleware(artistHandler))
	http.HandleFunc("/search", middleware(searchHandler))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %v", err)
				renderError(w, http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

func renderError(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	err := templates.ExecuteTemplate(w, "error.html", map[string]interface{}{
		"StatusCode": status,
		"Message":    http.StatusText(status),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	dataMutex.RLock()
	defer dataMutex.RUnlock()

	if len(artistsDetails) == 0 {
		renderError(w, http.StatusServiceUnavailable)
		return
	}

	if err := templates.ExecuteTemplate(w, "home.html", artistsDetails); err != nil {
		renderError(w, http.StatusInternalServerError)
	}
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/artist/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		renderError(w, http.StatusBadRequest)
		return
	}

	dataMutex.RLock()
	defer dataMutex.RUnlock()

	var artist *ArtistDetails
	for i := range artistsDetails {
		if artistsDetails[i].ID == id {
			artist = &artistsDetails[i]
			break
		}
	}

	if artist == nil {
		renderError(w, http.StatusNotFound)
		return
	}

	if err := templates.ExecuteTemplate(w, "artist.html", artist); err != nil {
		renderError(w, http.StatusInternalServerError)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
	if query == "" {
		json.NewEncoder(w).Encode([]Artist{})
		return
	}

	dataMutex.RLock()
	defer dataMutex.RUnlock()

	var results []Artist
	for _, ad := range artistsDetails {
		if strings.HasPrefix(strings.ToLower(ad.Name), query) {
			results = append(results, ad.Artist)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
