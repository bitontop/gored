package bitz

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

/* type JsonResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
} */

type JsonResponse struct {
	Status    int64           `json:"status"`
	Msg       string          `json:"msg"`
	Data      json.RawMessage `json:"data"`
	Time      int             `json:"time"`
	Microtime string          `json:"microtime"`
	Source    string          `json:"source"`
}

type PairsData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	CoinFrom    string `json:"coinFrom"`
	CoinTo      string `json:"coinTo"`
	NumberFloat string `json:"numberFloat"`
	PriceFloat  string `json:"priceFloat"`
	Status      string `json:"status"`
	MinTrade    string `json:"minTrade"`
	MaxTrade    string `json:"maxTrade"`
	BuyFree     string `json:"buyFree"`
	SellFree    string `json:"sellFree"`
}

type OrderBook struct {
	Asks     [][]string `json:"asks"`
	Bids     [][]string `json:"bids"`
	CoinPair string     `json:"coinPair"`
}

type AccountBalances []struct {
	Name string `json:"name"`
	Num  string `json:"num"`
	Over string `json:"over"`
	Lock string `json:"lock"`
	Btc  string `json:"btc"`
	Usd  string `json:"usd"`
	Cny  string `json:"cny"`
}

type UserInfo struct {
	CNY       interface{}     `json:"cny"` //This may be a string
	USD       interface{}     `json:"usd"`
	BTC_TOTAL interface{}     `json:"btc_total"`
	Info      json.RawMessage `json:"info"`
}

type PlaceOrder struct {
	ID         int64       `json:"id"`
	UID        int         `json:"uId"`
	Price      string      `json:"price"`
	Number     string      `json:"number"`
	NumberOver string      `json:"numberOver"`
	NumberDeal interface{} `json:"numberDeal"`
	Flag       string      `json:"flag"`
	Status     int         `json:"status"`
	CoinFrom   string      `json:"coinFrom"`
	CoinTo     string      `json:"coinTo"`
}

type OrderDetails struct {
	ID              string `json:"id"`
	UID             string `json:"uId"`
	Price           string `json:"price"`
	Number          string `json:"number"`
	Total           string `json:"total"`
	NumberOver      string `json:"numberOver"`
	NumberDeal      string `json:"numberDeal"`
	Flag            string `json:"flag"`
	Status          int    `json:"status"`
	CoinFrom        string `json:"coinFrom"`
	CoinTo          string `json:"coinTo"`
	AveragePrice    string `json:"averagePrice"`
	OrderTotalPrice string `json:"orderTotalPrice"`
	TradeType       string `json:"tradeType"`
	Created         string `json:"created"`
}

type CancelOrder struct {
	UpdateAssetsData struct {
		Coin string `json:"coin"`
		Over string `json:"over"`
		Lock string `json:"lock"`
	} `json:"updateAssetsData"`
	AssetsInfo struct {
		Coin string `json:"coin"`
		Over string `json:"over"`
		Lock string `json:"lock"`
	} `json:"assetsInfo"`
}

type Withdraw struct {
	ID         int    `json:"id"`
	Email      string `json:"email"`
	Coin       string `json:"coin"`
	NetworkFee string `json:"network_fee"`
	Eid        int    `json:"eid"`
}

type TradeHistory struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Data   []struct {
		ID int64 `json:"id"`
		// t  string `json:"t"`
		T0 string `json:"t"`
		T  int64  `json:"T"`
		P  string `json:"p"`
		N  string `json:"n"`
		S  string `json:"s"`
	} `json:"data"`
	Time      int    `json:"time"`
	Microtime string `json:"microtime"`
	Source    string `json:"source"`
}
