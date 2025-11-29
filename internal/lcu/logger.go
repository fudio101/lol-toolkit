package lcu

import (
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
// It executes the provided function, measures duration, and logs the result.
// Uses generics to maintain type safety - no type assertions needed.
func LoggedCall[T any](method, endpoint string, statusCode int, headers map[string]string, fn func() (T, error)) (T, error) {
	return logger.LoggedCall(apiType, method, endpoint, statusCode, headers, fn)
}

// LogSuccess logs a successful API call.
func LogSuccess(method, endpoint string, statusCode int, duration time.Duration, headers map[string]string, response string) {
	logger.LogSuccess(apiType, method, endpoint, statusCode, duration, headers, response)
}

// LogError logs a failed API call.
// If statusCode is 0, it will try to extract status code from the error, otherwise uses http.StatusInternalServerError as fallback.
func LogError(method, endpoint string, duration time.Duration, headers map[string]string, err error, statusCode ...int) {
	logger.LogError(apiType, method, endpoint, duration, headers, err, statusCode...)
}

// LogRequest logs an API call with full control over all fields.
func LogRequest(method, endpoint string, statusCode int, duration time.Duration, headers map[string]string, response string, err error) {
	logger.LogRequest(apiType, method, endpoint, statusCode, duration, headers, response, err)
}
