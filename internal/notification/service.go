package notification

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/Ademun/mining-lab-bot/internal/subscription"
	"github.com/Ademun/mining-lab-bot/pkg/cache"
	"github.com/Ademun/mining-lab-bot/pkg/event"
	"github.com/Ademun/mining-lab-bot/pkg/model"
)

type NotificationService interface {
	Start(ctx context.Context) error
}

type notificationService struct {
	eb         *event.Bus
	subService subscription.SubscriptionService
	cache      *cache.TTLCache[model.Slot]
}

func New(eb *event.Bus, subService subscription.SubscriptionService) NotificationService {
	return &notificationService{
		eb:         eb,
		subService: subService,
		cache:      cache.NewTTLCache[model.Slot](time.Minute*5, time.Minute*10),
	}
}

func (s *notificationService) Start(ctx context.Context) error {
	slog.Info("[Notification service] Starting...")
	event.Subscribe(s.eb, s.handleNewSlot)
	slog.Info("[Notification service] Started")
	return nil
}

func (s *notificationService) handleNewSlot(ctx context.Context, slot model.Slot) {
	_, exists := s.cache.Get(strconv.Itoa(slot.ID))
	if !exists {
		subs, err := s.subService.FindSubscriptionsBySlotInfo(ctx, slot)
		if err != nil {
			slog.Error("[Notification service] Failed to handle new slot", err)
		}

		for _, sub := range subs {
			notif := model.Notification{UserID: sub.UserID, Slot: slot}
			event.Publish(s.eb, ctx, &notif)
		}
	}

	s.cache.Set(strconv.Itoa(slot.ID), slot)
}
