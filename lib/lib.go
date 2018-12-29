package lib

import (
	"math"
	"strconv"
	"strings"
)

func Rad(d float64) float64 {
	pi := math.Pi / 180
	return d * pi
}

func StringToFloatArray(s string) ([]float64, error) {

	var res []float64
	val := strings.Split(s, ",")

	for idx := range val {
		out, err := strconv.ParseFloat(val[idx], 64)
		if err != nil {
			return res, err
		}
		res = append(res, out)
	}

	return res, nil
}
