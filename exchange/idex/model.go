package idex

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

type JsonResponse struct {
	Error  Error           `json:"error"`
	Result json.RawMessage `json:"result"`
	Cmd    string          `json:"cmd"`
	Ver    string          `json:"ver"`
}

type Error struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

type CoinsData map[string]CoinDetails

type CoinDetails struct {
	Name     string `json:"name"`
	Decimals int    `json:"decimals"`
	Address  string `json:"address"`
	Slug     string `json:"slug"`
}

type PairsData map[string]interface{}

type OrderBook struct {
	Asks []struct {
		Price     string `json:"price"`
		Amount    string `json:"amount"`
		Total     string `json:"total"`
		OrderHash string `json:"orderHash"`
		Params    struct {
			TokenBuy      string `json:"tokenBuy"`
			BuySymbol     string `json:"buySymbol"`
			BuyPrecision  int    `json:"buyPrecision"`
			AmountBuy     string `json:"amountBuy"`
			TokenSell     string `json:"tokenSell"`
			SellSymbol    string `json:"sellSymbol"`
			SellPrecision int    `json:"sellPrecision"`
			AmountSell    string `json:"amountSell"`
			Expires       int    `json:"expires"`
			Nonce         int64  `json:"nonce"`
			User          string `json:"user"`
		} `json:"params"`
	} `json:"asks"`
	Bids []struct {
		Price     string `json:"price"`
		Amount    string `json:"amount"`
		Total     string `json:"total"`
		OrderHash string `json:"orderHash"`
		Params    struct {
			TokenBuy      string `json:"tokenBuy"`
			BuySymbol     string `json:"buySymbol"`
			BuyPrecision  int    `json:"buyPrecision"`
			AmountBuy     string `json:"amountBuy"`
			TokenSell     string `json:"tokenSell"`
			SellSymbol    string `json:"sellSymbol"`
			SellPrecision int    `json:"sellPrecision"`
			AmountSell    string `json:"amountSell"`
			Expires       int    `json:"expires"`
			Nonce         int    `json:"nonce"`
			User          string `json:"user"`
		} `json:"params"`
	} `json:"bids"`
}

type AccountBalances []struct {
	Result struct {
		TotalBtc   string `json:"total_btc"`
		TotalCny   string `json:"total_cny"`
		TotalUsd   string `json:"total_usd"`
		AssetsList []struct {
			CoinSymbol string `json:"coin_symbol"`
			Balance    string `json:"balance"`
			Freeze     string `json:"freeze"`
			BTCValue   string `json:"BTCValue"`
			CNYValue   string `json:"CNYValue"`
			USDValue   string `json:"USDValue"`
		} `json:"assets_list"`
	} `json:"result"`
	Cmd string `json:"cmd"`
}

type PlaceOrder struct {
	Timestamp   int    `json:"timestamp"`
	Market      string `json:"market"`
	OrderNumber int    `json:"orderNumber"`
	OrderHash   string `json:"orderHash"`
	Price       string `json:"price"`
	Amount      string `json:"amount"`
	Total       string `json:"total"`
	Type        string `json:"type"`
	Params      struct {
		TokenBuy      string `json:"tokenBuy"`
		BuyPrecision  int    `json:"buyPrecision"`
		AmountBuy     string `json:"amountBuy"`
		TokenSell     string `json:"tokenSell"`
		SellPrecision int    `json:"sellPrecision"`
		AmountSell    string `json:"amountSell"`
		Expires       int    `json:"expires"`
		Nonce         string `json:"nonce"`
		User          string `json:"user"`
	} `json:"params"`
}

type OrderStatus []struct {
	Result struct {
		ID             int    `json:"id"`
		CreatedAt      int64  `json:"createdAt"`
		AccountType    int    `json:"account_type"`
		Pair           string `json:"pair"`
		CoinSymbol     string `json:"coin_symbol"`
		CurrencySymbol string `json:"currency_symbol"`
		OrderSide      int    `json:"order_side"`
		OrderType      int    `json:"order_type"`
		Price          string `json:"price"`
		DealPrice      string `json:"deal_price"`
		Amount         string `json:"amount"`
		Money          string `json:"money"`
		DealAmount     string `json:"deal_amount"`
		DealPercent    string `json:"deal_percent"`
		DealMoney      string `json:"deal_money"`
		Status         int    `json:"status"`
		Unexecuted     string `json:"unexecuted"`
		OrderFrom      int    `json:"order_from"`
	} `json:"result"`
	Cmd string `json:"cmd"`
}

type CancelOrder []struct {
	Result string `json:"result"`
	Cmd    string `json:"cmd"`
}
