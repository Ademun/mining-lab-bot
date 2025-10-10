package polling

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/Ademun/mining-lab-bot/pkg/model"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

func PollAvailableSlots(ctx context.Context, ids []int, fetchRate time.Duration) ([]model.Slot, error) {
	start := time.Now()

	slots := make([]model.Slot, 0)
	results := make(chan model.Slot)
	IDsChecked := atomic.Int64{}

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
					start := time.Now()
					IDsChecked.Add(1)
					slog.Info(fmt.Sprintf("Processing service. id = %d [%d/%d]", serviceID, IDsChecked.Load(), len(ids)))
					data, err := FetchServiceData(serviceID)
					if err != nil {
						return err
					}
					res, err := ParseServiceData(data)
					if err != nil {
						fmt.Println("error parsing id", serviceID)
						return err
					}
					for _, slot := range res {
						results <- slot
					}
					slog.Info(fmt.Sprintf("Finished processing service. id %d. Time elapsed %s\n", serviceID, time.Since(start)))
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

	slog.Info(fmt.Sprintf("Finished checking labs in %s. Total available slots %d", time.Since(start), len(slots)))
	return slots, processingErr
}
