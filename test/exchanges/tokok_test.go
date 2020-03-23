package test

import (
	"testing"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
	// "../exchange/tokok"
	// "./conf"
)

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

/********************Public API********************/

func Test_Tokok(t *testing.T) {
	e := InitEx(exchange.TOKOK)
	pair := pair.GetPairByKey("BTC|ETH")

	Test_Coins(e)
	Test_Pairs(e)
	Test_Pair(e, pair)
	Test_Orderbook(e, pair)
	Test_ConstraintFetch(e, pair)
	Test_Constraint(e, pair)

	Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.013, 0.01)
	// Test_Trading_Sell(e, pair, 0.029, 0.01)
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")
	// Test_CancelOrder(e, pair, "1907200026337682229018")
	// Test_Balance(e, pair)
}
