package handlers

import (
	"encoding/json"
	"errors"
	"github.com/xiaofeiqiu/data-preprocessor/lib/log"
	"github.com/xiaofeiqiu/data-preprocessor/lib/restutils"
	"github.com/xiaofeiqiu/data-preprocessor/services/dbservice"
	"github.com/xiaofeiqiu/data-preprocessor/services/utils"
	"io/ioutil"
	"math"
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

	SetNormalizedNDiffEma(entries, rawData, req.DiffLength)

	//update to db
	ct, err := api.DBService.UpdateDataInput(entries)
	if err != nil {
		return 500, err
	}
	log.Info("DataInputFillNDiffEma", strconv.Itoa(ct)+" data input inserted")
	restutils.ResponseWithJson(w, 200, strconv.Itoa(ct)+" updated")

	return 0, nil
}

func SetNormalizedNDiffEma(entires []dbservice.DataInputEntity, rawData []dbservice.RawDataEntity, DiffLength int) {

	SetNDiffEma(rawData, 20, DiffLength)
	SetNDiffEma(rawData, 50, DiffLength)
	SetNDiffEma(rawData, 100, DiffLength)
	SetNDiffEma(rawData, 200, DiffLength)
	SetNDiffCCI(rawData, DiffLength)
	SetNDiffAroon(rawData, DiffLength)
	SetNDiffMacd(rawData, DiffLength)
	SetNDiffMacdHist(rawData, DiffLength)
	SetNDiffMacdOSC(rawData, DiffLength)
	rawDataMap := RawDataArrayToMap(rawData)

	for i, v := range entires {
		if rawDataMap[v.Date.Format(time.RFC3339)] != nil {
			dt := v.Date.Format(time.RFC3339)
			rawData := rawDataMap[dt]
			entires[i].NDiff_EMA_20 = rawData.NormalizedDiffNEMA20
			entires[i].NDiff_EMA_50 = rawData.NormalizedDiffNEMA50
			entires[i].NDiff_EMA_100 = rawData.NormalizedDiffNEMA100
			entires[i].NDiff_EMA_200 = rawData.NormalizedDiffNEMA200
			entires[i].NDiff_CCI_100 = rawData.NormalizedDiffNCCI100
			entires[i].NDiff_AroonDown_50 = rawData.NormalizedDiffAroonDown
			entires[i].NDiff_AroonUp_50 = rawData.NormalizedDiffAroonUp
			entires[i].NDiff_Macd_20_200_200 = rawData.NormalizedDiffMacd
			entires[i].NDiff_Macd_Hist_20_200_200 = rawData.NormalizedDiffMacdHist
			entires[i].NDiff_Osc_10 = rawData.NormalizedDiffOSC
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

func SetNDiffEma(rawData []dbservice.RawDataEntity, period int, diffLength int) {

	var ema []*float64
	for _, v := range rawData {
		if period == 20 && v.EMA_20 != nil {
			ema = append(ema, v.EMA_20)
		} else if period == 50 && v.EMA_50 != nil {
			ema = append(ema, v.EMA_50)
		} else if period == 100 && v.EMA_100 != nil {
			ema = append(ema, v.EMA_100)
		} else if period == 200 && v.EMA_200 != nil {
			ema = append(ema, v.EMA_200)
		}
	}

	for i := 0; i+diffLength-1 < len(ema); i++ {

		avgdiffn := utils.AvgDiffNormalized(ema, i, i+diffLength-1)
		avgdiffn = math.Round(avgdiffn*10000) / 10000
		if period == 20 {
			rawData[i].NormalizedDiffNEMA20 = &avgdiffn
		} else if period == 50 {
			rawData[i].NormalizedDiffNEMA50 = &avgdiffn
		} else if period == 100 {
			rawData[i].NormalizedDiffNEMA100 = &avgdiffn
		} else if period == 200 {
			rawData[i].NormalizedDiffNEMA200 = &avgdiffn
		}

	}
}

func SetNDiffMacd(rawData []dbservice.RawDataEntity, diffLength int) {

	var macd []*float64
	for _, v := range rawData {
		if v.Macd_20_200_200 != nil {
			macd = append(macd, v.Macd_20_200_200)
		}
	}

	for i := 0; i+diffLength-1 < len(macd); i++ {
		avgdiffn := utils.AvgDiffNormalized(macd, i, i+diffLength-1)
		avgdiffn = math.Round(avgdiffn*10000) / 10000
		rawData[i].NormalizedDiffMacd = &avgdiffn
	}
}

func SetNDiffMacdHist(rawData []dbservice.RawDataEntity, diffLength int) {

	var hist []*float64
	for _, v := range rawData {
		if v.Macd_Hist_20_200_200 != nil {
			hist = append(hist, v.Macd_Hist_20_200_200)
		}
	}

	for i := 0; i+diffLength-1 < len(hist); i++ {
		avgdiffn := utils.AvgDiff(hist, i, i+diffLength-1)
		avgdiffn = math.Round(avgdiffn*10000) / 10000
		rawData[i].NormalizedDiffMacdHist = &avgdiffn
	}
}

func SetNDiffCCI(rawData []dbservice.RawDataEntity, diffLength int) {

	var ncci []*float64
	for _, v := range rawData {
		if v.CCI_100 != nil {
			ncci = append(ncci, v.GetNormalizedCCI())
		}
	}

	for i := 0; i+diffLength-1 < len(ncci); i++ {
		avgdiffn := utils.AvgNDiff(ncci, i, i+diffLength-1)
		avgdiffn = math.Round(avgdiffn*10000) / 10000
		rawData[i].NormalizedDiffNCCI100 = &avgdiffn
	}
}

func SetNDiffMacdOSC(rawData []dbservice.RawDataEntity, diffLength int) {

	var osc []*float64
	for _, v := range rawData {
		if v.GetNormalizedOSC() != nil {
			d := v.GetNormalizedOSC()
			osc = append(osc, d)
		}
	}

	for i := 0; i+diffLength-1 < len(osc); i++ {
		avgdiffn := utils.AvgDiff(osc, i, i+diffLength-1)
		avgdiffn = math.Round(avgdiffn*10000) / 10000
		rawData[i].NormalizedDiffOSC = &avgdiffn
	}
}

func SetNDiffAroon(rawData []dbservice.RawDataEntity, diffLength int) {

	var aroonUp []*float64
	var aroonDown []*float64
	for _, v := range rawData {
		if v.AroonUp_50 != nil && v.AroonDown_50 != nil {
			aroonUp = append(aroonUp, v.GetNormalizedAroonUp())
			aroonDown = append(aroonDown, v.GetNormalizedAroonDown())
		}
	}

	for i := 0; i+diffLength-1 < len(aroonUp); i++ {
		up := utils.AvgNDiff(aroonUp, i, i+diffLength-1)
		up = math.Round(up*10000) / 10000
		rawData[i].NormalizedDiffAroonUp = &up

		down := utils.AvgNDiff(aroonDown, i, i+diffLength-1)
		down = math.Round(down*10000) / 10000
		rawData[i].NormalizedDiffAroonDown = &down
	}
}
