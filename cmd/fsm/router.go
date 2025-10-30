package fsm

import (
	"context"
	"sync"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type HandlerFunc func(ctx context.Context, api *bot.Bot, update *models.Update, state *State)
type Router struct {
	fsm      *FSM
	handlers map[string]HandlerFunc
	mu       sync.RWMutex
}

func NewRouter(fsm *FSM) *Router {
	return &Router{
		fsm:      fsm,
		handlers: make(map[string]HandlerFunc),
		mu:       sync.RWMutex{},
	}
}

func (r *Router) RegisterHandler(name string, handler HandlerFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[name] = handler
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
