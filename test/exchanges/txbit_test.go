package test

import (
	"log"
	"testing"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"

	"github.com/bitontop/gored/exchange/txbit"
	"github.com/bitontop/gored/test/conf"
	// "../../exchange/txbit"
	// "../conf"
)

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

/********************Public API********************/

func Test_Txbit(t *testing.T) {
	e := InitTxbit()

	pair := pair.GetPairByKey("BTC|ETH")

	Test_Coins(e)
	Test_Pairs(e)
	Test_Pair(e, pair)
	Test_Orderbook(e, pair)
	Test_ConstraintFetch(e, pair)
	Test_Constraint(e, pair)

	// Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.001, 0.033)
	// Test_Trading_Sell(e, pair, 0.050012345678, 0.10001234567)
	// Test_OrderStatus(e, pair, "b5d7d18c-61fb-479e-8ee6-b222ced93e56")
	// Test_CancelOrder(e, pair, "b5d7d18c-61fb-479e-8ee6-b222ced93e56")
	// Test_Withdraw(e, pair.Target, 0.1, "0xf252be0c7758094a37bf10a4cbf4dec0d69b7bcc")
	// log.Println(e.GetTradingWebURL(pair))
}

func InitTxbit() exchange.Exchange {
	coin.Init()
	pair.Init()

	config := &exchange.Config{}
	config.Source = exchange.EXCHANGE_API
	// config.Source = exchange.JSON_FILE
	// config.SourceURI = "https://raw.githubusercontent.com/bitontop/gored/master/data"
	// utils.GetCommonDataFromJSON(config.SourceURI)

	conf.Exchange(exchange.TXBIT, config)

	ex := txbit.CreateTxbit(config)
	log.Printf("Initial [ %v ] ", ex.GetName())

	config = nil
	return ex
}
