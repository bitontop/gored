package oksim

import (
	"fmt"

	"github.com/bitontop/gored/exchange"
)

func (e *Oksim) DoAccountOperation(operation *exchange.AccountOperation) error {
	return fmt.Errorf("%s Operation type invalid: %s %v", operation.Ex, operation.Wallet, operation.Type)
}
