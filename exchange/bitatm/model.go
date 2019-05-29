package bitatm

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
	"time"
)

type JsonResponse struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Ts   int             `json:"ts"`
	Data json.RawMessage `json:"Data"`
}

type AccountBalances []struct {
	Currency string `json:"currency"`
	Balance  int    `json:"balance"`
	Frozen   int    `json:"frozen"`
}

type Withdrawal struct {
	Succeed    bool  `json:"succeed"`
	Withdrawid int64 `json:"withdrawid"`
}

type PlaceOrder struct {
	Orderid           int       `json:"orderid"`
	Ordertype         int       `json:"ordertype"`
	Direction         int       `json:"direction"`
	Price             int       `json:"price"`
	Amount            int       `json:"amount"`
	Transactionamount int       `json:"transactionamount"`
	Fee               int       `json:"fee"`
	Symbol            string    `json:"symbol"`
	Orderstatus       int       `json:"orderstatus"`
	Updatetime        time.Time `json:"updatetime"`
	Createtime        time.Time `json:"createtime"`
	Basecurrency      string    `json:"basecurrency"`
	Quotecurrency     string    `json:"quotecurrency"`
}

type PairsData []struct {
	ID              int    `json:"id"`
	Basecurrency    string `json:"basecurrency"`
	Quotecurrency   string `json:"quotecurrency"`
	Symbol          string `json:"symbol"`
	Priceprecision  string `json:"priceprecision"`
	Amountprecision string `json:"amountprecision"`
}

type CoinsData []struct {
	ID           int    `json:"id"`
	Currencyname string `json:"currencyname"`
}

type OrderBook struct {
	Ts   int `json:"ts"`
	Bids []struct {
		Price  int `json:"price"`
		Amount int `json:"amount"`
	} `json:"bids"`
	Asks []struct {
		Price  int `json:"price"`
		Amount int `json:"amount"`
	} `json:"asks"`
}
