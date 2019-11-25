package lbank

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

// v2
type CoinsData struct {
	Result string `json:"result"`
	Data   []struct {
		AmountScale string `json:"amountScale"`
		AssetCode   string `json:"assetCode"`
		CanWithDraw bool   `json:"canWithDraw"`
		Fee         string `json:"fee"`
		Type        string `json:"type"`
		Min         string `json:"min,omitempty"`
	} `json:"data"`
	ErrorCode int   `json:"error_code"`
	Ts        int64 `json:"ts"`
}

/* type CoinsData []struct {
	AssetCode   string `json:"assetCode"`
	Min         string `json:"min"`
	CanWithDraw bool   `json:"canWithDraw"`
	Fee         string `json:"fee"`
} */

type PairsData []struct {
	MinTranQua       string `json:"minTranQua"`
	PriceAccuracy    string `json:"priceAccuracy"`
	QuantityAccuracy string `json:"quantityAccuracy"`
	Symbol           string `json:"symbol"`
}

type OrderBook struct {
	Bids      [][]float64 `json:"bids"`
	Asks      [][]float64 `json:"asks"`
	Timestamp int64       `json:"timestamp"`
	// Message      type         `json:"msg"`
}

type AccountBalances struct {
	Result string `json:"result"`
	Info   struct {
		Asset  map[string]string `json:"asset"`
		Freeze map[string]string `json:"freeze"`
		Free   map[string]string `json:"free"`
	} `json:"info"`
	ErrorCode int `json:"error_code"`
}

type PlaceOrder struct {
	Result    string `json:"result"`
	OrderID   string `json:"order_id"`
	ErrorCode int    `json:"error_code"`
}

type OrderStatus struct {
	Result    string `json:"result"`
	ErrorCode int    `json:"error_code"`
	Orders    []struct {
		Symbol     string      `json:"symbol"`
		Amount     float64     `json:"amount"`
		CreateTime int64       `json:"create_time"`
		Price      float64     `json:"price"`
		CustomID   interface{} `json:"custom_id"`
		AvgPrice   float64     `json:"avg_price"`
		Type       string      `json:"type"`
		OrderID    string      `json:"order_id"`
		DealAmount float64     `json:"deal_amount"`
		Status     int         `json:"status"`
	} `json:"orders"`
}

type CancelOrders struct {
	Result    string `json:"result"`
	OrderID   string `json:"order_id"`
	Success   string `json:"success"`
	Error     string `json:"error"`
	ErrorCode int    `json:"error_code"`
}

type Withdraw struct {
	Result     string  `json:"result"`
	WithdrawID int     `json:"withdrawId"`
	Fee        float64 `json:"fee"`
}
