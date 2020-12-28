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

type DataInput struct {
	Symbol string    `json:"symbol" db:"symbol, primarykey, size:6"`
	Date   time.Time `json:"date" db:"dt, primarykey"`

	N_EMA         *float64 `json:"n_ema,omitempty" db:"-"`
	N_EMA_20      *float64 `json:"n_ema20,omitempty" db:"n_ema20"`
	N_EMA_50      *float64 `json:"n_ema50,omitempty" db:"n_ema50"`
	N_EMA_100     *float64 `json:"n_ema100,omitempty" db:"n_ema100"`
	N_EMA_200     *float64 `json:"n_ema200,omitempty" db:"n_ema200"`
	NDiff_EMA     *float64 `json:"ndiff_ema,omitempty" db:"-"`
	NDiff_EMA_20  *float64 `json:"ndiff_ema20,omitempty" db:"ndiff_ema20"`
	NDiff_EMA_50  *float64 `json:"ndiff_ema50,omitempty" db:"ndiff_ema50"`
	NDiff_EMA_100 *float64 `json:"ndiff_ema100,omitempty" db:"ndiff_ema100"`
	NDiff_EMA_200 *float64 `json:"ndiff_ema200,omitempty" db:"ndiff_ema200"`

	N_CCI         *float64 `json:"n_cci,omitempty" db:"-"`
	N_CCI_100     *float64 `json:"n_cci100,omitempty" db:"n_cci100"`
	NDiff_CCI     *float64 `json:"n_diff_cci,omitempty" db:"-"`
	NDiff_CCI_100 *float64 `json:"n_diff_cci100,omitempty" db:"ndiff_cci100"`

	N_AroonUp          *float64 `json:"n_aroonup,omitempty" db:"-"`
	N_AroonDown        *float64 `json:"n_aroondown,omitempty" db:"-"`
	N_AroonUp_50       *float64 `json:"n_aroonup50,omitempty" db:"n_aroonup50"`
	N_AroonDown_50     *float64 `json:"n_aroondown50,omitempty" db:"n_aroondown50"`
	NDiff_AroonUp      *float64 `json:"ndiff_aroonup,omitempty" db:"-"`
	NDiff_AroonDown    *float64 `json:"ndiff_aroondown,omitempty" db:"-"`
	NDiff_AroonUp_50   *float64 `json:"ndiff_aroonup50,omitempty" db:"ndiff_aroonup50"`
	NDiff_AroonDown_50 *float64 `json:"ndiff_aroondown50,omitempty" db:"ndiff_aroondown50"`

	N_Macd                       *float64 `json:"nmacd,omitempty" db:"-"`
	N_MacdHist                   *float64 `json:"nmacdhist,omitempty" db:"-"`
	N_MacdSignal                 *float64 `json:"nmacdsignal,omitempty" db:"-"`
	N_Macd_20_200_200            *float64 `json:"nmacd20200200,omitempty" db:"n_macd_20200200"`
	N_Macd_Signal_20_200_200     *float64 `json:"nmacd_signal20200200,omitempty" db:"n_macd_signal_20200200"`
	N_Macd_Hist_20_200_200       *float64 `json:"nmacd_hist20200200,omitempty" db:"n_macd_hist_20200200"`
	NDiff_Macd                   *float64 `json:"ndiff_macd,omitempty" db:"-"`
	NDiff_MacdHist               *float64 `json:"ndiff_macdhist,omitempty" db:"-"`
	NDiff_MacdSignal             *float64 `json:"ndiff_macdsignal,omitempty" db:"-"`
	NDiff_Macd_20_200_200        *float64 `json:"ndiff_macd20200200,omitempty" db:"ndiff_macd_20200200"`
	NDiff_Macd_Signal_20_200_200 *float64 `json:"ndiff_macd_signal20200200,omitempty" db:"ndiff_macd_signal_20200200"`
	NDiff_Macd_Hist_20_200_200   *float64 `json:"ndiff_macd_hist20200200,omitempty" db:"ndiff_macd_hist_20200200"`

	N_Osc    *float64 `json:"n_osc,omitempty" db:"-"`
	N_Osc_10 *float64 `json:"n_osc10,omitempty" db:"n_osc10"`

	NDiff_Osc    *float64 `json:"osc,omitempty" db:"-"`
	NDiff_Osc_10 *float64 `json:"ndiff_osc10,omitempty" db:"ndiff_osc10"`
}
