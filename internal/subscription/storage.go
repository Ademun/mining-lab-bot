package subscription

import (
	"context"

	"github.com/Ademun/mining-lab-bot/pkg/errs"
	"github.com/Ademun/mining-lab-bot/pkg/model"
	"github.com/jmoiron/sqlx"
)

type Repo interface {
	Create(ctx context.Context, sub model.Subscription) error
	Delete(ctx context.Context, UUID string) (bool, error)
	FindByUserID(ctx context.Context, userID int) ([]model.Subscription, error)
	FindBySlotInfo(ctx context.Context, labNumber, labAuditorium int) ([]model.Subscription, error)
	Count(ctx context.Context) (int, error)
}

type subscriptionRepo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repo {
	return &subscriptionRepo{db: db}
}

func (s *subscriptionRepo) Create(ctx context.Context, sub model.Subscription) error {
	query := `insert into subscriptions (uuid, user_id, chat_id, lab_number, lab_auditorium, weekday, day_time) values (:uuid, :user_id, :chat_id, :lab_number, :lab_auditorium, :weekday, :day_time)`
	_, err := s.db.NamedExecContext(ctx, query, sub)
	if err != nil {
		return &errs.ErrQueryExecution{Operation: "Create", Query: query, Err: err}
	}
	return nil
}

func (s *subscriptionRepo) Delete(ctx context.Context, uuid string) (bool, error) {
	query := `delete from subscriptions where uuid = ?`
	res, err := s.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return false, &errs.ErrQueryExecution{Operation: "Delete", Query: query, Err: err}
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return false, &errs.ErrQueryExecution{Operation: "Delete", Query: query, Err: err}
	}
	if affected == 0 {
		return false, nil
	}
	return true, nil
}

func (s *subscriptionRepo) FindByUserID(ctx context.Context, userID int) ([]model.Subscription, error) {
	query := `select uuid, user_id, chat_id, lab_number, lab_auditorium, weekday, day_time from subscriptions where user_id = ?`
	var subs []model.Subscription
	err := s.db.SelectContext(ctx, &subs, query, userID)
	if err != nil {
		return nil, &errs.ErrQueryExecution{Operation: "FindByUserID", Query: query, Err: err}
	}
	return subs, nil
}

func (s *subscriptionRepo) FindBySlotInfo(ctx context.Context, labNumber, labAuditorium int) ([]model.Subscription, error) {
	query := `select uuid, user_id, chat_id, lab_number, lab_auditorium, weekday, day_time from subscriptions where lab_number = ? and lab_auditorium = ?`
	var subs []model.Subscription
	err := s.db.SelectContext(ctx, &subs, query, labNumber, labAuditorium)
	if err != nil {
		return nil, &errs.ErrQueryExecution{Operation: "FindBySlotInfo", Query: query, Err: err}
	}
	return subs, nil
}

func (s *subscriptionRepo) Count(ctx context.Context) (int, error) {
	query := `select count(*) from subscriptions`
	var count int
	err := s.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, &errs.ErrQueryExecution{Operation: "Count", Query: query, Err: err}
	}
	return count, nil
}
