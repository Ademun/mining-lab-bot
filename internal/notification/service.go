package notification

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/Ademun/mining-lab-bot/internal/subscription"
	"github.com/Ademun/mining-lab-bot/pkg/cache"
	"github.com/Ademun/mining-lab-bot/pkg/event"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/Ademun/mining-lab-bot/pkg/model"
)

type NotificationService interface {
	Start() error
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

func (s *notificationService) handleNewSlot(ctx context.Context, slot model.Slot) {
	slog.Info("New slot", "data", slot, "service", logger.ServiceNotification)
	_, exists := s.cache.Get(strconv.Itoa(slot.ID))
	if !exists {
		subs, err := s.subService.FindSubscriptionsBySlotInfo(ctx, slot)
		if err != nil {
			slog.Error("Failed to find subscriptions for slot", "data", slot, "error", err, "service", logger.ServiceNotification)
		}

		for _, sub := range subs {
			notif := model.Notification{UserID: sub.UserID, Slot: slot}
			slog.Info("Sending notification", "data", notif, "service", logger.ServiceNotification)
			event.Publish(ctx, s.eventBus, &notif)
		}
	}

	s.cache.Set(strconv.Itoa(slot.ID), slot)
}
