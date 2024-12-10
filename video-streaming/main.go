package main

import (
	"log"
	"net/http"
	"path"
	"strings"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Serve the requested file with proper MIME types
		filePath := r.URL.Path
		if strings.HasSuffix(filePath, ".mpd") {
			w.Header().Set("Content-Type", "application/dash+xml")
		} else if strings.HasSuffix(filePath, ".m4s") {
			w.Header().Set("Content-Type", "video/iso.segment")
		} else {
			w.Header().Set("Content-Type", "application/octet-stream")
		}

		http.ServeFile(w, r, path.Join("./video", filePath))
	})

	log.Println("Starting DASH server with CORS on http://localhost:8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

