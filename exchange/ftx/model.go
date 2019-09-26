package ftx

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"

	"github.com/bitontop/gored/exchange"
)

type JsonResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

type CoinsData []struct { // not useful
	ID         string `json:"id"`
	Name       string `json:"name"`
	Index      bool   `json:"index,omitempty"`
	Collateral bool   `json:"collateral,omitempty"`
	Underlying string `json:"underlying,omitempty"`
}

type PairsData []struct {
	Ask             float64     `json:"ask"`
	Bid             float64     `json:"bid"`
	Change1H        float64     `json:"change1h"`
	Change24H       float64     `json:"change24h"`
	Description     string      `json:"description"`
	Enabled         bool        `json:"enabled"`
	Expired         bool        `json:"expired"`
	Expiry          interface{} `json:"expiry"`
	Index           float64     `json:"index"`
	IndexAdjustment float64     `json:"indexAdjustment"`
	Last            float64     `json:"last"`
	LowerBound      float64     `json:"lowerBound"`
	Mark            float64     `json:"mark"`
	Name            string      `json:"name"`
	Perpetual       bool        `json:"perpetual"`
	PostOnly        bool        `json:"postOnly"`
	PriceIncrement  float64     `json:"priceIncrement"`
	SizeIncrement   float64     `json:"sizeIncrement"`
	Type            string      `json:"type"`
	Underlying      string      `json:"underlying"`
	UpperBound      float64     `json:"upperBound"`
	VolumeUsd24H    float64     `json:"volumeUsd24h"`
}

// ========

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

type OrderBook struct {
	Buy  []exchange.Order `json:"buy"`
	Sell []exchange.Order `json:"sell"`
}
