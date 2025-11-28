package app

import (
	"lol-toolkit/internal/lcu"
)

// LCUStatus represents the League client connection status.
type LCUStatus struct {
	Connected bool   `json:"connected"`
	Error     string `json:"error,omitempty"`
}

// GetLCUStatus checks if the League client is running.
func (a *App) GetLCUStatus() *LCUStatus {
	if lcu.IsClientRunning() {
		return &LCUStatus{Connected: true}
	}
	return &LCUStatus{
		Connected: false,
		Error:     "League client not running",
	}
}

// GetCurrentSummoner returns the currently logged in summoner.
func (a *App) GetCurrentSummoner() (*lcu.CurrentSummoner, error) {
	client, err := lcu.NewClient()
	if err != nil {
		return nil, err
	}
	return client.GetCurrentSummoner()
}
