package latoken

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

type JsonResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

/* type JsonResponse struct {
	Error struct {
		Message    string `json:"message"`
		ErrorType  string `json:"errorType"`
		StatusCode int    `json:"statusCode"`
	} `json:"error"`
} */

type CoinsData []struct {
	CurrencyID int     `json:"currencyId"`
	Symbol     string  `json:"symbol"`
	Name       string  `json:"name"`
	Precission int     `json:"precission"`
	Type       string  `json:"type"`
	Fee        float64 `json:"fee"`
}

type PairsData []struct {
	PairID          int     `json:"pairId"`
	Symbol          string  `json:"symbol"`
	BaseCurrency    string  `json:"baseCurrency"`
	QuotedCurrency  string  `json:"quotedCurrency"`
	MakerFee        float64 `json:"makerFee"`
	TakerFee        float64 `json:"takerFee"`
	PricePrecision  int     `json:"pricePrecision"`
	AmountPrecision int     `json:"amountPrecision"`
	MinQty          float64 `json:"minQty"`
}

type OrderBook struct {
	PairID int     `json:"pairId"`
	Symbol string  `json:"symbol"`
	Spread float64 `json:"spread"`
	Asks   []struct {
		Price  float64 `json:"price"`
		Amount float64 `json:"amount"`
	} `json:"asks"`
	Bids []struct {
		Price  float64 `json:"price"`
		Amount float64 `json:"amount"`
	} `json:"bids"`
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
	ConditionTarget            string
}
