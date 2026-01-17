package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	url := fmt.Sprintf("http://localhost:%s", port)

	// Channel for HTTP-triggered shutdown
	shutdownChan := make(chan struct{}, 1)

	// Create HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/api/process", handleProcess)
	mux.HandleFunc("/api/shutdown", func(w http.ResponseWriter, r *http.Request) {
		handleShutdown(w, r, shutdownChan)
	})
	mux.HandleFunc("/static/", handleStatic)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Channel to signal when server is ready
	serverReady := make(chan bool, 1)

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s", url)
		serverReady <- true // Signal that we're attempting to start
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for server to be ready
	<-serverReady
	
	// Verify server is actually accepting connections before opening browser
	if !waitForServerReady(url, 5*time.Second) {
		log.Printf("Warning: Server may not be ready, but opening browser anyway")
		log.Printf("If the page doesn't load, please wait a moment and refresh")
	}

	// Open browser
	if err := openBrowser(url); err != nil {
		log.Printf("Warning: Failed to open browser: %v", err)
		log.Printf("Please open %s manually in your browser", url)
	}

	// Wait for interrupt signal or HTTP-triggered shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	log.Println("Server is running. Press Ctrl+C to stop.")
	
	// Wait for either signal-based shutdown or HTTP-triggered shutdown
	select {
	case <-sigChan:
		log.Println("Received shutdown signal (Ctrl+C)")
	case <-shutdownChan:
		log.Println("Received shutdown request from web interface")
	}

	// Graceful shutdown
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}

// waitForServerReady polls the server to verify it's accepting connections
func waitForServerReady(url string, timeout time.Duration) bool {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}
	
	deadline := time.Now().Add(timeout)
	
	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusMethodNotAllowed {
				return true
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	
	return false
}

// openBrowser opens the default browser to the specified URL
func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return cmd.Start()
}
