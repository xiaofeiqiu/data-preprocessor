package services

import (
	"errors"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

var validate = validator.New()

const TIME_SERIES_DAILY_ADJUSTED = "TIME_SERIES_DAILY_ADJUSTED"

func (api *AlphaVantageApi) GetDailyAdjusted(r *http.Request) (int, []byte, error) {
	req := DailyRequest{
		Function: TIME_SERIES_DAILY_ADJUSTED,
		DataType: DataTypeCsv,
	}

	err := decoder.Decode(&req, r.URL.Query())
	if err != nil {
		return 400, nil, errors.New("error decoding the request, " + err.Error())
	}

	if req.OutputSize == "" {
		req.OutputSize = OutputCompact
	}

	err = validate.Struct(req)
	if err != nil {
		return 400, nil, errors.New("request validation failed, " + err.Error())
	}

	url := api.GetUrl(req)
	status, body, err := api.HttpClient.DoGet(url, nil)
	if err != nil {
		return 500, nil, errors.New("error calling alpha vantage, " + err.Error())
	}

	return status, body, nil
}
