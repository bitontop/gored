package bilaxy

import "encoding/json"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

/********** Public API Structure**********/
type CoinsData struct {
	Currency                string `json:"currency"`
	FullName                string `json:"fullname"`
	MinWithdraw             string `json:"min_withdraw"`
	MaxWithdraw             string `json:"max_withdraw"`
	FixedWithdrawFee        string `json:"fixed_withdraw_fee"`
	PercentWithdrawFee      string `json:"percent_withdraw_fee"`
	WithdrawFeeCurrencyId   int    `json:"withdraw_fee_currency_id"`
	WithdrawFeeCurrencyName string `json:"withdraw_fee_currency_name"`
	WithdrawEnabled         bool   `json:"withdraw_enabled"`
	DepositEnabled          bool   `json:"deposit_enabled"`
}

type PairsData struct {
	PairId          int    `json:"pair_id"`
	Base            string `json:"base"`
	Quote           string `json:"quote"`
	PricePrecision  int    `json:"price_precision"`
	AmountPrecision int    `json:"amount_precision"`
	MinAmount       string `json:"min_amount"`
	MaxAmount       string `json:"max_amount"`
	MinTotal        string `json:"min_total"`
	MaxTotal        string `json:"max_total"`
	TradeEnabled    bool   `json:"trade_enabled"`
	Closed          bool   `json:"closed"`
}

type OrderBook struct {
	Bids [][]string `json:"bids"`
	Asks [][]string `json:"asks"`
}

/********** Private API Structure**********/
type AccountBalances []struct {
	Asset     string  `json:"asset"`
	Total     float64 `json:"total"`
	Available float64 `json:"available"`
	Locked    float64 `json:"locked"`
}

type WithdrawResponse struct {
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

type PlaceOrder struct {
	Symbol       string `json:"symbol"`
	OrderID      string `json:"orderId"`
	Side         string `json:"side"`
	Type         string `json:"type"`
	Price        string `json:"price"`
	AveragePrice string `json:"executedQty"`
	OrigQty      string `json:"origQty"`
	ExecutedQty  string `json:"executedQty"`
	Status       string `json:"status"`
	TimeInForce  string `json:"timeInForce"`
}
