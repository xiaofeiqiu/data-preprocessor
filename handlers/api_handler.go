package handlers

import (
	"github.com/xiaofeiqiu/data-preprocessor/lib/db"
	"github.com/xiaofeiqiu/data-preprocessor/services/alphavantage"
	"gopkg.in/go-playground/validator.v9"
	"github.com/xiaofeiqiu/data-preprocessor/services/dbservice"
)

var validate = validator.New()

// table names
const dailyRawData = "daily_raw_data"

type ApiHandler struct {
	AlphaVantageClient *alphavantage.AlphaVantageClient
	DBClient           *db.DBClient
}

func (api *ApiHandler) InitDBTableMapping() error {
	api.DBClient.DB.AddTableWithName(dbservice.RawDataEntity{}, dailyRawData)
	err := api.DBClient.DB.CreateTablesIfNotExists()
	if err != nil {
		return err
	}
	return nil
}
