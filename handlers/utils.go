package handlers

import (
	"bytes"
	"encoding/csv"
	"errors"
	"math"
	"strconv"
)

func ToDailyResponseArray(data []byte) ([]*DailyResponse, error) {

	var resps []*DailyResponse

	r := csv.NewReader(bytes.NewReader(data))
	lines, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	lines = lines[1:]

	for _, line := range lines {
		resp, err := readCandles(line)
		if err != nil {
			return nil, errors.New("error reading values, " + err.Error())
		}
		resps = append(resps, resp)
	}
	return resps, nil
}

func readCandles(line []string) (*DailyResponse, error) {
	resp := &DailyResponse{}
	resp.Timestamp = line[0]
	var err error
	resp.Open, err = strconv.ParseFloat(line[1], 32)
	resp.High, err = strconv.ParseFloat(line[2], 32)
	resp.Low, err = strconv.ParseFloat(line[3], 32)
	resp.Close, err = strconv.ParseFloat(line[4], 32)
	resp.AdjustedClose, err = strconv.ParseFloat(line[5], 32)
	resp.Volume, err = strconv.ParseInt(line[6], 10, 64)

	resp.Open = math.Round(resp.Open*100) / 100
	resp.High = math.Round(resp.High*100) / 100
	resp.Low = math.Round(resp.Low*100) / 100
	resp.Close = math.Round(resp.Close*100) / 100
	resp.AdjustedClose = math.Round(resp.AdjustedClose*100) / 100

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func SetStats(dailyResps []*DailyResponse) {
	for _, resp := range dailyResps {
		SetChange(resp)
	}
}

func SetChange(input *DailyResponse) {
	tmp := (input.Close - input.Open) * 100 / input.Open
	input.Change = math.Round(tmp*100) / 100
}
