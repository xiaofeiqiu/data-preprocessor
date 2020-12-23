package handlers

import (
	"github.com/xiaofeiqiu/data-preprocessor/services/alphavantage"
	"gopkg.in/go-playground/validator.v9"
	"github.com/xiaofeiqiu/data-preprocessor/lib/db"
)

var validate = validator.New()

type ApiHandler struct {
	AlphaVantageClient *alphavantage.AlphaVantageClient
	DBClient           *db.DBClient
}
