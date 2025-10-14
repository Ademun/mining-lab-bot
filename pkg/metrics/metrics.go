package metrics

import (
	"sync"
	"time"

	"github.com/Ademun/mining-lab-bot/pkg/config"
)

type Metrics struct {
	PollingMetrics      PollingMetrics
	NotificationMetrics NotificationMetrics
	SubscriptionMetrics SubscriptionMetrics
	StartTime           time.Time
	mu                  *sync.Mutex
}

type PollingMetrics struct {
	TotalPolls         int
	Mode               config.PollingMode
	ParsingErrors      int
	FetchErrors        int
	AveragePollingTime time.Duration
	AverageSlotNumber  int
}

type NotificationMetrics struct {
	TotalNotifications   int
	CacheLength          int
	AverageNotifications int
}

type SubscriptionMetrics struct {
	TotalSubscriptions int
}

var global = &Metrics{
	StartTime: time.Now(),
	mu:        &sync.Mutex{},
}

func Global() *Metrics {
	return global
}

func (m *Metrics) RecordPollResults(slotsLen int, parseErrs int, fetchErrs int, mode config.PollingMode, pollingDuration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.PollingMetrics.TotalPolls++
	m.PollingMetrics.Mode = mode

	m.PollingMetrics.ParsingErrors += parseErrs
	m.PollingMetrics.FetchErrors += fetchErrs

	prevAvg := int(m.PollingMetrics.AveragePollingTime.Microseconds())
	prevCount := m.PollingMetrics.TotalPolls - 1
	newDuration := int(pollingDuration.Microseconds())
	m.PollingMetrics.AveragePollingTime = time.Duration((prevAvg*prevCount + newDuration) / (prevCount + 1))

	prevAvg = m.PollingMetrics.AverageSlotNumber
	newSlots := slotsLen
	m.PollingMetrics.AverageSlotNumber = (prevAvg*prevCount + newSlots) / (prevCount + 1)
}

func (m *Metrics) RecordNotificationResults(notifLen int, cacheLen int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.NotificationMetrics.TotalNotifications += notifLen
	m.NotificationMetrics.CacheLength = cacheLen

	prevAvg := m.NotificationMetrics.AverageNotifications
	prevCount := m.NotificationMetrics.TotalNotifications
	m.NotificationMetrics.AverageNotifications = (prevAvg*prevCount + notifLen) / (prevCount + 1)
}

func (m *Metrics) RecordSubscriptionResults(subsLen int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.SubscriptionMetrics.TotalSubscriptions += subsLen
}

func (m *Metrics) Snapshot() Metrics {
	m.mu.Lock()
	defer m.mu.Unlock()
	return *m
}
