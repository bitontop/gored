package bitbns

import "encoding/json"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct {
	Status int             `json:"status"`
	Error  interface{}     `json:"error"`
	Code   int             `json:"code"`
	Data   json.RawMessage `json:"data"`
	// old
	Success bool   `json:"success"`
	Message string `json:"message"`
}

/********** Public API Structure**********/
type CoinsData map[string]*CoinsDetail
type CoinsDetail struct {
	HighestBuyBid   float64 `json:"highest_buy_bid"`
	LowestSellBid   float64 `json:"lowest_sell_bid"`
	LastTradedPrice float64 `json:"last_traded_price"`
	YesPrice        float64 `json:"yes_price"`
	Volume          struct {
		Max    float64 `json:"max"`
		Min    float64 `json:"min"`
		Volume float64 `json:"volume"`
	} `json:"volume"`
}

// type PairsData []struct {
// 	Symbol      string  `json:"symbol"`
// 	Status      string  `json:"status"`
// 	BaseAsset   string  `json:"baseAsset"`
// 	QuoteAsset  string  `json:"quoteAsset"`
// 	MakerFee    float64 `json:"makerFee"`
// 	TakerFee    float64 `json:"takerFee"`
// 	PriceFilter float64 `json:"priceFilter"`
// 	LotSize     float64 `json:"lotSize"`
// }

type OrderBook []struct {
	Rate float64 `json:"rate"`
	Btc  float64 `json:"btc"`
}

// type OrderBook struct {
// 	Bids [][]float64 `json:"bids"`
// 	Asks [][]float64 `json:"asks"`
// }

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
