package tradeogre

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type PairsData []map[string]*PairsDetail

type PairsDetail struct {
	Initialprice string `json:"initialprice"`
	Price        string `json:"price"`
	High         string `json:"high"`
	Low          string `json:"low"`
	Volume       string `json:"volume"`
	Bid          string `json:"bid"`
	Ask          string `json:"ask"`
}

type OrderBook struct {
	Success string            `json:"success"` //true false
	Buy     map[string]string `json:"buy"`
	Sell    map[string]string `json:"sell"`
}

type AccountBalances struct {
	Success  bool              `json:"success"`
	Balances map[string]string `json:"balances"`
}

type PlaceOrder struct {
	Success      bool   `json:"success"`
	UUID         string `json:"uuid"`
	Bnewbalavail string `json:"bnewbalavail"`
	Snewbalavail string `json:"snewbalavail"`
}

type OrderStatus struct {
	Success   bool   `json:"success"`
	Date      string `json:"date"`
	Type      string `json:"type"`
	Market    string `json:"market"`
	Price     string `json:"price"`
	Quantity  string `json:"quantity"`
	Fulfilled string `json:"fulfilled"`
	Error     string `json:"error"`
}

type CancelOrder struct {
	Success bool `json:"success"`
}
