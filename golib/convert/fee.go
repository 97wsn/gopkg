package convert

import "math"

func ToFen(fee float64) int {
	return int(math.Round(fee * 100))
}
