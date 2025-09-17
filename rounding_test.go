package decimal

import (
	"fmt"
	"math/big"
	"testing"
)

// TestRoundingMode_shouldRoundUp tests the shouldRoundUp logic for all rounding modes.
func TestRoundingMode_shouldRoundUp(t *testing.T) {
	i64 := func(i int64) *big.Int {
		return big.NewInt(i)
	}
	testCases := []struct {
		name       string
		mode       RoundingMode
		isNegative bool
		rem        *big.Int
		denom      *big.Int
		expected   bool
	}{
		// RoundDown
		{"RoundDown_Positive", RoundDown, false, i64(3), i64(10), false},
		{"RoundDown_Negative", RoundDown, true, i64(3), i64(10), false},

		// RoundUp
		{"RoundUp_Positive", RoundUp, false, i64(3), i64(10), true},
		{"RoundUp_Negative", RoundUp, true, i64(3), i64(10), true},

		// RoundCeiling
		{"RoundCeiling_Positive", RoundCeiling, false, i64(3), i64(10), true},
		{"RoundCeiling_Negative", RoundCeiling, true, i64(3), i64(10), false},

		// RoundFloor
		{"RoundFloor_Positive", RoundFloor, false, i64(3), i64(10), false},
		{"RoundFloor_Negative", RoundFloor, true, i64(3), i64(10), true},

		// RoundHalfUp
		{"RoundHalfUp_Positive_LessThanHalf", RoundHalfUp, false, i64(4), i64(10), false},
		{"RoundHalfUp_Positive_ExactlyHalf", RoundHalfUp, false, i64(5), i64(10), true},
		{"RoundHalfUp_Positive_MoreThanHalf", RoundHalfUp, false, i64(6), i64(10), true},
		{"RoundHalfUp_Negative_LessThanHalf", RoundHalfUp, true, i64(4), i64(10), false},
		{"RoundHalfUp_Negative_ExactlyHalf", RoundHalfUp, true, i64(5), i64(10), true},
		{"RoundHalfUp_Negative_MoreThanHalf", RoundHalfUp, true, i64(6), i64(10), true},

		// RoundHalfDown
		{"RoundHalfDown_Positive_LessThanHalf", RoundHalfDown, false, i64(4), i64(10), false},
		{"RoundHalfDown_Positive_ExactlyHalf", RoundHalfDown, false, i64(5), i64(10), false},
		{"RoundHalfDown_Positive_MoreThanHalf", RoundHalfDown, false, i64(6), i64(10), true},
		{"RoundHalfDown_Negative_LessThanHalf", RoundHalfDown, true, i64(4), i64(10), false},
		{"RoundHalfDown_Negative_ExactlyHalf", RoundHalfDown, true, i64(5), i64(10), false},
		{"RoundHalfDown_Negative_MoreThanHalf", RoundHalfDown, true, i64(6), i64(10), true},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.name), func(t *testing.T) {
			result := tc.mode.shouldRoundUp(tc.isNegative, tc.rem, tc.denom)
			if result != tc.expected {
				t.Errorf("shouldRoundUp(isNegative: %t, rem: %v, denom: %v) with mode %s = %t; want %t",
					tc.isNegative, tc.rem, tc.denom, tc.mode, result, tc.expected)
			}
		})
	}
}

func TestRoundingMode_String(t *testing.T) {
	tests := []struct {
		mode     RoundingMode
		expected string
	}{
		{RoundDown, "RoundDown"},
		{RoundUp, "RoundUp"},
		{RoundCeiling, "RoundCeiling"},
		{RoundFloor, "RoundFloor"},
		{RoundHalfUp, "RoundHalfUp"},
		{RoundHalfDown, "RoundHalfDown"},
		{RoundHalfEven, "RoundHalfEven"},
		{RoundUnnecessary, "RoundUnnecessary"},
		{RoundingMode(999), "RoundingMode(999)"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.mode.String() != tt.expected {
				t.Errorf("RoundingMode.String() = %v; want %v", tt.mode.String(), tt.expected)
			}
		})
	}
}
