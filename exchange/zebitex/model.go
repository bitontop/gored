package zebitex

import (
	"encoding/json"
)

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

/********** Public API Structure**********/
type Tick struct {
	Name       string      `json:"name"`
	BaseUnit   string      `json:"base_unit"`
	QuoteUnit  string      `json:"quote_unit"`
	AskFee     float64     `json:"ask_fee"`
	BidFee     float64     `json:"bid_fee"`
	Low        string      `json:"low"`
	High       string      `json:"high"`
	Last       string      `json:"last"`
	VisualLow  string      `json:"visualLow"`
	VisualHigh string      `json:"visualHigh"`
	VisualLast string      `json:"visualLast"`
	At         int         `json:"at"`
	Open       string      `json:"open"`
	Volume     interface{} `json:"volume"` //接口返回有时是浮点,有时是字符串
	Market     string      `json:"market"`
	// Buy          float64     `json:"buy"`
	IsUpTend bool `json:"isUpTend"`
	// Sell         float64     `json:"sell"`
	Percent      string `json:"percent"`
	Change       string `json:"change"`
	VisualOpen   string `json:"visualOpen"`
	VisualVolume string `json:"visualVolume"`
	VisualBuy    string `json:"visualBuy"`
	VisualSell   string `json:"visualSell"`
}

type CoinsData map[string]Tick

type PairsData map[string]Tick

type OrderBook struct {
	Bids [][]interface{} `json:"bids"`
	Asks [][]interface{} `json:"asks"`
}

type TradeHistory []struct {
	Tid               int    `json:"tid"`
	Type              string `json:"type"`
	Date              string `json:"date"`
	Price             string `json:"price"`
	Amount            string `json:"amount"`
	VisualPrice       string `json:"visualPrice"`
	VisualAmount      string `json:"visualAmount"`
	VisualQuoteAmount string `json:"visualQuoteAmount"`
}

/********** Private API Structure**********/
type Balance struct {
	IsFiat               bool        `json:"isFiat"`
	Code                 string      `json:"code"`
	Title                string      `json:"title"`
	PaymentAddress       string      `json:"paymentAddress"`
	Balance              string      `json:"balance"`
	LockedBalance        string      `json:"lockedBalance"`
	PaymentAddressQrCode string      `json:"paymentAddressQrCode"`
	BankAccounts         interface{} `json:"bankAccounts"`
	IsDisabled           bool        `json:"isDisabled"`
}

type AccountBalances map[string]Balance

type WithdrawResponse struct {
	Error string `json:"error"`
}

type PlaceOrder struct {
	Id          int64  `json:"id"`
	Bid         string `json:"bid"`
	Ask         string `json:"ask"`
	Price       string `json:"price"`
	Volume      string `json:"volume"`
	ExecutedQty string `json:"executedQty"`
	OrdType     string `json:"ord_type"`
}

type OrderDetail struct {
	Id        int64  `json:"id"`
	Side      string `json:"side"`
	State     string `json:"state"`
	OrdType   string `json:"ordType"`
	Currency  string `json:"currency"`
	Price     string `json:"price"`
	Filled    string `json:"filled"`
	Amount    string `json:"amount"`
	Avg       string `json:"avg"`
	Total     string `json:"total"`
	UpdatedAt string `json:"updatedAt"`
	Pair      string `json:"pair"`
	BaseUnit  string `json:"baseUnit"`
	QuoteUnit string `json:"quoteUnit"`
}

type OrdersPage struct {
	Per        int64         `json:"per"`
	Items      []OrderDetail `json:"items"`
	NextCursor interface{}   `json:"nextCursor"`
}

type FundSource struct {
	Id          int64       `json:"id"`
	Iban        interface{} `json:"iban"`
	Label       string      `json:"label"`
	Address     string      `json:"address"`
	Currency    string      `json:"currency"`
	AccountName string      `json:"accountName"`
	Confirmed   bool        `json:"confirmed"`
}

type FundSources []FundSource
