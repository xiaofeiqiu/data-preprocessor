package handlers

import (
	"bytes"
	"encoding/csv"
	"errors"
	"github.com/xiaofeiqiu/data-preprocessor/lib/log"
	"github.com/xiaofeiqiu/data-preprocessor/services/dbservice"
	"github.com/xiaofeiqiu/data-preprocessor/services/utils"
	"math"
	"strconv"
	"time"
)

type DataReader func(symbol string, line []string) (*dbservice.RawDataEntity, error)

func ReadCsvData(symbol string, data []byte, reader DataReader) ([]*dbservice.RawDataEntity, error) {

	log.Info("ReadCsvData", "Reading csv data")
	var resps []*dbservice.RawDataEntity

	r := csv.NewReader(bytes.NewReader(data))
	lines, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	lines = lines[1:]

	var theError error

	for _, line := range lines {
		resp, err := reader(symbol, line)
		if err != nil {
			theError = err
			continue
		}
		resps = append(resps, resp)
	}

	if theError != nil {
		return nil, errors.New("error reading values, " + err.Error())
	}

	log.Info("ReadCsvData", "Reading csv data successful")

	return resps, nil
}

func CandleReader(symbol string, line []string) (*dbservice.RawDataEntity, error) {
	var err error
	var theError error

	resp := &dbservice.RawDataEntity{}
	resp.Date, err = time.Parse("2006-01-02", line[0])
	if err != nil {
		theError = err
	}

	resp.Symbol = symbol

	resp.Open, err = strconv.ParseFloat(line[1], 32)
	if err != nil {
		theError = err
	}
	resp.High, err = strconv.ParseFloat(line[2], 32)
	if err != nil {
		theError = err
	}
	resp.Low, err = strconv.ParseFloat(line[3], 32)
	if err != nil {
		theError = err
	}
	resp.Close, err = strconv.ParseFloat(line[4], 32)
	if err != nil {
		theError = err
	}
	resp.AdjustedClose, err = strconv.ParseFloat(line[5], 32)
	if err != nil {
		theError = err
	}
	resp.Volume, err = strconv.ParseInt(line[6], 10, 64)
	if err != nil {
		theError = err
	}

	resp.Open = math.Round(resp.Open*100) / 100
	resp.High = math.Round(resp.High*100) / 100
	resp.Low = math.Round(resp.Low*100) / 100
	resp.Close = math.Round(resp.Close*100) / 100
	resp.AdjustedClose = math.Round(resp.AdjustedClose*100) / 100

	if theError != nil {
		return nil, err
	}

	return resp, nil
}

func EMA_Reader(symbol string, line []string) (*dbservice.RawDataEntity, error) {
	var err error
	resp := &dbservice.RawDataEntity{}
	resp.Date, err = time.Parse("2006-01-02", line[0])
	if err != nil {
		return nil, err
	}
	resp.Symbol = symbol
	tmp, err := strconv.ParseFloat(line[1], 32)
	tmp = math.Round(tmp*100) / 100
	resp.EMA = &tmp
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func CCI_Reader(symbol string, line []string) (*dbservice.RawDataEntity, error) {
	var err error
	resp := &dbservice.RawDataEntity{}
	resp.Date, err = time.Parse("2006-01-02", line[0])
	if err != nil {
		return nil, err
	}
	resp.Symbol = symbol
	tmp, err := strconv.ParseFloat(line[1], 32)
	tmp = math.Round(tmp*100) / 100
	resp.CCI = &tmp
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func Aroon_Reader(symbol string, line []string) (*dbservice.RawDataEntity, error) {
	var err error
	resp := &dbservice.RawDataEntity{}
	resp.Date, err = time.Parse("2006-01-02", line[0])
	if err != nil {
		return nil, err
	}
	resp.Symbol = symbol

	tmp1, err := strconv.ParseFloat(line[1], 32)
	tmp1 = math.Round(tmp1*100) / 100
	resp.AroonDown = &tmp1

	tmp2, err := strconv.ParseFloat(line[2], 32)
	tmp2 = math.Round(tmp2*100) / 100
	resp.AroonUp = &tmp2

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func MACD_Reader(symbol string, line []string) (*dbservice.RawDataEntity, error) {
	var err error
	resp := &dbservice.RawDataEntity{}
	resp.Date, err = time.Parse("2006-01-02", line[0])
	if err != nil {
		return nil, err
	}
	resp.Symbol = symbol

	tmp1, err := strconv.ParseFloat(line[1], 32)
	tmp1 = math.Round(tmp1*100) / 100
	resp.Macd = &tmp1

	tmp2, err := strconv.ParseFloat(line[2], 32)
	tmp2 = math.Round(tmp2*100) / 100
	resp.MacdHist = &tmp2

	tmp3, err := strconv.ParseFloat(line[3], 32)
	tmp3 = math.Round(tmp3*100) / 100
	resp.MacdSignal = &tmp3

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func Osc_Reader(symbol string, line []string) (*dbservice.RawDataEntity, error) {
	var err error
	resp := &dbservice.RawDataEntity{}
	resp.Date, err = time.Parse("2006-01-02", line[0])
	if err != nil {
		return nil, err
	}
	resp.Symbol = symbol

	tmp1, err := strconv.ParseFloat(line[1], 32)
	tmp1 = math.Round(tmp1*100) / 100
	resp.Osc = &tmp1

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func SetChanges(dailyResps []*dbservice.RawDataEntity) {
	for _, resp := range dailyResps {
		SetChange(resp)
	}
	log.Info("SetChanges", "SetChanges successful")
}

func SetChange(input *dbservice.RawDataEntity) {
	tmp := (input.Close - input.Open) * 100 / input.Open
	input.Change = math.Round(tmp*100) / 100
}

func RawDataPtrArrayToMap(data []*dbservice.RawDataEntity) map[string]*dbservice.RawDataEntity {
	result := map[string]*dbservice.RawDataEntity{}

	for _, v := range data {
		result[v.Date.Format(time.RFC3339)] = v
	}
	return result
}

func RawDataArrayToMap(data []dbservice.RawDataEntity) map[string]*dbservice.RawDataEntity {
	result := map[string]*dbservice.RawDataEntity{}

	for i, v := range data {
		result[v.Date.Format(time.RFC3339)] = &data[i]
	}
	return result
}

func validatePeriod(period string, validPeriods []string) bool {
	for _, v := range validPeriods {
		if period == v {
			return true
		}
	}
	return false
}

func LoadNormalizedNDiffEma(rawData []dbservice.RawDataEntity, period int, diffLength int) {

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
