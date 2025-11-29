package lol

import (
	"net/http"
	"time"
)

// ChampionMasteryInfo represents champion mastery data for the frontend
type ChampionMasteryInfo struct {
	ChampionID                   int    `json:"championId"`
	ChampionLevel                int    `json:"championLevel"`
	ChampionPoints               int    `json:"championPoints"`
	ChampionPointsSinceLastLevel int    `json:"championPointsSinceLastLevel"`
	ChampionPointsUntilNextLevel int    `json:"championPointsUntilNextLevel"`
	ChestGranted                 bool   `json:"chestGranted"`
	LastPlayTime                 int    `json:"lastPlayTime"`
	TokensEarned                 int    `json:"tokensEarned"`
	SummonerID                   string `json:"summonerId"`
}

// GetChampionMastery fetches champion mastery for a summoner and champion
func (c *Client) GetChampionMastery(summonerID string, championID string) (*ChampionMasteryInfo, error) {
	start := time.Now()
	mastery, err := c.golio.Riot.LoL.ChampionMastery.Get(summonerID, championID)
	duration := time.Since(start)
	if err != nil {
		LogError("GET", "champion-mastery/get", duration, c.getHeaders(), err)
		return nil, err
	}
	LogSuccess("GET", "champion-mastery/get", http.StatusOK, duration, c.getHeaders(), c.marshalResponse(mastery))

	return &ChampionMasteryInfo{
		ChampionID:                   mastery.ChampionID,
		ChampionLevel:                mastery.ChampionLevel,
		ChampionPoints:               mastery.ChampionPoints,
		ChampionPointsSinceLastLevel: mastery.ChampionPointsSinceLastLevel,
		ChampionPointsUntilNextLevel: mastery.ChampionPointsUntilNextLevel,
		ChestGranted:                 mastery.ChestGranted,
		LastPlayTime:                 mastery.LastPlayTime,
		TokensEarned:                 mastery.TokensEarned,
		SummonerID:                   mastery.SummonerID,
	}, nil
}

// GetAllChampionMasteries fetches all champion masteries for a summoner
func (c *Client) GetAllChampionMasteries(summonerID string) ([]*ChampionMasteryInfo, error) {
	start := time.Now()
	masteries, err := c.golio.Riot.LoL.ChampionMastery.List(summonerID)
	duration := time.Since(start)
	if err != nil {
		LogError("GET", "champion-mastery/list", duration, c.getHeaders(), err)
		return nil, err
	}
	LogSuccess("GET", "champion-mastery/list", http.StatusOK, duration, c.getHeaders(), c.marshalResponse(masteries))

	result := make([]*ChampionMasteryInfo, len(masteries))
	for i, m := range masteries {
		result[i] = &ChampionMasteryInfo{
			ChampionID:                   m.ChampionID,
			ChampionLevel:                m.ChampionLevel,
			ChampionPoints:               m.ChampionPoints,
			ChampionPointsSinceLastLevel: m.ChampionPointsSinceLastLevel,
			ChampionPointsUntilNextLevel: m.ChampionPointsUntilNextLevel,
			ChestGranted:                 m.ChestGranted,
			LastPlayTime:                 m.LastPlayTime,
			TokensEarned:                 m.TokensEarned,
			SummonerID:                   m.SummonerID,
		}
	}

	return result, nil
}

// GetTotalMasteryScore fetches the total mastery score for a summoner
func (c *Client) GetTotalMasteryScore(summonerID string) (int, error) {
	start := time.Now()
	total, err := c.golio.Riot.LoL.ChampionMastery.GetTotal(summonerID)
	duration := time.Since(start)
	if err != nil {
		LogError("GET", "champion-mastery/total", duration, c.getHeaders(), err)
		return 0, err
	}
	LogSuccess("GET", "champion-mastery/total", http.StatusOK, duration, c.getHeaders(), c.marshalResponse(total))
	return total, nil
}
