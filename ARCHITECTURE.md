# Architecture Documentation

## System Architecture

The Web Helper Framework is a simple, single-process web application built with Go's standard library. It follows a traditional request-response pattern with a clear separation of concerns.

## Component Overview

### 1. Main Application (`main.go`)

**Responsibilities:**
- Server initialization and configuration
- Route registration
- Browser opening logic
- Signal handling for graceful shutdown

**Key Functions:**
- `main()`: Application entry point, sets up HTTP server and handles lifecycle
- `openBrowser(url string)`: Cross-platform browser opening using OS-specific commands

**Design Decisions:**
- Uses environment variable `PORT` for configuration (default: 8080)
- Implements graceful shutdown with 5-second timeout
- Browser opening happens asynchronously after server start (500ms delay)

### 2. HTTP Handlers (`handlers.go`)

**Responsibilities:**
- Handle incoming HTTP requests
- Serve static files
- Process API requests

**Key Functions:**
- `handleIndex()`: Serves the main HTML page
- `handleProcess()`: Processes text input via JSON API
- `handleStatic()`: Serves static assets (CSS files)

**Design Decisions:**
- Uses standard library `net/http` handlers
- JSON-based API for data exchange
- File-based static asset serving (no embedded files)
- Proper HTTP status codes and error handling

### 3. Text Processing (`processor.go`)

**Responsibilities:**
- Text normalization and standardization
- Text transformation/processing

**Key Functions:**
- `StandardizeInput(input string) string`: Normalizes input text
- `ProcessInput(input string) string`: Transforms standardized text

**Design Decisions:**
- Pure functions (no side effects)
- Regex-based text processing for efficiency
- Designed for easy extensibility
- Current implementation:
  - Standardization: removes line breaks, collapses whitespace
  - Processing: converts to uppercase

### 4. Frontend (`static/index.html`, `static/style.css`)

**Responsibilities:**
- User interface rendering
- User interaction handling
- Visual feedback

**Key Components:**
- Three textareas (input, standardized, result)
- Submit button for processing
- Copy button for clipboard functionality
- JavaScript for API communication

**Design Decisions:**
- Vanilla JavaScript (no frameworks) for simplicity
- Modern CSS with gradient backgrounds and hover effects
- Responsive design using Flexbox
- Clipboard API for copy functionality
- Visual feedback for user actions

## Data Flow

```
┌─────────────┐
│   Browser   │
│  (Client)   │
└──────┬──────┘
       │
       │ HTTP GET /
       ▼
┌─────────────────┐
│  handleIndex()  │
│  (serves HTML)  │
└─────────────────┘
       │
       │ HTML + CSS + JS
       ▼
┌─────────────┐
│   Browser   │
│  (rendered) │
└──────┬──────┘
       │
       │ User enters text
       │ Clicks Submit
       │
       │ POST /api/process
       │ {"input": "text"}
       ▼
┌─────────────────┐
│ handleProcess() │
└────────┬────────┘
         │
         │ StandardizeInput()
         ▼
┌─────────────────┐
│ processor.go    │
└────────┬────────┘
         │
         │ ProcessInput()
         ▼
┌─────────────────┐
│ JSON Response   │
│ {standardized,  │
│  result}        │
└────────┬────────┘
         │
         │ Update textareas
         ▼
┌─────────────┐
│   Browser   │
│  (updated)   │
└─────────────┘
```

## Request/Response Patterns

### Index Page Request
```
GET /
→ 200 OK
→ Content-Type: text/html
→ <HTML content>
```

### Process Request
```
POST /api/process
Content-Type: application/json
Body: {"input": "hello\nworld"}

→ 200 OK
→ Content-Type: application/json
→ {"standardized": "hello world", "result": "HELLO WORLD"}
```

### Static File Request
```
GET /static/style.css
→ 200 OK
→ Content-Type: text/css
→ <CSS content>
```

## Error Handling

### HTTP Errors
- `400 Bad Request`: Invalid JSON or malformed request
- `404 Not Found`: File not found (static assets)
- `405 Method Not Allowed`: Wrong HTTP method
- `500 Internal Server Error`: Server-side errors

### Error Flow
- Handlers log errors using `log.Printf()`
- Return appropriate HTTP status codes
- Client-side JavaScript handles errors and shows alerts

## Testing Strategy

### Unit Tests
- **processor_test.go**: Tests for text processing functions
  - Edge cases (empty strings, whitespace-only)
  - Various input formats (newlines, tabs, mixed whitespace)
  - Unicode support

- **handlers_test.go**: Tests for HTTP handlers
  - Request validation
  - Response format verification
  - Error handling
  - Integration tests for full processing flow

- **main_test.go**: Tests for main functions
  - Browser opening function (limited due to OS constraints)

### Test Coverage
- All exported functions have test coverage
- Edge cases and error conditions are tested
- Integration tests verify end-to-end functionality

## Security Considerations

### Current Implementation
- No authentication/authorization (local use only)
- No input sanitization beyond standardization
- File serving from `static/` directory (no path traversal protection)
- No rate limiting

### Recommendations for Production
- Add input validation and sanitization
- Implement path traversal protection for static files
- Add rate limiting for API endpoints
- Consider HTTPS for production deployment
- Add CORS headers if needed for cross-origin access

## Performance Characteristics

### Current Implementation
- Single-threaded request handling (Go's HTTP server handles concurrency)
- In-memory processing (no database)
- File-based static asset serving
- No caching mechanisms

### Scalability
- Can handle multiple concurrent requests (Go's HTTP server)
- Stateless design allows horizontal scaling
- No shared state between requests
- Processing is CPU-bound (text operations)

## Deployment Considerations

### Current Setup
- Single binary deployment
- Requires `static/` directory at runtime
- Port configurable via environment variable
- No process management (use systemd/supervisor for production)

### Production Recommendations
- Embed static files using `embed` package
- Add health check endpoint
- Implement structured logging
- Add metrics/monitoring
- Use reverse proxy (nginx) for production
- Consider containerization (Docker)

## Extension Points

### Adding New Processors
1. Modify `ProcessInput()` function in `processor.go`
2. Add corresponding tests in `processor_test.go`
3. Update documentation

### Adding New Routes
1. Add handler function in `handlers.go`
2. Register route in `main.go`
3. Add tests in `handlers_test.go`

### Customizing UI
1. Modify `static/index.html` for structure
2. Modify `static/style.css` for styling
3. Add JavaScript functions as needed

## Dependencies

### Runtime Dependencies
- None (uses only Go standard library)

### Development Dependencies
- Go 1.19+ (for testing and build)

### External Commands (for browser opening)
- macOS: `open` (built-in)
- Linux: `xdg-open` (usually pre-installed)
- Windows: `cmd` (built-in)
