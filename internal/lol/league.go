package lol

import (
	"github.com/KnutZuidema/golio/riot/lol"
)

// RankedInfo represents ranked league data for the frontend
type RankedInfo struct {
	QueueType    string `json:"queueType"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	LeaguePoints int    `json:"leaguePoints"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
	HotStreak    bool   `json:"hotStreak"`
	Veteran      bool   `json:"veteran"`
	FreshBlood   bool   `json:"freshBlood"`
	Inactive     bool   `json:"inactive"`
	SummonerID   string `json:"summonerId"`
	SummonerName string `json:"summonerName"`
	PUUID        string `json:"puuid"`
}

// GetRankedStats fetches all ranked entries for a summoner
func (c *Client) GetRankedStats(summonerID string) ([]*RankedInfo, error) {
	entries, err := c.golio.Riot.LoL.League.ListBySummoner(summonerID)
	if err != nil {
		return nil, err
	}

	result := make([]*RankedInfo, len(entries))
	for i, e := range entries {
		result[i] = &RankedInfo{
			QueueType:    e.QueueType,
			Tier:         e.Tier,
			Rank:         e.Rank,
			LeaguePoints: e.LeaguePoints,
			Wins:         e.Wins,
			Losses:       e.Losses,
			HotStreak:    e.HotStreak,
			Veteran:      e.Veteran,
			FreshBlood:   e.FreshBlood,
			Inactive:     e.Inactive,
			SummonerID:   e.SummonerID,
			SummonerName: e.SummonerName,
			PUUID:        e.PUUID,
		}
	}

	return result, nil
}

// GetChallengers fetches the challenger league for a queue
func (c *Client) GetChallengers(queueType string) (*LeagueListInfo, error) {
	league, err := c.golio.Riot.LoL.League.GetChallenger(lol.QueueRankedSolo)
	if queueType == QueueRankedFlex {
		league, err = c.golio.Riot.LoL.League.GetChallenger(lol.QueueRankedFlex)
	}
	if err != nil {
		return nil, err
	}

	return toLeagueListInfo(league), nil
}

// GetGrandmasters fetches the grandmaster league for a queue
func (c *Client) GetGrandmasters(queueType string) (*LeagueListInfo, error) {
	league, err := c.golio.Riot.LoL.League.GetGrandmaster(lol.QueueRankedSolo)
	if queueType == QueueRankedFlex {
		league, err = c.golio.Riot.LoL.League.GetGrandmaster(lol.QueueRankedFlex)
	}
	if err != nil {
		return nil, err
	}

	return toLeagueListInfo(league), nil
}

// GetMasters fetches the master league for a queue
func (c *Client) GetMasters(queueType string) (*LeagueListInfo, error) {
	league, err := c.golio.Riot.LoL.League.GetMaster(lol.QueueRankedSolo)
	if queueType == QueueRankedFlex {
		league, err = c.golio.Riot.LoL.League.GetMaster(lol.QueueRankedFlex)
	}
	if err != nil {
		return nil, err
	}

	return toLeagueListInfo(league), nil
}

// LeagueListInfo represents a league list for the frontend
type LeagueListInfo struct {
	Tier     string        `json:"tier"`
	LeagueID string        `json:"leagueId"`
	Queue    string        `json:"queue"`
	Name     string        `json:"name"`
	Entries  []*RankedInfo `json:"entries"`
}

func toLeagueListInfo(l *lol.LeagueList) *LeagueListInfo {
	entries := make([]*RankedInfo, len(l.Entries))
	for i, e := range l.Entries {
		entries[i] = &RankedInfo{
			QueueType:    e.QueueType,
			Tier:         e.Tier,
			Rank:         e.Rank,
			LeaguePoints: e.LeaguePoints,
			Wins:         e.Wins,
			Losses:       e.Losses,
			HotStreak:    e.HotStreak,
			Veteran:      e.Veteran,
			FreshBlood:   e.FreshBlood,
			Inactive:     e.Inactive,
			SummonerID:   e.SummonerID,
			SummonerName: e.SummonerName,
			PUUID:        e.PUUID,
		}
	}

	return &LeagueListInfo{
		Tier:     l.Tier,
		LeagueID: l.LeagueID,
		Queue:    string(l.Queue),
		Name:     l.Name,
		Entries:  entries,
	}
}

// Queue type constants
const (
	QueueRankedSolo = "RANKED_SOLO_5x5"
	QueueRankedFlex = "RANKED_FLEX_SR"
)
