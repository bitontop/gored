package dcoin

import "encoding/json"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

/********** Public API Structure**********/
type PairsData []struct {
	Symbol          string `json:"symbol"`
	CountCoin       string `json:"count_coin"`
	AmountPrecision int    `json:"amount_precision"`
	BaseCoin        string `json:"base_coin"`
	PricePrecision  int    `json:"price_precision"`
}

type OrderBook struct {
	Asks [][]float64 `json:"asks"`
	Bids [][]float64 `json:"bids"`
}

/********** Private API Structure**********/
type AccountBalances struct {
	CoinList []struct {
		Coin   string  `json:"coin"`
		Normal float64 `json:"normal"`
		Locked float64 `json:"locked"`
	} `json:"coin_list"`
}

type PlaceOrder struct {
	OrderID int `json:"order_id"`
}

type OrderStatus struct {
	OrderInfo struct {
		ID         int     `json:"id"`
		Side       string  `json:"side"`
		Symbol     string  `json:"symbol"`
		Type       int     `json:"type"`
		Price      float64 `json:"price"`
		Volume     float64 `json:"volume"`
		Status     int     `json:"status"`
		DealVolume float64 `json:"deal_volume"`
		TotalPrice float64 `json:"total_price"`
		Fee        int     `json:"fee"`
		AgePrice   float64 `json:"average_price"`
		Ts         int64   `json:"ts"`
	} `json:"order_info"`
	TradeList []struct {
		ID        int     `json:"id"`
		Price     float64 `json:"price"`
		Volume    float64 `json:"volume"`
		Direction string  `json:"direction"`
		Ts        int64   `json:"ts"`
	} `json:"trade_list"`
}
