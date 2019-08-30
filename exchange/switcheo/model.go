package switcheo

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
type ListItem struct {
	Start  int  `json:"start"`
	End    int  `json:"end"`
	Paused bool `json:"paused"`
}

type ListData map[string]ListItem

type Token struct {
	Symbol           string      `json:"symbol"`
	Name             string      `json:"name"`
	Type             string      `json:"type"`
	Hash             string      `json:"hash"`
	Decimals         int         `json:"decimals"`
	TransferDecimals int         `json:"transfer_decimals"`
	Precision        int         `json:"precision"`
	MinimumQuantity  string      `json:"minimum_quantity"`
	TradingActive    bool        `json:"trading_active"`
	IsStablecoin     bool        `json:"is_stablecoin"`
	StablecoinType   interface{} `json:"stablecoin_type"`
	Active           bool        `json:"active"`
	ListingInfo      ListData    `json:"listing_info"`
}

type CoinsData map[string]Token

type Pair struct {
	Name             string `json:"name"`
	Precision        int    `json:"precision"`
	BaseAssetName    string `json:"baseAssetName"`
	BaseAssetSymbol  string `json:"baseAssetSymbol"`
	BaseContract     string `json:"baseContract"`
	QuoteAssetName   string `json:"quoteAssetName"`
	QuoteAssetSymbol string `json:"quoteAssetSymbol"`
	QuoteContract    string `json:"quoteContract"`
}

type PairsData []Pair

type OrderBook struct {
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
}

type OrderBooks struct {
	Asks []OrderBook `json:"asks"`
	Bids []OrderBook `json:"bids"`
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
