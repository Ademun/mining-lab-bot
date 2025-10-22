package polling

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func (s *pollingService) fetchServiceIDs(ctx context.Context) ([]int, error) {
	doc, err := s.fetchDocument(ctx, s.options.ServiceURL)
	if err != nil {
		return nil, err
	}

	serviceIDs := make([]int, 0)
	var parsingErr error
	doc.Find(".newrecord2").Each(func(_ int, s *goquery.Selection) {
		dataOptions, exists := s.Attr("data-options")
		if !exists {
			return
		}
		var pageOptions PageOptions
		err := json.Unmarshal([]byte(dataOptions), &pageOptions)
		if err != nil {
			newErr := &ErrParseData{err: err, data: dataOptions}
			parsingErr = errors.Join(parsingErr, newErr)
		}
		categories := pageOptions.StepData.List
		for _, category := range categories {
			for _, service := range category.Services {
				serviceIDs = append(serviceIDs, service.ID)
			}
		}
	})

	if parsingErr != nil {
		return nil, parsingErr
	}

	return serviceIDs, nil
}

func (s *pollingService) fetchDocument(ctx context.Context, url string) (*goquery.Document, error) {
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

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, &ErrFetch{url: url, err: err, msg: "Failed to parse document"}
	}

	return doc, nil
}
