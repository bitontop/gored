package test

import (
	"log"
	"testing"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"

	"github.com/bitontop/gored/exchange/kucoin"
	"github.com/bitontop/gored/test/conf"
	// "../../exchange/kucoin"
	// "../conf"
)

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

/********************Public API********************/

func Test_Kucoin(t *testing.T) {
	e := InitKucoin()

	pair := pair.GetPairByKey("BTC|ETH")

	// Test_Coins(e)
	// Test_Pairs(e)
	Test_Pair(e, pair)
	// Test_Orderbook(e, pair)
	// Test_ConstraintFetch(e, pair)
	// Test_Constraint(e, pair)

	// Test Transfer
	op := &exchange.AccountOperation{
		Type:                exchange.Transfer,
		Coin:                pair.Target,
		TransferAmount:      "0.1",
		TransferFrom:        exchange.AssetWallet,
		TransferDestination: exchange.SpotWallet,
	}
	err := e.DoAccoutOperation(op)
	if err != nil {
		log.Printf("%v", err)
	}

	// Test Balance
	op2 := &exchange.AccountOperation{
		Type:        exchange.Balance,
		Coin:        pair.Target,
		BalanceType: exchange.AssetWallet,
	}
	err = e.DoAccoutOperation(op2)
	if err != nil {
		log.Printf("%v", err)
	}
	log.Printf("Account available: %v, frozen: %v", op2.BalanceAvailable, op2.BalanceFrozen)

	// Test AllBalance
	op3 := &exchange.AccountOperation{
		Type:        exchange.BalanceList,
		BalanceType: exchange.AssetWallet,
	}
	err = e.DoAccoutOperation(op3)
	if err != nil {
		log.Printf("%v", err)
	}
	for _, balance := range op3.BalanceList {
		log.Printf("Account balance: Coin: %v, avaliable: %v, frozen: %v", balance.Coin.Code, balance.BalanceAvailable, balance.BalanceFrozen)
	}

	// Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.00000001, 100)
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")
}

func InitKucoin() exchange.Exchange {
	coin.Init()
	pair.Init()
	config := &exchange.Config{}
	config.Source = exchange.EXCHANGE_API
	// config.Source = exchange.JSON_FILE
	// config.SourceURI = "https://raw.githubusercontent.com/bitontop/gored/master/data"
	// utils.GetCommonDataFromJSON(config.SourceURI)

	conf.Exchange(exchange.KUCOIN, config)

	ex := kucoin.CreateKucoin(config)
	log.Printf("Initial [ %v ] ", ex.GetName())

	config = nil
	return ex
}
