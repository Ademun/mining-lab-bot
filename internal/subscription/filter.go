package subscription

import (
	"github.com/Ademun/mining-lab-bot/internal/polling"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type SubFilters struct {
	UserID        int
	Type          *polling.LabType
	LabNumber     int
	LabAuditorium int
	LabDomain     *polling.LabDomain
	Weekdays      []int
}

func (f *SubFilters) buildQuery() (string, []interface{}, error) {
	q := squirrel.Select("*").From("subscriptions")
	conditions := squirrel.And{}

	if f.UserID != 0 {
		conditions = append(conditions, squirrel.Eq{"user_id": f.UserID})
	}
	if f.Type != nil {
		conditions = append(conditions, squirrel.Eq{"lab_type": f.Type})
	}
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
	q := squirrel.Select("*").From("subscription_times")
	conditions := squirrel.And{}

	if len(f.SubUUIDs) > 0 {
		conditions = append(conditions, squirrel.Eq{"subscription_uuid": f.SubUUIDs})
	}

	if len(f.Includes) > 0 {
		orConditions := squirrel.Or{}
		for _, targetTime := range f.Includes {
			orConditions = append(orConditions, squirrel.And{
				squirrel.LtOrEq{"time_start": targetTime},
				squirrel.Gt{"time_end": targetTime},
			})
		}
		conditions = append(conditions, orConditions)
	}

	if len(conditions) > 0 {
		q = q.Where(conditions)
	}

	return q.ToSql()
}
