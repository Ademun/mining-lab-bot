package polling

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func FetchServiceData(ctx context.Context, serviceID int) (*ServiceData, error) {
	url := fmt.Sprintf("https://dikidi.net/ru/mobile/ajax/newrecord/get_datetimes/?company_id=550001&service_id%%5B%%5D=%d&with_first=1&day_month=", serviceID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, &ErrFetch{err: err, msg: "Failed to create request"}
	}

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, &ErrFetch{err: err, msg: "Failed to fetch document"}
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, &ErrFetch{err: errors.New("bad status code"), msg: fmt.Sprintf("Expected 200 but got %d", res.StatusCode)}
	}

	return unmarshalServiceData(res.Body)
}

func unmarshalServiceData(serviceData io.ReadCloser) (*ServiceData, error) {
	var data ServiceData
	if err := json.NewDecoder(serviceData).Decode(&data); err != nil {
		return nil, &ErrParseData{err: err}
	}
	return &data, nil
}
