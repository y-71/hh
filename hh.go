package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

func findAvailablePort(basePort int) (int, error) {
	for port := basePort; port < basePort+100; port++ {
		address := fmt.Sprintf(":%d", port)
		ln, err := net.Listen("tcp", address)
		if err == nil {
			ln.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available ports in the range %d-%d", basePort, basePort+100)
}

func serveDirectory(w http.ResponseWriter, r *http.Request, dir string) {
	// Build the file path from the request's URL and the provided directory
	filePath := filepath.Join(dir, r.URL.Path[1:]) // Remove the leading '/'

	// Try to open the file
	file, err := os.Open(filePath)
	if err != nil {
		http.NotFound(w, r)
		log.Printf(color.RedString("Error: %s %s - File not found: %v"), r.Method, r.URL.Path, err)
		return
	}
	defer file.Close()

	// Log the incoming request with the response code in green
	log.Printf(color.GreenString("Received request: %s %s - Status: %d"), r.Method, r.URL.Path, http.StatusOK)

	// Serve the file as the response
	http.ServeFile(w, r, filePath)
}

func main() {
	// Parse command-line arguments
	dir := flag.String("dir", ".", "Directory to serve files from")
	basePort := flag.Int("port", 8080, "Base port to listen on")
	flag.Parse()

	// Set up logging with colors
	log.SetFlags(0)

	// Find an available port
	port, err := findAvailablePort(*basePort)
	if err != nil {
		log.Fatalf(color.RedString("Error: %v"), err)
	}

	// Define the request handler function with the provided directory
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveDirectory(w, r, *dir)
	})

	// Start the HTTP server
	serverAddr := fmt.Sprintf(":%d", port)
	log.Printf(color.BlueString("Starting server... Serving directory: %s"), *dir)
	log.Printf(color.BlueString("Listening on http://localhost%s"), serverAddr)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf(color.RedString("Server error: %v"), err)
	}
}
