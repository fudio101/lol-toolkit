package lol

import (
	"fmt"
	"strings"
	"time"

	"github.com/KnutZuidema/golio/riot/lol"
)

// SummonerInfo represents basic summoner information for the frontend
type SummonerInfo struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
	PUUID         string `json:"puuid"`
	GameName      string `json:"gameName"`
	TagLine       string `json:"tagLine"`
	ProfileIconID int    `json:"profileIconId"`
	SummonerLevel int    `json:"summonerLevel"`
	RevisionDate  int    `json:"revisionDate"`
}

// SearchByRiotID searches for a summoner by Riot ID (gameName#tagLine)
func (c *Client) SearchByRiotID(riotID string) (*SummonerInfo, error) {
	// Parse Riot ID (format: gameName#tagLine)
	parts := strings.SplitN(riotID, "#", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid Riot ID format. Use: gameName#tagLine")
	}

	gameName := strings.TrimSpace(parts[0])
	tagLine := strings.TrimSpace(parts[1])

	if gameName == "" || tagLine == "" {
		return nil, fmt.Errorf("both game name and tag line are required")
	}

	// Get account by Riot ID
	start := time.Now()
	account, err := c.golio.Riot.Account.GetByRiotID(gameName, tagLine)
	duration := time.Since(start)
	if err != nil {
		logAPICall(APILogEntry{
			Type:       "riot",
			Method:     "GET",
			Endpoint:   "account/by-riot-id",
			Duration:   duration.Milliseconds(),
			Headers:    c.getHeaders(),
			Error:      err.Error(),
			StatusCode: 500,
		})
		return nil, fmt.Errorf("account not found: %w", err)
	}
	logAPICall(APILogEntry{
		Type:       "riot",
		Method:     "GET",
		Endpoint:   "account/by-riot-id",
		Duration:   duration.Milliseconds(),
		Headers:    c.getHeaders(),
		StatusCode: 200,
		Response:   c.marshalResponse(account),
	})

	// Get summoner by PUUID
	start = time.Now()
	summoner, err := c.golio.Riot.LoL.Summoner.GetByPUUID(account.Puuid)
	duration = time.Since(start)
	if err != nil {
		logAPICall(APILogEntry{
			Type:       "riot",
			Method:     "GET",
			Endpoint:   "summoner/by-puuid",
			Duration:   duration.Milliseconds(),
			Headers:    c.getHeaders(),
			Error:      err.Error(),
			StatusCode: 500,
		})
		return nil, fmt.Errorf("summoner not found: %w", err)
	}
	logAPICall(APILogEntry{
		Type:       "riot",
		Method:     "GET",
		Endpoint:   "summoner/by-puuid",
		Duration:   duration.Milliseconds(),
		Headers:    c.getHeaders(),
		StatusCode: 200,
		Response:   c.marshalResponse(summoner),
	})

	return &SummonerInfo{
		ID:            summoner.ID,
		AccountID:     summoner.AccountID,
		PUUID:         summoner.PUUID,
		GameName:      account.GameName,
		TagLine:       account.TagLine,
		ProfileIconID: summoner.ProfileIconID,
		SummonerLevel: summoner.SummonerLevel,
		RevisionDate:  summoner.RevisionDate,
	}, nil
}

// GetSummonerByPUUID fetches summoner info by PUUID
func (c *Client) GetSummonerByPUUID(puuid string) (*SummonerInfo, error) {
	start := time.Now()
	summoner, err := c.golio.Riot.LoL.Summoner.GetByPUUID(puuid)
	duration := time.Since(start)
	if err != nil {
		logAPICall(APILogEntry{
			Type:       "riot",
			Method:     "GET",
			Endpoint:   "summoner/by-puuid",
			Duration:   duration.Milliseconds(),
			Error:      err.Error(),
			StatusCode: 500,
		})
		return nil, err
	}
	logAPICall(APILogEntry{
		Type:       "riot",
		Method:     "GET",
		Endpoint:   "summoner/by-puuid",
		Duration:   duration.Milliseconds(),
		Headers:    c.getHeaders(),
		StatusCode: 200,
		Response:   c.marshalResponse(summoner),
	})

	// Also get account info for GameName and TagLine
	start = time.Now()
	account, err := c.golio.Riot.Account.GetByPUUID(puuid)
	duration = time.Since(start)
	if err != nil {
		logAPICall(APILogEntry{
			Type:       "riot",
			Method:     "GET",
			Endpoint:   "account/by-puuid",
			Duration:   duration.Milliseconds(),
			Headers:    c.getHeaders(),
			Error:      err.Error(),
			StatusCode: 500,
		})
		// Return summoner info without GameName/TagLine if account lookup fails
		return toSummonerInfo(summoner, "", ""), nil
	}
	logAPICall(APILogEntry{
		Type:       "riot",
		Method:     "GET",
		Endpoint:   "account/by-puuid",
		Duration:   duration.Milliseconds(),
		Headers:    c.getHeaders(),
		StatusCode: 200,
		Response:   c.marshalResponse(account),
	})

	return toSummonerInfo(summoner, account.GameName, account.TagLine), nil
}

// GetSummonerByID fetches summoner info by summoner ID
func (c *Client) GetSummonerByID(summonerID string) (*SummonerInfo, error) {
	start := time.Now()
	summoner, err := c.golio.Riot.LoL.Summoner.GetByID(summonerID)
	duration := time.Since(start)
	if err != nil {
		logAPICall(APILogEntry{
			Type:       "riot",
			Method:     "GET",
			Endpoint:   "summoner/by-id",
			Duration:   duration.Milliseconds(),
			Headers:    c.getHeaders(),
			Error:      err.Error(),
			StatusCode: 500,
		})
		return nil, err
	}
	logAPICall(APILogEntry{
		Type:       "riot",
		Method:     "GET",
		Endpoint:   "summoner/by-id",
		Duration:   duration.Milliseconds(),
		Headers:    c.getHeaders(),
		StatusCode: 200,
		Response:   c.marshalResponse(summoner),
	})

	// Also get account info for GameName and TagLine
	start = time.Now()
	account, err := c.golio.Riot.Account.GetByPUUID(summoner.PUUID)
	duration = time.Since(start)
	if err != nil {
		logAPICall(APILogEntry{
			Type:       "riot",
			Method:     "GET",
			Endpoint:   "account/by-puuid",
			Duration:   duration.Milliseconds(),
			Headers:    c.getHeaders(),
			Error:      err.Error(),
			StatusCode: 500,
		})
		return toSummonerInfo(summoner, "", ""), nil
	}
	logAPICall(APILogEntry{
		Type:       "riot",
		Method:     "GET",
		Endpoint:   "account/by-puuid",
		Duration:   duration.Milliseconds(),
		Headers:    c.getHeaders(),
		StatusCode: 200,
		Response:   c.marshalResponse(account),
	})

	return toSummonerInfo(summoner, account.GameName, account.TagLine), nil
}

// toSummonerInfo converts golio Summoner to our SummonerInfo
func toSummonerInfo(s *lol.Summoner, gameName, tagLine string) *SummonerInfo {
	return &SummonerInfo{
		ID:            s.ID,
		AccountID:     s.AccountID,
		PUUID:         s.PUUID,
		GameName:      gameName,
		TagLine:       tagLine,
		ProfileIconID: s.ProfileIconID,
		SummonerLevel: s.SummonerLevel,
		RevisionDate:  s.RevisionDate,
	}
}
