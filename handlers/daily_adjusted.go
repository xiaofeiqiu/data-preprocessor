package handlers

import (
	"bytes"
	"encoding/csv"
	"errors"
	"github.com/xiaofeiqiu/mlstock/lib/restutils"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
)

var validate = validator.New()

const TIME_SERIES_DAILY_ADJUSTED = "TIME_SERIES_DAILY_ADJUSTED"

type DailyResponse struct {
	Timestamp     string  `json:timestamp`
	Open          float64 `json:open`
	High          float64 `json:High`
	Low           float64 `json:Low`
	Close         float64 `json:Close`
	AdjustedClose float64 `json:adjusted_close`
	Volume        int64   `json:Volume`
}

type DailyResponseArray struct {
	Data []DailyResponse `data`
}

func (h *ApiHandler) GetDailyAdjusted(w http.ResponseWriter, r *http.Request) (int, error) {
	req := MLStockRequest{
		Function: TIME_SERIES_DAILY_ADJUSTED,
		DataType: DataTypeCsv,
	}

	err := decoder.Decode(&req, r.URL.Query())
	if err != nil {
		return 400, errors.New("error decoding query params, " + err.Error())
	}

	if req.OutputSize == "" {
		req.OutputSize = OutputCompact
	}

	err = validate.Struct(req)
	if err != nil {
		return 400, errors.New("request validation failed, " + err.Error())
	}

	url := h.GetUrl(req)
	status, body, err := h.HttpClient.DoGet(url, nil)
	if err != nil {
		return 500, errors.New("error calling downstream api, " + err.Error())
	}

	if restutils.Is2xxStatusCode(status) {
		resp, err := readDailyResponseCSV(body)
		if err != nil {
			return 500, errors.New("error reading response, " + err.Error())
		}
		restutils.ResponseWithJson(w, 200, resp)
		return 0, nil
	}

	return 500, errors.New("unexpected error occurred")
}

func readDailyResponseCSV(data []byte) ([]DailyResponse, error) {

	resps := []DailyResponse{}

	r := csv.NewReader(bytes.NewReader(data))
	lines, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	lines = lines[1:]

	for _, line := range lines {
		var resp DailyResponse
		resp.Timestamp = line[0]
		resp.Open, _ = strconv.ParseFloat(line[1], 32)
		resp.High, _ = strconv.ParseFloat(line[2], 32)
		resp.Low, _ = strconv.ParseFloat(line[3], 32)
		resp.Close, _ = strconv.ParseFloat(line[4], 32)
		resp.AdjustedClose, _ = strconv.ParseFloat(line[5], 32)
		resp.Volume, _ = strconv.ParseInt(line[6], 10, 64)
		resps = append(resps, resp)
	}
	return resps, nil
}
