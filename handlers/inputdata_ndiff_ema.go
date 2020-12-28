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

func (api *ApiHandler) DataInputFillNDiffEma(w http.ResponseWriter, r *http.Request) (int, error) {

	// create request and validate
	req, err := NewDataInputNDiffEmaRequest(r)
	if err != nil {
		return 400, errors.New("error creating new NDiffEma request, " + err.Error())
	}
	log.Info("DataInputFillNDiffEma", "Valid data input NDiffEma request")

	// find null col in data input table
	inputData := []dbservice.DataInputEntity{}
	api.DBService.FindNullDataInput(&inputData, req.Symbol, req.ColName)

	// get raw data
	rawData, err := api.DBService.FindRawData(inputData)
	if err != nil {
		return 500, errors.New("FindRawData failed, " + err.Error())
	}

	SetNormalizedNDiffEma(inputData, rawData, req.DiffLength)

	//update to db
	ct, err := api.DBService.UpdateDataInput(inputData)
	if err != nil {
		return 500, err
	}
	log.Info("DataInputFillNDiffEma", strconv.Itoa(ct)+" data input inserted")
	restutils.ResponseWithJson(w, 200, "successful")

	return 0, nil
}

func SetNormalizedNDiffEma(entires []dbservice.DataInputEntity, rawData []dbservice.RawDataEntity, DiffLength int) {

	var rawDataSlice dbservice.RawDataEntitySlice
	rawDataSlice = rawData
	rawDataSlice.LoadNormalizedNDiffEma(20, DiffLength)
	rawDataSlice.LoadNormalizedNDiffEma(50, DiffLength)
	rawDataSlice.LoadNormalizedNDiffEma(100, DiffLength)
	rawDataSlice.LoadNormalizedNDiffEma(200, DiffLength)
	rawDataMap := RawDataArrayToMap(rawData)

	for i, v := range entires {
		if rawDataMap[v.Date.Format(time.RFC3339)] != nil {
			dt := v.Date.Format(time.RFC3339)
			rawData := rawDataMap[dt]
			nNDiffEma20, err := rawData.GetNormalizedNDiffEma(20)
			if err != nil {
				log.Error("SetNormalizedNDiffEma", err, "nNDiffEma20")
			}
			nNDiffEma50, err := rawData.GetNormalizedNDiffEma(50)
			if err != nil {
				log.Error("SetNormalizedNDiffEma", err, "nNDiffEma50")
			}
			nNDiffEma100, err := rawData.GetNormalizedNDiffEma(100)
			if err != nil {
				log.Error("SetNormalizedNDiffEma", err, "nNDiffEma100")
			}
			nNDiffEma200, err := rawData.GetNormalizedNDiffEma(200)
			if err != nil {
				log.Error("SetNormalizedNDiffEma", err, "nNDiffEma200")
			}
			entires[i].NDiff_EMA_20 = nNDiffEma20
			entires[i].NDiff_EMA_50 = nNDiffEma50
			entires[i].NDiff_EMA_100 = nNDiffEma100
			entires[i].NDiff_EMA_200 = nNDiffEma200
		}
	}

	log.Info("SetNormalizedNDiffEma", "SetNormalizedNDiffEma successful")
}

func NewDataInputNDiffEmaRequest(r *http.Request) (*DataInputRequest, error) {
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
