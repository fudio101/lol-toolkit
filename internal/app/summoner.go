package app

import (
	"fmt"

	"lol-toolkit/internal/lol"
)

// SearchSummoner searches for a summoner by Riot ID (gameName#tagLine)
func (a *App) SearchSummoner(riotID string) (*lol.SummonerInfo, error) {
	if a.lolClient == nil {
		return nil, fmt.Errorf("API client not initialized. Please set your API key first")
	}

	return a.lolClient.SearchByRiotID(riotID)
}

// GetSummonerByPUUID searches for a summoner by PUUID
func (a *App) GetSummonerByPUUID(puuid string) (*lol.SummonerInfo, error) {
	if a.lolClient == nil {
		return nil, fmt.Errorf("API client not initialized. Please set your API key first")
	}

	return a.lolClient.GetSummonerByPUUID(puuid)
}

// GetSummonerByID searches for a summoner by summoner ID
func (a *App) GetSummonerByID(summonerID string) (*lol.SummonerInfo, error) {
	if a.lolClient == nil {
		return nil, fmt.Errorf("API client not initialized. Please set your API key first")
	}

	return a.lolClient.GetSummonerByID(summonerID)
}
