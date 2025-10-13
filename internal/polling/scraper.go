package polling

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func FetchServiceIDs(ctx context.Context, url string) ([]int, error) {
	doc, err := fetchDocument(ctx, url)
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
			parsingErr = &ErrParseData{err: err, data: dataOptions}
		}
		categories := pageOptions.StepData.List
		for _, category := range categories {
			for _, service := range category.Services {
				serviceIDs = append(serviceIDs, service.ID)
			}
		}
	})

	if parsingErr != nil {
		return nil, err
	}

	return serviceIDs, nil
}

func fetchDocument(ctx context.Context, url string) (*goquery.Document, error) {
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

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, &ErrFetch{err: err, msg: "Failed to parse document"}
	}

	return doc, nil
}
