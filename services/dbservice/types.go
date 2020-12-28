package dbservice

import (
	"time"
)

// make sure to update the date type in db when re-creating the table
type RawDataEntity struct {
	Symbol                 string    `json:"symbol" db:"symbol, primarykey, size:6"`
	Date                   time.Time `json:"date" db:"dt, primarykey"`
	Open                   float64   `json:"open" db:"open"`
	High                   float64   `json:"high" db:"high"`
	Low                    float64   `json:"low" db:"low"`
	Close                  float64   `json:"close" db:"close"`
	AdjustedClose          float64   `json:"adjusted_close" db:"adjusted_close"`
	Volume                 int64     `json:"volume" db:"volume"`
	Change                 float64   `json:"change" db:"change"`
	EMA                    *float64  `json:"ema,omitempty" db:"-"`
	EMA_20                 *float64  `json:"ema20,omitempty" db:"ema20"`
	EMA_50                 *float64  `json:"ema50,omitempty" db:"ema50"`
	EMA_100                *float64  `json:"ema100,omitempty" db:"ema100"`
	EMA_200                *float64  `json:"ema200,omitempty" db:"ema200"`
	CCI                    *float64  `json:"cci,omitempty" db:"-"`
	CCI_100                *float64  `json:"cci100,omitempty" db:"cci100"`
	AroonUp                *float64  `json:"aroonup,omitempty" db:"-"`
	AroonDown              *float64  `json:"aroondown,omitempty" db:"-"`
	AroonUp_50             *float64  `json:"aroonup50,omitempty" db:"aroonup50"`
	AroonDown_50           *float64  `json:"aroondown50,omitempty" db:"aroondown50"`
	Macd                   *float64  `json:"macd,omitempty" db:"-"`
	MacdHist               *float64  `json:"macdhist,omitempty" db:"-"`
	MacdSignal             *float64  `json:"macdsignal,omitempty" db:"-"`
	Macd_20_200_200        *float64  `json:"macd20200200,omitempty" db:"macd_20200200"`
	Macd_Signal_20_200_200 *float64  `json:"macd_signal20200200,omitempty" db:"macd_signal_20200200"`
	Macd_Hist_20_200_200   *float64  `json:"macd_hist20200200,omitempty" db:"macd_hist_20200200"`
	Osc                    *float64  `json:"osc,omitempty" db:"-"`
	Osc_10                 *float64  `json:"osc10,omitempty" db:"osc10"`
}
