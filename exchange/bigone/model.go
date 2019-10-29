package bigone

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
	"time"
)

type JsonResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type PairsData []struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	QuoteScale int    `json:"quote_scale"`
	QuoteAsset struct {
		ID     string `json:"id"`
		Symbol string `json:"symbol"`
		Name   string `json:"name"`
	} `json:"quote_asset"`
	BaseAsset struct {
		ID     string `json:"id"`
		Symbol string `json:"symbol"`
		Name   string `json:"name"`
	} `json:"base_asset"`
	BaseScale     int    `json:"base_scale"`
	MinQuoteValue string `json:"min_quote_value"`
}

type OrderBook struct {
	AssetPairName string `json:"asset_pair_name"`
	Bids          []struct {
		Price      string `json:"price"`
		OrderCount int    `json:"order_count"`
		Quantity   string `json:"quantity"`
	} `json:"bids"`
	Asks []struct {
		Price      string `json:"price"`
		OrderCount int    `json:"order_count"`
		Quantity   string `json:"quantity"`
	} `json:"asks"`
}

type AccountBalances []struct {
	AssetSymbol   string `json:"asset_symbol"`
	Balance       string `json:"balance"`
	LockedBalance string `json:"locked_balance"`
}

type PlaceOrder struct {
	ID            int       `json:"id"`
	AssetPairName string    `json:"asset_pair_name"`
	Price         string    `json:"price"`
	Amount        string    `json:"amount"`
	FilledAmount  string    `json:"filled_amount"`
	AvgDealPrice  string    `json:"avg_deal_price"`
	Side          string    `json:"side"`
	State         string    `json:"state"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Withdraw []struct {
	ID          int         `json:"id"`
	CustomerID  int         `json:"customer_id"`
	AssetUUID   string      `json:"asset_uuid"`
	Amount      string      `json:"amount"`
	Recipient   interface{} `json:"recipient"`
	State       string      `json:"state"`
	IsInternal  bool        `json:"is_internal"`
	Note        string      `json:"note"`
	Kind        string      `json:"kind"`
	Txid        string      `json:"txid"`
	Confirms    int         `json:"confirms"`
	InsertedAt  interface{} `json:"inserted_at"`
	UpdatedAt   interface{} `json:"updated_at"`
	CompletedAt interface{} `json:"completed_at"`
	Commision   interface{} `json:"commision"`
	Explain     string      `json:"explain"`
}
