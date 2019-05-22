package tradesatoshi

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

type JsonResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

type CoinsData []struct {
	Currency        string  `json:"currency"`
	CurrencyLong    string  `json:"currencyLong"`
	MinConfirmation int     `json:"minConfirmation"`
	TxFee           float64 `json:"txFee"`
	Status          string  `json:"status"`
	StatusMessage   string  `json:"statusMessage"`
	MinBaseTrade    float64 `json:"minBaseTrade"`
	IsTipEnabled    bool    `json:"isTipEnabled"`
	MinTip          float64 `json:"minTip"`
	MaxTip          float64 `json:"maxTip"`
}

type PairsData []struct {
	Market         string  `json:"market"`
	High           float64 `json:"high"`
	Low            float64 `json:"low"`
	Volume         float64 `json:"volume"`
	BaseVolume     float64 `json:"baseVolume"`
	Last           float64 `json:"last"`
	Bid            float64 `json:"bid"`
	Ask            float64 `json:"ask"`
	OpenBuyOrders  int     `json:"openBuyOrders"`
	OpenSellOrders int     `json:"openSellOrders"`
	MarketStatus   string  `json:"marketStatus"`
	Change         float64 `json:"change"`
}

type OrderBook struct {
	Buy []struct {
		Quantity float64 `json:"quantity"`
		Rate     float64 `json:"rate"`
	} `json:"buy"`
	Sell []struct {
		Quantity float64 `json:"quantity"`
		Rate     float64 `json:"rate"`
	} `json:"sell"`
}

type AccountBalances []struct {
	Currency        string  `json:"currency"`
	CurrencyLong    string  `json:"currencyLong"`
	Available       float64 `json:"available"`
	Total           float64 `json:"total"`
	HeldForTrades   float64 `json:"heldForTrades"`
	Unconfirmed     float64 `json:"unconfirmed"`
	PendingWithdraw float64 `json:"pendingWithdraw"`
	Address         string  `json:"address"`
}

type Withdraw struct {
	WithdrawalID int `json:"WithdrawalId"`
}

type PlaceOrder struct {
	OrderID int   `json:"OrderId"`
	Filled  []int `json:"Filled"`
}

type OrderStatus struct {
	ID        int     `json:"id"`
	Market    string  `json:"market"`
	Type      string  `json:"type"`
	Amount    float64 `json:"amount"`
	Rate      float64 `json:"rate"`
	Remaining float64 `json:"remaining"`
	Total     float64 `json:"total"`
	Status    string  `json:"status"`
	Timestamp string  `json:"timestamp"`
	IsAPI     bool    `json:"isApi"`
}

type CancelOrder struct {
	CanceledOrders []int `json:"CanceledOrders"`
}
