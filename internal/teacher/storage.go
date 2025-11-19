package teacher

import (
	"context"
	"fmt"
	"time"

	"github.com/Ademun/mining-lab-bot/pkg/errs"
	"github.com/jmoiron/sqlx"
)

type Repo interface {
	FindByWeekNumberWeekdayAuditorium(ctx context.Context, weekNumber int, weekday time.Weekday, auditorium int) ([]Teacher, error)
}

type teacherRepo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repo {
	return &teacherRepo{db: db}
}

func (t *teacherRepo) FindByWeekNumberWeekdayAuditorium(ctx context.Context, weekNumber int, weekday time.Weekday, auditorium int) ([]Teacher, error) {
	fmt.Println(weekNumber, weekday, auditorium)
	query := `select name, auditorium, week_number, weekday, time_start, time_end from teachers where week_number = :week_number and weekday = :weekday and auditorium = :auditorium`
	var teachers []Teacher
	err := t.db.SelectContext(ctx, &teachers, query, weekNumber, weekday, auditorium)
	if err != nil {
		return nil, &errs.ErrQueryExecution{Operation: "FindByWeekNumberWeekdayAuditorium", Query: query, Err: err}
	}
	for _, teacher := range teachers {
		fmt.Println(teacher.Name)
	}
	return teachers, nil
}
