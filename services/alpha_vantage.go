package services

import (
	"fmt"
	"github.com/gorilla/schema"
	"github.com/xiaofeiqiu/data-preprocessor/lib/restutils"
)

const Path = "/query"

const Symbol = "symbol"
const OutputSize = "outputsize"
const Function = "function"
const ApiKey = "apikey"
const DataType = "datatype"

const OutputFull = "full"
const OutputCompact = "compact"
const DataTypeCsv = "csv"

var decoder = schema.NewDecoder()

type AlphaVantageApi struct {
	Host       string
	ApiKey     string
	HttpClient *restutils.HttpClient
}

type DailyRequest struct {
	Symbol     string `validate:"required" schema:"symbol"`
	OutputSize string `schema:"outputsize"`
	Function   string `validate:"required"`
	DataType   string `validate:"required"`
}

func (h *AlphaVantageApi) GetUrl(req DailyRequest) string {
	return fmt.Sprintf("%s%s?%s=%s&%s=%s&%s=%s&%s=%s&%s=%s", h.Host, Path, Symbol, req.Symbol, Function, req.Function, OutputSize, req.OutputSize, ApiKey, h.ApiKey, DataType, req.DataType)
}
