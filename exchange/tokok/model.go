package tokok

import "encoding/json"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct {
	Result bool            `json:"result"`
	Code   int             `json:"code"`
	Data   json.RawMessage `json:"data"`
}

type PairsData []struct {
	Symbol              string `json:"symbol"`
	Status              string `json:"status"`
	BaseAsset           string `json:"baseAsset"`
	BaseAssetPrecision  int    `json:"baseAssetPrecision"`
	QuoteAsset          string `json:"quoteAsset"`
	QuoteAssetPrecision int    `json:"quoteAssetPrecision"`
	MinOrderAmount      string `json:"minOrderAmount"`
	MaxOrderAmount      string `json:"maxOrderAmount"`
}

type OrderBook struct {
	Bids [][]interface{} `json:"bids"`
	Asks [][]interface{} `json:"asks"`
}

type AccountBalances []struct {
	HotMoney  string `json:"hotMoney"`
	ColdMoney string `json:"coldMoney"`
	CoinCode  string `json:"coinCode"`
}

type PlaceOrder struct {
	Data string `json:"data"`
}

type OrderStatus struct {
	EntrustNum          string `json:"entrustNum"`
	Type                int    `json:"type"`
	Status              int    `json:"status"`
	EntrustPrice        string `json:"entrustPrice"`
	EntrustCount        string `json:"entrustCount"`
	EntrustSum          string `json:"entrustSum"`
	SurplusEntrustCount string `json:"surplusEntrustCount"`
	EntrustTimeLong     int64  `json:"entrustTime_long"`
	TransactionFee      string `json:"transactionFee"`
	ProcessedPrice      string `json:"processedPrice"`
	OpenTokFee          int    `json:"openTokFee"`
	Symbol              string `json:"symbol"`
}

type CancelOrder struct {
	Result bool `json:"result"`
	Code   int  `json:"code"`
}
