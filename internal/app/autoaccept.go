package app

import (
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"lol-toolkit/internal/lcu"
)

var (
	autoAcceptService *lcu.AutoAcceptService
	appInstance       *App
)

// AutoAcceptConfig represents the auto-accept configuration.
type AutoAcceptConfig struct {
	Enabled    bool `json:"enabled"`
	AutoAccept bool `json:"autoAccept"`
}

// StartAutoAccept starts the auto-accept service.
func (a *App) StartAutoAccept(config AutoAcceptConfig) error {
	appInstance = a

	if err := a.stopExistingService(); err != nil {
		return err
	}

	client, err := lcu.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create LCU client: %w", err)
	}

	autoAcceptService = lcu.NewAutoAcceptService(client)
	autoAcceptService.SetAutoAccept(config.AutoAccept)
	autoAcceptService.SetOnStopped(a.createStoppedCallback())
	autoAcceptService.Start()

	return nil
}

// StopAutoAccept stops the auto-accept service.
func (a *App) StopAutoAccept() {
	a.stopExistingService()
}

// UpdateAutoAcceptConfig updates the auto-accept configuration.
func (a *App) UpdateAutoAcceptConfig(config AutoAcceptConfig) error {
	if autoAcceptService == nil {
		return fmt.Errorf("auto-accept service not started")
	}

	autoAcceptService.SetAutoAccept(config.AutoAccept)
	return nil
}

// IsAutoAcceptRunning returns true if the auto-accept service is currently running.
func (a *App) IsAutoAcceptRunning() bool {
	return autoAcceptService != nil
}

// stopExistingService stops the existing auto-accept service if running.
func (a *App) stopExistingService() error {
	if autoAcceptService != nil {
		autoAcceptService.Stop()
		autoAcceptService = nil
	}
	return nil
}

// createStoppedCallback creates a callback to notify frontend when service stops.
func (a *App) createStoppedCallback() func() {
	return func() {
		if appInstance != nil && appInstance.ctx != nil {
			runtime.EventsEmit(appInstance.ctx, "auto-accept-stopped", map[string]interface{}{
				"reason": "connection_error",
			})
		}
	}
}
