package scraper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
)

var url string
var fileName string

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(fmt.Sprintf("Error loading .env file: %v", err))
	}
	url = os.Getenv("SERV_ID_SRC")
	fileName = os.Getenv("SERV_FILE_NAME")
}

func UpdateServiceIDs() error {
	doc := fetchDocument()

	serviceIDs := make([]int, 0)
	doc.Find(".newrecord2").Each(func(i int, s *goquery.Selection) {
		dataOptions, exists := s.Attr("data-options")
		if !exists {
			return
		}
		var pageOptions PageOptions
		err := json.Unmarshal([]byte(dataOptions), &pageOptions)
		if err != nil {
			log.Fatal(fmt.Sprintf("Failed to unmarshal json pageOptions %v", err))
		}
		categories := pageOptions.StepData.List
		for _, category := range categories {
			for _, service := range category.Services {
				serviceIDs = append(serviceIDs, service.ID)
			}
		}
	})

	return write(serviceIDs)
}

func fetchDocument() *goquery.Document {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to fetch card id page: %v", err))
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

func write(IDs []int) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(IDs)
}
