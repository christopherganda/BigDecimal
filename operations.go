package decimal

import (
	"math/big"
)

// rescale returns a new Decimal with its scale adjusted to a target scale.
// If the target scale is smaller than the current scale, it truncates the
// unscaled value, resulting in a loss of precision. It will never return an error.
func (d Decimal) rescale(targetScale int32) Decimal {
	// If scales are the same, return the original.
	if d.scale == targetScale {
		return d
	}

	var deltaScale int32
	var newUnscaled *big.Int

	if d.scale > targetScale {
		// Scaling down: division.
		deltaScale = d.scale - targetScale
		pow10 := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(deltaScale)), nil)
		newUnscaled = new(big.Int).Div(d.unscaledValue, pow10)
	} else {
		// Scaling up: multiplication.
		deltaScale = targetScale - d.scale
		pow10 := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(deltaScale)), nil)
		newUnscaled = new(big.Int).Mul(d.unscaledValue, pow10)
	}

	return Decimal{
		unscaledValue: newUnscaled,
		scale:         targetScale,
	}
}

func (d Decimal) Add(other Decimal) Decimal {
	finalScale := d.scale
	if other.scale > d.scale {
		finalScale = other.scale
	}

	d1 := d.rescale(finalScale)
	d2 := other.rescale(finalScale)

	return Decimal{
		unscaledValue: new(big.Int).Add(d1.unscaledValue, d2.unscaledValue),
		scale:         finalScale,
	}
}

func (d Decimal) Sub(other Decimal) Decimal {
	finalScale := d.scale
	if other.scale > d.scale {
		finalScale = other.scale
	}

	d1 := d.rescale(finalScale)
	d2 := other.rescale(finalScale)

	return Decimal{
		unscaledValue: new(big.Int).Sub(d1.unscaledValue, d2.unscaledValue),
		scale:         finalScale,
	}
}
