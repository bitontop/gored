package hitbtc

import "time"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type ErrResponse struct {
	Error struct {
		Code        int    `json:"code"`
		Message     string `json:"message"`
		Description string `json:"description"`
	} `json:"error"`
}

type CoinsData []struct {
	ID                 string `json:"id"`
	FullName           string `json:"fullName"`
	Crypto             bool   `json:"crypto"`
	PayinEnabled       bool   `json:"payinEnabled"`
	PayinPaymentID     bool   `json:"payinPaymentId"`
	PayinConfirmations int    `json:"payinConfirmations"`
	PayoutEnabled      bool   `json:"payoutEnabled"`
	PayoutIsPaymentID  bool   `json:"payoutIsPaymentId"`
	TransferEnabled    bool   `json:"transferEnabled"`
	Delisted           bool   `json:"delisted"`
	PayoutFee          string `json:"payoutFee"`
}

type PairsData []struct {
	ID                   string `json:"id"`
	BaseCurrency         string `json:"baseCurrency"`
	QuoteCurrency        string `json:"quoteCurrency"`
	QuantityIncrement    string `json:"quantityIncrement"`
	TickSize             string `json:"tickSize"`
	TakeLiquidityRate    string `json:"takeLiquidityRate"`
	ProvideLiquidityRate string `json:"provideLiquidityRate"`
	FeeCurrency          string `json:"feeCurrency"`
}

type OrderBook struct {
	Ask []struct {
		Price string `json:"price"`
		Size  string `json:"size"`
	} `json:"ask"`
	Bid []struct {
		Price string `json:"price"`
		Size  string `json:"size"`
	} `json:"bid"`
}

type AccountBalances []struct {
	Currency  string `json:"currency"`
	Available string `json:"available"`
	Reserved  string `json:"reserved"`
}

type PlaceOrder struct {
	ID            int64     `json:"id"`
	ClientOrderID string    `json:"clientOrderId"`
	Symbol        string    `json:"symbol"`
	Side          string    `json:"side"`
	Status        string    `json:"status"`
	Type          string    `json:"type"`
	TimeInForce   string    `json:"timeInForce"`
	Quantity      string    `json:"quantity"`
	Price         string    `json:"price"`
	CumQuantity   string    `json:"cumQuantity"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	PostOnly      bool      `json:"postOnly"`
}

type Withdraw struct {
	ID string `json:"id"`
}

type TradeHistory []struct {
	ID            int       `json:"id"`
	ClientOrderID string    `json:"clientOrderId"`
	OrderID       int       `json:"orderId"`
	Symbol        string    `json:"symbol"`
	Side          string    `json:"side"`
	Quantity      string    `json:"quantity"`
	Price         string    `json:"price"`
	Fee           string    `json:"fee"`
	Timestamp     time.Time `json:"timestamp"`
}
