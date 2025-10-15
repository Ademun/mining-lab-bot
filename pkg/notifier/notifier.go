package notifier

import (
	"context"

	"github.com/Ademun/mining-lab-bot/pkg/model"
)

type SlotNotifier interface {
	SendNotification(ctx context.Context, notif model.Notification)
}
