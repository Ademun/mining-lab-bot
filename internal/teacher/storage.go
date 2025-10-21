package teacher

import (
	"context"
	"database/sql"
	"time"

	"github.com/Ademun/mining-lab-bot/pkg/errs"
	"github.com/Ademun/mining-lab-bot/pkg/model"
)

type Repo interface {
	FindByWeekNumberWeekdayAuditorium(ctx context.Context, weekNumber int, weekday time.Weekday, auditorium int) ([]model.Teacher, error)
}

type teacherRepo struct {
	db *sql.DB
}

func NewRepo(ctx context.Context, db *sql.DB) (Repo, error) {
	query := `
create table if not exists teachers (
    id integer primary key,
    name text not null,
    auditorium integer not null,
    week_number integer not null,
    weekday integer not null,
    time_start text not null,
    time_end text not null
)
`
	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return nil, &errs.ErrQueryExecution{Operation: "Init", Query: query, Err: err}
	}

	return &teacherRepo{db: db}, nil
}

func (t *teacherRepo) FindByWeekNumberWeekdayAuditorium(ctx context.Context, weekNumber int, weekday time.Weekday, auditorium int) ([]model.Teacher, error) {
	query := `select id, name, auditorium, week_number, weekday, time_start, time_end from teachers where week_number = ? and weekday = ? and auditorium = ?`

	rows, err := t.db.QueryContext(ctx, query, weekNumber, weekday, auditorium)
	if err != nil {
		return nil, &errs.ErrQueryExecution{Operation: "FindByWeekNumberWeekday", Query: query, Err: err}
	}
	defer rows.Close()

	var teachers []model.Teacher
	for rows.Next() {
		var teacher model.Teacher
		err := rows.Scan(&teacher.ID, &teacher.Name, &teacher.Auditorium, &teacher.WeekNumber, &teacher.Weekday, &teacher.TimeStart, &teacher.TimeEnd)
		if err != nil {
			return nil, &errs.ErrRowIteration{Operation: "FindByWeekNumberWeekday", Query: query, Err: err}
		}
		teachers = append(teachers, teacher)
	}

	if err := rows.Err(); err != nil {
		return nil, &errs.ErrRowIteration{Operation: "FindByWeekNumberWeekday", Query: query, Err: err}
	}

	return teachers, nil
}
