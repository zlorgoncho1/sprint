package utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/zlorgoncho1/sprint/core"
)

// DictToJson takes a map[string]string and converts it into a JSON string.
// Useful for generating JSON from map data structures.
func DictToJson(__dict map[string]string) string {
	var items []string
	for k, v := range __dict {
		// Convert each key-value pair into a JSON string element.
		items = append(items, fmt.Sprintf(`"%s":"%s"`, k, v))
	}
	// Join all string elements to create a single JSON object.
	return "{" + strings.Join(items, ",") + "}"
}

// FormatStatusResponse formats an HTTP status line.
// Inputs include statusCode (e.g., 200), statusText (e.g., "OK"), and protocol (e.g., "HTTP/1.1").
func FormatStatusResponse(statusCode int, statusText string, protocol string) string {
	// Default to HTTP/1.1 if no protocol is specified.
	if protocol == "" {
		protocol = "HTTP/1.1"
	}
	// Returns the formatted status line.
	return fmt.Sprintf("%s %d %s", protocol, statusCode, statusText)
}

// DictToHTTPHeadersResponse takes a map of headers and formats them as HTTP headers.
func DictToHTTPHeadersResponse(headers map[string]string) string {
	var lines []string
	for k, v := range headers {
		// Each header is formatted as "Key: Value".
		lines = append(lines, fmt.Sprintf("%s: %s", k, v))
	}
	// Join all headers with a new line, mimicking HTTP header format.
	return strings.Join(lines, "\n")
}

// GetDefaultHeader generates and returns common default HTTP headers.
// It automatically calculates content length and sets the current date.
func GetDefaultHeader(content string, contentType core.ContentType) map[string]string {
	return map[string]string{
		"Content-Type":   string(contentType),
		"Content-Length": fmt.Sprintf("%d", len(content)),
		"Connection":     "close",
		"Date":           time.Now().Format(time.RFC1123Z),
	}
}

// FormatHTTPResponse constructs a complete HTTP response message.
// It combines the status line, headers, and content into a single byte slice.
func FormatHTTPResponse(status, headers, content string) []byte {
	var response strings.Builder
	response.WriteString(status + "\n")
	response.WriteString(headers + "\n\n")
	response.WriteString(content)
	return []byte(response.String())
}

// HandleHTML is a placeholder for future HTML response handling.
func HandleHTML(response *core.Response) {
	// Future implementation goes here.
}

// HandleJSON sets the ContentType of the response to "application/json"
// and converts the Content field of response to a JSON string.
func HandleJSON(response *core.Response) {
	response.ContentType = "application/json"
	jsonBytes, err := json.Marshal(response.Content)
	if err != nil {
		// If JSON marshaling fails, it panics. This might be replaced by better error handling in a real-world application.
		panic(err)
	}
	response.Content = string(jsonBytes)
}

// HandlePlainText simply sets the ContentType of the response to "text/plain".
// It assumes Content is already in string format.
func HandlePlainText(response *core.Response) {
	response.ContentType = "text/plain"
	response.Content = (response.Content).(string)
}

func JoinPaths(paths ...string) string {
	var buffer strings.Builder

	for i, path := range paths {
		// Trim slashes and then conditionally add one slash back.
		trimmedPath := strings.Trim(path, "/")
		if trimmedPath != "" {
			if i > 0 {
				buffer.WriteString("/")
			}
			buffer.WriteString(trimmedPath)
		}
	}
	return buffer.String()
}
