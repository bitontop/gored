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

type PlaceOrder struct {
	OrderID  string `json:"orderId"`
	ClientID string `json:"clientId"`
}

type OrderStatus struct {
	OrderID        string    `json:"orderId"`
	BaseAsset      string    `json:"baseAsset"`
	QuoteAsset     string    `json:"quoteAsset"`
	OrderDirection string    `json:"orderDirection"`
	Quantity       string    `json:"quantity"`
	Amount         string    `json:"amount"`
	FilledAmount   string    `json:"filledAmount"`
	AvgPrice       string    `json:"avgPrice"`
	OrderPrice     string    `json:"orderPrice"`
	TakerFeeRate   string    `json:"takerFeeRate"`
	MakerFeeRate   string    `json:"makerFeeRate"`
	OrderStatus    string    `json:"orderStatus"`
	OrderTime      time.Time `json:"orderTime"`
	TotalFee       string    `json:"totalFee"`
}

type Withdraw struct {
	Code int `json:"Code"`
	Data struct {
		ID      string `json:"Id"`
		Amount  string `json:"Amount"`
		Asset   string `json:"Asset"`
		Address string `json:"Address"`
		Tag     string `json:"Tag"`
		Chain   string `json:"Chain"`
	} `json:"Data"`
}

// TODO

// type Withdraw struct {
// 	Status     string `json:"status"`
// 	Timestamp  int64  `json:"timestamp"`
// 	WithdrawID int    `json:"withdrawId"`
// }

type TradeHistory [][]Trade

type Trade struct {
	Symbol    string    `json:"symbol"`
	Price     string    `json:"price"`
	Volume    string    `json:"volume"`
	Direction string    `json:"direction"`
	TradeTime time.Time `json:"tradeTime"`
}
