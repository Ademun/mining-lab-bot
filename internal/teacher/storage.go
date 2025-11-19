package teacher

import (
	"context"

	"github.com/Ademun/mining-lab-bot/pkg/errs"
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source=storage.go -destination=mocks/mock_storage.go -package=mocks
type Repo interface {
	FindBySchedule(ctx context.Context, filter Filter) ([]Teacher, error)
}

type teacherRepo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repo {
	return &teacherRepo{db: db}
}

func (t *teacherRepo) FindBySchedule(ctx context.Context, filter Filter) ([]Teacher, error) {
	query, args, err := filter.buildQuery()
	if err != nil {
		return nil, &errs.ErrQueryExecution{Operation: "FindBySchedule", Query: query, Err: err}
	}
	var teachers []Teacher
	if err := t.db.SelectContext(ctx, &teachers, query, args...); err != nil {
		return nil, &errs.ErrQueryExecution{Operation: "FindBySchedule", Query: query, Err: err}
	}
	return teachers, nil
}
