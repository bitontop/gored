package test

import (
	"testing"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
	// "../../exchange/poloniex"
	// "../conf"
)

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

/********************Public API********************/

func Test_Poloniex(t *testing.T) {
	e := InitEx(exchange.POLONIEX)
	pair := pair.GetPairByKey("BTC|ETH")

	// Test_Coins(e)
	// Test_Pairs(e)
	Test_Pair(e, pair)
	// Test_Orderbook(e, pair)
	// Test_ConstraintFetch(e, pair)
	// Test_Constraint(e, pair)

	Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.00000001, 100)
	// Test_Trading_Sell(e, pair, 0.05, 0.01)
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")

	// Test_TradeHistory(e, pair)
	// Test_NewOrderBook(e, pair)

	// ==============================================

	// spot Kline
	// interval options: 5min, 15min, 30min, 2hour, 4hour, 1day
	// opKline := &exchange.PublicOperation{
	// 	Wallet:         exchange.SpotWallet,
	// 	Type:           exchange.KLine,
	// 	EX:             e.GetName(),
	// 	Pair:           pair,
	// 	KlineInterval:  "15min", // default to 5min if not provided
	// 	KlineStartTime: 1530965420000,
	// 	KlineEndTime:   1530969020000,
	// 	DebugMode:      true,
	// }
	// err := e.LoadPublicData(opKline)
	// if err != nil {
	// 	log.Printf("%v", err)
	// }

	// for _, k := range opKline.Kline {
	// 	log.Printf("%s SpotKline %+v", e.GetName(), k)
	// }
	// ==============================================
}
