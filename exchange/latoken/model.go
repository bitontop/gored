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

type AccountBalances [][]struct {
	CurrencyID int     `json:"currencyId"`
	Symbol     string  `json:"symbol"`
	Name       string  `json:"name"`
	Amount     float64 `json:"amount"`
	Available  float64 `json:"available"`
	Frozen     int     `json:"frozen"`
	Pending    int     `json:"pending"`
}

type PlaceOrder struct {
	OrderID   string  `json:"orderId"`
	CliOrdID  string  `json:"cliOrdId"`
	PairID    int     `json:"pairId"`
	Symbol    string  `json:"symbol"`
	Side      string  `json:"side"`
	OrderType string  `json:"orderType"`
	Price     float64 `json:"price"`
	Amount    float64 `json:"amount"`
}

type OrderStatus struct {
	OrderID         string  `json:"orderId"`
	CliOrdID        string  `json:"cliOrdId"`
	PairID          int     `json:"pairId"`
	Symbol          string  `json:"symbol"`
	Side            string  `json:"side"`
	OrderType       string  `json:"orderType"`
	Price           float64 `json:"price"`
	Amount          float64 `json:"amount"`
	OrderStatus     string  `json:"orderStatus"`
	ExecutedAmount  float64 `json:"executedAmount"`
	ReaminingAmount float64 `json:"reaminingAmount"`
	TimeCreated     int64   `json:"timeCreated"`
	TimeFilled      int     `json:"timeFilled"`
}

type TestBuy struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	OrderID string `json:"orderId"`
}
