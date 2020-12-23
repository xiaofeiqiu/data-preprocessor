package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xiaofeiqiu/data-preprocessor/lib/log"
	"github.com/xiaofeiqiu/data-preprocessor/lib/restutils"
	"github.com/xiaofeiqiu/data-preprocessor/services/alphavantage"
	"io/ioutil"
	"net/http"
	"github.com/xiaofeiqiu/data-preprocessor/services/dbservice"
	"time"
)

var validPeriod = []string{"8", "30"}

func (api *ApiHandler) FillDailyEMA(w http.ResponseWriter, r *http.Request) (int, error) {

	// create request and validate
	req, err := NewEmaRequest(r)
	if err != nil {
		return 400, errors.New("error creating new ema request, " + err.Error())
	}
	log.Info("FillDailyEMA", "Create new ema request successful")

	err = validate.Struct(req)
	if err != nil {
		return 400, errors.New("request validation failed, " + err.Error())
	}
	log.Info("FillDailyEMA", "Valid daily ema request")

	// find nil ema entries
	var entries []dbservice.RawDataEntity
	api.DBService.FindNullEma(&entries, req.Symbol, req.TimePeriod)
	if len(entries) == 0 {
		err = fmt.Errorf("0 nil ema found")
		log.Error("FillDailyEMA", err, "")
		return 400, err
	}
	log.Info("FillDailyEMA", fmt.Sprintf("Find %d nil ema entries", len(entries)))

	// call alpha to get ema values, if 200, update entries in db
	status, body, err := api.AlphaVantageClient.Call(req)
	if err != nil {
		return status, errors.New("error calling FUNC_EMA, " + err.Error())
	}
	log.Info("FillDailyEMA", "Call alpha vantage successful")

	var emaResp []*dbservice.RawDataEntity
	if restutils.Is2xxStatusCode(status) {
		log.Info("FillDailyEMA", "Alpha vantage returns 200")
		emaResp, err = ReadCsvData(req.Symbol, body, EMA_Reader)
		if err != nil {
			return 500, errors.New("error reading response, " + err.Error())
		}
		log.Info("FillDailyEMA", "Read ema csv successful")

		SetEMA(entries, emaResp, req.TimePeriod)
		log.Info("SetEMA", "Set ema successful")

		c, err := api.DBService.UpdateEntries(entries)
		if err != nil {
			return 500, errors.New("error update ema entries, " + err.Error())
		}
		log.Info("UpdateEmaEntries", fmt.Sprintf("updated %d rows", c))
		restutils.ResponseWithJson(w, 200, "successful")
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
	if !isValid {
		err = fmt.Errorf("invalid period value")
		log.Error("NewEmaRequest", err, "")
		return alphavantage.DailyRequest{}, err
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

func SetEMA(entries []dbservice.RawDataEntity, emas []*dbservice.RawDataEntity, period string) {
	tmpMap := ToMap(emas)
	if period == "8" {
		for i, v := range entries {
			key := v.Date.Format(time.RFC3339)
			if tmpMap[key] != nil {
				entries[i].EMA_8 = tmpMap[key].EMA_8
			}
		}
	}
}
