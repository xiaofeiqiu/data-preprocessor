package handlers

import (
	"github.com/xiaofeiqiu/data-preprocessor/services"
	"math"
)

func SetChange(input *services.DailyResponse) {
	tmp := (input.Close - input.Open) * 100 / input.Open
	input.Change = math.Round(tmp*100) / 100
}

func SetNClose(input *services.DailyResponse) {
	tmp := (input.Close - input.Low) / (input.High - input.Low)
	input.N_Close = math.Round(tmp*100) / 100
}
