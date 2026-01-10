package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// handleIndex serves the main HTML page
func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	htmlPath := filepath.Join("static", "index.html")
	htmlContent, err := os.ReadFile(htmlPath)
	if err != nil {
		log.Printf("Error reading index.html: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(htmlContent)
}

// handleProcess handles the POST /api/process endpoint
func handleProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse JSON request
	var request struct {
		Input string `json:"input"`
	}

	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Standardize input
	standardized := StandardizeInput(request.Input)

	// Process input
	result := ProcessInput(standardized)

	// Create response
	response := struct {
		Standardized string `json:"standardized"`
		Result       string `json:"result"`
	}{
		Standardized: standardized,
		Result:       result,
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// handleStatic serves static files (CSS)
func handleStatic(w http.ResponseWriter, r *http.Request) {
	// Extract filename from path
	filename := r.URL.Path[len("/static/"):]
	if filename == "" {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	filePath := filepath.Join("static", filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Set appropriate content type based on file extension
	if filepath.Ext(filename) == ".css" {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	}

	// Serve file
	http.ServeFile(w, r, filePath)
}
