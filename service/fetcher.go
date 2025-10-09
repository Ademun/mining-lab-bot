package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchServiceData(serviceId int) (*ServiceData, error) {
	url := fmt.Sprintf("https://dikidi.net/ru/mobile/ajax/newrecord/get_datetimes/?company_id=550001&service_id%%5B%%5D=%d&with_first=1&day_month=", serviceId)
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch service data: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch service data, got response code: %d", res.StatusCode)
	}

	return unmarshalServiceData(res.Body)
}

func unmarshalServiceData(serviceData io.ReadCloser) (*ServiceData, error) {
	var data ServiceData
	if err := json.NewDecoder(serviceData).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal service data: %w", err)
	}
	return &data, nil
}
