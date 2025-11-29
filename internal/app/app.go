package app

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"lol-toolkit/internal/config"
	"lol-toolkit/internal/lcu"
	"lol-toolkit/internal/lol"
)

// App holds application state and dependencies.
type App struct {
	ctx       context.Context
	config    *config.Config
	lolClient *lol.Client
}

// New creates a new App instance.
func New() *App {
	return &App{}
}

// Startup initializes the app when Wails starts.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	// Wire backend LCU API logging to frontend via Wails events.
	lcu.SetAPILogger(func(entry lcu.APILogEntry) {
		runtime.EventsEmit(a.ctx, "api-call", entry)
	})

	// Wire backend Riot (golio) API logging to frontend as well.
	lol.SetAPILogger(func(entry lol.APILogEntry) {
		runtime.EventsEmit(a.ctx, "api-call", entry)
	})

	a.loadConfig()
	a.initLolClient()
}

// Shutdown cleans up resources when the app closes.
func (a *App) Shutdown(_ context.Context) {
	// Cleanup if needed
}

// loadConfig loads the configuration.
func (a *App) loadConfig() {
	cfg, err := config.Load()
	if err != nil {
		cfg = config.Default()
	}
	a.config = cfg
}

// initLolClient initializes the LoL API client.
func (a *App) initLolClient() {
	if a.config.RiotAPIKey == "" {
		return
	}

	client, err := lol.NewClient(a.config.RiotAPIKey, a.config.Region)
	if err != nil {
		return
	}
	a.lolClient = client
}

// GetConfig returns the current configuration.
func (a *App) GetConfig() *config.Config {
	return a.config
}

// SetAPIKey updates the Riot API key and reinitializes the client.
func (a *App) SetAPIKey(apiKey string) error {
	a.config.RiotAPIKey = apiKey

	if apiKey != "" {
		client, err := lol.NewClient(apiKey, a.config.Region)
		if err != nil {
			return err
		}
		a.lolClient = client
	} else {
		a.lolClient = nil
	}

	return config.Save(a.config)
}

// SetRegion updates the region and reinitializes the client.
func (a *App) SetRegion(region string) error {
	a.config.Region = region

	if a.config.RiotAPIKey != "" {
		client, err := lol.NewClient(a.config.RiotAPIKey, region)
		if err != nil {
			return err
		}
		a.lolClient = client
	}

	return config.Save(a.config)
}

// IsConfigured returns true if the API key is set.
func (a *App) IsConfigured() bool {
	return a.config != nil && a.config.RiotAPIKey != ""
}
