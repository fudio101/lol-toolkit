package app

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"lol-toolkit/internal/lcu"
)

// LCUStatus represents the League client connection status.
type LCUStatus struct {
	Connected bool   `json:"connected"`
	Port      string `json:"port,omitempty"`
	AuthToken string `json:"authToken,omitempty"`
	Error     string `json:"error,omitempty"`
}

// GetLCUStatus checks if the League client is running.
// This is treated as a backend API call for logging/telemetry purposes.
func (a *App) GetLCUStatus() *LCUStatus {
	start := time.Now()
	info := lcu.GetConnectionInfo()
	duration := time.Since(start)

	// Emit a synthetic backend API log so the Debug panel sees this call too.
	headers := buildLCUHeaders(info)

	var status *LCUStatus
	if info == nil {
		status = &LCUStatus{
			Connected: false,
			Error:     "League client not running",
		}
		statusJSON, _ := json.MarshalIndent(status, "", "  ")
		runtime.EventsEmit(a.ctx, "api-call", map[string]interface{}{
			"type":       "lcu",
			"method":     "GET",
			"endpoint":   "GetLCUStatus",
			"statusCode": http.StatusInternalServerError,
			"duration":   duration.Milliseconds(),
			"headers":    headers,
			"response":   string(statusJSON),
			"error":      "League client not running",
		})

		return status
	}

	status = &LCUStatus{
		Connected: true,
		Port:      info.Port,
		AuthToken: info.AuthToken,
	}
	statusJSON, _ := json.MarshalIndent(status, "", "  ")
	runtime.EventsEmit(a.ctx, "api-call", map[string]interface{}{
		"type":       "lcu",
		"method":     "GET",
		"endpoint":   "GetLCUStatus",
		"statusCode": http.StatusOK,
		"duration":   duration.Milliseconds(),
		"headers":    headers,
		"response":   string(statusJSON),
	})

	return status
}

// base64Encode encodes a string to base64.
func base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// GetCurrentSummoner returns the currently logged in summoner.
func (a *App) GetCurrentSummoner() (*lcu.CurrentSummoner, error) {
	start := time.Now()

	// Get connection info for headers
	info := lcu.GetConnectionInfo()
	headers := buildLCUHeaders(info)

	client, err := lcu.NewClient()
	if err != nil {
		// Log error when client initialization fails
		runtime.EventsEmit(a.ctx, "api-call", map[string]interface{}{
			"type":       "lcu",
			"method":     "GET",
			"endpoint":   "/lol-summoner/v1/current-summoner",
			"statusCode": http.StatusInternalServerError,
			"duration":   time.Since(start).Milliseconds(),
			"headers":    headers,
			"error":      err.Error(),
		})
		return nil, err
	}

	// GetCurrentSummoner will log its own errors and success
	return client.GetCurrentSummoner()
}

// buildLCUHeaders creates HTTP headers for LCU API requests.
func buildLCUHeaders(info *lcu.ConnectionInfo) map[string]string {
	headers := make(map[string]string)
	if info != nil && info.AuthToken != "" {
		headers["Authorization"] = fmt.Sprintf("Basic %s", base64Encode(fmt.Sprintf("riot:%s", info.AuthToken)))
	}
	headers["Accept"] = "application/json"
	return headers
}
