package test

import (
	"log"
	"testing"

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
	// e := InitExFromJson(exchange.HUOBI)
	// pair := pair.GetPairByKey("BTC|BSV")
	if pair == nil {
		log.Printf("got nil pair: %+v", pair)
	} else {
		log.Printf("got pair: %+v", pair)
	}

	// Test_CoinChainType(e, pair.Base)
	// Test_TradeHistory(e, pair)

	// Test_Coins(e)
	// Test_Pairs(e)
	Test_Pair(e, pair)
	// Test_Orderbook(e, pair)
	// Test_NewOrderBook(e, pair)
	// Test_ConstraintFetch(e, pair)
	Test_Constraint(e, pair)
	// Test_TickerPrice(e)

	Test_Balance(e, pair)
	// Test_CheckAllBalance(e, exchange.SpotWallet)
	// Test_Trading(e, pair, 0.00000001, 100)
	// Test_Trading_Sell(e, pair, 0.05, 0.01)
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")

	// ==============================================

	// spot Kline
	// interval options: 1min, 5min, 15min, 30min, 1hour, 4hour, 1day, 1mon, 1week, 1year
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
	// =====================================================================
	// TransferHistory
	// opCTransferHistory := &exchange.AccountOperation{
	// 	Type:      exchange.GetTransferHistory,
	// 	Wallet:    exchange.SpotWallet,
	// 	Ex:        e.GetName(),
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

	// SubAccount Transfer
	// opSubTransfer := &exchange.AccountOperation{
	// 	Wallet:          exchange.SpotWallet,
	// 	Type:            exchange.SubAccountTransfer,
	// 	Ex:              e.GetName(),
	// 	Coin:            pair.Target,
	// 	SubTransferFrom: "157709010", // ** Put subAccountId into 'SubTransferFrom' or 'SubTransferTo' **
	// 	// SubTransferTo:     "157709010",
	// 	SubTransferAmount: "0.00599",
	// 	DebugMode:         true,
	// }
	// err := e.DoAccountOperation(opSubTransfer)
	// if err != nil {
	// 	log.Printf("SubAccount Transfer error: %v", err)
	// }
	// log.Printf("SubAccount Transfer callResponse: %+v", opSubTransfer.CallResponce)

	// ===============================================

	// SubBalances(e, "8459451")
	// SubAccountList(e)
	// SubAllBalances(e)

	// Test Withdraw
	// Test_DoWithdraw(e, pair.Target, "1", "0x37E0Fc27C6cDB5035B2a3d0682B4E7C05A4e6C46", "tag")

	// Test_TradeHistory(e, pair)

	// Test_AOOpenOrder(e, pair)
	// time.Sleep(time.Second * 5)
	// Test_AOOrderHistory(e, pair)
	// time.Sleep(time.Second * 5)
	// Test_AODepositAddress(e, pair.Base)
	// time.Sleep(time.Second * 5)
	// Test_AODepositHistory(e, pair)
	// time.Sleep(time.Second * 5)
	// Test_AOWithdrawalHistory(e, pair)
}
