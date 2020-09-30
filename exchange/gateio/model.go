package gateio

import "encoding/json"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type CoinsData struct {
	Result string `json:"result"`
	Data   []struct {
		No          int         `json:"no"`
		Symbol      string      `json:"symbol"`
		Name        string      `json:"name"`
		NameEn      string      `json:"name_en"`
		NameCn      string      `json:"name_cn"`
		Pair        string      `json:"pair"`
		Rate        string      `json:"rate"`
		VolA        string      `json:"vol_a"`
		VolB        string      `json:"vol_b"`
		CurrA       string      `json:"curr_a"`
		CurrB       string      `json:"curr_b"`
		CurrSuffix  string      `json:"curr_suffix"`
		RatePercent string      `json:"rate_percent"`
		Trend       string      `json:"trend"`
		Supply      interface{} `json:"supply"`
		Marketcap   interface{} `json:"marketcap"`
		Lq          string      `json:"lq"`
	} `json:"data"`
}

type CoinsConstrain struct {
	Result string                  `json:"result"`
	Coins  []map[string]*Constrain `json:"coins"`
}

type Constrain struct {
	Delisted         int `json:"delisted"`
	WithdrawDisabled int `json:"withdraw_disabled"`
	WithdrawDelayed  int `json:"withdraw_delayed"`
	DepositDisabled  int `json:"deposit_disabled"`
	TradeDisabled    int `json:"trade_disabled"`
}

type PairsData struct {
	Result string             `json:"result"`
	Pairs  []map[string]*Pair `json:"pairs"`
}

type Pair struct {
	DecimalPlaces int     `json:"decimal_places"`
	MinAmount     float64 `json:"min_amount"`
	MinAmountA    float64 `json:"min_amount_a"`
	MinAmountB    float64 `json:"min_amount_b"`
	Fee           float64 `json:"fee"`
	TradeDisabled int     `json:"trade_disabled"`
}

type OrderBook struct {
	Elapsed string     `json:"elapsed"`
	Asks    [][]string `json:"asks"`
	Bids    [][]string `json:"bids"`
	Result  string     `json:"result"`
}

type AccountBalances struct {
	Result    string          `json:"result"`
	Available json.RawMessage `json:"available"`
	Locked    json.RawMessage `json:"locked"`
	Code      int             `json:"code"`
	Message   string          `json:"message"`
}

type FreeBalance struct {
	Result    string            `json:"result"`
	Available map[string]string `json:"available"`
	Locked    map[string]string `json:"locked"`
	Code      int               `json:"code"`
	Message   string            `json:"message"`
}

type PlaceOrder struct {
	Result        string  `json:"result"`
	Message       string  `json:"message"`
	Code          int     `json:"code"`
	Ctime         float64 `json:"ctime"`
	Side          int     `json:"side"`
	OrderNumber   int     `json:"orderNumber"`
	Rate          string  `json:"rate"`
	LeftAmount    string  `json:"leftAmount"`
	DealStock     string  `json:"deal_stock"`
	DealMoney     string  `json:"deal_money"`
	FilledAmount  string  `json:"filledAmount"`
	FilledRate    string  `json:"filledRate"`
	FeePercentage float64 `json:"feePercentage"`
	FeeValue      string  `json:"feeValue"`
	FeeCurrency   string  `json:"feeCurrency"`
	Fee           string  `json:"fee"`
}

type OrderStatus struct {
	Result string `json:"result"`
	Order  struct {
		ID            string `json:"id"`
		Status        string `json:"status"`
		Pair          string `json:"pair"`
		Type          string `json:"type"`
		Rate          string `json:"rate"`
		Amount        string `json:"amount"`
		InitialRate   string `json:"initial_rate"`
		InitialAmount string `json:"initial_amount"`
	} `json:"order"`
	Message string `json:"message"`
}

type CancelOrder struct {
	Result bool `json:"result"`
	Order  struct {
		OrderNumber   string  `json:"orderNumber"`
		Status        string  `json:"status"`
		CurrencyPair  string  `json:"currencyPair"`
		Type          string  `json:"type"`
		Rate          string  `json:"rate"`
		Amount        string  `json:"amount"`
		InitialRate   string  `json:"initialRate"`
		InitialAmount string  `json:"initialAmount"`
		FilledAmount  string  `json:"filledAmount"`
		FilledRate    string  `json:"filledRate"`
		FeePercentage float64 `json:"feePercentage"`
		FeeValue      string  `json:"feeValue"`
		FeeCurrency   string  `json:"feeCurrency"`
		Fee           string  `json:"fee"`
		Timestamp     int     `json:"timestamp"`
	} `json:"order"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Elapsed string `json:"elapsed"`
}

type WithdrawResponse struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

type TradeHistory struct {
	Elapsed string `json:"elapsed"`
	Result  string `json:"result"`
	Data    []struct {
		TradeID   string `json:"tradeID"`
		Total     string `json:"total"`
		Date      string `json:"date"`
		Rate      string `json:"rate"`
		Amount    string `json:"amount"`
		Timestamp string `json:"timestamp"`
		Type      string `json:"type"`
	} `json:"data"`
}
