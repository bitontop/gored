package bigone

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
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
	ID           int    `json:"id"`
	MarketID     string `json:"market_id"`
	Price        string `json:"price"`
	Amount       string `json:"amount"`
	FilledAmount string `json:"filled_amount"`
	AvgDealPrice string `json:"avg_deal_price"`
	Side         string `json:"side"`
	State        string `json:"state"`
}
