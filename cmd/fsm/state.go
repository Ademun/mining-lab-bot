package fsm

import (
	"github.com/Ademun/mining-lab-bot/internal/polling"
	"github.com/Ademun/mining-lab-bot/internal/subscription"
)

type ConversationStep string

const (
	StepIdle ConversationStep = "idle"
	// /sub chain
	StepAwaitingLabType                 ConversationStep = "awaiting_lab_type"
	StepAwaitingLabNumber               ConversationStep = "awaiting_lab_number"
	StepAwaitingLabAuditorium           ConversationStep = "awaiting_lab_auditorium"
	StepAwaitingLabDomain               ConversationStep = "awaiting_lab_domain"
	StepAwaitingLabWeekday              ConversationStep = "awaiting_lab_weekday"
	StepAwaitingLabLessons              ConversationStep = "awaiting_lab_lessons"
	StepAwaitingSubCreationConfirmation ConversationStep = "awaiting_sub_creation_confirmation"
	// /unsub chain
	StepAwaitingListingSubsAction ConversationStep = "awaiting_listing_action"
)

type StateData interface {
	StateData()
}

type IdleData struct{}

func (data *IdleData) StateData() {}

type SubscriptionCreationFlowData struct {
	UserID        int
	LabType       polling.LabType
	LabNumber     int
	LabAuditorium *int
	LabDomain     *polling.LabDomain
	Weekday       *int
	Lessons       []int
}

func (d *SubscriptionCreationFlowData) StateData() {}

type SubscriptionListingFlowData struct {
	UserSubs []subscription.ResponseSubscription
}

func (d *SubscriptionListingFlowData) StateData() {}

func dataTypeForStep(step ConversationStep) StateData {
	switch step {
	case StepIdle:
		return &IdleData{}
	case StepAwaitingLabType,
		StepAwaitingLabNumber,
		StepAwaitingLabAuditorium,
		StepAwaitingLabDomain,
		StepAwaitingLabWeekday,
		StepAwaitingLabLessons,
		StepAwaitingSubCreationConfirmation:
		return &SubscriptionCreationFlowData{}
	case StepAwaitingListingSubsAction:
		return &SubscriptionListingFlowData{}
	}
	return nil
}
