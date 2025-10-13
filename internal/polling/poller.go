package polling

import (
	"context"
	"time"

	"github.com/Ademun/mining-lab-bot/pkg/model"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

func PollAvailableSlots(ctx context.Context, ids []int, fetchRate time.Duration) ([]model.Slot, error) {
	slots := make([]model.Slot, 0)
	results := make(chan model.Slot)

	var processingErr error
	var eg errgroup.Group
	limiter := rate.NewLimiter(rate.Every(fetchRate), 1)
	go func() {
		for _, serviceID := range ids {
			eg.Go(func() error {
				select {
				case <-ctx.Done():
					return nil
				default:
					limiter.Wait(context.Background())
					data, err := FetchServiceData(ctx, serviceID)
					if err != nil {
						return err
					}
					res, err := ParseServiceData(data)
					if err != nil {
						return err
					}
					for _, slot := range res {
						results <- slot
					}
					return nil
				}
			})
		}
		if err := eg.Wait(); err != nil {
			processingErr = err
		}
		close(results)
	}()

	for slot := range results {
		slots = append(slots, slot)
	}

	return slots, processingErr
}
