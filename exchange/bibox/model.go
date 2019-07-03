package bibox

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

type JsonResponse struct {
	Error  Error           `json:"error"`
	Result json.RawMessage `json:"result"`
	Cmd    string          `json:"cmd"`
	Ver    string          `json:"ver"`
}

type Error struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

type CoinsData []struct {
	Result []struct {
		CoinSymbol       string  `json:"coin_symbol"`
		IsActive         int     `json:"is_active"`
		OriginalDecimals int     `json:"original_decimals"`
		EnableDeposit    int     `json:"enable_deposit"`
		EnableWithdraw   int     `json:"enable_withdraw"`
		WithdrawFee      float64 `json:"withdraw_fee"`
		WithdrawMin      float64 `json:"withdraw_min"`
	} `json:"result"`
	Cmd string `json:"cmd"`
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
	Result struct {
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
	} `json:"result"`
	Cmd string `json:"cmd"`
}

type PlaceOrder []struct {
	Result string `json:"result"`
	Cmd    string `json:"cmd"`
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

type CancelOrder []struct {
	Result string `json:"result"`
	Cmd    string `json:"cmd"`
}

type InnerTrans struct {
	Result struct {
		ID    string `json:"id"`
		State int    `json:"state"`
	} `json:"result"`
}

type Withdraw struct {
	Result int    `json:"result"`
	Cmd    string `json:"cmd"`
}
