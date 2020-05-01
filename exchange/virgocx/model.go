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

type PlaceOrder struct {
	OrderID string `json:"OrderId"`
}

type AccountBalances []struct {
	FreezingBalance float64 `json:"freezingBalance"`
	Total           float64 `json:"total"`
	Balance         float64 `json:"balance"` // available
	CoinName        string  `json:"coinName"`
}

type CancelOrder struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Data    string `json:"data"`
	Success bool   `json:"success"`
}

type OrderStatus []SingleOrder

type SingleOrder struct {
	CreateTime int64   `json:"createTime"`
	Price      float64 `json:"price"`
	Qty        float64 `json:"qty"`
	TradeQty   float64 `json:"tradeQty"`
	ID         int     `json:"id"`
	Type       int     `json:"type"`
	Direction  int     `json:"direction"`
	Status     int     `json:"status"`
}

type RawKline []struct {
	ID         int         `json:"id"`
	MarketID   int         `json:"marketId"`
	Open       float64     `json:"open"`
	High       float64     `json:"high"`
	Low        float64     `json:"low"`
	Close      float64     `json:"close"`
	Qty        float64     `json:"qty"`
	CreateTime int64       `json:"createTime"`
	CountTime  interface{} `json:"countTime"`
}
