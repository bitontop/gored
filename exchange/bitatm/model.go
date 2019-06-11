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
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
	Frozen   float64 `json:"frozen"`
}

type Withdrawal struct {
	Succeed    bool  `json:"succeed"`
	Withdrawid int64 `json:"withdrawid"`
}

type OrderStatus struct {
	Orderid           int64     `json:"orderid"`
	Ordertype         int32     `json:"ordertype"`
	Direction         int32     `json:"direction"`
	Price             float64   `json:"price"`
	Amount            float64   `json:"amount"`
	Transactionamount float64   `json:"transactionamount"`
	Fee               float64   `json:"fee"`
	Symbol            string    `json:"symbol"`
	Orderstatus       int32     `json:"orderstatus"`
	Updatetime        time.Time `json:"updatetime"`
	Createtime        time.Time `json:"createtime"`
	Basecurrency      string    `json:"basecurrency"`
	Quotecurrency     string    `json:"quotecurrency"`
}

type PairsData []struct {
	ID              string `json:"id"`
	Basecurrency    string `json:"basecurrency"`
	Quotecurrency   string `json:"quotecurrency"`
	Symbol          string `json:"symbol"`
	Priceprecision  string `json:"priceprecision"`
	Amountprecision string `json:"amountprecision"`
}

type CoinsData []struct {
	ID           string `json:"id"`
	Currencyname string `json:"currencyname"`
}

type OrderBook struct {
	Ts   int `json:"ts"`
	Bids []struct {
		Price  float64 `json:"price"`
		Amount float64 `json:"amount"`
		ID     int     `json:"id"`
	} `json:"bids"`
	Asks []struct {
		Price  float64 `json:"price"`
		Amount float64 `json:"amount"`
		ID     int     `json:"id"`
	} `json:"asks"`
}
