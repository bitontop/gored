package test

import (
	"testing"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
	// "../exchange/ftx"
	// "./conf"
)

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

/********************Public API********************/

func Test_Ftx(t *testing.T) {
	e := InitEx(exchange.FTX)
	// e := InitExFromJson(exchange.FTX)
	pair := pair.GetPairByKey("USD|BTC")

	// Test_Coins(e)
	// Test_Pairs(e)
	Test_Pair(e, pair)
	// Test_Orderbook(e, pair)
	// Test_ConstraintFetch(e, pair)
	Test_Constraint(e, pair)
	// Test_TickerPrice(e)

	// Test_AOOpenOrder(e, nil)
	// Test_AOOrderHistory(e, pair)
	// Test_AODepositAddress(e, pair.Target)
	// Test_AODepositHistory(e, pair)
	// Test_AOWithdrawalHistory(e, pair)
	// Test_AOTransferHistory(e)

	// ==============================================
	// spot Kline
	// interval options: 15s, 1min, 5min, 15min, 1hour, 4hour, 1day
	// opKline := &exchange.PublicOperation{
	// 	Wallet:         exchange.SpotWallet,
	// 	Type:           exchange.KLine,
	// 	EX:             e.GetName(),
	// 	Pair:           pair,
	// 	KlineInterval:  "1min", // default to 5min if not provided
	// 	KlineStartTime: 1593400000123,
	// 	KlineEndTime:   1593500000123,
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

	// private only balance done, POST auth not passed
	Test_Balance(e, pair)
	Test_CheckBalance(e, pair.Base, exchange.SpotWallet)
	Test_CheckAllBalance(e, exchange.SpotWallet)
	// Test_Trading(e, pair, 1234, 0.001)
	// Test_Trading_Sell(e, pair, 999999, 0.1)
	// Test_OrderStatus(e, pair, "9596912")
	// Test_CancelOrder(e, pair, "9596912")
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")
	// log.Println(e.GetTradingWebURL(pair))
}
