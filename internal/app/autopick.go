package app

import (
	"fmt"
	"lol-toolkit/internal/lcu"
)

var autoPickService *lcu.AutoPickService

// AutoPickConfig represents the auto-pick configuration.
type AutoPickConfig struct {
	Enabled      bool   `json:"enabled"`
	AutoAccept   bool   `json:"autoAccept"`
	AutoPick     bool   `json:"autoPick"`
	AutoLock     bool   `json:"autoLock"`
	ChampionID   int    `json:"championId"`
	ChampionName string `json:"championName,omitempty"`
}

// StartAutoPick starts the auto-pick service.
func (a *App) StartAutoPick(config AutoPickConfig) error {
	// Note: Auto-pick doesn't require Riot API key, only LCU connection

	// Stop existing service if running
	if autoPickService != nil {
		autoPickService.Stop()
		autoPickService = nil
	}

	client, err := lcu.NewClient()
	if err != nil {
		return fmt.Errorf("failed to create LCU client: %w", err)
	}

	autoPickService = lcu.NewAutoPickService(client)

	if config.ChampionName != "" {
		if err := autoPickService.SetChampion(config.ChampionName); err != nil {
			return err
		}
	} else if config.ChampionID > 0 {
		autoPickService.SetChampionID(config.ChampionID)
	}

	autoPickService.SetAutoAccept(config.AutoAccept)
	autoPickService.SetAutoPick(config.AutoPick)
	autoPickService.SetAutoLock(config.AutoLock)

	autoPickService.Start()
	return nil
}

// StopAutoPick stops the auto-pick service.
func (a *App) StopAutoPick() {
	if autoPickService != nil {
		autoPickService.Stop()
		autoPickService = nil
	}
}

// UpdateAutoPickConfig updates the auto-pick configuration.
func (a *App) UpdateAutoPickConfig(config AutoPickConfig) error {
	if autoPickService == nil {
		return fmt.Errorf("auto-pick service not started")
	}

	if config.ChampionName != "" {
		if err := autoPickService.SetChampion(config.ChampionName); err != nil {
			return err
		}
	} else if config.ChampionID > 0 {
		autoPickService.SetChampionID(config.ChampionID)
	}

	autoPickService.SetAutoAccept(config.AutoAccept)
	autoPickService.SetAutoPick(config.AutoPick)
	autoPickService.SetAutoLock(config.AutoLock)

	return nil
}

// GetOwnedChampions returns the list of owned champions.
func (a *App) GetOwnedChampions() ([]lcu.Champion, error) {
	client, err := lcu.NewClient()
	if err != nil {
		return nil, err
	}

	return client.GetOwnedChampions()
}
