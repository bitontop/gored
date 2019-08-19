package bitmax

import "encoding/json"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Status  string          `json:"status"`
	Email   string          `json:"email"`
	Data    json.RawMessage `json:"data"`
}

type CoinsData []struct {
	AssetCode        string  `json:"assetCode"`
	AssetName        string  `json:"assetName"`
	PrecisionScale   int     `json:"precisionScale"`
	NativeScale      int     `json:"nativeScale"`
	WithdrawalFee    float64 `json:"withdrawalFee"`
	MinWithdrawalAmt float64 `json:"minWithdrawalAmt"`
	Status           string  `json:"status"`
}

type PairsData []struct {
	Symbol        string `json:"symbol"`
	BaseAsset     string `json:"baseAsset"`
	QuoteAsset    string `json:"quoteAsset"`
	PriceScale    int    `json:"priceScale"`
	QtyScale      int    `json:"qtyScale"`
	NotionalScale int    `json:"notionalScale"`
	MinQty        string `json:"minQty"`
	MaxQty        string `json:"maxQty"`
	MinNotional   string `json:"minNotional"`
	MaxNotional   string `json:"maxNotional"`
	Status        string `json:"status"`
	MiningStatus  string `json:"miningStatus"`
}

type OrderBook struct {
	M      string     `json:"m"`
	S      string     `json:"s"`
	Seqnum int        `json:"seqnum"`
	Asks   [][]string `json:"asks"`
	Bids   [][]string `json:"bids"`
}

type AccountGroup struct {
	AccountGroup int `json:"accountGroup"`
}

type AccountBalances []struct {
	AssetCode       string `json:"assetCode"`
	AssetName       string `json:"assetName"`
	TotalAmount     string `json:"totalAmount"`
	AvailableAmount string `json:"availableAmount"`
	InOrderAmount   string `json:"inOrderAmount"`
}

type Withdrawal struct {
	Msg    string `json:"msg"`
	Status string `json:"status"`
}

type PlaceOrder struct {
	Coid    string `json:"coid"`
	Action  string `json:"action"`
	Success bool   `json:"success"`
}

type OrderStatus struct {
	Time       int64  `json:"time"`
	Coid       string `json:"coid"`
	Symbol     string `json:"symbol"`
	BaseAsset  string `json:"baseAsset"`
	QuoteAsset string `json:"quoteAsset"`
	Side       string `json:"side"`
	OrderPrice string `json:"orderPrice"`
	StopPrice  string `json:"stopPrice"`
	OrderQty   string `json:"orderQty"`
	Filled     string `json:"filled"`
	Fee        string `json:"fee"`
	FeeAsset   string `json:"feeAsset"`
	Status     string `json:"status"`
}
