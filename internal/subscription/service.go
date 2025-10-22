package subscription

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Ademun/mining-lab-bot/pkg/errs"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/Ademun/mining-lab-bot/pkg/metrics"
	"github.com/Ademun/mining-lab-bot/pkg/model"
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
	exists, err := s.subRepo.Exists(ctx, sub.UserID, sub.LabNumber, sub.LabAuditorium)
	if err != nil {
		slog.Error("Failed to check if subscription exists", "sub", sub, "err", err)
		return err
	}

	if exists {
		return errs.ErrSubscriptionExists
	}

	metrics.Global().RecordSubscriptionResults(1)

	return s.subRepo.Create(ctx, sub)
}

func (s *subscriptionService) Unsubscribe(ctx context.Context, subUUID string) error {
	metrics.Global().RecordSubscriptionResults(-1)

	err := s.subRepo.Delete(ctx, subUUID)
	if err != nil {
		slog.Error("Failed to delete subscription", "subUUID", subUUID, "err", err)
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
	fmt.Println("Subs", len(subs))
	if err != nil {
		slog.Error("Failed to find subscriptions", "slot", slot, "err", err)
	}

	res := make([]model.Subscription, 0)
	for _, sub := range subs {
		for _, available := range slot.Available {
			day := available.Time.Weekday()
			dayTime := available.Time.Format("15:04")

			if sub.Weekday == nil {
				res = append(res, sub)
				break
			}

			if day == *sub.Weekday && dayTime == *sub.DayTime {
				res = append(res, sub)
				break
			}
		}
	}

	return res, err
}
