package test

import (
	"log"
	"testing"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"

	"github.com/bitontop/gored/exchange/kraken"
	"github.com/bitontop/gored/test/conf"
	// "../../exchange/kraken"
	// "../conf"
)

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

/********************Public API********************/

func Test_Kraken(t *testing.T) {
	e := InitKraken()

	pair := pair.GetPairByKey("ETH|ETC")

	// Test_Coins(e)
	// Test_Pairs(e)
	// Test_Pair(e, pair)
	// Test_Orderbook(e, pair)
	// Test_ConstraintFetch(e, pair)
	// Test_Constraint(e, pair)

	Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.0001, 100)
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")

	Test_DoWithdraw(e, pair.Target, "1", "0x37E0Fc27C6cDB5035B2a3d0682B4E7C05A4e6C46", "tag")
}

func InitKraken() exchange.Exchange {
	coin.Init()
	pair.Init()
	config := &exchange.Config{}
	config.Source = exchange.EXCHANGE_API

	// config.Source = exchange.JSON_FILE
	// config.SourceURI = "https://raw.githubusercontent.com/bitontop/gored/master/data"
	// utils.GetCommonDataFromJSON(config.SourceURI)

	conf.Exchange(exchange.KRAKEN, config)

	ex := kraken.CreateKraken(config)
	log.Printf("Initial [ %v ] ", ex.GetName())

	config = nil
	return ex
}
