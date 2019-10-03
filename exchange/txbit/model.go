package txbit

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
	"time"
)

type JsonResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

type CoinsData []struct {
	Currency        string  `json:"Currency"`
	CurrencyLong    string  `json:"CurrencyLong"`
	MinConfirmation int     `json:"MinConfirmation"`
	CoinType        string  `json:"CoinType"`
	TxFee           float64 `json:"TxFee"`
	IsActive        bool    `json:"IsActive"`
	BaseAddress     string  `json:"BaseAddress"`
}

type PairsData []struct {
	MarketCurrency     string    `json:"MarketCurrency"`
	BaseCurrency       string    `json:"BaseCurrency"`
	MarketCurrencyLong string    `json:"MarketCurrencyLong"`
	BaseCurrencyLong   string    `json:"BaseCurrencyLong"`
	MinTradeSize       float64   `json:"MinTradeSize"`
	MarketName         string    `json:"MarketName"`
	IsActive           bool      `json:"IsActive"`
	Created            time.Time `json:"Created"`
}

type OrderBook struct {
	Buy []struct {
		Quantity float64 `json:"Quantity"`
		Rate     float64 `json:"Rate"`
	} `json:"buy"`
	Sell []struct {
		Quantity float64 `json:"Quantity"`
		Rate     float64 `json:"Rate"`
	} `json:"sell"`
}

// Old

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
