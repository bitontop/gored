package virgocx

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

type JsonResponse struct {
	Code    int             `json:"code"`
	Msg     string          `json:"msg"`
	Data    json.RawMessage `json:"data"`
	Success bool            `json:"success"`
}

type PairsData []struct {
	Volume     string  `json:"volume"`
	Symbol     string  `json:"symbol"`
	High       string  `json:"high"`
	Last       string  `json:"last"`
	Low        string  `json:"low"`
	Buy        float64 `json:"buy"`
	Sell       float64 `json:"sell"`
	ID         int     `json:"id"`
	ChangeRate string  `json:"changeRate"`
	Open       string  `json:"open"`
}

type OrderBook struct {
	Asks []struct {
		Price  float64 `json:"price"`
		Qty    float64 `json:"qty"`
		Volume float64 `json:"volume"`
	} `json:"asks"`
	Bids []struct {
		Price  float64 `json:"price"`
		Qty    float64 `json:"qty"`
		Volume float64 `json:"volume"`
	} `json:"bids"`
	Ts int64 `json:"ts"`
}

//--------------

type AccountBalances []struct {
	Currency      string  `json:"Currency"`
	Balance       float64 `json:"Balance"`
	Available     float64 `json:"Available"`
	Pending       float64 `json:"Pending"`
	CryptoAddress string  `json:"CryptoAddress"`
	Requested     bool    `json:"Requested"`
	Uuid          string  `json:"Uuid"`
}

type Uuid struct {
	Id string `json:"uuid"`
}

type PlaceOrder struct {
	AccountId                  string
	OrderUuid                  string `json:"OrderUuid"`
	Exchange                   string `json:"Exchange"`
	Type                       string
	Quantity                   float64 `json:"Quantity"`
	QuantityRemaining          float64 `json:"QuantityRemaining"`
	Limit                      float64 `json:"Limit"`
	Reserved                   float64
	ReserveRemaining           float64
	CommissionReserved         float64
	CommissionReserveRemaining float64
	CommissionPaid             float64
	Price                      float64 `json:"Price"`
	PricePerUnit               float64 `json:"PricePerUnit"`
	Opened                     string
	Closed                     string
	IsOpen                     bool
	Sentinel                   string
	CancelInitiated            bool
	ImmediateOrCancel          bool
	IsConditional              bool
	Condition                  string
	ConditionTarget            string
}
