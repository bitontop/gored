package bgogo

import "encoding/json"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

/********** Public API Structure**********/
type Ticker struct {
	LastPrice              string `json:"last_price"`
	LowestAskPrice         string `json:"lowest_ask_price"`
	HighestBidPrice        string `json:"highest_bid_price"`
	Past24hrsPriceChange   string `json:"past_24hrs_price_change"`
	Past24hrsBaseVolume    string `json:"past_24hrs_base_volume"`
	Past24hrsQuoteTurnover string `json:"past_24hrs_quote_turnover"`
	Past24hrsHighPrice     string `json:"past_24hrs_high_price"`
	Past24hrsLowPrice      string `json:"past_24hrs_low_price"`
}

type CoinsData map[string]Ticker

type PairsData map[string]Ticker

type OrderItem struct {
	Price  string `json:"price"`
	Amount string `json:"amount"`
}

type OrderBook struct {
	Bids []OrderItem `json:"bids"`
	Asks []OrderItem `json:"asks"`
}

type SnapshotJson struct {
	StatusCode int             `json:"status_code"`
	Message    string          `json:"message"`
	Time       int             `json:"time"`
	Data       json.RawMessage `json:"data"`
}

type SnapshotData struct {
	PriceStep                        string      `json:"price_step"`
	AmountStep                       string      `json:"amount_step"`
	AllSymbols                       interface{} `json:"all_symbols"`
	LastPrices                       interface{} `json:"last_prices"`
	Past24hrsPriceChanges            interface{} `json:"past_24hrs_price_changes"`
	Past24hrsHighPrice               interface{} `json:"past_24hrs_high_price"`
	Past24hrsLowPrice                interface{} `json:"past_24hrs_low_price"`
	Past24hrsVolumes                 interface{} `json:"past_24hrs_volumes"`
	Past24hrsTurnovers               interface{} `json:"past_24hrs_turnovers"`
	OrderBooks                       OrderBook   `json:"order_book"`
	TradeHistory                     interface{} `json:"trade_history"`
	QuoteCurrencyToFiatCurrencyPrice string      `json:"quote_currency_to_fiat_currency_price"`
	FiatCurrency                     string      `json:"fiat_currency"`
	MyAccountBalances                interface{} `json:"my_account_balances"`
	MyOrders                         interface{} `json:"my_orders"`
	Superpower                       bool        `json:"superpower"`
	FeeRate                          string      `json:"fee_rate"`
	EstimatedBtc                     string      `json:"estimated_btc"`
	EstimatedUsd                     string      `json:"estimated_usd"`
	MyFeeRate                        string      `json:"my_fee_rate"`
	BaseIntroLink                    string      `json:"base_intro_link"`
	Quota                            string      `json:"quota"`
	NextTimes                        int         `json:"next_time_s"`
	IeoStatus                        int         `json:"ieo_status"`
	QuotaCurrency                    string      `json:"quota_currency"`
	Category                         interface{} `json:"category"`
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
