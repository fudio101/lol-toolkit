package lcu

import (
	"fmt"
	"time"

	"lol-toolkit/internal/logger"
)

const apiType = "lcu"

// APILogEntry is an alias for logger.APILogEntry for backward compatibility.
type APILogEntry = logger.APILogEntry

// SetAPILogger registers a callback to receive LCU API logs.
func SetAPILogger(loggerFunc func(entry APILogEntry)) {
	logger.SetAPILogger(loggerFunc)
}

// LoggedCall wraps an API call with automatic timing and logging.
func LoggedCall[T any](method, endpoint string, statusCode int, headers map[string]string, fn func() (T, error)) (T, error) {
	if shouldBlockCall(endpoint) {
		return handleBlockedCall[T](method, endpoint, headers)
	}

	result, err := logger.LoggedCall(apiType, method, endpoint, statusCode, headers, fn)

	if err != nil {
		HandleConnectionError(err, endpoint)
	}

	return result, err
}

// shouldBlockCall checks if the call should be blocked due to disconnection.
func shouldBlockCall(endpoint string) bool {
	return endpoint != "GetLCUStatus" && !IsConnected()
}

// handleBlockedCall handles a blocked API call.
func handleBlockedCall[T any](method, endpoint string, headers map[string]string) (T, error) {
	var zero T
	err := fmt.Errorf("league client not connected")
	logger.LogError(apiType, method, endpoint, 0, headers, err)
	return zero, err
}

// LogSuccess logs a successful API call.
func LogSuccess(method, endpoint string, statusCode int, duration time.Duration, headers map[string]string, response string) {
	logger.LogSuccess(apiType, method, endpoint, statusCode, duration, headers, response)
}

// LogError logs a failed API call.
func LogError(method, endpoint string, duration time.Duration, headers map[string]string, err error, statusCode ...int) {
	logger.LogError(apiType, method, endpoint, duration, headers, err, statusCode...)
}

// LogRequest logs an API call with full control over all fields.
func LogRequest(method, endpoint string, statusCode int, duration time.Duration, headers map[string]string, response string, err error) {
	logger.LogRequest(apiType, method, endpoint, statusCode, duration, headers, response, err)
}
