package hibitex

import "encoding/json"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct {
	Code    string          `json:"code"`
	Msg     string          `json:"msg"`
	Message interface{}     `json:"message"`
	Data    json.RawMessage `json:"data"`
}

/********** Public API Structure**********/
type Ticker struct {
	Symbol          string `json:"symbol"`
	CountCoin       string `json:"count_coin"`
	AmountPrecision int    `json:"amount_precision"`
	BaseCoin        string `json:"base_coin"`
	PricePrecision  int    `json:"price_precision"`
}

type CoinsData []Ticker

type PairsData []Ticker

type Order struct {
	Amount float64 `json:"amount"`
	Price  float64 `json:"price"`
	Id     int     `json:"id"`
	Type   string  `json:"type"`
	Ts     int     `json:"ts"`
	Ds     string  `json:"ds"`
}

type OrderBook []Order

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
