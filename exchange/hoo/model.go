package hoo

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

type PairsData []struct {
	Amount string `json:"amount"`
	AmtNum int    `json:"amt_num"`
	Change string `json:"change"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Price  string `json:"price"`
	QtyNum int    `json:"qty_num"`
	Symbol string `json:"symbol"`
	Volume string `json:"volume"`
}

type OrderBook struct {
	Bids []struct {
		Price    string `json:"price"`
		Quantity string `json:"quantity"`
	} `json:"bids"`
	Asks []struct {
		Price    string `json:"price"`
		Quantity string `json:"quantity"`
	} `json:"asks"`
}

type AccountBalances []struct {
	Amount string `json:"amount"`
	Symbol string `json:"symbol"`
	Freeze string `json:"freeze"`
}

type PlaceOrder struct {
	OrderID string `json:"order_id"`
	TradeNo string `json:"trade_no"`
}

type OrderStatus struct {
	CreateAt   int64         `json:"create_at"`
	Fee        string        `json:"fee"`
	MatchAmt   string        `json:"match_amt"`
	MatchPrice string        `json:"match_price"`
	MatchQty   string        `json:"match_qty"`
	OrderID    string        `json:"order_id"`
	OrderType  int           `json:"order_type"`
	Price      string        `json:"price"`
	Quantity   string        `json:"quantity"`
	Side       int           `json:"side"`
	Status     int           `json:"status"`
	Ticker     string        `json:"ticker"`
	TradeNo    string        `json:"trade_no"`
	Trades     []interface{} `json:"trades"`
}
