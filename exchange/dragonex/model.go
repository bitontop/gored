package dragonex

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

type JsonResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Ok   bool            `json:"ok"`
	Data json.RawMessage `json:"data"`
}

// public
type CoinsData []struct {
	CoinID int    `json:"coin_id"`
	Code   string `json:"code"`
}

type PairsData struct {
	Columns []string        `json:"columns"`
	List    [][]interface{} `json:"list"`
}

type OrderBook struct {
	Buys []struct {
		Price  string `json:"price"`
		Volume string `json:"volume"`
	} `json:"buys"`
	Sells []struct {
		Price  string `json:"price"`
		Volume string `json:"volume"`
	} `json:"sells"`
}

// private
type Token struct {
	ExpireTime int    `json:"expire_time"`
	Token      string `json:"token"`
}

type AccountBalances []struct {
	Code   string `json:"code"`
	CoinID int    `json:"coin_id"`
	Frozen string `json:"frozen"`
	Volume string `json:"volume"`
}

type PlaceOrder struct {
	OrderID     string `json:"order_id"`
	Price       string `json:"price"`
	Status      int    `json:"status"`
	Timestamp   int    `json:"timestamp"`
	TradeVolume string `json:"trade_volume"`
	Volume      string `json:"volume"`
}

type OrderStatus struct {
	OrderID      string `json:"order_id"`
	OrderType    int    `json:"order_type"`
	Price        string `json:"price"`
	Status       int    `json:"status"`
	SymbolID     int    `json:"symbol_id"`
	Timestamp    int    `json:"timestamp"`
	TradeVolume  string `json:"trade_volume"`
	Volume       string `json:"volume"`
	ActualAmount string `json:"actual_amount"`
	ActualFee    string `json:"actual_fee"`
}
