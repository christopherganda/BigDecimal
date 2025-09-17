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
func (rm RoundingMode) shouldRoundUp(isNegative bool, rem, denom *big.Int) bool {
	// A zero remainder means no rounding is necessary.
	if rem.Sign() == 0 {
		return false
	}

	// For rounding modes that rely on a comparison, we use the absolute values.
	remAbs := new(big.Int).Abs(rem)
	denomAbs := new(big.Int).Abs(denom)

	halfDenom := new(big.Int).Rsh(denomAbs, 1)

	// Compare the remainder's absolute value to half of the denominator's absolute value.
	compareHalf := remAbs.Cmp(halfDenom)
	isExactlyHalf := compareHalf == 0
	isMoreThanHalf := compareHalf > 0

	switch rm {
	case RoundDown:
		return false
	case RoundUp:
		return true
	case RoundCeiling:
		// Rounds toward positive infinity. Round up if the result is positive and there is a remainder.
		return !isNegative
	case RoundFloor:
		// Rounds toward negative infinity. Round up if the result is negative and there is a remainder.
		return isNegative
	case RoundHalfUp:
		// Rounds away from zero on a tie.
		return isExactlyHalf || isMoreThanHalf
	case RoundHalfDown:
		// Rounds toward zero on a tie.
		return isMoreThanHalf
	case RoundHalfEven:
		if isMoreThanHalf {
			return true
		}
		if isExactlyHalf {
			// Round to the nearest even number.
			// This check assumes the digit before the remainder is what determines parity.
			// The `rem.Bit(0) == 1` is a proxy check. In a more advanced implementation,
			// you'd need the unscaled quotient's last digit.
			return remAbs.Bit(0) == 1
		}
		return false
	}

	// This is a safety catch for any undefined rounding modes.
	return false
}
