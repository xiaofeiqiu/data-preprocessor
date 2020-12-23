package handlers

import (
	"github.com/xiaofeiqiu/data-preprocessor/services/alphavantage"
	"gopkg.in/go-playground/validator.v9"
	"github.com/xiaofeiqiu/data-preprocessor/lib/db"
)

var validate = validator.New()

// table names
const dailyRawData = "daily_raw_data"

type RawDataEntity struct {
	Symbol        string   `json:"symbol" db:"symbol, primarykey"`
	Date          string   `json:"date" db:"dt, primarykey"`
	Open          float64  `json:"open" db:"open"`
	High          float64  `json:"high" db:"high"`
	Low           float64  `json:"low" db:"low"`
	Close         float64  `json:"close" db:"close"`
	AdjustedClose float64  `json:"adjusted_close" db:"adjusted_close"`
	Volume        int64    `json:"volume" db:"volume"`
	Change        float64  `json:"change" db:"change"`
	EMA_8         *float64 `json:"ema_8,omitempty" db:"ema_8"`
}

type ApiHandler struct {
	AlphaVantageClient *alphavantage.AlphaVantageClient
	DBClient           *db.DBClient
}

func (api *ApiHandler) InitDBTableMapping() {
	api.DBClient.DB.AddTableWithName(RawDataEntity{}, dailyRawData)
	api.DBClient.DB.CreateTablesIfNotExists()
}
