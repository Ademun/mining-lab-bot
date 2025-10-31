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
	Step ConversationStep       `json:"step"`
	Data map[string]interface{} `json:"data"`
}

func NewFSM(client *redis.Client) *FSM {
	slog.Info("FSM instance created", "service", logger.TelegramBot)
	return &FSM{client: client}
}

func (f *FSM) GetState(ctx context.Context, userID int64) (*State, error) {
	key := f.makeKey(userID)

	slog.Debug("Getting user state",
		"user_id", userID,
		"key", key,
		"service", logger.TelegramBot)

	data, err := f.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		slog.Debug("State not found, returning default",
			"user_id", userID,
			"service", logger.TelegramBot)
		return &State{
			Step: "",
			Data: make(map[string]interface{}),
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

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		slog.Error("Failed to unmarshal state",
			"error", err,
			"user_id", userID,
			"data", string(data),
			"service", logger.TelegramBot)
		return nil, err
	}

	slog.Debug("State retrieved successfully",
		"user_id", userID,
		"step", state.Step,
		"service", logger.TelegramBot)

	return &state, nil
}

func (f *FSM) SetStep(ctx context.Context, userID int64, step ConversationStep) error {
	key := f.makeKey(userID)

	slog.Debug("Setting conversation step",
		"user_id", userID,
		"step", step,
		"service", logger.TelegramBot)

	state, err := f.GetState(ctx, userID)
	if err != nil {
		slog.Error("Failed to get current state before setting step",
			"error", err,
			"user_id", userID,
			"service", logger.TelegramBot)
		return err
	}

	oldStep := state.Step
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

	slog.Info("Conversation step updated",
		"user_id", userID,
		"old_step", oldStep,
		"new_step", step,
		"service", logger.TelegramBot)

	return nil
}

func (f *FSM) UpdateData(ctx context.Context, userID int64, data map[string]interface{}) error {
	key := f.makeKey(userID)

	slog.Debug("Updating state data",
		"user_id", userID,
		"data_keys", getMapKeys(data),
		"service", logger.TelegramBot)

	state, err := f.GetState(ctx, userID)
	if err != nil {
		slog.Error("Failed to get current state before updating data",
			"error", err,
			"user_id", userID,
			"service", logger.TelegramBot)
		return err
	}

	for k, v := range data {
		state.Data[k] = v
	}

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

	slog.Info("State data updated successfully",
		"user_id", userID,
		"updated_keys", getMapKeys(data),
		"service", logger.TelegramBot)

	return nil
}

func (f *FSM) GetData(ctx context.Context, userID int64, key string) (interface{}, error) {
	slog.Debug("Getting specific data key",
		"user_id", userID,
		"key", key,
		"service", logger.TelegramBot)

	state, err := f.GetState(ctx, userID)
	if err != nil {
		slog.Error("Failed to get state for data retrieval",
			"error", err,
			"user_id", userID,
			"key", key,
			"service", logger.TelegramBot)
		return nil, err
	}

	value := state.Data[key]
	slog.Debug("Data key retrieved",
		"user_id", userID,
		"key", key,
		"has_value", value != nil,
		"service", logger.TelegramBot)

	return value, nil
}

func (f *FSM) ResetState(ctx context.Context, userID int64) error {
	key := f.makeKey(userID)

	slog.Info("Resetting user state",
		"user_id", userID,
		"service", logger.TelegramBot)

	state, err := f.GetState(ctx, userID)
	if err != nil {
		slog.Error("Failed to get current state before reset",
			"error", err,
			"user_id", userID,
			"service", logger.TelegramBot)
		return err
	}

	oldStep := state.Step
	state.Step = StepIdle
	state.Data = make(map[string]interface{})

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

	slog.Info("User state reset successfully",
		"user_id", userID,
		"old_step", oldStep,
		"service", logger.TelegramBot)

	return nil
}

func (f *FSM) makeKey(userID int64) string {
	return fmt.Sprintf("fsm:%d:state", userID)
}

func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
