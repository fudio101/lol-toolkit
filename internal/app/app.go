package app

import (
	"context"
	"fmt"

	"lol-toolkit/internal/config"
	"lol-toolkit/internal/lol"
)

// App struct holds the application state and dependencies
type App struct {
	ctx       context.Context
	config    *config.Config
	lolClient *lol.Client
}

// New creates a new App instance
func New() *App {
	return &App{}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		// Use default config if loading fails
		cfg = config.Default()
	}
	a.config = cfg

	// Initialize LoL client if API key is configured
	if cfg.RiotAPIKey != "" {
		client, err := lol.NewClient(cfg.RiotAPIKey, cfg.Region)
		if err == nil {
			a.lolClient = client
		}
	}
}

// Shutdown is called when the app is closing
func (a *App) Shutdown(ctx context.Context) {
	// Cleanup resources if needed
}

// GetConfig returns the current configuration (for frontend)
func (a *App) GetConfig() *config.Config {
	return a.config
}

// SetAPIKey updates the Riot API key
func (a *App) SetAPIKey(apiKey string) error {
	a.config.RiotAPIKey = apiKey

	// Reinitialize client with new key
	if apiKey != "" {
		client, err := lol.NewClient(apiKey, a.config.Region)
		if err != nil {
			return err
		}
		a.lolClient = client
	}

	// Save config
	return config.Save(a.config)
}

// SetRegion updates the region setting
func (a *App) SetRegion(region string) error {
	a.config.Region = region

	// Reinitialize client with new region
	if a.config.RiotAPIKey != "" {
		client, err := lol.NewClient(a.config.RiotAPIKey, region)
		if err != nil {
			return err
		}
		a.lolClient = client
	}

	// Save config
	return config.Save(a.config)
}

// IsConfigured returns true if the app has a valid API key
func (a *App) IsConfigured() bool {
	return a.config != nil && a.config.RiotAPIKey != ""
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
