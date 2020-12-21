package services

import (
	"errors"
	"net/http"
)

type DailyRequest struct {
	Symbol     string `validate:"required" schema:"symbol" url:"symbol"`
	OutputSize string `schema:"outputsize" url:"outputsize"`
	Function   string `validate:"required" url:"function"`
	DataType   string `validate:"required" url:"datatype"`
	Interval   string `schema:interval url:"interval"`
	TimePeriod string `schema:time_period url:"time_period"`
	SeriesType string `schema:series_type url:"series_type"`
}

type DailyResponse struct {
	Timestamp     string  `json:timestamp`
	Open          float64 `json:open`
	High          float64 `json:high`
	Low           float64 `json:low`
	Close         float64 `json:close`
	AdjustedClose float64 `json:adjusted_close`
	Volume        int64   `json:volume`
	Change        float64 `json:change`
	N_Close       float64 `json:n_close`
}

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
		req.OutputSize = OutputSizeCompact
	}

	err = validate.Struct(req)
	if err != nil {
		return 400, nil, errors.New("request validation failed, " + err.Error())
	}

	url, err := api.GetUrl(req)
	if err!=nil {
		return 400, nil, errors.New("error getting url, " + err.Error())
	}
	status, body, err := api.HttpClient.DoGet(url, nil)
	if err != nil {
		return 500, nil, errors.New("error calling alpha vantage, " + err.Error())
	}

	return status, body, nil
}
