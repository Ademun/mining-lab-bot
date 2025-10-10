package subscription

import (
	"context"
	"log/slog"

	"github.com/Ademun/mining-lab-bot/pkg/event"
	"github.com/Ademun/mining-lab-bot/pkg/model"
)

type SubscriptionService interface {
	Start(ctx context.Context) error
	Subscribe(ctx context.Context, sub model.Subscription) error
	Unsubscribe(ctx context.Context, subUUID string) error
	List(ctx context.Context) ([]model.Subscription, error)
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
	slog.Info("[SubscriptionService] Starting...]")
	slog.Info("[SubscriptionService] Started")
	return nil
}

func (s *subscriptionService) Subscribe(ctx context.Context, sub model.Subscription) error {
	slog.Info("[SubscriptionService] New subscription")
	return s.subRepo.Create(ctx, sub)
}

func (s *subscriptionService) Unsubscribe(ctx context.Context, subUUID string) error {
	slog.Info("[SubscriptionService] Removing subscription")
	return s.subRepo.Delete(ctx, subUUID)
}

func (s *subscriptionService) List(ctx context.Context) ([]model.Subscription, error) {
	slog.Info("[SubscriptionService] Listing subscriptions")
	return s.subRepo.List(ctx)
}
