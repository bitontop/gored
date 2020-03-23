package test

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"testing"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
	// "../../exchange/bitfinex"
	// "../conf"
)

/********************Public API********************/
func Test_Bitfinex(t *testing.T) {
	e := InitEx(exchange.BITFINEX)
	pair := pair.GetPairByKey("BTC|ETH")

	// Test_Coins(e)
	// Test_Pairs(e)
	Test_Pair(e, pair)
	// Test_Orderbook(e, pair)
	// Test_ConstraintFetch(e, pair)
	// Test_Constraint(e, pair)

	Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.0001000123, 100)
	// Test_Trading_Sell(e, pair, 0.06123456789, 0.02)
	// Test_OrderStatus(e, pair, "40797839810") // 40797839810 an open order
	// Test_CancelOrder(e, pair, "41094576119")
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")
	// log.Printf("Url: %v", e.GetTradingWebURL(pair))
	// Test_DoWithdraw(e, pair.Target, "1", "0x37E0Fc27C6cDB5035B2a3d0682B4E7C05A4e6C46", "tag")

	Test_TradeHistory(e, pair)
}
