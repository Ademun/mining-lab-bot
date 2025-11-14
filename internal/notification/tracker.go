package notification

import (
	"context"
	"log/slog"

	"github.com/Ademun/mining-lab-bot/internal/polling"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
)

func (s *notificationService) trackSlot(ctx context.Context, slot polling.Slot) {
	exists, err := s.cache.Exists(ctx, slot.Key())
	if err != nil {
		slog.Error("Redis error", "error", err, "service", logger.ServiceNotification)
		return
	}

	if exists {
		return
	}

	key := "track:" + slot.Key()

	if err = s.cache.Set(ctx, slot, key, s.options.CacheTTL); err != nil {
		slog.Error("Redis error", "error", err, "service", logger.ServiceNotification)
	}
	s.mu.Lock()
	recordSlot(slot.Type)
	s.mu.Unlock()
	return
}

func (s *notificationService) resetUniqueSlots(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.cache.DeleteAll(ctx, "track:"); err != nil {
		slog.Error("Redis error", "error", err, "service", logger.ServiceNotification)
	}

	uniqueSlotsMetrics.Reset()
}
