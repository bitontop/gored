package bkex

import "encoding/json"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
	Msg  string          `json:"msg"`
}

/********** Public API Structure**********/
type Coin struct {
	CoinType          string  `json:"coinType"`
	MaxWithdrawOneDay int     `json:"maxWithdrawOneDay"`
	MaxWithdrawSingle int     `json:"maxWithdrawSingle"`
	MinWithdrawSingle float64 `json:"minWithdrawSingle"`
	SupportDeposit    bool    `json:"supportDeposit"`
	SupportTrade      bool    `json:"supportTrade"`
	SupportWithdraw   bool    `json:"supportWithdraw"`
	WithdrawFee       float64 `json:"withdrawFee"`
}

type Pair struct {
	AmountPrecision  int    `json:"amountPrecision"`
	DefaultPrecision int    `json:"defaultPrecision"`
	Pair             string `json:"pair"`
	SupportTrade     bool   `json:"supportTrade"`
}

type CoinsData []Coin

type PairsData []Pair

type ExchangeData struct {
	CoinTypes CoinsData
	Pairs     PairsData
}

type OrderItem struct {
	Amt       float64 `json:"amt"`
	Direction string  `json:"direction"`
	Pair      string  `json:"pair"`
	Price     float64 `json:"price"`
}

type OrderBook struct {
	Asks []OrderItem `json:"asks"`
	Bids []OrderItem `json:"bids"`
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
