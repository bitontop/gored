package digifinex

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

type JsonResponse struct {
	Data json.RawMessage `json:"data"`
	Date int             `json:"date"`
	Code int             `json:"code"`
}

type PairsData []struct {
	VolumePrecision int     `json:"volume_precision"`
	PricePrecision  int     `json:"price_precision"`
	Market          string  `json:"market"`
	MinAmount       float64 `json:"min_amount"`
	MinVolume       float64 `json:"min_volume"`
}

type OrderBook struct {
	Bids [][]float64 `json:"bids"`
	Asks [][]float64 `json:"asks"`
	Date int         `json:"date"`
	Code int         `json:"code"`
}

type AccountBalances struct {
	Code   int                `json:"code"`
	Date   int                `json:"date"`
	Free   map[string]float64 `json:"free"`
	Frozen map[string]float64 `json:"frozen"`
}

type PlaceOrder struct {
	Code    int    `json:"code"`
	OrderID string `json:"order_id"`
}

type OrderStatus []struct {
	OrderID        string  `json:"order_id"`
	CreatedDate    int     `json:"created_date"`
	FinishedDate   int     `json:"finished_date"`
	Price          float64 `json:"price"`
	Amount         float64 `json:"amount"`
	ExecutedAmount float64 `json:"executed_amount"`
	CashAmount     int     `json:"cash_amount"`
	AvgPrice       float64 `json:"avg_price"`
	Type           string  `json:"type"`
	Status         int     `json:"status"`
}

type CancelOrder struct {
	Success []interface{}   `json:"success"`
	Fail    [][]interface{} `json:"fail"`
}
