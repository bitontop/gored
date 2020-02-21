package bittrex

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"

	"github.com/bitontop/gored/exchange"
)

type JsonResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

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

type PairsData []struct {
	MarketCurrency     string      `json:"MarketCurrency"`
	BaseCurrency       string      `json:"BaseCurrency"`
	MarketCurrencyLong string      `json:"MarketCurrencyLong"`
	BaseCurrencyLong   string      `json:"BaseCurrencyLong"`
	MinTradeSize       float64     `json:"MinTradeSize"`
	MarketName         string      `json:"MarketName"`
	IsActive           bool        `json:"IsActive"`
	Created            string      `json:"Created"`
	Notice             interface{} `json:"Notice"`
	IsSponsored        interface{} `json:"IsSponsored"`
	LogoURL            string      `json:"LogoUrl"`
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

type OrderBook struct {
	Buy  []exchange.Order `json:"buy"`
	Sell []exchange.Order `json:"sell"`
}

type TradeHistory []struct {
	ID        int     `json:"Id"`
	TimeStamp string  `json:"TimeStamp"`
	Quantity  float64 `json:"Quantity"`
	Price     float64 `json:"Price"`
	Total     float64 `json:"Total"`
	FillType  string  `json:"FillType"`
	OrderType string  `json:"OrderType"`
	UUID      string  `json:"Uuid"`
}
