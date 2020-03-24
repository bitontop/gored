package test

import (
	"testing"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
	// "../../exchange/txbit"
	// "../conf"
)

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

/********************Public API********************/

func Test_Txbit(t *testing.T) {
	e := InitEx(exchange.TXBIT)
	pair := pair.GetPairByKey("BTC|AIB")

	// log.Printf("coin.LastID: %v", coin.LastID)
	// Test_Coins(e)
	// Test_Pairs(e)
	Test_Pair(e, pair)
	// Test_Orderbook(e, pair)
	// Test_ConstraintFetch(e, pair)
	// Test_Constraint(e, pair)

	// Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.001, 0.033)
	// Test_Trading_Sell(e, pair, 0.050012345678, 0.10001234567)
	// Test_OrderStatus(e, pair, "b5d7d18c-61fb-479e-8ee6-b222ced93e56")
	// Test_CancelOrder(e, pair, "b5d7d18c-61fb-479e-8ee6-b222ced93e56")
	// Test_Withdraw(e, pair.Target, 0.1, "0xf252be0c7758094a37bf10a4cbf4dec0d69b7bcc")
	// log.Println(e.GetTradingWebURL(pair))

	// // Test Withdraw
	// opWithdraw := &exchange.AccountOperation{
	// 	Type:            exchange.Withdraw,
	// 	Coin:            pair.Target,
	// 	WithdrawAmount:  "1",
	// 	WithdrawAddress: "addr",
	// 	DebugMode:       true,
	// }
	// err := e.DoAccoutOperation(opWithdraw)
	// if err != nil {
	// 	log.Printf("%v", err)
	// }
	// log.Printf("WithdrawID: %v, err: %v", opWithdraw.WithdrawID, opWithdraw.Error)

	Test_TradeHistory(e, pair)
	Test_NewOrderBook(e, pair)
}
