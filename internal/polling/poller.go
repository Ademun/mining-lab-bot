package polling

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

func (s *pollingService) pollServerData(ctx context.Context) (chan ServerData, chan error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make(chan ServerData)

	fetchRate := s.getFetchRate()
	limiter := rate.NewLimiter(rate.Every(fetchRate), 1)
	errChan := make(chan error)
	var wg sync.WaitGroup

	go func() {
		for _, serviceID := range s.serviceIDs {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := limiter.Wait(ctx); err != nil {
					return
				}
				data, err := s.fetchServerData(ctx, serviceID)
				if err != nil {
					errChan <- err
					return
				}
				results <- *data
			}()
		}
		wg.Wait()
		close(errChan)
		close(results)
	}()

	return results, errChan
}

func (s *pollingService) fetchServerData(ctx context.Context, serviceID int) (*ServerData, error) {
	url := fmt.Sprintf("https://dikidi.net/ru/mobile/ajax/newrecord/get_datetimes/?company_id=550001&service_id%%5B%%5D=%d&with_first=1&day_month=", serviceID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, &ErrFetch{url: url, err: err, msg: "Failed to create request"}
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		return nil, &ErrFetch{url: url, err: err, msg: "Failed to fetch document"}
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, &ErrFetch{url: url, err: errors.New("bad status code"), msg: fmt.Sprintf("Expected 200 but got %d", res.StatusCode)}
	}

	return unmarshalServiceData(res.Body, serviceID)
}

func unmarshalServiceData(serverData io.ReadCloser, serviceID int) (*ServerData, error) {
	var data ServerData
	if err := json.NewDecoder(serverData).Decode(&data); err != nil {
		return nil, &ErrParseData{msg: fmt.Sprintf("serviceID: %d", serviceID), err: err}
	}
	return &data, nil
}
