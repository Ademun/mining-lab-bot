package subscription

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Ademun/mining-lab-bot/internal/polling"
	"github.com/Ademun/mining-lab-bot/pkg/errs"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
)

type Service interface {
	Subscribe(ctx context.Context, sub RequestSubscription) error
	Unsubscribe(ctx context.Context, subUUID uuid.UUID) error
	FindSubscriptionsByUserID(ctx context.Context, userID int) ([]ResponseSubscription, error)
	FindUsersBySlotInfo(ctx context.Context, slot polling.Slot) ([]ResponseUser, error)
}

type subscriptionService struct {
	subRepo Repo
}

func New(repo Repo) Service {
	return &subscriptionService{
		subRepo: repo,
	}
}

func (s *subscriptionService) Subscribe(ctx context.Context, sub RequestSubscription) error {
	err := s.subRepo.Create(ctx, sub)
	if err != nil {
		if isDuplicateError(err) {
			return errs.ErrSubscriptionExists
		}
		slog.Error("Failed to create subscription", "sub", sub, "err", err)
		return err
	}
	return nil
}

func isDuplicateError(err error) bool {
	var queryErr *errs.ErrQueryExecution
	if errors.As(err, &queryErr) {
		var sqliteErr sqlite3.Error
		if errors.As(queryErr.Err, &sqliteErr) {
			return errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique)
		}
	}
	return false
}

func (s *subscriptionService) Unsubscribe(ctx context.Context, subUUID uuid.UUID) error {
	_, err := s.subRepo.Delete(ctx, subUUID)
	if err != nil {
		slog.Error("Failed to delete subscription", "subUUID", subUUID, "err", err)
	}
	return err
}

func (s *subscriptionService) FindSubscriptionsByUserID(ctx context.Context, userID int) ([]ResponseSubscription, error) {
	subFilters, timeFilters := SubFilters{UserID: userID}, TimeFilters{}
	subs, err := s.subRepo.Find(ctx, subFilters, timeFilters)
	if err != nil {
		slog.Error("Failed to find subscriptions", "userID", userID, "err", err)
	}

	return subs, err
}

func (s *subscriptionService) FindUsersBySlotInfo(ctx context.Context, slot polling.Slot) ([]ResponseUser, error) {
	weekdays := make([]int, 0, len(slot.TimesTeachers))
	for t := range slot.TimesTeachers {
		weekdays = append(weekdays, int(t.Weekday()))
	}
	subFilters := SubFilters{
		Type:          &slot.Type,
		LabNumber:     slot.Number,
		LabAuditorium: slot.Auditorium,
		Weekdays:      weekdays,
	}
	if slot.Type == polling.LabTypeDefence {
		subFilters.LabDomain = &slot.Domain
	}
	times := make([]string, 0, len(slot.TimesTeachers))
	for t := range slot.TimesTeachers {
		times = append(times, t.Format("15:04"))
	}
	timeFilters := TimeFilters{
		Includes: times,
	}
	subs, err := s.subRepo.Find(ctx, subFilters, timeFilters)
	if err != nil {
		slog.Error("Failed to find subscriptions", "slot", slot, "err", err)
	}

	userIDSubs := make(map[int][]ResponseSubscription)
	for _, sub := range subs {
		userIDSubs[sub.UserID] = append(userIDSubs[sub.UserID], sub)
	}

	users := make([]ResponseUser, 0, len(userIDSubs))
	for userID, userSubs := range userIDSubs {
		prefTimes := make(map[time.Weekday][]TimeRange)
		for _, sub := range userSubs {
			if sub.Weekday == nil {
				continue
			}
			prefTimes[time.Weekday(*sub.Weekday)] = append(prefTimes[time.Weekday(*sub.Weekday)], sub.PreferredTimes...)
		}
		users = append(users, ResponseUser{
			UserID:         userID,
			PreferredTimes: prefTimes,
		})
	}

	return users, err
}
