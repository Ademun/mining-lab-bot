package polling

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/Ademun/mining-lab-bot/pkg/config"
	"golang.org/x/time/rate"
)

func (s *pollingService) pollSlots(ctx context.Context) (chan serverData, chan error) {
	results := make(chan serverData)

	var fetchRate time.Duration
	switch s.options.Mode {
	case config.ModeNormal:
		fetchRate = time.Second * 1
	case config.ModeAggressive:
		fetchRate = time.Millisecond * 250
	}

	limiter := rate.NewLimiter(rate.Every(fetchRate), 1)
	errChan := make(chan error, len(s.serviceIDs))
	var wg sync.WaitGroup

	go func() {
		for _, serviceID := range s.serviceIDs {
			wg.Add(1)
			go func() {
				defer wg.Done()
				select {
				case <-ctx.Done():
					return
				default:
					limiter.Wait(ctx)
					data, err := fetchServerData(ctx, serviceID)
					if err != nil {
						errChan <- err
						return
					}
					results <- *data
				}
			}()
		}
		wg.Wait()
		close(errChan)
		close(results)
	}()

	return results, errChan
}

func fetchServerData(ctx context.Context, serviceID int) (*serverData, error) {
	url := fmt.Sprintf("https://dikidi.net/ru/mobile/ajax/newrecord/get_datetimes/?company_id=550001&service_id%%5B%%5D=%d&with_first=1&day_month=", serviceID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, &ErrFetch{err: err, msg: "Failed to create request"}
	}

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, &ErrFetch{err: err, msg: "Failed to fetch document"}
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, &ErrFetch{err: errors.New("bad status code"), msg: fmt.Sprintf("Expected 200 but got %d", res.StatusCode)}
	}

	return unmarshalServiceData(res.Body)
}

func unmarshalServiceData(resData io.ReadCloser) (*serverData, error) {
	var data serverData
	if err := json.NewDecoder(resData).Decode(&data); err != nil {
		return nil, &ErrParseData{err: err}
	}
	return &data, nil
}
