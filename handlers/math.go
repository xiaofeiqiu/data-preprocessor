package handlers

import (
	"github.com/xiaofeiqiu/data-preprocessor/services/alphavantage"
	"math"
)

func SetChange(input *alphavantage.DailyResponse) {
	tmp := (input.Close - input.Open) * 100 / input.Open
	input.Change = math.Round(tmp*100) / 100
}

func SetNClose(input *alphavantage.DailyResponse) {
	input.N_Close = Normalize(input.Close,input.Low,input.High)
}

func Normalize(x float64, low float64, high float64) float64 {
	tmp := (x - low) / (high - low)
	return math.Round(tmp*100) / 100
}
