package polling

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/Ademun/mining-lab-bot/pkg/event"
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
	eb   *event.Bus
	ids  []int
	opts ServiceOptions
	mu   *sync.RWMutex
}

func New(eb *event.Bus, opts *ServiceOptions) PollingService {
	if opts == nil {
		opts = &defaultServiceOptions
	}

	return &pollingService{
		eb:   eb,
		ids:  make([]int, 0),
		opts: *opts,
		mu:   &sync.RWMutex{},
	}
}

func (s *pollingService) Start(ctx context.Context) error {
	slog.Info("[Polling service] Starting...")
	if err := s.initIDUpdates(ctx); err != nil {
		return err
	}

	go s.startPollingLoop(ctx)

	slog.Info("[Polling service] Started")
	return nil
}

func (s *pollingService) GetPollingMode() PollingMode {
	return s.opts.Mode
}

func (s *pollingService) SetPollingMode(mode PollingMode) {
	s.opts.Mode = mode
}

func (s *pollingService) startPollingLoop(ctx context.Context) {
	if err := s.poll(ctx); err != nil {
		slog.Error(fmt.Sprintf("[Polling service] Failed to poll slots: %v", err))
	}
	go func() {
		for {
			var pollRate time.Duration
			switch s.opts.Mode {
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
					slog.Error(fmt.Sprintf("[Polling service] Failed to poll slots: %v", err))
				}
			}
		}
	}()
}

func (s *pollingService) poll(ctx context.Context) error {
	slog.Info("[Polling service] Polling...")
	var fetchRate time.Duration
	switch s.opts.Mode {
	case ModeNormal:
		fetchRate = time.Second * 1
	case ModeAggressive:
		fetchRate = time.Millisecond * 500
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	slots, err := PollAvailableSlots(ctx, s.ids, fetchRate)

	for _, slot := range slots {
		slotEvent := event.NewSlotEvent{Slot: slot}
		event.Publish(s.eb, slotEvent)
	}

	slog.Info("[Polling service] Polling finished")
	return err
}

func (s *pollingService) initIDUpdates(ctx context.Context) error {
	if err := s.updateIDs(ctx); err != nil {
		slog.Error(fmt.Sprintf("[PollingService] Init ID Updates failed: %v]", err))
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
					slog.Error(fmt.Sprintf("[PollingService] Failed to update ID list: %v]", err))
				}
			}
		}
	}()

	return nil
}

func (s *pollingService) updateIDs(ctx context.Context) error {
	slog.Info("[PollingService] Updating IDs...")
	ids, err := FetchServiceIDs(ctx)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.ids = ids
	s.mu.Unlock()
	slog.Info("[PollingService] Finished updating IDs")
	return nil
}
