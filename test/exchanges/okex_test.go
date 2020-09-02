package test

import (
	"testing"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
	// "../../exchange/okex"
	// "../conf"
)

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

/********************Public API********************/

func Test_Okex(t *testing.T) {
	// e := InitEx(exchange.OKEX)
	// pair := pair.GetPairByKey("BTC|ETH")
	e := InitExFromJson(exchange.OKEX)
	pair := pair.GetPairByKey("BTC|EOS")

	// Test_Coins(e)
	// Test_Pairs(e)
	Test_Pair(e, pair)
	// Test_Orderbook(e, pair)
	// Test_NewOrderBook(e, pair)
	// Test_ConstraintFetch(e, pair)
	Test_Constraint(e, pair)
	// Test_TickerPrice(e)

	// new interface methods
	// Test_DoWithdraw(e, pair.Target, "1", "0x37E0Fc27C6cDB5035B2a3d0682B4E7C05A4e6C46", "tag")
	// Test_DoTransfer(e, pair.Target, "10", exchange.AssetWallet, exchange.SpotWallet)
	// Test_CheckBalance(e, pair.Target, exchange.AssetWallet)
	Test_CheckAllBalance(e, exchange.SpotWallet)

	// Test_Trading_Sell(e, pair, 0.001, 0.1)
	// okex.Socket(pair)
	Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.00000001, 100)
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")

	// Test_AOOpenOrder(e, pair)

	// =====================================================================
	// TransferHistory
	// opCTransferHistory := &exchange.AccountOperation{
	// 	Type:   exchange.GetTransferHistory,
	// 	Wallet: exchange.SpotWallet,
	// 	Ex:     e.GetName(),
	// 	// TransferStartTime: 123456,
	// 	// TransferEndTime:   234567,
	// 	DebugMode: true,
	// }

	// if err := e.DoAccountOperation(opCTransferHistory); err != nil {
	// 	log.Printf("%+v", err)
	// } else {
	// 	for _, o := range opCTransferHistory.TransferInHistory {
	// 		log.Printf("%s TransferInHistory %+v", e.GetName(), o)
	// 	}
	// 	for _, o := range opCTransferHistory.TransferOutHistory {
	// 		log.Printf("%s TransferOutHistory %+v", e.GetName(), o)
	// 	}
	// 	log.Printf("Spot TransferHistory CallResponse: %+v", opCTransferHistory.CallResponce)
	// }
	// =====================================================================

	// ==============================================

	// spot Kline
	// interval options: 1min, 5min, 15min, 30min, 1hour, 2hour, 4hour, 6hour, 12hour, 1day, 1week
	// can only get 24h data
	// opKline := &exchange.PublicOperation{
	// 	Wallet:         exchange.SpotWallet,
	// 	Type:           exchange.KLine,
	// 	EX:             e.GetName(),
	// 	Pair:           pair,
	// 	KlineInterval:  "1hour", // default to 5min if not provided
	// 	KlineStartTime: 1592650051001,
	// 	KlineEndTime:   1592660051002,
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

	// SubAccount Transfer
	// opSubTransfer := &exchange.AccountOperation{
	// 	Wallet: exchange.SpotWallet,
	// 	Type:   exchange.SubAccountTransfer,
	// 	Ex:     e.GetName(),
	// 	Coin:   pair.Target,
	// 	// SubTransferFrom: "bitontop", // ** Put subAccountId into 'SubTransferFrom' or 'SubTransferTo' **
	// 	SubTransferTo:     "bitontop",
	// 	SubTransferAmount: "0.001",
	// 	DebugMode:         true,
	// }
	// err := e.DoAccountOperation(opSubTransfer)
	// if err != nil {
	// 	log.Printf("SubAccount Transfer error: %v", err)
	// }
	// log.Printf("SubAccount Transfer callResponse: %+v", opSubTransfer.CallResponce)

	// ===============================================

	// Test_TradeHistory(e, pair)
}
