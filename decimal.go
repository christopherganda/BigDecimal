package decimal

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"sync"
)

type Decimal struct {
	unscaledValue *big.Int
	scale         int32
}

// Add more cached powers and thread safety
var (
	powersOfTenMutex sync.RWMutex
	powersOfTen      = make(map[int32]*big.Int, 128) // Increase initial capacity
)

func init() {
	// Pre-calculate more powers
	for i := int32(0); i <= 38; i++ { // Common powers for uint128
		pow10(i)
	}
}

func pow10(n int32) *big.Int {
	if n < 0 {
		// For negative powers, we actually need 1 / 10^(-n).
		// This function primarily provides positive powers for multiplication.
		// Division by 10^N is handled by dividing by pow10(N).
		panic(fmt.Sprintf("pow10 does not support negative exponents for direct multiplication: %d", n))
	}
	powersOfTenMutex.RLock()
	if p, ok := powersOfTen[n]; ok {
		powersOfTenMutex.RUnlock()
		return p
	}
	powersOfTenMutex.RUnlock()

	powersOfTenMutex.Lock()
	defer powersOfTenMutex.Unlock()

	// Double-check after acquiring write lock
	if p, ok := powersOfTen[n]; ok {
		return p
	}

	p := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(n)), nil)
	powersOfTen[n] = p
	return p
}

func New(val int64, scale int32) Decimal {
	return Decimal{
		unscaledValue: big.NewInt(val),
		scale:         scale,
	}
}

func NewFromInt(val int32) Decimal {
	return Decimal{
		unscaledValue: big.NewInt(int64(val)),
		scale:         0,
	}
}

func NewFromInt64(val int64) Decimal {
	return Decimal{
		unscaledValue: big.NewInt(val),
		scale:         0,
	}
}

func NewFromUint64(val uint64) Decimal {
	return Decimal{
		unscaledValue: new(big.Int).SetUint64(val), // Correctly uses SetUint64
		scale:         0,
	}
}

func NewFromBigInt(val *big.Int, scale int32) (Decimal, error) {
	if val == nil {
		return Decimal{}, fmt.Errorf("nil big.Int value")
	}
	return Decimal{
		unscaledValue: new(big.Int).Set(val),
		scale:         scale,
	}, nil
}

// NewFromString parses a string representation of a decimal number into a Decimal.
// It supports formats like "123", "123.45", "-123.45", and scientific notation like "1.23e+5", "-4.5E-2".
// For example: 1.23e+5
// mantissaStr = "1.23"
// exponentStr = "5"
// unscaledStr = "123", mantissaScale = 2(2 digits after the decimal point)
// finalScale = mantissaScale - exponent = 2 - 5 = -3
// Decimal{unscaledValue: 123, scale: -3} representing 123 * 10^3 = 123000
// if scale is positive, it means the number has that many digits after the decimal point.
// if scale is negative, we times the number by 10^(scale) to get the actual value.
func NewFromString(val string) (Decimal, error) {
	if val == "" {
		return Decimal{}, fmt.Errorf("cannot parse empty string to Decimal")
	}

	originalVal := val

	isNegative := false
	switch val[0] {
	case '-':
		isNegative = true
		val = val[1:]
	case '+':
		val = val[1:]
	default:
		// No sign, continue with the original value
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
	var exponent int64 = 0

	if eIndex != -1 {
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

// NewFromFloat64 creates a new Decimal from a float64 value.
// This conversion aims for the most precise decimal representation of the float64's binary value.
// It converts the float64 to a *big.Rat and then uses NewFromRat.
// The default precision for this conversion is set to 64 decimal places, which is usually
// sufficient to capture the full precision of a float64 (approx 15-17 digits).
func NewFromFloat64(val float64) (Decimal, error) {
	if math.IsInf(val, 0) {
		return Decimal{}, fmt.Errorf("cannot convert infinity to Decimal")
	}
	if math.IsNaN(val) {
		return Decimal{}, fmt.Errorf("cannot convert NaN to Decimal")
	}
	if val == 0 {
		return Decimal{unscaledValue: big.NewInt(0), scale: 0}, nil
	}

	// Convert float64 to its exact rational representation.
	// This captures the exact binary value of the float64.
	rat := new(big.Rat).SetFloat64(val)
	if rat == nil {
		// This case should ideally not happen for valid non-NaN/Inf floats
		return Decimal{}, fmt.Errorf("failed to convert float64 to *big.Rat: %v", val)
	}

	// Convert the *big.Rat to Decimal with a sufficiently high precision
	// and a default rounding mode. 64 decimal places is chosen as it's
	// more than enough to exactly represent any float64 in decimal form.
	// RoundHalfEven is a good default for general numerical conversions.
	return NewFromRat(rat, 64, RoundHalfEven)
}

// NewFromRat creates a new Decimal from a *big.Rat (rational number).
// It converts the rational number to a Decimal with the specified precision and rounding mode.
// This is a robust conversion that handles non-terminating decimals by rounding.
func NewFromRat(val *big.Rat, precision int32, roundingMode RoundingMode) (Decimal, error) {
	if val == nil {
		return Decimal{}, fmt.Errorf("cannot create Decimal from nil *big.Rat")
	}
	if precision < 0 {
		return Decimal{}, fmt.Errorf("precision must be non-negative for NewFromRat, got %d", precision)
	}

	// Handle zero explicitly
	if val.Num().Cmp(big.NewInt(0)) == 0 {
		return Decimal{unscaledValue: big.NewInt(0), scale: precision}, nil
	}

	// If the rational number is an exact integer and target precision is 0,
	// we can take a fast path.
	if val.IsInt() && precision == 0 {
		return Decimal{
			unscaledValue: new(big.Int).Set(val.Num()),
			scale:         0,
		}, nil
	}

	// Numerator and Denominator
	num := new(big.Int).Set(val.Num())
	den := new(big.Int).Set(val.Denom())

	// Multiply numerator by 10^precision to shift the decimal point for division
	tempNum := new(big.Int).Mul(num, pow10(precision))

	// Perform division: quotient and remainder
	quotient := new(big.Int)
	remainder := new(big.Int)
	quotient.QuoRem(tempNum, den, remainder)

	// Determine if rounding is needed based on the remainder
	// If remainder is zero, no rounding needed.
	if remainder.Cmp(big.NewInt(0)) == 0 {
		return Decimal{
			unscaledValue: quotient,
			scale:         precision,
		}, nil
	}

	// If RoundUnnecessary and remainder is not zero, return error
	if roundingMode == RoundUnnecessary {
		return Decimal{}, fmt.Errorf("rounding necessary for NewFromRat with RoundUnnecessary mode: %s at precision %d", val.String(), precision)
	}

	// absRemainder is used for comparison against half of the denominator
	absRemainder := new(big.Int).Abs(remainder)

	// halfDivisor is divisor / 2
	halfDivisor := new(big.Int).Rsh(den, 1) // Efficiently divide by 2

	// Compare abs(remainder) with (denominator / 2)
	cmpResult := absRemainder.Cmp(halfDivisor) // -1 if absR < halfDivisor, 0 if absR == halfDivisor, 1 if absR > halfDivisor

	// Check if exactly halfway (abs(remainder) * 2 == denominator)
	// This is important for HalfUp, HalfDown, HalfEven.
	isHalfway := (cmpResult == 0) && (new(big.Int).Mul(absRemainder, big.NewInt(2)).Cmp(den) == 0)

	// Determine if we need to increment the quotient based on rounding mode
	shouldIncrement := false
	quotientSign := quotient.Sign() // -1 for negative, 0 for zero, 1 for positive

	switch roundingMode {
	case RoundUp:
		// Round away from zero
		shouldIncrement = true
	case RoundDown:
		// Round towards zero (no increment needed if remainder exists)
		shouldIncrement = false // Already truncated by QuoRem
	case RoundCeiling:
		// Round towards positive infinity
		shouldIncrement = quotientSign >= 0 // Increment if positive or zero
	case RoundFloor:
		// Round towards negative infinity
		shouldIncrement = quotientSign < 0 // Increment if negative
	case RoundHalfUp:
		// Round towards nearest neighbor; if equidistant, round up.
		shouldIncrement = cmpResult == 1 || (cmpResult == 0 && isHalfway)
	case RoundHalfDown:
		// Round towards nearest neighbor; if equidistant, round down.
		shouldIncrement = cmpResult == 1
	case RoundHalfEven:
		// Round towards nearest neighbor; if equidistant, round to even.
		if cmpResult == 1 { // Remainder > halfDivisor
			shouldIncrement = true
		} else if cmpResult == 0 && isHalfway { // Exactly halfway
			// Check if the last digit of the quotient is odd
			// Note: This checks if the *quotient* is odd/even, not the remainder.
			// For RoundHalfEven, you check the *truncated* digit before rounding.
			// If the truncated digit is even, round down; if odd, round up.
			// This is equivalent to checking if quotient is even/odd.
			lastDigitOfQuotient := new(big.Int).Mod(quotient, big.NewInt(2))
			shouldIncrement = lastDigitOfQuotient.Cmp(big.NewInt(1)) == 0 // Increment if odd
		}
	default:
		return Decimal{}, fmt.Errorf("unsupported rounding mode: %v", roundingMode)
	}

	// Apply increment if determined
	if shouldIncrement {
		if quotientSign >= 0 { // Positive or zero quotient
			quotient.Add(quotient, big.NewInt(1))
		} else { // Negative quotient
			quotient.Sub(quotient, big.NewInt(1))
		}
	}

	return Decimal{
		unscaledValue: quotient,
		scale:         precision,
	}, nil
}

func NewFromBytes(val []byte) (Decimal, error) {
	if len(val) == 0 {
		return Decimal{}, fmt.Errorf("cannot parse empty bytes to Decimal")
	}
	// Reuse NewString logic but work directly with bytes
	// to avoid string conversion
	return NewFromString(string(val))
}

// Scan implements the sql.Scanner interface.
// It allows our Decimal type to be scanned directly from a database query.
func (d *Decimal) Scan(value interface{}) error {
	if value == nil {
		// Handle a NULL value from the database
		d.unscaledValue = new(big.Int)
		d.scale = 0
		return nil
	}

	var parsedDecimal Decimal
	var err error
	switch v := value.(type) {
	case string:
		parsedDecimal, err = NewFromString(v)
	case []byte:
		parsedDecimal, err = NewFromBytes(v)
	default:
		// Return an error for unsupported types
		return fmt.Errorf("unsupported type for Decimal Scan: %T", value)
	}

	if err != nil {
		return fmt.Errorf("failed to scan string to Decimal: %w", err)
	}

	// Set the receiver's fields to the newly parsed value
	*d = parsedDecimal

	return nil
}

// String returns the string representation of the decimal.
func (d Decimal) String() string {
	if d.unscaledValue == nil {
		return "<nil>"
	}

	numStr := d.unscaledValue.String()
	scale := d.scale

	if scale == 0 {
		return numStr
	}

	// Handle negative scale (exponent)
	if scale < 0 {
		// Multiply by 10^(-scale)
		exponent := -scale
		multiplier := pow10(exponent)
		// Create a new big.Int to avoid modifying the original
		tempInt := new(big.Int).Set(d.unscaledValue)
		tempInt.Mul(tempInt, multiplier)
		return tempInt.String()
	}

	// Handle positive scale (fractional part)
	if scale > 0 {
		// Insert decimal point
		integerPart := numStr[:len(numStr)-int(scale)]
		fractionalPart := numStr[len(numStr)-int(scale):]

		// Pad with leading zeros if integer part is empty
		if integerPart == "" {
			integerPart = "0"
		}

		// Pad with trailing zeros if fractional part is too short
		if len(fractionalPart) < int(scale) {
			fractionalPart = fractionalPart + strings.Repeat("0", int(scale)-len(fractionalPart))
		}

		return integerPart + "." + fractionalPart
	}

	return ""
}
