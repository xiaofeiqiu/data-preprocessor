package handlers

import (
	"github.com/xiaofeiqiu/data-preprocessor/lib/restutils"
	"github.com/xiaofeiqiu/data-preprocessor/services/alphavantage"
	"gopkg.in/go-playground/validator.v9"
	"github.com/xiaofeiqiu/data-preprocessor/services/dbservice"
)

var validate = validator.New()

type ApiHandler struct {
	AlphaVantageClient *alphavantage.AlphaVantageClient
	DBService          *dbservice.DBService
	DefaultClient      *restutils.HttpClient
}

type DataInputRequest struct {
	Symbol     string `json:"symbol" validate:"required"`
	ColName    string `json:colName validate:"required"`
	DiffLength int    `json:diffLength`
}
