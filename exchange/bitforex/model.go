package bitforex

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

type JsonResponse struct {
	Success bool            `json:"success"`
	Code    string          `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type PairsData []struct {
	AmountPrecision int     `json:"amountPrecision"`
	MinOrderAmount  float64 `json:"minOrderAmount"`
	PricePrecision  int     `json:"pricePrecision"`
	Symbol          string  `json:"symbol"`
}

type OrderBook struct {
	Asks []struct {
		Amount float64 `json:"amount"`
		Price  float64 `json:"price"`
	} `json:"asks"`
	Bids []struct {
		Amount float64 `json:"amount"`
		Price  float64 `json:"price"`
	} `json:"bids"`
}

type AccountBalances []struct {
	Fix      string `json:"fix"`
	Frozen   string `json:"frozen"`
	Active   string `json:"active"`
	Currency string `json:"currency"`
}

type PlaceOrder struct {
	OrderID int `json:"orderId"`
}

type OrderStatus struct {
	AvgPrice    int         `json:"avgPrice"`
	CreateTime  int64       `json:"createTime"`
	DealAmount  interface{} `json:"dealAmount"`
	LastTime    int64       `json:"lastTime"`
	OrderAmount interface{} `json:"orderAmount"`
	OrderID     int         `json:"orderId"`
	OrderPrice  interface{} `json:"orderPrice"`
	OrderState  int         `json:"orderState"`
	Symbol      string      `json:"symbol"`
	TradeFee    int         `json:"tradeFee"`
	TradeType   int         `json:"tradeType"`
}

type CancelOrder struct {
	Cancel bool `json:"data"`
}

/* type PairsData []struct {
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
} */

/* type AccountBalances []struct {
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
} */

/* type OrderBook struct {
	Buy  []exchange.Order `json:"buy"`
	Sell []exchange.Order `json:"sell"`
} */
