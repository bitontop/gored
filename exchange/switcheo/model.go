package switcheo

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
type Token struct {
	Symbol           string      `json:"symbol"`
	Name             string      `json:"name"`
	Type             string      `json:"type"`
	Hash             string      `json:"hash"`
	Decimals         int64       `json:"decimals"`
	TransferDecimals int64       `json:"transfer_decimals"`
	Precision        int64       `json:"precision"`
	MinimumQuantity  string      `json:"minimum_quantity"`
	TradingActive    bool        `json:"trading_active"`
	IsStablecoin     bool        `json:"is_stablecoin"`
	StablecoinType   interface{} `json:"stablecoin_type"`
}

type CoinsData map[string]Token

type PairsData []string

type OrderBook struct {
	Bids [][]float64 `json:"bids"`
	Asks [][]float64 `json:"asks"`
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
