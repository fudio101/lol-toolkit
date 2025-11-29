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
func (a *App) GetLCUStatus() *LCUStatus {
	start := time.Now()
	info := lcu.GetConnectionInfo()
	duration := time.Since(start)

	lcu.SetConnectionStatus(info != nil)

	status := a.createStatus(info)
	a.emitStatusLog(status, duration, info)

	return status
}

// GetCurrentSummoner returns the currently logged in summoner.
func (a *App) GetCurrentSummoner() (*lcu.CurrentSummoner, error) {
	start := time.Now()

	client, err := lcu.NewClient()
	if err != nil {
		a.emitSummonerError(err, time.Since(start))
		return nil, err
	}

	return client.GetCurrentSummoner()
}

// createStatus creates an LCUStatus from connection info.
func (a *App) createStatus(info *lcu.ConnectionInfo) *LCUStatus {
	if info == nil {
		return &LCUStatus{
			Connected: false,
			Error:     "League client not running",
		}
	}

	return &LCUStatus{
		Connected: true,
		Port:      info.Port,
		AuthToken: info.AuthToken,
	}
}

// emitStatusLog emits a status check log event to the frontend.
func (a *App) emitStatusLog(status *LCUStatus, duration time.Duration, info *lcu.ConnectionInfo) {
	statusJSON, _ := json.MarshalIndent(status, "", "  ")
	headers := buildLCUHeaders(info)

	logEntry := map[string]interface{}{
		"type":     "lcu",
		"method":   "GET",
		"endpoint": "GetLCUStatus",
		"duration": duration.Milliseconds(),
		"headers":  headers,
		"response": string(statusJSON),
	}

	if status.Connected {
		logEntry["statusCode"] = http.StatusOK
	} else {
		logEntry["statusCode"] = http.StatusInternalServerError
		logEntry["error"] = status.Error
	}

	runtime.EventsEmit(a.ctx, "api-call", logEntry)
}

// emitSummonerError emits an error log event for summoner fetch failures.
func (a *App) emitSummonerError(err error, duration time.Duration) {
	info := lcu.GetConnectionInfo()
	headers := buildLCUHeaders(info)

	runtime.EventsEmit(a.ctx, "api-call", map[string]interface{}{
		"type":       "lcu",
		"method":     "GET",
		"endpoint":   "/lol-summoner/v1/current-summoner",
		"statusCode": http.StatusInternalServerError,
		"duration":   duration.Milliseconds(),
		"headers":    headers,
		"error":      err.Error(),
	})
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

// base64Encode encodes a string to base64.
func base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
