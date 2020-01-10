package homiex

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

// TODO

type AccountBalances []struct {
	Currency      string  `json:"Currency"`
	Balance       float64 `json:"Balance"`
	Available     float64 `json:"Available"`
	Pending       float64 `json:"Pending"`
	CryptoAddress string  `json:"CryptoAddress"`
	Requested     bool    `json:"Requested"`
	Uuid          string  `json:"Uuid"`
}

type Uuid struct {
	Id string `json:"uuid"`
}

type PlaceOrder struct {
	AccountId                  string
	OrderUuid                  string `json:"OrderUuid"`
	Exchange                   string `json:"Exchange"`
	Type                       string
	Quantity                   float64 `json:"Quantity"`
	QuantityRemaining          float64 `json:"QuantityRemaining"`
	Limit                      float64 `json:"Limit"`
	Reserved                   float64
	ReserveRemaining           float64
	CommissionReserved         float64
	CommissionReserveRemaining float64
	CommissionPaid             float64
	Price                      float64 `json:"Price"`
	PricePerUnit               float64 `json:"PricePerUnit"`
	Opened                     string
	Closed                     string
	IsOpen                     bool
	Sentinel                   string
	CancelInitiated            bool
	ImmediateOrCancel          bool
	IsConditional              bool
	Condition                  string
	ConditionTarget            float64
}
