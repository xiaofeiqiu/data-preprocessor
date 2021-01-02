package handlers

import (
	"encoding/json"
	"errors"
	"github.com/xiaofeiqiu/data-preprocessor/lib/log"
	"github.com/xiaofeiqiu/data-preprocessor/lib/restutils"
	"github.com/xiaofeiqiu/data-preprocessor/services/alphavantage"
	"github.com/xiaofeiqiu/data-preprocessor/services/dbservice"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"
)

func (api *ApiHandler) InsertDailyCandle(w http.ResponseWriter, r *http.Request) (int, error) {

	req := alphavantage.DailyRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return 400, errors.New("error reading request body, " + err.Error())
	}
	json.Unmarshal(body, &req)

	req.Function = alphavantage.FUNC_TIME_SERIES_DAILY_ADJUSTED
	req.DataType = alphavantage.CSV
	if req.OutputSize == "" {
		req.OutputSize = alphavantage.Compact
	}

	err = validate.Struct(req)
	if err != nil {
		return 400, errors.New("request validation failed, " + err.Error())
	}
	log.Info("InsertDailyCandle", "Valid InsertDailyCandle request")

	status, body, err := api.AlphaVantageClient.Call(req)
	if err != nil {
		return status, errors.New("error calling FUNC_TIME_SERIES_DAILY_ADJUSTED, " + err.Error())
	}
	log.Info("InsertDailyCandle", "Call alpha vantage successful")

	resp := []*dbservice.RawDataEntity{}
	if restutils.Is2xxStatusCode(status) {
		log.Info("InsertDailyCandle", "Alpha vantage returns 200 status code")
		resp, err = ReadCsvData(req.Symbol, body, CandleReader)
		if err != nil {
			return 500, errors.New("error reading response, " + err.Error())
		}
		SetChanges(resp)
		err = api.DBService.BulkInsertRawDataEntity(resp)
		if err != nil {
			return 500, err
		}
		restutils.ResponseWithJson(w, 200, "Insert successful")
		return 0, nil
	}

	if restutils.Is4xxStatusCode(status) || restutils.Is5xxStatusCode(status) {
		log.Info("InsertDailyCandle", "Alpha vantage returns "+strconv.Itoa(status))
		restutils.ResponseWithJson(w, status, string(body))
		return 0, nil
	}

	return 500, errors.New("unexpected error occurred")
}

func (api *ApiHandler) InsertMissingDailyCandle(w http.ResponseWriter, r *http.Request) (int, error) {

	req := alphavantage.DailyRequest{}
	body, err := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &req)

	req.Function = alphavantage.FUNC_TIME_SERIES_DAILY_ADJUSTED
	req.DataType = alphavantage.CSV
	if req.OutputSize == "" {
		req.OutputSize = alphavantage.Compact
	}

	err = validate.Struct(req)
	if err != nil {
		return 400, errors.New("request validation failed, " + err.Error())
	}
	log.Info("InsertMissingDailyCandle", "Valid InsertMissingDailyCandle request")

	status, body, err := api.AlphaVantageClient.Call(req)
	if err != nil {
		return status, errors.New("error calling FUNC_TIME_SERIES_DAILY_ADJUSTED, " + err.Error())
	}
	log.Info("InsertMissingDailyCandle", "Call alpha vantage successful")

	resp := []*dbservice.RawDataEntity{}
	if restutils.Is2xxStatusCode(status) {
		log.Info("InsertMissingDailyCandle", "Alpha vantage returns 200 status code")
		resp, err = ReadCsvData(req.Symbol, body, CandleReader)
		if err != nil {
			return 500, errors.New("error reading response, " + err.Error())
		}
		SetChanges(resp)
		FixPrice(req.Symbol, resp)
		api.DBService.InsertRawDataEntityIgnoreError(resp)
		restutils.ResponseWithJson(w, 200, "Insert missing successful")
		return 0, nil
	}

	if restutils.Is4xxStatusCode(status) || restutils.Is5xxStatusCode(status) {
		log.Info("InsertMissingDailyCandle", "Alpha vantage returns "+strconv.Itoa(status))
		restutils.ResponseWithJson(w, status, string(body))
		return 0, nil
	}

	return 500, errors.New("unexpected error occurred")
}

func FixPrice(symbol string, data []*dbservice.RawDataEntity) error {

	if symbol == "AAPL" {
		str := "2020-08-29"
		t, err := time.Parse("2006-01-02", str)
		if err != nil {
			return err
		}
		for _, v := range data {
			if v.Date.Before(t) {
				v.Close = v.Close / 4
				v.Close = math.Round(v.Close*math.Pow(10, 2)) / math.Pow(10, 2)
			}
		}
	}
	return nil
}
