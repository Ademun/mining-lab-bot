package notification

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Ademun/mining-lab-bot/internal/subscription"
	"github.com/Ademun/mining-lab-bot/pkg/cache"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/Ademun/mining-lab-bot/pkg/metrics"
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

	prefTimes := getSubscriptionPreferredTimes(subs...)
	userIDsMap := make(map[int]struct{})
	for _, sub := range subs {
		if _, ok := userIDsMap[sub.UserID]; ok {
			continue
		}
		userIDsMap[sub.UserID] = struct{}{}
	}
	for userID := range userIDsMap {
		notif := model.Notification{
			UserID:         userID,
			PreferredTimes: prefTimes[userID],
			Slot:           slot,
		}
		s.notifier.SendNotification(ctx, notif)
	}
	metrics.Global().RecordNotificationResults(len(userIDsMap), len(s.cache.List()))

	slog.Info("Finished sending notifications", "total", len(subs), "slot", slot, "service", logger.ServiceNotification)
}

func (s *notificationService) NotifyNewSubscription(ctx context.Context, sub model.Subscription) {
	slots := s.findSlotsBySubscriptionInfo(sub)

	prefTimes := getSubscriptionPreferredTimes(sub)
	for _, slot := range slots {
		notif := model.Notification{
			UserID:         sub.UserID,
			PreferredTimes: prefTimes[sub.UserID],
			Slot:           slot,
		}
		s.notifier.SendNotification(ctx, notif)
	}
	metrics.Global().RecordNotificationResults(len(slots), len(s.cache.List()))

	slog.Info("Finished sending notifications", "total", len(slots), "sub", sub, "service", logger.ServiceNotification)
}

func getSubscriptionPreferredTimes(subs ...model.Subscription) map[int][]model.PreferredTime {
	userPrefTimes := make(map[int][]model.PreferredTime)
	for _, sub := range subs {
		userID := sub.UserID
		var prefTime model.PreferredTime
		if sub.Weekday != nil && sub.DayTime != nil {
			prefTime.Weekday = *sub.Weekday
			prefTime.DayTime = *sub.DayTime
		}
		if _, exists := userPrefTimes[userID]; !exists {
			userPrefTimes[userID] = []model.PreferredTime{prefTime}
			continue
		}
		userPrefTimes[userID] = append(userPrefTimes[userID], prefTime)
	}
	return userPrefTimes
}

func (s *notificationService) findSlotsBySubscriptionInfo(sub model.Subscription) []model.Slot {
	slots := s.cache.List()
	items := make([]model.Slot, 0)
	for _, slot := range slots {
		if slot.LabNumber != sub.LabNumber || slot.LabAuditorium != sub.LabAuditorium {
			continue
		}
		if sub.Weekday == nil || sub.DayTime == nil {
			items = append(items, slot)
		}
		prefTime := fmt.Sprintf("%d-%s", *sub.Weekday, *sub.DayTime)
		for _, available := range slot.Available {
			slotTime := fmt.Sprintf("%d-%s", available.Time.Weekday(), available.Time.Format("15:04"))
			if prefTime == slotTime {
				items = append(items, slot)
			}
		}
	}
	return items
}
