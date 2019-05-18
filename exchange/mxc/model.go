package mxc

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

type JsonResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

type PairsData struct {
	PriceScale    int     `json:"priceScale"`
	QuantityScale int     `json:"quantityScale"`
	MinAmount     float64 `json:"minAmount"`
	BuyFeeRate    float64 `json:"buyFeeRate"`
	SellFeeRate   float64 `json:"sellFeeRate"`
}

type OrderBook struct {
	Asks []struct {
		Price    string `json:"price"`
		Quantity string `json:"quantity"`
	} `json:"asks"`
	Bids []struct {
		Price    string `json:"price"`
		Quantity string `json:"quantity"`
	} `json:"bids"`
}
