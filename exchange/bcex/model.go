package bcex

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

type JsonResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"msg"`
	Data    json.RawMessage `json:"data"`
}

type CoinsData struct {
	CalcPrecision  int    `json:"calc_precision"`
	BizPrecision   int    `json:"biz_precision"`
	PhyPrecision   int    `json:"phy_precision"`
	OrderAmountMin string `json:"order_amount_min"`
	OrderAmountMax string `json:"order_amount_max"`
}

type PairsData struct {
	Main map[string][]*PairDetail `json:"main"`
	// Hot       []*PairDetail            `json:"hot"`
	// ByPercent []*PairDetail            `json:"by_percent"`
	// ByTime    []*PairDetail            `json:"by_time"`
}

type PairDetail struct {
	ID            int         `json:"id"`
	Vol           interface{} `json:"vol"`
	Market        string      `json:"market"`
	MarketAs      string      `json:"market_as"`
	Token         string      `json:"token"`
	TokenAs       string      `json:"token_as"`
	Amount        interface{} `json:"amount"`
	MaxPrice      string      `json:"max_price"`
	MinPrice      string      `json:"min_price"`
	OpenPrice     string      `json:"open_price"`
	IsUp          string      `json:"is_up"`
	Last          string      `json:"last"`
	PrevPrice     string      `json:"prev_price"`
	PPrecision    string      `json:"p_precision"`
	NPrecision    string      `json:"n_precision"`
	LatestPrice   string      `json:"latest_price"`
	PercentChange string      `json:"percent_change"`
	Currency      struct {
		CNY string `json:"CNY"`
		KRW string `json:"KRW"`
	} `json:"currency"`
}

type OrderBook struct {
	Bids [][]string `json:"bids"`
	Asks [][]string `json:"asks"`
}

// this exchange often return different type for the same field
type AccountBalances struct {
	Count int `json:"count"`
	Data  []struct {
		// UserID  interface{} `json:"user_id"`
		// OrgID   int         `json:"org_id"`
		Token string `json:"token"`
		// TokenAs string      `json:"token_as"`
		Usable string `json:"usable"`
		Locked string `json:"locked"`
		Total  string `json:"total"`
	} `json:"data"`
}

type PlaceOrder struct {
	ID            int    `json:"id"`
	Status        int    `json:"status"`
	OrderNo       string `json:"order_no"`
	Type          string `json:"type"`
	MarketType    string `json:"market_type"`
	AccountID     int    `json:"account_id"`
	OrgID         int    `json:"org_id"`
	UserID        int    `json:"user_id"`
	Market        string `json:"market"`
	Token         string `json:"token"`
	Price         string `json:"price"`
	Amount        string `json:"amount"`
	Volume        string `json:"volume"`
	MatchedAmount string `json:"matched_amount"`
	MatchedVolume string `json:"matched_volume"`
	AvgPrice      string `json:"avg_price"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}
