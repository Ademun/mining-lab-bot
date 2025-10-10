package polling

import (
	errors2 "errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Ademun/mining-lab-bot/pkg/model"
)

var labNameRegexp = regexp.MustCompile(`\p{L}+\s+\p{L}+\s+№\s*(\d+).*?\((\d+)\s*\p{L}+\.\)(?:.*?\))?\s*(\p{L}.+)$`)

func ParseServiceData(data *ServiceData) ([]model.Slot, error) {
	masters := data.Data.Masters

	if len(masters.MasterMap) == 0 {
		return nil, nil
	}

	times := data.Data.Times

	slots := make([]model.Slot, 0, len(masters.MasterMap))

	errors := make([]error, 0)
	for id, master := range masters.MasterMap {
		labNumber, labAuditorium, labName, err := parseMasterName(master.Username)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		var labType model.LabType
		switch {
		case strings.Contains(master.ServiceName, "Выполнение"):
			labType = model.LabPerformance
		case strings.Contains(master.ServiceName, "Защита"):
			labType = model.LabDefence
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

		for _, dateTime := range availableTimes {
			slot := model.Slot{
				ID:            id,
				LabNumber:     labNumber,
				LabName:       labName,
				LabAuditorium: labAuditorium,
				LabType:       labType,
				DateTime:      dateTime,
			}

			slots = append(slots, slot)
		}
	}

	if len(errors) > 0 {
		return nil, errors2.Join(errors...)
	}

	return slots, nil
}

func parseMasterName(masterName string) (int, int, string, error) {
	matches := labNameRegexp.FindStringSubmatch(masterName)
	if len(matches) < 4 {
		return 0, 0, "", fmt.Errorf("failed to parse lab name: %s", masterName)
	}
	labNumber, _ := strconv.Atoi(matches[1])
	labAuditory, _ := strconv.Atoi(matches[2])
	labName := matches[3]
	return labNumber, labAuditory, labName, nil
}

func parseTimeString(timeString string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", timeString)
}
