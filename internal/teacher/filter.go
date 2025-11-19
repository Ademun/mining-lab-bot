package teacher

import (
	"time"

	"github.com/Masterminds/squirrel"
)

type Filter struct {
	WeekNumber int
	Weekday    time.Weekday
	Auditorium int
	TargetTime time.Time
}

func (f *Filter) buildQuery() (string, []interface{}, error) {
	q := squirrel.Select("*").From("teachers")
	conditions := squirrel.And{}

	if f.WeekNumber > 0 {
		conditions = append(conditions, squirrel.Eq{"week_number": f.WeekNumber})
	}
	if f.Weekday >= 0 {
		conditions = append(conditions, squirrel.Eq{"weekday": f.Weekday})
	}
	if f.Auditorium > 0 {
		conditions = append(conditions, squirrel.Eq{"auditorium": f.Auditorium})
	}

	targetTime := f.TargetTime.Format("15:04")
	conditions = append(conditions, squirrel.LtOrEq{"time_start": targetTime})
	conditions = append(conditions, squirrel.Gt{"time_end": targetTime})

	if len(conditions) > 0 {
		q = q.Where(conditions)
	}

	return q.ToSql()
}
