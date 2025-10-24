package polling

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer close(errChan)
		defer close(results)

		for _, serviceID := range s.serviceIDs {
			wg.Add(1)
			serviceID := serviceID
			go func() {
				defer wg.Done()
				if err := limiter.Wait(ctx); err != nil {
					select {
					case errChan <- err:
					case <-ctx.Done():
					}
					return
				}
				if ctx.Err() != nil {
					return
				}
				initialData, err := s.fetchServerData(ctx, serviceID, nil)
				if err != nil {
					select {
					case errChan <- err:
					case <-ctx.Done():
					}
					return
				}
				initialData.Data.ServiceID = serviceID
				dates := initialData.Data.DatesTrue
				if len(dates) == 0 {
					select {
					case results <- *initialData:
						return
					case <-ctx.Done():
						return
					}
				}
				// API includes data of the first date, so we can skip it
				dates = dates[1:]
				for _, date := range dates {
					if ctx.Err() != nil {
						return
					}
					if err := limiter.Wait(ctx); err != nil {
						select {
						case errChan <- err:
						case <-ctx.Done():
						}
						return
					}
					newData, err := s.fetchServerData(ctx, serviceID, &date)
					if err != nil {
						select {
						case errChan <- err:
						case <-ctx.Done():
						}
						return
					}
					for k, v := range newData.Data.Masters {
						initialData.Data.Masters[k] = v
					}
					for k, v := range newData.Data.Times {
						if _, exists := initialData.Data.Times[k]; !exists {
							initialData.Data.Times[k] = v
							continue
						}
						initialData.Data.Times[k] = append(initialData.Data.Times[k], newData.Data.Times[k]...)
					}
					newData.Data.ServiceID = serviceID
				}
				select {
				case results <- *initialData:
					return
				case <-ctx.Done():
					return
				}
			}()
		}
		wg.Wait()
	}()

	return results, errChan
}

func (s *pollingService) fetchServerData(ctx context.Context, serviceID int, date *string) (*ServerData, error) {
	u, err := url.Parse("https://dikidi.net/ru/mobile/ajax/newrecord/get_datetimes/?company_id=550001")
	if err != nil {
		return nil, &ErrFetch{err: err, msg: "Failed to build url"}
	}
	q := u.Query()
	if date != nil {
		q.Set("date", *date)
	}
	q.Set("service_id[]", fmt.Sprintf("%d", serviceID))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, &ErrFetch{url: u.String(), err: err, msg: "Failed to create request"}
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		return nil, &ErrFetch{url: u.String(), err: err, msg: "Failed to fetch document"}
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, &ErrFetch{url: u.String(), err: errors.New("bad status code"), msg: fmt.Sprintf("Expected 200 but got %d", res.StatusCode)}
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
