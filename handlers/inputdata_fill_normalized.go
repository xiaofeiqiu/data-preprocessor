package handlers

import (
	"encoding/json"
	"errors"
	"github.com/xiaofeiqiu/data-preprocessor/lib/log"
	"github.com/xiaofeiqiu/data-preprocessor/lib/restutils"
	"github.com/xiaofeiqiu/data-preprocessor/services/dbservice"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func (api *ApiHandler) DataInputFillNEma(w http.ResponseWriter, r *http.Request) (int, error) {

	// create request and validate
	req, err := NewDataInputRequest(r)
	if err != nil {
		return 400, errors.New("error creating new ema request, " + err.Error())
	}
	log.Info("DataInputFillNEma", "Valid data input ema request")

	// find entries to fill
	entries, err := api.findEntriesToFill(req)
	if err != nil {
		return 500, err
	}

	// get raw data
	rawData, err := api.DBService.FindRawData(entries)
	if err != nil {
		return 500, errors.New("FindRawData failed, " + err.Error())
	}

	rawDataMap := RawDataArrayToMap(rawData)

	SetNormalizedEma(entries, rawDataMap)
	SetNormalizedCCI(entries, rawDataMap)
	SetNormalizedAroon(entries, rawDataMap)
	SetNormalizedMacd(entries, rawDataMap)

	//update to db
	ct, err := api.DBService.UpdateDataInput(entries)
	if err != nil {
		return 500, err
	}
	log.Info("DataInputFillNEma", strconv.Itoa(ct)+" data input inserted")
	restutils.ResponseWithJson(w, 200, strconv.Itoa(ct)+" updated")

	return 0, nil
}

func (api *ApiHandler) findEntriesToFill(req *DataInputRequest) ([]dbservice.DataInputEntity, error) {
	entries := []dbservice.DataInputEntity{}
	inputDataMap := map[string]*dbservice.DataInputEntity{}
	for _, colName := range req.GetColNames() {
		tmp := []dbservice.DataInputEntity{}
		err := api.DBService.FindNullDataInput(&tmp, req.Symbol, colName)
		if err != nil {
			return nil, errors.New("find null" + colName + " failed")
		}
		for i, v := range tmp {
			if inputDataMap[v.Date.Format(time.RFC3339)] == nil {
				entries = append(entries, tmp[i])
				inputDataMap[v.Date.Format(time.RFC3339)] = &tmp[i]
			}
		}
	}
	return entries, nil
}

func SetNormalizedMacd(entries []dbservice.DataInputEntity, rawDataMap map[string]*dbservice.RawDataEntity) {
	for i, v := range entries {
		if rawDataMap[v.Date.Format(time.RFC3339)] != nil {
			entries[i].N_Macd_20_200_200 = rawDataMap[v.Date.Format(time.RFC3339)].GetNormalizedMacd()
		}
	}
}

func SetNormalizedAroon(entries []dbservice.DataInputEntity, rawDataMap map[string]*dbservice.RawDataEntity) {
	for i, v := range entries {
		if rawDataMap[v.Date.Format(time.RFC3339)] != nil {
			entries[i].N_AroonUp_50 = rawDataMap[v.Date.Format(time.RFC3339)].GetNormalizedAroonUp()
			entries[i].N_AroonDown_50 = rawDataMap[v.Date.Format(time.RFC3339)].GetNormalizedAroonDown()
		}
	}
}

func SetNormalizedCCI(entries []dbservice.DataInputEntity, rawDataMap map[string]*dbservice.RawDataEntity) {
	for i, v := range entries {
		if rawDataMap[v.Date.Format(time.RFC3339)] != nil {
			entries[i].N_CCI_100 = rawDataMap[v.Date.Format(time.RFC3339)].GetNormalizedCCI()
		}
	}
}

func SetNormalizedEma(entries []dbservice.DataInputEntity, rawDataMap map[string]*dbservice.RawDataEntity) {

	for i, v := range entries {
		if rawDataMap[v.Date.Format(time.RFC3339)] != nil {
			dt := v.Date.Format(time.RFC3339)
			rawData := rawDataMap[dt]
			nema20 := rawData.GetNormalizedEMA(20)
			nema50 := rawData.GetNormalizedEMA(50)
			nema100 := rawData.GetNormalizedEMA(100)
			nema200 := rawData.GetNormalizedEMA(200)

			entries[i].N_EMA_20 = nema20
			entries[i].N_EMA_50 = nema50
			entries[i].N_EMA_100 = nema100
			entries[i].N_EMA_200 = nema200
		}
	}

	log.Info("SetNormalizedEma", "SetNormalizedEma successful")
}

func NewDataInputRequest(r *http.Request) (*DataInputRequest, error) {
	req := DataInputRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(body, &req)

	err = validate.Struct(req)
	if err != nil {
		return nil, errors.New("validation failed for data input, " + err.Error())
	}
	return &req, nil
}
