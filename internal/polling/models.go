package polling

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"time"
)

// Модель записи на лабораторную работу
// ======================================================

type LabDomain int

const (
	LabDomainElectricity LabDomain = iota
	LabDomainMechanics
	LabDomainVirtual
)

func (ld LabDomain) String() string {
	switch ld {
	case LabDomainElectricity:
		return "Электричество"
	case LabDomainMechanics:
		return "Механика"
	case LabDomainVirtual:
		return "Виртуалка"
	default:
		return "Неизвестно"
	}
}

type LabType int

const (
	LabTypePerformance LabType = iota
	LabTypeDefence
)

func (lt LabType) String() string {
	switch lt {
	case LabTypePerformance:
		return "Выполнение"
	case LabTypeDefence:
		return "Защита"
	default:
		return "Неизвестно"
	}
}

type Slot struct {
	Type          LabType
	Name          string
	Number        int
	Auditorium    int
	Order         int // 0 if no order
	Domain        LabDomain
	TimesTeachers map[time.Time][]string
	URL           string
}

func (s Slot) Key() string {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(s)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(buf.Bytes())
	return fmt.Sprintf("%x", hash)
}

// ======================================================

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
