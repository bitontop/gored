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

type CoinsData []struct {
	Currency        string      `json:"Currency"`
	CurrencyLong    string      `json:"CurrencyLong"`
	MinConfirmation int         `json:"MinConfirmation"`
	TxFee           float64     `json:"TxFee"`
	IsActive        bool        `json:"IsActive"`
	IsRestricted    bool        `json:"IsRestricted"`
	CoinType        string      `json:"CoinType"`
	BaseAddress     string      `json:"BaseAddress"`
	Notice          interface{} `json:"Notice"`
}
