package cointiger

import "encoding/json"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

/********** Public API Structure**********/
type PairsData map[string]PairsDetails

type PairsDetails []struct {
	BaseCurrency    string  `json:"baseCurrency"`
	QuoteCurrency   string  `json:"quoteCurrency"`
	PricePrecision  int     `json:"pricePrecision"`
	AmountPrecision int     `json:"amountPrecision"`
	AmountMin       float64 `json:"amountMin"`
	WithdrawFeeMin  float64 `json:"withdrawFeeMin"`
	WithdrawFeeMax  float64 `json:"withdrawFeeMax"`
	WithdrawOneMin  float64 `json:"withdrawOneMin"`
	WithdrawOneMax  float64 `json:"withdrawOneMax"`
	DepthSelect     struct {
		Step0 string `json:"step0"`
		Step1 string `json:"step1"`
		Step2 string `json:"step2"`
	} `json:"depthSelect"`
}

type OrderBook struct {
	Symbol    string `json:"symbol"`
	DepthData struct {
		Tick struct {
			Buys [][]interface{} `json:"buys"`
			Asks [][]interface{} `json:"asks"`
		} `json:"tick"`
		Ts int64 `json:"ts"`
	} `json:"depth_data"`
}

/********** Private API Structure**********/
type AccountBalances []struct {
	Normal string `json:"normal"`
	Lock   string `json:"lock"`
	Coin   string `json:"coin"`
}

type PlaceOrder struct {
	OrderID int `json:"order_id"`
}

type OrderStatus struct {
	Symbol     string `json:"symbol"`
	Fee        string `json:"fee"`
	AvgPrice   string `json:"avg_price"`
	Source     int    `json:"source"`
	Type       string `json:"type"`
	Mtime      int64  `json:"mtime"`
	Volume     string `json:"volume"`
	UserID     int    `json:"user_id"`
	Price      string `json:"price"`
	Ctime      int64  `json:"ctime"`
	DealVolume string `json:"deal_volume"`
	ID         int    `json:"id"`
	DealMoney  string `json:"deal_money"`
	Status     int    `json:"status"`
}

type CancelOrder struct {
	Success []int `json:"success"`
	Failed  []struct {
		ErrMsg  string `json:"err-msg"`
		OrderID string `json:"order-id"`
		ErrCode string `json:"err-code"`
	} `json:"failed"`
}
