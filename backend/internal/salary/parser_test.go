package salary_test

import (
	"fmt"
	"jobstream/internal/salary"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedMin *int64
		expectedMax *int64
	}{
		{
			name:        "Empty string",
			input:       "",
			expectedMin: nil,
			expectedMax: nil,
		},
		{
			name:        "Single exact salary in dollars",
			input:       "$80000",
			expectedMin: int64Ptr(80000),
			expectedMax: nil,
		},
		{
			name:        "Single salary with k suffix",
			input:       "120k USD",
			expectedMin: int64Ptr(120000),
			expectedMax: nil,
		},
		{
			name:        "Salary range with k suffix",
			input:       "$120k - $150k",
			expectedMin: int64Ptr(120000),
			expectedMax: int64Ptr(150000),
		},
		{
			name:        "Salary range with decimals",
			input:       "120.5k - 150.5k",
			expectedMin: int64Ptr(120500),
			expectedMax: int64Ptr(150500),
		},
		{
			name:        "Salary in millions",
			input:       "1.2m",
			expectedMin: int64Ptr(1200000),
			expectedMax: nil,
		},
		{
			name:        "Reversed range order",
			input:       "180k - 130k",
			expectedMin: int64Ptr(130000),
			expectedMax: int64Ptr(180000),
		},
		{
			name:        "No numbers",
			input:       "Competitive salary / DOE",
			expectedMin: nil,
			expectedMax: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := salary.Parse(tt.input)

			if !compareInt64Ptr(result.Min, tt.expectedMin) {
				t.Errorf("salary.Parse(%q) min = %v, expected %v", tt.input, valStr(result.Min), valStr(tt.expectedMin))
			}

			if !compareInt64Ptr(result.Max, tt.expectedMax) {
				t.Errorf("salary.Parse(%q) max = %v, expected %v", tt.input, valStr(result.Max), valStr(tt.expectedMax))
			}
		})
	}
}

func int64Ptr(v int64) *int64 {
	return &v
}

func compareInt64Ptr(a, b *int64) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func valStr(p *int64) string {
	if p == nil {
		return "nil"
	}
	return fmt.Sprintf("%d", *p)
}
