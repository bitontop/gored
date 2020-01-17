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

// TODO

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
	ConditionTarget            float64
}
