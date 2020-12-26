package handlers

import (
	"bytes"
	"encoding/csv"
	"errors"
	"github.com/xiaofeiqiu/data-preprocessor/lib/log"
	"github.com/xiaofeiqiu/data-preprocessor/services/dbservice"
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

func SetChanges(dailyResps []*dbservice.RawDataEntity) {
	for _, resp := range dailyResps {
		SetChange(resp)
	}
	log.Info("SetChanges", "SetChanges successful")
}

func ToMap(data []*dbservice.RawDataEntity) map[string]*dbservice.RawDataEntity {
	result := map[string]*dbservice.RawDataEntity{}

	for _, v := range data {
		result[v.Date.Format(time.RFC3339)] = v
	}
	return result
}
