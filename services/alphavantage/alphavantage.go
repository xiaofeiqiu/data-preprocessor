package alphavantage

import (
	"errors"
	"github.com/google/go-querystring/query"
	"github.com/gorilla/schema"
	"github.com/xiaofeiqiu/data-preprocessor/lib/restutils"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

const Path = "/query"

// data
const OutputSizeFull = "full"
const Compact = "compact"
const CSV = "csv"

// time
const Daily = "daily"
const TimePeriod8 = "8"
const TimePeriod60 = "60"
const Close = "close"

// functions
const EMA = "EMA"
const TIME_SERIES_DAILY_ADJUSTED = "TIME_SERIES_DAILY_ADJUSTED"

var decoder = schema.NewDecoder()
var validate = validator.New()

type AlphaVantageApi struct {
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
	EMA_Daily_8   float64 `json:ema_daily_8`
}

func (api *AlphaVantageApi) Call(function string, r *http.Request) (int, []byte, error) {

	req := DailyRequest{
		Function: function,
	}

	err := decoder.Decode(&req, r.URL.Query())
	if err != nil {
		return 400, nil, errors.New("error decoding the request, " + err.Error())
	}

	setDefaultParams(function, &req)

	err = validate.Struct(req)
	if err != nil {
		return 400, nil, errors.New("request validation failed, " + err.Error())
	}

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

func setDefaultParams(function string, request *DailyRequest) {

	request.DataType = CSV

	if request.OutputSize == "" {
		request.OutputSize = Compact
	}

	if function == EMA {
		if request.Interval == "" {
			request.Interval = Daily
		}

		if request.SeriesType == "" {
			request.SeriesType = Close
		}

		if request.TimePeriod == "" {
			request.TimePeriod = "8"
		}
	}
}

func (h *AlphaVantageApi) GetUrl(req DailyRequest) (string, error) {
	params, err := query.Values(req)
	if err != nil {
		return "", err
	}
	url := h.Host + Path + "?" + params.Encode() + "&apikey="+h.ApiKey
	return url, nil
}
