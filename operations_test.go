package decimal

import (
	"testing"
)

func TestOperations_Add(t *testing.T) {
	tests := []struct {
		name     string
		a        Decimal
		b        Decimal
		expected Decimal
	}{
		{
			"AddPositive",
			New(5, 0),
			New(3, 0),
			New(8, 0),
		},
		{
			"AddNegative",
			New(-5, 0),
			New(-3, 0),
			New(-8, 0),
		},
		{
			"AddMixed",
			New(5, 0),
			New(-3, 0),
			New(2, 0),
		},
		{
			"AddZero",
			New(5, 0),
			New(0, 0),
			New(5, 0),
		},
		{
			"AddDifferentScales_A_Smaller",
			New(1234, 2),
			New(567, 1),
			New(6904, 2),
		}, // 12.34 + 56.70 = 69.04
		{
			"AddDifferentScales_B_Smaller",
			New(567, 1),
			New(1234, 2),
			New(6904, 2),
		}, // 56.70 + 12.34 = 69.04
		{
			"AddDifferentScales_With_Negatives",
			New(-123, 1),
			New(456, 2),
			New(-774, 2),
		}, // -12.30 + 4.56 = -7.74
		{
			"AddDifferentScales_ManyPlaces_Case1",
			New(999, 3),
			New(1, 1),
			New(1099, 3),
		}, // 0.999 + 0.100 = 1.099
		{
			"AddDifferentScales_ManyPlaces_Case2",
			New(999, 3),
			New(1, 2),
			New(1009, 3),
		}, // 0.999 + 0.010 = 1.009
		{
			"AddWithCarryOver",
			New(99, 2),
			New(2, 1),
			New(119, 2),
		}, // Corrected: 0.99 + 0.20 = 1.19
		{
			"AddWithDifferentScales_MoreZeros",
			New(12500, 4),
			New(35, 1),
			New(47500, 4),
		}, // 1.2500 + 3.5000 = 4.7500
		{
			"AddWithNegativeResult",
			New(100, 2),
			New(-250, 2),
			New(-150, 2),
		}, // 1.00 + (-2.50) = -1.50
		{
			"AddScalesToZero",
			New(100, 2),
			New(-100, 2),
			New(0, 2),
		}, // 1.00 + (-1.00) = 0.00
		{
			"AddWithLargeScaleDifference",
			New(1, 6),
			New(1, 1),
			New(100001, 6),
		}, // 0.000001 + 0.100000 = 0.100001
		{
			"AddLargeNumbers",
			Decimal{unscaledValue: bigIntFromString("9223372036854775807"),
				scale: 0},
			Decimal{unscaledValue: bigIntFromString("1"),
				scale: 0},
			Decimal{unscaledValue: bigIntFromString("9223372036854775808"),
				scale: 0}},
		{
			"	AddLargeNegativeNumbers",
			Decimal{unscaledValue: bigIntFromString("-9223372036854775807"),
				scale: 0},
			Decimal{unscaledValue: bigIntFromString("-1"),
				scale: 0},
			Decimal{unscaledValue: bigIntFromString("-9223372036854775808"),
				scale: 0}},
		{
			"AddLargeMixedNumbers",
			Decimal{unscaledValue: bigIntFromString("9223372036854775807"),
				scale: 0},
			Decimal{unscaledValue: bigIntFromString("-9223372036854775807"),
				scale: 0},
			Decimal{unscaledValue: bigIntFromString("0"),
				scale: 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Add(tt.b)
			if result.unscaledValue.Cmp(tt.expected.unscaledValue) != 0 || result.scale != tt.expected.scale {
				t.Errorf("Add(%v, %v) = %v, want %v",
					tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestOperations_Sub(t *testing.T) {
	tests := []struct {
		name     string
		a        Decimal
		b        Decimal
		expected Decimal
	}{
		{
			"SubtractPositive",
			New(8, 0),
			New(3, 0),
			New(5, 0),
		},
		{
			"SubtractNegative",
			New(-8, 0),
			New(-3, 0),
			New(-5, 0),
		},
		{
			"SubtractMixed",
			New(5, 0),
			New(-3, 0),
			New(8, 0),
		},
		{
			"SubtractZero",
			New(5, 0),
			New(0, 0),
			New(5, 0),
		},
		{
			"SubtractDifferentScales_A_Smaller",
			New(567, 1),  // 56.7
			New(1234, 2), // 12.34
			New(4436, 2), // 44.36
		},
		{
			"SubtractDifferentScales_B_Smaller",
			New(1234, 2),  // 12.34
			New(567, 1),   // 56.7
			New(-4436, 2), // -44.36
		},
		{
			"SubtractDifferentScales_With_Negatives",
			New(-123, 1), // -12.3
			New(-456, 2), // -4.56
			New(-774, 2), // -7.74
		},
		{
			"SubtractWithBorrow",
			New(119, 2), // 1.19
			New(2, 1),   // 0.2
			New(99, 2),  // 0.99
		},
		{
			"SubtractToZero",
			New(1250, 2), // 12.50
			New(125, 1),  // 12.5
			New(0, 2),
		},
		{
			"SubtractWithLargeScaleDifference",
			New(1, 1),     // 0.1
			New(1, 6),     // 0.000001
			New(99999, 6), // 0.099999
		},
		{
			"SubtractLargeNumbers",
			Decimal{unscaledValue: bigIntFromString("9223372036854775808"), scale: 0},
			Decimal{unscaledValue: bigIntFromString("1"), scale: 0},
			Decimal{unscaledValue: bigIntFromString("9223372036854775807"), scale: 0},
		},
		{
			"SubtractLargeNegativeNumbers",
			Decimal{unscaledValue: bigIntFromString("-9223372036854775808"), scale: 0},
			Decimal{unscaledValue: bigIntFromString("-1"), scale: 0},
			Decimal{unscaledValue: bigIntFromString("-9223372036854775807"), scale: 0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.a.Subtract(tc.b)

			if result.unscaledValue.Cmp(tc.expected.unscaledValue) != 0 {
				t.Errorf("unscaledValue mismatch for %s: got %s, want %s", tc.name, result.unscaledValue, tc.expected.unscaledValue)
			}
			if result.scale != tc.expected.scale {
				t.Errorf("scale mismatch for %s: got %d, want %d", tc.name, result.scale, tc.expected.scale)
			}
		})
	}
}

func TestOperations_Multiply(t *testing.T) {
	tests := []struct {
		name     string
		a        Decimal
		b        Decimal
		expected Decimal
	}{
		{
			"MultiplyTwoPositiveNumbers",
			New(123, 2),   // 1.23
			New(456, 2),   // 4.56
			New(56088, 4), // 5.6088
		},
		{
			"MultiplyWithZero",
			New(1234, 3), // 1.234
			New(0, 1),    // 0.0
			New(0, 4),    // 0.0000
		},
		{
			"MultiplyWithNegative",
			New(125, 2),  // 1.25
			New(-5, 1),   // -0.5
			New(-625, 3), // -0.625
		},
		{
			"MultiplyTwoNegatives",
			New(-125, 2), // -1.25
			New(-5, 1),   // -0.5
			New(625, 3),  // 0.625
		},
		{
			"MultiplyDifferentScales_Simple",
			New(123, 2), // 1.23
			New(4, 0),   // 4
			New(492, 2), // 4.92
		},
		{
			"MultiplyDifferentScales_Complex",
			New(1234, 2), // 12.34
			New(5, 1),    // 0.5
			New(6170, 3), // 6.170
		},
		{
			"MultiplySmallNumbers",
			New(1, 2), // 0.01
			New(1, 2), // 0.01
			New(1, 4), // 0.0001
		},
		{
			"MultiplySmallNumbers",
			New(2, 1), // 0.2
			New(3, 1), // 0.3
			New(6, 2), // 0.06
		},
		{
			"MultiplyLargeNumbers",
			Decimal{unscaledValue: bigIntFromString("9223372036854775807"), scale: 0}, // A large number
			Decimal{unscaledValue: bigIntFromString("100"), scale: 0},                 // 100
			Decimal{unscaledValue: bigIntFromString("922337203685477580700"), scale: 0},
		},
		{
			"MultiplyLargeAndScaled",
			Decimal{unscaledValue: bigIntFromString("123456789"), scale: 4},          // 12345.6789
			Decimal{unscaledValue: bigIntFromString("987654321"), scale: 5},          // 9876.54321
			Decimal{unscaledValue: bigIntFromString("121932631112635269"), scale: 9}, // 121932631.112635269
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.a.Multiply(tc.b)

			if result.unscaledValue.Cmp(tc.expected.unscaledValue) != 0 {
				t.Errorf("unscaledValue mismatch for %s: got %s, want %s", tc.name, result.unscaledValue, tc.expected.unscaledValue)
			}
			if result.scale != tc.expected.scale {
				t.Errorf("scale mismatch for %s: got %d, want %d", tc.name, result.scale, tc.expected.scale)
			}
		})
	}
}
