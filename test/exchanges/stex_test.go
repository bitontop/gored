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

	"github.com/bitontop/gored/exchange/stex"
	"github.com/bitontop/gored/test/conf"
	// "../../exchange/stex"
	// "../conf"
)

/********************Public API********************/
func Test_Stex(t *testing.T) {
	e := InitStex()

	pair := pair.GetPairByKey("BTC|ETH") // "ETH|AIB"

	// Test_Coins(e)
	// Test_Pairs(e)
	Test_Pair(e, pair)
	// Test_Orderbook(e, pair)
	// Test_ConstraintFetch(e, pair)
	// Test_Constraint(e, pair)

	Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.00000001, 100)
	// Test_Trading_Sell(e, pair, 0.06, 0.01)
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")

	// // Test Withdraw
	// opWithdraw := &exchange.AccountOperation{
	// 	Type:            exchange.Withdraw,
	// 	Coin:            pair.Target,
	// 	WithdrawAmount:  "1",
	// 	WithdrawAddress: "addr",
	// 	DebugMode:       true,
	// }
	// err := e.DoAccoutOperation(opWithdraw)
	// if err != nil {
	// 	log.Printf("%v", err)
	// }
	// log.Printf("WithdrawID: %v, err: %v", opWithdraw.WithdrawID, opWithdraw.Error)
}

func Test_STEX_TradeHistory(t *testing.T) {
	e := InitStex()
	p := pair.GetPairByKey("USDT|ETH")

	opTradeHistory := &exchange.PublicOperation{
		Type:      exchange.TradeHistory,
		Wallet:    exchange.SpotWallet,
		EX:        e.GetName(),
		Pair:      p,
		DebugMode: true,
	}

	err := e.LoadPublicData(opTradeHistory)
	if err != nil {
		log.Printf("%v", err)
	}

	log.Printf("TradeHistory: %s::%s", opTradeHistory.EX, opTradeHistory.Pair.Name)

	for _, d := range opTradeHistory.TradeHistory {
		log.Printf(">> %+v ", d)
	}
}

func Test_STEX_OrderBook(t *testing.T) {
	e := InitStex()
	p := pair.GetPairByKey("USDT|ETH")

	opSpotOrderBook := &exchange.PublicOperation{
		Type:   exchange.Orderbook,
		Wallet: exchange.SpotWallet,

		EX:   e.GetName(),
		Pair: p,

		Proxy:     "http://207.180.236.225:3128",
		DebugMode: true,
	}

	err := e.LoadPublicData(opSpotOrderBook)
	if err != nil {
		log.Printf("%v", err)
	}

	log.Printf("OrderBook: %s::%s", opSpotOrderBook.EX, opSpotOrderBook.Pair.Name)
	log.Printf("%+v", opSpotOrderBook.Maker)
}

func InitStex() exchange.Exchange {
	coin.Init()
	pair.Init()

	config := &exchange.Config{}
	config.Source = exchange.EXCHANGE_API
	conf.Exchange(exchange.STEX, config)

	ex := stex.CreateStex(config)
	log.Printf("Initial [ %v ] ", ex.GetName())

	config = nil
	return ex
}
