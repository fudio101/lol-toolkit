package app

import (
	"fmt"

	"lol-toolkit/internal/lol"
)

// GetRankedStats gets ranked stats for a summoner
func (a *App) GetRankedStats(summonerID string) ([]*lol.RankedInfo, error) {
	if a.lolClient == nil {
		return nil, fmt.Errorf("API client not initialized. Please set your API key first")
	}

	return a.lolClient.GetRankedStats(summonerID)
}

// GetChallengers gets the challenger leaderboard
func (a *App) GetChallengers(queueType string) (*lol.LeagueListInfo, error) {
	if a.lolClient == nil {
		return nil, fmt.Errorf("API client not initialized. Please set your API key first")
	}

	return a.lolClient.GetChallengers(queueType)
}

// GetGrandmasters gets the grandmaster leaderboard
func (a *App) GetGrandmasters(queueType string) (*lol.LeagueListInfo, error) {
	if a.lolClient == nil {
		return nil, fmt.Errorf("API client not initialized. Please set your API key first")
	}

	return a.lolClient.GetGrandmasters(queueType)
}

// GetMasters gets the master leaderboard
func (a *App) GetMasters(queueType string) (*lol.LeagueListInfo, error) {
	if a.lolClient == nil {
		return nil, fmt.Errorf("API client not initialized. Please set your API key first")
	}

	return a.lolClient.GetMasters(queueType)
}

