package lcu

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ReadyCheckState represents the state of a ready check.
type ReadyCheckState string

const (
	ReadyCheckInProgress ReadyCheckState = "InProgress"
	ReadyCheckAccepted   ReadyCheckState = "Accepted"
	ReadyCheckDeclined   ReadyCheckState = "Declined"
)

// ReadyCheckResource represents a matchmaking ready check.
type ReadyCheckResource struct {
	State                          ReadyCheckState `json:"state"`
	PlayerResponse                 string          `json:"playerResponse"`
	DeclinerIds                    []int64         `json:"declinerIds"`
	DodgeWarning                   string          `json:"dodgeWarning"`
	Timer                          float64         `json:"timer"`
	SuppressUx                     bool            `json:"suppressUx"`
	ResponseRequired               bool            `json:"responseRequired"`
	EstimatedMatchmakingTimeMillis int64           `json:"estimatedMatchmakingTimeMillis"`
}

// Champion represents a champion in the collection.
type Champion struct {
	ID               int    `json:"id"`
	Owned            bool   `json:"owned"`
	Rented           bool   `json:"rented"`
	FreeToPlay       bool   `json:"freeToPlay"`
	FreeToPlayReward bool   `json:"freeToPlayReward"`
	ChampionID       int    `json:"championId"`
	Purchased        int64  `json:"purchased"`
	Alias            string `json:"alias"`
	Name             string `json:"name"`
	Active           bool   `json:"active"`
	BotEnabled       bool   `json:"botEnabled"`
	BotMmEnabled     bool   `json:"botMmEnabled"`
}

// ChampSelectAction represents an action in champion select.
type ChampSelectAction struct {
	ID                 int    `json:"id"`
	ActorCellID        int64  `json:"actorCellId"`
	ChampionID         int    `json:"championId"`
	Completed          bool   `json:"completed"`
	IsAllyAction       bool   `json:"isAllyAction"`
	IsInProgress       bool   `json:"isInProgress"`
	Type               string `json:"type"`
	ChampionPickIntent int    `json:"championPickIntent"`
}

// ChampSelectSession represents a champion select session.
type ChampSelectSession struct {
	Actions            [][]ChampSelectAction `json:"actions"`
	AllowBattleBoost   bool                  `json:"allowBattleBoost"`
	AllowReroll        bool                  `json:"allowReroll"`
	AllowSkinSelection bool                  `json:"allowSkinSelection"`
	Bans               interface{}           `json:"bans"`
	ChatDetails        interface{}           `json:"chatDetails"`
	IsSpectating       bool                  `json:"isSpectating"`
	LocalPlayerCellID  int64                 `json:"localPlayerCellId"`
	MyTeam             []interface{}         `json:"myTeam"`
	RerollsRemaining   int                   `json:"rerollsRemaining"`
	TheirTeam          []interface{}         `json:"theirTeam"`
	Timer              interface{}           `json:"timer"`
	Trades             []interface{}         `json:"trades"`
}

// AcceptMatch accepts a ready check match.
// TEMP: For testing, this will DECLINE the ready check instead of accepting it.
func (c *Client) AcceptMatch() error {
	// Use decline endpoint for safe testing
	_, err := c.Request("POST", "/lol-matchmaking/v1/ready-check/decline", nil)
	return err
}

// GetReadyCheck gets the current ready check status.
func (c *Client) GetReadyCheck() (*ReadyCheckResource, error) {
	data, err := c.Request("GET", "/lol-matchmaking/v1/ready-check", nil)
	if err != nil {
		return nil, err
	}

	var readyCheck ReadyCheckResource
	if err := json.Unmarshal(data, &readyCheck); err != nil {
		return nil, fmt.Errorf("failed to parse ready check: %w", err)
	}

	return &readyCheck, nil
}

// GetOwnedChampions gets the list of owned champions.
func (c *Client) GetOwnedChampions() ([]Champion, error) {
	data, err := c.Request("GET", "/lol-champions/v1/owned-champions-minimal", nil)
	if err != nil {
		return nil, err
	}

	var champions []Champion
	if err := json.Unmarshal(data, &champions); err != nil {
		return nil, fmt.Errorf("failed to parse champions: %w", err)
	}

	return champions, nil
}

// GetChampSelectSession gets the current champion select session.
func (c *Client) GetChampSelectSession() (*ChampSelectSession, error) {
	data, err := c.Request("GET", "/lol-champ-select/v1/session", nil)
	if err != nil {
		return nil, err
	}

	var session ChampSelectSession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to parse champ select session: %w", err)
	}

	return &session, nil
}

// SelectChampion selects a champion in champion select.
func (c *Client) SelectChampion(actionID int, championID int) error {
	url := fmt.Sprintf("/lol-champ-select/v1/session/actions/%d", actionID)

	body := map[string]int{
		"championId": championID,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	_, err = c.Request("PATCH", url, bytes.NewReader(jsonBody))
	return err
}

// LockChampion locks in the selected champion.
func (c *Client) LockChampion(actionID int) error {
	url := fmt.Sprintf("/lol-champ-select/v1/session/actions/%d/complete", actionID)
	_, err := c.Request("POST", url, nil)
	return err
}

// FindChampionID finds a champion ID by name or alias.
func (c *Client) FindChampionID(name string) (int, error) {
	champions, err := c.GetOwnedChampions()
	if err != nil {
		return -1, err
	}

	// Normalize search name
	searchName := normalizeChampionName(name)

	for _, champ := range champions {
		if normalizeChampionName(champ.Name) == searchName ||
			normalizeChampionName(champ.Alias) == searchName {
			return champ.ID, nil
		}
	}

	return -1, fmt.Errorf("champion not found: %s", name)
}

// GetMyActionID gets the action ID for the current player in champ select.
func (c *Client) GetMyActionID() (int, error) {
	session, err := c.GetChampSelectSession()
	if err != nil {
		return -1, err
	}

	if len(session.Actions) == 0 {
		return -1, fmt.Errorf("no actions in session")
	}

	// Actions is a 2D array: [phase][action]
	// First phase is usually pick phase
	actions := session.Actions[0]

	for _, action := range actions {
		if action.ActorCellID == session.LocalPlayerCellID {
			return action.ID, nil
		}
	}

	return -1, fmt.Errorf("action not found for local player")
}

// normalizeChampionName normalizes a champion name for comparison.
func normalizeChampionName(name string) string {
	// Convert to lowercase and remove spaces
	result := ""
	for _, r := range name {
		if r != ' ' && r != '\'' && r != '-' {
			result += string(r)
		}
	}
	return result
}
