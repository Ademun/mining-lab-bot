package polling

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/Ademun/mining-lab-bot/internal/notification"
	"github.com/Ademun/mining-lab-bot/internal/teacher"
	"github.com/Ademun/mining-lab-bot/pkg/config"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/Ademun/mining-lab-bot/pkg/metrics"
)

type Service interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context)
}

type pollingService struct {
	notifService   notification.Service
	teacherService teacher.Service
	options        config.PollingConfig
	serviceIDs     []int
	httpClient     *http.Client
	wg             sync.WaitGroup
	mu             sync.RWMutex
}

func New(notifService notification.Service, teacherService teacher.Service, opts *config.PollingConfig) Service {
	return &pollingService{
		notifService:   notifService,
		teacherService: teacherService,
		options:        *opts,
		serviceIDs:     make([]int, 0),
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
		wg: sync.WaitGroup{},
		mu: sync.RWMutex{},
	}
}

func (s *pollingService) Start(ctx context.Context) error {
	slog.Info("Starting", "options", s.options, "service", logger.ServicePolling)

	s.startIDUpdateLoop(ctx)
	s.startPollingLoop(ctx)

	slog.Info("Started", "service", logger.ServicePolling)
	return nil
}

func (s *pollingService) Stop(ctx context.Context) {
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		slog.Info("Stopped", "service", logger.ServicePolling)
	case <-ctx.Done():
		slog.Error("Stopped by timeout", "service", logger.ServicePolling)
	}
}

func (s *pollingService) startPollingLoop(ctx context.Context) {
	s.poll(ctx)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		pollRate := s.getPolRate()
		ticker := time.NewTicker(pollRate)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.poll(ctx)
				newRate := s.getPolRate()
				ticker.Reset(newRate)
			}
		}
	}()
}

func (s *pollingService) getPolRate() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	switch s.options.Mode {
	case config.ModeNormal:
		return s.options.NormalPollRate
	case config.ModeAggressive:
		return s.options.AggressivePollRate
	}
	slog.Warn("Unknown service mode", "mode", s.options.Mode, "service", logger.ServicePolling)
	return s.options.NormalPollRate
}

func (s *pollingService) getFetchRate() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	switch s.options.Mode {
	case config.ModeNormal:
		return s.options.NormalFetchRate
	case config.ModeAggressive:
		return s.options.AggressiveFetchRate
	}
	slog.Warn("Unknown service mode", "mode", s.options.Mode, "service", logger.ServicePolling)
	return s.options.NormalFetchRate
}

func (s *pollingService) poll(ctx context.Context) {
	s.wg.Add(1)
	defer s.wg.Done()

	dataChan, errChan := s.pollServerData(ctx)
	var dataLen, fetchErrs, parseErrs int
	start := time.Now()
	for dataChan != nil || errChan != nil {
		select {
		case <-ctx.Done():
			return
		case data, ok := <-dataChan:
			if !ok {
				dataChan = nil
				continue
			}
			serviceID := data.Data.ServiceID
			slots, err := s.ParseServerData(ctx, &data, serviceID)
			if err != nil {
				parseErrs++
				slog.Warn("Parsing error", "error", err, "service", logger.ServicePolling)
			}
			dataLen += len(slots)
			for _, slot := range slots {
				s.notifService.SendNotification(ctx, slot)
			}
		case err, ok := <-errChan:
			if !ok {
				fetchErrs++
				errChan = nil
				continue
			}
			slog.Warn("Polling error", "error", err, "service", logger.ServicePolling)
		}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	metrics.Global().RecordPollResults(dataLen, len(s.serviceIDs), parseErrs, fetchErrs, s.options.Mode, time.Since(start))
}

func (s *pollingService) startIDUpdateLoop(ctx context.Context) {
	s.updateIDs(ctx)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(s.options.ServiceIDUpdateRate)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.updateIDs(ctx)
			}
		}
	}()
}

func (s *pollingService) updateIDs(ctx context.Context) {
	s.wg.Add(1)
	defer s.wg.Done()

	ids, err := s.fetchServiceIDs(ctx)
	if err != nil {
		slog.Error("Failed to fetch service IDs", "error", err, "service", logger.ServicePolling)
		return
	}

	s.mu.Lock()
	s.serviceIDs = ids
	s.mu.Unlock()
}
