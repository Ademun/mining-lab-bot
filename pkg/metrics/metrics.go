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
	TotalPolls      int
	Mode            config.PollingMode
	ParsingErrors   int
	FetchErrors     int
	LastPollingTime time.Duration
	LastSlotNumber  int
	LastIDNumber    int
}

type NotificationMetrics struct {
	TotalNotifications int
	CacheLength        int
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

func (m *Metrics) RecordPollResults(slotsLen int, idLen int, parseErrs int, fetchErrs int, mode config.PollingMode, pollingDuration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.PollingMetrics.TotalPolls++
	m.PollingMetrics.Mode = mode
	m.PollingMetrics.ParsingErrors += parseErrs
	m.PollingMetrics.FetchErrors += fetchErrs
	m.PollingMetrics.LastPollingTime = pollingDuration
	m.PollingMetrics.LastSlotNumber = slotsLen
	m.PollingMetrics.LastIDNumber = idLen
}

func (m *Metrics) RecordNotificationResults(notifLen int, cacheLen int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.NotificationMetrics.TotalNotifications += notifLen
	m.NotificationMetrics.CacheLength = cacheLen
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
