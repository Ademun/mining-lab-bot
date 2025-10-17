package cmd

import "sync"

type conversationStep string

const (
	stepAwaitingLabNumber  conversationStep = "awaiting_lab_number"
	stepAwaitingAuditorium conversationStep = "awaiting_auditorium"
	stepAwaitingWeekday    conversationStep = "awaiting_weekday"
	stepAwaitingTime       conversationStep = "awaiting_time"
	stepAwaitingTeacher    conversationStep = "awaiting_teacher"
	stepConfirming         conversationStep = "confirming"
)

type subscriptionData struct {
	LabNumber  int
	Auditorium int
	Weekday    string
	TimeInput  string
	Teacher    string
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

func (sm *stateManager) exists(userID int64) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	_, exists := sm.states[userID]
	return exists
}
