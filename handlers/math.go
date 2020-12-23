package handlers

import (
	"math"
	"github.com/xiaofeiqiu/data-preprocessor/services/dbservice"
)

func SetChange(input *dbservice.RawDataEntity) {
	tmp := (input.Close - input.Open) * 100 / input.Open
	input.Change = math.Round(tmp*100) / 100
}

func Normalize(x float64, low float64, high float64) float64 {
	tmp := (x - low) / (high - low)
	return math.Round(tmp*100) / 100
}
