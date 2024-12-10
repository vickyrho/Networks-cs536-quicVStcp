package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"github.com/quic-go/quic-go"
)

func main() {
	// Create a QUIC listener on port 8080
	listener, err := quic.ListenAddr("localhost:8080", generateTLSConfig(), nil)
	if err != nil {
		log.Fatalf("Failed to start QUIC server: %v", err)
	}
	fmt.Println("QUIC server is running on localhost:8080")

	for {
		// Accept a new session
		session, err := listener.Accept(context.Background())
		if err != nil {
			log.Printf("Failed to accept session: %v", err)
			continue
		}
		go handleSession(session)
	}
}

func handleSession(session quic.Connection) {
	fmt.Printf("New session from: %s\n", session.RemoteAddr())

	// Accept a new stream
	stream, err := session.AcceptStream(context.Background())
	if err != nil {
		log.Printf("Failed to accept stream: %v", err)
		return
	}
	defer stream.Close()

	// Read data from the stream
	buf := make([]byte, 1024)
	for {
		n, err := stream.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Stream closed by client")
				break
			}
			log.Printf("Failed to read from stream: %v", err)
			break
		}
		fmt.Printf("Received: %s\n", string(buf[:n]))

		// Write a response
		_, err = stream.Write([]byte("Hello from QUIC server!"))
		if err != nil {
			log.Printf("Failed to write to stream: %v", err)
			break
		}
	}
}

func generateTLSConfig() *tls.Config {
	// Check if certificate and key files exist
	certFile := "server.crt"
	keyFile := "server.key"
	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		log.Fatalf("TLS certificate not found: %s", certFile)
	}
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		log.Fatalf("TLS key not found: %s", keyFile)
	}

	// Load the certificate and key
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Failed to load TLS certificate and key: %v", err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"quic-echo-example"},
	}
}
