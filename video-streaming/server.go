package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/quic-go/quic-go"
	"io"
)

func main() {
	// Generate TLS configuration
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		fmt.Printf("TLS certificate loading error: %v\n", err)
		return
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"quic-test"},
	}

	// Start QUIC listener
	listener, err := quic.ListenAddr("localhost:4242", tlsConfig, nil)
	if err != nil {
		fmt.Printf("Error starting QUIC listener: %v\n", err)
		return
	}
	fmt.Println("Server is listening on localhost:4242")

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			fmt.Printf("Error during Accept: %v\n", err)
			continue
		}
		fmt.Printf("Connection accepted: %s\n", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn quic.Connection) {
	defer func() {
		err := conn.CloseWithError(0, "Server shutting down")
		if err != nil {
			fmt.Printf("Error closing connection: %v\n", err)
		}
	}()

	// Accept the stream from the client
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		fmt.Printf("Error accepting stream: %v\n", err)
		return
	}
	defer stream.Close()

	// Read the message from the client
	buffer := make([]byte, 1024)
	n, err := stream.Read(buffer)
	if err != nil && err != io.EOF {
		fmt.Printf("Error reading from stream: %v\n", err)
		return
	}
	message := string(buffer[:n])
	fmt.Printf("Received message: %s\n", message)

	// Echo the message back to the client
	_, err = stream.Write([]byte("Message received: " + message))
	if err != nil {
		fmt.Printf("Error writing to stream: %v\n", err)
	}
}
