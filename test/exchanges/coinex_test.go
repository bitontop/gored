package test

import (
	"testing"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
	// "../../exchange/coinex"
	// "../conf"
)

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

/********************Public API********************/

func Test_Coinex(t *testing.T) {
	e := InitEx(exchange.COINEX)
	pair := pair.GetPairByKey("BTC|ETH")

	// Test_Coins(e)
	// Test_Pairs(e)
	Test_Pair(e, pair)
	// Test_Orderbook(e, pair)
	// Test_NewOrderBook(e, pair)
	// Test_ConstraintFetch(e, pair)
	// Test_Constraint(e, pair)

	// Test_TickerPrice(e)

	Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.00000001, 100)
	// Test_Trading_Sell(e, pair, 0.06, 100)
	// Test_OrderStatus(e, pair, "1234567890")
	// Test_CancelOrder(e, pair, "1234567890")
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")
	// log.Println(e.GetTradingWebURL(pair))

	// New Interface
	// Test_CheckBalance(e, pair.Target, exchange.AssetWallet)
	Test_CheckAllBalance(e, exchange.SpotWallet)
	// Test_DoTransfer(e, pair.Target, "1", exchange.AssetWallet, exchange.SpotWallet)
	// Test_DoWithdraw(e, pair.Target, "1", "0x37E0Fc27C6cDB5035B2a3d0682B4E7C05A4e6C46", "tag")

	Test_AOOpenOrder(e, pair)
	Test_AOWithdrawalHistory(e, pair)
	Test_AODepositHistory(e, pair)

	// SubBalances(e, "test_sub2")
	// SubAllBalances(e)

	// Test_TradeHistory(e, pair)
	// Test_CoinChainType(e, pair.Base)
}
