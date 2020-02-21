package bitfinex

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type CoinsData [][][]string

type WithdrawFee struct {
	Withdraw map[string]interface{}
}

type PairsData []struct {
	Pair             string `json:"pair"`
	PricePrecision   int    `json:"price_precision"`
	InitialMargin    string `json:"initial_margin"`
	MinimumMargin    string `json:"minimum_margin"`
	MaximumOrderSize string `json:"maximum_order_size"`
	MinimumOrderSize string `json:"minimum_order_size"`
	Expiration       string `json:"expiration"`
	Margin           bool   `json:"margin"`
}

type OrderBook struct {
	Bids []struct {
		Price     string `json:"price"`
		Amount    string `json:"amount"`
		Timestamp string `json:"timestamp"`
	} `json:"bids"`
	Asks []struct {
		Price     string `json:"price"`
		Amount    string `json:"amount"`
		Timestamp string `json:"timestamp"`
	} `json:"asks"`
}

type AccountBalances []struct {
	Type      string `json:"type"`
	Currency  string `json:"currency"`
	Amount    string `json:"amount"`
	Available string `json:"available"`
}

type PlaceOrder struct {
	ID                int64       `json:"id"`
	Cid               int         `json:"cid"`
	CidDate           string      `json:"cid_date"`
	Gid               interface{} `json:"gid"`
	Symbol            string      `json:"symbol"`
	Exchange          string      `json:"exchange"`
	Price             string      `json:"price"`
	AvgExecutionPrice string      `json:"avg_execution_price"`
	Side              string      `json:"side"`
	Type              string      `json:"type"`
	Timestamp         string      `json:"timestamp"`
	IsLive            bool        `json:"is_live"`
	IsCancelled       bool        `json:"is_cancelled"`
	IsHidden          bool        `json:"is_hidden"`
	OcoOrder          interface{} `json:"oco_order"`
	WasForced         bool        `json:"was_forced"`
	OriginalAmount    string      `json:"original_amount"`
	RemainingAmount   string      `json:"remaining_amount"`
	ExecutedAmount    string      `json:"executed_amount"`
	Src               string      `json:"src"`
	OrderID           int64       `json:"order_id"`
	Message           string      `json:"message"`
}

type Withdraw []struct {
	Status       string `json:"status"`
	Message      string `json:"message"`
	WithdrawalID int    `json:"withdrawal_id"`
}

type TradeHistory [][]float64
