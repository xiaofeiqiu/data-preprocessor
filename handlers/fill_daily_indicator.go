package handlers

import (
	"errors"
	"github.com/xiaofeiqiu/data-preprocessor/lib/restutils"
	"github.com/xiaofeiqiu/data-preprocessor/services/alphavantage"
	"net/http"
)

func (api *ApiHandler) FillDailyIndicator(w http.ResponseWriter, r *http.Request) (int, error) {
	status, body, err := api.AlphaVantageApi.Call(alphavantage.EMA, r)
	if err != nil {
		return status, errors.New("error calling EMA, " + err.Error())
	}

	if restutils.Is2xxStatusCode(status) {
		resp, err := ToDailyResponseArray(body, EMA_8_Reader)
		if err != nil {
			return 500, errors.New("error reading response, " + err.Error())
		}

		restutils.ResponseWithJson(w, 200, resp)
		return 0, nil
	}

	return 500, errors.New("unexpected error occurred")
}
