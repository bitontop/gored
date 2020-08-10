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

type TickerPrice struct {
	Time   int64 `json:"time"`
	Ticker []struct {
		Symbol       string `json:"symbol"`
		SymbolName   string `json:"symbolName"`
		Buy          string `json:"buy"`
		Sell         string `json:"sell"`
		ChangeRate   string `json:"changeRate"`
		ChangePrice  string `json:"changePrice"`
		High         string `json:"high"`
		Low          string `json:"low"`
		Vol          string `json:"vol"`
		VolValue     string `json:"volValue"`
		Last         string `json:"last"`
		AveragePrice string `json:"averagePrice"`
	} `json:"ticker"`
}

type KLine struct {
	Code string     `json:"code"`
	Data [][]string `json:"data"`
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

type TradeHistory []struct {
	Sequence string `json:"sequence"`
	Side     string `json:"side"`
	Size     string `json:"size"`
	Price    string `json:"price"`
	Time     int64  `json:"time"`
}

type SubAllAccountBalances []struct {
	SubUserID    string `json:"subUserId"`
	SubName      string `json:"subName"`
	MainAccounts []struct {
		Currency          string `json:"currency"`
		Balance           string `json:"balance"`
		Available         string `json:"available"`
		Holds             string `json:"holds"`
		BaseCurrency      string `json:"baseCurrency"`
		BaseCurrencyPrice string `json:"baseCurrencyPrice"`
		BaseAmount        string `json:"baseAmount"`
	} `json:"mainAccounts"`
	TradeAccounts []struct {
		Currency          string `json:"currency"`
		Balance           string `json:"balance"`
		Available         string `json:"available"`
		Holds             string `json:"holds"`
		BaseCurrency      string `json:"baseCurrency"`
		BaseCurrencyPrice string `json:"baseCurrencyPrice"`
		BaseAmount        string `json:"baseAmount"`
	} `json:"tradeAccounts"`
	MarginAccounts []struct {
		Currency          string `json:"currency"`
		Balance           string `json:"balance"`
		Available         string `json:"available"`
		Holds             string `json:"holds"`
		BaseCurrency      string `json:"baseCurrency"`
		BaseCurrencyPrice string `json:"baseCurrencyPrice"`
		BaseAmount        string `json:"baseAmount"`
	} `json:"marginAccounts"`
}

type SubAccountBalances struct {
	SubUserID    string `json:"subUserId"`
	SubName      string `json:"subName"`
	MainAccounts []struct {
		Currency          string `json:"currency"`
		Balance           string `json:"balance"`
		Available         string `json:"available"`
		Holds             string `json:"holds"`
		BaseCurrency      string `json:"baseCurrency"`
		BaseCurrencyPrice string `json:"baseCurrencyPrice"`
		BaseAmount        string `json:"baseAmount"`
	} `json:"mainAccounts"`
	TradeAccounts []struct {
		Currency          string `json:"currency"`
		Balance           string `json:"balance"`
		Available         string `json:"available"`
		Holds             string `json:"holds"`
		BaseCurrency      string `json:"baseCurrency"`
		BaseCurrencyPrice string `json:"baseCurrencyPrice"`
		BaseAmount        string `json:"baseAmount"`
	} `json:"tradeAccounts"`
	MarginAccounts []struct {
		Currency          string `json:"currency"`
		Balance           string `json:"balance"`
		Available         string `json:"available"`
		Holds             string `json:"holds"`
		BaseCurrency      string `json:"baseCurrency"`
		BaseCurrencyPrice string `json:"baseCurrencyPrice"`
		BaseAmount        string `json:"baseAmount"`
	} `json:"marginAccounts"`
}

type SubAccountList []struct {
	UserID  string `json:"userId"`
	SubName string `json:"subName"`
	Remarks string `json:"remarks"`
}

type OpenOrders struct {
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
	TotalNum    int `json:"totalNum"`
	TotalPage   int `json:"totalPage"`
	Items       []struct {
		ID            string      `json:"id"`
		Symbol        string      `json:"symbol"`
		OpType        string      `json:"opType"`
		Type          string      `json:"type"`
		Side          string      `json:"side"`
		Price         string      `json:"price"`
		Size          string      `json:"size"`
		Funds         string      `json:"funds"`
		DealFunds     string      `json:"dealFunds"`
		DealSize      string      `json:"dealSize"`
		Fee           string      `json:"fee"`
		FeeCurrency   string      `json:"feeCurrency"`
		Stp           string      `json:"stp"`
		Stop          string      `json:"stop"`
		StopTriggered bool        `json:"stopTriggered"`
		StopPrice     string      `json:"stopPrice"`
		TimeInForce   string      `json:"timeInForce"`
		PostOnly      bool        `json:"postOnly"`
		Hidden        bool        `json:"hidden"`
		Iceberg       bool        `json:"iceberg"`
		VisibleSize   string      `json:"visibleSize"`
		CancelAfter   int         `json:"cancelAfter"`
		Channel       string      `json:"channel"`
		ClientOid     string      `json:"clientOid"`
		Remark        interface{} `json:"remark"`
		Tags          interface{} `json:"tags"`
		IsActive      bool        `json:"isActive"`
		CancelExist   bool        `json:"cancelExist"`
		CreatedAt     int64       `json:"createdAt"`
		TradeType     string      `json:"tradeType"`
	} `json:"items"`
}
