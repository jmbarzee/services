package testhelper

import (
	"fmt"
	"math"

	"github.com/go-test/deep"
	"github.com/jmbarzee/services/lightorchestrator/service/color"
)

const (
	MinErrColor = 0.000001
)

// ShiftsEqual compares and diffs floats (shifts)
func ShiftsEqual(a, b float64, err float64) bool {
	return float64(math.Abs(float64(a-b))) < err
}

// ColorsEqual compares and diffs colors
func ColorsEqual(a, b color.HSLA) bool {
	if !ShiftsEqual(a.H, b.H, MinErrColor) {
		if a.H > 0.99 {
			if !ShiftsEqual(1-a.H, b.H, MinErrColor) {
				return false
			}
		} else if b.H > 0.99 {
			if !ShiftsEqual(a.H, 1-b.H, MinErrColor) {
				return false
			}
		} else {
			return false
		}
	}
	if !ShiftsEqual(a.S, b.S, MinErrColor) {
		return false
	}
	if !ShiftsEqual(a.L, b.L, MinErrColor) {
		return false
	}
	if !ShiftsEqual(a.A, b.A, MinErrColor) {
		return false
	}
	return true
}

// StructsEqual compares and diffs structs
func StructsEqual(expected, actual interface{}) bool {
	if diffs := deep.Equal(expected, actual); len(diffs) > 0 {
		for _, diff := range diffs {
			fmt.Printf("%s\n", diff)
		}
		return false
	}
	return true
}