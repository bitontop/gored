package belfries

import "encoding/json"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct {
	Data      json.RawMessage `json:"data"`
	Message   string          `json:"message"`
	Timestamp int             `json:"timestamp"`
	Status    int             `json:"status"`
	IsSuccess bool            `json:"isSuccess"`
}

/********** Public API Structure**********/
type Country struct {
	NickName   string `json:"nickName"`
	Icon       string `json:"icon"`
	Iso        string `json:"iso"`
	PrivacyUrl string `json:"privacyUrl"`
	TermsUrl   string `json:"termsUrl"`
	PhoneCode  int    `json:"phoneCode"`
	ID         int    `json:"id"`
}

type CountryData struct {
	ExchangeCountries []Country
	Message           string `json:"message"`
	IsSuccess         bool   `json:"isSuccess"`
}

type Currency struct {
	Id                           int         `json:"id"`
	Name                         string      `json:"name"`
	Symbol                       string      `json:"symbol"`
	Type                         string      `json:"type"`
	NetworkFees                  float64     `json:"networkFees"`
	NoConfirmations              int         `json:"noConfirmations"`
	NoConfirmationsMerchant      int         `json:"noConfirmationsMerchant"`
	TxnFeePerKB                  float64     `json:"txnFeePerKB"`
	MinThreshold                 int         `json:"minThreshold"`
	MinTolerance                 int         `json:"minTolerance"`
	MaxThreshold                 int         `json:"maxThreshold"`
	MaxTolerance                 int         `json:"maxTolerance"`
	DeepFreezeTransferFeeAccount int         `json:"deepFreezeTransferFeeAccount"`
	MinWithdrawLimit             float64     `json:"minWithdrawLimit"`
	MaxWithdrawLimit             int         `json:"maxWithdrawLimit"`
	MaxWithdrawLimitPer24Hrs     int         `json:"maxWithdrawLimitPer24Hrs"`
	MaxDepositAmountLimitPerday  interface{} `json:"maxDepositAmountLimitPerday"`
	Image                        string      `json:"image"`
	Scale                        int         `json:"scale"`
	IsActive                     bool        `json:"isActive"`
	Priority                     interface{} `json:"priority"`
	ExchangeCode                 string      `json:"exchangeCode"`
	IsBaseCurrency               bool        `json:"isBaseCurrency"`
	UpdatedDate                  int         `json:"updatedDate"`
	CanCreateWallet              bool        `json:"canCreateWallet"`
	IsDeposit                    bool        `json:"isDeposit"`
	IsWithdrawal                 bool        `json:"isWithdrawal"`
	ErcToken                     bool        `json:"ercToken"`
}

type CoinsData []Currency

type Market struct {
	ViewScale1      string  `json:"viewScale1"`
	ViewScale2      string  `json:"viewScale2"`
	Instrument   string  `json:"instrument"`
	Id  string  `json:"id"`
}

type PairsData []Market

type Order struct {
	OrderType  string  `json:"orderType"`
	Instrument string  `json:"instrument"`
	MarketId   int     `json:"marketId"`
	Quantity   float64 `json:"quantity"`
	Price      float64 `json:"price"`
}

type OrderBook struct {
	SELL []Order `json:"SELL"`
	BUY  []Order `json:"BUY"`
}

/********** Private API Structure**********/
type AccountBalances []struct {
	Asset     string  `json:"asset"`
	Total     float64 `json:"total"`
	Available float64 `json:"available"`
	Locked    float64 `json:"locked"`
}

type WithdrawResponse struct {
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

type PlaceOrder struct {
	Symbol       string `json:"symbol"`
	OrderID      string `json:"orderId"`
	Side         string `json:"side"`
	Type         string `json:"type"`
	Price        string `json:"price"`
	AveragePrice string `json:"executedQty"`
	OrigQty      string `json:"origQty"`
	ExecutedQty  string `json:"executedQty"`
	Status       string `json:"status"`
	TimeInForce  string `json:"timeInForce"`
}
