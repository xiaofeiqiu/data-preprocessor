package handlers

import (
	"errors"
	"github.com/gorilla/schema"
	"github.com/xiaofeiqiu/mlstock/lib/restutils"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

var validate = validator.New()
var decoder = schema.NewDecoder()

const TIME_SERIES_DAILY_ADJUSTED = "TIME_SERIES_DAILY_ADJUSTED"

type DailyAdjustedRequest struct {
	Symbol     string `validate:"required" schema:"symbol"`
	OutputSize string `schema:"outputsize"`
}

func (h *ApiHandler) GetDailyAdjusted(w http.ResponseWriter, r *http.Request) (int, error) {
	var req DailyAdjustedRequest

	err := decoder.Decode(&req, r.URL.Query())
	if err != nil {
		return 400, errors.New("error decoding query params, " + err.Error())
	}

	err = validate.Struct(req)
	if err != nil {
		return 400, errors.New("request validation failed, " + err.Error())
	}

	restutils.ResponseWithJson(w, 200, "test")
	return 0, nil
}
