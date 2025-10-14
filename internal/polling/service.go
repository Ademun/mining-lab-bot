package polling

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/Ademun/mining-lab-bot/pkg/config"
	"github.com/Ademun/mining-lab-bot/pkg/event"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/Ademun/mining-lab-bot/pkg/metrics"
)

type PollingService interface {
	Start(ctx context.Context) error
	GetPollingMode() config.PollingMode
	SetPollingMode(mode config.PollingMode)
}

type pollingService struct {
	eventBus   *event.Bus
	serviceIDs []int
	options    config.PollingConfig
	mutex      *sync.RWMutex
}

func New(eb *event.Bus, opts *config.PollingConfig) PollingService {
	return &pollingService{
		eventBus:   eb,
		serviceIDs: make([]int, 0),
		options:    *opts,
		mutex:      &sync.RWMutex{},
	}
}

func (s *pollingService) Start(ctx context.Context) error {
	slog.Info("Starting", "options", s.options, "service", logger.ServicePolling)
	if err := s.initIDUpdates(ctx); err != nil {
		return err
	}

	go s.startPollingLoop(ctx)

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
	if errs := s.poll(ctx); len(errs) > 0 {
		for _, err := range errs {
			slog.Warn("Polling error", "error", err, "service", logger.ServicePolling)
		}
	}
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
				if errs := s.poll(ctx); len(errs) > 0 {
					for _, err := range errs {
						slog.Warn("Polling error", "error", err, "service", logger.ServicePolling)
					}
				}
			}
		}
	}()
}

func (s *pollingService) poll(ctx context.Context) []error {
	slog.Info("Polling", "service", logger.ServicePolling)
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
		}
		if errors.Is(err, fetchErr) {
			fetchErrs++
		}
	}

	metrics.Global().RecordPollResults(len(slots), parseErrs, fetchErrs, s.GetPollingMode(), total)

	for _, slot := range slots {
		slotEvent := event.NewSlotEvent{Slot: slot}
		event.Publish(ctx, s.eventBus, slotEvent)
	}
	slog.Info("Polling finished", "service", logger.ServicePolling)
	return errs
}

func (s *pollingService) initIDUpdates(ctx context.Context) error {
	if err := s.updateIDs(ctx); err != nil {
		slog.Error("Failed to update IDs", "error", err, "service", logger.ServicePolling)
		return err
	}

	ticker := time.Tick(time.Hour * 24)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker:
				err := s.updateIDs(ctx)
				if err != nil {
					slog.Error("Failed to update IDs", "error", err, "service", logger.ServicePolling)
				}
			}
		}
	}()

	return nil
}

func (s *pollingService) updateIDs(ctx context.Context) error {
	slog.Info("Updating IDs", "service", logger.ServicePolling)
	ids, err := FetchServiceIDs(ctx, s.options.ServiceURL)
	if err != nil {
		return err
	}
	s.mutex.Lock()
	s.serviceIDs = ids
	s.mutex.Unlock()
	slog.Info("Finished updating IDs", "service", logger.ServicePolling)
	return nil
}
