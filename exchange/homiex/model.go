package homiex

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

type JsonResponse struct {
	// error section
	Code int    `json:"code"`
	Msg  string `json:"msg"`

	Success bool            `json:"success"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

type PairsData struct {
	Timezone      string        `json:"timezone"`
	ServerTime    string        `json:"serverTime"`
	BrokerFilters []interface{} `json:"brokerFilters"`
	Symbols       []struct {
		Filters []struct {
			MinPrice    string `json:"minPrice,omitempty"`
			MaxPrice    string `json:"maxPrice,omitempty"`
			TickSize    string `json:"tickSize,omitempty"`
			FilterType  string `json:"filterType"`
			MinQty      string `json:"minQty,omitempty"`
			MaxQty      string `json:"maxQty,omitempty"`
			StepSize    string `json:"stepSize,omitempty"`
			MinNotional string `json:"minNotional,omitempty"`
		} `json:"filters"`
		ExchangeID         string `json:"exchangeId"`
		Symbol             string `json:"symbol"`
		SymbolName         string `json:"symbolName"`
		Status             string `json:"status"`
		BaseAsset          string `json:"baseAsset"`
		BaseAssetPrecision string `json:"baseAssetPrecision"`
		QuoteAsset         string `json:"quoteAsset"`
		QuotePrecision     string `json:"quotePrecision"`
		IcebergAllowed     bool   `json:"icebergAllowed"`
	} `json:"symbols"`
	// Aggregates []interface{} `json:"aggregates"`
	// RateLimits []struct {
	// 	RateLimitType string `json:"rateLimitType"`
	// 	Interval      string `json:"interval"`
	// 	IntervalUnit  int    `json:"intervalUnit"`
	// 	Limit         int    `json:"limit"`
	// } `json:"rateLimits"`
	// Options   []interface{} `json:"options"`
	// Contracts []struct {
	// 	Filters []struct {
	// 		MinPrice    string `json:"minPrice,omitempty"`
	// 		MaxPrice    string `json:"maxPrice,omitempty"`
	// 		TickSize    string `json:"tickSize,omitempty"`
	// 		FilterType  string `json:"filterType"`
	// 		MinQty      string `json:"minQty,omitempty"`
	// 		MaxQty      string `json:"maxQty,omitempty"`
	// 		StepSize    string `json:"stepSize,omitempty"`
	// 		MinNotional string `json:"minNotional,omitempty"`
	// 	} `json:"filters"`
	// 	ExchangeID          string `json:"exchangeId"`
	// 	Symbol              string `json:"symbol"`
	// 	SymbolName          string `json:"symbolName"`
	// 	Status              string `json:"status"`
	// 	BaseAsset           string `json:"baseAsset"`
	// 	BaseAssetPrecision  string `json:"baseAssetPrecision"`
	// 	QuoteAsset          string `json:"quoteAsset"`
	// 	QuoteAssetPrecision string `json:"quoteAssetPrecision"`
	// 	IcebergAllowed      bool   `json:"icebergAllowed"`
	// 	Inverse             bool   `json:"inverse"`
	// 	Index               string `json:"index"`
	// 	MarginToken         string `json:"marginToken"`
	// 	MarginPrecision     string `json:"marginPrecision"`
	// 	ContractMultiplier  string `json:"contractMultiplier"`
	// 	Underlying          string `json:"underlying"`
	// 	RiskLimits          []struct {
	// 		RiskLimitID   string `json:"riskLimitId"`
	// 		Quantity      string `json:"quantity"`
	// 		InitialMargin string `json:"initialMargin"`
	// 		MaintMargin   string `json:"maintMargin"`
	// 	} `json:"riskLimits"`
	// } `json:"contracts"`
}

type OrderBook struct {
	Time int64      `json:"time"`
	Bids [][]string `json:"bids"`
	Asks [][]string `json:"asks"`
}

type AccountBalances struct {
	CanTrade    bool `json:"canTrade"`
	CanWithdraw bool `json:"canWithdraw"`
	CanDeposit  bool `json:"canDeposit"`
	UpdateTime  int  `json:"updateTime"`
	Balances    []struct {
		Asset   string `json:"asset"`
		Total   string `json:"total"`
		Free    string `json:"free"`
		Locked  string `json:"locked"`
		AssetID string `json:"assetId"`
	} `json:"balances"`
}

type PlaceOrder struct {
	OrderID       string `json:"orderId"`
	ClientOrderID string `json:"clientOrderId"`
}

type OrderStatus struct {
	Symbol              string `json:"symbol"`
	OrderID             string `json:"orderId"`
	ClientOrderID       string `json:"clientOrderId"`
	Price               string `json:"price"`
	OrigQty             string `json:"origQty"`
	ExecutedQty         string `json:"executedQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status              string `json:"status"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	Side                string `json:"side"`
	StopPrice           string `json:"stopPrice"`
	IcebergQty          string `json:"icebergQty"`
	Time                string `json:"time"`
	UpdateTime          string `json:"updateTime"`
	IsWorking           bool   `json:"isWorking"`
}

type CancelOrder struct {
	Symbol        string `json:"symbol"`
	ClientOrderID string `json:"clientOrderId"`
	OrderID       string `json:"orderId"`
	Status        string `json:"status"`
}

type TradeHistory []struct {
	Price        string `json:"price"`
	Time         int64  `json:"time"`
	Qty          string `json:"qty"`
	IsBuyerMaker bool   `json:"isBuyerMaker"`
}

type Withdraw struct {
	Ret int `json:"ret"`
}
