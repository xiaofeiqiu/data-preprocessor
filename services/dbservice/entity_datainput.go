package dbservice

import "time"

type DataInputEntity struct {
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
	NDiff_Macd                   *float64 `json:"ndiff_macd,omitempty" db:"-"`
	NDiff_MacdHist               *float64 `json:"ndiff_macdhist,omitempty" db:"-"`
	NDiff_MacdSignal             *float64 `json:"ndiff_macdsignal,omitempty" db:"-"`
	NDiff_Macd_20_200_200        *float64 `json:"ndiff_macd20200200,omitempty" db:"ndiff_macd_20200200"`
	NDiff_Macd_Hist_20_200_200   *float64 `json:"ndiff_macd_hist20200200,omitempty" db:"ndiff_macd_hist_20200200"`

	N_Osc    *float64 `json:"n_osc,omitempty" db:"-"`
	N_Osc_10 *float64 `json:"n_osc10,omitempty" db:"n_osc10"`

	NDiff_Osc    *float64 `json:"osc,omitempty" db:"-"`
	NDiff_Osc_10 *float64 `json:"ndiff_osc10,omitempty" db:"ndiff_osc10"`

	Buy  *bool `json:"buy,omitempty" db:"buy"`
	Sell *bool `json:"sell,omitempty" db:"sell"`
	Hold *bool `json:"hold,omitempty" db:"hold"`
}
