package lcu

import (
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"sync"
	"syscall"
	"time"

	lcuclient "github.com/its-haze/lcu-gopher"
)

// Connection cache settings.
const cacheTTL = 30 * time.Second

var (
	cache      connectionCache
	cacheMutex sync.RWMutex
)

type connectionCache struct {
	port    string
	token   string
	updated time.Time
}

// Client is a thin wrapper around the lcu-gopher client.
type Client struct {
	client *lcuclient.Client
}

// CurrentSummoner represents the currently logged in summoner (DTO exposed to the frontend).
type CurrentSummoner struct {
	AccountID                   int64        `json:"accountId"`
	DisplayName                 string       `json:"displayName"`
	GameName                    string       `json:"gameName"`
	TagLine                     string       `json:"tagLine"`
	InternalName                string       `json:"internalName"`
	NameChangeFlag              bool         `json:"nameChangeFlag"`
	PercentCompleteForNextLevel int          `json:"percentCompleteForNextLevel"`
	ProfileIconID               int          `json:"profileIconId"`
	PUUID                       string       `json:"puuid"`
	RerollPoints                RerollPoints `json:"rerollPoints"`
	SummonerID                  int64        `json:"summonerId"`
	SummonerLevel               int          `json:"summonerLevel"`
	XpSinceLastLevel            int          `json:"xpSinceLastLevel"`
	XpUntilNextLevel            int          `json:"xpUntilNextLevel"`
}

// RerollPoints represents ARAM reroll points.
type RerollPoints struct {
	CurrentPoints    int `json:"currentPoints"`
	MaxRolls         int `json:"maxRolls"`
	NumberOfRolls    int `json:"numberOfRolls"`
	PointsCostToRoll int `json:"pointsCostToRoll"`
	PointsToReroll   int `json:"pointsToReroll"`
}

// ConnectionInfo holds connection details for external use.
type ConnectionInfo struct {
	Port      string `json:"port"`
	AuthToken string `json:"authToken"`
}

// NewClient creates a new LCU client using lcu-gopher.
func NewClient() (*Client, error) {
	config := lcuclient.DefaultConfig()
	config.AwaitConnection = false
	config.Timeout = 5 * time.Second

	client, err := lcuclient.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("league client not running: %w", err)
	}

	// Connect to the client
	if err := client.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &Client{client: client}, nil
}

// GetCurrentSummoner returns the currently logged in summoner.
func (c *Client) GetCurrentSummoner() (*CurrentSummoner, error) {
	// Use lcu-gopher's built-in method
	lcuSummoner, err := c.client.GetCurrentSummoner()
	if err != nil {
		return nil, err
	}

	// Convert to our exported type
	summoner := &CurrentSummoner{
		AccountID:                   lcuSummoner.AccountID,
		DisplayName:                 lcuSummoner.DisplayName,
		GameName:                    lcuSummoner.GameName,
		TagLine:                     lcuSummoner.TagLine,
		InternalName:                lcuSummoner.InternalName,
		NameChangeFlag:              lcuSummoner.NameChangeFlag,
		PercentCompleteForNextLevel: lcuSummoner.PercentCompleteForNextLevel,
		ProfileIconID:               lcuSummoner.ProfileIconID,
		PUUID:                       lcuSummoner.Puuid,
		RerollPoints: RerollPoints{
			CurrentPoints:    lcuSummoner.RerollPoints.CurrentPoints,
			MaxRolls:         lcuSummoner.RerollPoints.MaxRolls,
			NumberOfRolls:    lcuSummoner.RerollPoints.NumberOfRolls,
			PointsCostToRoll: lcuSummoner.RerollPoints.PointsCostToRoll,
			PointsToReroll:   lcuSummoner.RerollPoints.PointsToReroll,
		},
		SummonerID:       lcuSummoner.SummonerID,
		SummonerLevel:    lcuSummoner.SummonerLevel,
		XpSinceLastLevel: lcuSummoner.XpSinceLastLevel,
		XpUntilNextLevel: lcuSummoner.XpUntilNextLevel,
	}

	return summoner, nil
}

// GetConnectionInfo returns the current LCU connection info, or nil if not connected.
func GetConnectionInfo() *ConnectionInfo {
	port, token, err := getConnectionInfo()
	if err != nil {
		return nil
	}

	return &ConnectionInfo{
		Port:      port,
		AuthToken: token,
	}
}

// IsClientRunning checks if the League client is running.
func IsClientRunning() bool {
	_, _, err := getConnectionInfo()
	return err == nil
}

// ClearCache clears the cached connection info.
func ClearCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	cache = connectionCache{}
}

// getConnectionInfo returns cached or fresh connection info.
func getConnectionInfo() (string, string, error) {
	// Check cache
	cacheMutex.RLock()
	if cache.port != "" && time.Since(cache.updated) < cacheTTL {
		port, token := cache.port, cache.token
		cacheMutex.RUnlock()
		return port, token, nil
	}
	cacheMutex.RUnlock()

	// Fetch new connection info
	port, token, err := findFromProcess()
	if err != nil {
		ClearCache()
		return "", "", err
	}

	// Update cache
	cacheMutex.Lock()
	cache = connectionCache{port: port, token: token, updated: time.Now()}
	cacheMutex.Unlock()

	return port, token, nil
}

// findFromProcess extracts LCU connection info from running process.
func findFromProcess() (string, string, error) {
	cmd := exec.Command("wmic", "process", "where", "name='LeagueClientUx.exe'", "get", "commandline")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
	}

	output, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("leagueClientUx.exe not running")
	}

	return parseProcessArgs(string(output))
}

// parseProcessArgs extracts port and token from process arguments.
func parseProcessArgs(output string) (string, string, error) {
	portRe := regexp.MustCompile(`--app-port=(\d+)`)
	tokenRe := regexp.MustCompile(`--remoting-auth-token=([\w-]+)`)

	portMatch := portRe.FindStringSubmatch(output)
	if len(portMatch) < 2 {
		return "", "", fmt.Errorf("league client port not found")
	}

	tokenMatch := tokenRe.FindStringSubmatch(output)
	if len(tokenMatch) < 2 {
		return "", "", fmt.Errorf("auth token not found")
	}

	return portMatch[1], tokenMatch[1], nil
}

// Request makes a raw HTTP request to the LCU API.
func (c *Client) Request(method, endpoint string, body io.Reader) ([]byte, error) {
	resp, err := c.client.Request(method, endpoint, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("lcu api error: %s - %s", resp.Status, string(data))
	}

	return data, nil
}
