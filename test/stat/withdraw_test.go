package test

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"log"
	"testing"

	"github.com/bitontop/gored/pair"
	// "github.com/davecgh/go-spew/spew"
)

func Test_WithdrawEnableExchanges(t *testing.T) {
	// var wg sync.WaitGroup

	// coin.Init()
	// pair.Init()

	// utils.GetCommonDataFromJSON("https://raw.githubusercontent.com/bitontop/gored/master/data") //e.config.SourceURI)

	// exman := exchange.CreateExchangeManager()
	// exchanges := exman.GetSupportExchanges()
	// initMan := initial.CreateInitManager()

	// for _, ex := range exchanges {
	// 	InitEx(&wg, ex, initMan)
	// }

	// wg.Wait()

	// INITEX := fmt.Sprintf("Init DONE !")
	// log.Printf("%s %s %s ", COLOR_B_PINK, INITEX, COLOR_NO)

	init := Init()

	p := pair.GetPairByKey("BTC|ETH")

	for _, ex := range init.ExMan.GetExchanges() {
		status := ex.GetConstraintFetchMethod(p)
		// spew.Dump(status)

		if status.HasWithdraw {
			log.Printf("%s  [ %12s ]  Withdraw Enabled!  %s", COLOR_GREEN,
				ex.GetName(),
				COLOR_NO)
		}

	}

}
