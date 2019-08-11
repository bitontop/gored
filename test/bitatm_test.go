package test

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"log"
	"testing"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"

	"github.com/bitontop/gored/exchange/bitatm"
	"github.com/bitontop/gored/pair"
	"github.com/bitontop/gored/test/conf"
	// "../exchange/bitatm"
)

/********************Public API********************/
func Test_BitATM(t *testing.T) {
	e := InitBitATM()

	pair := pair.GetPairByKey("BHD|BTC")

	// Test_Coins(e)
	// Test_Pairs(e)
	// Test_Pair(e, pair)
	// Test_Orderbook(e, pair)
	// Test_ConstraintFetch(e, pair)
	// Test_Constraint(e, pair)

	Test_Balance(e, pair)
	Test_Trading(e, pair, 0.001, 100)
	//Test_Withdraw(e, pair.Base, 1, "ADDRESS")
}

func Test_BitatmOrder(t *testing.T) {
	e := InitBitATM()

	pair := pair.GetPairByKey("BHD|BTC")

	order := &exchange.Order{
		Pair:     pair,
		OrderID:  "123456",
		Rate:     0.001,
		Quantity: 100,
		Side:     "Buy",
		Status:   exchange.New,
	}

	err := e.OrderStatus(order)
	if err == nil {
		log.Printf("%s Order Status: %v", e.GetName(), order)
	} else {
		log.Printf("%s Order Status Err: %s", e.GetName(), err)
	}
}

func InitBitATM() exchange.Exchange {
	coin.Init()
	pair.Init()
	config := &exchange.Config{}
	config.Source = exchange.EXCHANGE_API
	conf.Exchange(exchange.BITATM, config)

	ex := bitatm.CreateBitATM(config)
	log.Printf("Initial [ %v ] ", ex.GetName())

	config = nil
	return ex
}
