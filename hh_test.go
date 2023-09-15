package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestServeDirectory(t *testing.T) {
	// Create a temporary directory for testing
	testDir := t.TempDir()

	// Create a test file in the temporary directory
	testFile := filepath.Join(testDir, "testfile.txt")
	if err := os.WriteFile(testFile, []byte("Test content"), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testFile)

	// Create a request to the test server
	req := httptest.NewRequest(http.MethodGet, "/testfile.txt", nil)
	w := httptest.NewRecorder()

	// Call the serveDirectory function with the temporary directory
	serveDirectory(w, req, testDir)

	// Check the response status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check the response body
	expectedBody := "Test content"
	if w.Body.String() != expectedBody {
		t.Errorf("Expected response body %q, got %q", expectedBody, w.Body.String())
	}
}

func TestServeDirectoryFileNotFound(t *testing.T) {
	// Create a temporary directory for testing
	testDir := t.TempDir()

	// Create a request to the test server for a non-existing file
	req := httptest.NewRequest(http.MethodGet, "/non_existing.txt", nil)
	w := httptest.NewRecorder()

	// Call the serveDirectory function with the temporary directory
	serveDirectory(w, req, testDir)

	// Check the response status code
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}
