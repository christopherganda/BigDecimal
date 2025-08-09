package decimal

import "math/big"

func (d Decimal) Add(other Decimal) Decimal {
	// Determine the final scale.
	finalScale := d.scale
	if other.scale > d.scale {
		finalScale = other.scale
	}

	// Create new Decimals with aligned scales.
	d1 := d // In a real implementation, you'd call a rescale function here.
	d2 := other

	// Example of a naive rescale implementation for demonstration.
	if d1.scale < finalScale {
		pow10 := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(finalScale-d1.scale)), nil)
		d1.unscaledValue.Mul(d1.unscaledValue, pow10)
		d1.scale = finalScale
	}
	if d2.scale < finalScale {
		pow10 := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(finalScale-d2.scale)), nil)
		d2.unscaledValue.Mul(d2.unscaledValue, pow10)
		d2.scale = finalScale
	}

	// Perform the addition on the unscaled values.
	result := new(big.Int).Add(d1.unscaledValue, d2.unscaledValue)

	return Decimal{
		unscaledValue: result,
		scale:         finalScale,
	}
}
