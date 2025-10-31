package polling

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	numRe      = regexp.MustCompile(`№\s*(\d+)`)
	audRe      = regexp.MustCompile(`\((\d+)\s*\p{L}+\.\)`)
	orderRe    = regexp.MustCompile(`\((\d+)-?\p{L}*\s*место\)`)
	domainRe   = regexp.MustCompile(`\b(Электричество|Механика|Виртуальная\s*лаб\.?)\b`)
	typePrefix = "Аудиторная"
)

func (s *pollingService) ParseServerData(ctx context.Context, data *ServerData, serviceID int) ([]Slot, error) {
	dataMasters := data.Data.Masters
	if len(dataMasters) == 0 {
		return nil, nil
	}

	slots := make([]Slot, 0, len(dataMasters))
	errs := make([]error, 0)

	for id, master := range dataMasters {
		slot, err := parseSlotInfo(master.Username, master.ServiceName)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		dataTimes := data.Data.Times
		timesTeachers := make(map[time.Time][]string, len(dataTimes))
		for _, timeString := range dataTimes[id] {
			timestamp, err := parseTimeString(timeString)
			if err != nil {
				errs = append(errs, &ErrParseData{
					data: timeString,
					err:  err,
				})
				continue
			}
			teachers := s.teacherService.FindTeachersForTime(ctx, timestamp, slot.Auditorium)
			teacherNames := make([]string, 0, len(teachers))
			for _, teacher := range teachers {
				teacherNames = append(teacherNames, teacher.Name)
			}
			timesTeachers[timestamp.Round(0)] = teacherNames
		}

		slot.TimesTeachers = timesTeachers
		slot.URL = buildURL(serviceID)

		slots = append(slots, *slot)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return slots, nil
}

func parseSlotInfo(username, serviceName string) (*Slot, error) {
	number, err := parseNumber(username, serviceName)
	if err != nil {
		return nil, err
	}
	auditorium, err := parseAuditorium(username, serviceName)
	if err != nil {
		return nil, err
	}
	order, err := parseOrder(username, serviceName)
	if err != nil {
		return nil, err
	}
	domain := parseDomain(serviceName)
	labType := parseType(username)

	name := parseName(username)

	return &Slot{
			Name:       name,
			Number:     number,
			Auditorium: auditorium,
			Order:      order,
			Domain:     domain,
			Type:       labType,
		},
		nil
}

func parseName(username string) string {
	name := numRe.ReplaceAllString(username, "")
	name = audRe.ReplaceAllString(name, "")
	name = orderRe.ReplaceAllString(name, "")
	name = strings.TrimPrefix(name, "Лабораторная работа")
	name = strings.TrimSpace(name)
	return strings.Join(strings.Fields(name), " ")
}

func parseNumber(username string, serviceName string) (int, error) {
	if match := numRe.FindStringSubmatch(username); match != nil {
		labNum, _ := strconv.Atoi(match[1])
		return labNum, nil
	} else if match := numRe.FindStringSubmatch(serviceName); match != nil {
		labNum, _ := strconv.Atoi(match[1])
		return labNum, nil
	}
	return 0, &ErrParseData{data: username + " " + serviceName, msg: "lab number not found", err: errors.New("invalid lab name format")}
}

func parseAuditorium(username string, serviceName string) (int, error) {
	if match := audRe.FindStringSubmatch(username); match != nil {
		labAud, _ := strconv.Atoi(match[1])
		return labAud, nil
	} else if match := audRe.FindStringSubmatch(serviceName); match != nil {
		labAud, _ := strconv.Atoi(match[1])
		return labAud, nil
	}
	return 0, &ErrParseData{data: username + " " + serviceName, msg: "lab auditorium not found", err: errors.New("invalid lab name format")}
}

func parseOrder(username string, serviceName string) (int, error) {
	if match := orderRe.FindStringSubmatch(username); match != nil {
		labOrder, _ := strconv.Atoi(match[1])
		return labOrder, nil
	} else if match := orderRe.FindStringSubmatch(serviceName); match != nil {
		labOrder, _ := strconv.Atoi(match[1])
		return labOrder, nil
	}
	return 0, &ErrParseData{data: username + " " + serviceName, msg: "lab order not found", err: errors.New("invalid lab name format")}
}

func parseDomain(serviceName string) LabDomain {
	if match := domainRe.FindStringSubmatch(serviceName); match != nil {
		switch match[1] {
		case "Электричество":
			return LabDomainElectricity
		case "Механика":
			return LabDomainMechanics
		case "Виртуальная лаб.":
			return LabDomainVirtual
		}
	}
	return LabDomainVirtual
}

func parseType(username string) LabType {
	if strings.Contains(username, typePrefix) {
		return LabTypeDefence
	}
	return LabTypePerformance
}

func parseTimeString(timeString string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", timeString)
}

func buildURL(serviceID int) string {
	return "https://dikidi.net/550001?p=3.pi-po-ssm-sd&o=7&s=" + strconv.Itoa(serviceID) + "&rl=0_undefined"
}
