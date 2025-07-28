package decimal

import (
	"fmt"
	"math/big"
	"strings"
)

type Decimal struct {
	unscaledValue *big.Int
	scale         int32
}

func New(val int64, scale int32) Decimal {
	return Decimal{
		unscaledValue: big.NewInt(val),
		scale:         scale,
	}
}

func NewInt(val int32) Decimal {
	return Decimal{
		unscaledValue: big.NewInt(int64(val)),
		scale:         0,
	}
}

func NewInt64(val int64) Decimal {
	return Decimal{
		unscaledValue: big.NewInt(val),
		scale:         0,
	}
}

func NewUint64(val uint64) Decimal {
	return Decimal{
		unscaledValue: big.NewInt(int64(val)),
		scale:         0,
	}
}

func NewBigInt(val *big.Int, scale int32) Decimal {
	return Decimal{
		unscaledValue: new(big.Int).Set(val),
		scale:         scale,
	}
}

// NewString parses a string representation of a decimal number into a Decimal.
// It supports formats like "123", "123.45", "-123.45", and scientific notation like "1.23e+5", "-4.5E-2".
func NewString(val string) (Decimal, error) {
	if val == "" {
		return Decimal{}, fmt.Errorf("cannot parse empty string to Decimal")
	}

	originalVal := val

	isNegative := false
	if val[0] == '-' {
		isNegative = true
		val = val[1:]
	} else if val[0] == '+' {
		val = val[1:]
	}

	// Check for scientific notation 'e' or 'E'
	eIndex := -1
	for i, r := range val {
		if r == 'e' || r == 'E' {
			eIndex = i
			break
		}
	}

	var mantissaStr string
	var exponent int64 = 0 // Default exponent is 0

	if eIndex != -1 {
		// Scientific notation found
		mantissaStr = val[:eIndex]
		exponentStr := val[eIndex+1:]

		if exponentStr == "" {
			return Decimal{}, fmt.Errorf("invalid scientific notation: missing exponent after 'e' in %q", originalVal)
		}

		// Parse exponent
		expBigInt := new(big.Int)
		_, ok := expBigInt.SetString(exponentStr, 10)
		if !ok {
			return Decimal{}, fmt.Errorf("invalid exponent in scientific notation: %q", originalVal)
		}
		// Convert to int64, checking for overflow if scale can be int32
		if !expBigInt.IsInt64() {
			return Decimal{}, fmt.Errorf("exponent out of int64 range: %q", originalVal)
		}
		exponent = expBigInt.Int64()

	} else {
		mantissaStr = val
	}

	// Parse the mantissa part (which might still contain a decimal point)
	parts := strings.Split(mantissaStr, ".")
	var unscaledStr string
	var mantissaScale int32

	switch len(parts) {
	case 1:
		unscaledStr = parts[0]
		mantissaScale = 0
	case 2:
		integerPart := parts[0]
		fractionalPart := parts[1]

		if fractionalPart == "" {
			// e.g., "123." or "123.e+5"
			unscaledStr = integerPart
			mantissaScale = 0
		} else {
			// Ensure fractional part contains only digits
			for _, r := range fractionalPart {
				if r < '0' || r > '9' {
					return Decimal{}, fmt.Errorf("invalid character in fractional part: %q", originalVal)
				}
			}
			unscaledStr = integerPart + fractionalPart
			mantissaScale = int32(len(fractionalPart))
		}
	default:
		return Decimal{}, fmt.Errorf("invalid decimal string format: %q (multiple decimal points in mantissa)", originalVal)
	}

	unscaledValue := new(big.Int)
	_, ok := unscaledValue.SetString(unscaledStr, 10)
	if !ok {
		return Decimal{}, fmt.Errorf("invalid characters in number part: %q", originalVal)
	}

	if isNegative {
		unscaledValue.Neg(unscaledValue)
	}

	finalScale := mantissaScale - int32(exponent)

	return Decimal{
		unscaledValue: unscaledValue,
		scale:         finalScale,
	}, nil
}
