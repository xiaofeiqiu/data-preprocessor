package utils

import (
	"math"
)

func Normalize(x float64, min float64, max float64, around float64) *float64 {
	tmp := (x - min) / (max - min)
	result := math.Round(tmp*math.Pow(10, around)) / math.Pow(10, around)
	return &result
}

func AvgDiffNormalized(data []*float64, from, to int) float64 {
	return (*data[from] - *data[to]) * 100 / *data[to] / float64(to-from+1)
}
