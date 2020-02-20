package bitstamp

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

/* type JsonResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
} */

type CoinsData []struct {
	ID                  string `json:"id"`
	FullName            string `json:"fullName"`
	Crypto              bool   `json:"crypto"`
	DepositStatus       bool   `json:"depositStatus"`
	DepositConfirmation int    `json:"depositConfirmation"`
	WithdrawStatus      bool   `json:"withdrawStatus"`
	WithdrawFee         string `json:"withdrawFee"`
}

type PairsData []struct {
	BaseDecimals    int    `json:"base_decimals"`
	MinimumOrder    string `json:"minimum_order"`
	Name            string `json:"name"`
	CounterDecimals int    `json:"counter_decimals"`
	Trading         string `json:"trading"`
	URLSymbol       string `json:"url_symbol"`
	Description     string `json:"description"`
}

type OrderBook struct {
	Timestamp string     `json:"timestamp"`
	Bids      [][]string `json:"bids"`
	Asks      [][]string `json:"asks"`
}

type TradeHistory []struct {
	Date   string `json:"date"`
	Tid    string `json:"tid"`
	Price  string `json:"price"`
	Type   string `json:"type"`
	Amount string `json:"amount"`
}
