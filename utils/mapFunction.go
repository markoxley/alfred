package utils

import "math"

// Map remaps the value v from the range of mn1 to mx1
// to the range of mn2 to mx2
func Map(v, mn1, mx1, mn2, mx2 float64) float64 {
	df1 := math.Abs(mx1 - mn1)
	df2 := math.Abs(mx2 - mn2)
	nv := (((v - mn1) / df1) * df2) + mn2
	return math.Max(mn2, math.Min(mx2, nv))
}
