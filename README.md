# Web Helper Framework

A Go-based web application framework that provides a three-panel interface for text input processing. The application automatically opens a browser window on startup and provides a clean, modern interface for standardizing and processing text input.

## Features

- **Automatic Browser Opening**: Launches the default browser automatically when the server starts
- **Three-Panel Interface**:
  - Input textarea for user text entry
  - Standardized textarea showing normalized input
  - Result textarea showing processed output
- **Text Standardization**: Removes line breaks and normalizes whitespace for natural text wrapping
- **Text Processing**: Converts text to uppercase (extensible for other transformations)
- **Copy to Clipboard**: One-click copying of processed results
- **Graceful Shutdown**: Handles Ctrl+C signals for clean server shutdown

## Project Structure

```
webhelper/
├── main.go              # Entry point, server setup, browser opening
├── handlers.go          # HTTP handlers for web routes
├── processor.go         # Input standardization and processing logic
├── processor_test.go    # Unit tests for processor functions
├── handlers_test.go     # Unit tests for HTTP handlers
├── main_test.go        # Unit tests for main functions
├── static/              # Static assets
│   ├── index.html      # Main web interface
│   └── style.css       # Modern CSS styling
├── go.mod              # Go module definition
└── README.md           # This file
```

## Installation

1. Ensure Go is installed (version 1.19 or later)
2. Clone or navigate to the project directory
3. Build the application:
   ```bash
   go build -o webhelper
   ```

## Usage

### Running the Application

```bash
./webhelper
```

Or using `go run`:

```bash
go run .
```

The server will start on port 8080 by default. To use a different port, set the `PORT` environment variable:

```bash
PORT=3000 ./webhelper
```

### Using the Web Interface

1. The browser will automatically open to `http://localhost:8080`
2. Enter text in the "Input Text" textarea
3. Click "Submit" to process the text
4. View the standardized input in the "Standardized Input" textarea
5. View the processed result in the "Result" textarea
6. Click "Copy to Clipboard" to copy the result

### Keyboard Shortcuts

- `Ctrl+Enter` (or `Cmd+Enter` on macOS) in the input textarea submits the form

## API Endpoints

### GET `/`
Serves the main HTML page.

### POST `/api/process`
Processes input text and returns standardized and processed results.

**Request Body:**
```json
{
  "input": "your text here"
}
```

**Response:**
```json
{
  "standardized": "normalized text",
  "result": "PROCESSED TEXT"
}
```

### GET `/static/*`
Serves static files (CSS, etc.).

## Testing

Run all tests:

```bash
go test -v ./...
```

Run specific test suites:

```bash
go test -v -run TestStandardizeInput
go test -v -run TestProcessInput
go test -v -run TestHandleProcess
```

## Implementation Details

### Text Standardization

The `StandardizeInput()` function:
- Trims leading and trailing whitespace
- Replaces all newlines (`\n`) and carriage returns (`\r`) with spaces
- Collapses multiple spaces and tabs into single spaces
- Allows text to wrap naturally in the textarea

### Text Processing

The `ProcessInput()` function currently converts text to uppercase. This function is designed to be easily extensible for other processing operations.

## Extending the Framework

### Adding New Processing Functions

Modify the `ProcessInput()` function in `processor.go` to add new transformations:

```go
func ProcessInput(input string) string {
    // Add your processing logic here
    return strings.ToUpper(input)
}
```

### Customizing the UI

- Modify `static/index.html` for HTML structure changes
- Modify `static/style.css` for styling changes

## Technical Stack

- **Backend**: Go standard library (`net/http`)
- **Frontend**: Vanilla JavaScript (no frameworks)
- **Styling**: Modern CSS with Flexbox
- **Testing**: Go testing package with `net/http/httptest`

## License

This project is provided as-is for educational and development purposes.
