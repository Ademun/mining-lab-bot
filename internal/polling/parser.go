package polling

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Ademun/mining-lab-bot/pkg/model"
)

var (
	numRe   = regexp.MustCompile(`№\s*(\d+)`)
	audRe   = regexp.MustCompile(`\((\d+)\s*\p{L}+\.\)`)
	orderRe = regexp.MustCompile(`\((\d+)-?\p{L}*\s*место\)`)
)

func (s *pollingService) ParseServerData(ctx context.Context, data *ServerData, serviceID int) ([]model.Slot, error) {
	dataMasters := data.Data.Masters
	if len(dataMasters.MasterMap) == 0 {
		return nil, nil
	}

	dataTimes := data.Data.Times

	slots := make([]model.Slot, 0, len(dataMasters.MasterMap))

	errs := make([]error, 0)
	for id, master := range dataMasters.MasterMap {
		slot, err := parseLabName(master.Username, master.ServiceName)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		// TODO: Implement correct Defence detection when new slots are opened in November
		var labType = model.LabPerformance

		available := make([]model.TimeTeachers, 0)
		for _, timeString := range dataTimes.TimesMap[id] {
			timestamp, err := parseTimeString(timeString)
			if err != nil {
				errs = append(errs, &ErrParseData{
					data: timeString,
					err:  err,
				})
				continue
			}
			teachers := s.teacherService.FindTeachersForTime(ctx, timestamp, slot.LabAuditorium)
			available = append(available, model.TimeTeachers{
				Time:     timestamp,
				Teachers: teachers,
			})
		}

		slot.ID = id
		slot.LabType = labType
		slot.Available = available
		slot.URL = buildURL(serviceID)

		slots = append(slots, *slot)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return slots, nil
}

func parseLabName(username, serviceName string) (*model.Slot, error) {
	labNumber, err := parseLabNumber(username, serviceName)
	if err != nil {
		return nil, err
	}

	labAuditorium, err := parseLabAuditorium(username, serviceName)
	if err != nil {
		return nil, err
	}

	labOrder, err := parseLabOrder(username, serviceName)
	if err != nil {
		labOrder = 0
	}

	labName := username
	labName = numRe.ReplaceAllString(labName, "")
	labName = audRe.ReplaceAllString(labName, "")
	labName = orderRe.ReplaceAllString(labName, "")
	labName = strings.TrimPrefix(labName, "Лабораторная работа")
	labName = strings.TrimSpace(labName)

	return &model.Slot{
		LabName:       labName,
		LabNumber:     labNumber,
		LabAuditorium: labAuditorium,
		LabOrder:      labOrder,
	}, nil
}

func parseLabNumber(username string, serviceName string) (int, error) {
	if match := numRe.FindStringSubmatch(username); match != nil {
		labNum, _ := strconv.Atoi(match[1])
		return labNum, nil
	} else if match := numRe.FindStringSubmatch(serviceName); match != nil {
		labNum, _ := strconv.Atoi(match[1])
		return labNum, nil
	} else {
		return 0, &ErrParseData{data: username + " " + serviceName, msg: "lab number not found", err: errors.New("invalid lab name format")}
	}
}

func parseLabAuditorium(username string, serviceName string) (int, error) {
	if match := audRe.FindStringSubmatch(username); match != nil {
		labAud, _ := strconv.Atoi(match[1])
		return labAud, nil
	} else if match := audRe.FindStringSubmatch(serviceName); match != nil {
		labAud, _ := strconv.Atoi(match[1])
		return labAud, nil
	} else {
		return 0, &ErrParseData{data: username + " " + serviceName, msg: "lab auditorium not found", err: errors.New("invalid lab name format")}
	}
}

func parseLabOrder(username string, serviceName string) (int, error) {
	if match := orderRe.FindStringSubmatch(username); match != nil {
		labOrder, _ := strconv.Atoi(match[1])
		return labOrder, nil
	} else if match := orderRe.FindStringSubmatch(serviceName); match != nil {
		labOrder, _ := strconv.Atoi(match[1])
		return labOrder, nil
	} else {
		return 0, &ErrParseData{data: username + " " + serviceName, msg: "lab order not found", err: errors.New("invalid lab name format")}
	}
}

func parseTimeString(timeString string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", timeString)
}

func buildURL(serviceID int) string {
	return "https://dikidi.net/550001?p=3.pi-po-ssm-sd&o=7&s=" + strconv.Itoa(serviceID) + "&rl=0_undefined"
}
