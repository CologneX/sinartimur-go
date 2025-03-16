package utils

import "math"

// FloatEquals checks if two float64 values are equal within a small epsilon
func FloatEquals(a, b float64) bool {
	epsilon := 0.000001 // Small threshold for floating point comparison
	return math.Abs(a-b) < epsilon
}
