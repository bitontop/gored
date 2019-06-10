package huobidm

import "encoding/json"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct {
	Status string          `json:"status"`
	Data   json.RawMessage `json:"data"`
	Ts     int64           `json:"ts"`
}

type JsonResponse2 struct {
	Ch     string          `json:"ch"`
	Status string          `json:"status"`
	Tick   json.RawMessage `json:"tick"`
	Ts     int64           `json:"ts"`
}

/********** Public API Structure**********/
type ContractsData []struct {
	Symbol         string  `json:"symbol"`
	ContractCode   string  `json:"contract_code"`
	ContractType   string  `json:"contract_type"`
	ContractSize   float64 `json:"contract_size"`
	PriceTick      float64 `json:"price_tick"`
	DeliveryDate   string  `json:"delivery_date"`
	CreateDate     string  `json:"create_date"`
	ContractStatus int     `json:"contract_status"`
}

type PairsData []struct {
	Symbol      string  `json:"symbol"`
	Status      string  `json:"status"`
	BaseAsset   string  `json:"baseAsset"`
	QuoteAsset  string  `json:"quoteAsset"`
	MakerFee    float64 `json:"makerFee"`
	TakerFee    float64 `json:"takerFee"`
	PriceFilter float64 `json:"priceFilter"`
	LotSize     float64 `json:"lotSize"`
}

type OrderBook struct {
	Bids [][]float64 `json:"bids"`
	Asks [][]float64 `json:"asks"`
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
