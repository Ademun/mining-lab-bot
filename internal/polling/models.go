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
	ServiceID int
	Masters   Masters  `json:"masters"`
	DatesTrue []string `json:"dates_true"`
	Times     Times    `json:"times"`
}

type Masters map[int]MasterData

func (m *Masters) UnmarshalJSON(b []byte) error {
	var emptySlice []interface{}
	if err := json.Unmarshal(b, &emptySlice); err == nil {
		*m = make(map[int]MasterData)
		return nil
	}

	var masterMap map[int]MasterData
	if err := json.Unmarshal(b, &masterMap); err == nil {
		*m = masterMap
		return nil
	}

	return fmt.Errorf("unknown masters format")
}

type MasterData struct {
	Username    string `json:"username"`
	ServiceName string `json:"service_name"`
}

type Times map[int][]string

func (t *Times) UnmarshalJSON(b []byte) error {
	var emptySlice []interface{}
	if err := json.Unmarshal(b, &emptySlice); err == nil {
		*t = make(map[int][]string)
		return nil
	}
	var timesMap map[int][]string
	if err := json.Unmarshal(b, &timesMap); err == nil {
		*t = timesMap
		return nil
	}
	return fmt.Errorf("unknown times format")
}

// ======================================================
