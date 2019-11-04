package test

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"fmt"
	"log"
	"testing"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"

	"github.com/bitontop/gored/exchange/bgogo"
	"github.com/bitontop/gored/test/conf"
	// "../exchange/bgogo"
	// "./conf"
)

/********************Public API********************/
func Test_Bgogo(t *testing.T) {
	e := InitBgogo()

	//"BTC|ETH" 即 ETH_BTC
	pair := pair.GetPairByKey("BTC|ETH")

	Test_Coins(e)
	Test_Pairs(e)
	Test_Pair(e, pair)
	Test_Orderbook(e, pair)
	Test_ConstraintFetch(e, pair)
	Test_Constraint(e, pair)

	log.Println(e.GetTradingWebURL(pair))
	// Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.00000001, 100)
	// Test_Trading_Sell(e, pair, 0.00000001, 667)
	// Test_Withdraw(e, pair.Base, 1.5, "1HB5XMLmzFVj8ALj6mfBsbifRoD4miY36v")
}

func InitBgogo() exchange.Exchange {
	coin.Init()
	pair.Init()
	config := &exchange.Config{}
	config.Source = exchange.EXCHANGE_API
	conf.Exchange(exchange.BITCOIN, config)
	fmt.Printf("%+v\n", config)

	ex := bgogo.CreateBgogo(config)
	log.Printf("Initial [ %v ] ", ex.GetName())

	config = nil
	return ex
}
