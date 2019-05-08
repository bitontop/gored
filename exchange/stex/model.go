package stex

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

type JsonResponse struct {
	Success int             `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   string          `json:"error"`
	Notice  string          `json:"notice"`
	Message string          `json:"message"`
}

type JsonResponseV3 struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
	Err     interface{}     `json:"errors"`
}

type CoinsData []struct {
	ID                        int    `json:"id"`
	Code                      string `json:"code"`
	Name                      string `json:"name"`
	Active                    bool   `json:"active"`
	Delisted                  bool   `json:"delisted"`
	Precision                 int    `json:"precision"`
	MinimumWithdrawalAmount   string `json:"minimum_withdrawal_amount"`
	MinimumDepositAmount      string `json:"minimum_deposit_amount"`
	DepositFeeCurrencyID      int    `json:"deposit_fee_currency_id"`
	DepositFeeCurrencyCode    string `json:"deposit_fee_currency_code"`
	DepositFeeConst           string `json:"deposit_fee_const"`
	DepositFeePercent         string `json:"deposit_fee_percent"`
	WithdrawalFeeCurrencyID   int    `json:"withdrawal_fee_currency_id"`
	WithdrawalFeeCurrencyCode string `json:"withdrawal_fee_currency_code"`
	WithdrawalFeeConst        string `json:"withdrawal_fee_const"`
	WithdrawalFeePercent      string `json:"withdrawal_fee_percent"`
	BlockExplorerURL          string `json:"block_explorer_url"`
}

type PairsData []struct {
	ID                int         `json:"id"`
	CurrencyID        int         `json:"currency_id"`
	CurrencyCode      string      `json:"currency_code"`
	CurrencyName      string      `json:"currency_name"`
	MarketCurrencyID  int         `json:"market_currency_id"`
	MarketCode        string      `json:"market_code"`
	MarketName        string      `json:"market_name"`
	MinOrderAmount    string      `json:"min_order_amount"`
	MinBuyPrice       string      `json:"min_buy_price"`
	MinSellPrice      string      `json:"min_sell_price"`
	BuyFeePercent     string      `json:"buy_fee_percent"`
	SellFeePercent    string      `json:"sell_fee_percent"`
	Active            bool        `json:"active"`
	Delisted          bool        `json:"delisted"`
	PairMessage       interface{} `json:"pair_message"`
	CurrencyPrecision int         `json:"currency_precision"`
	MarketPrecision   int         `json:"market_precision"`
	Symbol            string      `json:"symbol"`
	GroupName         string      `json:"group_name"`
	GroupID           int         `json:"group_id"`
}

type OrderBook struct {
	Ask []struct {
		CurrencyPairID   int     `json:"currency_pair_id"`
		Amount           string  `json:"amount"`
		Price            string  `json:"price"`
		Amount2          string  `json:"amount2"`
		Count            int     `json:"count"`
		CumulativeAmount float64 `json:"cumulative_amount"`
	} `json:"ask"`
	Bid []struct {
		CurrencyPairID   int     `json:"currency_pair_id"`
		Amount           string  `json:"amount"`
		Price            string  `json:"price"`
		Amount2          string  `json:"amount2"`
		Count            int     `json:"count"`
		CumulativeAmount float64 `json:"cumulative_amount"`
	} `json:"bid"`
}

type AccountBalances struct {
	Email        string            `json:"email"`
	Username     string            `json:"username"`
	Hash         string            `json:"hash"`
	IntercomHash string            `json:"intercom_hash"`
	Funds        map[string]string `json:"funds"`
	OpenOrders   int               `json:"open_orders"`
	ServerTime   int               `json:"server_time"`
}

type EmptyAccount struct {
	Success int `json:"success"`
	Data    struct {
		Email        string        `json:"email"`
		Username     string        `json:"username"`
		Hash         string        `json:"hash"`
		IntercomHash string        `json:"intercom_hash"`
		Funds        []interface{} `json:"funds"`
		OpenOrders   int           `json:"open_orders"`
		ServerTime   int           `json:"server_time"`
	} `json:"data"`
}

type Withdraw struct {
	Code                  string `json:"code"`
	ID                    int    `json:"id"`
	Amount                string `json:"amount"`
	Address               string `json:"address"`
	WithdrawalFee         string `json:"withdrawal_fee"`
	WithdrawalFeeCurrency string `json:"withdrawal_fee_currency"`
	Token                 string `json:"token"`
	Date                  struct {
		Date         string `json:"date"`
		TimezoneType int    `json:"timezone_type"`
		Timezone     string `json:"timezone"`
	} `json:"date"`
	Msg string `json:"msg"`
}

type CancelOrder struct {
	Funds   map[string]string `json:"funds"`
	OrderID string            `json:"order_id"`
}

type TradeDetail struct {
	Funds   map[string]string `json:"funds"`
	OrderID int64             `json:"order_id"`
}

type ActiveOrder map[string]*OrderDetail

type OrderDetail struct {
	Pair           string                    `json:"pair"`
	Type           string                    `json:"type"`
	OriginalAmount string                    `json:"original_amount"`
	BuyAmount      interface{}               `json:"buy_amount"`
	SellAmount     interface{}               `json:"sell_amount"`
	IsYourOrder    int                       `json:"is_your_order"`
	Timestamp      int                       `json:"timestamp"`
	Rates          map[string]*BuySellAmount `json:"rates"`
}

type BuySellAmount struct {
	BuyAmount  interface{} `json:"buy_amount"`
	SellAmount interface{} `json:"sell_amount"`
}
