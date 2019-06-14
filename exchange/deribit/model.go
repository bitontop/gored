package deribit

import "encoding/json"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct {
	Jsonrpc string          `json:"jsonrpc"`
	Data    json.RawMessage `json:"result"`
	UsIn    int64           `json:"usIn"`
	UsOut   int64           `json:"usOut"`
	UsDiff  int             `json:"usDiff"`
	Testnet bool            `json:"testnet"`
}

/********** Public API Structure**********/
type ContractsData []struct {
	TickSize            float64 `json:"tick_size"`
	Strike              float64 `json:"strike,omitempty"`
	SettlementPeriod    string  `json:"settlement_period"`
	QuoteCurrency       string  `json:"quote_currency"`
	OptionType          string  `json:"option_type,omitempty"`
	MinTradeAmount      float64 `json:"min_trade_amount"`
	Kind                string  `json:"kind"`
	IsActive            bool    `json:"is_active"`
	InstrumentName      string  `json:"instrument_name"`
	ExpirationTimestamp int64   `json:"expiration_timestamp"`
	CreationTimestamp   int64   `json:"creation_timestamp"`
	ContractSize        float64 `json:"contract_size"`
	BaseCurrency        string  `json:"base_currency"`
}

type OrderBook struct {
	Timestamp int64 `json:"timestamp"`
	Stats     struct {
		Volume float64 `json:"volume"`
		Low    float64 `json:"low"`
		High   float64 `json:"high"`
	} `json:"stats"`
	State           string      `json:"state"`
	SettlementPrice float64     `json:"settlement_price"`
	OpenInterest    int         `json:"open_interest"`
	MinPrice        float64     `json:"min_price"`
	MaxPrice        float64     `json:"max_price"`
	MarkPrice       float64     `json:"mark_price"`
	LastPrice       float64     `json:"last_price"`
	InstrumentName  string      `json:"instrument_name"`
	IndexPrice      float64     `json:"index_price"`
	Funding8H       float64     `json:"funding_8h"`
	CurrentFunding  float64     `json:"current_funding"`
	ChangeID        int64       `json:"change_id"`
	Bids            [][]float64 `json:"bids"`
	BestBidPrice    float64     `json:"best_bid_price"`
	BestBidAmount   float64     `json:"best_bid_amount"`
	BestAskPrice    float64     `json:"best_ask_price"`
	BestAskAmount   float64     `json:"best_ask_amount"`
	Asks            [][]float64 `json:"asks"`
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
