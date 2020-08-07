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
	e := InitEx(exchange.KUCOIN)
	pair := pair.GetPairByKey("BTC|ETH")

	// Test_TradeHistory(e, pair)
	// Test_NewOrderBook(e, pair)

	// Test_Coins(e)
	// Test_Pairs(e)
	Test_Pair(e, pair)
	// Test_Orderbook(e, pair)
	// Test_ConstraintFetch(e, pair)
	// Test_Constraint(e, pair)

	Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.00000001, 100)
	// Test_Trading_Sell(e, pair, 0.06, 0.01)
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")

	// SubBalances(e, "5cbd31ab9c93e9280cd36a0a")
	// SubAccountList(e)
	// SubAllBalances(e)
	// ===============================================
	op := &exchange.AccountOperation{
		Type:    exchange.GetOpenOrder,
		Wallet:  exchange.SpotWallet,
		Sandbox: true,
		Ex:      e.GetName(),
		Pair:    pair,
		// StartTime: 1596838564000,
		// EndTime:   1596838564001,
		DebugMode: true,
	}

	if err := e.DoAccountOperation(op); err != nil {
		log.Printf("%+v", err)
	} else {
		for _, o := range op.OpenOrders {
			log.Printf("%s OpenOrders: %v %+v", e.GetName(), o.Pair.Name, o)
		}
		if len(op.OpenOrders) == 0 {
			log.Printf("%s OpenOrder Response: %v", e.GetName(), op.CallResponce)
		}
	}
	// ===============================================
	// Test_AOOpenOrder(e, pair)
	Test_TickerPrice(e)

	// Test_CheckBalance(e, pair.Target, exchange.AssetWallet)
	// Test_CheckAllBalance(e, exchange.SpotWallet)
	// Test_DoTransfer(e, pair.Target, "1", exchange.AssetWallet, exchange.SpotWallet)
	// Test_DoWithdraw(e, pair.Target, "1", "0x37E0Fc27C6cDB5035B2a3d0682B4E7C05A4e6C46", "tag")
}

func Test_KUCOIN_TradeHistory(t *testing.T) {
	e := InitKucoin()
	p := pair.GetPairByKey("USDT|LTC")

	opTradeHistory := &exchange.PublicOperation{
		Type:      exchange.TradeHistory,
		EX:        e.GetName(),
		Pair:      p,
		DebugMode: true,
	}

	err := e.LoadPublicData(opTradeHistory)
	if err != nil {
		log.Printf("%v", err)
	}

	log.Printf("TradeHistory: %s::%s", opTradeHistory.EX, opTradeHistory.Pair.Name)

	for _, d := range opTradeHistory.TradeHistory {
		log.Printf(">> %+v ", d)
	}
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
