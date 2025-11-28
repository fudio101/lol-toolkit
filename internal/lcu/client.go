package lcu

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"regexp"
	"sync"
	"syscall"
	"time"
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

// Client represents the League Client Update API client.
type Client struct {
	port      string
	authToken string
	http      *http.Client
}

// RerollPoints represents ARAM reroll points.
type RerollPoints struct {
	CurrentPoints    int `json:"currentPoints"`
	MaxRolls         int `json:"maxRolls"`
	NumberOfRolls    int `json:"numberOfRolls"`
	PointsCostToRoll int `json:"pointsCostToRoll"`
	PointsToReroll   int `json:"pointsToReroll"`
}

// CurrentSummoner represents the currently logged in summoner.
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

// NewClient creates a new LCU client.
func NewClient() (*Client, error) {
	port, token, err := getConnectionInfo()
	if err != nil {
		return nil, fmt.Errorf("league client not running: %w", err)
	}

	return &Client{
		port:      port,
		authToken: token,
		http:      newHTTPClient(),
	}, nil
}

// newHTTPClient creates an HTTP client for LCU (skips TLS verification).
func newHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 5 * time.Second,
	}
}

// GetCurrentSummoner returns the currently logged in summoner.
func (c *Client) GetCurrentSummoner() (*CurrentSummoner, error) {
	data, err := c.request("GET", "/lol-summoner/v1/current-summoner")
	if err != nil {
		return nil, err
	}

	var summoner CurrentSummoner
	if err := json.Unmarshal(data, &summoner); err != nil {
		return nil, fmt.Errorf("failed to parse summoner: %w", err)
	}

	return &summoner, nil
}

// IsClientRunning checks if the League client is running.
func IsClientRunning() bool {
	_, _, err := getConnectionInfo()
	return err == nil
}

// ConnectionInfo holds connection details for external use.
type ConnectionInfo struct {
	Port      string `json:"port"`
	AuthToken string `json:"authToken"`
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

// request makes a request to the LCU API.
func (c *Client) request(method, endpoint string) ([]byte, error) {
	url := fmt.Sprintf("https://127.0.0.1:%s%s", c.port, endpoint)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	// Basic auth
	auth := base64.StdEncoding.EncodeToString([]byte("riot:" + c.authToken))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		ClearCache()
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("lcu api error: %s - %s", resp.Status, string(body))
	}

	return body, nil
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
		return "", "", fmt.Errorf("league client not found")
	}

	tokenMatch := tokenRe.FindStringSubmatch(output)
	if len(tokenMatch) < 2 {
		return "", "", fmt.Errorf("auth token not found")
	}

	return portMatch[1], tokenMatch[1], nil
}
