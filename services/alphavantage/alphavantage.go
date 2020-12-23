package alphavantage

import (
	"errors"
	"github.com/google/go-querystring/query"
	"github.com/xiaofeiqiu/data-preprocessor/lib/restutils"
)

const Path = "/query"

// data
const OutputSizeFull = "full"
const Compact = "compact"
const CSV = "csv"

// time
const IntervalDaily = "daily"
const TimePeriod8 = "8"
const TimePeriod60 = "60"
const SeriesTypeClose = "close"

// functions
const FUNC_EMA = "EMA"
const FUNC_TIME_SERIES_DAILY_ADJUSTED = "TIME_SERIES_DAILY_ADJUSTED"

type AlphaVantageClient struct {
	Host       string
	ApiKey     string
	HttpClient *restutils.HttpClient
}

type DailyRequest struct {
	Symbol     string `validate:"required" schema:"symbol" url:"symbol"`
	OutputSize string `schema:"outputsize" url:"outputsize"`
	Function   string `validate:"required" url:"function"`
	DataType   string `validate:"required" url:"datatype"`
	Interval   string `schema:"interval" url:"interval"`
	TimePeriod string `schema:"time_period" url:"time_period"`
	SeriesType string `schema:"series_type" url:"series_type"`
}

type RawDataEntity struct {
	Symbol        string   `json:"symbol" db:"symbol, primarykey"`
	Timestamp     string   `json:"timestamp" db:"timestamp, primarykey"`
	Open          float64  `json:"open" db:"open"`
	High          float64  `json:"high" db:"high"`
	Low           float64  `json:"low" db:"low"`
	Close         float64  `json:"close" db:"close"`
	AdjustedClose float64  `json:"adjusted_close" db:"adjusted_close"`
	Volume        int64    `json:"volume" db:"volume"`
	Change        float64  `json:"change" db:"change"`
	EMA_8         *float64 `json:"ema_8,omitempty" db:"ema_8"`
}

func (api *AlphaVantageClient) Call(req DailyRequest) (int, []byte, error) {

	url, err := api.GetUrl(req)
	if err != nil {
		return 400, nil, errors.New("error getting url, " + err.Error())
	}
	status, body, err := api.HttpClient.DoGet(url, nil)
	if err != nil {
		return 500, nil, errors.New("error calling alpha vantage, " + err.Error())
	}

	return status, body, nil
}

func (h *AlphaVantageClient) GetUrl(req DailyRequest) (string, error) {
	params, err := query.Values(req)
	if err != nil {
		return "", err
	}
	url := h.Host + Path + "?" + params.Encode() + "&apikey=" + h.ApiKey
	return url, nil
}
