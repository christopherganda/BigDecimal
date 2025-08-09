package decimal

import (
	"fmt"
	"math/big"
)

// RoundingMode determines how decimal numbers are rounded
type RoundingMode int

const (
	// RoundDown rounds toward zero
	RoundDown RoundingMode = iota

	// RoundUp rounds away from zero
	RoundUp

	// RoundCeiling rounds toward positive infinity
	RoundCeiling

	// RoundFloor rounds toward negative infinity
	RoundFloor

	// RoundHalfUp rounds toward nearest neighbor, ties toward positive infinity
	RoundHalfUp

	// RoundHalfDown rounds toward nearest neighbor, ties toward negative infinity
	RoundHalfDown

	// RoundHalfEven rounds toward nearest neighbor, ties toward even neighbor (Default)
	RoundHalfEven

	// RoundUnnecessary throws error if rounding is necessary
	RoundUnnecessary
)

// String returns the string representation of the rounding mode
func (rm RoundingMode) String() string {
	switch rm {
	case RoundDown:
		return "RoundDown"
	case RoundUp:
		return "RoundUp"
	case RoundCeiling:
		return "RoundCeiling"
	case RoundFloor:
		return "RoundFloor"
	case RoundHalfUp:
		return "RoundHalfUp"
	case RoundHalfDown:
		return "RoundHalfDown"
	case RoundHalfEven:
		return "RoundHalfEven"
	case RoundUnnecessary:
		return "RoundUnnecessary"
	default:
		return fmt.Sprintf("RoundingMode(%d)", rm)
	}
}

// shouldRoundUp determines if we should round up based on the remainder and denominator
func (rm RoundingMode) shouldRoundUp(rem, denom *big.Int) bool {
	// If remainder is zero, no rounding needed
	if rem.Sign() == 0 {
		return false
	}

	// Get absolute values for comparison
	remAbs := new(big.Int).Abs(rem)
	denomAbs := new(big.Int).Abs(denom)

	// Calculate half of denominator for comparison
	half := new(big.Int).Rsh(denomAbs, 1)

	// Compare remainder with half of denominator
	compareHalf := remAbs.Cmp(half)

	switch rm {
	case RoundDown:
		return false

	case RoundUp:
		return true

	case RoundCeiling:
		return rem.Sign() > 0

	case RoundFloor:
		return rem.Sign() < 0

	case RoundHalfUp:
		return compareHalf >= 0

	case RoundHalfDown:
		return compareHalf > 0

	case RoundHalfEven:
		if compareHalf > 0 {
			return true
		}
		if compareHalf < 0 {
			return false
		}
		// If exactly half, round to even
		return rem.Bit(0) == 1

	case RoundUnnecessary:
		if rem.Sign() != 0 {
			panic("rounding necessary but RoundUnnecessary specified")
		}
		return false

	default:
		panic("unknown rounding mode")
	}
}
