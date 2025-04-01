package pincodesgenerator

import "errors"

var (
	ErrNumberOfPinCodesZero      = errors.New("number of pin codes must be greater than zero")
	ErrNumberOfPinCodesExceedsMax = errors.New("number of pin codes requested exceeds the maximum possible pin codes")
)
