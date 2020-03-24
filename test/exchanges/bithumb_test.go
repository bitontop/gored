package test

import (
	"testing"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
	// "../../exchange/bithumb"
	// "../conf"
)

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

/********************Public API********************/

func Test_Bithumb(t *testing.T) {
	e := InitEx(exchange.BITHUMB)
	pair := pair.GetPairByKey("BTC|ETH")

	// Test_Coins(e)
	// Test_Pairs(e)
	Test_Pair(e, pair)
	// Test_Orderbook(e, pair)
	// Test_ConstraintFetch(e, pair)
	// Test_Constraint(e, pair)

	// Test_Balance(e, pair)
	// Test_Trading_Sell(e, pair, 0.1, 100)
	// Test_Trading(e, pair, 0.01, 100)
	// Test_OrderStatus(e, pair, "1234567890")
	// Test_CancelOrder(e, pair, "1234567890")
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")
	// log.Println(e.GetTradingWebURL(pair))

	// Test_CheckBalance(e, pair.Target, exchange.AssetWallet)
	// Test_CheckAllBalance(e, exchange.SpotWallet)
	// Test_DoTransfer(e, pair.Target, "0.1", exchange.SpotWallet, exchange.AssetWallet)
	// Test_DoWithdraw(e, pair.Target, "1", "0x37E0Fc27C6cDB5035B2a3d0682B4E7C05A4e6C46", "tag")

	Test_TradeHistory(e, pair)
}
