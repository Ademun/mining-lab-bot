package fsm

import (
	"context"
	"log/slog"
	"sync"

	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type ConversationStep string

const (
	StepIdle                    ConversationStep = "idle"
	StepAwaitingLabType         ConversationStep = "awaiting_lab_type"
	StepAwaitingLabNumber       ConversationStep = "awaiting_lab_number"
	StepAwaitingLabAuditorium   ConversationStep = "awaiting_lab_auditorium"
	StepAwaitingLabDomain       ConversationStep = "awaiting_lab_domain"
	StepAwaitingLabWeekday      ConversationStep = "awaiting_lab_weekday"
	StepAwaitingLabLessons      ConversationStep = "awaiting_lab_lessons"
	StepAwaitingLabConfirmation ConversationStep = "awaiting_lab_confirmation"
)

type HandlerFunc func(ctx context.Context, api *bot.Bot, update *models.Update, state *State)
type Router struct {
	fsm      *FSM
	handlers map[ConversationStep]HandlerFunc
	mu       *sync.RWMutex
}

func NewRouter(fsm *FSM) *Router {
	return &Router{
		fsm:      fsm,
		handlers: make(map[ConversationStep]HandlerFunc),
		mu:       &sync.RWMutex{},
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

		if err := r.fsm.ResetState(ctx, userID); err != nil {
			return
		}
		next(ctx, b, update)
	}
}

func (r *Router) Transition(ctx context.Context, userID int64, nextStep ConversationStep, data map[string]interface{}) error {
	if err := r.fsm.SetStep(ctx, userID, nextStep); err != nil {
		slog.Error("Failed to update conversation step", "error", err, logger.TelegramBot)
		if err := r.fsm.ResetState(ctx, userID); err != nil {
			slog.Error("Fatal redis error when clearing conversation state", "error", err, "service", logger.TelegramBot)
		}
		return err
	}
	if data == nil {
		return nil
	}
	if err := r.fsm.UpdateData(ctx, userID, data); err != nil {
		slog.Error("Failed to update conversation data", "error", err, "service", logger.TelegramBot)
		if err := r.fsm.ResetState(ctx, userID); err != nil {
			slog.Error("Fatal redis error when clearing conversation state", "error", err, "service", logger.TelegramBot)
		}
		return err
	}
	return nil
}
