package hibitex

import "encoding/json"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct {
	Code    string          `json:"code"`
	Msg     string          `json:"msg"`
	Message interface{}     `json:"message"`
	Data    json.RawMessage `json:"data"`
}

/********** Public API Structure**********/
type Ticker struct {
	Symbol          string `json:"symbol"`
	CountCoin       string `json:"count_coin"`
	AmountPrecision int    `json:"amount_precision"`
	BaseCoin        string `json:"base_coin"`
	PricePrecision  int    `json:"price_precision"`
}

type CoinsData []Ticker

type PairsData []Ticker

type OrderBook struct {
	Tick struct {
		Asks [][]float64 `json:"asks"`
		Bids [][]float64 `json:"bids"`
		Time interface{} `json:"time"`
	} `json:"tick"`
}

// type Order struct {
// 	Amount float64 `json:"amount"`
// 	Price  float64 `json:"price"`
// 	Id     int     `json:"id"`
// 	Type   string  `json:"type"`
// 	Ts     int     `json:"ts"`
// 	Ds     string  `json:"ds"`
// }

// type OrderBook []Order

/********** Private API Structure**********/
type AccountBalances struct {
	TotalAsset float64 `json:"total_asset"`
	CoinList   []struct {
		Coin        string  `json:"coin"`
		Normal      float64 `json:"normal"`
		Locked      float64 `json:"locked"`
		BtcValuatin float64 `json:"btcValuatin"`
	} `json:"coin_list"`
}

type PlaceOrder struct {
	OrderID int `json:"order_id"`
}

type OrderStatus struct {
	Count      int `json:"count"`
	ResultList []struct {
		Side         string `json:"side"`
		TotalPrice   string `json:"total_price"`
		CreatedAt    int64  `json:"created_at"`
		AvgPrice     string `json:"avg_price"`
		CountCoin    string `json:"countCoin"`
		Source       int    `json:"source"`
		Type         int    `json:"type"`
		SideMsg      string `json:"side_msg"`
		Volume       string `json:"volume"`
		Price        string `json:"price"`
		SourceMsg    string `json:"source_msg"`
		StatusMsg    string `json:"status_msg"`
		DealVolume   string `json:"deal_volume"`
		ID           int    `json:"id"`
		RemainVolume string `json:"remain_volume"`
		BaseCoin     string `json:"baseCoin"`
		Status       int    `json:"status"`
	} `json:"resultList"`
}

// type OrderStatus struct {
// 	OrderInfo struct {
// 		ID         int     `json:"id"`
// 		Side       string  `json:"side"`
// 		SideMsg    string  `json:"side_msg"`
// 		CreatedAt  string  `json:"created_at"`
// 		Price      float64 `json:"price"`
// 		Volume     float64 `json:"volume"`
// 		DealVolume float64 `json:"deal_volume"`
// 		TotalPrice float64 `json:"total_price"`
// 		Fee        float64 `json:"fee"`
// 		AvgPrice   float64 `json:"avg_price"`
// 	} `json:"order_info"`
// 	TradeList []struct {
// 		ID        int     `json:"id"`
// 		CreatedAt string  `json:"created_at"`
// 		Price     float64 `json:"price"`
// 		Volume    float64 `json:"volume"`
// 		DealPrice float64 `json:"deal_price"`
// 		Fee       float64 `json:"fee"`
// 	} `json:"trade_list"`
// }

// ================ Old

type WithdrawResponse struct {
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	ID      string `json:"id"`
}
