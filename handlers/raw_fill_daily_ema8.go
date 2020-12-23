package handlers

import (
	"encoding/json"
	"errors"
	"github.com/xiaofeiqiu/data-preprocessor/lib/restutils"
	"github.com/xiaofeiqiu/data-preprocessor/services/alphavantage"
	"io/ioutil"
	"net/http"
)

var validPeriod = []string{"8", "30"}

func (api *ApiHandler) FillDailyEMA(w http.ResponseWriter, r *http.Request) (int, error) {

	req, err := NewEmaRequest(r)
	if err != nil {
		return 400, errors.New("error creating new ema request")
	}

	err = validate.Struct(req)
	if err != nil {
		return 400, errors.New("request validation failed, " + err.Error())
	}
	status, body, err := api.AlphaVantageClient.Call(req)
	if err != nil {
		return status, errors.New("error calling FUNC_EMA, " + err.Error())
	}

	var ema8Resp []*RawDataEntity
	if restutils.Is2xxStatusCode(status) {
		ema8Resp, err = ReadCsvData(req.Symbol, body, EMA_8_Reader)
		if err != nil {
			return 500, errors.New("error reading response, " + err.Error())
		}
		restutils.ResponseWithJson(w, 200, ema8Resp)
		return 0, nil
	}

	return 500, errors.New("unexpected error occurred")
}

func NewEmaRequest(r *http.Request) (alphavantage.DailyRequest, error) {
	req := alphavantage.DailyRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return alphavantage.DailyRequest{}, err
	}
	json.Unmarshal(body, &req)

	req.Function = alphavantage.FUNC_EMA
	req.DataType = alphavantage.CSV
	req.SeriesType = alphavantage.SeriesTypeClose
	req.Interval = alphavantage.IntervalDaily

	if req.OutputSize == "" {
		req.OutputSize = alphavantage.Compact
	}

	isValid := validatePeriod(req.TimePeriod)
	if !isValid{
		return alphavantage.DailyRequest{}, errors.New("invalid period value")
	}

	return req, nil
}

func validatePeriod(period string) bool {
	for _, v := range validPeriod {
		if period == v {
			return true
		}
	}
	return false
}
