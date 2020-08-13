package test

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"log"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/initial"
	"github.com/bitontop/gored/pair"
	"github.com/bitontop/gored/test/conf"
	"github.com/bitontop/gored/utils"
	"github.com/davecgh/go-spew/spew"
)

func InitEx(exName exchange.ExchangeName) exchange.Exchange {
	coin.Init()
	pair.Init()
	config := &exchange.Config{}
	config.Source = exchange.EXCHANGE_API
	// config.Source = exchange.JSON_FILE
	// config.SourceURI = "https://raw.githubusercontent.com/bitontop/gored/master/data"
	// utils.GetCommonDataFromJSON(config.SourceURI)
	conf.Exchange(exName, config)

	inMan := initial.CreateInitManager()
	e := inMan.Init(config)
	log.Printf("Initial [ %v ] ", e.GetName())

	config = nil

	return e
}

func InitExFromJson(exName exchange.ExchangeName) exchange.Exchange {
	coin.Init()
	pair.Init()
	config := &exchange.Config{}
	// config.Source = exchange.EXCHANGE_API
	config.Source = exchange.JSON_FILE
	config.SourceURI = "https://raw.githubusercontent.com/bitontop/gored/master/data"
	utils.GetCommonDataFromJSON(config.SourceURI)
	conf.Exchange(exName, config)

	inMan := initial.CreateInitManager()
	e := inMan.Init(config)
	log.Printf("Initial [ %v ] from JSON", e.GetName())

	config = nil

	return e
}

/********************Public API********************/
func Test_Coins(e exchange.Exchange) {
	coins := e.GetCoins()
	if len(coins) > 0 {
		for _, coin := range coins {
			log.Printf("%s Coin %+v", e.GetName(), coin)
		}
	} else {
		log.Panicf("%s didn't get coins' data.", e.GetName())
	}
}

func Test_Pairs(e exchange.Exchange) {
	pairs := e.GetPairs()
	if len(pairs) > 0 {
		for _, pair := range pairs {
			log.Printf("%s Pair %+v", e.GetName(), pair)
		}
	} else {
		log.Panicf("%s didn't get pairs' data.", e.GetName())
	}
}

func Test_Pair(e exchange.Exchange, pair *pair.Pair) {
	log.Printf("%s Pair: %+v", e.GetName(), pair)
	log.Printf("%s Pair Code: %s", e.GetName(), e.GetSymbolByPair(pair))
	log.Printf("%s Coin Codes: %s, %s", e.GetName(), e.GetSymbolByCoin(pair.Base), e.GetSymbolByCoin(pair.Target))
}

func Test_Orderbook(e exchange.Exchange, p *pair.Pair) {
	maker, err := e.OrderBook(p)
	log.Printf("%s OrderBook %+v   error:%v", e.GetName(), maker, err)
}

/********************Private API********************/
func Test_Balance(e exchange.Exchange, p *pair.Pair) {
	e.UpdateAllBalances()

	base := e.GetBalance(p.Base)
	target := e.GetBalance(p.Target)
	log.Printf("Pair: %12s  Base %s: %f | Target %s: %f", p.Name, p.Base.Code, base, p.Target.Code, target)
}

func Test_Trading(e exchange.Exchange, p *pair.Pair, rate, quantity float64) {
	order, err := e.LimitBuy(p, quantity, rate)
	if err == nil {
		log.Printf("%s Limit Buy: %v", e.GetName(), order)

		err = e.OrderStatus(order)
		if err == nil {
			log.Printf("%s Order Status: %+v", e.GetName(), order)
		} else {
			log.Printf("%s Order Status Err: %s", e.GetName(), err)
		}

		err = e.CancelOrder(order)
		if err == nil {
			log.Printf("%s Cancel Order: %+v", e.GetName(), order)
		} else {
			log.Printf("%s Cancel Err: %s", e.GetName(), err)
		}

		err = e.OrderStatus(order)
		if err == nil {
			log.Printf("%s Order Status: %+v", e.GetName(), order)
		} else {
			log.Printf("%s Order Status Err: %s", e.GetName(), err)
		}
	} else {
		log.Printf("%s Limit Buy Err: %s", e.GetName(), err)
	}
}

func Test_Trading_Sell(e exchange.Exchange, p *pair.Pair, rate, quantity float64) {
	order, err := e.LimitSell(p, quantity, rate)
	if err == nil {
		log.Printf("%s Limit Sell: %+v", e.GetName(), order)

		err = e.OrderStatus(order)
		if err == nil {
			log.Printf("%s Order Status: %+v", e.GetName(), order)
		} else {
			log.Printf("%s Order Status Err: %s", e.GetName(), err)
		}

		err = e.CancelOrder(order)
		if err == nil {
			log.Printf("%s Cancel Order: %+v", e.GetName(), order)
		} else {
			log.Printf("%s Cancel Err: %s", e.GetName(), err)
		}

		err = e.OrderStatus(order)
		if err == nil {
			log.Printf("%s Order Status: %+v", e.GetName(), order)
		} else {
			log.Printf("%s Order Status Err: %s", e.GetName(), err)
		}
	} else {
		log.Printf("%s Limit Sell Err: %s", e.GetName(), err)
	}
}

// check auth only
func Test_OrderStatus(e exchange.Exchange, p *pair.Pair, orderID string) {
	order := &exchange.Order{
		Pair:      p,
		OrderID:   orderID,
		Rate:      0.001,
		Quantity:  100,
		Direction: exchange.Buy,
		Status:    exchange.New,
	}

	err := e.OrderStatus(order)
	if err == nil {
		log.Printf("%s Order Status: %v", e.GetName(), order)
	} else {
		log.Printf("%s Order Status Err: %s", e.GetName(), err)
	}
}

func Test_CancelOrder(e exchange.Exchange, p *pair.Pair, orderID string) {
	order := &exchange.Order{
		Pair:      p,
		OrderID:   orderID,
		Rate:      0.001,
		Quantity:  10,
		Direction: exchange.Buy,
		Status:    exchange.New,
	}

	err := e.CancelOrder(order)
	if err == nil {
		log.Printf("%s Cancel Order: %v", e.GetName(), order)
	} else {
		log.Printf("%s Cancel Order Err: %s", e.GetName(), err)
	}
}

func Test_Withdraw(e exchange.Exchange, c *coin.Coin, amount float64, addr string) {
	if e.Withdraw(c, amount, addr, "") {
		log.Printf("%s %s Withdraw Successful!", e.GetName(), c.Code)
	} else {
		log.Printf("%s %s Withdraw Failed!", e.GetName(), c.Code)
	}
}

func Test_DoWithdraw(e exchange.Exchange, c *coin.Coin, amount string, addr string, tag string) {
	opWithdraw := &exchange.AccountOperation{
		Type:            exchange.Withdraw,
		Coin:            c,
		WithdrawAmount:  amount,
		WithdrawAddress: addr,
		WithdrawTag:     tag,
		DebugMode:       true,
	}
	err := e.DoAccountOperation(opWithdraw)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	log.Printf("WithdrawID: %v, err: %v", opWithdraw.WithdrawID, opWithdraw.Error)
}

func Test_DoTransfer(e exchange.Exchange, c *coin.Coin, amount string, from, to exchange.WalletType) {
	opTransfer := &exchange.AccountOperation{
		Type:                exchange.Transfer,
		Coin:                c,
		TransferAmount:      amount,
		TransferFrom:        from,
		TransferDestination: to,
		DebugMode:           true,
	}
	err := e.DoAccountOperation(opTransfer)
	if err != nil {
		log.Printf("%v", err)
		return
	}
}

func Test_CheckBalance(e exchange.Exchange, c *coin.Coin, balanceType exchange.WalletType) {
	opBalance := &exchange.AccountOperation{
		Type:      exchange.Balance,
		Coin:      c,
		Wallet:    balanceType,
		DebugMode: true,
	}
	err := e.DoAccountOperation(opBalance)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	log.Printf("%v Account available: %v, frozen: %v", opBalance.Coin.Code, opBalance.BalanceAvailable, opBalance.BalanceFrozen)
}

func Test_CheckAllBalance(e exchange.Exchange, balanceType exchange.WalletType) {
	opAllBalance := &exchange.AccountOperation{
		Type:      exchange.BalanceList,
		Wallet:    balanceType,
		DebugMode: true,
	}
	err := e.DoAccountOperation(opAllBalance)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	for i, balance := range opAllBalance.BalanceList {
		if balance.BalanceAvailable+balance.BalanceFrozen == 0 {
			continue
		}
		log.Printf("AllAccount balance: %v Coin: %v, avaliable: %v, frozen: %v", i+1, balance.Coin.Code, balance.BalanceAvailable, balance.BalanceFrozen)
	}
	if len(opAllBalance.BalanceList) == 0 {
		log.Println("AllAccount 0 balance")
	}
}

func Test_TradeHistory(e exchange.Exchange, pair *pair.Pair) {
	opTradeHistory := &exchange.PublicOperation{
		Type:      exchange.TradeHistory,
		EX:        e.GetName(),
		Pair:      pair,
		DebugMode: true,
	}
	err := e.LoadPublicData(opTradeHistory)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	for _, trade := range opTradeHistory.TradeHistory {
		log.Printf("TradeHistory: %+v", trade)
	}
}

func Test_NewOrderBook(e exchange.Exchange, pair *pair.Pair) {
	opOrderBook := &exchange.PublicOperation{
		Type:      exchange.Orderbook,
		Wallet:    exchange.SpotWallet,
		EX:        e.GetName(),
		Pair:      pair,
		DebugMode: true,
	}
	err := e.LoadPublicData(opOrderBook)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	log.Printf("%s OrderBook %+v   error:%v", e.GetName(), opOrderBook.Maker, opOrderBook.Error)
}

func Test_CoinChainType(e exchange.Exchange, coin *coin.Coin) {
	opCoinChainType := &exchange.PublicOperation{
		Type:      exchange.CoinChainType,
		EX:        e.GetName(),
		Coin:      coin,
		DebugMode: true,
	}

	err := e.LoadPublicData(opCoinChainType)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	log.Printf("%s %s Chain Type: %s", opCoinChainType.EX, opCoinChainType.Coin.Code, opCoinChainType.CoinChainType)
}

func Test_DoOrderbook(e exchange.Exchange, pair *pair.Pair) {
	opTradeHistory := &exchange.PublicOperation{
		Type:      exchange.Orderbook,
		EX:        e.GetName(),
		Pair:      pair,
		DebugMode: true,
	}
	err := e.LoadPublicData(opTradeHistory)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	log.Printf("%s OrderBook %+v", e.GetName(), opTradeHistory.Maker)
}

func Test_AOOpenOrder(e exchange.Exchange, pair *pair.Pair) {
	op := &exchange.AccountOperation{
		Type:      exchange.GetOpenOrder,
		Wallet:    exchange.SpotWallet,
		Ex:        e.GetName(),
		Pair:      pair,
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
}

func Test_AOOrderHistory(e exchange.Exchange, pair *pair.Pair) {
	op := &exchange.AccountOperation{
		Type:      exchange.GetOrderHistory,
		Wallet:    exchange.SpotWallet,
		Ex:        e.GetName(),
		Coin:      pair.Base,
		Pair:      pair,
		DebugMode: true,
	}

	if err := e.DoAccountOperation(op); err != nil {
		log.Printf("%+v", err)
	} else {
		for _, o := range op.OrderHistory {
			log.Printf("%s OrderHistory %+v", e.GetName(), o)
		}
		if len(op.OrderHistory) == 0 {
			log.Printf("%s OrderHistory Response: %v", e.GetName(), op.CallResponce)
		}
	}
}

func Test_AODepositAddress(e exchange.Exchange, coin *coin.Coin) {
	op := &exchange.AccountOperation{
		Type:      exchange.GetDepositAddress,
		Wallet:    exchange.SpotWallet,
		Ex:        e.GetName(),
		Coin:      coin,
		DebugMode: true,
	}

	if err := e.DoAccountOperation(op); err != nil {
		log.Printf("%+v", err)
	} else {
		for chain, addr := range op.DepositAddresses {
			log.Printf("%s DepositAddresses: %v - %v %+v", e.GetName(), chain, addr.Coin.Code, addr)
		}
	}
}

func Test_AODepositHistory(e exchange.Exchange, pair *pair.Pair) {
	op := &exchange.AccountOperation{
		Type:      exchange.GetDepositHistory,
		Wallet:    exchange.SpotWallet,
		Ex:        e.GetName(),
		Coin:      pair.Base,
		Pair:      pair,
		DebugMode: true,
	}

	if err := e.DoAccountOperation(op); err != nil {
		log.Printf("%+v", err)
	} else {
		if len(op.DepositHistory) == 0 {
			log.Printf("%s DepositHistory Response: %v", e.GetName(), op.CallResponce)
		}
		for i, his := range op.DepositHistory {
			log.Printf("%s DepositHistory: %v %+v", e.GetName(), i, his)
		}
	}
}

func Test_AOWithdrawalHistory(e exchange.Exchange, pair *pair.Pair) {
	op := &exchange.AccountOperation{
		Type:      exchange.GetWithdrawalHistory,
		Wallet:    exchange.SpotWallet,
		Ex:        e.GetName(),
		Coin:      pair.Base,
		Pair:      pair,
		DebugMode: true,
	}

	if err := e.DoAccountOperation(op); err != nil {
		log.Printf("%+v", err)
	} else {
		if len(op.WithdrawalHistory) == 0 {
			log.Printf("%s WithdrawalHistory Response: %v", e.GetName(), op.CallResponce)
		}
		for i, his := range op.WithdrawalHistory {
			log.Printf("%s WithdrawalHistory: %v %+v", e.GetName(), i, his)
		}
	}
}

func Test_AOTransferHistory(e exchange.Exchange) {
	op := &exchange.AccountOperation{
		Type:        exchange.GetTransferHistory,
		Wallet:      exchange.SpotWallet,
		SubUserName: "sub1", // coinex only
		Ex:          e.GetName(),
		DebugMode:   true,
	}

	if err := e.DoAccountOperation(op); err != nil {
		log.Printf("%+v", err)
	} else {
		if len(op.TransferInHistory)+len(op.TransferOutHistory) == 0 {
			log.Printf("%s TransferInHistory Response: %v", e.GetName(), op.CallResponce)
		}
		for i, tIn := range op.TransferInHistory {
			log.Printf("%s TransferInHistory: %v %+v", e.GetName(), i, tIn)
		}
		for i, tOut := range op.TransferOutHistory {
			log.Printf("%s TransferOutHistory: %v %+v", e.GetName(), i, tOut)
		}
	}
}

func Test_TickerPrice(e exchange.Exchange) {
	opTickerPrice := &exchange.PublicOperation{
		Type:      exchange.GetTickerPrice,
		EX:        e.GetName(),
		Wallet:    exchange.SpotWallet,
		DebugMode: true,
	}
	err := e.LoadPublicData(opTickerPrice)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	for _, ticker := range opTickerPrice.TickerPrice {
		log.Printf("TickerPrice: %v, %v", ticker.Pair.Name, ticker.Price)
	}
}

func SubBalances(e exchange.Exchange, subID string) {
	// Sub Spot AllBalance
	opSubBalance := &exchange.AccountOperation{
		Wallet:       exchange.SpotWallet,
		Type:         exchange.SubBalanceList,
		SubAccountID: subID,
		Ex:           e.GetName(),
		DebugMode:    true,
	}
	err := e.DoAccountOperation(opSubBalance)
	if err != nil {
		log.Printf("==%v", err)
		return
	}
	for _, balance := range opSubBalance.BalanceList {
		log.Printf("SubBalances balance: Coin: %v, avaliable: %v, frozen: %v", balance.Coin.Code, balance.BalanceAvailable, balance.BalanceFrozen)
	}
	if len(opSubBalance.BalanceList) == 0 {
		log.Println("SubBalances 0 balance")
	}
	log.Printf("SubBalances JSON RESPONSE: %v", opSubBalance.CallResponce)
	log.Printf("SubBalances done")
}

func SubAllBalances(e exchange.Exchange) {
	// Sub All Spot AllBalance
	opSubAllBalance := &exchange.AccountOperation{
		Wallet:    exchange.SpotWallet,
		Type:      exchange.SubAllBalanceList,
		Ex:        e.GetName(),
		DebugMode: true,
	}
	err := e.DoAccountOperation(opSubAllBalance)
	if err != nil {
		log.Printf("==%v", err)
		return
	}
	for _, balance := range opSubAllBalance.BalanceList {
		log.Printf("SubAllBalances balance: Coin: %v, avaliable: %v, frozen: %v", balance.Coin.Code, balance.BalanceAvailable, balance.BalanceFrozen)
	}
	if len(opSubAllBalance.BalanceList) == 0 {
		log.Println("SubAllBalances 0 balance")
	}
	log.Printf("SubAllBalances JSON RESPONSE: %v", opSubAllBalance.CallResponce)
	log.Printf("SubAllBalances done")
}

func SubAccountList(e exchange.Exchange) {
	// Sub account list
	opSubAccountList := &exchange.AccountOperation{
		Wallet:    exchange.SpotWallet,
		Type:      exchange.GetSubAccountList,
		Ex:        e.GetName(),
		DebugMode: true,
	}
	err := e.DoAccountOperation(opSubAccountList)
	if err != nil {
		log.Printf("==%v", err)
		return
	}
	for _, account := range opSubAccountList.SubAccountList {
		log.Printf("AllSubAccount account: %+v", account)
	}
	if len(opSubAccountList.SubAccountList) == 0 {
		log.Println("No Sub Account Info")
	}
	log.Printf("SubAccountList JSON RESPONSE: %v", opSubAccountList.CallResponce)
	log.Printf("AllSubAccount done")
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

	log.Printf("%s %s Coin Constraint: %+v, %v", e.GetName(), p.Base.Code, baseConstraint, baseConstraint.Coin)
	log.Printf("%s %s Coin Constraint: %+v, %v", e.GetName(), p.Target.Code, targerConstraint, targerConstraint.Coin)
	log.Printf("%s %s Pair Constraint: %+v", e.GetName(), p.Name, pairConstrinat)
}
