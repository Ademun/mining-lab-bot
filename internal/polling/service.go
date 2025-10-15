package polling

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/Ademun/mining-lab-bot/internal/notification"
	"github.com/Ademun/mining-lab-bot/pkg/config"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/Ademun/mining-lab-bot/pkg/metrics"
)

type Service interface {
	Start(ctx context.Context) error
	GetPollingMode() config.PollingMode
	SetPollingMode(mode config.PollingMode)
}

type pollingService struct {
	notifService notification.Service
	options      config.PollingConfig
	serviceIDs   []int
	mutex        *sync.RWMutex
}

func New(notifService notification.Service, opts *config.PollingConfig) Service {
	return &pollingService{
		notifService: notifService,
		options:      *opts,
		serviceIDs:   make([]int, 0),
		mutex:        &sync.RWMutex{},
	}
}

func (s *pollingService) Start(ctx context.Context) error {
	slog.Info("Starting", "options", s.options, "service", logger.ServicePolling)

	s.startIDUpdateLoop(ctx)
	s.startPollingLoop(ctx)

	slog.Info("Started", "service", logger.ServicePolling)
	return nil
}

func (s *pollingService) GetPollingMode() config.PollingMode {
	return s.options.Mode
}

func (s *pollingService) SetPollingMode(mode config.PollingMode) {
	s.options.Mode = mode
}

func (s *pollingService) startPollingLoop(ctx context.Context) {
	s.poll(ctx)

	go func() {
		for {
			var pollRate time.Duration
			switch s.options.Mode {
			case config.ModeNormal:
				pollRate = time.Minute * 1
			case config.ModeAggressive:
				pollRate = time.Second * 25
			}
			ticker := time.Tick(pollRate)
			select {
			case <-ctx.Done():
				return
			case <-ticker:
				s.poll(ctx)
			}
		}
	}()
}

func (s *pollingService) poll(ctx context.Context) {
	var fetchRate time.Duration
	switch s.options.Mode {
	case config.ModeNormal:
		fetchRate = time.Second * 2
	case config.ModeAggressive:
		fetchRate = time.Millisecond * 500
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	start := time.Now()
	slots, errs := PollAvailableSlots(ctx, s.serviceIDs, fetchRate)
	total := time.Since(start)

	parseErrs, fetchErrs := 0, 0
	var parseErr *ErrParseData
	var fetchErr *ErrFetch
	for _, err := range errs {
		if errors.Is(err, parseErr) {
			parseErrs++
			slog.Warn("Parsing error", "error", err, "service", logger.ServicePolling)
		}
		if errors.Is(err, fetchErr) {
			fetchErrs++
			slog.Warn("Fetching error", "error", err, "service", logger.ServicePolling)
		}
	}

	metrics.Global().RecordPollResults(len(slots), parseErrs, fetchErrs, s.GetPollingMode(), total)

	for _, slot := range slots {
		s.notifService.SendNotification(ctx, slot)
	}
}

func (s *pollingService) startIDUpdateLoop(ctx context.Context) {
	s.updateIDs(ctx)

	ticker := time.Tick(time.Hour * 24)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker:
				s.updateIDs(ctx)
			}
		}
	}()
}

func (s *pollingService) updateIDs(ctx context.Context) {
	ids, err := FetchServiceIDs(ctx, s.options.ServiceURL)
	if err != nil {
		slog.Warn("Failed to fetch service IDs", "error", err, "service", logger.ServicePolling)
	}

	s.mutex.Lock()
	s.serviceIDs = ids
	s.mutex.Unlock()

	for _, id := range ids {
		fmt.Println(id)
	}
}
