package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

var fileName string

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(fmt.Sprintf("Error loading .env file: %v", err))
	}
	fileName = os.Getenv("SERV_FILE_NAME")
}

func CheckAvailableLabs() ([]*LabData, error) {
	serviceIDs, err := readServiceIDs()
	if err != nil {
		return nil, fmt.Errorf("failed to read service IDs: %v", err)
	}

	idChan := make(chan int)
	resultChan := make(chan *LabData)
	limiter := rate.NewLimiter(rate.Every(2*time.Second), 1)

	results := make([]*LabData, 0)

	go func() {
		for _, serviceID := range serviceIDs {
			idChan <- serviceID
		}
		close(idChan)
	}()

	go func() {
		for range serviceIDs {
			result := <-resultChan
			results = append(results, result)
		}
	}()

	var g errgroup.Group
	for range 10 {
		g.Go(func() error {
			for id := range idChan {
				limiter.Wait(context.Background())
				data, err := FetchServiceData(id)
				if err != nil {
					return err
				}
				parsed, err := ParseServiceData(data)
				for _, labData := range parsed {
					resultChan <- labData
				}
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return results, nil
}

func readServiceIDs() ([]int, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var numbers []int
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&numbers)

	return numbers, err
}
