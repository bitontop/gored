package test

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"log"
	"testing"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/exchange/binance"
	"github.com/bitontop/gored/pair"
	"github.com/bitontop/gored/test/conf"
)

/********************Public API********************/
func Test_Binance(t *testing.T) {
	e := InitBinance()

	pair := pair.GetPairByKey("BTC|ETH")

	Test_Coins(e)
	Test_Pairs(e)
	Test_Pair(e, pair)
	Test_Orderbook(e, pair)
	// Test_Trading(e, pair, 0.00000001, 100)
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")
	Test_ConstraintFetch(e, pair)
	Test_Constraint(e, pair)
}

func InitBinance() exchange.Exchange {
	coin.Init()
	pair.Init()
	config := &exchange.Config{}
	config.Source = exchange.EXCHANGE_API
	conf.Exchange(exchange.BINANCE, config)

	ex := binance.CreateBinance(config)
	log.Printf("Initial [ %v ] ", ex.GetName())

	config = nil
	return ex
}
