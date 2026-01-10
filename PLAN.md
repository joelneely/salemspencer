# Web Helper Framework - Implementation Plan

## Architecture Overview

The application is a single Go program that:

1. Starts an HTTP server on a configurable port (default: 8080)
2. Automatically opens the default browser to the server URL
3. Serves a web interface with three text areas and a submit button
4. Processes form submissions through an API endpoint
5. Runs until interrupted (Ctrl+C)

## Project Structure

```
webhelper/
├── main.go              # Entry point, server setup, browser opening
├── handlers.go          # HTTP handlers for web routes
├── processor.go         # Input standardization and processing logic
├── static/              # Static assets
│   ├── index.html      # Main web interface
│   └── style.css       # Modern CSS styling
└── go.mod              # Go module definition
```

## Implementation Details

### main.go

- Initialize HTTP server on port 8080 (configurable via environment variable or flag)
- Register routes:
  - `GET /` - serves the main HTML page
  - `POST /api/process` - handles form submission
  - `/static/` - serves CSS files
- Open browser using OS-specific commands:
  - macOS: `open http://localhost:8080`
  - Linux: `xdg-open http://localhost:8080`
  - Windows: `cmd /c start http://localhost:8080`
- Graceful shutdown on SIGINT/SIGTERM

### handlers.go

- `handleIndex()` - serves `index.html` with proper content-type
- `handleProcess()` - JSON API endpoint that:
  - Accepts POST request with JSON body: `{"input": "user text"}`
  - Calls standardization function
  - Calls process function (uppercase)
  - Returns JSON: `{"standardized": "...", "result": "..."}`
- `handleStatic()` - serves CSS files

### processor.go

- `StandardizeInput(input string) string` - normalizes input:
  - Trim leading/trailing whitespace
  - Replace newlines with spaces to allow natural text wrapping
  - Normalize internal whitespace (collapse multiple spaces/tabs)
  - Return standardized string
- `ProcessInput(input string) string` - processes standardized input:
  - Initially: convert to uppercase
  - Designed to be easily extensible for future processing functions

### static/index.html

- Three textarea elements:
  1. Input textarea (5+ lines, 60+ chars width, editable)
  2. Standardized textarea (readonly, shows standardized input)
  3. Result textarea (readonly, shows processed result)
- Submit button (above or between textareas)
- Copy button below the result textarea
- Modern, clean layout using CSS Grid or Flexbox
- JavaScript to handle:
  - Form submission via fetch API
  - Real-time updates of standardized and result areas
  - Clipboard copy functionality using Clipboard API (`navigator.clipboard.writeText()`)

### static/style.css

- Modern, professional styling
- Responsive layout
- Clear visual separation between the three text areas
- Styled submit button with hover effects
- Styled copy button with hover effects and visual feedback
- Readonly textareas visually distinct from input textarea

## Data Flow

```
User Input → Submit Button → JavaScript fetch() → POST /api/process
                                                         ↓
                                    StandardizeInput() → ProcessInput()
                                                         ↓
                                    JSON Response → Update Textareas
```

## Technical Decisions

- **HTTP Server**: Standard library `net/http` (no external dependencies)
- **Browser Opening**: `os/exec` with OS detection
- **Frontend**: Vanilla JavaScript (no frameworks)
- **Styling**: Modern CSS with CSS Grid/Flexbox
- **API**: JSON-based REST endpoint
- **Port**: Default 8080, configurable via `PORT` environment variable

## Future Extensibility

The `ProcessInput()` function is designed as a pluggable interface, making it easy to:

- Add new processing functions
- Switch between different processors
- Extend the framework for more complex helper workflows

## Implementation Status

All planned features have been implemented:

- ✅ Project structure created
- ✅ Processor functions implemented (StandardizeInput, ProcessInput)
- ✅ HTTP handlers implemented (handleIndex, handleProcess, handleStatic)
- ✅ HTML interface created with three textareas and buttons
- ✅ CSS styling implemented
- ✅ Main server setup with browser opening
- ✅ Comprehensive unit tests added
- ✅ Copy to clipboard functionality added
- ✅ Text standardization updated to remove line breaks
