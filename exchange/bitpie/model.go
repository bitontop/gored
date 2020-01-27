package bitpie

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"

	"github.com/bitontop/gored/exchange"
)

type JsonResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

type CoinsData []struct {
	Coins []struct {
		Chain                  string `json:"chain"`
		NumOfConfirmations     int    `json:"numOfConfirmations"`
		NumOfFastConfirmations int    `json:"numOfFastConfirmations"`
		DepositStatus          string `json:"depositStatus"`
		MinDepositAmt          string `json:"minDepositAmt"`
		WithdrawStatus         string `json:"withdrawStatus"`
		MinWithdrawAmt         string `json:"minWithdrawAmt"`
		WithdrawPrecision      int    `json:"withdrawPrecision"`
		MaxWithdrawAmt         string `json:"maxWithdrawAmt"`
		WithdrawQuotaPerDay    string `json:"withdrawQuotaPerDay"`
		WithdrawQuotaPerYear   string `json:"withdrawQuotaPerYear"`
		WithdrawQuotaTotal     string `json:"withdrawQuotaTotal"`
		WithdrawFeeType        string `json:"withdrawFeeType"`
		TransactFeeWithdraw    string `json:"transactFeeWithdraw"`
	} `json:"coins"`
}

type PairsData []struct {
	Name                string `json:"name"`
	MarketOrderMinMoney string `json:"market_order_min_money"`
	StockPrecision      int    `json:"stock_precision"`
	Money               string `json:"money"`
	MakerFeeRate        int    `json:"maker_fee_rate"`
	OrderMinVol         string `json:"order_min_vol"`
	MoneyPrecision      int    `json:"money_precision"`
	MarketOrderEnabled  bool   `json:"market_order_enabled"`
	Stock               string `json:"stock"`
	Enabled             bool   `json:"enabled"`
	TakerFeeRate        int    `json:"taker_fee_rate"`
}

//-------

// type PairsData []struct {
// 	MarketCurrency     string      `json:"MarketCurrency"`
// 	BaseCurrency       string      `json:"BaseCurrency"`
// 	MarketCurrencyLong string      `json:"MarketCurrencyLong"`
// 	BaseCurrencyLong   string      `json:"BaseCurrencyLong"`
// 	MinTradeSize       float64     `json:"MinTradeSize"`
// 	MarketName         string      `json:"MarketName"`
// 	IsActive           bool        `json:"IsActive"`
// 	Created            string      `json:"Created"`
// 	Notice             interface{} `json:"Notice"`
// 	IsSponsored        interface{} `json:"IsSponsored"`
// 	LogoURL            string      `json:"LogoUrl"`
// }

type AccountBalances []struct {
	Currency      string  `json:"Currency"`
	Balance       float64 `json:"Balance"`
	Available     float64 `json:"Available"`
	Pending       float64 `json:"Pending"`
	CryptoAddress string  `json:"CryptoAddress"`
	Requested     bool    `json:"Requested"`
	Uuid          string  `json:"Uuid"`
}

type Uuid struct {
	Id string `json:"uuid"`
}

type PlaceOrder struct {
	AccountId                  string
	OrderUuid                  string `json:"OrderUuid"`
	Exchange                   string `json:"Exchange"`
	Type                       string
	Quantity                   float64 `json:"Quantity"`
	QuantityRemaining          float64 `json:"QuantityRemaining"`
	Limit                      float64 `json:"Limit"`
	Reserved                   float64
	ReserveRemaining           float64
	CommissionReserved         float64
	CommissionReserveRemaining float64
	CommissionPaid             float64
	Price                      float64 `json:"Price"`
	PricePerUnit               float64 `json:"PricePerUnit"`
	Opened                     string
	Closed                     string
	IsOpen                     bool
	Sentinel                   string
	CancelInitiated            bool
	ImmediateOrCancel          bool
	IsConditional              bool
	Condition                  string
	ConditionTarget            float64
}

type OrderBook struct {
	Buy  []exchange.Order `json:"buy"`
	Sell []exchange.Order `json:"sell"`
}
