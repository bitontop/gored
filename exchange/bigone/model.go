package bigone

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
	"time"
)

type JsonResponse struct {
	Errors []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
	Data json.RawMessage `json:"data"`
}

type PairsData []struct {
	UUID       string `json:"uuid"`
	QuoteScale int    `json:"quoteScale"`
	QuoteAsset struct {
		UUID   string `json:"uuid"`
		Symbol string `json:"symbol"`
		Name   string `json:"name"`
	} `json:"quoteAsset"`
	Name      string `json:"name"`
	BaseScale int    `json:"baseScale"`
	BaseAsset struct {
		UUID   string `json:"uuid"`
		Symbol string `json:"symbol"`
		Name   string `json:"name"`
	} `json:"baseAsset"`
}

type OrderBook struct {
	MarketUUID string `json:"market_uuid"`
	MarketID   string `json:"market_id"`
	Bids       []struct {
		Price      string `json:"price"`
		OrderCount int    `json:"order_count"`
		Amount     string `json:"amount"`
	} `json:"bids"`
	Asks []struct {
		Price      string `json:"price"`
		OrderCount int    `json:"order_count"`
		Amount     string `json:"amount"`
	} `json:"asks"`
}

type AccountBalances []struct {
	LockedBalance string `json:"locked_balance"`
	Balance       string `json:"balance"`
	AssetUUID     string `json:"asset_uuid"`
	AssetID       string `json:"asset_id"`
}

type PlaceOrder struct {
	ID           string `json:"id"`
	MarketID     string `json:"market_id"`
	Price        string `json:"price"`
	Amount       string `json:"amount"`
	FilledAmount string `json:"filled_amount"`
	AvgDealPrice string `json:"avg_deal_price"`
	Side         string `json:"side"`
	State        string `json:"state"`
}

type Withdraw struct {
	Edges []struct {
		Node struct {
			ID            int       `json:"id"`
			CustomerID    string    `json:"customer_id"`
			AssetID       string    `json:"asset_id"`
			Amount        string    `json:"amount"`
			State         string    `json:"state"`
			Note          time.Time `json:"note"`
			Txid          string    `json:"txid"`
			CompletedAt   time.Time `json:"completed_at"`
			InsertedAt    time.Time `json:"inserted_at"`
			IsInternal    bool      `json:"is_internal"`
			TargetAddress string    `json:"target_address"`
		} `json:"node"`
		Cursor string `json:"cursor"`
	} `json:"edges"`
	PageInfo struct {
		EndCursor       string `json:"end_cursor"`
		StartCursor     string `json:"start_cursor"`
		HasNextPage     bool   `json:"has_next_page"`
		HasPreviousPage bool   `json:"has_previous_page"`
	} `json:"page_info"`
}
