package notification

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/Ademun/mining-lab-bot/internal/subscription"
	"github.com/Ademun/mining-lab-bot/pkg/cache"
	"github.com/Ademun/mining-lab-bot/pkg/event"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/Ademun/mining-lab-bot/pkg/metrics"
	"github.com/Ademun/mining-lab-bot/pkg/model"
)

type NotificationService interface {
	Start() error
	CheckCurrentSlots(ctx context.Context, sub model.Subscription)
}

type notificationService struct {
	eventBus   *event.Bus
	subService subscription.SubscriptionService
	cache      *cache.TTLCache[model.Slot]
}

func New(eb *event.Bus, subService subscription.SubscriptionService) NotificationService {
	return &notificationService{
		eventBus:   eb,
		subService: subService,
		cache:      cache.NewTTLCache[model.Slot](time.Minute*5, time.Minute*10),
	}
}

func (s *notificationService) Start() error {
	slog.Info("Starting", "service", logger.ServiceNotification)
	event.Subscribe(s.eventBus, s.handleNewSlot)
	slog.Info("Started", "service", logger.ServiceNotification)
	return nil
}

var logCounter atomic.Int64

func (s *notificationService) handleNewSlot(ctx context.Context, slotEvent event.NewSlotEvent) {
	_, exists := s.cache.Get(slotEvent.Slot.Key())
	notifCounter := 0
	if !exists {
		logCounter.Add(1)
		slog.Info("New slot", "seq", logCounter.Load(), "data", slotEvent.Slot, "service", logger.ServiceNotification)
		subs, err := s.subService.FindSubscriptionsBySlotInfo(ctx, slotEvent.Slot)
		if err != nil {
			slog.Error("Failed to find subscriptions for slot", "seq", logCounter.Load(), "data", slotEvent.Slot, "error", err, "service", logger.ServiceNotification)
		}

		for _, sub := range subs {
			notifCounter++
			notif := model.Notification{UserID: sub.UserID, ChatID: sub.ChatID, Slot: slotEvent.Slot}
			slog.Info("Sending notification", "seq", logCounter.Load(), "data", notif, "service", logger.ServiceNotification)
			event.Publish(ctx, s.eventBus, event.NewNotificationEvent{Notification: notif})
		}
	}

	s.cache.Set(slotEvent.Slot.Key(), slotEvent.Slot)

	metrics.Global().RecordNotificationResults(notifCounter, len(s.cache.List()))
}

func (s *notificationService) CheckCurrentSlots(ctx context.Context, sub model.Subscription) {
	notifCounter := 0
	slots := s.findSlotsBySubscriptionInfo(sub)
	for _, slot := range slots {
		notifCounter++
		notif := model.Notification{UserID: sub.UserID, ChatID: sub.ChatID, Slot: slot}
		slog.Info("Sending notification", "data", notif, "service", logger.ServiceNotification)
		event.Publish(ctx, s.eventBus, event.NewNotificationEvent{Notification: notif})
	}

	metrics.Global().RecordNotificationResults(notifCounter, len(s.cache.List()))
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
