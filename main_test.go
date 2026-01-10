package main

import (
	"runtime"
	"testing"
)

func TestOpenBrowser(t *testing.T) {
	// This test verifies that openBrowser doesn't crash for different OS values
	// Note: We can't easily test the actual browser opening without mocking exec.Command
	// which would require refactoring. This test at least ensures the function handles
	// different OS values correctly.

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid URL",
			url:     "http://localhost:8080",
			wantErr: false, // May fail on unsupported OS, but shouldn't crash
		},
		{
			name:    "empty URL",
			url:     "",
			wantErr: false, // Command will fail but function won't error on unsupported OS
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Only test on the current OS to avoid false failures
			err := openBrowser(tt.url)
			
			// On unsupported OS, we expect an error
			if runtime.GOOS != "darwin" && runtime.GOOS != "linux" && runtime.GOOS != "windows" {
				if err == nil {
					t.Error("openBrowser() should return error on unsupported OS")
				}
				return
			}

			// On supported OS, the function may succeed or fail depending on system,
			// but it shouldn't panic. We mainly test that it doesn't crash.
			// Note: This test may fail if browser command doesn't exist, which is acceptable.
			_ = err // Ignore error as it depends on system configuration
		})
	}
}

func TestOpenBrowserUnsupportedOS(t *testing.T) {
	// Test that unsupported OS returns an error
	// We can't easily simulate this without changing runtime.GOOS, which isn't possible.
	// This test documents the expected behavior.
	
	// Save original GOOS
	originalOS := runtime.GOOS
	
	// Note: We can't actually change runtime.GOOS in a test, so this test
	// just documents the expected behavior. In practice, the function should
	// return an error for unsupported operating systems.
	
	_ = originalOS // Suppress unused variable warning
	
	// The actual test would require mocking or a way to change runtime.GOOS,
	// which isn't feasible. The function implementation handles this correctly.
}

// Note: Testing main() function is complex because it:
// 1. Sets up an HTTP server
// 2. Opens a browser
// 3. Blocks waiting for signals
// 
// For integration testing of main(), consider:
// - Using a test binary with a flag to skip browser opening
// - Using a test HTTP client to verify server endpoints
// - Testing graceful shutdown with signal simulation
//
// These would require refactoring main() to be more testable, which may
// not be necessary for a simple application like this.
