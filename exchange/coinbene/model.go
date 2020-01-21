package coinbene

import (
	"encoding/json"
	"time"
)

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct {
	Code    int             `json:"code"`
	Message interface{}     `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type PairsData []struct {
	Symbol           string `json:"symbol"`
	BaseAsset        string `json:"baseAsset"`
	QuoteAsset       string `json:"quoteAsset"`
	PricePrecision   string `json:"pricePrecision"`
	AmountPrecision  string `json:"amountPrecision"`
	TakerFeeRate     string `json:"takerFeeRate"`
	MakerFeeRate     string `json:"makerFeeRate"`
	MinAmount        string `json:"minAmount"`
	PriceFluctuation string `json:"priceFluctuation"`
	Site             string `json:"site"`
}

type OrderBook struct {
	Asks      [][]string `json:"asks"`
	Bids      [][]string `json:"bids"`
	Timestamp time.Time  `json:"timestamp"`
}

type AccountBalances []struct {
	Asset         string `json:"asset"`
	Available     string `json:"available"`
	RrozenBalance string `json:"rrozenBalance"`
	TotalBalance  string `json:"totalBalance"`
}

// TODO

type OrderBookDetail []struct {
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
}

type PlaceOrder struct {
	Orderid     string `json:"orderid"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Timestamp   int64  `json:"timestamp"`
}

type OrderStatus struct {
	Order struct {
		Createtime     int64  `json:"createtime"`
		Filledamount   string `json:"filledamount"`
		Filledquantity string `json:"filledquantity"`
		Orderid        string `json:"orderid"`
		Orderquantity  string `json:"orderquantity"`
		Orderstatus    string `json:"orderstatus"`
		Price          string `json:"price"`
		Symbol         string `json:"symbol"`
		Type           string `json:"type"`
	} `json:"order"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Timestamp   int64  `json:"timestamp"`
}

type Withdraw struct {
	Status     string `json:"status"`
	Timestamp  int64  `json:"timestamp"`
	WithdrawID int    `json:"withdrawId"`
}
