package polling

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/PuerkitoBio/goquery"
)

func (s *pollingService) fetchServiceIDs(ctx context.Context) ([]int, error) {
	doc, err := s.fetchDocument(ctx)
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

func (s *pollingService) fetchDocument(ctx context.Context) (*goquery.Document, error) {
	res, err := s.fetchData(ctx, s.options.ServiceURL)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, &ErrParseData{msg: "failed to parse id list document", err: err}
	}
	return doc, nil
}
