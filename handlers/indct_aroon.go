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

var validAroonPeriod = []string{"50"}

func (api *ApiHandler) FillDailyAroon(w http.ResponseWriter, r *http.Request) (int, error) {

	// create request and validate
	req, err := NewAroonRequest(r)
	if err != nil {
		return 400, errors.New("error creating new Aroon request, " + err.Error())
	}
	log.Info("FillDailyAroon", "Create new Aroon request successful")

	err = validate.Struct(req)
	if err != nil {
		return 400, errors.New("request validation failed, " + err.Error())
	}
	log.Info("FillDailyAroon", "Valid daily Aroon request")

	// find nil Aroon entries
	var entries []dbservice.RawDataEntity
	api.DBService.FindNullColEntries(&entries, req.Symbol, "Aroonup"+req.TimePeriod)
	if len(entries) == 0 {
		err = fmt.Errorf("0 nil Aroon found")
		log.Error("FillDailyAroon", err, "")
		return 400, err
	}
	log.Info("FillDailyAroon", fmt.Sprintf("Find %d nil Aroon entries", len(entries)))

	// call alpha to get Aroon values, if 200, update entries in db
	status, body, err := api.AlphaVantageClient.Call(req)
	if err != nil {
		return status, errors.New("error calling FUNC_Aroon, " + err.Error())
	}
	log.Info("FillDailyAroon", "Call alpha vantage successful")

	var AroonResp []*dbservice.RawDataEntity
	if restutils.Is2xxStatusCode(status) {
		log.Info("FillDailyAroon", "Alpha vantage returns 200")
		AroonResp, err = ReadCsvData(req.Symbol, body, Aroon_Reader)
		if err != nil {
			return 500, errors.New("error reading response, " + err.Error())
		}
		log.Info("FillDailyAroon", "Read Aroon csv successful")

		ct := SetAroon(entries, AroonResp, req.TimePeriod)
		log.Info("SetAroon", "Set Aroon successful, count: "+strconv.Itoa(ct))

		c, err := api.DBService.UpdateEntries(entries)
		if err != nil {
			return 500, errors.New("error update Aroon entries, " + err.Error())
		}
		log.Info("UpdateAroonEntries", fmt.Sprintf("updated %d rows", c))
		restutils.ResponseWithJson(w, 200, "successful")
		return 0, nil
	}

	return 500, errors.New("unexpected error occurred")
}

func NewAroonRequest(r *http.Request) (alphavantage.DailyRequest, error) {
	req := alphavantage.DailyRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return alphavantage.DailyRequest{}, err
	}
	json.Unmarshal(body, &req)

	req.Function = alphavantage.FUNC_AROON
	req.DataType = alphavantage.CSV
	req.Interval = alphavantage.IntervalDaily

	if req.OutputSize == "" {
		req.OutputSize = alphavantage.Compact
	}

	isValid := validatePeriod(req.TimePeriod, validAroonPeriod)
	if !isValid {
		err = fmt.Errorf("invalid period value")
		log.Error("NewAroonRequest", err, "")
		return alphavantage.DailyRequest{}, err
	}

	return req, nil
}

func SetAroon(entries []dbservice.RawDataEntity, Aroons []*dbservice.RawDataEntity, period string) int {
	tmpMap := ToMap(Aroons)
	count := 0
	if period == "50" {
		for i, v := range entries {
			key := v.Date.Format(time.RFC3339)
			if tmpMap[key] != nil {
				entries[i].AroonUp_50 = tmpMap[key].AroonUp
				entries[i].AroonDown_50 = tmpMap[key].AroonDown
				count++
			}
		}
	}
	return count
}
