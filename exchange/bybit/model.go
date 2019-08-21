package bybit

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

type JsonResponse struct {
	RetCode int             `json:"ret_code"`
	RetMsg  string          `json:"ret_msg"`
	ExtCode string          `json:"ext_code"`
	ExtInfo string          `json:"ext_info"`
	TimeNow string          `json:"time_now"`
	Result  json.RawMessage `json:"result"`
}

type PairsData []struct {
	Name          string `json:"name"`
	BaseCurrency  string `json:"base_currency"`
	QuoteCurrency string `json:"quote_currency"`
	PriceScale    int    `json:"price_scale"`
	PriceFilter   struct {
		MinPrice string `json:"min_price"`
		MaxPrice string `json:"max_price"`
		TickSize string `json:"tick_size"`
	} `json:"price_filter"`
	LotSizeFilter struct {
		MaxTradingQty int     `json:"max_trading_qty"`
		MinTradingQty int     `json:"min_trading_qty"`
		QtyStep       float64 `json:"qty_step"`
	} `json:"lot_size_filter"`
}

/* type OrderBook struct {
	Buy  []exchange.Order `json:"buy"`
	Sell []exchange.Order `json:"sell"`
} */
