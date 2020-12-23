package handlers

import (
	"math"
)

func SetChange(input *RawDataEntity) {
	tmp := (input.Close - input.Open) * 100 / input.Open
	input.Change = math.Round(tmp*100) / 100
}

func Normalize(x float64, low float64, high float64) float64 {
	tmp := (x - low) / (high - low)
	return math.Round(tmp*100) / 100
}
