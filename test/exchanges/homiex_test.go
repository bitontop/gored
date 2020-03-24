package test

import (
	"testing"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
	// "../../exchange/homiex"
	// "../conf"
)

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

/********************Public API********************/

func Test_Homiex(t *testing.T) {
	e := InitEx(exchange.HOMIEX)
	pair := pair.GetPairByKey("BTC|ETH") // USDT|VBCC

	// Test_TradeHistory(e, pair)

	// Test_Coins(e)
	// Test_Pairs(e)
	Test_Pair(e, pair)
	// Test_Orderbook(e, pair)
	// Test_ConstraintFetch(e, pair)
	// Test_Constraint(e, pair)

	// Test_Balance(e, pair)
	// Test_Trading_Sell(e, pair, 0.06, 0.01)
	// Test_Trading(e, pair, 0.0001, 0.1)
	// Test_OrderStatus(e, pair, "1234567890")
	// Test_CancelOrder(e, pair, "539389524336195328")
	// log.Println(e.GetTradingWebURL(pair))

	Test_CheckBalance(e, pair.Target, exchange.AssetWallet)
	Test_CheckAllBalance(e, exchange.SpotWallet)
	// Test_DoWithdraw(e, pair.Target, "0.075", "0xd3ceb35d6fa3dcc11cf7ea70f2d3bdf141b1e82f", "tag")

}
