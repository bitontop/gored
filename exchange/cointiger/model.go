package cointiger

import "encoding/json"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

/********** Public API Structure**********/
type PairsData map[string]PairsDetails

type PairsDetails []struct {
	BaseCurrency    string  `json:"baseCurrency"`
	QuoteCurrency   string  `json:"quoteCurrency"`
	PricePrecision  int     `json:"pricePrecision"`
	AmountPrecision int     `json:"amountPrecision"`
	WithdrawFeeMin  float64 `json:"withdrawFeeMin"`
	WithdrawFeeMax  float64 `json:"withdrawFeeMax"`
	WithdrawOneMin  float64 `json:"withdrawOneMin"`
	WithdrawOneMax  float64 `json:"withdrawOneMax"`
	DepthSelect     struct {
		Step0 string `json:"step0"`
		Step1 string `json:"step1"`
		Step2 string `json:"step2"`
	} `json:"depthSelect"`
}

type OrderBook struct {
	Symbol    string `json:"symbol"`
	DepthData struct {
		Tick struct {
			Buys [][]interface{} `json:"buys"`
			Asks [][]interface{} `json:"asks"`
		} `json:"tick"`
		Ts int64 `json:"ts"`
	} `json:"depth_data"`
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
