package logger

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const (
	// minHTTPStatusCode is the minimum valid HTTP status code (100)
	minHTTPStatusCode = 100
	// maxHTTPStatusCode is the validation upper bound (600) - not a standard HTTP status code
	maxHTTPStatusCode = 600
)

// APILogEntry represents a single API call for logging/telemetry.
// Duration is expressed in milliseconds for easy display in the frontend.
type APILogEntry struct {
	Type       string            `json:"type"`               // "lcu" or "riot"
	Method     string            `json:"method"`             // GET, POST, etc.
	Endpoint   string            `json:"endpoint"`           // API endpoint path
	StatusCode int               `json:"statusCode"`         // HTTP status code
	Duration   int64             `json:"duration"`           // request duration in ms
	Headers    map[string]string `json:"headers,omitempty"`  // request headers
	Response   string            `json:"response,omitempty"` // optional response body (JSON)
	Error      string            `json:"error,omitempty"`
}

// apiLogger is an optional callback set by the app to receive API logs.
var apiLogger func(entry APILogEntry)

// SetAPILogger registers a callback to receive API logs.
func SetAPILogger(logger func(entry APILogEntry)) {
	apiLogger = logger
}

// logAPICall sends an API log entry to the registered logger if present.
func logAPICall(entry APILogEntry) {
	if apiLogger != nil {
		apiLogger(entry)
	}
}

// extractStatusCode tries to extract HTTP status code from an error message.
// Returns the status code if found, otherwise returns http.StatusInternalServerError as fallback.
func extractStatusCode(err error) int {
	if err == nil {
		return http.StatusInternalServerError
	}

	errMsg := err.Error()

	// Try to match common patterns:
	// - "status code: 404"
	// - "404"
	// - "unexpected status code: 404"
	// - "HTTP 404"
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`status code[:\s]+(\d{3})`),
		regexp.MustCompile(`HTTP[:\s]+(\d{3})`),
		regexp.MustCompile(`\b(\d{3})\b`), // Any 3-digit number
	}

	for _, pattern := range patterns {
		matches := pattern.FindStringSubmatch(errMsg)
		if len(matches) >= 2 {
			if code, parseErr := strconv.Atoi(matches[1]); parseErr == nil {
				// Validate it's a valid HTTP status code (100-599)
				if code >= minHTTPStatusCode && code < maxHTTPStatusCode {
					return code
				}
			}
		}
	}

	// Fallback to http.StatusInternalServerError if no status code found
	return http.StatusInternalServerError
}

// LoggedCall wraps an API call with automatic timing and logging.
// It executes the provided function, measures duration, and logs the result.
// Uses generics to maintain type safety - no type assertions needed.
func LoggedCall[T any](apiType, method, endpoint string, statusCode int, headers map[string]string, fn func() (T, error)) (T, error) {
	start := time.Now()
	result, err := fn()
	duration := time.Since(start)

	var response string
	if err == nil {
		// Try to marshal the result to JSON
		if jsonData, marshalErr := json.MarshalIndent(result, "", "  "); marshalErr == nil {
			response = string(jsonData)
		}
	}

	// Build log entry
	entry := APILogEntry{
		Type:       apiType,
		Method:     method,
		Endpoint:   endpoint,
		StatusCode: statusCode,
		Duration:   duration.Milliseconds(),
		Headers:    headers,
		Response:   response,
	}

	if err != nil {
		// Try to extract status code from error, fallback to http.StatusInternalServerError
		entry.StatusCode = extractStatusCode(err)
		entry.Error = err.Error()
		entry.Response = "" // Clear response on error
	}

	logAPICall(entry)
	return result, err
}

// LogSuccess logs a successful API call.
func LogSuccess(apiType, method, endpoint string, statusCode int, duration time.Duration, headers map[string]string, response string) {
	logAPICall(APILogEntry{
		Type:       apiType,
		Method:     method,
		Endpoint:   endpoint,
		StatusCode: statusCode,
		Duration:   duration.Milliseconds(),
		Headers:    headers,
		Response:   response,
	})
}

// LogError logs a failed API call.
// If statusCode is 0, it will try to extract status code from the error, otherwise uses http.StatusInternalServerError as fallback.
func LogError(apiType, method, endpoint string, duration time.Duration, headers map[string]string, err error, statusCode ...int) {
	if err == nil {
		return
	}

	code := http.StatusInternalServerError
	if len(statusCode) > 0 && statusCode[0] > 0 {
		code = statusCode[0]
	} else {
		// Try to extract from error message
		code = extractStatusCode(err)
	}

	logAPICall(APILogEntry{
		Type:       apiType,
		Method:     method,
		Endpoint:   endpoint,
		StatusCode: code,
		Duration:   duration.Milliseconds(),
		Headers:    headers,
		Error:      err.Error(),
	})
}

// LogRequest logs an API call with full control over all fields.
func LogRequest(apiType, method, endpoint string, statusCode int, duration time.Duration, headers map[string]string, response string, err error) {
	entry := APILogEntry{
		Type:       apiType,
		Method:     method,
		Endpoint:   endpoint,
		StatusCode: statusCode,
		Duration:   duration.Milliseconds(),
		Headers:    headers,
		Response:   response,
	}
	if err != nil {
		entry.Error = err.Error()
	}
	logAPICall(entry)
}
