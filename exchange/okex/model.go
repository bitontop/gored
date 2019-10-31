package okex

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"time"
)

type ErrorMsg struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

type CoinsData []struct {
	CanDeposit    string `json:"can_deposit"`
	CanWithdraw   string `json:"can_withdraw"`
	Currency      string `json:"currency"`
	MinWithdrawal string `json:"min_withdrawal"`
	Name          string `json:"name"`
}

type PairsData []struct {
	BaseCurrency  string `json:"base_currency"`
	InstrumentID  string `json:"instrument_id"`
	MinSize       string `json:"min_size"`
	QuoteCurrency string `json:"quote_currency"`
	SizeIncrement string `json:"size_increment"`
	TickSize      string `json:"tick_size"`
}

type OrderBook struct {
	Asks      [][]string `json:"asks"`
	Bids      [][]string `json:"bids"`
	Timestamp time.Time  `json:"timestamp"`
}

type AccountBalances []struct {
	Frozen    string `json:"frozen"`
	Hold      string `json:"hold"`
	ID        string `json:"id"`
	Currency  string `json:"currency"`
	Balance   string `json:"balance"`
	Available string `json:"available"`
	Holds     string `json:"holds"`
}

type WithdrawResponse struct {
	Result       bool   `json:"result"`
	Amount       string `json:"amount"`
	WithdrawalID string `json:"withdrawal_id"`
	Currency     string `json:"currency"`
	Code         int    `json:"code"`
	Message      string `json:"message"`
}

type WithdrawFee []struct {
	Currency string `json:"currency"`
	MaxFee   string `json:"max_fee"`
	MinFee   string `json:"min_fee"`
}

type Transfer struct {
	TransferID int     `json:"transfer_id"`
	Currency   string  `json:"currency"`
	From       int     `json:"from"`
	Amount     float64 `json:"amount"`
	To         int     `json:"to"`
	Result     bool    `json:"result"`
	Code       int     `json:"code"`
	Message    string  `json:"message"`
}

type PlaceOrder struct {
	OrderID   string `json:"order_id"`
	ClientOid string `json:"client_oid"`
	Result    bool   `json:"result"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
}

type OrderStatus struct {
	OrderID        string    `json:"order_id"`
	Notional       string    `json:"notional"`
	Price          string    `json:"price"`
	Size           string    `json:"size"`
	InstrumentID   string    `json:"instrument_id"`
	Side           string    `json:"side"`
	Type           string    `json:"type"`
	Timestamp      time.Time `json:"timestamp"`
	FilledSize     string    `json:"filled_size"`
	FilledNotional string    `json:"filled_notional"`
	Status         string    `json:"status"`
	Code           int       `json:"code"`
	Message        string    `json:"message"`
}

type WSOrderBook struct {
	Table  string `json:"table"`
	Action string `json:"action"`
	Data   []struct {
		InstrumentID string     `json:"instrument_id"`
		Asks         [][]string `json:"asks"`
		Bids         [][]string `json:"bids"`
		Timestamp    time.Time  `json:"timestamp"`
		Checksum     int        `json:"checksum"`
	} `json:"data"`
}
