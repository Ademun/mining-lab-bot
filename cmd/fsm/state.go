package fsm

import (
	"github.com/Ademun/mining-lab-bot/internal/polling"
	"github.com/Ademun/mining-lab-bot/internal/subscription"
)

type ConversationStep string

const (
	StepIdle                            ConversationStep = "idle"
	StepAwaitingLabType                 ConversationStep = "awaiting_lab_type"
	StepAwaitingLabNumber               ConversationStep = "awaiting_lab_number"
	StepAwaitingLabAuditorium           ConversationStep = "awaiting_lab_auditorium"
	StepAwaitingLabDomain               ConversationStep = "awaiting_lab_domain"
	StepAwaitingLabWeekday              ConversationStep = "awaiting_lab_weekday"
	StepAwaitingLabLessons              ConversationStep = "awaiting_lab_lessons"
	StepAwaitingSubCreationConfirmation ConversationStep = "awaiting_sub_creation_confirmation"
	StepAwaitingListingSubsAction       ConversationStep = "awaiting_listing_action"
	StepAwaitingFeedbackMsg             ConversationStep = "awaiting_feedback_msg"
	StepAwaitingFeedbackReaction        ConversationStep = "awaiting_feedback_reaction"
	StepAwaitingTeacherAuditorium       ConversationStep = "awaiting_teacher_auditorium"
	StepAwaitingTeacherWeekParity       ConversationStep = "awaiting_teacher_week_parity"
	StepAwaitingTeacherWeekday          ConversationStep = "awaiting_teacher_weekday"
	StepAwaitingTeacherLesson           ConversationStep = "awaiting_teacher_lesson"
	StepAwaitingTeacherSurname          ConversationStep = "awaiting_teacher_surname"
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

type TeacherReportFlowData struct {
	UserID     int64
	Auditorium int
	WeekParity string // "even" или "odd"
	Weekday    int
	LessonNum  int
	Surname    string
}

func (d *TeacherReportFlowData) StateData() {}

func dataTypeForStep(step ConversationStep) StateData {
	switch step {
	case StepIdle, StepAwaitingFeedbackMsg, StepAwaitingFeedbackReaction:
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
	case StepAwaitingTeacherAuditorium,
		StepAwaitingTeacherWeekParity,
		StepAwaitingTeacherWeekday,
		StepAwaitingTeacherLesson,
		StepAwaitingTeacherSurname:
		return &TeacherReportFlowData{}
	}
	return nil
}
