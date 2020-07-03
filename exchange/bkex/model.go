package bkex

import "encoding/json"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
	Msg  string          `json:"msg"`
}

/********** Public API Structure**********/
type Coin struct {
	CoinType          string  `json:"coinType"`
	MaxWithdrawOneDay float64 `json:"maxWithdrawOneDay"`
	MaxWithdrawSingle float64 `json:"maxWithdrawSingle"`
	MinWithdrawSingle float64 `json:"minWithdrawSingle"`
	SupportDeposit    bool    `json:"supportDeposit"`
	SupportTrade      bool    `json:"supportTrade"`
	SupportWithdraw   bool    `json:"supportWithdraw"`
	WithdrawFee       float64 `json:"withdrawFee"`
}

type Pair struct {
	AmountPrecision  int    `json:"amountPrecision"`
	DefaultPrecision int    `json:"defaultPrecision"`
	Pair             string `json:"pair"`
	SupportTrade     bool   `json:"supportTrade"`
}

type CoinsData []Coin

type PairsData []Pair

type ExchangeData struct {
	CoinTypes CoinsData
	Pairs     PairsData
}

type OrderItem struct {
	Amt       float64 `json:"amt"`
	Direction string  `json:"direction"`
	Pair      string  `json:"pair"`
	Price     float64 `json:"price"`
}

type OrderBook struct {
	Asks []OrderItem `json:"asks"`
	Bids []OrderItem `json:"bids"`
}

type TradeHistory struct {
	Code int `json:"code"`
	Data []struct {
		CreatedTime        int64   `json:"createdTime"`
		DealAmount         float64 `json:"dealAmount"`
		Pair               string  `json:"pair"`
		Price              float64 `json:"price"`
		TradeDealDirection string  `json:"tradeDealDirection"`
	} `json:"data"`
	Msg string `json:"msg"`
}

/********** Private API Structure**********/
// type Balance struct {
// 	Available float64 `json:"available"`
// 	CoinType  string  `json:"coinType"`
// 	Frozen    float64 `json:"frozen"`
// 	Total     float64 `json:"total"`
// }

// type AccountBalances []Balance

type AccountBalances struct {
	WALLET []struct {
		Available float64 `json:"available"`
		CoinType  string  `json:"coinType"`
		Frozen    float64 `json:"frozen"`
		Total     float64 `json:"total"`
	} `json:"WALLET"`
	OTC []interface{} `json:"OTC"`
}

type WithdrawResponse struct {
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

type OrderStatus struct {
	CreatedTime         int64       `json:"createdTime"`
	DealAmount          float64     `json:"dealAmount"`
	DealAvgPrice        float64     `json:"dealAvgPrice"`
	Direction           string      `json:"direction"`
	FrozenAmountByOrder float64     `json:"frozenAmountByOrder"`
	ID                  string      `json:"id"`
	OrderType           string      `json:"orderType"`
	Pair                string      `json:"pair"`
	Price               float64     `json:"price"`
	Status              int         `json:"status"`
	TotalAmount         float64     `json:"totalAmount"`
	UpdateTime          interface{} `json:"updateTime"`
}

type OrderDetail struct {
	CreatedTime         int         `json:"createdTime"`
	DealAmount          float64     `json:"dealAmount"`
	DealAvgPrice        float64     `json:"dealAvgPrice"`
	Direction           string      `json:"direction"`
	FrozenAmountByOrder float64     `json:"frozenAmountByOrder"`
	Id                  string      `json:"id"`
	OrderType           string      `json:"orderType"`
	Pair                string      `json:"pair"`
	Price               float64     `json:"price"`
	Status              int         `json:"status"`
	TotalAmount         float64     `json:"totalAmount"`
	UpdateTime          interface{} `json:"updateTime"`
}

type OrdersPage struct {
	Data        []OrderDetail `json:"data"`
	PageRequest interface{}   `json:"pageRequest"`
	Total       int           `json:"total"`
}
