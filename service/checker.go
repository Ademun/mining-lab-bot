package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync/atomic"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
)

var fileName string

func init() {
	if err := godotenv.Load(); err != nil {
		slog.Error(fmt.Sprintf("Error loading .env file: %v", err))
	}
	fileName = os.Getenv("SERV_FILE_NAME")
}

func CheckAvailableLabs() ([]*LabData, error) {
	slog.Info("Checking labs")
	start := time.Now()
	serviceIDs, err := readServiceIDs()
	if err != nil {
		return nil, fmt.Errorf("failed to read service IDs: %v", err)
	}

	labList := make([]*LabData, 0)
	results := make(chan *LabData)
	IDsChecked := atomic.Int64{}

	var processingErr error
	var eg errgroup.Group
	limiter := rate.NewLimiter(rate.Every(time.Second*1), 1)
	go func() {
		for _, serviceID := range serviceIDs {
			eg.Go(func() error {
				limiter.Wait(context.Background())
				start := time.Now()
				IDsChecked.Add(1)
				slog.Info(fmt.Sprintf("Processing service. id = %d [%d/%d]", serviceID, IDsChecked.Load(), len(serviceIDs)))
				data, err := FetchServiceData(serviceID)
				if err != nil {
					return err
				}
				res, err := ParseServiceData(data)
				if err != nil {
					return err
				}
				for _, lab := range res {
					results <- lab
				}
				slog.Info(fmt.Sprintf("Finished processing service. id %d. Time elapsed %s\n", serviceID, time.Since(start)))
				return nil
			})
		}
		if err := eg.Wait(); err != nil {
			processingErr = err
		}
		close(results)
	}()

	for result := range results {
		labList = append(labList, result)
	}

	slog.Info(fmt.Sprintf("Finished checking labs in %s. Total available labs %d", time.Since(start), len(labList)))
	return labList, processingErr
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
