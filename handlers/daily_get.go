package handlers

import (
	"errors"
	"github.com/xiaofeiqiu/data-preprocessor/lib/restutils"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

var validate = validator.New()

const TIME_SERIES_DAILY_ADJUSTED = "TIME_SERIES_DAILY_ADJUSTED"

type DailyResponse struct {
	Timestamp     string  `json:timestamp`
	Open          float64 `json:open`
	High          float64 `json:high`
	Low           float64 `json:low`
	Close         float64 `json:close`
	AdjustedClose float64 `json:adjusted_close`
	Volume        int64   `json:volume`
	Change        float64 `json:change`
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
		resp, err := ToDailyResponseArray(body)
		if err != nil {
			return 500, errors.New("error reading response, " + err.Error())
		}
		restutils.ResponseWithJson(w, 200, resp)
		return 0, nil
	}

	return 500, errors.New("unexpected error occurred")
}
