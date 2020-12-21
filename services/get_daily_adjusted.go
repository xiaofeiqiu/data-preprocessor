package services

import (
	"errors"
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

var validate = validator.New()

const TIME_SERIES_DAILY_ADJUSTED = "TIME_SERIES_DAILY_ADJUSTED"

type DailyRequest struct {
	Symbol     string `validate:"required" schema:"symbol"`
	OutputSize string `schema:"outputsize"`
	Function   string `validate:"required"`
	DataType   string `validate:"required"`
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

func (h *AlphaVantageApi) GetUrl(req DailyRequest) string {
	return fmt.Sprintf("%s%s?%s=%s&%s=%s&%s=%s&%s=%s&%s=%s", h.Host, Path, Symbol, req.Symbol, Function, req.Function, OutputSize, req.OutputSize, ApiKey, h.ApiKey, DataType, req.DataType)
}
