package lib

import "math"

func Rad(d float64) float64 {
	pi := math.Pi / 180
	return d * pi
}
