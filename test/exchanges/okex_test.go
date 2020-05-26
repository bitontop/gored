package test

import (
	"log"
	"testing"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
	// "../../exchange/okex"
	// "../conf"
)

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

/********************Public API********************/

func Test_Okex(t *testing.T) {
	e := InitEx(exchange.OKEX)
	pair := pair.GetPairByKey("BTC|ETH")

	// Test_Coins(e)
	// Test_Pairs(e)
	Test_Pair(e, pair)
	// Test_Orderbook(e, pair)
	// Test_NewOrderBook(e, pair)
	// Test_ConstraintFetch(e, pair)
	// Test_Constraint(e, pair)
	// Test_TickerPrice(e)

	// new interface methods
	// Test_DoWithdraw(e, pair.Target, "1", "0x37E0Fc27C6cDB5035B2a3d0682B4E7C05A4e6C46", "tag")
	// Test_DoTransfer(e, pair.Target, "10", exchange.AssetWallet, exchange.SpotWallet)
	// Test_CheckBalance(e, pair.Target, exchange.AssetWallet)
	Test_CheckAllBalance(e, exchange.SpotWallet)

	// okex.Socket(pair)
	Test_Balance(e, pair)
	// Test_Trading(e, pair, 0.00000001, 100)
	// Test_Withdraw(e, pair.Base, 1, "ADDRESS")

	// =====================================================================
	// TransferHistory
	opCTransferHistory := &exchange.AccountOperation{
		Type:      exchange.GetTransferHistory,
		Wallet:    exchange.SpotWallet,
		Ex:        e.GetName(),
		DebugMode: true,
	}

	if err := e.DoAccountOperation(opCTransferHistory); err != nil {
		log.Printf("%+v", err)
	} else {
		for _, o := range opCTransferHistory.TransferInHistory {
			log.Printf("%s TransferInHistory %+v", e.GetName(), o)
		}
		for _, o := range opCTransferHistory.TransferOutHistory {
			log.Printf("%s TransferOutHistory %+v", e.GetName(), o)
		}
		log.Printf("Contract TransferHistory CallResponse: %+v", opCTransferHistory.CallResponce)
	}
	// =====================================================================

	// Test_TradeHistory(e, pair)
}
