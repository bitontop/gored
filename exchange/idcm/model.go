package idcm

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

// type JsonResponse struct {
// 	Success bool            `json:"success"`
// 	Message string          `json:"message"`
// 	Result  json.RawMessage `json:"result"`
// }

type JsonResponse struct {
	Result int             `json:"result"`
	Code   string          `json:"code"`
	Data   json.RawMessage `json:"data"`
	// old
	// Success bool   `json:"success"`
	// Message string `json:"message"`
}

type PairsData struct {
	Data []struct {
		TradePairID   string  `json:"TradePairID"`
		TradePairCode string  `json:"TradePairCode"`
		LastPrice     float64 `json:"LastPrice"`
		Change        float64 `json:"Change"`
		Rose          float64 `json:"Rose"`
		Volume        float64 `json:"Volume"`
		High          float64 `json:"High"`
		Low           float64 `json:"Low"`
		Open          float64 `json:"Open"`
		Close         float64 `json:"Close"`
		Turnover      float64 `json:"Turnover"`
		Sort          int     `json:"Sort"`
		PriceDigit    int     `json:"PriceDigit"`
		QuantityDigit int     `json:"QuantityDigit"`
	} `json:"Data"`
	NeedLang   bool        `json:"NeedLang"`
	Status     bool        `json:"Status"`
	Msg        interface{} `json:"Msg"`
	URL        interface{} `json:"Url"`
	StatusCode string      `json:"StatusCode"`
	Extra      interface{} `json:"Extra"`
}

type OrderBook struct {
	Asks []struct {
		Symbol string  `json:"symbol"`
		Price  float64 `json:"price"`
		Amount float64 `json:"amount"`
	} `json:"asks"`
	Bids []struct {
		Symbol string  `json:"symbol"`
		Price  float64 `json:"price"`
		Amount float64 `json:"amount"`
	} `json:"bids"`
}

type AccountBalances []struct {
	Code    string  `json:"code"`
	Free    float64 `json:"free"`
	Freezed float64 `json:"freezed"`
}

type PlaceOrder struct {
	Orderid string `json:"orderid"`
}

type OrderStatus []struct {
	Orderid        string  `json:"orderid"`
	Symbol         string  `json:"symbol"`
	Price          float64 `json:"price"`
	Avgprice       float64 `json:"avgprice"`
	Side           int     `json:"side"`
	Type           int     `json:"type"`
	Timestamp      string  `json:"timestamp"`
	Amount         int     `json:"amount"`
	Executedamount int     `json:"executedamount"`
	Status         int     `json:"status"`
}

type Withdraw struct {
	WithdrawID string `json:"WithdrawID"`
}
