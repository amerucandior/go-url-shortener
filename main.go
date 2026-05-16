package main

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"strings"
	"sync"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateCode(length int) (string, error) {
	code := make([]byte, length)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		code[i] = alphabet[n.Int64()]
	}
	return string(code), nil
}

// package-level store — declared so both handlers can reach it
var store = struct {
	mu   sync.RWMutex
	urls map[string]string
}{urls: make(map[string]string)}

func urlShortenerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		URL string `json:"url"`
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	code, err := generateCode(6)
	if err != nil {
		http.Error(w, "failed to generate code", http.StatusInternalServerError)
		return
	}

	store.mu.Lock()
	store.urls[code] = request.URL
	store.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"short_url": "http://localhost:8080/" + code,
	})
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/")

	store.mu.RLock()
	original, ok := store.urls[code]
	store.mu.RUnlock()

	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, original, http.StatusFound)
}

func main() {
	http.HandleFunc("/shorten", urlShortenerHandler)
	http.HandleFunc("/", redirectHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
