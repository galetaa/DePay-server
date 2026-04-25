package validation

import (
	"errors"
	"math/big"
	"regexp"
)

var evmAddressPattern = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)

func EVMAddress(address string) error {
	if !evmAddressPattern.MatchString(address) {
		return errors.New("invalid address format")
	}
	return nil
}

func PositiveAmount(amount string) error {
	value, ok := new(big.Rat).SetString(amount)
	if !ok {
		return errors.New("invalid amount")
	}
	if value.Sign() <= 0 {
		return errors.New("amount must be greater than zero")
	}
	return nil
}
