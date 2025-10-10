package subscription

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Ademun/mining-lab-bot/pkg/model"
)

type SubscriptionRepo interface {
	Create(ctx context.Context, sub model.Subscription) error
	Delete(ctx context.Context, UUID string) error
	List(ctx context.Context) ([]model.Subscription, error)
}

type subscriptionRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) (SubscriptionRepo, error) {
	query := `
create table if not exists subscriptions (
    uuid text not null primary key,
    user_id integer not null,
    lab_number integer not null,
    lab_auditorium integer not null,    
)`
	_, err := db.Exec(query)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscriptions table: %w", err)
	}
	return &subscriptionRepo{db: db}, nil
}

func (s *subscriptionRepo) Create(ctx context.Context, sub model.Subscription) error {
	query := `insert into subscriptions (uuid, user_id, lab_number, lab_auditorium) values (?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query, sub.UUID, sub.UserID, sub.LabNumber, sub.LabAuditorium)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}
	return nil
}

func (s *subscriptionRepo) Delete(ctx context.Context, UUID string) error {
	query := `delete from subscriptions where uuid = ?`
	_, err := s.db.ExecContext(ctx, query, UUID)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}
	return nil
}

func (s *subscriptionRepo) List(ctx context.Context) ([]model.Subscription, error) {
	query := `select uuid, user_id, lab_number, lab_auditorium from subscriptions`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	defer rows.Close()

	var subs []model.Subscription
	for rows.Next() {
		var sub model.Subscription
		err := rows.Scan(&sub.UUID, &sub.UserID, &sub.LabNumber, &sub.LabAuditorium)
		if err != nil {
			return nil, fmt.Errorf("failed to list subscriptions: %w", err)
		}
		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}

	return subs, nil
}
