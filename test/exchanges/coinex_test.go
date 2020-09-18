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
	Test_Constraint(e, pair)

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

	// Test_AOOpenOrder(e, pair)
	// Test_AOWithdrawalHistory(e, pair)
	// Test_AODepositHistory(e, pair)
	// Test_AOTransferHistory(e)

	// SubBalances(e, "test_sub2")
	// SubAllBalances(e)

	// Test_TradeHistory(e, pair)
	// Test_CoinChainType(e, pair.Base)

	// ==============================================

	// spot Kline
	// interval options: 1min, 3min, 5min, 15min, 30min, 1hour, 2hour, 4hour, 6hour, 12hour, 1day, 3day, 1week
	// opKline := &exchange.PublicOperation{
	// 	Wallet:        exchange.SpotWallet,
	// 	Type:          exchange.KLine,
	// 	EX:            e.GetName(),
	// 	Pair:          pair,
	// 	KlineInterval: "1min", // default to 5min if not provided
	// 	DebugMode:     true,
	// }
	// err := e.LoadPublicData(opKline)
	// if err != nil {
	// 	log.Printf("%v", err)
	// }

	// for _, k := range opKline.Kline {
	// 	log.Printf("%s SpotKline %+v", e.GetName(), k)
	// }
	// ==============================================

	// SubAccount Transfer
	// opSubTransfer := &exchange.AccountOperation{
	// 	Wallet: exchange.SpotWallet,
	// 	Type:   exchange.SubAccountTransfer,
	// 	Ex:     e.GetName(),
	// 	Coin:   pair.Target,
	// 	// SubTransferFrom: "sub1", // ** Put subAccountId into 'SubTransferFrom' or 'SubTransferTo' **
	// 	SubTransferTo:     "sub1",
	// 	SubTransferAmount: "0.0001",
	// 	DebugMode:         true,
	// }
	// err := e.DoAccountOperation(opSubTransfer)
	// if err != nil {
	// 	log.Printf("SubAccount Transfer error: %v", err)
	// }
	// log.Printf("SubAccount Transfer callResponse: %+v", opSubTransfer.CallResponce)

	// ===============================================
}
