package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandleIndex(t *testing.T) {
	// Ensure static directory exists
	if _, err := os.Stat("static"); os.IsNotExist(err) {
		t.Skip("static directory does not exist, skipping test")
	}

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		checkContent   bool
	}{
		{
			name:           "GET request",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			checkContent:   true,
		},
		{
			name:           "POST request",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
			checkContent:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/", nil)
			w := httptest.NewRecorder()

			handleIndex(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("handleIndex() status = %d, want %d", w.Code, tt.expectedStatus)
			}

			if tt.checkContent {
				contentType := w.Header().Get("Content-Type")
				if contentType != "text/html; charset=utf-8" {
					t.Errorf("handleIndex() Content-Type = %q, want %q", contentType, "text/html; charset=utf-8")
				}

				if w.Body.Len() == 0 {
					t.Error("handleIndex() returned empty body")
				}
			}
		})
	}
}

func TestHandleProcess(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		requestBody    string
		expectedStatus int
		validateJSON   bool
	}{
		{
			name:           "valid POST request",
			method:         http.MethodPost,
			requestBody:    `{"input":"hello world"}`,
			expectedStatus: http.StatusOK,
			validateJSON:   true,
		},
		{
			name:           "POST with empty input",
			method:         http.MethodPost,
			requestBody:    `{"input":""}`,
			expectedStatus: http.StatusOK,
			validateJSON:   true,
		},
		{
			name:           "POST with multiline input",
			method:         http.MethodPost,
			requestBody:    `{"input":"hello\nworld"}`,
			expectedStatus: http.StatusOK,
			validateJSON:   true,
		},
		{
			name:           "invalid JSON",
			method:         http.MethodPost,
			requestBody:    `{"input":invalid}`,
			expectedStatus: http.StatusBadRequest,
			validateJSON:   false,
		},
		{
			name:           "missing input field",
			method:         http.MethodPost,
			requestBody:    `{}`,
			expectedStatus: http.StatusOK,
			validateJSON:   true,
		},
		{
			name:           "GET request",
			method:         http.MethodGet,
			requestBody:    "",
			expectedStatus: http.StatusMethodNotAllowed,
			validateJSON:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/api/process", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handleProcess(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("handleProcess() status = %d, want %d", w.Code, tt.expectedStatus)
			}

			if tt.validateJSON && w.Code == http.StatusOK {
				contentType := w.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("handleProcess() Content-Type = %q, want %q", contentType, "application/json")
				}

				var response struct {
					Standardized string `json:"standardized"`
					Result       string `json:"result"`
				}
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("handleProcess() returned invalid JSON: %v", err)
				}

				// Verify that result is uppercase of standardized input
				expectedResult := ProcessInput(response.Standardized)
				if response.Result != expectedResult {
					t.Errorf("handleProcess() result = %q, want %q", response.Result, expectedResult)
				}
			}
		})
	}
}

func TestHandleStatic(t *testing.T) {
	// Ensure static directory exists
	if _, err := os.Stat("static"); os.IsNotExist(err) {
		t.Skip("static directory does not exist, skipping test")
	}

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		checkContentType bool
	}{
		{
			name:             "CSS file",
			path:             "/static/style.css",
			expectedStatus:   http.StatusOK,
			checkContentType: true,
		},
		{
			name:             "non-existent file",
			path:             "/static/nonexistent.css",
			expectedStatus:   http.StatusNotFound,
			checkContentType: false,
		},
		{
			name:             "empty filename",
			path:             "/static/",
			expectedStatus:   http.StatusNotFound,
			checkContentType: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			handleStatic(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("handleStatic() status = %d, want %d", w.Code, tt.expectedStatus)
			}

			if tt.checkContentType && w.Code == http.StatusOK {
				contentType := w.Header().Get("Content-Type")
				expectedType := "text/css; charset=utf-8"
				if contentType != expectedType {
					t.Errorf("handleStatic() Content-Type = %q, want %q", contentType, expectedType)
				}
			}
		})
	}
}

// TestHandleProcessIntegration tests the full flow of standardization and processing
func TestHandleProcessIntegration(t *testing.T) {
	testCases := []struct {
		name             string
		input            string
		expectedStandard string
		expectedResult   string
	}{
		{
			name:             "simple text",
			input:            "hello world",
			expectedStandard: "hello world",
			expectedResult:   "HELLO WORLD",
		},
		{
			name:             "text with newlines",
			input:            "hello\nworld",
			expectedStandard: "hello world",
			expectedResult:   "HELLO WORLD",
		},
		{
			name:             "text with extra whitespace",
			input:            "  hello   world  ",
			expectedStandard: "hello world",
			expectedResult:   "HELLO WORLD",
		},
		{
			name:             "text with Unicode line separator",
			input:            "hello\u2028world",
			expectedStandard: "hello world",
			expectedResult:   "HELLO WORLD",
		},
		{
			name:             "text with Unicode paragraph separator",
			input:            "hello\u2029world",
			expectedStandard: "hello world",
			expectedResult:   "HELLO WORLD",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requestBody := map[string]string{"input": tc.input}
			jsonBody, _ := json.Marshal(requestBody)

			req := httptest.NewRequest(http.MethodPost, "/api/process", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handleProcess(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("handleProcess() status = %d, want %d", w.Code, http.StatusOK)
			}

			var response struct {
				Standardized string `json:"standardized"`
				Result       string `json:"result"`
			}
			if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if response.Standardized != tc.expectedStandard {
				t.Errorf("Standardized = %q, want %q", response.Standardized, tc.expectedStandard)
			}
			if response.Result != tc.expectedResult {
				t.Errorf("Result = %q, want %q", response.Result, tc.expectedResult)
			}
		})
	}
}

// TestGetStaticDir tests the getStaticDir function
func TestGetStaticDir(t *testing.T) {
	staticDir := getStaticDir()
	
	// Verify that the returned path exists and is a directory
	if staticDir == "" {
		t.Error("getStaticDir() returned empty string")
	}
	
	// Check if the static directory exists at the returned path
	info, err := os.Stat(staticDir)
	if err != nil {
		t.Errorf("getStaticDir() returned path %q that does not exist: %v", staticDir, err)
		return
	}
	
	if !info.IsDir() {
		t.Errorf("getStaticDir() returned path %q that is not a directory", staticDir)
	}
	
	// Verify that index.html exists in the static directory
	indexPath := filepath.Join(staticDir, "index.html")
	if _, err := os.Stat(indexPath); err != nil {
		t.Errorf("getStaticDir() returned path %q but index.html not found: %v", staticDir, err)
	}
	
	// Verify that style.css exists in the static directory
	cssPath := filepath.Join(staticDir, "style.css")
	if _, err := os.Stat(cssPath); err != nil {
		t.Errorf("getStaticDir() returned path %q but style.css not found: %v", staticDir, err)
	}
}

// TestCopyButtonHTMLStructure tests that the copy button and related HTML elements exist in the index page
func TestCopyButtonHTMLStructure(t *testing.T) {
	// Ensure static directory exists
	if _, err := os.Stat("static"); os.IsNotExist(err) {
		t.Skip("static directory does not exist, skipping test")
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handleIndex(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("handleIndex() status = %d, want %d", w.Code, http.StatusOK)
	}

	htmlContent := w.Body.String()

	// Verify copy button exists
	if !strings.Contains(htmlContent, `id="copy-btn"`) {
		t.Error("Copy button with id 'copy-btn' not found in HTML")
	}

	// Verify copy button has correct type
	if !strings.Contains(htmlContent, `<button id="copy-btn" type="button"`) {
		t.Error("Copy button should have type='button'")
	}

	// Verify copy button has aria-label for accessibility
	if !strings.Contains(htmlContent, `aria-label="Copy result to clipboard"`) {
		t.Error("Copy button should have aria-label attribute for accessibility")
	}

	// Verify copy feedback element exists
	if !strings.Contains(htmlContent, `id="copy-feedback"`) {
		t.Error("Copy feedback element with id 'copy-feedback' not found in HTML")
	}

	// Verify copy feedback has accessibility attributes
	if !strings.Contains(htmlContent, `role="status"`) {
		t.Error("Copy feedback should have role='status' for accessibility")
	}
	if !strings.Contains(htmlContent, `aria-live="polite"`) {
		t.Error("Copy feedback should have aria-live='polite' for accessibility")
	}
	if !strings.Contains(htmlContent, `aria-atomic="true"`) {
		t.Error("Copy feedback should have aria-atomic='true' for accessibility")
	}

	// Verify result textarea exists (required for copy functionality)
	if !strings.Contains(htmlContent, `id="result-text"`) {
		t.Error("Result textarea with id 'result-text' not found in HTML")
	}
}

// TestCopyButtonJavaScript tests that the copy functionality JavaScript is present
func TestCopyButtonJavaScript(t *testing.T) {
	// Ensure static directory exists
	if _, err := os.Stat("static"); os.IsNotExist(err) {
		t.Skip("static directory does not exist, skipping test")
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handleIndex(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("handleIndex() status = %d, want %d", w.Code, http.StatusOK)
	}

	htmlContent := w.Body.String()

	// Verify copy button event listener is present
	if !strings.Contains(htmlContent, `copyBtn.addEventListener('click'`) && !strings.Contains(htmlContent, `copyBtn.addEventListener("click"`) {
		t.Error("Copy button click event listener not found in JavaScript")
	}

	// Verify clipboard API check function exists
	if !strings.Contains(htmlContent, `isClipboardAPIAvailable`) {
		t.Error("Clipboard API availability check function not found")
	}

	// Verify execCommand fallback function exists
	if !strings.Contains(htmlContent, `copyWithExecCommand`) {
		t.Error("execCommand fallback copy function not found")
	}

	// Verify main copy function exists
	if !strings.Contains(htmlContent, `copyToClipboard`) {
		t.Error("Main copyToClipboard function not found")
	}

	// Verify error handling exists
	if !strings.Contains(htmlContent, `catch`) {
		t.Error("Error handling (catch block) not found in copy functionality")
	}

	// Verify feedback function exists
	if !strings.Contains(htmlContent, `showCopyFeedback`) {
		t.Error("showCopyFeedback function not found")
	}

	// Verify keyboard support exists
	if !strings.Contains(htmlContent, `addEventListener('keydown'`) && !strings.Contains(htmlContent, `addEventListener("keydown"`) {
		// Check if it's on the copy button specifically
		if !strings.Contains(htmlContent, `copyBtn.addEventListener`) {
			t.Error("Keyboard event listener for copy button not found")
		}
	}

	// Verify browser compatibility check exists
	if !strings.Contains(htmlContent, `isExecCommandAvailable`) {
		t.Error("execCommand availability check function not found")
	}
}
