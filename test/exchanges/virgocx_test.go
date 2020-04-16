package test

import (
	"log"
	"testing"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
	// "../../exchange/virgocx"
	// "../conf"
)

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

/********************Public API********************/

func Test_Virgocx(t *testing.T) {
	e := InitEx(exchange.VIRGOCX)
	pair := pair.GetPairByKey("CAD|BTC")

	// Test_Coins(e)
	// Test_Pairs(e)
	Test_Pair(e, pair)
	// Test_Orderbook(e, pair)
	// Test_ConstraintFetch(e, pair)
	// Test_Constraint(e, pair)

	Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.00000001, 100)
	// Test_Trading_Sell(e, pair, 13000, 0.001)
	// Test_OrderStatus(e, pair, "1234567890")
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")
	// log.Println(e.GetTradingWebURL(pair))

	var err error

	// contract Kline
	opKline := &exchange.PublicOperation{
		Wallet:        exchange.ContractWallet,
		Type:          exchange.KLine,
		EX:            e.GetName(),
		Pair:          pair,
		KlineInterval: "5", // default to 5 if not provided
		DebugMode:     true,
	}
	err = e.LoadPublicData(opKline)
	if err != nil {
		log.Printf("%v", err)
	}

	for _, k := range opKline.Kline {
		log.Printf("%s ContractKline %+v", e.GetName(), k)
	}

	// Limit/Market PlaceOrder
	// opPlaceOrder := &exchange.AccountOperation{
	// 	Wallet:         exchange.SpotWallet,
	// 	Type:           exchange.PlaceOrder,
	// 	Ex:             e.GetName(),
	// 	Pair:           pair,
	// 	OrderDirection: exchange.BUY,
	// 	TradeType:      exchange.TRADE_LIMIT, // TRADE_MARKET
	// 	Rate:           5000,
	// 	Quantity:       0.01,
	// 	DebugMode:      true,
	// }
	// err = e.DoAccountOperation(opPlaceOrder)
	// if err != nil {
	// 	log.Printf("==%v", err)
	// }

	// Test_CancelOrder(e, pair, "2395474881")
	Test_CheckBalance(e, pair.Target, exchange.AssetWallet)
	Test_CheckAllBalance(e, exchange.SpotWallet)
}
