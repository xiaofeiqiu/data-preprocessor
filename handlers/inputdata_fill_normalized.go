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

	// find null col in data input table
	inputData := []dbservice.DataInputEntity{}
	for _, colName := range req.GetColNames() {
		tmp := []dbservice.DataInputEntity{}
		err = api.DBService.FindNullDataInput(&tmp, req.Symbol, colName)
		if err != nil {
			return 500, errors.New("find null" + colName + " failed")
		}
		inputData = append(inputData, tmp...)
	}

	// get raw data
	rawData, err := api.DBService.FindRawData(inputData)
	if err != nil {
		return 500, errors.New("FindRawData failed, " + err.Error())
	}

	rawDataMap := RawDataArrayToMap(rawData)

	SetNormalizedEma(inputData, rawDataMap)
	SetNormalizedCCI(inputData, rawDataMap)

	//update to db
	ct, err := api.DBService.UpdateDataInput(inputData)
	if err != nil {
		return 500, err
	}
	log.Info("DataInputFillNEma", strconv.Itoa(ct)+" data input inserted")
	restutils.ResponseWithJson(w, 200, strconv.Itoa(ct)+" updated")

	return 0, nil
}

func SetNormalizedCCI(entries []dbservice.DataInputEntity, rawDataMap map[string]*dbservice.RawDataEntity) {
	for i, v := range entries {
		if rawDataMap[v.Date.Format(time.RFC3339)] != nil {
			if rawDataMap[v.Date.Format(time.RFC3339)].CCI_100 != nil {
				entries[i].N_CCI_100 = rawDataMap[v.Date.Format(time.RFC3339)].GetNormalizedCCI()
			}
		}
	}
}

func SetNormalizedEma(entries []dbservice.DataInputEntity, rawDataMap map[string]*dbservice.RawDataEntity) {

	for i, v := range entries {
		if rawDataMap[v.Date.Format(time.RFC3339)] != nil {
			dt := v.Date.Format(time.RFC3339)
			rawData := rawDataMap[dt]
			nema20, err := rawData.GetNormalizedEMA(20)
			if err != nil {
				log.Error("SetNormalizedEma", err, "nema20")
			}
			nema50, err := rawData.GetNormalizedEMA(50)
			if err != nil {
				log.Error("SetNormalizedEma", err, "nema50")
			}
			nema100, err := rawData.GetNormalizedEMA(100)
			if err != nil {
				log.Error("SetNormalizedEma", err, "nema100")
			}
			nema200, err := rawData.GetNormalizedEMA(200)
			if err != nil {
				log.Error("SetNormalizedEma", err, "nema200")
			}
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
