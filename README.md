# Go BigDecimal

A comprehensive and precise arbitrary-precision decimal number library for Go, inspired by Java's `BigDecimal`. This library provides a robust, complete set of features for financial, scientific, and general-purpose decimal arithmetic.

## ✨ Features

  * **Arbitrary Precision**: Handles numbers of any size without precision loss.
  * **Comprehensive Operations**: Supports all standard arithmetic operations (`Add`, `Sub`, `Mul`, `Div`), as well as `Abs`, `Neg`, and `Rem`.
  * **Advanced Rounding**: Provides a full suite of rounding modes for precise control over calculations, essential for financial applications.
  * **Database & JSON Integration**: Seamlessly integrates with SQL databases and JSON encoding/decoding via standard Go interfaces (`sql.Scanner`, `driver.Valuer`, `json.Marshaler`, `json.Unmarshaler`).
  * **Multiple Initializers**: Create `Decimal` values from a wide variety of types, including `int64`, `float64`, `string`, `big.Int`, and `big.Rat`.
  * **Zero-Allocation Conversion**: High-performance methods for converting to and from `string` and `[]byte` without extra memory allocations.

## 📦 Installation

To start using `go-bigdecimal`, simply run the following command:

```bash
go get github.com/christopherganda/go-bigdecimal
```

## 🚀 Usage

### Initializing a Decimal

You can create a new `Decimal` from various data types.

```go
package main

import (
	"fmt"
	"math/big"
	bigdecimal "github.com/your-username/go-bigdecimal"
)

func main() {
	// From int64 and scale (123.45)
	d1 := bigdecimal.New(12345, 2) 
	
	// From string ("-987.654321")
	d2, err := bigdecimal.NewFromString("-987.654321") 
	if err != nil {
		panic(err)
	}

	// From float64 (123.456)
	d3 := bigdecimal.NewFromFloat64(123.456) 

	// From big.Int and scale (123456789.0)
	bi := big.NewInt(123456789)
	d4 := bigdecimal.NewFromBigInt(bi, 0) 

	fmt.Println(d1, d2, d3, d4)
}
```

### Performing Operations

All operations are immutable and return a new `Decimal` value.

```go
package main

import (
	"fmt"
	"github.com/your-username/go-bigdecimal"
)

func main() {
	a := bigdecimal.New(10, 1) // 1.0
	b := bigdecimal.New(3, 1)  // 0.3

	// Addition: 1.0 + 0.3 = 1.3
	sum := a.Add(b)
	fmt.Printf("Sum: %s\n", sum.String()) 

	// Subtraction: 1.0 - 0.3 = 0.7
	diff := a.Sub(b)
	fmt.Printf("Difference: %s\n", diff.String()) 

	// Multiplication: 1.0 * 0.3 = 0.3
	prod := a.Mul(b)
	fmt.Printf("Product: %s\n", prod.String()) 

	// Division: 1.0 / 0.3 = 3.333... (rounded to 5 decimal places)
	// Precision is required for non-terminating decimals.
	// You must define your rounding modes, e.g., bigdecimal.RoundHalfUp
	// quotient := a.Div(b, 5, bigdecimal.RoundHalfUp) 
	// fmt.Printf("Quotient: %s\n", quotient.String()) 
}
```

## 📖 API Documentation

### Initializers

  * `New(val int64, scale int32)`: Creates a new `Decimal` from a scaled `int64`.
  * `NewFromInt(val int)`: Creates a new `Decimal` from an `int`.
  * `NewFromInt64(val int64)`: Creates a new `Decimal` from an `int64`.
  * `NewFromUint64(val uint64)`: Creates a new `Decimal` from a `uint64`.
  * `NewFromBigInt(val *big.Int, scale int32)`: Creates a `Decimal` from a `big.Int` and a scale.
  * `NewFromString(s string)`: Creates a `Decimal` from a string.
  * `NewFromFloat64(val float64)`: Creates a `Decimal` from a `float64`.
  * `NewFromRat(val *big.Rat)`: Creates a `Decimal` from a `big.Rat`.
  * `NewFromBytes(b []byte)`: Creates a `Decimal` from a byte slice.

### Operations

  * `Add(other Decimal) Decimal`: Returns the sum of two decimals.
  * `Sub(other Decimal) Decimal`: Returns the difference.
  * `Mul(other Decimal) Decimal`: Returns the product.
  * `Div(other Decimal, precision int32) Decimal`: Returns the quotient. **Note**: This function requires precision and may have an optional rounding mode argument depending on your final implementation.
  * `Rem(other Decimal) Decimal`: Returns the remainder.
  * `DivRem(other Decimal) (Decimal, Decimal)`: Returns the quotient and remainder.
  * `Abs() Decimal`: Returns the absolute value.
  * `Neg() Decimal`: Returns the negation.
  * `Cmp(other Decimal) int`: Compares two decimals, returning -1, 0, or 1.
  * `Sign() int`: Returns the sign of the decimal.

### Converters & Utilities

  * `String() string`: Returns a canonical string representation.
  * `StringFixed(places int32) string`: Returns a string with a fixed number of decimal places, with rounding.
  * `Int64() (int64, bool)`: Converts the decimal to `int64`, returning a boolean indicating if the conversion was exact.
  * `Float64() (float64, bool)`: Converts the decimal to `float64`, returning a boolean indicating if the conversion was exact.
  * `BigInt() *big.Int`: Returns the unscaled value as a `big.Int`.
  * `Scale() int32`: Returns the current scale of the number.
  * `Rescale(newScale int32, roundingMode RoundingMode) Decimal`: Changes the scale, rounding the number if necessary.
  * `IsZero()`: Returns true if the decimal is zero.

## 🤝 Contributing

Contributions are welcome\! Please open an issue or a pull request for any features, bug fixes, or documentation improvements.
