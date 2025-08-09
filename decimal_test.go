package decimal

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name  string
		val   int64
		scale int32
	}{
		{"valid scale", 123, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.val, tt.scale)
			if got.unscaledValue.String() != fmt.Sprintf("%d", tt.val) {
				t.Errorf("New(%v) = %v, want %v", tt.val, got.unscaledValue, tt.val)
			}
			if got.scale != tt.scale {
				t.Errorf("New(%v) scale = %v, want %v", tt.val, got.scale, tt.scale)
			}
		})
	}
}

func TestNewInt(t *testing.T) {
	tests := []struct {
		input   int32
		wantVal string
	}{
		{0, "0"},
		{123, "123"},
		{-123, "-123"},
		{math.MaxInt32, "2147483647"},
		{math.MinInt32, "-2147483648"},
	}
	for _, tt := range tests {
		t.Run(tt.wantVal, func(t *testing.T) {
			got := NewFromInt(tt.input)
			if got.unscaledValue.String() != tt.wantVal {
				t.Errorf("NewInt(%v) = %v, want %v", tt.input, got.unscaledValue, tt.wantVal)
			}
			if got.scale != 0 {
				t.Errorf("NewInt(%v) scale = %v, want 0", tt.input, got.scale)
			}
		})
	}
}

func TestNewInt64(t *testing.T) {
	tests := []struct {
		input   int64
		wantVal string
	}{
		{0, "0"},
		{123, "123"},
		{-123, "-123"},
		{math.MaxInt64, "9223372036854775807"},
		{math.MinInt64, "-9223372036854775808"},
	}
	for _, tt := range tests {
		t.Run(tt.wantVal, func(t *testing.T) {
			got := NewFromInt64(tt.input)
			if got.unscaledValue.String() != tt.wantVal {
				t.Errorf("NewInt64(%v) = %v, want %v", tt.input, got.unscaledValue, tt.wantVal)
			}
			if got.scale != 0 {
				t.Errorf("NewInt64(%v) scale = %v, want 0", tt.input, got.scale)
			}
		})
	}
}

func TestNewUint64(t *testing.T) {
	tests := []struct {
		input   uint64
		wantVal string
	}{
		{0, "0"},
		{123, "123"},
		{math.MaxUint32, "4294967295"},
		{math.MaxUint64, "18446744073709551615"},
	}
	for _, tt := range tests {
		t.Run(tt.wantVal, func(t *testing.T) {
			got := NewFromUint64(tt.input)
			if got.unscaledValue.String() != tt.wantVal {
				t.Errorf("NewUint64(%v) = %v, want %v", tt.input, got.unscaledValue, tt.wantVal)
			}
			if got.scale != 0 {
				t.Errorf("NewUint64(%v) scale = %v, want 0", tt.input, got.scale)
			}
		})
	}
}

func TestNewBigInt(t *testing.T) {
	tests := []struct {
		name    string
		input   *big.Int
		scale   int32
		wantVal string
		wantErr bool
	}{
		{"nil", nil, 0, "", true},
		{"zero", big.NewInt(0), 0, "0", false},
		{"positive", big.NewInt(123), 2, "123", false},
		{"negative", big.NewInt(-123), 2, "-123", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFromBigInt(tt.input, tt.scale)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBigInt() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got.unscaledValue.String() != tt.wantVal {
				t.Errorf("NewBigInt() = %v, want %v", got.unscaledValue, tt.wantVal)
			}
		})
	}
}

func TestNewString(t *testing.T) {
	tests := []struct {
		input     string
		wantVal   string
		wantScale int32
		wantErr   bool
	}{
		{"0", "0", 0, false},
		{"123", "123", 0, false},
		{"-123", "-123", 0, false},
		{"123.45", "12345", 2, false},
		{"-123.45", "-12345", 2, false},
		{"123.", "123", 0, false},
		{"0.123", "123", 3, false},
		{"1.23e+2", "123", 0, false},
		{"-1.23e-2", "-123", 4, false},
		{"123.45.67", "", 0, true},
		{"123.4x5", "", 0, true},
		{"", "", 0, true},
		{".", "", 0, true},
		{"1.23e+100", "123", -98, false},
		{"1.23e-100", "123", 102, false},
		{"000123.45", "12345", 2, false},
		{"123.450", "123450", 3, false},
		{"0.0", "0", 1, false},
		{"+", "", 0, true},
		{"-", "", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := NewFromString(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewString(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr {
				if got.unscaledValue.String() != tt.wantVal {
					t.Errorf("NewString(%q) = %v, want %v", tt.input, got.unscaledValue, tt.wantVal)
				}
				if got.scale != tt.wantScale {
					t.Errorf("NewString(%q) scale = %v, want %v", tt.input, got.scale, tt.wantScale)
				}
			}
		})
	}
}

func TestNewFromFloat64(t *testing.T) {
	tests := []struct {
		input   float64
		wantErr bool
	}{
		{0.0, false},
		{123.0, false},
		{-123.0, false},
		{123.45, false},
		{-123.45, false},
		{1e-20, false},
		{1e20, false},
		{math.NaN(), true},
		{math.Inf(1), true},
		{math.Inf(-1), true},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%g", tt.input), func(t *testing.T) {
			got, err := NewFromFloat64(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFromFloat64(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr {
				// Convert back to float64 for comparison
				rat := new(big.Rat).SetFrac(got.unscaledValue, pow10(got.scale))
				f, _ := rat.Float64()
				if math.Abs(f-tt.input) > 1e-12 {
					t.Errorf("NewFromFloat64(%v) = %v, want %v", tt.input, f, tt.input)
				}
			}
		})
	}
}

func TestNewFromRat(t *testing.T) {
	tests := []struct {
		name         string
		num          int64
		denom        int64
		precision    int32
		roundingMode RoundingMode
		wantVal      string
		wantScale    int32
		wantErr      bool
	}{
		{"integer", 5, 1, 0, RoundHalfEven, "5", 0, false},
		{"half", 1, 2, 1, RoundHalfEven, "5", 1, false},
		{"third", 1, 3, 2, RoundHalfEven, "33", 2, false},
		{"third up", 1, 3, 2, RoundUp, "34", 2, false},
		{"negative half", -1, 2, 1, RoundHalfEven, "-5", 1, false},
		{"zero", 0, 1, 0, RoundHalfEven, "0", 0, false},
		{"nil rat", 0, 0, 0, RoundHalfEven, "", 0, true},
		{"negative precision", 1, 1, -1, RoundHalfEven, "", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rat *big.Rat
			if tt.denom != 0 {
				rat = new(big.Rat).SetFrac(big.NewInt(tt.num), big.NewInt(tt.denom))
			}
			got, err := NewFromRat(rat, tt.precision, tt.roundingMode)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFromRat() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if got.unscaledValue.String() != tt.wantVal {
					t.Errorf("NewFromRat() = %v, want %v", got.unscaledValue, tt.wantVal)
				}
				if got.scale != tt.wantScale {
					t.Errorf("NewFromRat() scale = %v, want %v", got.scale, tt.wantScale)
				}
			}
		})
	}
}

func TestNewFromBytes(t *testing.T) {
	tests := []struct {
		input     []byte
		wantVal   string
		wantScale int32
		wantErr   bool
	}{
		{[]byte("0"), "0", 0, false},
		{[]byte("123"), "123", 0, false},
		{[]byte("-123"), "-123", 0, false},
		{[]byte("123.45"), "12345", 2, false},
		{[]byte("-123.45"), "-12345", 2, false},
		{[]byte("1.23e+2"), "123", 0, false},
		{[]byte(""), "", 0, true},
		{[]byte("abc"), "", 0, true},
		{[]byte("123.45.67"), "", 0, true},
		{[]byte("123.4x5"), "", 0, true},
		{[]byte("+"), "", 0, true},
		{[]byte("-"), "", 0, true},
		{[]byte("."), "", 0, true},
	}
	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			got, err := NewFromBytes(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFromBytes(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr {
				if got.unscaledValue.String() != tt.wantVal {
					t.Errorf("NewFromBytes(%q) = %v, want %v", tt.input, got.unscaledValue, tt.wantVal)
				}
				if got.scale != tt.wantScale {
					t.Errorf("NewFromBytes(%q) scale = %v, want %v", tt.input, got.scale, tt.wantScale)
				}
			}
		})
	}
}

func TestPow10(t *testing.T) {
	tests := []struct {
		input int32
		want  string
	}{
		{0, "1"},
		{1, "10"},
		{2, "100"},
		{10, "10000000000"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("pow10_%d", tt.input), func(t *testing.T) {
			got := pow10(tt.input)
			if got.String() != tt.want {
				t.Errorf("pow10(%d) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
	// Negative power should panic
	t.Run("negative", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("pow10(-1) did not panic")
			}
		}()
		pow10(-1)
	})
}

func TestPow10Cache(t *testing.T) {
	// Clear cache
	powersOfTen = make(map[int32]*big.Int, 128)
	p1 := pow10(10)
	p2 := pow10(10)
	if p1 != p2 {
		t.Error("pow10 cache not working: got different instances for same power")
	}
}

func TestDecimal_Scan(t *testing.T) {
	tests := []struct {
		input     string
		wantVal   string
		wantScale int32
		wantErr   bool
	}{
		{"123", "123", 0, false},
		{"-123.45", "-12345", 2, false},
		{"0.00123", "123", 5, false},
		{"1.23e+2", "123", 0, false},
		{"invalid", "", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			var d Decimal
			err := d.Scan(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scan(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr {
				if d.unscaledValue.String() != tt.wantVal {
					t.Errorf("Scan(%q) = %v, want %v", tt.input, d.unscaledValue, tt.wantVal)
				}
				if d.scale != tt.wantScale {
					t.Errorf("Scan(%q) scale = %v, want %v", tt.input, d.scale, tt.wantScale)
				}
			}
		})
	}
}

func TestDecimal_String(t *testing.T) {
	tests := []struct {
		input Decimal
		want  string
	}{
		{New(123, 0), "123"},
		{New(-12345, 2), "-123.45"},
		{New(0, 0), "0"},
		{New(100000, 2), "1000.00"},
		{New(100000, -2), "10000000"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.input.String()
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
