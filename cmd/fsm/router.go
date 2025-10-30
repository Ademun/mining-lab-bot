package fsm

import (
	"context"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type ConversationStep string

const (
	StepAwaitingLabType      ConversationStep = "awaiting_lab_type"
	StepWaitingLabNumber     ConversationStep = "awaiting_lab_number"
	StepWaitingLabAuditorium ConversationStep = "awaiting_lab_auditorium"
	StepAwaitingLabDomain    ConversationStep = "awaiting_lab_domain"
	StepAwaitingLabWeekday   ConversationStep = "awaiting_lab_weekday"
	StepAwaitingLabLessons   ConversationStep = "awaiting_lab_lessons"
)

type HandlerFunc func(ctx context.Context, api *bot.Bot, update *models.Update, state *State)
type Router struct {
	fsm      *FSM
	handlers map[ConversationStep]HandlerFunc
	mu       sync.RWMutex
}

func NewRouter(fsm *FSM) *Router {
	return &Router{
		fsm:      fsm,
		handlers: make(map[ConversationStep]HandlerFunc),
		mu:       sync.RWMutex{},
	}
}

func (r *Router) RegisterHandler(step ConversationStep, handler HandlerFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[step] = handler
}

func (r *Router) Middleware(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		var userID int64
		if update.Message != nil {
			userID = update.Message.From.ID
		} else if update.CallbackQuery != nil {
			userID = update.CallbackQuery.From.ID
		} else {
			return
		}

		state, err := r.fsm.GetState(ctx, userID)
		if err != nil {
			return
		}

		r.mu.RLock()
		handler, exists := r.handlers[state.Step]
		r.mu.RUnlock()

		if exists {
			handler(ctx, b, update, state)
			return
		}

		next(ctx, b, update)
	}
}
