package subscription

import (
	"context"

	"github.com/Ademun/mining-lab-bot/pkg/errs"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repo interface {
	Create(ctx context.Context, subReq RequestSubscription) error
	Delete(ctx context.Context, uuid uuid.UUID) (bool, error)
	Find(ctx context.Context, subFilters SubFilters, timeFilters TimeFilters) ([]ResponseSubscription, error)
}

type subscriptionRepo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repo {
	return &subscriptionRepo{db: db}
}

func (s *subscriptionRepo) Create(ctx context.Context, subReq RequestSubscription) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return errs.ErrBeginTransaction
	}
	defer tx.Rollback()
	sub, subTimes := subReq.toDBModels()

	subInsert := `
insert into subscriptions 
(uuid, user_id, lab_type, lab_number, lab_auditorium, lab_domain, weekday) 
values 
(:uuid, :user_id, :lab_type, :lab_number, :lab_auditorium, :lab_domain, :weekday)`
	_, err = tx.NamedExecContext(ctx, subInsert, sub)
	if err != nil {
		return &errs.ErrQueryExecution{Operation: "Create", Query: subInsert, Err: err}
	}

	if len(subTimes) == 0 {
		return tx.Commit()
	}

	timesInsert := `
insert into subscription_times 
(subscription_uuid, time_start, time_end) 
values 
(:subscription_uuid, :time_start, :time_end)
`
	if _, err = tx.NamedExecContext(ctx, timesInsert, subTimes); err != nil {
		return &errs.ErrQueryExecution{Operation: "Create", Query: timesInsert, Err: err}
	}

	return tx.Commit()
}

func (s *subscriptionRepo) Delete(ctx context.Context, uuid uuid.UUID) (bool, error) {
	query := `delete from subscriptions where uuid = ?`
	res, err := s.db.ExecContext(ctx, query, uuid.String())
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

func (s *subscriptionRepo) Find(ctx context.Context, subFilters SubFilters, timeFilters TimeFilters) ([]ResponseSubscription, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, errs.ErrBeginTransaction
	}
	defer tx.Rollback()
	query, args, err := subFilters.buildQuery()
	if err != nil {
		return nil, &errs.ErrQueryCreation{Operation: "Find", Query: query, Err: err}
	}
	var subs []DBSubscription
	if err = tx.SelectContext(ctx, &subs, query, args...); err != nil {
		return nil, &errs.ErrQueryExecution{Operation: "Find", Query: query, Err: err}
	}

	response, err := s.convertDBSubsToResponse(ctx, tx, subs, timeFilters)
	if err != nil {
		return nil, &errs.ErrQueryExecution{Operation: "Find", Query: query, Err: err}
	}

	return response, tx.Commit()
}

func (s *subscriptionRepo) convertDBSubsToResponse(ctx context.Context, tx *sqlx.Tx, subs []DBSubscription, timeFilters TimeFilters) ([]ResponseSubscription, error) {
	subUUIDs := make([]uuid.UUID, len(subs))
	for idx, sub := range subs {
		subUUIDs[idx] = sub.UUID
	}

	timeFilters.SubUUIDs = subUUIDs
	prefTimes, err := s.findPreferredTimes(ctx, tx, timeFilters)
	if err != nil {
		return nil, err
	}

	subTimes := make(map[uuid.UUID][]DBSubscriptionTimes, len(subs))
	for _, time := range prefTimes {
		subTimes[time.SubscriptionUUID] = append(subTimes[time.SubscriptionUUID], time)
	}

	response := make([]ResponseSubscription, len(subs))
	for idx, sub := range subs {
		response[idx] = toResponse(sub, subTimes[sub.UUID])
	}
	return response, nil
}

func (s *subscriptionRepo) findPreferredTimes(ctx context.Context, tx *sqlx.Tx, timeFilters TimeFilters) ([]DBSubscriptionTimes, error) {
	query, args, err := timeFilters.buildQuery()
	if err != nil {
		return nil, &errs.ErrQueryCreation{Operation: "findPreferredTimes", Query: query, Err: err}
	}
	var times []DBSubscriptionTimes
	if err = tx.SelectContext(ctx, &times, query, args...); err != nil {
		return nil, &errs.ErrQueryExecution{Operation: "findPreferredTimes", Query: query, Err: err}
	}
	return times, nil
}
