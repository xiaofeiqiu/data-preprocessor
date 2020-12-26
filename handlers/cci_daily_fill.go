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

var validPeriod = []string{"20", "50", "100", "200"}

func (api *ApiHandler) FillDailyCCI(w http.ResponseWriter, r *http.Request) (int, error) {

	// create request and validate
	req, err := NewCciRequest(r)
	if err != nil {
		return 400, errors.New("error creating new cci request, " + err.Error())
	}
	log.Info("FillDailyCCI", "Create new cci request successful")

	err = validate.Struct(req)
	if err != nil {
		return 400, errors.New("request validation failed, " + err.Error())
	}
	log.Info("FillDailyCCI", "Valid daily cci request")

	// find nil cci entries
	var entries []dbservice.RawDataEntity
	api.DBService.FindNullColEntries(&entries, req.Symbol, "cci"+req.TimePeriod)
	if len(entries) == 0 {
		err = fmt.Errorf("0 nil cci found")
		log.Error("FillDailyCCI", err, "")
		return 400, err
	}
	log.Info("FillDailyCCI", fmt.Sprintf("Find %d nil cci entries", len(entries)))

	// call alpha to get cci values, if 200, update entries in db
	status, body, err := api.AlphaVantageClient.Call(req)
	if err != nil {
		return status, errors.New("error calling FUNC_cci, " + err.Error())
	}
	log.Info("FillDailyCCI", "Call alpha vantage successful")

	var cciResp []*dbservice.RawDataEntity
	if restutils.Is2xxStatusCode(status) {
		log.Info("FillDailyCCI", "Alpha vantage returns 200")
		cciResp, err = ReadCsvData(req.Symbol, body, CCI_Reader)
		if err != nil {
			return 500, errors.New("error reading response, " + err.Error())
		}
		log.Info("FillDailyCCI", "Read cci csv successful")

		ct := SetCci(entries, cciResp, req.TimePeriod)
		log.Info("SetCCI", "Set cci successful, count: "+strconv.Itoa(ct))

		c, err := api.DBService.UpdateEntries(entries)
		if err != nil {
			return 500, errors.New("error update cci entries, " + err.Error())
		}
		log.Info("UpdatecciEntries", fmt.Sprintf("updated %d rows", c))
		restutils.ResponseWithJson(w, 200, "successful")
		return 0, nil
	}

	return 500, errors.New("unexpected error occurred")
}

func NewCciRequest(r *http.Request) (alphavantage.DailyRequest, error) {
	req := alphavantage.DailyRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return alphavantage.DailyRequest{}, err
	}
	json.Unmarshal(body, &req)

	req.Function = alphavantage.FUNC_CCI
	req.DataType = alphavantage.CSV
	req.Interval = alphavantage.IntervalDaily

	if req.OutputSize == "" {
		req.OutputSize = alphavantage.Compact
	}

	isValid := validatePeriod(req.TimePeriod)
	if !isValid {
		err = fmt.Errorf("invalid period value")
		log.Error("NewcciRequest", err, "")
		return alphavantage.DailyRequest{}, err
	}

	return req, nil
}


func SetCci(entries []dbservice.RawDataEntity, ccis []*dbservice.RawDataEntity, period string) int {
	tmpMap := ToMap(ccis)
	count := 0
	if period == "100" {
		for i, v := range entries {
			key := v.Date.Format(time.RFC3339)
			if tmpMap[key] != nil {
				entries[i].CCI_100 = tmpMap[key].CCI
				count++
			}
		}
	}
	return count
}
