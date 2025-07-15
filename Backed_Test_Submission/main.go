package main

import (
	"UrlShortner/Logging_Middleware"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

type URLInfo struct {
	OriginalURL string
	CreatedAt   time.Time
	Expiry      time.Time
	Clicks      int
	ClickLogs   []time.Time
}

var (
	store = make(map[string]*URLInfo)
	mu    sync.Mutex
)

func generateShortcode() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid Method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		URL       string `json:"url"`
		Validity  int    `json:"validity"`
		Shortcode string `json:"shortcode"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Log.Log("backend", "error", "handler", "invalid request body")
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		Log.Log("backend", "warn", "handler", "missing url field")
		http.Error(w, "Missing URL", http.StatusBadRequest)
		return
	}

	if req.Validity <= 0 {
		req.Validity = 30
	}

	if req.Shortcode == "" {
		req.Shortcode = generateShortcode()
	}

	mu.Lock()
	if _, exists := store[req.Shortcode]; exists {
		mu.Unlock()
		Log.Log("backend", "warn", "handler", "shortcode already exists")
		http.Error(w, "Shortcode already exists", http.StatusConflict)
		return
	}

	store[req.Shortcode] = &URLInfo{
		OriginalURL: req.URL,
		CreatedAt:   time.Now(),
		Expiry:      time.Now().Add(time.Duration(req.Validity) * time.Minute),
	}
	mu.Unlock()

	resp := map[string]string{
		"shortLink": fmt.Sprintf("http://localhost:8080/%s", req.Shortcode),
		"expiry":    store[req.Shortcode].Expiry.Format(time.RFC3339),
	}
	Log.Log("backend", "info", "handler", fmt.Sprintf("Shortcode %s created", req.Shortcode))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	shortcode := strings.TrimPrefix(r.URL.Path, "/")
	if shortcode == "" {
		http.NotFound(w, r)
		return
	}

	mu.Lock()
	info, ok := store[shortcode]
	if !ok {
		mu.Unlock()
		Log.Log("backend", "warn", "handler", "shortcode not found")
		http.NotFound(w, r)
		return
	}
	if time.Now().After(info.Expiry) {
		mu.Unlock()
		Log.Log("backend", "info", "handler", "shortcode expired")
		http.Error(w, "Link expired", http.StatusGone)
		return
	}
	info.Clicks++
	info.ClickLogs = append(info.ClickLogs, time.Now())
	mu.Unlock()

	Log.Log("backend", "info", "handler", fmt.Sprintf("Redirected to %s", info.OriginalURL))
	http.Redirect(w, r, info.OriginalURL, http.StatusFound)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	shortcode := strings.TrimPrefix(r.URL.Path, "/shorturls/")
	mu.Lock()
	info, ok := store[shortcode]
	mu.Unlock()

	if !ok {
		Log.Log("backend", "warn", "handler", "stats requested for non-existent shortcode")
		http.NotFound(w, r)
		return
	}

	resp := map[string]interface{}{
		"originalUrl": info.OriginalURL,
		"createdAt":   info.CreatedAt.Format(time.RFC3339),
		"expiry":      info.Expiry.Format(time.RFC3339),
		"clicks":      info.Clicks,
	}
	Log.Log("backend", "debug", "handler", fmt.Sprintf("Stats returned for %s", shortcode))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/shorturls", shortenHandler)
	http.HandleFunc("/shorturls/", statsHandler)
	http.HandleFunc("/", redirectHandler)

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
