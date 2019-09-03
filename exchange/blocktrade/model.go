package blocktrade

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
type CoinsData []struct {
	ID                      int    `json:"id"`
	FullName                string `json:"full_name"`
	IsoCode                 string `json:"iso_code"`
	IconPath                string `json:"icon_path"`
	IconPathPng             string `json:"icon_path_png"`
	Color                   string `json:"color"`
	Sign                    string `json:"sign"`
	CurrencyType            string `json:"currency_type"`
	MinimalWithdrawalAmount string `json:"minimal_withdrawal_amount"`
	MinimalOrderValue       string `json:"minimal_order_value"`
	MaximumOrderValue       string `json:"maximum_order_value"`
	LotSize                 string `json:"lot_size"`
	DecimalPrecision        int    `json:"decimal_precision"`
}

type PairsData []struct {
	ID               int    `json:"id"`
	BaseAssetId      int    `json:"base_asset_id"`
	QuoteAssetId     int    `json:"quote_asset_id"`
	DecimalPrecision int    `json:"decimal_precision"`
	LotSize          string `json:"lot_size"`
	TickSize         string `json:"tick_size"`
	Active           bool   `json:"active"`
}

type OrderBook struct {
	Amount string `json:"amount"`
	Price  string `json:"price"`
	Value  string `json:"value"`
}

type OrderBooks struct {
	Asks []OrderBook `json:"asks"`
	Bids []OrderBook `json:"bids"`
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
