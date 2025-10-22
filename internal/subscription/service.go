package subscription

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Ademun/mining-lab-bot/pkg/errs"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/Ademun/mining-lab-bot/pkg/metrics"
	"github.com/Ademun/mining-lab-bot/pkg/model"
	"github.com/mattn/go-sqlite3"
)

type Service interface {
	Start(ctx context.Context) error
	Subscribe(ctx context.Context, sub model.Subscription) error
	Unsubscribe(ctx context.Context, subUUID string) error
	FindSubscriptionsByUserID(ctx context.Context, chatID int) ([]model.Subscription, error)
	FindSubscriptionsBySlotInfo(ctx context.Context, slot model.Slot) ([]model.Subscription, error)
}

type subscriptionService struct {
	subRepo Repo
}

func New(repo Repo) Service {
	return &subscriptionService{
		subRepo: repo,
	}
}

func (s *subscriptionService) Start(ctx context.Context) error {
	slog.Info("Starting", "service", logger.ServiceSubscription)
	subs, err := s.subRepo.Count(ctx)
	if err != nil {
		return err
	}
	metrics.Global().RecordSubscriptionResults(subs)
	slog.Info("Started", "service", logger.ServiceSubscription)
	return nil
}

func (s *subscriptionService) Subscribe(ctx context.Context, sub model.Subscription) error {
	err := s.subRepo.Create(ctx, sub)
	if err != nil {
		if errors.Is(err, sqlite3.ErrConstraintUnique) {
			return errs.ErrSubscriptionExists
		}
		slog.Error("Failed to create subscription", "sub", sub, "err", err)
		return err
	}
	metrics.Global().RecordSubscriptionResults(1)
	return nil
}

func (s *subscriptionService) Unsubscribe(ctx context.Context, subUUID string) error {
	success, err := s.subRepo.Delete(ctx, subUUID)
	if err != nil {
		slog.Error("Failed to delete subscription", "subUUID", subUUID, "err", err)
	}

	if success {
		metrics.Global().RecordSubscriptionResults(-1)
	}

	return err
}

func (s *subscriptionService) FindSubscriptionsByUserID(ctx context.Context, userID int) ([]model.Subscription, error) {
	subs, err := s.subRepo.FindByUserID(ctx, userID)
	if err != nil {
		slog.Error("Failed to find subscriptions", "userID", userID, "err", err)
	}

	return subs, err
}

func (s *subscriptionService) FindSubscriptionsBySlotInfo(ctx context.Context, slot model.Slot) ([]model.Subscription, error) {
	subs, err := s.subRepo.FindBySlotInfo(ctx, slot.LabNumber, slot.LabAuditorium)
	if err != nil {
		slog.Error("Failed to find subscriptions", "slot", slot, "err", err)
	}

	res := make([]model.Subscription, 0)
	availableMap := make(map[string]bool)
	for _, available := range slot.Available {
		key := fmt.Sprintf("%d-%s", available.Time.Weekday(), available.Time.Format("15:04"))
		availableMap[key] = true
	}

	for _, sub := range subs {
		if sub.Weekday == nil || sub.DayTime == nil {
			res = append(res, sub)
			continue
		}

		key := fmt.Sprintf("%d-%s", *sub.Weekday, *sub.DayTime)
		if availableMap[key] {
			res = append(res, sub)
		}
	}

	return res, err
}
