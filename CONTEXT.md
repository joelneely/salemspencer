# Project Context and Current State

## Project Overview

The Web Helper Framework is a Go-based web application that provides a simple interface for text processing. It was built as a foundation for creating web-based helper tools with a clean, extensible architecture.

## Current Implementation Status

### ✅ Completed Features

1. **Core Functionality**
   - HTTP server with automatic browser opening
   - Server readiness check before opening browser (prevents "localhost can't be reached" errors)
   - Three-panel web interface (input, standardized, result)
   - Text standardization (removes all line breaks including Unicode separators, normalizes whitespace)
   - Text processing (uppercase conversion)
   - Copy to clipboard functionality
   - Robust static file path resolution (works regardless of execution context)

2. **Code Quality**
   - Comprehensive unit tests (35+ test cases)
   - All tests passing
   - No linter errors
   - Clean code structure
   - Test-driven development (TDD) methodology

3. **Documentation**
   - README.md with usage instructions
   - PLAN.md with implementation plan
   - ARCHITECTURE.md with technical details
   - CONTEXT.md (this file)

### Implementation Details

#### Text Standardization Behavior

The `StandardizeInput()` function:
- Trims leading/trailing whitespace
- **Replaces all newlines with spaces** (allows natural text wrapping)
  - Handles standard newlines (`\r`, `\n`)
  - Handles Unicode line separator (U+2028)
  - Handles Unicode paragraph separator (U+2029)
- Collapses multiple spaces and tabs into single spaces
- Returns a single-line string suitable for textarea display

Example:
```
Input:  "hello\n\nworld"
Output: "hello world"

Input:  "hello\u2028world\u2029test"
Output: "hello world test"
```

#### Textarea Display Behavior

The standardized and result textareas are configured to:
- Wrap text naturally based on width (`wrap="soft"` attribute)
- Use CSS properties (`white-space: normal`, `word-wrap: break-word`) to ensure proper wrapping
- Display text without preserving original line breaks
- Only break to new lines when necessary due to text length

#### Text Processing Behavior

The `ProcessInput()` function:
- Currently converts text to uppercase
- Designed to be easily extensible
- Operates on standardized input

Example:
```
Input:  "hello world"
Output: "HELLO WORLD"
```

### File Structure

```
webhelper/
├── .gitignore              # Git ignore rules
├── README.md               # User documentation
├── PLAN.md                 # Implementation plan
├── ARCHITECTURE.md         # Technical architecture
├── CONTEXT.md              # This file - project context
├── go.mod                  # Go module definition
├── main.go                 # Entry point (server setup, browser opening, readiness check)
├── main_test.go            # Tests for main.go (browser opening, server readiness)
├── handlers.go             # HTTP handlers (includes getStaticDir function)
├── handlers_test.go        # Tests for handlers.go (includes getStaticDir test)
├── processor.go            # Text processing (Unicode separator support)
├── processor_test.go       # Tests for processor.go (Unicode separator tests)
└── static/
    ├── index.html          # Web interface (with wrap="soft" attributes)
    └── style.css           # Styling (text wrapping properties)
```

### Test Coverage

**processor_test.go:**
- `TestStandardizeInput`: 15 test cases
  - Empty strings
  - Simple text
  - Whitespace handling
  - Newline removal (`\r`, `\n`)
  - Unicode line separator (U+2028)
  - Unicode paragraph separator (U+2029)
  - Mixed newlines and Unicode separators
  - Mixed whitespace
  - Unicode support

- `TestProcessInput`: 7 test cases
  - Case conversion
  - Special characters
  - Numbers
  - Unicode

**handlers_test.go:**
- `TestHandleIndex`: GET/POST method validation
- `TestHandleProcess`: 6 test cases
  - Valid requests
  - Invalid JSON
  - Empty input
  - Multiline input
  - Wrong HTTP methods
- `TestHandleStatic`: Static file serving
- `TestHandleProcessIntegration`: End-to-end flow tests
  - Includes Unicode separator tests
- `TestGetStaticDir`: Static directory path resolution
  - Verifies static directory can be found
  - Checks for required files (index.html, style.css)

**main_test.go:**
- `TestOpenBrowser`: Browser opening function tests
- `TestOpenBrowserUnsupportedOS`: OS compatibility tests
- `TestWaitForServerReady`: Server readiness verification
  - Tests server ready immediately
  - Tests timeout scenarios
  - Tests various HTTP status codes

### Recent Improvements

1. **Unicode Separator Support**
   - Added support for Unicode line separator (U+2028) and paragraph separator (U+2029)
   - `StandardizeInput()` now handles all newline types comprehensively
   - Added test cases for Unicode separators

2. **Textarea Wrapping**
   - Updated CSS with `white-space: normal`, `word-wrap: break-word`, `overflow-wrap: break-word`
   - Added `wrap="soft"` attribute to textareas
   - Ensures text wraps naturally without preserving original line breaks

3. **Static File Path Resolution**
   - Implemented `getStaticDir()` function for robust path resolution
   - Works regardless of execution context (development, production, subdirectories)
   - Added comprehensive test coverage

4. **Server Readiness Check**
   - Added `waitForServerReady()` function to verify server is accepting connections
   - Browser opens only after server is ready (prevents connection errors)
   - Polls server with HTTP requests before opening browser
   - Added test coverage for readiness scenarios

### Git History

1. **Initial Commit** (`b05fdc7`)
   - Implement web helper framework with comprehensive unit tests
   - All source files and static assets

2. **Second Commit** (`316cc6a`)
   - Add .gitignore for Go project

### Key Design Decisions

1. **No External Dependencies**
   - Uses only Go standard library
   - Easier to build and deploy
   - Smaller binary size

2. **Vanilla JavaScript**
   - No frontend frameworks
   - Simpler to understand and modify
   - Faster page loads

3. **File-Based Static Assets**
   - Static files served from filesystem
   - `getStaticDir()` function finds static directory reliably:
     - Checks current working directory first
     - Falls back to executable directory
     - Handles binaries in subdirectories (e.g., `bin/`)
     - Works with `go run` during development
   - Easy to modify without recompiling
   - Could be embedded later if needed

4. **Pure Functions for Processing**
   - `StandardizeInput()` and `ProcessInput()` are pure functions
   - Easy to test
   - No side effects

5. **JSON API**
   - Simple request/response format
   - Easy to extend
   - Standard web API pattern

### Known Limitations

1. **No Authentication**
   - Designed for local use only
   - Not suitable for public deployment without security

2. **No Input Validation**
   - Processes any text input
   - No size limits
   - No content filtering

3. **Static File Serving**
   - No path traversal protection
   - Serves files from `static/` directory

4. **Browser Opening**
   - Server readiness check ensures browser opens only after server is accepting connections
   - `waitForServerReady()` polls server with HTTP requests before opening browser
   - Prevents "localhost can't be reached" errors
   - May fail on systems without browser commands
   - Fallback: logs warning and provides manual URL if browser opening fails

### Future Enhancement Ideas

1. **Processing Functions**
   - Add more processing options (lowercase, title case, etc.)
   - Allow user to select processing type
   - Chain multiple processors

2. **UI Improvements**
   - Add processing options selector
   - Show processing history
   - Add undo/redo functionality

3. **Performance**
   - Add request caching
   - Optimize regex operations
   - Consider streaming for large inputs

4. **Features**
   - Save/load functionality
   - Export results to file
   - Multiple input formats support

5. **Testing**
   - Add integration tests with real browser
   - Add performance benchmarks
   - Add fuzzing tests

### Development Workflow

**Test-Driven Development (TDD) Approach:**
- This project follows a **test-driven development** methodology
- Changes are **described via tests first** before implementation
- Proposed changes are **verified by running tests** throughout development
- Write tests that describe the desired behavior, then implement to make tests pass

1. **Making Changes (TDD Process)**
   - **Write tests first** that describe the desired behavior/change
   - Run `go test -v ./...` to see tests fail (red phase)
   - Implement the minimal code to make tests pass (green phase)
   - Refactor if needed while keeping tests passing
   - Run `go test -v ./...` to verify all tests pass
   - Run `go build` to check compilation
   - Test manually by running `./webhelper` if applicable

2. **Adding Tests**
   - Add test cases to appropriate `*_test.go` file
   - Follow existing test patterns
   - Ensure all edge cases are covered
   - Tests should clearly describe the expected behavior

3. **Extending Processing**
   - Write tests in `processor_test.go` that describe the new processing behavior
   - Run tests to verify they fail initially
   - Modify `ProcessInput()` in `processor.go` to implement the feature
   - Run tests again to verify they pass
   - Update documentation

4. **Committing Changes**
   - **MANDATORY**: Run `go test -v ./...` before committing any changes
   - **DO NOT COMMIT** if any tests fail
   - Following TDD, tests should already be passing (green phase)
   - Verify all tests pass before creating a commit
   - This ensures code quality and prevents regressions

### Environment Setup

**Required:**
- Go 1.19 or later
- Web browser (for testing)

**Optional:**
- Git (for version control)
- IDE with Go support (for development)

### Build and Run

```bash
# Build
go build -o webhelper

# Run
./webhelper

# Run tests
go test -v ./...

# Run with custom port
PORT=3000 ./webhelper
```

### Troubleshooting

**Browser doesn't open:**
- Check that browser command exists (`open`, `xdg-open`, or `cmd`)
- The server now waits for readiness before opening browser (up to 5 seconds)
- If browser still doesn't open, manually navigate to `http://localhost:8080`
- Check server logs for any startup errors

**Static files not found:**
- The `getStaticDir()` function automatically finds the static directory
- Ensure you're running from the project root or have the `static/` directory accessible
- The function checks multiple locations automatically

**Port already in use:**
- Change port using `PORT` environment variable
- Or kill the process using the port

**Tests fail:**
- Ensure `static/` directory exists
- Check that all files are present
- Run `go mod tidy` to ensure dependencies are correct

### Project Goals

The project was designed to be:
- **Simple**: Easy to understand and modify
- **Extensible**: Easy to add new features
- **Testable**: Well-tested with comprehensive coverage
- **Documented**: Clear documentation for users and developers

This foundation can be extended to create more complex web-based helper tools while maintaining the same clean architecture.
