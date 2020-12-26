package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/xiaofeiqiu/data-preprocessor/lib/log"
	"github.com/xiaofeiqiu/data-preprocessor/lib/restutils"
	"github.com/xiaofeiqiu/data-preprocessor/services/alphavantage"
	"io/ioutil"
	"net/http"
)

type DoallRequest struct {
	Symbol     string `json:"symbol" validate:"required"`
	OutputSize string `json:"outputsize"`
}

var c = &restutils.HttpClient{
	Client: &http.Client{},
}

func (api *ApiHandler) Doall(w http.ResponseWriter, r *http.Request) (int, error) {
	req, err := NewDoAllRequest(r)
	if err != nil {
		return 400, errors.New("error creating new do all request, " + err.Error())
	}

	// post daily missing
	url := "http://localhost:8080/preprocessor/candle/missingdailyadjusted"
	body := []byte(fmt.Sprintf(`{"symbol":"%s","outputsize": "%s"}`, req.Symbol, req.OutputSize))
	status, _, err := api.DefaultClient.DoPost(url, nil, body)
	if err != nil {
		return 500, errors.New("error calling data pre processor, post daily missing failed, " + err.Error())
	}

	if !restutils.Is2xxStatusCode(status) {
		return 500, errors.New("post daily missing failed")
	}
	log.Info("Doall", "Post daily successful")

	// fill ema 20
	period := "20"
	status, err = api.fillEma(req, period)
	if !restutils.Is2xxStatusCode(status) {
		return 500, err
	}
	log.Info("Doall", fmt.Sprintf("fill ema %s successful", period))

	// fill ema 50
	period = "50"
	status, err = api.fillEma(req, period)
	if !restutils.Is2xxStatusCode(status) {
		return 500, err
	}
	log.Info("Doall", fmt.Sprintf("fill ema %s successful", period))

	// fill ema 100
	period = "100"
	status, err = api.fillEma(req, period)
	if !restutils.Is2xxStatusCode(status) {
		return 500, err
	}
	log.Info("Doall", fmt.Sprintf("fill ema %s successful", period))

	// fill ema 200
	period = "200"
	status, err = api.fillEma(req, period)
	if !restutils.Is2xxStatusCode(status) {
		return 500, err
	}
	log.Info("Doall", fmt.Sprintf("fill ema %s successful", period))

	restutils.ResponseWithJson(w, 200, "successful")
	return 0, nil
}

func (api *ApiHandler) fillEma(req DoallRequest, period string) (int, error) {
	url := "http://localhost:8080/preprocessor/ema/dailyadjusted"
	body := []byte(fmt.Sprintf(`{"symbol":"%s","time_period":"%s"}`, req.Symbol, period))
	status, resp, err := api.DefaultClient.DoPut(url, nil, body)
	if err != nil {
		return 500, errors.New("error calling data pre processor, fill ema " + period + " failed, " + err.Error())
	}

	if !restutils.Is2xxStatusCode(status) {
		return 500, errors.New(fmt.Sprintf("fill ema %s failed, %s", period, string(resp)))
	}
	return 200, nil
}

func NewDoAllRequest(r *http.Request) (DoallRequest, error) {
	req := DoallRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return DoallRequest{}, err
	}
	json.Unmarshal(body, &req)

	if req.OutputSize == "" {
		req.OutputSize = alphavantage.Compact
	}

	return req, nil
}
