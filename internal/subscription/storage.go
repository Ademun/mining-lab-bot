package subscription

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Ademun/mining-lab-bot/pkg/errs"
	"github.com/Ademun/mining-lab-bot/pkg/model"
)

type Repo interface {
	Create(ctx context.Context, sub model.Subscription) error
	Delete(ctx context.Context, UUID string) error
	Exists(ctx context.Context, userID, labNumber, labAuditorium int) (bool, error)
	FindByUserID(ctx context.Context, userID int) ([]model.Subscription, error)
	FindBySlotInfo(ctx context.Context, labNumber, labAuditorium int) ([]model.Subscription, error)
	Count(ctx context.Context) (int, error)
}

type subscriptionRepo struct {
	db *sql.DB
}

func NewRepo(ctx context.Context, db *sql.DB) (Repo, error) {
	query := `
create table if not exists subscriptions (
    uuid text not null primary key,
    user_id integer not null,
    chat_id integer not null,
    lab_number integer not null,
    lab_auditorium integer not null
    weekday integer,
    day_time text,
)`
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return nil, &errs.ErrQueryExecution{Operation: "Init", Query: query, Err: err}
	}

	return &subscriptionRepo{db: db}, nil
}

func (s *subscriptionRepo) Create(ctx context.Context, sub model.Subscription) error {
	query := `insert into subscriptions (uuid, user_id, chat_id, lab_number, lab_auditorium, weekday, day_time, teacher) values (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.ExecContext(ctx, query, sub.UUID, sub.UserID, sub.ChatID, sub.LabNumber, sub.LabAuditorium, sub.Weekday, sub.DayTime, sub.Teacher)
	if err != nil {
		return &errs.ErrQueryExecution{Operation: "Create", Query: query, Err: err}
	}

	return nil
}

func (s *subscriptionRepo) Delete(ctx context.Context, uuid string) error {
	query := `delete from subscriptions where uuid = ?`

	_, err := s.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return &errs.ErrQueryExecution{Operation: "Delete", Query: query, Err: err}
	}

	return nil
}

func (s *subscriptionRepo) Exists(ctx context.Context, userID, labNumber, labAuditorium int) (bool, error) {
	query := `select exists (select 1 from subscriptions where user_id = ? and lab_number = ? and lab_auditorium = ?)`

	var exists bool
	err := s.db.QueryRowContext(ctx, query, userID, labNumber, labAuditorium).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, &errs.ErrQueryExecution{Operation: "Exists", Query: query, Err: err}
	}

	return exists, nil
}

func (s *subscriptionRepo) FindByUserID(ctx context.Context, userID int) ([]model.Subscription, error) {
	query := `select uuid, user_id, chat_id, lab_number, lab_auditorium, weekday, day_time, teacher from subscriptions where user_id = ?`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, &errs.ErrQueryExecution{Operation: "FindByUserID", Query: query, Err: err}
	}
	defer rows.Close()

	var subs []model.Subscription
	for rows.Next() {
		var sub model.Subscription
		err := rows.Scan(&sub.UUID, &sub.UserID, &sub.ChatID, &sub.LabNumber, &sub.LabAuditorium, &sub.Weekday, &sub.DayTime, &sub.Teacher)
		if err != nil {
			return nil, &errs.ErrRowIteration{Operation: "FindByUserID", Query: query, Err: err}
		}
		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, &errs.ErrRowIteration{Operation: "FindByUserID", Query: query, Err: err}
	}

	return subs, nil
}

func (s *subscriptionRepo) FindBySlotInfo(ctx context.Context, labNumber, labAuditorium int) ([]model.Subscription, error) {
	query := `select uuid, user_id, chat_id, lab_number, lab_auditorium, weekday, day_time, teacher from subscriptions where lab_number = ? and lab_auditorium = ?`
	rows, err := s.db.QueryContext(ctx, query, labNumber, labAuditorium)
	if err != nil {
		return nil, &errs.ErrQueryExecution{Operation: "FindBySlotInfo", Query: query, Err: err}
	}
	defer rows.Close()

	var subs []model.Subscription
	for rows.Next() {
		var sub model.Subscription
		err := rows.Scan(&sub.UUID, &sub.UserID, &sub.ChatID, &sub.LabNumber, &sub.LabAuditorium, &sub.Weekday, &sub.DayTime)
		if err != nil {
			return nil, &errs.ErrRowIteration{Operation: "FindBySlotInfo", Query: query, Err: err}
		}
		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, &errs.ErrRowIteration{Operation: "FindBySlotInfo", Query: query, Err: err}
	}

	return subs, nil
}

func (s *subscriptionRepo) Count(ctx context.Context) (int, error) {
	query := `select count(*) from subscriptions`

	var count int
	err := s.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return -1, &errs.ErrQueryExecution{Operation: "Count", Query: query, Err: err}
	}

	return count, nil
}
