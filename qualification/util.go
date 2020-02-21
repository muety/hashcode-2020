package main

import "math"

func intMax(x, y int) int {
	return int(math.Max(float64(x), float64(y)))
}

func intMin(x, y int) int {
	return int(math.Min(float64(x), float64(y)))

}
