package bibox

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

/* type JsonResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
} */

type Error struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

type JsonResponse struct {
	Error struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
	} `json:"error"`
	Result json.RawMessage `json:"result"`
	Cmd    string          `json:"cmd"`
	Ver    string          `json:"ver"`
}

type PairsData []struct {
	ID       int    `json:"id"`
	Pair     string `json:"pair"`
	PairType int    `json:"pair_type"`
	AreaID   int    `json:"area_id"`
	IsHide   int    `json:"is_hide"`
}

type OrderBook struct {
	UpdateTime int64 `json:"update_time"`
	Asks       []struct {
		Volume string `json:"volume"`
		Price  string `json:"price"`
	} `json:"asks"`
	Bids []struct {
		Volume string `json:"volume"`
		Price  string `json:"price"`
	} `json:"bids"`
	Pair string `json:"pair"`
}

type AccountBalances []struct {
	TotalBtc   string `json:"total_btc"`
	TotalCny   string `json:"total_cny"`
	TotalUsd   string `json:"total_usd"`
	AssetsList []struct {
		CoinSymbol string `json:"coin_symbol"`
		Balance    string `json:"balance"`
		Freeze     string `json:"freeze"`
		BTCValue   string `json:"BTCValue"`
		CNYValue   string `json:"CNYValue"`
		USDValue   string `json:"USDValue"`
	} `json:"assets_list"`
}

type Asset struct {
	Cmd  string       `json:"cmd"`
	Body *AssetDetail `json:"body"`
}

type AssetDetail struct {
	Select int `json:"select"`
}

type OrderParam struct {
	Cmd         string             `json:"cmd"`
	Index       int                `json:"index"`
	BodyDetails *OrderParamDetails `json:"body"`
}

type OrderParamDetails struct {
	Pair        string  `json:"pair"`
	AccountType int     `json:"account_type"`
	OrderType   int     `json:"order_type"`
	OrderSide   int     `json:"order_side"`
	Price       float64 `json:"price"`
	Amount      float64 `json:"amount"`
}

type PlaceOrder []struct {
	Error struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
	} `json:"error"`
	Result int    `json:"result"`
	Cmd    string `json:"cmd"`
	Index  int    `json:"index"`
}

type OrderStatus []struct {
	Result struct {
		ID             int    `json:"id"`
		CreatedAt      int64  `json:"createdAt"`
		AccountType    int    `json:"account_type"`
		Pair           string `json:"pair"`
		CoinSymbol     string `json:"coin_symbol"`
		CurrencySymbol string `json:"currency_symbol"`
		OrderSide      int    `json:"order_side"`
		OrderType      int    `json:"order_type"`
		Price          string `json:"price"`
		DealPrice      string `json:"deal_price"`
		Amount         string `json:"amount"`
		Money          string `json:"money"`
		DealAmount     string `json:"deal_amount"`
		DealPercent    string `json:"deal_percent"`
		DealMoney      string `json:"deal_money"`
		Status         int    `json:"status"`
		Unexecuted     string `json:"unexecuted"`
		OrderFrom      int    `json:"order_from"`
	} `json:"result"`
	Cmd string `json:"cmd"`
}

type StatusParam struct {
	Cmd         string              `json:"cmd"`
	BodyDetails *StatusParamDetails `json:"body"`
}
type StatusParamDetails struct {
	Id int `json:"id"`
}

type CancelOrder []struct {
	Result string `json:"result"`
	Cmd    string `json:"cmd"`
	Index  int    `json:"index"`
}

type CancelParam struct {
	Cmd         string              `json:"cmd"`
	Index       int                 `json:"index"`
	BodyDetails *CancelParamDetails `json:"body"`
}
type CancelParamDetails struct {
	OrdersId int `json:"orders_id"`
}
