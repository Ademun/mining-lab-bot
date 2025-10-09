package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type ServiceData struct {
	Data Data `json:"data"`
}

type Data struct {
	Masters   Masters  `json:"masters"`
	DatesTrue []string `json:"dates_true"`
	Times     Times    `json:"times"`
}

type Masters struct {
	EmptySlice []interface{}
	MasterMap  map[int]MasterData
}

func (m *Masters) UnmarshalJSON(b []byte) error {
	var emptySlice []interface{}
	if err := json.Unmarshal(b, &emptySlice); err == nil {
		m.EmptySlice = emptySlice
		return nil
	}

	var masterMap map[int]MasterData
	if err := json.Unmarshal(b, &masterMap); err == nil {
		m.MasterMap = masterMap
		return nil
	}

	return fmt.Errorf("unknown masters format")
}

type MasterData struct {
	Username    string `json:"username"`
	ServiceName string `json:"service_name"`
}

type Times struct {
	EmptySlice []interface{}
	TimesMap   map[int][]string
}

func (t *Times) UnmarshalJSON(b []byte) error {
	var emptySlice []interface{}
	if err := json.Unmarshal(b, &emptySlice); err == nil {
		t.EmptySlice = emptySlice
		return nil
	}
	var timesMap map[int][]string
	if err := json.Unmarshal(b, &timesMap); err == nil {
		t.TimesMap = timesMap
		return nil
	}
	return fmt.Errorf("unknown times format")
}

type LabData struct {
	Name      string
	Number    int
	Auditory  int
	Type      LabType
	Available []time.Time
}

func (ld LabData) String() string {
	result := strings.Builder{}
	result.WriteString(fmt.Sprintf("№%d. %s. %s\n", ld.Number, ld.Name, ld.Type.String()))
	result.WriteString(fmt.Sprintf("Аудитория %d\n", ld.Auditory))
	result.WriteString("Доступные даты:\n")
	for i, timestamp := range ld.Available {
		result.WriteString(timestamp.Format("2006-01-02 15:04:05"))
		if i != len(ld.Available)-1 {
			result.WriteString("\n")
		}
	}
	return result.String()
}

type LabType int

const (
	LabPerformance LabType = iota
	LabDefence
)

func (t LabType) String() string {
	switch t {
	case LabPerformance:
		return "Выполнение"
	case LabDefence:
		return "Защита"
	}
	return "Unknown"
}
