package coinbase

import "time"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type CoinsData []struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	MinSize string `json:"min_size"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Details struct {
		Type                  string   `json:"type"`
		Symbol                string   `json:"symbol"`
		NetworkConfirmations  int      `json:"network_confirmations"`
		SortOrder             int      `json:"sort_order"`
		CryptoAddressLink     string   `json:"crypto_address_link"`
		CryptoTransactionLink string   `json:"crypto_transaction_link"`
		PushPaymentMethods    []string `json:"push_payment_methods"`
		ProcessingTimeSeconds int      `json:"processing_time_seconds"`
		MinWithdrawalAmount   float64  `json:"min_withdrawal_amount"`
	} `json:"details"`
	MaxPrecision  string   `json:"max_precision"`
	ConvertibleTo []string `json:"convertible_to,omitempty"`
}

type PairsData []struct {
	ID             string `json:"id"`
	BaseCurrency   string `json:"base_currency"`
	QuoteCurrency  string `json:"quote_currency"`
	BaseMinSize    string `json:"base_min_size"`
	BaseMaxSize    string `json:"base_max_size"`
	QuoteIncrement string `json:"quote_increment"`
	BaseIncrement  string `json:"base_increment"`
	DisplayName    string `json:"display_name"`
	MinMarketFunds string `json:"min_market_funds"`
	MaxMarketFunds string `json:"max_market_funds"`
	MarginEnabled  bool   `json:"margin_enabled"`
	PostOnly       bool   `json:"post_only"`
	LimitOnly      bool   `json:"limit_only"`
	CancelOnly     bool   `json:"cancel_only"`
	Status         string `json:"status"`
	StatusMessage  string `json:"status_message"`
}

type TradeHistory []struct {
	Time    time.Time `json:"time"`
	TradeID int       `json:"trade_id"`
	Price   string    `json:"price"`
	Size    string    `json:"size"`
	Side    string    `json:"side"`
}

// ====================================

type AccountBalances []struct {
	ID        string `json:"id"`
	Currency  string `json:"currency"`
	Balance   string `json:"balance"`
	Available string `json:"available"`
	Hold      string `json:"hold"`
	ProfileID string `json:"profile_id"`
}

type WithdrawResponse struct {
	ID       string    `json:"id"`
	Amount   string    `json:"amount"`
	Currency string    `json:"currency"`
	PayoutAt time.Time `json:"payout_at"`
}

type PlaceOrder struct {
	ID             string    `json:"id"`
	Price          string    `json:"price"`
	Size           string    `json:"size"`
	ProductID      string    `json:"product_id"`
	Side           string    `json:"side"`
	Stp            string    `json:"stp"`
	Funds          string    `json:"funds"`
	SpecifiedFunds string    `json:"specified_funds"`
	Type           string    `json:"type"`
	TimeInForce    string    `json:"time_in_force"`
	PostOnly       bool      `json:"post_only"`
	CreatedAt      time.Time `json:"created_at"`
	DoneAt         time.Time `json:"done_at"`
	DoneReason     string    `json:"done_reason"`
	FillFees       string    `json:"fill_fees"`
	FilledSize     string    `json:"filled_size"`
	ExecutedValue  string    `json:"executed_value"`
	Status         string    `json:"status"`
	Settled        bool      `json:"settled"`
	Message        string    `json:"message"`
}

type OrderBook struct {
	Sequence int64      `json:"sequence"`
	Bids     [][]string `json:"bids"`
	Asks     [][]string `json:"asks"`
}
