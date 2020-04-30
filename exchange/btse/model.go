package btse

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
type CoinsData []struct {
	symbol string `json:"symbol"`
	Base   string `json:"base"`
	Quote  string `json:"quote"`
}

type PairsData []struct {
	symbol string `json:"symbol"`
	Base   string `json:"base"`
	Quote  string `json:"quote"`
}

type OrderBook struct {
	buyQuote []struct {
		price string `json:"price"`
		size  string `json:"size"`
	} `json:"buyQuote"`
	sellQuote []struct {
		price string `json:"price"`
		size  string `json:"size"`
	} `json:"sellQuote"`
	symbol    string `json:"symbol"`
	timestamp string `json:"timestamp"`
}

/********** Private API Structure**********/
type AccountBalances []struct {
	Total     float64 `json:"total"`
	Available float64 `json:"available"`
	Currency  string  `json:"currency"`
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
