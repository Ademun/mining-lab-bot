package subscription

import (
	"time"

	"github.com/Ademun/mining-lab-bot/internal/polling"
	"github.com/google/uuid"
)

type TimeRange struct {
	TimeStart string
	TimeEnd   string
}

// ResponseUser represents a unique user id, with a map of preferred subscription times based on search by slot info
// It is needed to collect all distinct preferred times per weekday for a specific subscription group, that belongs to a user
// This way, we can group n subscription that target one slot, but a different times, and send only 1 notification instead of n
type ResponseUser struct {
	UserID         int
	PreferredTimes map[time.Weekday][]TimeRange
}

type ResponseSubscription struct {
	UUID           uuid.UUID
	UserID         int
	LabType        polling.LabType
	LabNumber      int
	LabAuditorium  *int
	LabDomain      *polling.Domain
	Weekday        *int
	PreferredTimes []TimeRange
}

type DBSubscription struct {
	UUID          uuid.UUID       `db:"uuid"`
	UserID        int             `db:"user_id"`
	LabType       polling.LabType `db:"lab_type"`
	LabNumber     int             `db:"lab_number"`
	LabAuditorium *int            `db:"lab_auditorium"`
	LabDomain     *polling.Domain `db:"lab_domain"`
	Weekday       *int            `db:"weekday"`
}

type DBSubscriptionTimes struct {
	SubscriptionUUID uuid.UUID `db:"subscription_uuid"`
	TimeStart        string    `db:"time_start"`
	TimeEnd          string    `db:"time_end"`
}

func toResponse(sub DBSubscription, subTimes []DBSubscriptionTimes) ResponseSubscription {
	prefTimes := make([]TimeRange, len(subTimes))
	for idx, t := range subTimes {
		prefTimes[idx] = TimeRange{
			TimeStart: t.TimeStart,
			TimeEnd:   t.TimeEnd,
		}
	}
	return ResponseSubscription{
		UUID:           sub.UUID,
		UserID:         sub.UserID,
		LabType:        sub.LabType,
		LabNumber:      sub.LabNumber,
		LabAuditorium:  sub.LabAuditorium,
		LabDomain:      sub.LabDomain,
		Weekday:        sub.Weekday,
		PreferredTimes: prefTimes,
	}
}

type RequestSubscription struct {
	UserID        int
	Type          polling.LabType
	LabNumber     int
	LabAuditorium *int
	LabDomain     *polling.Domain
	Weekday       *int
	Lessons       []int
}

func (rs RequestSubscription) toDBModels() (DBSubscription, []DBSubscriptionTimes) {
	dbSub := DBSubscription{
		UUID:          uuid.New(),
		UserID:        rs.UserID,
		LabType:       rs.Type,
		LabNumber:     rs.LabNumber,
		LabAuditorium: rs.LabAuditorium,
		LabDomain:     rs.LabDomain,
		Weekday:       rs.Weekday,
	}
	dbTimes := make([]DBSubscriptionTimes, len(rs.Lessons))
	for idx, timeRange := range LessonsToTimeRanges(rs.Lessons...) {
		dbTimes[idx] = DBSubscriptionTimes{
			SubscriptionUUID: dbSub.UUID,
			TimeStart:        timeRange.TimeStart,
			TimeEnd:          timeRange.TimeEnd,
		}
	}
	return dbSub, dbTimes
}
