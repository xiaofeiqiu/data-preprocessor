package utils

import (
	"math"
)

func Normalize(x float64, min float64, max float64) *float64 {
	tmp := (x - min) / (max - min)
	result := math.Round(tmp*100) / 100
	return &result
}

func AvgDiffNormalized(data []*float64, from, to int) float64 {
	return (*data[from] - *data[to]) * 100 / *data[to] / float64(to-from+1)
}
