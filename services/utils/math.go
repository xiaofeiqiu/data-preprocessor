package utils

import "math"

func Normalize(x float64, min float64, max float64) *float64 {
	tmp := (x - min) / (max - min)
	result := math.Round(tmp*100) / 100
	return &result
}
