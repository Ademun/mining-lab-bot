package polling

import (
	"encoding/json"
	"fmt"
)

// Структура HTML документа для извлечения айди сервисов
// ======================================================

type PageOptions struct {
	StepData StepData `json:"step_data"`
}

type StepData struct {
	List []Category `json:"list"`
}

type Category struct {
	Services []LabService `json:"services"`
}

type LabService struct {
	ID int `json:"id"`
}

// ======================================================
// Структура JSON ответа от сервиса
// ======================================================

type ServerData struct {
	Data ServiceData `json:"data"`
}

type ServiceData struct {
	Company   Company  `json:"company"`
	Masters   Masters  `json:"masters"`
	DatesTrue []string `json:"dates_true"`
	Times     Times    `json:"times"`
}

type Company struct {
	ID string `json:"id"`
}

type Masters struct {
	// For some reason if the Masters field is empty, it is returned as an empty json array, but if it contains data, it is returned as a map
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
	// For some reason if the Times field is empty, it is returned as an empty json array, but if it contains data, it is returned as a map
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

// ======================================================
