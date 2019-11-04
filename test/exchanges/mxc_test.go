package test

import (
	"log"
	"testing"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"

	"github.com/bitontop/gored/exchange/mxc"
	"github.com/bitontop/gored/test/conf"
	// "../exchange/mxc"
	// "./conf"
)

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

/********************Public API********************/

func Test_Mxc(t *testing.T) {
	e := InitMxc()

	pair := pair.GetPairByKey("BTC|BTT")

	Test_Coins(e)
	Test_Pairs(e)
	Test_Pair(e, pair)
	Test_Orderbook(e, pair)
	Test_ConstraintFetch(e, pair)
	Test_Constraint(e, pair)

	// Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.00000001, 100)
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")
}

// test url and signature
func Test_MxcOrder(t *testing.T) {
	e := InitMxc()

	var order *exchange.Order
	err := e.OrderStatus(order)
	if err == nil {
		log.Printf("%s Order Status: %v", e.GetName(), order)
	} else {
		log.Printf("%s Order Status Err: %s", e.GetName(), err)
	}
}

func InitMxc() exchange.Exchange {
	coin.Init()
	pair.Init()

	config := &exchange.Config{}
	config.Source = exchange.EXCHANGE_API
	// config.Source = exchange.JSON_FILE
	// config.SourceURI = "https://raw.githubusercontent.com/bitontop/gored/master/data"
	// utils.GetCommonDataFromJSON(config.SourceURI)

	conf.Exchange(exchange.MXC, config)

	ex := mxc.CreateMxc(config)
	log.Printf("Initial [ %v ] ", ex.GetName())

	config = nil
	return ex
}
