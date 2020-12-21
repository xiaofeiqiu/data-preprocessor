package handlers

import (
	"encoding/json"
	"errors"
	"github.com/xiaofeiqiu/data-preprocessor/lib/restutils"
	"github.com/xiaofeiqiu/data-preprocessor/services/alphavantage"
	"io/ioutil"
	"net/http"
)

func (api *ApiHandler) InsertDailyCandle(w http.ResponseWriter, r *http.Request) (int, error) {

	req := alphavantage.DailyRequest{}
	body, err := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &req)

	req.Function = alphavantage.FUNC_TIME_SERIES_DAILY_ADJUSTED
	req.DataType = alphavantage.CSV
	if req.OutputSize == ""{
		req.OutputSize = alphavantage.Compact
	}

	err = validate.Struct(req)
	if err != nil {
		return 400, errors.New("request validation failed, " + err.Error())
	}

	status, body, err := api.AlphaVantageApi.Call(req)
	if err != nil {
		return status, errors.New("error calling FUNC_TIME_SERIES_DAILY_ADJUSTED, " + err.Error())
	}

	resp := []*alphavantage.DailyResponse{}
	if restutils.Is2xxStatusCode(status) {
		resp, err = ReadCsvData(req.Symbol, body, CandleReader)
		if err != nil {
			return 500, errors.New("error reading response, " + err.Error())
		}
		SetStats(resp)
		restutils.ResponseWithJson(w, 200, resp)
		return 0, nil
	}

	return 500, errors.New("unexpected error occurred")
}
