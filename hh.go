package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func serveDirectory(w http.ResponseWriter, r *http.Request, dir string) {
	// Log the incoming request
	log.Printf("Received request: %s %s", r.Method, r.URL.Path)

	// Build the file path from the request's URL and the provided directory
	filePath := filepath.Join(dir, r.URL.Path[1:]) // Remove the leading '/'

	// Try to open the file
	file, err := os.Open(filePath)
	if err != nil {
		http.NotFound(w, r)
		log.Printf("File not found: %v", err)
		return
	}
	defer file.Close()

	// Serve the file as the response
	http.ServeFile(w, r, filePath)
}

func main() {
	// Parse command-line arguments
	dir := flag.String("dir", ".", "Directory to serve files from")
	addr := flag.String("addr", ":8080", "Address to listen on")
	flag.Parse()

	// Set up logging
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Printf("Starting server... Serving directory: %s", *dir)

	// Validate the provided directory
	if _, err := os.Stat(*dir); err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Define the request handler function with the provided directory
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveDirectory(w, r, *dir)
	})

	// Start the HTTP server
	log.Printf("Listening on http://localhost%s", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
