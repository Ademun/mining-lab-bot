package service

import (
	errors2 "errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var labNameRegexp = regexp.MustCompile(`\p{L}+\s+\p{L}+\s+№\s*(\d+).*?\((\d+)\s*\p{L}+\.\)(?:.*?\))?\s*(\p{L}.+)$`)

func ParseServiceData(data *ServiceData) (map[int]*LabData, error) {
	masters := data.Data.Masters

	if len(masters.MasterMap) == 0 {
		return nil, nil
	}

	times := data.Data.Times

	labs := make(map[int]*LabData, len(masters.MasterMap))

	errors := make([]error, 0)
	for id, master := range masters.MasterMap {
		labNumber, labAuditory, labName, err := parseLabName(master)
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to parse lab name %s: %w", master.Username, err))
			continue
		}

		var labType LabType
		switch {
		case strings.Contains(master.ServiceName, "Выполнение"):
			labType = LabPerformance
		case strings.Contains(master.ServiceName, "Защита"):
			labType = LabDefence
		default:
			errors = append(errors, fmt.Errorf("failed to parse lab type: lab name should contain either 'Выполнение' or 'Защита'"))
			continue
		}

		availableTimes := make([]time.Time, 0)
		for _, timeString := range times.TimesMap[id] {
			timestamp, err := parseTimeString(timeString)
			if err != nil {
				errors = append(errors, fmt.Errorf("failed to parse timestamp %s: %w", timeString, err))
				continue
			}
			availableTimes = append(availableTimes, timestamp)
		}

		lab := &LabData{
			Name:      labName,
			Number:    labNumber,
			Auditory:  labAuditory,
			Type:      labType,
			Available: availableTimes,
		}

		labs[id] = lab
	}

	if len(errors) > 0 {
		return nil, errors2.Join(errors...)
	}

	return labs, nil
}

func parseLabName(master MasterData) (int, int, string, error) {
	matches := labNameRegexp.FindStringSubmatch(master.Username)
	if len(matches) < 4 {
		return 0, 0, "", fmt.Errorf("failed to parse lab name: %s", master.Username)
	}
	labNumber, _ := strconv.Atoi(matches[1])
	labAuditory, _ := strconv.Atoi(matches[2])
	labName := matches[3]
	return labNumber, labAuditory, labName, nil
}

func parseTimeString(timeString string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", timeString)
}
