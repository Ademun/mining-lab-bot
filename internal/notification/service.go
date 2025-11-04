package notification

import (
	"context"
	"log/slog"
	"time"

	"github.com/Ademun/mining-lab-bot/internal/polling"
	"github.com/Ademun/mining-lab-bot/internal/subscription"
	"github.com/Ademun/mining-lab-bot/pkg/config"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

type Service interface {
	SendNotification(ctx context.Context, slot polling.Slot)
	NotifyNewSubscription(ctx context.Context, sub subscription.RequestSubscription)
}

type notificationService struct {
	subService subscription.Service
	notifier   SlotNotifier
	options    config.NotificationConfig
	limiter    *rate.Limiter
	cache      SlotCache
}

func New(subService subscription.Service, notifier SlotNotifier, client *redis.Client, opts *config.NotificationConfig) Service {
	return &notificationService{
		subService: subService,
		notifier:   notifier,
		options:    *opts,
		limiter:    rate.NewLimiter(opts.NotificationRate, 1),
		cache:      *NewSlotCache(client, opts.RedisPrefix, opts.CacheTTL),
	}
}

func (s *notificationService) SendNotification(ctx context.Context, slot polling.Slot) {
	exists, err := s.cache.Exists(ctx, slot.Key())
	if err != nil {
		slog.Error("Redis error", "error", err, "service", logger.ServiceNotification)
	}

	if exists {
		if err := s.cache.Refresh(ctx, slot.Key()); err != nil {
			slog.Error("Redis error", "error", err, "service", logger.ServiceNotification)
		}
		return
	}

	if err := s.cache.Set(ctx, slot); err != nil {
		slog.Error("Redis error", "error", err, "service", logger.ServiceNotification)
	}

	users, err := s.subService.FindUsersBySlotInfo(ctx, slot)
	if err != nil {
		slog.Error("Failed to find users", "slot", slot, "err", err, "service", logger.ServiceNotification)
	}

	for _, user := range users {
		notif := Notification{
			UserID:         user.UserID,
			PreferredTimes: user.PreferredTimes,
			Slot:           slot,
		}
		if err = s.limiter.Wait(ctx); err != nil {
			slog.Error("Limiter error", "err", err, "service", logger.ServiceNotification)
			return
		}
		recordNotification()
		s.notifier.SendNotification(ctx, notif)
	}

	if len(users) > 0 {
		slog.Info("Finished sending notifications", "total", len(users), "slot", slot, "service", logger.ServiceNotification)
	}
}

func (s *notificationService) NotifyNewSubscription(ctx context.Context, sub subscription.RequestSubscription) {
	slots, err := s.findSlotsBySubscriptionInfo(ctx, sub)
	if err != nil {
		slog.Error("Failed to find slots for subscription", "error", err, "service", logger.ServiceNotification)
	}

	prefTimes := subscription.GetSubscriptionPreferredTimes(sub)
	for _, slot := range slots {
		notif := Notification{
			UserID:         sub.UserID,
			PreferredTimes: prefTimes,
			Slot:           slot,
		}
		if err = s.limiter.Wait(ctx); err != nil {
			slog.Error("Limiter error", "err", err, "service", logger.ServiceNotification)
			return
		}
		recordNotification()
		s.notifier.SendNotification(ctx, notif)
	}

	slog.Info("Finished sending notifications", "total", len(slots), "sub", sub, "service", logger.ServiceNotification)
}

func (s *notificationService) findSlotsBySubscriptionInfo(ctx context.Context, sub subscription.RequestSubscription) ([]polling.Slot, error) {
	items := make([]polling.Slot, 0)
	cacheSlots, errChan := s.cache.ListSlots(ctx)
	prefTimes := subscription.LessonsToTimeRanges(sub.Lessons...)
	for cacheSlots != nil || errChan != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case slot, ok := <-cacheSlots:
			if !ok {
				cacheSlots = nil
				continue
			}
			if slot.Type != sub.Type {
				continue
			}
			if slot.Number != sub.LabNumber {
				continue
			}
			if sub.LabAuditorium != nil {
				if slot.Auditorium != *sub.LabAuditorium {
					continue
				}
			}
			if sub.LabDomain != nil {
				if slot.Domain != *sub.LabDomain {
					continue
				}
			}
			if matchesPreferredTimes(slot.TimesTeachers, sub.Weekday, prefTimes) {
				items = append(items, slot)
			}
		case err, ok := <-errChan:
			if !ok {
				errChan = nil
				continue
			}
			return nil, err
		}
	}
	return items, nil
}

func matchesPreferredTimes(slotTimes map[time.Time][]string, subWeekday *int, prefTimes []subscription.TimeRange) bool {
	if subWeekday == nil {
		return true
	}
	for slotTime := range slotTimes {
		slotWeekday := int(slotTime.Weekday())
		if slotWeekday != *subWeekday {
			continue
		}
		slotTimeStr := slotTime.Format("15:04")
		for _, prefTime := range prefTimes {
			if slotTimeStr >= prefTime.TimeStart && slotTimeStr < prefTime.TimeEnd {
				return true
			}
		}
	}
	return false
}
