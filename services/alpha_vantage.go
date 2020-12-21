package services

import (
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