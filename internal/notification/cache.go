package notification

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Ademun/mining-lab-bot/internal/polling"
	"github.com/redis/go-redis/v9"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrMarshal   = errors.New("failed to marshal data")
	ErrUnmarshal = errors.New("failed to unmarshal data")
)

type SlotCache struct {
	client    *redis.Client
	keyPrefix string
	TTL       time.Duration
}

func (c *SlotCache) Set(ctx context.Context, slot polling.Slot) error {
	data, err := json.Marshal(slot)
	if err != nil {
		return ErrMarshal
	}
	return c.client.Set(ctx, c.makeKey(slot.Key()), data, c.TTL).Err()
}

func (c *SlotCache) Get(ctx context.Context, key string) (*polling.Slot, error) {
	data, err := c.client.Get(ctx, c.makeKey(key)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	result := &polling.Slot{}
	if err := json.Unmarshal([]byte(data), result); err != nil {
		return nil, ErrUnmarshal
	}
	return result, nil
}

func (c *SlotCache) Exists(ctx context.Context, key string) (bool, error) {
	_, err := c.client.Get(ctx, c.makeKey(key)).Result()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (c *SlotCache) Refresh(ctx context.Context, key string) error {
	_, err := c.client.Expire(ctx, c.makeKey(key), c.TTL).Result()
	if err != nil {
		return err
	}
	return err
}

func (c *SlotCache) makeKey(key string) string {
	return c.keyPrefix + key
}
