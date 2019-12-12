package tagz

import "encoding/json"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct { // delete
	Success bool            `json:"success"`
	Code    string          `json:"code"`
	Msg     string          `json:"msg"`
	Retry   bool            `json:"retry"`
	Data    json.RawMessage `json:"data"`
}

type PairsData []string

type OrderBook struct {
	Instrument string `json:"instrument"`
	Bids       []struct {
		Amount float64 `json:"amount"`
		Price  float64 `json:"price"`
	} `json:"bids"`
	Asks []struct {
		Amount float64 `json:"amount"`
		Price  float64 `json:"price"`
	} `json:"asks"`
	Version        int     `json:"version"`
	AskTotalAmount float64 `json:"askTotalAmount"`
	BidTotalAmount float64 `json:"bidTotalAmount"`
	Snapshot       bool    `json:"snapshot"`
}

//------------ OLD

type AccountBalance []struct {
	Balance   string `json:"balance"`
	Available string `json:"available"`
	Holds     string `json:"holds"`
	Currency  string `json:"currency"`
	ID        string `json:"id"`
	Type      string `json:"type"`
}

// v1 api
/* type InnerTransIDs struct {
	clientOid    string
	payAccountId string
	recAccountId string
	freeAmount   float64
} */

type Withdraw struct {
	WithdrawalID string `json:"withdrawalId"`
}

type OrderDetail struct {
	OrderID string `json:"orderId"`
}

type OrderStatus struct {
	ID            string `json:"id"`
	Symbol        string `json:"symbol"`
	OpType        string `json:"opType"`
	Type          string `json:"type"`
	Side          string `json:"side"`
	Price         string `json:"price"`
	Size          string `json:"size"`
	Funds         string `json:"funds"`
	DealFunds     string `json:"dealFunds"`
	DealSize      string `json:"dealSize"`
	Fee           string `json:"fee"`
	FeeCurrency   string `json:"feeCurrency"`
	Stp           string `json:"stp"`
	Stop          string `json:"stop"`
	StopTriggered bool   `json:"stopTriggered"`
	StopPrice     string `json:"stopPrice"`
	TimeInForce   string `json:"timeInForce"`
	PostOnly      bool   `json:"postOnly"`
	Hidden        bool   `json:"hidden"`
	Iceberg       bool   `json:"iceberg"`
	VisibleSize   string `json:"visibleSize"`
	CancelAfter   int    `json:"cancelAfter"`
	Channel       string `json:"channel"`
	ClientOid     string `json:"clientOid"`
	Remark        string `json:"remark"`
	Tags          string `json:"tags"`
	IsActive      bool   `json:"isActive"`
	CancelExist   bool   `json:"cancelExist"`
	CreatedAt     int64  `json:"createdAt"`
}

type CancelOrder struct {
	CancelledOrderIds []string `json:"cancelledOrderIds"`
}

type InnerTrans struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	OrderID string `json:"orderId"`
}

type AccountID []struct {
	ID        string `json:"id"`
	Currency  string `json:"currency"`
	Type      string `json:"type"`
	Balance   string `json:"balance"`
	Available string `json:"available"`
	Holds     string `json:"holds"`
}
