package services

import (
	"github.com/google/go-querystring/query"
	"github.com/gorilla/schema"
	"github.com/xiaofeiqiu/data-preprocessor/lib/restutils"
	"gopkg.in/go-playground/validator.v9"
)

const Path = "/query"

// data
const OutputSizeFull = "full"
const OutputSizeCompact = "compact"
const DataTypeCsv = "csv"

// time
const TimePeriodDaily = "daily"
const TimePeriod8 = "8"
const TimePeriod60 = "60"
const SeriesTypeClose = "close"

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

func (h *AlphaVantageApi) GetUrl(req DailyRequest) (string, error) {
	params, err := query.Values(req)
	if err != nil {
		return "", err
	}
	url := h.Host + Path + "?" + params.Encode() + "&apikey="+h.ApiKey
	return url, nil
}
