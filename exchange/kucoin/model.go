package kucoin

import "encoding/json"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct {
	Success bool            `json:"success"`
	Code    string          `json:"code"`
	Msg     string          `json:"msg"`
	Retry   bool            `json:"retry"`
	Data    json.RawMessage `json:"data"`
}

type CoinsData []struct {
	WithdrawalMinFee  string `json:"withdrawalMinFee"`
	Precision         int    `json:"precision"`
	Name              string `json:"name"`
	FullName          string `json:"fullName"`
	Currency          string `json:"currency"`
	WithdrawalMinSize string `json:"withdrawalMinSize"`
	IsWithdrawEnabled bool   `json:"isWithdrawEnabled"`
	IsDepositEnabled  bool   `json:"isDepositEnabled"`
}

type PairsData []struct {
	Symbol         string `json:"symbol"`
	QuoteMaxSize   string `json:"quoteMaxSize"`
	EnableTrading  bool   `json:"enableTrading"`
	PriceIncrement string `json:"priceIncrement"`
	BaseMaxSize    string `json:"baseMaxSize"`
	BaseCurrency   string `json:"baseCurrency"`
	QuoteCurrency  string `json:"quoteCurrency"`
	Market         string `json:"market"`
	QuoteIncrement string `json:"quoteIncrement"`
	BaseMinSize    string `json:"baseMinSize"`
	QuoteMinSize   string `json:"quoteMinSize"`
	Name           string `json:"name"`
	BaseIncrement  string `json:"baseIncrement"`
}

type OrderBook struct {
	Sequence string     `json:"sequence"`
	Asks     [][]string `json:"asks"`
	Bids     [][]string `json:"bids"`
}

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

type InnerTrans struct {
	OrderID string `json:"orderId"`
}

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
