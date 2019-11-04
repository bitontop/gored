package test

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"fmt"
	"log"
	"sync"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
	"github.com/bitontop/gored/test/conf"

	initial "github.com/bitontop/gored/initial"
	utils "github.com/bitontop/gored/utils"
	// "github.com/davecgh/go-spew/spew"
)

func Init() *InitHandler {
	handler := &InitHandler{}
	var wg sync.WaitGroup

	coin.Init()
	pair.Init()

	utils.GetCommonDataFromJSON("https://raw.githubusercontent.com/bitontop/gored/master/data") //e.config.SourceURI)

	handler.ExMan = exchange.CreateExchangeManager()
	exchanges := handler.ExMan.GetSupportExchanges()
	initMan := initial.CreateInitManager()

	for _, ex := range exchanges {
		initEx(&wg, ex, initMan)
	}

	wg.Wait()

	INITEX := fmt.Sprintf("Init DONE !")
	log.Printf("%s %s %s ", COLOR_B_PINK, INITEX, COLOR_NO)

	// p := pair.GetPairByKey("BTC|ETH")

	// for _, ex := range exman.GetExchanges() {
	// 	status := ex.GetConstraintFetchMethod(p)
	// 	// spew.Dump(status)

	// 	if status.HasWithdraw {
	// 		log.Printf("%s  [ %12s ]  Withdraw Enabled!  %s", COLOR_GREEN,
	// 			ex.GetName(),
	// 			COLOR_NO)
	// 	}

	// }
	return handler
}

func initEx(wg *sync.WaitGroup, name exchange.ExchangeName, initMan *initial.InitManager) {

	wg.Add(1)

	config := &exchange.Config{}
	config.ExName = name
	config.Source = exchange.JSON_FILE
	config.SourceURI = "https://raw.githubusercontent.com/bitontop/gored/master/data"

	conf.Exchange(name, config)

	go func() {

		ex := initMan.Init(config)
		config = nil

		if ex != nil {
			log.Printf("%s Initial [ %12s ]  Coin:%d Pair:%d  %s", COLOR_B_PINK,
				name,
				len(ex.GetCoins()), len(ex.GetPairs()),
				COLOR_NO)
		} else {

			coinQ := 0
			if ex != nil {
				coinQ = len(ex.GetCoins())
			}
			pairQ := 0
			if ex != nil {
				pairQ = len(ex.GetPairs())
			}

			log.Printf("%s FAIL to Initial [ %v ]  Coin:%d Pair:%d  %s ", COLOR_RED,
				name,
				coinQ, pairQ,
				COLOR_NO)

		}

		if wg != nil {
			wg.Done()

		}

	}()
}
