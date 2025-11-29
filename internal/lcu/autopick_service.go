package lcu

import (
	"fmt"
	"sync"
	"time"
)

// AutoPickService manages auto-accept and auto-pick functionality.
type AutoPickService struct {
	client          *Client
	enabled         bool
	autoAccept      bool
	autoPick        bool
	autoLock        bool
	selectedChampID int
	lastAcceptTime  time.Time
	accepted        bool
	picked          bool
	locked          bool
	mu              sync.Mutex
	stop            chan struct{}
	wg              sync.WaitGroup
}

// NewAutoPickService creates a new auto-pick service.
func NewAutoPickService(client *Client) *AutoPickService {
	return &AutoPickService{
		client:     client,
		enabled:    false,
		autoAccept: true,
		autoPick:   true,
		autoLock:   true,
		stop:       make(chan struct{}),
	}
}

// Start starts the auto-pick service.
func (s *AutoPickService) Start() {
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

// Stop stops the auto-pick service.
func (s *AutoPickService) Stop() {
	s.mu.Lock()
	if !s.enabled {
		s.mu.Unlock()
		return
	}
	s.enabled = false
	s.mu.Unlock()

	close(s.stop)
	s.wg.Wait()
}

// SetAutoAccept enables/disables auto-accept.
func (s *AutoPickService) SetAutoAccept(enabled bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.autoAccept = enabled
}

// SetAutoPick enables/disables auto-pick.
func (s *AutoPickService) SetAutoPick(enabled bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.autoPick = enabled
}

// SetAutoLock enables/disables auto-lock.
func (s *AutoPickService) SetAutoLock(enabled bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.autoLock = enabled
}

// SetChampion sets the champion to auto-pick by name.
func (s *AutoPickService) SetChampion(championName string) error {
	champID, err := s.client.FindChampionID(championName)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.selectedChampID = champID
	s.picked = false
	s.locked = false
	return nil
}

// SetChampionID sets the champion to auto-pick by ID.
func (s *AutoPickService) SetChampionID(championID int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.selectedChampID = championID
	s.picked = false
	s.locked = false
}

// run is the main loop for the auto-pick service.
func (s *AutoPickService) run() {
	defer s.wg.Done()

	ticker := time.NewTicker(500 * time.Millisecond) // Check every 500ms
	defer ticker.Stop()

	for {
		select {
		case <-s.stop:
			return
		case <-ticker.C:
			s.checkAndProcess()
		}
	}
}

// checkAndProcess checks for ready check and champ select, then processes them.
func (s *AutoPickService) checkAndProcess() {
	s.mu.Lock()
	enabled := s.enabled
	autoAccept := s.autoAccept
	autoPick := s.autoPick
	autoLock := s.autoLock
	selectedChampID := s.selectedChampID
	lastAcceptTime := s.lastAcceptTime
	accepted := s.accepted
	picked := s.picked
	locked := s.locked
	s.mu.Unlock()

	if !enabled {
		return
	}

	// Check for ready check
	if autoAccept {
		readyCheck, err := s.client.GetReadyCheck()
		if err == nil && readyCheck != nil {
			if readyCheck.State == ReadyCheckInProgress && !accepted {
				// 10 second cooldown between accepts
				if time.Since(lastAcceptTime) > 10*time.Second {
					if err := s.client.AcceptMatch(); err == nil {
						s.mu.Lock()
						s.accepted = true
						s.lastAcceptTime = time.Now()
						s.picked = false
						s.locked = false
						s.mu.Unlock()
					}
				}
			} else if readyCheck.State != ReadyCheckInProgress {
				s.mu.Lock()
				s.accepted = false
				s.mu.Unlock()
			}
		}
	}

	// Check for champ select
	if autoPick || autoLock {
		session, err := s.client.GetChampSelectSession()
		if err == nil && session != nil {
			// Reset accepted when champ select starts
			s.mu.Lock()
			s.accepted = false
			s.mu.Unlock()

			// Get action ID
			actionID, err := s.getActionID(session)
			if err == nil && actionID > 0 {
				// Auto pick
				if autoPick && selectedChampID > 0 && !picked {
					if err := s.client.SelectChampion(actionID, selectedChampID); err == nil {
						s.mu.Lock()
						s.picked = true
						s.mu.Unlock()
					}
				}

				// Auto lock
				if autoLock && !locked {
					if err := s.client.LockChampion(actionID); err == nil {
						s.mu.Lock()
						s.locked = true
						s.picked = false
						s.mu.Unlock()
					}
				}
			}
		} else {
			// Reset states when not in champ select
			s.mu.Lock()
			s.picked = false
			s.locked = false
			s.mu.Unlock()
		}
	}
}

// getActionID gets the action ID from a session.
func (s *AutoPickService) getActionID(session *ChampSelectSession) (int, error) {
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
