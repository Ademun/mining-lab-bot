package fsm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type FSM struct {
	client *redis.Client
}

type State struct {
	Step ConversationStep `json:"step"`
	Data StateData        `json:"data"`
}

func NewFSM(client *redis.Client) *FSM {
	return &FSM{client: client}
}

func (f *FSM) GetState(ctx context.Context, userID int64) (*State, error) {
	key := f.makeKey(userID)

	data, err := f.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return &State{
			Step: StepIdle,
			Data: &IdleData{},
		}, nil
	}
	if err != nil {
		slog.Error("Failed to get state from Redis",
			"error", err,
			"user_id", userID,
			"key", key,
			"service", logger.TelegramBot)
		return nil, err
	}

	var wrapper struct {
		Step ConversationStep `json:"step"`
		Data json.RawMessage  `json:"data"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		slog.Error("Failed to unmarshal state",
			"error", err,
			"user_id", userID,
			"data", string(data),
			"service", logger.TelegramBot)
		return nil, err
	}

	stateData := dataTypeForStep(wrapper.Step)
	if err := json.Unmarshal(wrapper.Data, stateData); err != nil {
		slog.Error("Failed to unmarshal state",
			"error", err,
			"user_id", userID,
			"data", string(data),
			"service", logger.TelegramBot)
		return nil, err
	}

	return &State{
		Step: wrapper.Step,
		Data: stateData,
	}, nil
}

func (f *FSM) SetStep(ctx context.Context, userID int64, step ConversationStep) error {
	key := f.makeKey(userID)

	state, err := f.GetState(ctx, userID)
	if err != nil {
		slog.Error("Failed to get current state before setting step",
			"error", err,
			"user_id", userID,
			"service", logger.TelegramBot)
		return err
	}
	state.Step = step

	data, err := json.Marshal(state)
	if err != nil {
		slog.Error("Failed to marshal state",
			"error", err,
			"user_id", userID,
			"step", step,
			"service", logger.TelegramBot)
		return err
	}

	if err := f.client.Set(ctx, key, data, 24*time.Hour).Err(); err != nil {
		slog.Error("Failed to save step to Redis",
			"error", err,
			"user_id", userID,
			"step", step,
			"service", logger.TelegramBot)
		return err
	}

	return nil
}

func (f *FSM) UpdateData(ctx context.Context, userID int64, data StateData) error {
	key := f.makeKey(userID)

	state, err := f.GetState(ctx, userID)
	if err != nil {
		slog.Error("Failed to get current state before updating data",
			"error", err,
			"user_id", userID,
			"service", logger.TelegramBot)
		return err
	}

	state.Data = data

	newData, err := json.Marshal(state)
	if err != nil {
		slog.Error("Failed to marshal updated state",
			"error", err,
			"user_id", userID,
			"service", logger.TelegramBot)
		return err
	}

	if err := f.client.Set(ctx, key, newData, 24*time.Hour).Err(); err != nil {
		slog.Error("Failed to save updated data to Redis",
			"error", err,
			"user_id", userID,
			"service", logger.TelegramBot)
		return err
	}

	return nil
}

func (f *FSM) ResetState(ctx context.Context, userID int64) error {
	key := f.makeKey(userID)

	state, err := f.GetState(ctx, userID)
	if err != nil {
		slog.Error("Failed to get current state before reset",
			"error", err,
			"user_id", userID,
			"service", logger.TelegramBot)
		return err
	}

	state.Step = StepIdle
	state.Data = &IdleData{}

	data, err := json.Marshal(state)
	if err != nil {
		slog.Error("Failed to marshal reset state",
			"error", err,
			"user_id", userID,
			"service", logger.TelegramBot)
		return err
	}

	if err := f.client.Set(ctx, key, data, 24*time.Hour).Err(); err != nil {
		slog.Error("Failed to save reset state to Redis",
			"error", err,
			"user_id", userID,
			"service", logger.TelegramBot)
		return err
	}

	return nil
}

func (f *FSM) makeKey(userID int64) string {
	return fmt.Sprintf("fsm:%d:state", userID)
}
