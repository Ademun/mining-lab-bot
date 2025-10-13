package subscription

import (
	"context"
	"log/slog"

	"github.com/Ademun/mining-lab-bot/pkg/event"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/Ademun/mining-lab-bot/pkg/model"
)

type SubscriptionService interface {
	Subscribe(ctx context.Context, sub model.Subscription) error
	Unsubscribe(ctx context.Context, subUUID string) error
	FindSubscriptionsByChatID(ctx context.Context, chatID int) ([]model.Subscription, error)
	FindSubscriptionsBySlotInfo(ctx context.Context, slot model.Slot) ([]model.Subscription, error)
}

type subscriptionService struct {
	eventBus *event.Bus
	subRepo  SubscriptionRepo
}

func New(eb *event.Bus, repo SubscriptionRepo) SubscriptionService {
	return &subscriptionService{
		eventBus: eb,
		subRepo:  repo,
	}
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

	return s.subRepo.Create(ctx, sub)
}

func (s *subscriptionService) Unsubscribe(ctx context.Context, subUUID string) error {
	slog.Info("Deleting subscription", "uuid", subUUID, "service", logger.ServiceSubscription)
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
