package polling

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/Ademun/mining-lab-bot/pkg/event"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
)

type PollingService interface {
	Start(ctx context.Context) error
	GetPollingMode() PollingMode
	SetPollingMode(PollingMode)
}

type ServiceOptions struct {
	Mode PollingMode
}

type PollingMode int

const (
	ModeNormal PollingMode = iota
	ModeAggressive
)

var defaultServiceOptions = ServiceOptions{Mode: ModeNormal}

type pollingService struct {
	eventBus   *event.Bus
	serviceURL string
	serviceIDs []int
	options    ServiceOptions
	mutex      *sync.RWMutex
}

func New(eb *event.Bus, serviceURL string, opts *ServiceOptions) PollingService {
	if opts == nil {
		opts = &defaultServiceOptions
	}

	return &pollingService{
		eventBus:   eb,
		serviceURL: serviceURL,
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

func (s *pollingService) GetPollingMode() PollingMode {
	return s.options.Mode
}

func (s *pollingService) SetPollingMode(mode PollingMode) {
	s.options.Mode = mode
}

func (s *pollingService) startPollingLoop(ctx context.Context) {
	if err := s.poll(ctx); err != nil {
		slog.Error("Polling errors", "errors", err, "service", logger.ServicePolling)
	}
	go func() {
		for {
			var pollRate time.Duration
			switch s.options.Mode {
			case ModeNormal:
				pollRate = time.Minute * 1
			case ModeAggressive:
				pollRate = time.Second * 25
			}
			ticker := time.Tick(pollRate)
			select {
			case <-ctx.Done():
				return
			case <-ticker:
				if err := s.poll(ctx); err != nil {
					slog.Error("Polling errors", "errors", err, "service", logger.ServicePolling)
				}
			}
		}
	}()
}

func (s *pollingService) poll(ctx context.Context) error {
	slog.Info("Polling", "service", logger.ServicePolling)
	var fetchRate time.Duration
	switch s.options.Mode {
	case ModeNormal:
		fetchRate = time.Second * 2
	case ModeAggressive:
		fetchRate = time.Millisecond * 500
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	slots, err := PollAvailableSlots(ctx, s.serviceIDs, fetchRate)

	for _, slot := range slots {
		slotEvent := event.NewSlotEvent{Slot: slot}
		event.Publish(s.eventBus, ctx, slotEvent)
	}

	slog.Info("Polling finished", "service", logger.ServicePolling)
	return err
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
	ids, err := FetchServiceIDs(ctx, s.serviceURL)
	if err != nil {
		return err
	}
	s.mutex.Lock()
	s.serviceIDs = ids
	s.mutex.Unlock()
	slog.Info("Finished updating IDs", "service", logger.ServicePolling)
	return nil
}
