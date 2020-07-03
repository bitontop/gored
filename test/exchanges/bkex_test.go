package test

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"testing"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
	// "../../exchange/bkex"
	// "../conf"
)

/********************Public API********************/
func Test_Bkex(t *testing.T) {
	e := InitEx(exchange.BKEX)
	pair := pair.GetPairByKey("BTC|ETH")

	Test_Coins(e)
	Test_Pairs(e)
	Test_Pair(e, pair)
	// Test_Orderbook(e, pair)
	Test_NewOrderBook(e, pair)
	Test_ConstraintFetch(e, pair)
	Test_Constraint(e, pair)

	// Test_TradeHistory(e, pair)

	// Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.00000001, 100)
	// Test_Trading_Sell(e, pair, 0.06, 0.0015)
	// Test_OrderStatus(e, pair, "2020030605420418411071353")
	// Test_CancelOrder(e, pair, "1234567890")
	// // Test_Withdraw(e, pair.Base, 1.5, "1HB5XMLmzFVj8ALj6mfBsbifRoD4miY36v")

	// Test_CheckBalance(e, pair.Target, exchange.AssetWallet)
	// Test_CheckAllBalance(e, exchange.AssetWallet)
	// Test_DoTransfer(e, pair.Target, "1", exchange.AssetWallet, exchange.SpotWallet)
	// Test_DoWithdraw(e, pair.Target, "1", "0x37E0Fc27C6cDB5035B2a3d0682B4E7C05A4e6C46", "tag")
}
