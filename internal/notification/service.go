package notification

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/Ademun/mining-lab-bot/internal/polling"
	"github.com/Ademun/mining-lab-bot/internal/subscription"
	"github.com/Ademun/mining-lab-bot/pkg/config"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"golang.org/x/time/rate"
)

type Service interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context)
	SendNotification(ctx context.Context, slot polling.Slot)
	NotifyNewSubscription(ctx context.Context, sub subscription.RequestSubscription)
}

type notificationService struct {
	subService    subscription.Service
	notifier      SlotNotifier
	options       config.NotificationConfig
	limiter       *rate.Limiter
	cache         SlotCache
	cronScheduler *cron.Cron
	mu            sync.Mutex
}

func New(subService subscription.Service, notifier SlotNotifier, client *redis.Client, opts *config.NotificationConfig) Service {
	return &notificationService{
		subService: subService,
		notifier:   notifier,
		options:    *opts,
		limiter:    rate.NewLimiter(opts.NotificationRate, 1),
		cache:      *NewSlotCache(client),
		mu:         sync.Mutex{},
	}
}

func (s *notificationService) Start(ctx context.Context) error {
	slog.Info("Starting", "service", logger.ServiceNotification)
	c := cron.New(cron.WithLocation(time.Local))
	_, err := c.AddFunc("0 0 * * *", func() {
		slog.Info("Resetting unique cache", "service", logger.ServiceNotification)
		s.resetUniqueSlots(ctx)
	})
	if err != nil {
		slog.Info("Cron error", "error", err, "service", logger.ServiceNotification)
	}
	c.Start()
	s.cronScheduler = c
	slog.Info("Started", "service", logger.ServiceNotification)
	return nil
}

func (s *notificationService) Stop(ctx context.Context) {
	<-ctx.Done()
	s.cronScheduler.Stop()
	slog.Info("Stopped", "service", logger.ServiceNotification)
}

func (s *notificationService) SendNotification(ctx context.Context, slot polling.Slot) {
	defer s.trackSlot(ctx, slot)
	exists, err := s.cache.Exists(ctx, s.options.RedisPrefix+slot.Key())
	if err != nil {
		slog.Error("Redis error", "error", err, "service", logger.ServiceNotification)
	}

	if exists {
		// Since the slot can have different available times, we should update the old one
		if err := s.cache.Set(ctx, slot, s.options.RedisPrefix+slot.Key(), s.options.CacheTTL); err != nil {
			slog.Error("Redis error", "error", err, "service", logger.ServiceNotification)
		}
		return
	}

	recordSlot(slot.Type)

	if err := s.cache.Set(ctx, slot, s.options.RedisPrefix+slot.Key(), s.options.CacheTTL); err != nil {
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
	slog.Info("sub", "sub", sub)
	items := make([]polling.Slot, 0)
	cacheSlots, errChan := s.cache.ListSlots(ctx, s.options.RedisPrefix)
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
