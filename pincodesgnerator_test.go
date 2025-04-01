package pincodesgenerator

import "testing"

func TestGenerateCode(t *testing.T) {
	tests := []struct {
		name           string
		input          int
		pinCodeLength  int
		expectedOutput PinCode
	}{
		{
			name:           "input 1, length 6",
			input:          1,
			pinCodeLength:  6,
			expectedOutput: PinCode{0, 0, 0, 0, 0, 1},
		},
		{
			name:           "input 13579, length 6",
			input:          13579,
			pinCodeLength:  6,
			expectedOutput: PinCode{0, 1, 3, 5, 7, 9},
		},
		{
			name:           "input 245687, length 6",
			input:          245687,
			pinCodeLength:  6,
			expectedOutput: PinCode{2, 4, 5, 6, 8, 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := PinCodesGenerator{
				options: GeneratorOptions{
					pinCodeLength: allowedPinCodeLengths(tt.pinCodeLength),
				},
			}
			result := gen.generateCode(tt.input)
			for i := range result {
				if result[i] != tt.expectedOutput[i] {
					t.Errorf("generateCode(%d) = %v; want %v", tt.input, result, tt.expectedOutput)
					break
				}
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	type testCase struct {
		name       string
		code       PinCode
		options    GeneratorOptions
		wantResult bool
	}

	tests := []testCase{
		{
			name:       "1234 should be invalid due to incrementing sequence",
			code:       PinCode{1, 2, 3, 4},
			options:    GeneratorOptions{removeRepeats: true, removeIncrements: true},
			wantResult: false,
		},
		{
			name:       "1211 should be invalid due to repeated digits",
			code:       PinCode{1, 2, 1, 1},
			options:    GeneratorOptions{removeRepeats: true, removeIncrements: true},
			wantResult: false,
		},
		{
			name:       "1210 should be valid",
			code:       PinCode{1, 2, 1, 0},
			options:    GeneratorOptions{removeRepeats: true, removeIncrements: true},
			wantResult: true,
		},
		{
			name:       "0010 should be invalid due to repeated digits",
			code:       PinCode{0, 0, 1, 0},
			options:    GeneratorOptions{removeRepeats: true, removeIncrements: true},
			wantResult: false,
		},
		{
			name:       "2123 should be invalid due to increment",
			code:       PinCode{2, 1, 2, 3},
			options:    GeneratorOptions{removeRepeats: true, removeIncrements: true},
			wantResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := PinCodesGenerator{
				options: tt.options,
			}
			got := gen.isValid(tt.code)
			if got != tt.wantResult {
				t.Errorf("isValid(%v) = %v; want %v", tt.code, got, tt.wantResult)
			}
		})
	}
}

// tests that we get the expected number of codes
// and that they are unique
func TestGenerate_UniqueValidFourDigitCodes(t *testing.T) {
	gen := NewPinCodesGenerator()
	gen.DisallowRepeats()
	gen.DisallowIncrements()

	n := 1000
	codes, err := gen.Generate(n)
	if err != nil {
		t.Fatalf("Generate(%d) returned error: %v", n, err)
	}

	if len(codes) != n {
		t.Fatalf("expected %d codes, got %d", n, len(codes))
	}

	// Uniqueness check using a map
	seen := make(map[string]bool)
	for _, code := range codes {
		if seen[code.String()] {
			t.Errorf("duplicate code generated: %v", code)
		}
		seen[code.String()] = true
	}
}

// tests that we don't have invalid codes included in the pool of possible codes
func TestGenerateSetOfAllPossibleCodes_ExcludesInvalidPatterns(t *testing.T) {
	gen := NewPinCodesGenerator()
	gen.DisallowRepeats()
	gen.DisallowIncrements()

	codes := gen.generateSetOfAllPossibleCodes()
	codeSet := make(map[string]struct{})
	for _, code := range codes {
		codeSet[code.String()] = struct{}{}
	}

	// List of known bad codes (as PinCodes)
	badCodes := []PinCode{
		{1, 2, 3, 4}, // incrementing sequence
		{0, 1, 2, 3},
		{2, 2, 5, 6}, // repeat
		{4, 5, 6, 7},
		{1, 1, 1, 1}, // all same digit
		{7, 8, 9, 0}, // wrapping increment
		{1, 1, 5, 6}, // example used in directions
		{1, 2, 3, 6}, // example used in directions
	}

	for _, badCode := range badCodes {
		if _, exists := codeSet[badCode.String()]; exists {
			t.Errorf("invalid code was included in output: %v", badCode)
		}
	}
}
