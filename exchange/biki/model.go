package biki

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

type JsonResponse struct {
	Code    string          `json:"code"`
	Message string          `json:"msg"`
	Result  json.RawMessage `json:"data"`
}

type PairsData []struct {
	Symbol          string `json:"symbol"`
	CountCoin       string `json:"count_coin"`
	AmountPrecision int    `json:"amount_precision"`
	BaseCoin        string `json:"base_coin"`
	PricePrecision  int    `json:"price_precision"`
}

type OrderBook struct {
	Tick struct {
		Asks [][]float64 `json:"asks"`
		Bids [][]float64 `json:"bids"`
		Time interface{} `json:"time"`
	} `json:"tick"`
}

type AccountBalances struct {
	TotalAsset string `json:"total_asset"`
	CoinList   []struct {
		Normal      string `json:"normal"`
		BtcValuatin string `json:"btcValuatin"`
		Locked      string `json:"locked"`
		Coin        string `json:"coin"`
	} `json:"coin_list"`
}

type PlaceOrder struct {
	OrderID int `json:"order_id"`
}

type OrderStatus struct {
	TradeList []interface{} `json:"trade_list"`
	OrderInfo struct {
		Side         string        `json:"side"`
		TotalPrice   string        `json:"total_price"`
		Fee          float64       `json:"fee"`
		CreatedAt    int64         `json:"created_at"`
		DealPrice    float64       `json:"deal_price"`
		AvgPrice     string        `json:"avg_price"`
		CountCoin    string        `json:"countCoin"`
		Source       int           `json:"source"`
		Type         int           `json:"type"`
		SideMsg      string        `json:"side_msg"`
		Volume       string        `json:"volume"`
		Price        string        `json:"price"`
		SourceMsg    string        `json:"source_msg"`
		StatusMsg    string        `json:"status_msg"`
		DealVolume   string        `json:"deal_volume"`
		FeeCoin      string        `json:"fee_coin"`
		ID           int           `json:"id"`
		RemainVolume string        `json:"remain_volume"`
		BaseCoin     string        `json:"baseCoin"`
		TradeList    []interface{} `json:"tradeList"`
		Status       int           `json:"status"`
	} `json:"order_info"`
}
