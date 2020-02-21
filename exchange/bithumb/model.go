package bithumb

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

type JsonResponse struct {
	Data    json.RawMessage `json:"data"`
	Success bool            `json:"success"`
	Msg     string          `json:"msg"`
	Code    string          `json:"code"`
	Params  []interface{}   `json:"params"`
}

type CoinsData struct {
	CoinConfig []struct {
		// MakerFeeRate   string `json:"makerFeeRate"`
		MinTxAmt      string `json:"minTxAmt,omitempty"`
		Name          string `json:"name"`
		DepositStatus string `json:"depositStatus"`
		FullName      string `json:"fullName"`
		// TakerFeeRate   string `json:"takerFeeRate"`
		MinWithdraw    string `json:"minWithdraw"`
		WithdrawFee    string `json:"withdrawFee"`
		WithdrawStatus string `json:"withdrawStatus"`
	} `json:"coinConfig"`
	ContractConfig []struct {
		Symbol       string `json:"symbol"`
		MakerFeeRate string `json:"makerFeeRate"`
		TakerFeeRate string `json:"takerFeeRate"`
	} `json:"contractConfig"`
	SpotConfig []struct {
		Symbol   string   `json:"symbol"`
		Accuracy []string `json:"accuracy"`
		OpenTime int64    `json:"openTime"`
	} `json:"spotConfig"`
}

type OrderBook struct {
	Symbol string     `json:"symbol"`
	B      [][]string `json:"b"`
	Ver    string     `json:"ver"`
	S      [][]string `json:"s"`
}

type AccountBalances []struct {
	CoinType    string `json:"coinType"`
	Count       string `json:"count"`
	Frozen      string `json:"frozen"`
	BtcQuantity string `json:"btcQuantity"`
	Type        string `json:"type"`
}

type PlaceOrder struct {
	OrderID string `json:"orderId"`
	Symbol  string `json:"symbol"`
}

type OrderStatus struct {
	OrderID    string `json:"orderId"`
	Symbol     string `json:"symbol"`
	Price      string `json:"price"`
	TradedNum  string `json:"tradedNum"`
	Quantity   string `json:"quantity"`
	AvgPrice   string `json:"avgPrice"`
	Status     string `json:"status"`
	Type       string `json:"type"`
	Side       string `json:"side"`
	CreateTime int64  `json:"createTime"`
	TradeTotal string `json:"tradeTotal"`
}

type TradeHistory struct {
	Data []struct {
		P   string `json:"p"`
		Ver string `json:"ver"`
		S   string `json:"s"`
		T   string `json:"t"`
		V   string `json:"v"`
	} `json:"data"`
	Code      string `json:"code"`
	Msg       string `json:"msg"`
	Timestamp int64  `json:"timestamp"`
}
