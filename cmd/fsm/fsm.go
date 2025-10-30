package fsm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

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
	return &FSM{client: client}
}

func (f *FSM) GetState(ctx context.Context, userID int64) (*State, error) {
	key := f.makeKey(userID)
	data, err := f.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return &State{
			Step: "",
			Data: make(map[string]interface{}),
		}, nil
	}
	if err != nil {
		return nil, err
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}

func (f *FSM) SetStep(ctx context.Context, userID int64, step ConversationStep) error {
	key := f.makeKey(userID)
	state, err := f.GetState(ctx, userID)
	if err != nil {
		return err
	}
	state.Step = step
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return f.client.Set(ctx, key, data, 24*time.Hour).Err()
}

func (f *FSM) UpdateData(ctx context.Context, userID int64, field string, value interface{}) error {
	key := f.makeKey(userID)
	state, err := f.GetState(ctx, userID)
	if err != nil {
		return err
	}
	state.Data[field] = value
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return f.client.Set(ctx, key, data, 24*time.Hour).Err()
}

func (f *FSM) GetData(ctx context.Context, userID int64, key string) (interface{}, error) {
	state, err := f.GetState(ctx, userID)
	if err != nil {
		return nil, err
	}
	return state.Data[key], nil
}

func (f *FSM) ResetState(ctx context.Context, userID int64) error {
	key := f.makeKey(userID)
	state, err := f.GetState(ctx, userID)
	if err != nil {
		return err
	}
	state.Step = ""
	state.Data = make(map[string]interface{})
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return f.client.Set(ctx, key, data, 24*time.Hour).Err()
}

func (f *FSM) makeKey(userID int64) string {
	return fmt.Sprintf("fsm:%d:state", userID)
}
