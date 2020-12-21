package handlers

import (
	"github.com/xiaofeiqiu/data-preprocessor/services/alphavantage"
	"gopkg.in/go-playground/validator.v9"
)

var validate = validator.New()

type ApiHandler struct {
	AlphaVantageApi *alphavantage.AlphaVantageApi
}




