package test

import (
	"testing"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
	// "../../exchange/hoo"
	// "../conf"
)

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

/********************Public API********************/

func Test_Hoo(t *testing.T) {
	e := InitEx(exchange.HOO)
	pair := pair.GetPairByKey("BTC|ETH")

	// Test_Coins(e)
	// Test_Pairs(e)
	Test_Pair(e, pair)
	// Test_Orderbook(e, pair)
	// Test_ConstraintFetch(e, pair)
	// Test_Constraint(e, pair)

	Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.0001, 100)
	// Test_Trading_Sell(e, pair, 0.06, 0.00001)

	// ****************** Cancel Order
	// order := &exchange.Order{
	// 	Pair:         pair,
	// 	OrderID:      "11579192813104769",
	// 	Rate:         0.001,
	// 	Quantity:     10,
	// 	Side:         "Buy",
	// 	Status:       exchange.New,
	// 	CancelStatus: "500788212430320131636",
	// }

	// err := e.CancelOrder(order)
	// if err == nil {
	// 	log.Printf("%s Cancel Order: %v", e.GetName(), order)
	// } else {
	// 	log.Printf("%s Cancel Order Err: %s", e.GetName(), err)
	// }
	// ******************

	// Test_OrderStatus(e, pair, "11579192813472132")
	// Test_CancelOrder(e, pair, "11579192813104769")
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")
	// log.Println(e.GetTradingWebURL(pair))

	Test_CheckBalance(e, pair.Target, exchange.AssetWallet)
	Test_CheckAllBalance(e, exchange.SpotWallet)
	// Test_DoTransfer(e, pair.Target, "1", exchange.AssetWallet, exchange.SpotWallet)
	// Test_DoWithdraw(e, pair.Target, "1", "0x37E0Fc27C6cDB5035B2a3d0682B4E7C05A4e6C46", "tag")

}
