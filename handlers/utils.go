package handlers

import (
	"bytes"
	"encoding/csv"
	"errors"
	"math"
	"strconv"
)

type DataReader func(symbol string, line []string) (*RawDataEntity, error)

func ReadCsvData(symbol string, data []byte, reader DataReader) ([]*RawDataEntity, error) {

	var resps []*RawDataEntity

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

	return resps, nil
}

func CandleReader(symbol string, line []string) (*RawDataEntity, error) {
	resp := &RawDataEntity{}
	resp.Date = line[0]
	resp.Symbol = symbol

	var err error
	var theError error

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

func EMA_8_Reader(symbol string, line []string) (*RawDataEntity, error) {
	var err error
	resp := &RawDataEntity{}
	resp.Date = line[0]
	resp.Symbol = symbol
	tmp, err := strconv.ParseFloat(line[1], 32)
	resp.EMA_8 = &tmp
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func SetChanges(dailyResps []*RawDataEntity) {
	for _, resp := range dailyResps {
		SetChange(resp)
	}
}

func ToInterfaceArray(data []*RawDataEntity) []interface{} {
	result := make([]interface{}, len(data))
	for i, s := range data {
		result[i] = s
	}
	return result
}
