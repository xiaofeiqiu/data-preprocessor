package handlers

import (
	"github.com/xiaofeiqiu/data-preprocessor/lib/restutils"
	"github.com/xiaofeiqiu/data-preprocessor/services/alphavantage"
	"gopkg.in/go-playground/validator.v9"
	"github.com/xiaofeiqiu/data-preprocessor/services/dbservice"
	"strings"
)

var validate = validator.New()

type ApiHandler struct {
	AlphaVantageClient *alphavantage.AlphaVantageClient
	DBService          *dbservice.DBService
	DefaultClient      *restutils.HttpClient
}

// data input
var allCols = "n_ema20,n_ema50,n_ema100,n_ema200,ndiff_ema20,ndiff_ema50,ndiff_ema100,ndiff_ema200,n_cci100,ndiff_cci100,n_aroonup50,n_aroondown50,ndiff_aroonup50,ndiff_aroondown50,n_macd_20200200,n_macd_signal_20200200,n_macd_hist_20200200,ndiff_macd_20200200,ndiff_macd_signal_20200200,ndiff_macd_hist_20200200,n_osc10,ndiff_osc10,buy,sell,hold,"

type DataInputRequest struct {
	Symbol     string `json:"symbol" validate:"required"`
	ColNames   string `json:colNames`
	DiffLength int    `json:diffLength`
}

func (r *DataInputRequest) GetColNames() []string {
	if r.ColNames != "" {
		return strings.Split(r.ColNames, ",")
	}
	return strings.Split(allCols, ",")
}
