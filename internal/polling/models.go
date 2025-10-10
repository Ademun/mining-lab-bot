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
	Services []Service `json:"services"`
}

type Service struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ======================================================

// Структура JSON ответа от сервиса
// ======================================================

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

// ======================================================
