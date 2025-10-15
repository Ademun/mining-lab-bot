package notification

import (
	"context"
	"log/slog"
	"time"

	"github.com/Ademun/mining-lab-bot/internal/subscription"
	"github.com/Ademun/mining-lab-bot/pkg/cache"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/Ademun/mining-lab-bot/pkg/model"
	"github.com/Ademun/mining-lab-bot/pkg/notifier"
)

type Service interface {
	SendNotification(ctx context.Context, slot model.Slot)
	NotifyNewSubscription(ctx context.Context, sub model.Subscription)
}

type notificationService struct {
	subService subscription.Service
	notifier   notifier.SlotNotifier
	cache      *cache.TTLCache[model.Slot]
}

func New(subService subscription.Service, notifier notifier.SlotNotifier) Service {
	return &notificationService{
		subService: subService,
		notifier:   notifier,
		cache:      cache.NewTTLCache[model.Slot](time.Minute*5, time.Minute*10),
	}
}

func (s *notificationService) SendNotification(ctx context.Context, slot model.Slot) {
	_, exists := s.cache.Get(slot.Key())
	s.cache.Set(slot.Key(), slot)

	if exists {
		return
	}

	subs, err := s.subService.FindSubscriptionsBySlotInfo(ctx, slot)
	if err != nil {
		slog.Error("Failed to find subscriptions", "slot", slot, "err", err, "service", logger.ServiceNotification)
	}

	for _, sub := range subs {
		notif := model.Notification{
			UserID: sub.UserID,
			ChatID: sub.ChatID,
			Slot:   slot,
		}
		slog.Info("Sending notification", "notification", notif, "service", logger.ServiceNotification)
		s.notifier.SendNotification(ctx, notif)
	}
}

func (s *notificationService) NotifyNewSubscription(ctx context.Context, sub model.Subscription) {
	slots := s.findSlotsBySubscriptionInfo(sub)

	for _, slot := range slots {
		notif := model.Notification{
			UserID: sub.UserID,
			ChatID: sub.ChatID,
			Slot:   slot,
		}
		slog.Info("Sending notification", "notification", notif, "service", logger.ServiceNotification)
		s.notifier.SendNotification(ctx, notif)
	}
}

func (s *notificationService) findSlotsBySubscriptionInfo(sub model.Subscription) []model.Slot {
	slots := s.cache.List()
	items := make([]model.Slot, 0)
	for _, slot := range slots {
		if slot.LabNumber == sub.LabNumber && slot.LabAuditorium == sub.LabAuditorium {
			items = append(items, slot)
		}
	}
	return items
}
