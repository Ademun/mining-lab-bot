package event

import "github.com/Ademun/mining-lab-bot/pkg/model"

type Event struct{}

func (e Event) Type() string {
	return "event"
}

type NewSlotEvent struct {
	Event
	Slot model.Slot
}

func (e NewSlotEvent) Type() string {
	return "event:slot:new"
}

type NewNotificationEvent struct {
	Event
	Notification model.Notification
}

func (e NewNotificationEvent) Type() string {
	return "event:notification:new"
}

type NewSubscriptionEvent struct {
	Event
	Subscription model.Subscription
}

func (e NewSubscriptionEvent) Type() string {
	return "event:subscription:new"
}
