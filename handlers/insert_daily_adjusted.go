package handlers

import (
	"errors"
	"github.com/xiaofeiqiu/data-preprocessor/lib/restutils"
	"github.com/xiaofeiqiu/data-preprocessor/services/alphavantage"
	"net/http"
)

func (api *ApiHandler) InsertDailyAdjusted(w http.ResponseWriter, r *http.Request) (int, error) {

	status, body, err := api.AlphaVantageApi.Call(alphavantage.TIME_SERIES_DAILY_ADJUSTED, r)
	if err != nil {
		return status, errors.New("error calling TIME_SERIES_DAILY_ADJUSTED, " + err.Error())
	}

	resp := []*alphavantage.DailyResponse{}
	if restutils.Is2xxStatusCode(status) {
		resp, err = ToDailyResponseArray(body, CandleReader)
		if err != nil {
			return 500, errors.New("error reading response, " + err.Error())
		}
		SetStats(resp)
		restutils.ResponseWithJson(w, 200, resp)
		return 0, nil
	}

	return 500, errors.New("unexpected error occurred")
}
