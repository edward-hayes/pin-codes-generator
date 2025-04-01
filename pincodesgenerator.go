package pincodesgenerator

import (
	"math"
	"math/rand/v2"
	"strconv"
)

type PinCodesGenerator struct {
	options GeneratorOptions
}

type GeneratorOptions struct {
	removeRepeats    bool                  // Allow/Disallow repeating digits in the pin code
	removeIncrements bool                  // Allow/Disallow incrementing digits in the pin code
	pinCodeLength    allowedPinCodeLengths // Length of the pin code (e.g., FourDigits or SixDigits)
}

type PinCode []int

type allowedPinCodeLengths int

const (
	FourDigits allowedPinCodeLengths = 4
	SixDigits  allowedPinCodeLengths = 6
)

func NewPinCodesGenerator() *PinCodesGenerator {
	return &PinCodesGenerator{
		options: GeneratorOptions{
			pinCodeLength: FourDigits, // Default pin code length
		},
	}
}

// determines if the pin code generator should allow repeating digits
func (g *PinCodesGenerator) DisallowRepeats() {
	g.options.removeRepeats = true
}

// determines if the pin code generator should allow incrementing digits
func (g *PinCodesGenerator) DisallowIncrements() {
	g.options.removeIncrements = true
}

func (g *PinCodesGenerator) SetPinCodeToSix() {
	g.options.pinCodeLength = SixDigits
}

// Generate generates a list of n pin codes
func (g *PinCodesGenerator) Generate(n int) ([]PinCode, error) {
	if n <= 0 {
		return nil, ErrNumberOfPinCodesZero
	}

	allPossibleCodes := g.generateSetOfAllPossibleCodes()
	if len(allPossibleCodes) < n {
		return nil, ErrNumberOfPinCodesExceedsMax
	}

	return g.chooseNCodes(n, allPossibleCodes), nil
}

// makes the full set of possible codes from which we can randomly select
// this strategy makes more sense then trying to fill a bucket
// of a random, unique codes. As n goes gets close to the max possible codes
// the time to generate a random code increases
func (g *PinCodesGenerator) generateSetOfAllPossibleCodes() []PinCode {
	codes := make([]PinCode, 0)
	max := int(math.Pow10(int(g.options.pinCodeLength))) // 10^n where n is the pin code length
	for i := range max {
		code := g.generateCode(i)
		if g.isValid(code) {
			codes = append(codes, code)
		}
	}
	return codes
}

func (g *PinCodesGenerator) generateCode(i int) PinCode {
	code := make(PinCode, g.options.pinCodeLength)
	for j := g.options.pinCodeLength - 1; j >= 0; j-- {
		code[j] = i % 10
		i /= 10
	}
	return code
}

func (g *PinCodesGenerator) isValid(code PinCode) bool {
	if g.options.removeRepeats {
		for i := 1; i < len(code); i++ {
			if code[i] == code[i-1] {
				return false
			}
		}
	}

	if g.options.removeIncrements {
		for i := 2; i < len(code); i++ {
			if code[i] == code[i-1]+1 && code[i-1] == code[i-2]+1 {
				return false
			}
		}
	}
	return true
}

func (g *PinCodesGenerator) chooseNCodes(n int, codes []PinCode) []PinCode {
	g.shuffleCodes(codes)
	return codes[:n]
}

// shuffles the codes in place using the Fisher-Yates algorithm
// took this from the examples in the Go docs https://pkg.go.dev/math/rand#Shuffle
func (g *PinCodesGenerator) shuffleCodes(codes []PinCode) {
	rand.Shuffle(len(codes), func(i, j int) {
		codes[i], codes[j] = codes[j], codes[i]
	})
}

func (p *PinCode) String() string {
	s := ""
	for _, digit := range *p {
		s += strconv.Itoa(digit)
	}
	return s
}
