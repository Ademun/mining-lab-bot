package polling

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Ademun/mining-lab-bot/pkg/model"
)

var labNameRegexp = regexp.MustCompile(`\p{L}+\s+\p{L}+\s+№\s*(\d+).*?\((\d+)\s*\p{L}+\.\)(?:.*?\))?\s*(\p{L}.+)$`)

func ParseServiceData(data *ServiceData, serviceID int) ([]model.Slot, error) {
	masters := data.Data.Masters

	if len(masters.MasterMap) == 0 {
		return nil, nil
	}

	times := data.Data.Times

	slots := make([]model.Slot, 0, len(masters.MasterMap))

	errs := make([]error, 0)
	for id, master := range masters.MasterMap {
		labNumber, labAuditorium, labName, err := parseMasterName(master.Username)
		if err != nil {
			labNumber, labAuditorium, labName, err = parseMasterName(master.ServiceName)
			if err != nil {
				errs = append(errs, err)
				continue
			}
		}

		var labType model.LabType
		switch {
		case strings.Contains(master.ServiceName, "Выполнение"):
			labType = model.LabPerformance
		case strings.Contains(master.ServiceName, "Защита"):
			labType = model.LabDefence
		default:
			errs = append(errs, &ErrParseData{
				data: master.ServiceName,
				err:  errors.New("invalid lab type in service name"),
			})
			continue
		}

		availableTimes := make([]time.Time, 0)
		for _, timeString := range times.TimesMap[id] {
			timestamp, err := parseTimeString(timeString)
			if err != nil {
				errs = append(errs, &ErrParseData{
					data: timeString,
					err:  err,
				})
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
				URL:           buildURL(serviceID),
			}

			slots = append(slots, slot)
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return slots, nil
}

func parseMasterName(masterName string) (int, int, string, error) {
	matches := labNameRegexp.FindStringSubmatch(masterName)
	if len(matches) < 4 {
		return 0, 0, "", &ErrParseData{
			data: masterName,
			err:  errors.New("invalid lab name format"),
		}
	}
	labNumber, _ := strconv.Atoi(matches[1])
	labAuditory, _ := strconv.Atoi(matches[2])
	labName := matches[3]
	return labNumber, labAuditory, labName, nil
}

func parseTimeString(timeString string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", timeString)
}

func buildURL(serviceID int) string {
	return "https://dikidi.net/550001?p=3.pi-po-ssm-sd&o=7&s=" + strconv.Itoa(serviceID) + "&rl=0_undefined"
}
