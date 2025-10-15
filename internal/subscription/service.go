package subscription

import (
	"context"
	"log/slog"

	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/Ademun/mining-lab-bot/pkg/metrics"
	"github.com/Ademun/mining-lab-bot/pkg/model"
)

type Service interface {
	Start(ctx context.Context) error
	Subscribe(ctx context.Context, sub model.Subscription) error
	Unsubscribe(ctx context.Context, subUUID string) error
	FindSubscriptionsByChatID(ctx context.Context, chatID int) ([]model.Subscription, error)
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
	subs, err := s.subRepo.List(ctx)
	if err != nil {
		return err
	}
	metrics.Global().RecordSubscriptionResults(len(subs))
	slog.Info("Started", "service", logger.ServiceSubscription)
	return nil
}

func (s *subscriptionService) Subscribe(ctx context.Context, sub model.Subscription) error {
	slog.Info("Creating new subscription", "data", sub, "service", logger.ServiceSubscription)
	exists, err := s.subRepo.Exists(ctx, sub.UserID, sub.LabNumber, sub.LabAuditorium)
	if err != nil {
		return err
	}

	if exists {
		slog.Warn("Subscription already exists", "data", sub, "service", logger.ServiceSubscription)
		return ErrSubscriptionExists
	}

	metrics.Global().RecordSubscriptionResults(1)

	return s.subRepo.Create(ctx, sub)
}

func (s *subscriptionService) Unsubscribe(ctx context.Context, subUUID string) error {
	slog.Info("Deleting subscription", "uuid", subUUID, "service", logger.ServiceSubscription)

	metrics.Global().RecordSubscriptionResults(-1)

	return s.subRepo.Delete(ctx, subUUID)
}

func (s *subscriptionService) FindSubscriptionsByChatID(ctx context.Context, chatID int) ([]model.Subscription, error) {
	slog.Info("Finding subscriptions by chat ID", "chatID", chatID, "service", logger.ServiceSubscription)
	return s.subRepo.FindByChatID(ctx, chatID)
}

func (s *subscriptionService) FindSubscriptionsBySlotInfo(ctx context.Context, slot model.Slot) ([]model.Subscription, error) {
	slog.Info("Finding subscriptions by slot info", "data", slot, "service", logger.ServiceSubscription)
	return s.subRepo.FindBySlotInfo(ctx, slot.LabNumber, slot.LabAuditorium)
}
