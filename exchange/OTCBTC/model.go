package otcbtc

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"time"
)

type ErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type PairsData []struct {
	ID                 string  `json:"id"`
	TickerID           string  `json:"ticker_id"`
	Name               string  `json:"name"`
	MinimalTotalVolume float64 `json:"minimal_total_volume"`
	TradingRule        struct {
		MinAmount      float64 `json:"min_amount"`
		MinPrice       float64 `json:"min_price"`
		MinOrderVolume float64 `json:"min_order_volume"`
	} `json:"trading_rule"`
}

type OrderBook struct {
	Timestamp int        `json:"timestamp"`
	Bids      [][]string `json:"bids"`
	Asks      [][]string `json:"asks"`
}

type AccountBalance struct {
	UserName      string `json:"user_name"`
	Email         string `json:"email"`
	Icon          string `json:"icon"`
	OtbFeeEnabled bool   `json:"otb_fee_enabled"`
	Accounts      []struct {
		Currency string `json:"currency"`
		Balance  string `json:"balance"`
		Locked   string `json:"locked"`
		Saving   string `json:"saving"`
	} `json:"accounts"`
}

type PlaceOrder struct {
	ID              int       `json:"id"`
	Side            string    `json:"side"`
	OrdType         string    `json:"ord_type"`
	Price           string    `json:"price"`
	AvgPrice        string    `json:"avg_price"`
	State           string    `json:"state"`
	Market          string    `json:"market"`
	CreatedAt       time.Time `json:"created_at"`
	Volume          string    `json:"volume"`
	RemainingVolume string    `json:"remaining_volume"`
	ExecutedVolume  string    `json:"executed_volume"`
	TradesCount     int       `json:"trades_count"`
}
