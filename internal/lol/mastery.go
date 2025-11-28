package lol

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
	c.waitForRateLimit()
	mastery, err := c.golio.Riot.LoL.ChampionMastery.Get(summonerID, championID)
	if err != nil {
		return nil, err
	}

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
	c.waitForRateLimit()
	masteries, err := c.golio.Riot.LoL.ChampionMastery.List(summonerID)
	if err != nil {
		return nil, err
	}

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
	c.waitForRateLimit()
	return c.golio.Riot.LoL.ChampionMastery.GetTotal(summonerID)
}
