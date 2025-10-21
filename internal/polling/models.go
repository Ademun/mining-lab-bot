package polling

import (
	"encoding/json"
	"fmt"
)

// Структура HTML документа для извлечения айди сервисов
// ======================================================

type pageOptions struct {
	StepData stepData `json:"step_data"`
}

type stepData struct {
	List []category `json:"list"`
}

type category struct {
	Services []labService `json:"services"`
}

type labService struct {
	ID int `json:"id"`
}

// ======================================================
// Структура JSON ответа от сервиса
// ======================================================

type serverData struct {
	Data serviceData `json:"data"`
}

type serviceData struct {
	Masters   masters  `json:"masters"`
	DatesTrue []string `json:"dates_true"`
	Times     times    `json:"times"`
}

type masters struct {
	EmptySlice []interface{}
	MasterMap  map[int]MasterData
}

func (m *masters) UnmarshalJSON(b []byte) error {
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

type times struct {
	EmptySlice []interface{}
	TimesMap   map[int][]string
}

func (t *times) UnmarshalJSON(b []byte) error {
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
