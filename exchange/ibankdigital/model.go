package ibankdigital

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

type JsonResponse struct {
	Status  string          `json:"status"`
	ErrCode string          `json:"err-code"`
	ErrMsg  string          `json:"err-msg"`
	Data    json.RawMessage `json:"data"`
}

type PairsData []struct {
	BaseCurrency    string `json:"base-currency"`
	QuoteCurrency   string `json:"quote-currency"`
	PricePrecision  int    `json:"price-precision"`
	AmountPrecision int    `json:"amount-precision"`
	SymbolPartition string `json:"symbol-partition"`
	Symbol          string `json:"symbol"`
	State           string `json:"state"`
}

type HuobiJsonResponse struct {
	Status  string          `json:"status"`
	Data    json.RawMessage `json:"data"`
	Tick    json.RawMessage `json:"tick"`
	ErrCode string          `json:"err-code"`
	ErrMsg  string          `json:"err-msg"`
}

type HuobiCoins []struct {
	CurrencyAddrWithTag     bool          `json:"currency-addr-with-tag"`
	FastConfirms            int           `json:"fast-confirms"`
	SafeConfirms            int           `json:"safe-confirms"`
	CurrencyType            string        `json:"currency-type"`
	SupportSites            []interface{} `json:"support-sites"`
	OtcEnable               int           `json:"otc-enable"`
	CountryDisabled         bool          `json:"country-disabled"`
	Tags                    interface{}   `json:"tags"`
	DepositEnabled          bool          `json:"deposit-enabled"`
	WithdrawEnabled         bool          `json:"withdraw-enabled"`
	WhiteEnabled            bool          `json:"white-enabled"`
	WithdrawPrecision       int           `json:"withdraw-precision"`
	CurrencyPartition       string        `json:"currency-partition"`
	QuoteCurrency           bool          `json:"quote-currency"`
	WithdrawMinAmount       string        `json:"withdraw-min-amount"`
	Weight                  int           `json:"weight"`
	Visible                 bool          `json:"visible"`
	ShowPrecision           string        `json:"show-precision"`
	DepositMinAmount        string        `json:"deposit-min-amount"`
	VisibleAssetsTimestamp  int64         `json:"visible-assets-timestamp"`
	DepositEnableTimestamp  int64         `json:"deposit-enable-timestamp"`
	WithdrawEnableTimestamp int64         `json:"withdraw-enable-timestamp"`
	Name                    string        `json:"name"`
	State                   string        `json:"state"`
	DisplayName             string        `json:"display-name"`
	DepositDesc             string        `json:"deposit-desc"`
	WithdrawDesc            string        `json:"withdraw-desc"`
	SuspendVisibleDesc      string        `json:"suspend-visible-desc"`
	SuspendDepositDesc      string        `json:"suspend-deposit-desc"`
	SuspendWithdrawDesc     string        `json:"suspend-withdraw-desc"`
	CurrencyAddrOneoff      bool          `json:"currency-addr-oneoff,omitempty"`
	Blockchains             string        `json:"blockchains,omitempty"`
}

type OrderBook struct {
	Status string `json:"status"`
	Ch     string `json:"ch"`
	Ts     int64  `json:"ts"`
	Tick   struct {
		Bids    [][]float64 `json:"bids"`
		Asks    [][]float64 `json:"asks"`
		Ts      int64       `json:"ts"`
		Version int64       `json:"version"`
	} `json:"tick"`
}

type AccountID []struct {
	ID     int    `json:"id"`
	Type   string `json:"type"`
	State  string `json:"state"`
	UserID int    `json:"user-id"`
}

type AccountBalances struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	State string `json:"state"`
	List  []struct {
		Currency string `json:"currency"`
		Type     string `json:"type"`
		Balance  string `json:"balance"`
	} `json:"list"`
	UserID int `json:"user-id"`
}

type PlaceOrder struct {
	Data string `json:"data"`
}

type OrderStatus struct {
	ID              int    `json:"id"`
	Symbol          string `json:"symbol"`
	AccountID       int    `json:"account-id"`
	Amount          string `json:"amount"`
	Price           string `json:"price"`
	CreatedAt       int64  `json:"created-at"`
	Type            string `json:"type"`
	FieldAmount     string `json:"field-amount"`
	FieldCashAmount string `json:"field-cash-amount"`
	FieldFees       string `json:"field-fees"`
	FinishedAt      int64  `json:"finished-at"`
	UserID          int    `json:"user-id"`
	Source          string `json:"source"`
	State           string `json:"state"`
	CanceledAt      int    `json:"canceled-at"`
	Exchange        string `json:"exchange"`
	Batch           string `json:"batch"`
}
