package coindeal

import "encoding/json"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

/********** Public API Structure**********/
type Tick struct {
	ID            int    `json:"id"`
	Last          string `json:"last"`
	LowestAsk     string `json:"lowestAsk"`
	HighestBid    string `json:"highestBid"`
	PercentChange string `json:"percentChange"`
	BaseVolume    string `json:"baseVolume"`
	QuoteVolume   string `json:"quoteVolume"`
	IsFrozen      string `json:"isFrozen"`
	High24hr      string `json:"high24hr"`
	Low24hr       string `json:"low24hr"`
	CurrencyName  string `json:"currencyName"`
	WalletStatus  string `json:"walletStatus"`
}

type CoinsData map[string]Tick

type PairsData map[string]Tick

type OrderInfo struct {
	Price  string `json:"price"`
	Amount string `json:"amount"`
}
type OrderBook struct {
	Asks []OrderInfo `json:"ask"`
	Bids []OrderInfo `json:"bid"`
}

/********** Private API Structure**********/
type AccountBalances []struct {
	Asset     string  `json:"asset"`
	Total     float64 `json:"total"`
	Available float64 `json:"available"`
	Locked    float64 `json:"locked"`
}

type WithdrawResponse struct {
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

type PlaceOrder struct {
	Symbol       string `json:"symbol"`
	OrderID      string `json:"orderId"`
	Side         string `json:"side"`
	Type         string `json:"type"`
	Price        string `json:"price"`
	AveragePrice string `json:"executedQty"`
	OrigQty      string `json:"origQty"`
	ExecutedQty  string `json:"executedQty"`
	Status       string `json:"status"`
	TimeInForce  string `json:"timeInForce"`
}
