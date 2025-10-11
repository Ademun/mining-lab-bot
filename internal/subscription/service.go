package subscription

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Ademun/mining-lab-bot/pkg/event"
	"github.com/Ademun/mining-lab-bot/pkg/model"
)

type SubscriptionService interface {
	Start(ctx context.Context) error
	Subscribe(ctx context.Context, sub model.Subscription) error
	Unsubscribe(ctx context.Context, subUUID string) error
	ListForUser(ctx context.Context, userID int) ([]model.Subscription, error)
}

type subscriptionService struct {
	eb      *event.Bus
	subRepo SubscriptionRepo
}

func New(eb *event.Bus, repo SubscriptionRepo) SubscriptionService {
	return &subscriptionService{
		eb:      eb,
		subRepo: repo,
	}
}

func (s *subscriptionService) Start(ctx context.Context) error {
	slog.Info("[SubscriptionService] Starting...")
	slog.Info("[SubscriptionService] Started")
	return nil
}

func (s *subscriptionService) Subscribe(ctx context.Context, sub model.Subscription) error {
	slog.Info("[SubscriptionService] New subscription")

	exists, err := s.subRepo.Exists(ctx, sub.UserID, sub.LabNumber, sub.LabAuditorium)
	if err != nil {
		slog.Error("[SubscriptionService] Error checking if subscription exists: ", err)
		return err
	}

	if exists {
		slog.Info("[SubscriptionService] Subscription already exists")
		return errors.New("вы уже подписаны на эту лабу")
	}

	return s.subRepo.Create(ctx, sub)
}

func (s *subscriptionService) Unsubscribe(ctx context.Context, subUUID string) error {
	slog.Info("[SubscriptionService] Removing subscription")
	return s.subRepo.Delete(ctx, subUUID)
}

func (s *subscriptionService) ListForUser(ctx context.Context, userID int) ([]model.Subscription, error) {
	slog.Info("[SubscriptionService] Listing subscriptions")
	return s.subRepo.ListForUser(ctx, userID)
}
