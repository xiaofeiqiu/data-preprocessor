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
	"strconv"
	"time"
)

var validMacdPeriod = []string{"20200200"}

func (api *ApiHandler) FillDailyMacd(w http.ResponseWriter, r *http.Request) (int, error) {

	// create request and validate
	req, err := NewMacdRequest(r)
	if err != nil {
		return 400, errors.New("error creating new Macd request, " + err.Error())
	}
	log.Info("FillDailyMacd", "Create new Macd request successful")

	err = validate.Struct(req)
	if err != nil {
		return 400, errors.New("request validation failed, " + err.Error())
	}
	log.Info("FillDailyMacd", "Valid daily Macd request")

	// find nil Macd entries
	var entries []dbservice.RawDataEntity
	p := req.FastPeriod + req.SlowPeriod + req.SignalPeriod
	api.DBService.FindNullColEntries(&entries, req.Symbol, "macd_"+p)
	if len(entries) == 0 {
		err = fmt.Errorf("0 nil Macd found")
		log.Error("FillDailyMacd", err, "")
		return 400, err
	}
	log.Info("FillDailyMacd", fmt.Sprintf("Find %d nil Macd entries", len(entries)))

	// call alpha to get Macd values, if 200, update entries in db
	status, body, err := api.AlphaVantageClient.Call(req)
	if err != nil {
		return status, errors.New("error calling FUNC_Macd, " + err.Error())
	}
	log.Info("FillDailyMacd", "Call alpha vantage successful")

	var MacdResp []*dbservice.RawDataEntity
	if restutils.Is2xxStatusCode(status) {
		log.Info("FillDailyMacd", "Alpha vantage returns 200")
		MacdResp, err = ReadCsvData(req.Symbol, body, MACD_Reader)
		if err != nil {
			return 500, errors.New("error reading response, " + err.Error())
		}
		log.Info("FillDailyMacd", "Read Macd csv successful")

		ct := SetMacd(entries, MacdResp, p)
		log.Info("SetMacd", "Set Macd successful, count: "+strconv.Itoa(ct))

		c, err := api.DBService.UpdateEntries(entries)
		if err != nil {
			return 500, errors.New("error update Macd entries, " + err.Error())
		}
		log.Info("UpdateMacdEntries", fmt.Sprintf("updated %d rows", c))
		restutils.ResponseWithJson(w, 200, "successful")
		return 0, nil
	}

	return 500, errors.New("unexpected error occurred")
}

func NewMacdRequest(r *http.Request) (alphavantage.DailyRequest, error) {
	req := alphavantage.DailyRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return alphavantage.DailyRequest{}, err
	}
	json.Unmarshal(body, &req)

	req.Function = alphavantage.FUNC_MACD
	req.DataType = alphavantage.CSV
	req.SeriesType = alphavantage.SeriesTypeClose
	req.Interval = alphavantage.IntervalDaily

	if req.OutputSize == "" {
		req.OutputSize = alphavantage.Compact
	}

	p := req.FastPeriod + req.SlowPeriod + req.SignalPeriod
	isValid := validatePeriod(p, validMacdPeriod)
	if !isValid {
		err = fmt.Errorf("invalid period value, " + p)
		log.Error("NewMacdRequest", err, "")
		return alphavantage.DailyRequest{}, err
	}

	return req, nil
}

func SetMacd(entries []dbservice.RawDataEntity, Macds []*dbservice.RawDataEntity, period string) int {
	tmpMap := RawDataPtrArrayToMap(Macds)
	count := 0
	if period == "20200200" {
		for i, v := range entries {
			key := v.Date.Format(time.RFC3339)
			if tmpMap[key] != nil {
				entries[i].Macd_20_200_200 = tmpMap[key].Macd
				entries[i].Macd_Signal_20_200_200 = tmpMap[key].MacdSignal
				entries[i].Macd_Hist_20_200_200 = tmpMap[key].MacdHist
				count++
			}
		}
	}
	return count
}
