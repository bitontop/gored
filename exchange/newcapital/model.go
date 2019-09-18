package newcapital

import "encoding/json"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

/********** Public API Structure**********/
type Coin struct {
	Symbol             string      `json:"symbol"`
	BaseAsset          string      `json:"baseAsset"`
	BaseAssetPrecision int         `json:"baseAssetPrecision"`
	BaseAssetName      string      `json:"baseAssetName"`
	QuoteAsset         string      `json:"quoteAsset"`
	QuotePrecision     int         `json:"quotePrecision"`
	QuoteAssetName     string      `json:"quoteAssetName"`
	OrderTypes         interface{} `json:"orderTypes"`
}

type CoinsData struct {
	Timezone   string            `json:"timezone"`
	ServerTime int               `json:"serverTime"`
	Symbols    []Coin            `json:"symbols"`
	Volumes    map[string]string `json:"24h_volume"`
	UsdPrice   map[string]string `json:"usd_price"`
}

type Pair struct {
	Symbol string `json:"symbol"`
	// PriceChange        string      `json:"priceChange"`
	PriceChangePercent string      `json:"priceChangePercent"`
	LastPrice          string      `json:"lastPrice"`
	BidPrice           interface{} `json:"bidPrice"`
	AskPrice           interface{} `json:"askPrice"`
	OpenPrice          string      `json:"openPrice"`
	HighPrice          string      `json:"highPrice"`
	LowPrice           string      `json:"lowPrice"`
	Volume             string      `json:"volume"`
	QuoteVolume        string      `json:"quoteVolume"`
	OpenTime           int         `json:"openTime"`
	CloseTime          int         `json:"closeTime"`
	FirstId            int         `json:"firstId"`
	LastId             int         `json:"lastId"`
	Count              int         `json:"count"`
}

type PairsData []Pair

type OrderBook struct {
	Bids [][]string `json:"bids"`
	Asks [][]string `json:"asks"`
}

/********** Private API Structure**********/
type AccountBalances []struct {
	Asset     string  `json:"asset"`
	Total     float64 `json:"total"`
	Available float64 `json:"available"`
	Locked    float64 `json:"locked"`
}

type WithdrawResponse struct {
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

type PlaceOrder struct {
	Symbol       string `json:"symbol"`
	OrderID      string `json:"orderId"`
	Side         string `json:"side"`
	Type         string `json:"type"`
	Price        string `json:"price"`
	AveragePrice string `json:"executedQty"`
	OrigQty      string `json:"origQty"`
	ExecutedQty  string `json:"executedQty"`
	Status       string `json:"status"`
	TimeInForce  string `json:"timeInForce"`
}
