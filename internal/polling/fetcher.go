package polling

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

func (s *pollingService) fetchData(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, &ErrFetch{url: url, msg: "failed to create request", err: err}
	}

	if err := s.fetchRateLimiter.Wait(ctx); err != nil {
		return nil, &ErrFetch{url: url, msg: "rate limiting error", err: err}
	}

	res, err := s.httpClient.Do(req)
	if err != nil {
		return nil, &ErrFetch{url: url, msg: "failed to fetch data", err: err}
	}

	if res.StatusCode != http.StatusOK {
		err := s.processBadHTTPResponse(res)
		return nil, &ErrFetch{url: url, msg: "bad response", err: err}
	}
	s.increaseFetchRate()

	return res, nil
}

func (s *pollingService) processBadHTTPResponse(res *http.Response) error {
	switch res.StatusCode {
	case http.StatusTooManyRequests:
		s.decreaseFetchRate()
		return fmt.Errorf("too many requests")
	case http.StatusInternalServerError:
		return fmt.Errorf("internal server error")
	}
	return fmt.Errorf("unexpected status code: %d", res.StatusCode)
}

func (s *pollingService) increaseFetchRate() {
	newRateFloat := float64(s.options.GetFetchRate().Milliseconds()) * s.options.RecoveryFactor
	newRateFloat = math.Min(float64(s.options.MaxFetchRate.Milliseconds()), newRateFloat)
	newRate := time.Millisecond * time.Duration(math.Round(newRateFloat))
	s.mu.Lock()
	defer s.mu.Unlock()
	s.options.SetFetchRate(newRate)
	s.fetchRateLimiter.SetLimit(rate.Every(newRate))
}

func (s *pollingService) decreaseFetchRate() {
	newRateFloat := float64(s.options.GetFetchRate().Milliseconds()) / s.options.BackoffFactor
	newRateFloat = math.Max(float64(s.options.MinFetchRate.Milliseconds()), newRateFloat)
	newRate := time.Millisecond * time.Duration(math.Round(newRateFloat))
	s.mu.Lock()
	defer s.mu.Unlock()
	s.options.SetFetchRate(newRate)
	s.fetchRateLimiter.SetLimit(rate.Every(newRate))
}
