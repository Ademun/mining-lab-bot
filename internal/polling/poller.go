package polling

import (
	"context"
	"sync"
	"time"

	"github.com/Ademun/mining-lab-bot/pkg/model"
	"golang.org/x/time/rate"
)

func (s *pollingService) PollAvailableSlots(ctx context.Context, ids []int, fetchRate time.Duration) ([]model.Slot, []error) {
	slots := make([]model.Slot, 0)
	results := make(chan model.Slot)

	var wg sync.WaitGroup
	limiter := rate.NewLimiter(rate.Every(fetchRate), 1)
	errChan := make(chan error, len(ids))
	go func() {
		for _, serviceID := range ids {
			wg.Add(1)
			go func() {
				defer wg.Done()
				select {
				case <-ctx.Done():
					return
				default:
					limiter.Wait(ctx)
					data, err := FetchServiceData(ctx, serviceID)
					if err != nil {
						errChan <- err
						return
					}
					res, err := s.ParseServiceData(ctx, data, serviceID)
					if err != nil {
						errChan <- err
						return
					}
					for _, slot := range res {
						results <- slot
					}
				}
			}()
		}
		wg.Wait()
		close(errChan)
		close(results)
	}()

	for slot := range results {
		slots = append(slots, slot)
	}

	errList := make([]error, 0)
	for err := range errChan {
		errList = append(errList, err)
	}

	return slots, errList
}
