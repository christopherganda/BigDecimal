package decimal

import (
	"fmt"
	"math/big"
	"testing"
)

// TestRoundingMode_shouldRoundUp tests the shouldRoundUp logic for all rounding modes.
func TestRoundingMode_shouldRoundUp(t *testing.T) {
	// A helper function to create a big.Int from an int64 for cleaner test cases.
	i64 := func(i int64) *big.Int {
		return big.NewInt(i)
	}

	testCases := []struct {
		name     string
		mode     RoundingMode
		rem      *big.Int
		denom    *big.Int
		expected bool
	}{
		// RoundDown: Always returns false unless the remainder is non-zero, in which case it is an error
		{"RoundDown_NoRounding", RoundDown, i64(0), i64(10), false},
		{"RoundDown_PositiveRemainder", RoundDown, i64(3), i64(10), false},
		{"RoundDown_NegativeRemainder", RoundDown, i64(-3), i64(10), false},

		// RoundUp: Always returns true
		{"RoundUp_PositiveRemainder", RoundUp, i64(3), i64(10), true},
		{"RoundUp_NegativeRemainder", RoundUp, i64(-3), i64(10), true},
		{"RoundUp_NoRemainder", RoundUp, i64(0), i64(10), false},

		// RoundCeiling: Rounds up toward positive infinity
		{"RoundCeiling_PositiveRemainder", RoundCeiling, i64(3), i64(10), true},
		{"RoundCeiling_NegativeRemainder", RoundCeiling, i64(-3), i64(10), false},
		{"RoundCeiling_NoRemainder", RoundCeiling, i64(0), i64(10), false},

		// RoundFloor: Rounds up toward negative infinity
		{"RoundFloor_PositiveRemainder", RoundFloor, i64(3), i64(10), false},
		{"RoundFloor_NegativeRemainder", RoundFloor, i64(-3), i64(10), true},
		{"RoundFloor_NoRemainder", RoundFloor, i64(0), i64(10), false},

		// RoundHalfUp: Rounds toward nearest neighbor, ties go up (positive infinity)
		{"RoundHalfUp_LessThanHalf", RoundHalfUp, i64(4), i64(10), false},
		{"RoundHalfUp_ExactlyHalf", RoundHalfUp, i64(5), i64(10), true},
		{"RoundHalfUp_MoreThanHalf", RoundHalfUp, i64(6), i64(10), true},
		{"RoundHalfUp_NegativeLessThanHalf", RoundHalfUp, i64(-4), i64(10), false},
		{"RoundHalfUp_NegativeExactlyHalf", RoundHalfUp, i64(-5), i64(10), true},
		{"RoundHalfUp_NegativeMoreThanHalf", RoundHalfUp, i64(-6), i64(10), true},

		// RoundHalfDown: Rounds toward nearest neighbor, ties go down (negative infinity)
		{"RoundHalfDown_LessThanHalf", RoundHalfDown, i64(4), i64(10), false},
		{"RoundHalfDown_ExactlyHalf", RoundHalfDown, i64(5), i64(10), false},
		{"RoundHalfDown_MoreThanHalf", RoundHalfDown, i64(6), i64(10), true},
		{"RoundHalfDown_NegativeLessThanHalf", RoundHalfDown, i64(-4), i64(10), false},
		{"RoundHalfDown_NegativeExactlyHalf", RoundHalfDown, i64(-5), i64(10), false},
		{"RoundHalfDown_NegativeMoreThanHalf", RoundHalfDown, i64(-6), i64(10), true},

		// RoundHalfEven: Rounds toward nearest neighbor, ties go to even neighbor
		{"RoundHalfEven_LessThanHalf", RoundHalfEven, i64(4), i64(10), false},
		{"RoundHalfEven_MoreThanHalf", RoundHalfEven, i64(6), i64(10), true},
		{"RoundHalfEven_ExactlyHalf_OddRemainder", RoundHalfEven, i64(5), i64(10), true},  // 0.5 rounds up to 1 (odd)
		{"RoundHalfEven_ExactlyHalf_EvenRemainder", RoundHalfEven, i64(2), i64(4), false}, // 0.5 rounds down to 0 (even)
		{"RoundHalfEven_NegativeLessThanHalf", RoundHalfEven, i64(-4), i64(10), false},
		{"RoundHalfEven_NegativeMoreThanHalf", RoundHalfEven, i64(-6), i64(10), true},
		{"RoundHalfEven_NegativeExactlyHalf_OddRemainder", RoundHalfEven, i64(-5), i64(10), true},
		{"RoundHalfEven_NegativeExactlyHalf_EvenRemainder", RoundHalfEven, i64(-2), i64(4), false},
		{"RoundUnncessary", RoundUnnecessary, i64(3), i64(10), false},
		{"Default", RoundingMode(999), i64(3), i64(10), false}, // Unknown mode defaults to no rounding
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.name), func(t *testing.T) {
			result := tc.mode.shouldRoundUp(tc.rem, tc.denom)
			if result != tc.expected {
				t.Errorf("shouldRoundUp(%v, %v) with mode %s = %t; want %t", tc.rem, tc.denom, tc.mode, result, tc.expected)
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
