package polling

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"sync"
)

func (s *pollingService) pollServerData(ctx context.Context) (chan ServerData, chan error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make(chan ServerData)
	errChan := make(chan error)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		defer close(errChan)
		defer close(results)

		wg := sync.WaitGroup{}

		for _, serviceID := range s.serviceIDs {
			wg.Add(1)
			serviceID := serviceID
			go func() {
				result, err := s.processSingleService(ctx, serviceID)
				if err != nil {
					select {
					case errChan <- err:
					case <-ctx.Done():
					}
				}
				select {
				case results <- *result:
				case <-ctx.Done():
				}
			}()
		}
		wg.Wait()
	}()

	return results, errChan
}

// The initial request retrieves a list of dates, which is used to request all available slots for the serviceID
func (s *pollingService) processSingleService(ctx context.Context, serviceID int) (*ServerData, error) {
	initialData, err := s.fetchServerData(ctx, serviceID, nil)
	if err != nil {
		return nil, err
	}
	initialData.Data.ServiceID = serviceID

	dates := initialData.Data.DatesTrue
	if len(dates) == 0 {
		return initialData, nil
	}

	// API includes data of the first date, so we can skip it
	for _, date := range dates[1:] {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		newData, err := s.fetchServerData(ctx, serviceID, &date)
		if err != nil {
			return nil, err
		}

		for k, v := range newData.Data.Masters {
			initialData.Data.Masters[k] = v
		}
		for k, v := range newData.Data.Times {
			if existing, exists := initialData.Data.Times[k]; exists {
				initialData.Data.Times[k] = append(existing, v...)
			} else {
				initialData.Data.Times[k] = v
			}
		}
	}

	return initialData, nil
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

	res, err := s.fetchData(ctx, u.String())
	if err != nil {
		return nil, err
	}

	return unmarshalServerData(res.Body, serviceID)
}

func unmarshalServerData(serverData io.ReadCloser, serviceID int) (*ServerData, error) {
	var data ServerData
	if err := json.NewDecoder(serverData).Decode(&data); err != nil {
		return nil, &ErrParseData{msg: fmt.Sprintf("serviceID: %d", serviceID), err: err}
	}
	return &data, nil
}
