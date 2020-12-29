package dbservice

import (
	"errors"
	"github.com/xiaofeiqiu/data-preprocessor/services/utils"
	"math"
	"sort"
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

	maxEma                *float64 `json:"-"`
	minEma                *float64 `json:"-"`
	NormalizedDiffNEMA20  *float64 `json:"-"`
	NormalizedDiffNEMA50  *float64 `json:"-"`
	NormalizedDiffNEMA100 *float64 `json:"-"`
	NormalizedDiffNEMA200 *float64 `json:"-"`
}

func (e *RawDataEntity) GetNormalizedEMA(period int) *float64 {

	err := e.loadMaxMixEma()
	if err != nil {
		return nil
	}

	switch period {
	case 20:
		return utils.Normalize(*e.EMA_20, *e.minEma, *e.maxEma, 4)
	case 50:
		return utils.Normalize(*e.EMA_50, *e.minEma, *e.maxEma, 4)
	case 100:
		return utils.Normalize(*e.EMA_100, *e.minEma, *e.maxEma, 4)
	case 200:
		return utils.Normalize(*e.EMA_200, *e.minEma, *e.maxEma, 4)
	default:
		return nil
	}
}

func (e *RawDataEntity) GetNormalizedMacd() *float64 {
	if e.Macd_20_200_200 == nil || e.Macd_Signal_20_200_200 == nil {
		return nil
	}

	n := *e.Macd_20_200_200 / *e.Macd_Signal_20_200_200
	n = math.Round(n*10000) / 10000
	return &n
}

func (e *RawDataEntity) GetNormalizedAroonUp() *float64 {
	if e.AroonUp_50 == nil {
		return nil
	}
	n := *e.AroonUp_50 / float64(100)
	n = math.Round(n*10000) / 10000
	return &n
}

func (e *RawDataEntity) GetNormalizedAroonDown() *float64 {
	if e.AroonDown_50 == nil {
		return nil
	}
	n := *e.AroonDown_50 / float64(100)
	n = math.Round(n*10000) / 10000
	return &n
}

func (e *RawDataEntity) GetNormalizedCCI() *float64 {
	if e.CCI_100 == nil {
		return nil
	}
	n := *e.CCI_100 / float64(100)
	n = math.Round(n*10000) / 10000
	return &n
}

func (e *RawDataEntity) loadMaxMixEma() error {
	if e.isAnyEmaNull() {
		return errors.New("nil ema")
	}
	arr := []float64{*e.EMA_20, *e.EMA_50, *e.EMA_100, *e.EMA_200}
	sort.Float64s(arr)

	e.maxEma = &arr[3]
	e.minEma = &arr[0]
	return nil
}

func (e *RawDataEntity) isAnyEmaNull() bool {
	return e.EMA_20 == nil || e.EMA_50 == nil || e.EMA_100 == nil || e.EMA_200 == nil
}

func (e *RawDataEntity) isMaxMinEmaNil() bool {
	return e.minEma == nil || e.maxEma == nil
}
