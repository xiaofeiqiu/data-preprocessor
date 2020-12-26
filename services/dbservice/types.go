package dbservice

import (
	"encoding/json"
	"time"
)

type RawDataEntity struct {
	Symbol        string    `json:"symbol" db:"symbol, primarykey, size:6"`
	Date          time.Time `json:"date" db:"dt, primarykey"`
	Open          float64   `json:"open" db:"open"`
	High          float64   `json:"high" db:"high"`
	Low           float64   `json:"low" db:"low"`
	Close         float64   `json:"close" db:"close"`
	AdjustedClose float64   `json:"adjusted_close" db:"adjusted_close"`
	Volume        int64     `json:"volume" db:"volume"`
	Change        float64   `json:"change" db:"change"`
	EMA           *float64  `json:"ema,omitempty" db:"-"`
	EMA_20        *float64  `json:"ema20,omitempty" db:"ema20"`
	EMA_50        *float64  `json:"ema50,omitempty" db:"ema50"`
	EMA_100       *float64  `json:"ema100,omitempty" db:"ema100"`
	EMA_200       *float64  `json:"ema200,omitempty" db:"ema200"`
}

func (e *RawDataEntity) ToString() string {
	str, _ := json.Marshal(e)
	return string(str)
}
