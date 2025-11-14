package notification

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Ademun/mining-lab-bot/internal/polling"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrMarshal   = errors.New("failed to marshal data")
	ErrUnmarshal = errors.New("failed to unmarshal data")
)

type SlotCache struct {
	client *redis.Client
}

func NewSlotCache(client *redis.Client) *SlotCache {
	return &SlotCache{
		client: client,
	}
}

func (c *SlotCache) Set(ctx context.Context, slot polling.Slot, key string, ttl time.Duration) error {
	data, err := json.Marshal(slot)
	if err != nil {
		return ErrMarshal
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *SlotCache) Get(ctx context.Context, key string) (*polling.Slot, error) {
	data, err := c.client.Get(ctx, key).Result()
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
	_, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (c *SlotCache) DeleteAll(ctx context.Context, prefix string) error {
	var cursor uint64
	for {
		keys, nextCursor, err := c.client.Scan(ctx, cursor, prefix+"*", 20).Result()
		if err != nil {
			return err
		}
		if len(keys) == 0 {
			return nil
		}
		if err := c.client.Del(ctx, keys...).Err(); err != nil {
			return err
		}
		if nextCursor == 0 {
			return nil
		}
		cursor = nextCursor
	}
}

func (c *SlotCache) ListSlots(ctx context.Context, key string) (chan polling.Slot, chan error) {
	slots := make(chan polling.Slot)
	errChan := make(chan error)
	eg, errCtx := errgroup.WithContext(ctx)
	go func() {
		eg.Go(func() error {
			var cursor uint64
			for {
				select {
				case <-errCtx.Done():
					return errCtx.Err()
				default:
					keys, nextCursor, err := c.client.Scan(ctx, cursor, key+"*", 20).Result()
					if err != nil {
						return err
					}
					if len(keys) == 0 {
						return nil
					}
					values, err := c.client.MGet(ctx, keys...).Result()
					if err != nil {
						return err
					}
					for _, value := range values {
						value, ok := value.(string)
						if !ok {
							return ErrUnmarshal
						}
						slot := &polling.Slot{}
						if err := json.Unmarshal([]byte(value), slot); err != nil {
							return ErrUnmarshal
						}
						slots <- *slot
					}
					cursor = nextCursor
					if cursor == 0 {
						return nil
					}
				}
			}
		})
		if err := eg.Wait(); err != nil {
			errChan <- err
		}
		close(slots)
		close(errChan)
	}()

	return slots, errChan
}

func (c *SlotCache) makeKey(prefix string, slot polling.Slot) string {
	return prefix + slot.Key()
}
