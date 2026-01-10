package main

import (
	"regexp"
	"strings"
)

// StandardizeInput normalizes the input text by:
// - Trimming leading/trailing whitespace
// - Replacing newlines with spaces to allow natural text wrapping
// - Normalizing internal whitespace (collapsing multiple spaces/tabs)
func StandardizeInput(input string) string {
	// Trim leading and trailing whitespace
	standardized := strings.TrimSpace(input)

	// Replace all newlines, carriage returns, and Unicode line/paragraph separators with spaces
	// First replace standard newlines (\r, \n)
	newlineRegex := regexp.MustCompile(`[\r\n]+`)
	standardized = newlineRegex.ReplaceAllString(standardized, " ")
	// Then replace Unicode line separator (U+2028) and paragraph separator (U+2029)
	standardized = strings.ReplaceAll(standardized, "\u2028", " ")
	standardized = strings.ReplaceAll(standardized, "\u2029", " ")

	// Collapse multiple spaces and tabs into single space
	spaceRegex := regexp.MustCompile(`[ \t]+`)
	standardized = spaceRegex.ReplaceAllString(standardized, " ")

	return standardized
}

// ProcessInput processes the standardized input.
// Initially converts text to uppercase.
// Designed to be easily extensible for future processing functions.
func ProcessInput(input string) string {
	return strings.ToUpper(input)
}
