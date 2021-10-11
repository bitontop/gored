package ftx

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
	"time"
)

type JsonResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

type CoinsData []struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Index      bool   `json:"index,omitempty"`
	Collateral bool   `json:"collateral,omitempty"`
	Underlying string `json:"underlying,omitempty"`
}

type PairsData []struct {
	Ask            float64 `json:"ask"`
	BaseCurrency   string  `json:"baseCurrency"`
	Bid            float64 `json:"bid"`
	Change1H       float64 `json:"change1h,omitempty"`
	Change24H      float64 `json:"change24h,omitempty"`
	ChangeBod      float64 `json:"changeBod,omitempty"`
	Enabled        bool    `json:"enabled"`
	Last           float64 `json:"last"`
	MinProvideSize float64 `json:"minProvideSize"`
	Name           string  `json:"name"`
	Price          float64 `json:"price"`
	PriceIncrement float64 `json:"priceIncrement"`
	QuoteCurrency  string  `json:"quoteCurrency"`
	QuoteVolume24H float64 `json:"quoteVolume24h,omitempty"`
	Restricted     bool    `json:"restricted"`
	SizeIncrement  float64 `json:"sizeIncrement"`
	Type           string  `json:"type"`
	Underlying     string  `json:"underlying"`
	VolumeUsd24H   float64 `json:"volumeUsd24h,omitempty"`
}

type OrderBook struct {
	Asks [][]float64 `json:"asks"`
	Bids [][]float64 `json:"bids"`
}

type RawKline []struct {
	Close     float64   `json:"close"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Open      float64   `json:"open"`
	StartTime time.Time `json:"startTime"`
	Time      float64   `json:"time"`
	Volume    float64   `json:"volume"`
}

type AccountBalances []struct {
	Coin  string  `json:"coin"`
	Free  float64 `json:"free"`
	Total float64 `json:"total"`
}

type AccountAllBalances struct {
	Main []struct {
		Coin  string  `json:"coin"`
		Free  float64 `json:"free"`
		Total float64 `json:"total"`
	} `json:"main"`
	BattleRoyale []struct {
		Coin  string  `json:"coin"`
		Free  float64 `json:"free"`
		Total float64 `json:"total"`
	} `json:"Battle Royale"`
}

type PlaceOrder struct {
	CreatedAt     time.Time   `json:"createdAt"`
	FilledSize    float64     `json:"filledSize"`
	Future        string      `json:"future"`
	ID            int         `json:"id"`
	Market        string      `json:"market"`
	Price         float64     `json:"price"`
	RemainingSize float64     `json:"remainingSize"`
	Side          string      `json:"side"`
	Size          float64     `json:"size"`
	Status        string      `json:"status"`
	Type          string      `json:"type"`
	ReduceOnly    bool        `json:"reduceOnly"`
	Ioc           bool        `json:"ioc"`
	PostOnly      bool        `json:"postOnly"`
	ClientID      interface{} `json:"clientId"`
}

type OrderStatus struct {
	CreatedAt     time.Time   `json:"createdAt"`
	FilledSize    float64     `json:"filledSize"`
	Future        string      `json:"future"`
	ID            int         `json:"id"`
	Market        string      `json:"market"`
	Price         float64     `json:"price"`
	AvgFillPrice  float64     `json:"avgFillPrice"`
	RemainingSize float64     `json:"remainingSize"`
	Side          string      `json:"side"`
	Size          float64     `json:"size"`
	Status        string      `json:"status"`
	Type          string      `json:"type"`
	ReduceOnly    bool        `json:"reduceOnly"`
	Ioc           bool        `json:"ioc"`
	PostOnly      bool        `json:"postOnly"`
	ClientID      interface{} `json:"clientId"`
}

type OpenOrders []struct {
	CreatedAt     time.Time   `json:"createdAt"`
	FilledSize    float64     `json:"filledSize"`
	Future        string      `json:"future"`
	ID            int         `json:"id"`
	Market        string      `json:"market"`
	Price         float64     `json:"price"`
	AvgFillPrice  float64     `json:"avgFillPrice"`
	RemainingSize float64     `json:"remainingSize"`
	Side          string      `json:"side"`
	Size          float64     `json:"size"`
	Status        string      `json:"status"`
	Type          string      `json:"type"`
	ReduceOnly    bool        `json:"reduceOnly"`
	Ioc           bool        `json:"ioc"`
	PostOnly      bool        `json:"postOnly"`
	ClientID      interface{} `json:"clientId"`
}

type CloseOrders []struct {
	AvgFillPrice  float64     `json:"avgFillPrice"`
	ClientID      interface{} `json:"clientId"`
	CreatedAt     time.Time   `json:"createdAt"`
	FilledSize    float64     `json:"filledSize"`
	Future        string      `json:"future"`
	ID            int         `json:"id"`
	Ioc           bool        `json:"ioc"`
	Market        string      `json:"market"`
	PostOnly      bool        `json:"postOnly"`
	Price         float64     `json:"price"`
	ReduceOnly    bool        `json:"reduceOnly"`
	RemainingSize float64     `json:"remainingSize"`
	Side          string      `json:"side"`
	Size          float64     `json:"size"`
	Status        string      `json:"status"`
	Type          string      `json:"type"`

	// HasMoreData bool `json:"hasMoreData"`
}

// need test array or single
type WithdrawHistory []struct {
	Coin    string      `json:"coin"`
	Notes   string      `json:"notes"`
	Address string      `json:"address"`
	Tag     interface{} `json:"tag"`
	Fee     float64     `json:"fee"`
	ID      int         `json:"id"`
	Size    float64     `json:"size"`
	Status  string      `json:"status"`
	Time    time.Time   `json:"time"`
	Txid    string      `json:"txid"`
}

type DepositHistory []struct {
	Coin          string    `json:"coin"`
	Notes         string    `json:"notes"`
	Confirmations int       `json:"confirmations"`
	ConfirmedTime time.Time `json:"confirmedTime"`
	Fee           float64   `json:"fee"`
	ID            int       `json:"id"`
	SentTime      time.Time `json:"sentTime"`
	Size          float64   `json:"size"`
	Status        string    `json:"status"`
	Time          time.Time `json:"time"`
	Txid          string    `json:"txid"`
}

type DepositAddress struct {
	Address string `json:"address"`
	Tag     string `json:"tag"`
}

// =========

type Uuid struct {
	Id string `json:"uuid"`
}

// type OpenPositions []struct {
// 	Future                       string  `json:"future"`
// 	Size                         float64 `json:"size"`
// 	Side                         string  `json:"side"`
// 	NetSize                      float64 `json:"netSize"`
// 	LongOrderSize                float64 `json:"longOrderSize"`
// 	ShortOrderSize               float64 `json:"shortOrderSize"`
// 	Cost                         float64 `json:"cost"`
// 	EntryPrice                   float64 `json:"entryPrice"`
// 	UnrealizedPnl                float64 `json:"unrealizedPnl"`
// 	RealizedPnl                  float64 `json:"realizedPnl"`
// 	InitialMarginRequirement     float64 `json:"initialMarginRequirement"`
// 	MaintenanceMarginRequirement float64 `json:"maintenanceMarginRequirement"`
// 	OpenSize                     float64 `json:"openSize"`
// 	CollateralUsed               float64 `json:"collateralUsed"`
// 	EstimatedLiquidationPrice    float64 `json:"estimatedLiquidationPrice"`
// }
