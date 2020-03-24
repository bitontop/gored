package test

import (
	"log"
	"testing"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
	// "../../exchange/mxc"
	// "../conf"
)

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

/********************Public API********************/

func Test_Mxc(t *testing.T) {
	e := InitEx(exchange.MXC)
	// api only accept usdt based pair
	pair := pair.GetPairByKey("USDT|ETH") //"USDT|EOS"

	// Test_Coins(e)
	// Test_Pairs(e)
	Test_Pair(e, pair)
	Test_Orderbook(e, pair)
	// Test_ConstraintFetch(e, pair)
	// Test_Constraint(e, pair)

	var err error
	// Test Balance
	op2 := &exchange.AccountOperation{
		Type:   exchange.Balance,
		Coin:   pair.Target,
		Wallet: exchange.AssetWallet,
	}
	err = e.DoAccoutOperation(op2)
	if err != nil {
		log.Printf("%v", err)
	}
	log.Printf("Account available: %v, frozen: %v", op2.BalanceAvailable, op2.BalanceFrozen)

	// Test AllBalance
	op3 := &exchange.AccountOperation{
		Type:   exchange.BalanceList,
		Wallet: exchange.SpotWallet,
	}
	err = e.DoAccoutOperation(op3)
	if err != nil {
		log.Printf("%v", err)
	}
	for _, balance := range op3.BalanceList {
		log.Printf("Account balance: Coin: %v, avaliable: %v, frozen: %v", balance.Coin.Code, balance.BalanceAvailable, balance.BalanceFrozen)
	}

	// Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.00001, 100)
	// Test_Trading_Sell(e, pair, 220, 0.05)
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")
}
