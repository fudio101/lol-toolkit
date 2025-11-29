package lcu

import (
	"strings"
	"sync"
)

// ConnectionStatus represents the current LCU connection status.
type ConnectionStatus int

const (
	ConnectionStatusUnknown ConnectionStatus = iota
	ConnectionStatusConnected
	ConnectionStatusDisconnected
)

// ConnectionStatusManager manages the global LCU connection status.
type ConnectionStatusManager struct {
	status         ConnectionStatus
	onStatusChange func(connected bool)
	mu             sync.RWMutex
}

var globalStatusManager = &ConnectionStatusManager{
	status: ConnectionStatusUnknown,
}

// GetConnectionStatus returns the current connection status.
func GetConnectionStatus() ConnectionStatus {
	globalStatusManager.mu.RLock()
	defer globalStatusManager.mu.RUnlock()
	return globalStatusManager.status
}

// IsConnected returns true if the client is connected.
// Returns true if status is Unknown (hasn't been checked yet) to allow initial calls.
func IsConnected() bool {
	status := GetConnectionStatus()
	// If status is unknown, allow calls (will be checked on first GetLCUStatus)
	return status == ConnectionStatusConnected || status == ConnectionStatusUnknown
}

// SetConnectionStatus updates the connection status.
func SetConnectionStatus(connected bool) {
	newStatus := getStatusFromBool(connected)

	globalStatusManager.mu.Lock()
	oldStatus := globalStatusManager.status
	statusChanged := oldStatus != newStatus
	globalStatusManager.status = newStatus
	callback := globalStatusManager.onStatusChange
	globalStatusManager.mu.Unlock()

	if statusChanged && callback != nil {
		callback(connected)
	}
}

// getStatusFromBool converts a boolean to ConnectionStatus.
func getStatusFromBool(connected bool) ConnectionStatus {
	if connected {
		return ConnectionStatusConnected
	}
	return ConnectionStatusDisconnected
}

// SetOnStatusChange sets a callback to be called when connection status changes.
func SetOnStatusChange(callback func(connected bool)) {
	globalStatusManager.mu.Lock()
	defer globalStatusManager.mu.Unlock()
	globalStatusManager.onStatusChange = callback
}

// CheckAndUpdateStatus checks the connection status by calling GetConnectionInfo.
// This is used when a connection refused error is detected.
func CheckAndUpdateStatus() {
	info := GetConnectionInfo()
	SetConnectionStatus(info != nil)
}

// IsConnectionRefusedError checks if an error indicates the client connection was refused.
func IsConnectionRefusedError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())

	patterns := []struct {
		indicators []string
		all        bool // true = all indicators must be present, false = any
	}{
		{indicators: []string{"dial tcp", "actively refused"}, all: true},
		{indicators: []string{"connectex: no connection could be made"}, all: false},
		{indicators: []string{"connection refused"}, all: false},
		{indicators: []string{"league client not running"}, all: false},
		{indicators: []string{"failed to connect"}, all: false},
	}

	for _, pattern := range patterns {
		if matchesPattern(errStr, pattern.indicators, pattern.all) {
			// Special case: exclude 404 errors for "failed to connect"
			if strings.Contains(pattern.indicators[0], "failed to connect") && strings.Contains(errStr, "404") {
				continue
			}
			return true
		}
	}

	return false
}

// matchesPattern checks if error string matches a pattern.
func matchesPattern(errStr string, indicators []string, all bool) bool {
	if all {
		for _, indicator := range indicators {
			if !strings.Contains(errStr, indicator) {
				return false
			}
		}
		return true
	}

	for _, indicator := range indicators {
		if strings.Contains(errStr, indicator) {
			return true
		}
	}
	return false
}

// HandleConnectionError handles a connection error by checking status and updating it.
func HandleConnectionError(err error, endpoint string) bool {
	if endpoint == "GetLCUStatus" || !IsConnectionRefusedError(err) {
		return false
	}

	// Connection refused - HTTP server is not reachable
	SetConnectionStatus(false)
	return true
}
