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
	AverageFillPrice float64 `json:"averageFillPrice"`
	ClOrderID        string  `json:"clOrderID"`
	FillSize         float64 `json:"fillSize"`
	Message          string  `json:"message"`
	OrderID          string  `json:"orderID"`
	OrderType        float64 `json:"orderType"`
	Price            float64 `json:"price"`
	Side             string  `json:"side"`
	Size             float64 `json:"size"`
	Status           float64 `json:"status"`
	StopPrice        float64 `json:"stopPrice"`
	Symbol           string  `json:"symbol"`
	Timestamp        float64 `json:"timestamp"`
	Trigger          bool    `json:"trigger"`
	TriggerPrice     float64 `json:"triggerPrice"`
}
