package subscription

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Ademun/mining-lab-bot/pkg/model"
)

type Repo interface {
	Create(ctx context.Context, sub model.Subscription) error
	Delete(ctx context.Context, UUID string) error
	List(ctx context.Context) ([]model.Subscription, error)
	FindByChatID(ctx context.Context, chatID int) ([]model.Subscription, error)
	FindBySlotInfo(ctx context.Context, labNumber, labAuditorium int) ([]model.Subscription, error)
	Exists(ctx context.Context, userID, labNumber, labAuditorium int) (bool, error)
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
)`
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return nil, &ErrQueryExecution{"Init", query, err}
	}
	return &subscriptionRepo{db: db}, nil
}

func (s *subscriptionRepo) Create(ctx context.Context, sub model.Subscription) error {
	query := `insert into subscriptions (uuid, user_id, chat_id, lab_number, lab_auditorium) values (?, ?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, sub.UUID, sub.UserID, sub.ChatID, sub.LabNumber, sub.LabAuditorium)
	if err != nil {
		return &ErrQueryExecution{"Create", query, err}
	}
	return nil
}

func (s *subscriptionRepo) Delete(ctx context.Context, uuid string) error {
	query := `delete from subscriptions where uuid = ?`
	_, err := s.db.ExecContext(ctx, query, uuid)
	if err != nil {
		return &ErrQueryExecution{"Delete", query, err}
	}
	return nil
}

func (s *subscriptionRepo) List(ctx context.Context) ([]model.Subscription, error) {
	query := `select uuid, user_id, chat_id, lab_number, lab_auditorium from subscriptions`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, &ErrQueryExecution{"List", query, err}
	}
	defer rows.Close()

	var subs []model.Subscription
	for rows.Next() {
		var sub model.Subscription
		err = rows.Scan(&sub.UUID, &sub.UserID, &sub.ChatID, &sub.LabNumber, &sub.LabAuditorium)
		if err != nil {
			return nil, &ErrRowIteration{"List", query, err}
		}
		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, &ErrRowIteration{"List", query, err}
	}

	return subs, nil
}

func (s *subscriptionRepo) FindByChatID(ctx context.Context, chatID int) ([]model.Subscription, error) {
	query := `select uuid, user_id, chat_id, lab_number, lab_auditorium from subscriptions where chat_id = ?`
	rows, err := s.db.QueryContext(ctx, query, chatID)
	if err != nil {
		return nil, &ErrQueryExecution{"FindByUserID", query, err}
	}
	defer rows.Close()

	var subs []model.Subscription
	for rows.Next() {
		var sub model.Subscription
		err := rows.Scan(&sub.UUID, &sub.UserID, &sub.ChatID, &sub.LabNumber, &sub.LabAuditorium)
		if err != nil {
			return nil, &ErrRowIteration{"FindByUserID", query, err}
		}
		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, &ErrRowIteration{"FindByUserID", query, err}
	}

	return subs, nil
}

func (s *subscriptionRepo) FindBySlotInfo(ctx context.Context, labNumber, labAuditorium int) ([]model.Subscription, error) {
	query := `select uuid, user_id, chat_id, lab_number, lab_auditorium from subscriptions where lab_number = ? and lab_auditorium = ?`
	rows, err := s.db.QueryContext(ctx, query, labNumber, labAuditorium)
	if err != nil {
		return nil, &ErrQueryExecution{"FindBySlotInfo", query, err}
	}
	defer rows.Close()

	var subs []model.Subscription
	for rows.Next() {
		var sub model.Subscription
		err := rows.Scan(&sub.UUID, &sub.UserID, &sub.ChatID, &sub.LabNumber, &sub.LabAuditorium)
		if err != nil {
			return nil, &ErrRowIteration{"FindBySlotInfo", query, err}
		}
		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, &ErrRowIteration{"FindBySlotInfo", query, err}
	}

	return subs, nil
}

func (s *subscriptionRepo) Exists(ctx context.Context, userID, labNumber, labAuditorium int) (bool, error) {
	query := `select exists (select 1 from subscriptions where user_id = ? and lab_number = ? and lab_auditorium = ?)`

	var exists bool
	err := s.db.QueryRowContext(ctx, query, userID, labNumber, labAuditorium).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, &ErrQueryExecution{"Exists", query, err}
	}

	return exists, nil
}
