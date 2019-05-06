package test

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"log"

	"../coin"
	"../exchange"
	"../pair"
	"github.com/davecgh/go-spew/spew"
)

/********************Public API********************/
func Test_Coins(e exchange.Exchange) {
	for _, coin := range e.GetCoins() {
		log.Printf("%s Coin %+v", e.GetName(), coin)
	}
}

func Test_Pairs(e exchange.Exchange) {
	for _, pair := range e.GetPairs() {
		log.Printf("%s Pair %+v", e.GetName(), pair)
	}
}

func Test_Pair(e exchange.Exchange, pair *pair.Pair) {
	log.Printf("%s Pair: %+v", e.GetName(), pair)
	log.Printf("%s Pair Code: %s", e.GetName(), e.GetPairCode(pair))
}

func Test_Orderbook(e exchange.Exchange, p *pair.Pair) {
	maker, err := e.OrderBook(p)
	log.Printf("%s OrderBook %+v   error:%v", e.GetName(), maker, err)
}

/********************Private API********************/
func Test_Trading(e exchange.Exchange, p *pair.Pair, rate, quantity float64) {
	order, err := e.LimitBuy(p, quantity, rate)
	if err == nil {
		log.Printf("%s Limit Buy: %v", e.GetName(), order)
	} else {
		log.Printf("%s Limit Buy Err: %s", e.GetName(), err)
	}

	err = e.OrderStatus(order)
	if err == nil {
		log.Printf("%s Order Status: %v", e.GetName(), order)
	} else {
		log.Printf("%s Order Status Err: %s", e.GetName(), err)
	}

	err = e.CancelOrder(order)
	if err == nil {
		log.Printf("%s Cancel Order: %v", e.GetName(), order)
	} else {
		log.Printf("%s Cancel Err: %s", e.GetName(), err)
	}
}

func Test_Withdraw(e exchange.Exchange, c *coin.Coin, amount float64, addr string) {
	if e.Withdraw(c, amount, addr, "tag") {
		log.Printf("%s %s Withdraw Successful!", e.GetName(), c.Code)
	} else {
		log.Printf("%s %s Withdraw Failed!", e.GetName(), c.Code)
	}
}

/********************General********************/
func Test_ConstraintFetch(e exchange.Exchange, p *pair.Pair) {
	status := e.GetConstraintFetchMethod(p)
	spew.Dump(status)
}

func Test_Constraint(e exchange.Exchange, p *pair.Pair) {
	baseConstraint := e.GetCoinConstraint(p.Base)
	targerConstraint := e.GetCoinConstraint(p.Target)
	pairConstrinat := e.GetPairConstraint(p)

	log.Printf("%s %s Coin Constraint: %+v", e.GetName(), p.Base.Code, baseConstraint)
	log.Printf("%s %s Coin Constraint: %+v", e.GetName(), p.Target.Code, targerConstraint)
	log.Printf("%s %s Pair Constraint: %+v", e.GetName(), p.Name, pairConstrinat)
}
