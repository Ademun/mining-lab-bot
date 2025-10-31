package cmd

import (
	"sync"
	"time"
)

type subscriptionData struct {
	LabNumber     int
	LabAuditorium int
	Weekday       *time.Weekday
	Daytime       *string
}

type userState struct {
	Step conversationStep
	Data subscriptionData
}

type stateManager struct {
	states map[int64]*userState
	mu     *sync.RWMutex
}

func newStateManager() *stateManager {
	return &stateManager{
		states: make(map[int64]*userState),
		mu:     &sync.RWMutex{},
	}
}

func (sm *stateManager) get(userID int64) (*userState, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	state, exists := sm.states[userID]
	return state, exists
}

func (sm *stateManager) set(userID int64, state *userState) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.states[userID] = state
}

func (sm *stateManager) clear(userID int64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.states, userID)
}
