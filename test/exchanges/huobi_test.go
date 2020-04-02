package test

import (
	"testing"
	"time"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
	// "../../exchange/huobi"
	// "../conf"
)

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

/********************Public API********************/

func Test_Huobi(t *testing.T) {
	e := InitEx(exchange.HUOBI)
	pair := pair.GetPairByKey("BTC|ETH")

	// Test_CoinChainType(e, pair.Base)
	// Test_TradeHistory(e, pair)

	// Test_Coins(e)
	// Test_Pairs(e)
	// Test_Pair(e, pair)
	// Test_Orderbook(e, pair)
	// Test_NewOrderBook(e, pair)
	// Test_ConstraintFetch(e, pair)
	// Test_Constraint(e, pair)

	// Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.00000001, 100)
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")

	// Test Withdraw
	// Test_DoWithdraw(e, pair.Target, "1", "0x37E0Fc27C6cDB5035B2a3d0682B4E7C05A4e6C46", "tag")

	// Test_TradeHistory(e, pair)

	// Test_AOOpenOrder(e, pair)
	// time.Sleep(time.Second * 5)
	// Test_AOOrderHistory(e, pair)
	// time.Sleep(time.Second * 5)
	Test_AODepositAddress(e, pair.Base)
	time.Sleep(time.Second * 5)
	Test_AODepositHistory(e, pair)
	time.Sleep(time.Second * 5)
	Test_AOWithdrawalHistory(e, pair)
}
