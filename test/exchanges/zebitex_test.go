package test

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"testing"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
	// "../../exchange/zebitex"
	// "../conf"
)

/********************Public API********************/
func Test_Zebitex(t *testing.T) {
	e := InitEx(exchange.ZEBITEX)
	pair := pair.GetPairByKey("BTC|ETH")

	// Test_Coins(e)
	// Test_Pairs(e)
	// Test_Pair(e, pair)
	// Test_Orderbook(e, pair)
	// Test_ConstraintFetch(e, pair)
	// Test_Constraint(e, pair)

	// Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.001, 100)
	// Test_Trading_Sell(e, pair, 0.05, 0.001)
	// Test_OrderStatus(e, pair, "88410563")
	// Test_CancelOrder(e, pair, "88408135")
	// Test_Withdraw(e, pair.Target, 0.001, "0xaC05f7b683b14e5997d288a8C031c5143533F9e3") // 1HB5XMLmzFVj8ALj6mfBsbifRoD4miY36v
	// self deposit "0x13e1689aff770c5d011799925259df4a10198020"

	// Test_CheckBalance(e, pair.Target, exchange.AssetWallet)
	// Test_CheckAllBalance(e, exchange.SpotWallet)
	// Test_DoWithdraw(e, pair.Target, "1", "0xaC05f7b683b14e5997d288a8C031c5143533F9e3", "tag")

	Test_TradeHistory(e, pair)
	Test_NewOrderBook(e, pair)
}
