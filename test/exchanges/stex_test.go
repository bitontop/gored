package test

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"testing"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
	// "../../exchange/stex"
	// "../conf"
)

/********************Public API********************/
func Test_Stex(t *testing.T) {
	e := InitEx(exchange.STEX)
	pair := pair.GetPairByKey("BTC|ETH") // "ETH|AIB"

	// Test_Coins(e)
	// Test_Pairs(e)
	Test_Pair(e, pair)
	// Test_Orderbook(e, pair)
	// Test_ConstraintFetch(e, pair)
	// Test_Constraint(e, pair)

	Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.00000001, 100)
	// Test_Trading_Sell(e, pair, 0.06, 0.01)
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")

	// Test_AOOpenOrder(e, pair)
	// Test_TradeHistory(e, pair)
	// Test_NewOrderBook(e, pair)

	// // Test Withdraw
	// opWithdraw := &exchange.AccountOperation{
	// 	Type:            exchange.Withdraw,
	// 	Coin:            pair.Target,
	// 	WithdrawAmount:  "1",
	// 	WithdrawAddress: "addr",
	// 	DebugMode:       true,
	// }
	// err := e.DoAccountOperation(opWithdraw)
	// if err != nil {
	// 	log.Printf("%v", err)
	// }
	// log.Printf("WithdrawID: %v, err: %v", opWithdraw.WithdrawID, opWithdraw.Error)
}
