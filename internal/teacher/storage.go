package teacher

import (
	"context"
	"database/sql"
	"time"

	"github.com/Ademun/mining-lab-bot/pkg/model"
)

type Repo interface {
	FindByWeekNumberWeekdayTimeStart(ctx context.Context, weekNumber int, weekday time.Weekday, startTime string) (model.Teacher, error)
}

type teacherRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) (Repo, error) {
	query := `
create table if not exists teachers (
    uuid text not null primary key,
    name text not null,
    weekday integer not null,
    
)
`
}
