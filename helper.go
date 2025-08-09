package decimal

import (
	"math/big"
)

func bigIntFromString(s string) *big.Int {
	if s == "" {
		return nil
	}
	val := new(big.Int)
	_, ok := val.SetString(s, 10)
	if !ok {
		return nil
	}
	return val
}
