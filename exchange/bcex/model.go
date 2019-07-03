package bcex

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

type JsonResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"msg"`
	Data    json.RawMessage `json:"data"`
}

type CoinsData struct {
	CalcPrecision  int    `json:"calc_precision"`
	BizPrecision   int    `json:"biz_precision"`
	PhyPrecision   int    `json:"phy_precision"`
	OrderAmountMin string `json:"order_amount_min"`
	OrderAmountMax string `json:"order_amount_max"`
}

type PairsData struct {
	Main      map[string][]*PairDetail `json:"main"`
	Hot       []*PairDetail            `json:"hot"`
	ByPercent []*PairDetail            `json:"by_percent"`
	ByTime    []*PairDetail            `json:"by_time"`
}

type PairDetail struct {
	ID            int         `json:"id"`
	Vol           interface{} `json:"vol"`
	Market        string      `json:"market"`
	MarketAs      string      `json:"market_as"`
	Token         string      `json:"token"`
	TokenAs       string      `json:"token_as"`
	Amount        interface{} `json:"amount"`
	MaxPrice      string      `json:"max_price"`
	MinPrice      string      `json:"min_price"`
	OpenPrice     string      `json:"open_price"`
	IsUp          string      `json:"is_up"`
	Last          string      `json:"last"`
	PrevPrice     string      `json:"prev_price"`
	PPrecision    string      `json:"p_precision"`
	NPrecision    string      `json:"n_precision"`
	LatestPrice   string      `json:"latest_price"`
	PercentChange string      `json:"percent_change"`
	Currency      struct {
		CNY string `json:"CNY"`
		KRW string `json:"KRW"`
	} `json:"currency"`
}

type OrderBook struct {
	Bids [][]string `json:"bids"`
	Asks [][]string `json:"asks"`
}

type AccountBalances struct {
	Count int `json:"count"`
	Data  []struct {
		UserID  int    `json:"user_id"`
		OrgID   int    `json:"org_id"`
		Token   string `json:"token"`
		TokenAs string `json:"token_as"`
		Usable  string `json:"usable"`
		Locked  string `json:"locked"`
		Total   string `json:"total"`
	} `json:"data"`
}

//---

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
	ConditionTarget            string
}
