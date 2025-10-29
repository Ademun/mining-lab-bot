package notification

import (
	"context"
	"time"

	"github.com/Ademun/mining-lab-bot/internal/polling"
	"github.com/Ademun/mining-lab-bot/internal/subscription"
)

type PreferredTimes map[time.Weekday][]subscription.TimeRange

type Notification struct {
	UserID         int
	PreferredTimes PreferredTimes
	Slot           polling.Slot
}

type SlotNotifier interface {
	SendNotification(ctx context.Context, notif Notification)
}
