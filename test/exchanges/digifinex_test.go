package test

import (
	"log"
	"testing"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"

	"github.com/bitontop/gored/exchange/digifinex"
	"github.com/bitontop/gored/test/conf"
	// "../../exchange/digifinex"
	// "../conf"
)

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

/********************Public API********************/

func Test_Digifinex(t *testing.T) {
	e := InitDigifinex()

	pair := pair.GetPairByKey("BTC|BTT")

	Test_Coins(e)
	Test_Pairs(e)
	Test_Pair(e, pair)
	Test_Orderbook(e, pair)
	Test_ConstraintFetch(e, pair)
	Test_Constraint(e, pair)

	// Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.0001, 10)
	// Test_OrderStatus(e, pair, "1234567890")
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")
}

func InitDigifinex() exchange.Exchange {
	coin.Init()
	pair.Init()

	config := &exchange.Config{}
	config.Source = exchange.EXCHANGE_API
	// config.Source = exchange.JSON_FILE
	// config.SourceURI = "https://raw.githubusercontent.com/bitontop/gored/master/data"
	// utils.GetCommonDataFromJSON(config.SourceURI)

	conf.Exchange(exchange.DIGIFINEX, config)

	ex := digifinex.CreateDigifinex(config)
	log.Printf("Initial [ %v ] ", ex.GetName())

	config = nil
	return ex
}
