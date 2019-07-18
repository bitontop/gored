package test

import (
	"log"
	"testing"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"

	"github.com/bitontop/gored/exchange/tokok"
	"github.com/bitontop/gored/test/conf"
	// "../exchange/tokok"
	// "./conf"
)

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

/********************Public API********************/

func Test_Tokok(t *testing.T) {
	e := InitTokok()

	pair := pair.GetPairByKey("BTC|ETH")

	Test_Coins(e)
	Test_Pairs(e)
	Test_Pair(e, pair)
	Test_Orderbook(e, pair)
	Test_ConstraintFetch(e, pair)
	Test_Constraint(e, pair)

	Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.01, 0.01)
	// Test_Trading_Sell(e, pair, 0.04, 0.01)
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")
}

func InitTokok() exchange.Exchange {
	coin.Init()
	pair.Init()
	config := &exchange.Config{}
	config.Source = exchange.EXCHANGE_API
	// config.Source = exchange.JSON_FILE
	// config.SourceURI = "https://raw.githubusercontent.com/bitontop/gored/master/data"
	// utils.GetCommonDataFromJSON(config.SourceURI)
	conf.Exchange(exchange.TOKOK, config)

	ex := tokok.CreateTokok(config)
	log.Printf("Initial [ %v ] ", ex.GetName())

	config = nil
	return ex
}
