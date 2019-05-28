package bitmart

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type CoinsData []struct {
	Name            string `json:"name"`
	WithdrawEnabled bool   `json:"withdraw_enabled"`
	ID              string `json:"id"`
	DepositEnabled  bool   `json:"deposit_enabled"`
}

type PairsData []struct {
	ID                string `json:"id"`
	BaseCurrency      string `json:"base_currency"`
	QuoteCurrency     string `json:"quote_currency"`
	QuoteIncrement    string `json:"quote_increment"`
	BaseMinSize       string `json:"base_min_size"`
	BaseMaxSize       string `json:"base_max_size"`
	PriceMinPrecision int    `json:"price_min_precision"`
	PriceMaxPrecision int    `json:"price_max_precision"`
	Expiration        string `json:"expiration"`
}

type OrderBook struct {
	Buys []struct {
		Amount string `json:"amount"`
		Total  string `json:"total"`
		Price  string `json:"price"`
		Count  string `json:"count"`
	} `json:"buys"`
	Sells []struct {
		Amount string `json:"amount"`
		Total  string `json:"total"`
		Price  string `json:"price"`
		Count  string `json:"count"`
	} `json:"sells"`
}

type AccountBalances []struct {
	Name      string `json:"name"`
	Available string `json:"available"`
	Frozen    string `json:"frozen"`
	ID        string `json:"id"`
}

type PlaceOrder struct {
	Message   string `json:"message"`
	EntrustID int    `json:"entrust_id"`
}

type OrderStatus struct {
	EntrustID       int    `json:"entrust_id"`
	Symbol          string `json:"symbol"`
	Timestamp       int64  `json:"timestamp"`
	Side            string `json:"side"`
	Price           string `json:"price"`
	Fees            string `json:"fees"`
	OriginalAmount  string `json:"original_amount"`
	ExecutedAmount  string `json:"executed_amount"`
	RemainingAmount string `json:"remaining_amount"`
	Status          int    `json:"status"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
