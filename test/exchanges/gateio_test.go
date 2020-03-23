package test

import (
	"testing"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
	// "../../exchange/gateio"
	// "../conf"
)

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

/********************Public API********************/

func Test_Gateio(t *testing.T) {
	e := InitEx(exchange.GATEIO)
	pair := pair.GetPairByKey("BTC|ETH")

	Test_Coins(e)
	Test_Pairs(e)
	Test_Pair(e, pair)
	// Test_Orderbook(e, pair)
	Test_NewOrderBook(e, pair)
	Test_ConstraintFetch(e, pair)
	Test_Constraint(e, pair)
	// log.Println(e.GetTradingWebURL(pair))

	// Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.00000001, 100)
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")

	// Test_DoWithdraw(e, pair.Target, "0.2", "0x2d1a6a1d65ae08502a5e0ddda0be8df9874f7c14", "tag")

	// Test_TradeHistory(e, pair)
}
