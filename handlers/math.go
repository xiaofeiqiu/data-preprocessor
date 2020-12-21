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
	tmp := (input.Close - input.Low) / (input.High - input.Low)
	input.N_Close = math.Round(tmp*100) / 100
}
