package polling

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
)

var url string

func init() {
	if err := godotenv.Load(); err != nil {
		slog.Error(fmt.Sprintf("Error loading .env file: %v", err))
		os.Exit(1)
	}
	url = os.Getenv("SERV_ID_SRC")
}

func FetchServiceIDs(ctx context.Context) ([]int, error) {
	doc := fetchDocument(ctx)

	serviceIDs := make([]int, 0)
	doc.Find(".newrecord2").Each(func(i int, s *goquery.Selection) {
		dataOptions, exists := s.Attr("data-options")
		if !exists {
			return
		}
		var pageOptions PageOptions
		err := json.Unmarshal([]byte(dataOptions), &pageOptions)
		if err != nil {
			slog.Error(fmt.Sprintf("Failed to unmarshal json pageOptions %v", err))
			return
		}
		categories := pageOptions.StepData.List
		for _, category := range categories {
			for _, service := range category.Services {
				serviceIDs = append(serviceIDs, service.ID)
			}
		}
	})

	return serviceIDs, nil
}

func fetchDocument(ctx context.Context) *goquery.Document {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to construct a request: %v", err))
	}

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to fetch document: %v", err))
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatal(fmt.Sprintf("Failed to fetch card id page. Got status code %d", res.StatusCode))
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to parse response body: %v", err))
	}
	return doc
}
