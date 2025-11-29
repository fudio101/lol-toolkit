package lcu

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
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
	headers := buildClientHeaders()

	return LoggedCall("GET", "/lol-summoner/v1/current-summoner", http.StatusOK, headers, func() (*CurrentSummoner, error) {
		lcuSummoner, err := c.client.GetCurrentSummoner()
		if err != nil {
			return nil, err
		}

		return &CurrentSummoner{
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
		}, nil
	})
}

// buildClientHeaders creates HTTP headers for LCU API requests.
func buildClientHeaders() map[string]string {
	_, token, _ := getConnectionInfo()
	headers := make(map[string]string)
	if token != "" {
		headers["Authorization"] = fmt.Sprintf("Basic %s", base64Encode(fmt.Sprintf("riot:%s", token)))
	}
	headers["Accept"] = "application/json"
	return headers
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

// base64Encode encodes a string to base64.
func base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// Request makes a raw HTTP request to the LCU API and logs it.
func (c *Client) Request(method, endpoint string, body io.Reader) ([]byte, error) {
	if endpoint != "GetLCUStatus" && !IsConnected() {
		return c.handleDisconnected(method, endpoint)
	}

	start := time.Now()
	headers := c.buildRequestHeaders(body)

	resp, err := c.client.Request(method, endpoint, body)
	duration := time.Since(start)

	if err != nil {
		return c.handleRequestError(method, endpoint, duration, headers, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.handleReadError(method, endpoint, duration, headers, err)
	}

	return c.handleResponse(method, endpoint, resp, data, duration, headers)
}

// handleDisconnected handles requests when client is disconnected.
func (c *Client) handleDisconnected(method, endpoint string) ([]byte, error) {
	err := fmt.Errorf("league client not connected")
	headers := buildClientHeaders()
	LogError(method, endpoint, 0, headers, err)
	return nil, err
}

// buildRequestHeaders builds headers for an HTTP request.
func (c *Client) buildRequestHeaders(body io.Reader) map[string]string {
	headers := buildClientHeaders()
	if body != nil {
		headers["Content-Type"] = "application/json"
	}
	return headers
}

// handleRequestError handles errors from the HTTP request.
func (c *Client) handleRequestError(method, endpoint string, duration time.Duration, headers map[string]string, err error) ([]byte, error) {
	LogError(method, endpoint, duration, headers, err)
	HandleConnectionError(err, endpoint)
	return nil, err
}

// handleReadError handles errors from reading the response body.
func (c *Client) handleReadError(method, endpoint string, duration time.Duration, headers map[string]string, err error) ([]byte, error) {
	LogError(method, endpoint, duration, headers, err)
	HandleConnectionError(err, endpoint)
	return nil, err
}

// handleResponse handles the HTTP response.
func (c *Client) handleResponse(method, endpoint string, resp *http.Response, data []byte, duration time.Duration, headers map[string]string) ([]byte, error) {
	responseBody := string(data)

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		err := fmt.Errorf("lcu api error: %s - %s", resp.Status, responseBody)
		LogRequest(method, endpoint, resp.StatusCode, duration, headers, responseBody, err)
		return nil, err
	}

	LogSuccess(method, endpoint, resp.StatusCode, duration, headers, responseBody)
	return data, nil
}
