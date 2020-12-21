package handlers

import (
	"errors"
	"github.com/xiaofeiqiu/data-preprocessor/lib/restutils"
	"net/http"
)

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

func (api *ApiHandler) GetDailyAdjusted(w http.ResponseWriter, r *http.Request) (int, error) {

	status, body, err := api.AlphaVantageApi.GetDailyAdjusted(r)
	if err != nil {
		return status, errors.New("error calling GetDailyAdjusted, " + err.Error())
	}

	if restutils.Is2xxStatusCode(status) {
		resp, err := ToDailyResponseArray(body)
		if err != nil {
			return 500, errors.New("error reading response, " + err.Error())
		}
		SetStats(resp)
		restutils.ResponseWithJson(w, 200, resp)
	}



	return 500, errors.New("unexpected error occurred")
}
