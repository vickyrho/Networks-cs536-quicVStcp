package main

import (
	"log"
	"net/http"

	"github.com/quic-go/quic-go/http3"
)

func main() {
	// Create a simple HTTP handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Log the request
		log.Printf("Request: %s %s\n", r.Method, r.URL.Path)

		// Serve files from the ./video directory
		http.ServeFile(w, r, "./video"+r.URL.Path)
	})

	// Create an HTTP/3 server
	server := http3.Server{
		Addr:    ":9000",
		Handler: handler,
	}

	certFile := "server.crt"
	keyFile := "server.key"

	log.Println("Starting HTTP/3 server on https://localhost:9000...")
	err := server.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to start HTTP/3 server: %v", err)
	}
}