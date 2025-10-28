package subscription

import (
	"strings"

	"github.com/Ademun/mining-lab-bot/internal/polling"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type SubFilters struct {
	UserID        int
	Type          polling.LabType
	LabNumber     int
	LabAuditorium int
	LabDomain     *polling.Domain
	Weekdays      []int
}

func (f *SubFilters) buildQuery() (string, []interface{}, error) {
	q := squirrel.Select("*").From("subscriptions")
	conditions := squirrel.And{}

	if f.UserID != 0 {
		conditions = append(conditions, squirrel.Eq{"user_id": f.UserID})
	}
	conditions = append(conditions, squirrel.Eq{"type": f.Type})
	if f.LabNumber != 0 {
		conditions = append(conditions, squirrel.Eq{"lab_number": f.LabNumber})
	}
	if f.LabAuditorium != 0 {
		conditions = append(conditions, squirrel.Eq{"lab_auditorium": f.LabAuditorium})
	}
	if f.LabDomain != nil {
		conditions = append(conditions, squirrel.Eq{"lab_domain": f.LabDomain})
	}
	if len(f.Weekdays) > 0 {
		conditions = append(conditions, squirrel.Eq{"weekday": f.Weekdays})
	}
	if len(conditions) > 0 {
		q = q.Where(conditions)
	}

	return q.ToSql()
}

type TimeFilters struct {
	SubUUIDs []uuid.UUID
	Includes []string
}

func (f *TimeFilters) buildQuery() (string, []interface{}, error) {
	q := squirrel.Select("*").From("subscriptions_times")
	conditions := squirrel.And{}

	if len(f.SubUUIDs) > 0 {
		conditions = append(conditions, squirrel.Eq{"subscription_uuid": f.SubUUIDs})
	}
	if len(f.Includes) > 0 {
		placeholders := make([]string, len(f.Includes))
		args := make([]interface{}, len(f.Includes))

		for idx, t := range f.Includes {
			placeholders[idx] = "?"
			args[idx] = t
		}

		conditions = append(conditions, squirrel.Expr("EXISTS (SELECT 1 FROM (SELECT "+strings.Join(placeholders, " UNION ALL SELECT ")+") AS t(target_time) WHERE target_time >= time_start AND target_time < time_end)"))
	}
	if len(conditions) > 0 {
		q = q.Where(conditions)
	}

	return q.ToSql()
}
