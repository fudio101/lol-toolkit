package lcu

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
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

// GameflowPhase represents the current gameflow phase.
type GameflowPhase string

const (
	GameflowPhaseNone            GameflowPhase = "None"
	GameflowPhaseLobby           GameflowPhase = "Lobby"
	GameflowPhaseMatchmaking     GameflowPhase = "Matchmaking"
	GameflowPhaseReadyCheck      GameflowPhase = "ReadyCheck"
	GameflowPhaseChampSelect     GameflowPhase = "ChampSelect"
	GameflowPhaseInProgress      GameflowPhase = "InProgress"
	GameflowPhaseReconnect       GameflowPhase = "Reconnect"
	GameflowPhaseWaitingForStats GameflowPhase = "WaitingForStats"
	GameflowPhasePreEndOfGame    GameflowPhase = "PreEndOfGame"
)

// ClientState represents the current client state for auto-accept purposes.
type ClientState string

const (
	ClientStateNotInQueue  ClientState = "NotInQueue"
	ClientStateInQueue     ClientState = "InQueue"
	ClientStateMatchFound  ClientState = "MatchFound"
	ClientStateChampSelect ClientState = "ChampSelect"
	ClientStateInGame      ClientState = "InGame"
	ClientStateUnknown     ClientState = "Unknown"
)

// Polling intervals for different client states
const (
	pollIntervalVeryFast = 200 * time.Millisecond // Match found - need fast acceptance
	pollIntervalFast     = 500 * time.Millisecond // In queue - waiting for match
	pollIntervalSlow     = 3 * time.Second        // Idle/not relevant - minimal overhead
)

// buildLCUHeaders creates HTTP headers for LCU API requests.
func buildLCUHeaders() map[string]string {
	_, token, _ := getConnectionInfo()
	headers := make(map[string]string)
	if token != "" {
		headers["Authorization"] = fmt.Sprintf("Basic %s", base64Encode(fmt.Sprintf("riot:%s", token)))
	}
	headers["Accept"] = "application/json"
	return headers
}

// AcceptMatch accepts a ready check match.
func (c *Client) AcceptMatch() error {
	headers := buildLCUHeaders()

	// Wrap with LoggedCall - use interface{} as result type since POST returns empty body
	_, err := LoggedCall("POST", "/lol-matchmaking/v1/ready-check/accept", http.StatusOK, headers, func() (interface{}, error) {
		// Use lcu-gopher's client Post method directly
		resp, err := c.client.Post("/lol-matchmaking/v1/ready-check/accept", nil)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			return nil, fmt.Errorf("failed to accept match: status %d", resp.StatusCode)
		}

		return nil, nil
	})
	return err
}

// GetReadyCheck gets the current ready check status.
func (c *Client) GetReadyCheck() (*ReadyCheckResource, error) {
	headers := buildLCUHeaders()

	return LoggedCall("GET", "/lol-matchmaking/v1/ready-check", http.StatusOK, headers, func() (*ReadyCheckResource, error) {
		// Use lcu-gopher's client Get method directly
		resp, err := c.client.Get("/lol-matchmaking/v1/ready-check")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to get ready check: status %d", resp.StatusCode)
		}

		// Read body to check for "None" response before decoding
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read ready check response: %w", err)
		}

		// Check for "None" response - indicates ready check not available
		responseStr := strings.TrimSpace(string(data))
		if responseStr == "" || responseStr == `"None"` || responseStr == "None" || responseStr == "null" {
			return nil, fmt.Errorf("lcu api error: 404 Not Found - ready check not available")
		}

		// Decode JSON response
		var readyCheck ReadyCheckResource
		if err := json.Unmarshal(data, &readyCheck); err != nil {
			return nil, fmt.Errorf("failed to decode ready check: %w", err)
		}

		return &readyCheck, nil
	})
}

// GetGameflowPhase gets the current gameflow phase.
func (c *Client) GetGameflowPhase() (GameflowPhase, error) {
	headers := buildLCUHeaders()

	return LoggedCall("GET", "/lol-gameflow/v1/gameflow-phase", http.StatusOK, headers, func() (GameflowPhase, error) {
		// Use lcu-gopher's client directly - same pattern as their GetCurrentSummoner
		resp, err := c.client.Get("/lol-gameflow/v1/gameflow-phase")
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("failed to get gameflow phase: status %d", resp.StatusCode)
		}

		var phase string
		if err := json.NewDecoder(resp.Body).Decode(&phase); err != nil {
			return "", fmt.Errorf("failed to decode gameflow phase: %w", err)
		}

		// "None" is a valid state (NotInQueue/Idle)
		return GameflowPhase(phase), nil
	})
}

// GetClientState determines the current client state based on gameflow phase and other endpoints.
// This provides a high-level abstraction for auto-accept logic.
func (c *Client) GetClientState() (ClientState, error) {
	phase, err := c.GetGameflowPhase()
	if err != nil {
		return ClientStateUnknown, err
	}

	switch phase {
	case GameflowPhaseNone, GameflowPhaseLobby:
		return ClientStateNotInQueue, nil
	case GameflowPhaseMatchmaking:
		return ClientStateInQueue, nil
	case GameflowPhaseReadyCheck:
		return ClientStateMatchFound, nil
	case GameflowPhaseChampSelect:
		return ClientStateChampSelect, nil
	case GameflowPhaseInProgress, GameflowPhaseReconnect:
		return ClientStateInGame, nil
	case GameflowPhaseWaitingForStats, GameflowPhasePreEndOfGame:
		// Post-game states - consider as not in queue
		return ClientStateNotInQueue, nil
	default:
		return ClientStateUnknown, nil
	}
}

// AutoAcceptService manages auto-accept functionality.
type AutoAcceptService struct {
	client          *Client
	enabled         bool
	autoAccept      bool
	consecutive404s int // Track consecutive 404 errors to reduce polling
	mu              sync.Mutex
	stop            chan struct{}
	wg              sync.WaitGroup
	onStopped       func() // Callback when service stops due to connection error
	// Client state tracking
	lastClientState ClientState // Last detected client state
}

// NewAutoAcceptService creates a new auto-accept service.
func NewAutoAcceptService(client *Client) *AutoAcceptService {
	return &AutoAcceptService{
		client:     client,
		enabled:    false,
		autoAccept: true,
		stop:       make(chan struct{}),
	}
}

// Start starts the auto-accept service.
func (s *AutoAcceptService) Start() {
	s.mu.Lock()
	if s.enabled {
		s.mu.Unlock()
		return
	}
	s.enabled = true
	s.mu.Unlock()

	s.wg.Add(1)
	go s.run()
}

// Stop stops the auto-accept service.
func (s *AutoAcceptService) Stop() {
	s.mu.Lock()
	if !s.enabled {
		s.mu.Unlock()
		return
	}
	s.enabled = false
	stopChan := s.stop
	s.mu.Unlock()

	close(stopChan)
	s.wg.Wait()

	// Recreate stop channel for potential restart
	s.mu.Lock()
	s.stop = make(chan struct{})
	s.mu.Unlock()
}

// SetAutoAccept enables/disables auto-accept.
func (s *AutoAcceptService) SetAutoAccept(enabled bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.autoAccept = enabled
}

// SetOnStopped sets a callback to be called when the service stops due to connection error.
func (s *AutoAcceptService) SetOnStopped(callback func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onStopped = callback
}

// run is the main loop for the auto-accept service.
func (s *AutoAcceptService) run() {
	defer s.wg.Done()

	ticker := time.NewTicker(pollIntervalFast)
	defer ticker.Stop()

	lastInterval := pollIntervalFast

	for {
		select {
		case <-s.stop:
			return
		case <-ticker.C:
			s.checkAndProcess()

			// Adjust polling interval based on client state
			newInterval := s.getPollingInterval()
			if newInterval != lastInterval {
				ticker.Stop()
				ticker = time.NewTicker(newInterval)
				lastInterval = newInterval
			}
		}
	}
}

// getPollingInterval returns the appropriate polling interval based on current client state.
func (s *AutoAcceptService) getPollingInterval() time.Duration {
	s.mu.Lock()
	state := s.lastClientState
	s.mu.Unlock()

	switch state {
	case ClientStateMatchFound:
		return pollIntervalVeryFast
	case ClientStateInQueue:
		return pollIntervalFast
	default:
		return pollIntervalSlow
	}
}

// checkAndProcess checks for ready check and accepts matches.
func (s *AutoAcceptService) checkAndProcess() {
	if !s.shouldProcess() {
		return
	}

	if !IsConnected() {
		return
	}

	state, err := s.client.GetClientState()
	if err != nil {
		s.incrementConsecutive404s()
		return
	}

	s.updateClientState(state)

	// Only poll ready-check when in relevant states
	if state == ClientStateInQueue || state == ClientStateMatchFound {
		s.checkReadyCheck()
	} else {
		s.resetConsecutive404s()
	}
}

// shouldProcess checks if the service should process (enabled and auto-accept on).
func (s *AutoAcceptService) shouldProcess() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.enabled && s.autoAccept
}

// updateClientState updates the client state and handles state transitions.
func (s *AutoAcceptService) updateClientState(newState ClientState) {
	s.mu.Lock()
	lastState := s.lastClientState
	s.lastClientState = newState
	s.mu.Unlock()

	// Handle state transitions
	if lastState != newState && lastState != "" && newState == ClientStateInGame {
		s.disableAutoAccept()
	}
}

// disableAutoAccept disables auto-accept and notifies the frontend.
func (s *AutoAcceptService) disableAutoAccept() {
	s.mu.Lock()
	s.autoAccept = false
	onStopped := s.onStopped
	s.mu.Unlock()

	if onStopped != nil {
		onStopped()
	}
}

// incrementConsecutive404s increments the 404 error counter.
func (s *AutoAcceptService) incrementConsecutive404s() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.consecutive404s++
}

// resetConsecutive404s resets the 404 error counter.
func (s *AutoAcceptService) resetConsecutive404s() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.consecutive404s = 0
}

// checkReadyCheck checks for ready check and accepts matches.
func (s *AutoAcceptService) checkReadyCheck() {
	readyCheck, err := s.client.GetReadyCheck()
	if err != nil {
		s.handleReadyCheckError(err)
		return
	}

	if readyCheck == nil {
		return
	}

	s.resetConsecutive404s()

	// Accept if ready check is in progress and we haven't responded yet
	if readyCheck.State == ReadyCheckInProgress && s.shouldAcceptReadyCheck(readyCheck) {
		s.client.AcceptMatch()
	}
}

// shouldAcceptReadyCheck returns true if we should accept the ready check.
func (s *AutoAcceptService) shouldAcceptReadyCheck(readyCheck *ReadyCheckResource) bool {
	return readyCheck.PlayerResponse == "" || readyCheck.PlayerResponse == "None"
}

// handleReadyCheckError handles errors from GetReadyCheck.
func (s *AutoAcceptService) handleReadyCheckError(err error) {
	if is404Error(err) {
		s.incrementConsecutive404s()
		return
	}

	if IsConnectionRefusedError(err) {
		s.handleConnectionRefused()
	}
}

// handleConnectionRefused handles connection refused errors by stopping the service.
func (s *AutoAcceptService) handleConnectionRefused() {
	s.mu.Lock()
	wasEnabled := s.enabled
	s.enabled = false
	s.autoAccept = false
	stopChan := s.stop
	onStopped := s.onStopped
	s.mu.Unlock()

	if wasEnabled {
		select {
		case <-stopChan:
			// Channel already closed
		default:
			close(stopChan)
		}

		if onStopped != nil {
			onStopped()
		}
	}
}

// is404Error checks if an error is a 404 (expected when not in queue/champ select).
func is404Error(err error) bool {
	if err == nil {
		return false
	}
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "404") ||
		strings.Contains(errStr, "not attached") ||
		strings.Contains(errStr, "no active delegate")
}
