package app

import (
	"fmt"

	"lol-toolkit/internal/lol"
)

// GetChampionMastery gets mastery for a specific champion
func (a *App) GetChampionMastery(summonerID string, championID string) (*lol.ChampionMasteryInfo, error) {
	if a.lolClient == nil {
		return nil, fmt.Errorf("API client not initialized. Please set your API key first")
	}

	return a.lolClient.GetChampionMastery(summonerID, championID)
}

// GetAllChampionMasteries gets all champion masteries for a summoner
func (a *App) GetAllChampionMasteries(summonerID string) ([]*lol.ChampionMasteryInfo, error) {
	if a.lolClient == nil {
		return nil, fmt.Errorf("API client not initialized. Please set your API key first")
	}

	return a.lolClient.GetAllChampionMasteries(summonerID)
}

// GetTotalMasteryScore gets the total mastery score
func (a *App) GetTotalMasteryScore(summonerID string) (int, error) {
	if a.lolClient == nil {
		return 0, fmt.Errorf("API client not initialized. Please set your API key first")
	}

	return a.lolClient.GetTotalMasteryScore(summonerID)
}

