package test

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"log"
	"testing"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"

	"github.com/bitontop/gored/exchange/binance"
	"github.com/bitontop/gored/test/conf"
	// "../../exchange/binance"
	// "../conf"
)

/********************Public API********************/
func Test_Binance(t *testing.T) {
	e := InitBinance()

	pair := pair.GetPairByKey("USDT|BTC")

	// Test_Coins(e)
	// Test_Pairs(e)
	Test_Pair(e, pair)
	// Test_Orderbook(e, pair)
	// Test_ConstraintFetch(e, pair)
	// Test_Constraint(e, pair)

	var err error
	// ==============================================

	// contract orderbook
	// opOrderBook := &exchange.PublicOperation{
	// 	OperationType: exchange.ContractWallet,
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

	// contract PlaceOrder
	opPlaceOrder := &exchange.AccountOperation{
		OperationType:  exchange.ContractWallet,
		Type:           exchange.PlaceOrder,
		Ex:             e.GetName(),
		Pair:           pair,
		OrderDirection: exchange.Sell,
		Rate:           9000,
		Quantity:       0.01,
		DebugMode:      true,
	}
	err = e.DoAccoutOperation(opPlaceOrder)
	if err != nil {
		log.Printf("==%v", err)
	}
	// ==============================================

	// contract OrderStatus
	order := &exchange.Order{
		Pair:    pair,
		OrderID: "1573346959",
	}
	opOrderStatus := &exchange.AccountOperation{
		OperationType: exchange.ContractWallet,
		Type:          exchange.OrderStatusOp,
		Ex:            e.GetName(),
		Pair:          pair,
		Order:         order,
		DebugMode:     true,
	}
	err = e.DoAccoutOperation(opOrderStatus)
	if err != nil {
		log.Printf("==%v", err)
	}
	// ==============================================

	// contract CancelOrder
	order = &exchange.Order{
		Pair:    pair,
		OrderID: "1573346959",
	}
	opCancelOrder := &exchange.AccountOperation{
		OperationType: exchange.ContractWallet,
		Type:          exchange.CancelOrder,
		Ex:            e.GetName(),
		Pair:          pair,
		Order:         order,
		DebugMode:     true,
	}
	err = e.DoAccoutOperation(opCancelOrder)
	if err != nil {
		log.Printf("==%v", err)
	}
	// ==============================================

	// contract AllBalance
	opAllBalance := &exchange.AccountOperation{
		OperationType: exchange.ContractWallet,
		Type:          exchange.BalanceList,
		Ex:            e.GetName(),
		DebugMode:     true,
	}
	err = e.DoAccoutOperation(opAllBalance)
	if err != nil {
		log.Printf("==%v", err)
	}
	for _, balance := range opAllBalance.BalanceList {
		log.Printf("AllAccount balance: Coin: %v, avaliable: %v, frozen: %v", balance.Coin.Code, balance.BalanceAvailable, balance.BalanceFrozen)
	}
	if len(opAllBalance.BalanceList) == 0 {
		log.Println("AllAccount 0 balance")
	}
	// ==============================================

	Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.00000001, 100)
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")
	// Test_DoWithdraw(e, pair.Target, "1", "0x37E0Fc27C6cDB5035B2a3d0682B4E7C05A4e6C46", "tag")

	// Test_TradeHistory(e, pair)
}

func InitBinance() exchange.Exchange {
	coin.Init()
	pair.Init()
	config := &exchange.Config{}
	config.Source = exchange.EXCHANGE_API
	// config.Source = exchange.JSON_FILE
	// config.SourceURI = "https://raw.githubusercontent.com/bitontop/gored/master/data"
	// utils.GetCommonDataFromJSON(config.SourceURI)

	conf.Exchange(exchange.BINANCE, config)

	ex := binance.CreateBinance(config)
	log.Printf("Initial [ %v ] ", ex.GetName())

	config = nil
	return ex
}

func Test_Binance_TradeHistory(t *testing.T) {
	e := InitBinance()
	p := pair.GetPairByKey("BTC|ETH")

	opTradeHistory := &exchange.PublicOperation{
		Type:      exchange.TradeHistory,
		EX:        e.GetName(),
		Pair:      p,
		DebugMode: true,
	}

	err := e.LoadPublicData(opTradeHistory)
	if err != nil {
		log.Printf("%v", err)
	}
}
