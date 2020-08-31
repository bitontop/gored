package test

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"log"
	"testing"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
	// "../../exchange/binance"
	// "../conf"
)

/********************Public API********************/
func Test_Binance(t *testing.T) {
	// e := InitEx(exchange.BINANCE)
	e := InitExFromJson(exchange.BINANCE)
	// e := InitExFromJson(exchange.BINANCE)
	pair := pair.GetPairByKey("BTC|ETH")

	// Test_Coins(e)
	// Test_Pairs(e)
	Test_Pair(e, pair)
	// Test_Orderbook(e, pair)
	// Test_ConstraintFetch(e, pair)
	Test_Constraint(e, pair)

	// Test_NewOrderBook(e, pair)
	// Test_TickerPrice(e)

	Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.01, 0.01)
	// Test_Trading_Sell(e, pair, 0.04, 0.01)
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")
	// Test_DoWithdraw(e, pair.Target, "1", "0x37E0Fc27C6cDB5035B2a3d0682B4E7C05A4e6C46", "tag")

	// Test_AOOpenOrder(e, pair)
	// Test_AOOrderHistory(e, pair) // not tested with asset
	// Test_AODepositAddress(e, pair.Base)
	// Test_AODepositHistory(e, pair)
	// Test_AOWithdrawalHistory(e, pair) // not tested with asset

	// var err error

	// SubBalances(e, "example@bitontop.com")
	// SubAccountList(e)
	// SubAllBalances(e)

	// Spot AllBalance
	// opAllBalance := &exchange.AccountOperation{
	// 	Wallet:    exchange.SpotWallet,
	// 	Type:      exchange.BalanceList,
	// 	Ex:        e.GetName(),
	// 	DebugMode: true,
	// }
	// err = e.DoAccountOperation(opAllBalance)
	// if err != nil {
	// 	log.Printf("==%v", err)
	// }
	// for _, balance := range opAllBalance.BalanceList {
	// 	log.Printf("AllAccount balance: Coin: %v, avaliable: %v, frozen: %v", balance.Coin.Code, balance.BalanceAvailable, balance.BalanceFrozen)
	// }
	// if len(opAllBalance.BalanceList) == 0 {
	// 	log.Println("AllAccount 0 balance")
	// }
	// log.Printf("AllAccount done")

	// ======================================================
	// contract orderbook
	// opOrderBook := &exchange.PublicOperation{
	// 	Wallet: exchange.ContractWallet,
	// 	Type:          exchange.Orderbook,
	// 	EX:            e.GetName(),
	// 	Pair:          pair,
	// 	DebugMode:     true,
	// }
	// err = e.LoadPublicData(opOrderBook)
	// if err != nil {
	// 	log.Printf("%v", err)
	// }

	// log.Printf("%s ContractOrderBook %+v", e.GetName(), opOrderBook.Maker)

	// ==============================================

	// Contract OpenOrder
	// opCOpen := &exchange.AccountOperation{
	// 	Type:      exchange.GetOpenOrder,
	// 	Wallet:    exchange.ContractWallet,
	// 	Ex:        e.GetName(),
	// 	Pair:      pair,
	// 	DebugMode: true,
	// }

	// if err := e.DoAccountOperation(opCOpen); err != nil {
	// 	log.Printf("error: %+v", err)
	// } else {
	// 	for _, o := range opCOpen.OpenOrders {
	// 		log.Printf("%s OpenOrders %+v", e.GetName(), o)
	// 	}
	// 	log.Printf("callResponse: %+v", opCOpen.CallResponce)
	// }
	// ==============================================
	// Contract OrderHistory
	// opCOrderHistory := &exchange.AccountOperation{
	// 	Type:      exchange.GetOrderHistory,
	// 	Wallet:    exchange.ContractWallet,
	// 	Ex:        e.GetName(),
	// 	Pair:      pair,
	// 	DebugMode: true,
	// }

	// if err := e.DoAccountOperation(opCOrderHistory); err != nil {
	// 	log.Printf("%+v", err)
	// } else {
	// 	for _, o := range opCOrderHistory.OrderHistory {
	// 		log.Printf("%s OrderHistory %+v", e.GetName(), o)
	// 	}
	// 	log.Printf("Contract OrderHistory CallResponse: %+v", opCOrderHistory.CallResponce)
	// }
	// ==============================================
	// Contract TransferHistory
	// opCTransferHistory := &exchange.AccountOperation{
	// 	Type:              exchange.GetTransferHistory,
	// 	Wallet:            exchange.ContractWallet,
	// 	Ex:                e.GetName(),
	// 	Coin:              pair.Base,
	// 	TransferStartTime: 1555056425000,
	// 	DebugMode:         true,
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
	// 	log.Printf("Contract TransferHistory CallResponse: %+v", opCTransferHistory.CallResponce)
	// }
	// ==============================================

	// Contract Position Information
	// opCPositionInfo := &exchange.AccountOperation{
	// 	Type:      exchange.GetPositionInfo,
	// 	Wallet:    exchange.ContractWallet,
	// 	Ex:        e.GetName(),
	// 	DebugMode: true,
	// }

	// if err := e.DoAccountOperation(opCPositionInfo); err != nil {
	// 	log.Printf("%+v", err)
	// } else {
	// 	log.Printf("Contract Position Information CallResponse: %+v", opCPositionInfo.CallResponce)
	// }
	// ==============================================

	// spot Kline
	// interval options: 1min, 5min, 15min, 30min, 1hour, 2hour, 4hour, 6hour, 8hour, 12hour, 1day, 3day, 1week, 1month
	// opKline := &exchange.PublicOperation{
	// 	Wallet:         exchange.SpotWallet,
	// 	Type:           exchange.KLine,
	// 	EX:             e.GetName(),
	// 	Pair:           pair,
	// 	KlineInterval:  "1min", // default to 5min if not provided
	// 	KlineStartTime: 1530965420000,
	// 	KlineEndTime:   1530966020000,
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

	// contract AllBalance
	// opContractAllBalance := &exchange.AccountOperation{
	// 	Wallet:    exchange.ContractWallet,
	// 	Type:      exchange.BalanceList,
	// 	Ex:        e.GetName(),
	// 	DebugMode: true,
	// }
	// err = e.DoAccountOperation(opContractAllBalance)
	// if err != nil {
	// 	log.Printf("==%v", err)
	// }
	// for _, balance := range opContractAllBalance.BalanceList {
	// 	log.Printf("AllAccount balance: Coin: %v, avaliable: %v, frozen: %v", balance.Coin.Code, balance.BalanceAvailable, balance.BalanceFrozen)
	// }
	// if len(opContractAllBalance.BalanceList) == 0 {
	// 	log.Println("AllAccount 0 balance")
	// }
	// ==============================================

	// OrderType: GTC, IOC, FOK, GTX
	// TradeType: TRADE_LIMIT, TRADE_MARKET, Trade_STOP_LIMIT, Trade_STOP_MARKET
	// Stop order need 'StopRate' param
	// Contract PlaceOrder
	// opPlaceOrder := &exchange.AccountOperation{
	// 	Wallet:         exchange.ContractWallet,
	// 	Type:           exchange.PlaceOrder,
	// 	Ex:             e.GetName(),
	// 	Pair:           pair,
	// 	OrderDirection: exchange.Buy,
	// 	TradeType:      exchange.TRADE_LIMIT,
	// 	OrderType:      exchange.GTC,
	// 	Rate:           3500,
	// 	Quantity:       0.01,
	// 	DebugMode:      true,
	// }
	// if err := e.DoAccountOperation(opPlaceOrder); err != nil {
	// 	log.Printf("==%v", err)
	// }
	// log.Printf("Contract PlaceOrder CallResponse: %+v", opPlaceOrder.CallResponce)

	// ===============================================

	// contract CancelOrder
	// cancelOrder := &exchange.Order{
	// 	Pair:    pair,
	// 	OrderID: "6257716142",
	// }
	// opCancelOrder := &exchange.AccountOperation{
	// 	Wallet:    exchange.ContractWallet,
	// 	Type:      exchange.CancelOrder,
	// 	Ex:        e.GetName(),
	// 	Pair:      pair,
	// 	Order:     cancelOrder,
	// 	DebugMode: true,
	// }
	// if err := e.DoAccountOperation(opCancelOrder); err != nil {
	// 	log.Printf("==%v", err)
	// } else {
	// 	log.Printf("Contract opCancelOrder callResponse: %+v", opCancelOrder.CallResponce)
	// }
	// ==============================================

	// contract OrderStatus
	// order := &exchange.Order{
	// 	Pair:    pair,
	// 	OrderID: "6257716142",
	// }
	// opOrderStatus := &exchange.AccountOperation{
	// 	Wallet:    exchange.ContractWallet,
	// 	Type:      exchange.GetOrderStatus,
	// 	Ex:        e.GetName(),
	// 	Pair:      pair,
	// 	Order:     order,
	// 	DebugMode: true,
	// }
	// err := e.DoAccountOperation(opOrderStatus)
	// if err != nil {
	// 	log.Printf("==%v", err)
	// } else {
	// 	log.Printf("Contract OrderStatus: %+v", opOrderStatus.Order)
	// 	log.Printf("Contract OrderStatus callResponse: %+v", opOrderStatus.CallResponce)
	// }

	// // ===============================================

	// SubAccount Transfer
	// opSubTransfer := &exchange.AccountOperation{
	// 	Wallet:            exchange.SpotWallet,
	// 	Type:              exchange.SubAccountTransfer,
	// 	Ex:                e.GetName(),
	// 	Coin:              pair.Target,
	// 	SubTransferFrom:   "exapi@bitontop.com",
	// 	SubTransferTo:     "tonywei@bitontop.com",
	// 	SubTransferAmount: "0.01",
	// 	DebugMode:         true,
	// }
	// err := e.DoAccountOperation(opSubTransfer)
	// if err != nil {
	// 	log.Printf("SubAccount Transfer error: %v", err)
	// }
	// log.Printf("SubAccount Transfer callResponse: %+v", opSubTransfer.CallResponce)

	// // ===============================================

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
	// time.Sleep(time.Second * 5)
	// Test_AOTransferHistory(e)
}

func Withdraw(e exchange.Exchange, pair *pair.Pair) {
	opWithdraw := &exchange.AccountOperation{
		Wallet:          exchange.SpotWallet,
		Type:            exchange.Withdraw,
		Coin:            pair.Target,
		WithdrawAddress: "35263",
		WithdrawAmount:  "0.1",
		Ex:              e.GetName(),
		DebugMode:       true,
	}
	err := e.DoAccountOperation(opWithdraw)
	if err != nil {
		log.Printf("==%v", err)
		return
	}
	log.Printf("Withdraw JSON RESPONSE: %v", opWithdraw.CallResponce)
}
