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

var validOscPeriod = []string{"10"}

func (api *ApiHandler) FillDailyOSC(w http.ResponseWriter, r *http.Request) (int, error) {

	// create request and validate
	req, err := ƒ(r)
	if err != nil {
		return 400, errors.New("error creating new Osc request, " + err.Error())
	}
	log.Info("FillDailyOsc", "Create new Osc request successful")

	err = validate.Struct(req)
	if err != nil {
		return 400, errors.New("request validation failed, " + err.Error())
	}
	log.Info("FillDailyOsc", "Valid daily Osc request")

	// find nil Osc entries
	var entries []dbservice.RawDataEntity
	api.DBService.FindNullColEntries(&entries, req.Symbol, "osc"+req.TimePeriod)
	if len(entries) == 0 {
		err = fmt.Errorf("0 nil Osc found")
		log.Error("FillDailyOsc", err, "")
		return 400, err
	}
	log.Info("FillDailyOsc", fmt.Sprintf("Find %d nil Osc entries", len(entries)))

	// call alpha to get Osc values, if 200, update entries in db
	status, body, err := api.AlphaVantageClient.Call(req)
	if err != nil {
		return status, errors.New("error calling FUNC_Osc, " + err.Error())
	}
	log.Info("FillDailyOsc", "Call alpha vantage successful")

	var OscResp []*dbservice.RawDataEntity
	if restutils.Is2xxStatusCode(status) {
		log.Info("FillDailyOsc", "Alpha vantage returns 200")
		OscResp, err = ReadCsvData(req.Symbol, body, Osc_Reader)
		if err != nil {
			return 500, errors.New("error reading response, " + err.Error())
		}
		log.Info("FillDailyOsc", "Read Osc csv successful")

		ct := SetOsc(entries, OscResp, req.TimePeriod)
		log.Info("SetOsc", "Set Osc successful, count: "+strconv.Itoa(ct))

		c, err := api.DBService.UpdateEntries(entries)
		if err != nil {
			return 500, errors.New("error update Osc entries, " + err.Error())
		}
		log.Info("UpdateOscEntries", fmt.Sprintf("updated %d rows", c))
		restutils.ResponseWithJson(w, 200, "successful")
		return 0, nil
	}

	return 500, errors.New("unexpected error occurred")
}

func ƒ(r *http.Request) (alphavantage.DailyRequest, error) {
	req := alphavantage.DailyRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return alphavantage.DailyRequest{}, err
	}
	json.Unmarshal(body, &req)

	req.Function = alphavantage.FUNC_AROON_OSC
	req.DataType = alphavantage.CSV
	req.Interval = alphavantage.IntervalDaily

	if req.OutputSize == "" {
		req.OutputSize = alphavantage.Compact
	}

	isValid := validatePeriod(req.TimePeriod, validOscPeriod)
	if !isValid {
		err = fmt.Errorf("invalid period value")
		log.Error("ƒ", err, "")
		return alphavantage.DailyRequest{}, err
	}

	return req, nil
}

func SetOsc(entries []dbservice.RawDataEntity, Oscs []*dbservice.RawDataEntity, period string) int {
	tmpMap := RawDataPtrArrayToMap(Oscs)
	count := 0
	if period == "10" {
		for i, v := range entries {
			key := v.Date.Format(time.RFC3339)
			if tmpMap[key] != nil {
				entries[i].Osc_10 = tmpMap[key].Osc
				count++
			}
		}
	}
	return count
}
