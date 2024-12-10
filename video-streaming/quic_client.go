package main

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/quic-go/quic-go/http3"
)

// DASHManifest represents a simplified MPEG-DASH manifest structure
type DASHManifest struct {
	XMLName   xml.Name  `xml:"MPD"`
	Periods   []Period  `xml:"Period"`
	BaseURL   string    `xml:"BaseURL"`
}

type Period struct {
	AdaptationSets []AdaptationSet `xml:"AdaptationSet"`
}

type AdaptationSet struct {
	Representations []Representation `xml:"Representation"`
}

type Representation struct {
	ID          string `xml:"id,attr"`
	SegmentList SegmentList `xml:"SegmentList"`
}

type SegmentList struct {
	Initialization Segment `xml:"Initialization"`
	Segments       []Segment `xml:"SegmentURL"`
}

type Segment struct {
	SourceURL string `xml:"sourceURL,attr"`
}

func main() {
	// Create an HTTP/3 client
	client := &http.Client{
		Transport: &http3.RoundTripper{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Skip TLS verification for development
			},
		},
	}

	// URL to the DASH manifest
	manifestURL := "https://localhost:9000/manifest.mpd"

	// Step 1: Download and parse the manifest
	manifest, err := fetchManifest(client, manifestURL)
	if err != nil {
		log.Fatalf("Failed to fetch manifest: %v", err)
	}

	// Step 2: Download initialization segments for each Representation
	for _, period := range manifest.Periods {
		for _, adaptationSet := range period.AdaptationSets {
			for _, representation := range adaptationSet.Representations {
				initSegmentURL := buildFullURL(manifest.BaseURL, representation.SegmentList.Initialization.SourceURL)
				log.Printf("Downloading initialization segment: %s\n", initSegmentURL)
				if err := downloadFile(client, initSegmentURL, path.Base(initSegmentURL)); err != nil {
					log.Fatalf("Failed to download initialization segment: %v", err)
				}

				// Step 3: Download video chunks for this Representation
				for _, segment := range representation.SegmentList.Segments {
					segmentURL := buildFullURL(manifest.BaseURL, segment.SourceURL)
					log.Printf("Downloading segment: %s\n", segmentURL)
					if err := downloadFile(client, segmentURL, path.Base(segmentURL)); err != nil {
						log.Fatalf("Failed to download segment: %v", err)
					}
				}
			}
		}
	}

	log.Println("Video streaming complete.")
}

// fetchManifest downloads and parses the DASH manifest file
func fetchManifest(client *http.Client, url string) (*DASHManifest, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch manifest: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch manifest, status: %s", resp.Status)
	}

	var manifest DASHManifest
	decoder := xml.NewDecoder(resp.Body)
	if err := decoder.Decode(&manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	// If the manifest contains a BaseURL, prepend it to segment URLs
	if manifest.BaseURL == "" && strings.HasSuffix(url, "/manifest.mpd") {
		manifest.BaseURL = strings.TrimSuffix(url, "/manifest.mpd")
	}

	return &manifest, nil
}

// buildFullURL constructs the full URL for a segment
func buildFullURL(baseURL, segmentURL string) string {
	if strings.HasPrefix(segmentURL, "http") {
		return segmentURL
	}
	return fmt.Sprintf("%s/%s", strings.TrimSuffix(baseURL, "/"), strings.TrimPrefix(segmentURL, "/"))
}

// downloadFile downloads a file and saves it locally
func downloadFile(client *http.Client, url, filename string) error {
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch file, status: %s", resp.Status)
	}

	// Create the local file
	out, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Write the response body to the file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
