package huobiotc

import "encoding/json"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct {
	Code       int             `json:"code"`
	Message    string          `json:"message"`
	TotalCount int             `json:"totalCount"`
	PageSize   int             `json:"pageSize"`
	TotalPage  int             `json:"totalPage"`
	CurrPage   int             `json:"currPage"`
	Data       json.RawMessage `json:"data"`
	Success    bool            `json:"success"`
}

type OrderBook []struct {
	ID                int     `json:"id"`
	UID               int     `json:"uid"`
	UserName          string  `json:"userName"`
	MerchantLevel     int     `json:"merchantLevel"`
	CoinID            int     `json:"coinId"`
	Currency          int     `json:"currency"`
	TradeType         int     `json:"tradeType"`
	BlockType         int     `json:"blockType"`
	PayMethod         string  `json:"payMethod"`
	PayTerm           int     `json:"payTerm"`
	PayName           string  `json:"payName"`
	MinTradeLimit     float64 `json:"minTradeLimit"`
	MaxTradeLimit     float64 `json:"maxTradeLimit"`
	Price             float64 `json:"price"`
	TradeCount        float64 `json:"tradeCount"`
	IsOnline          bool    `json:"isOnline"`
	TradeMonthTimes   int     `json:"tradeMonthTimes"`
	OrderCompleteRate int     `json:"orderCompleteRate"`
	TakerLimit        int     `json:"takerLimit"`
	GmtSort           int64   `json:"gmtSort"`
}
