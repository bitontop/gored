package bw

import "encoding/json"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct {
	Datas  json.RawMessage `json:"datas"`
	ResMsg struct {
		Message string `json:"message"`
		Method  string `json:"method"`
		Code    string `json:"code"`
	}
}

/********** Public API Structure**********/
// type CoinsData []struct {
// 	CurrencyID         string `json:"currencyId"`
// 	Name               string `json:"name"`
// 	Alias              string `json:"alias"`
// 	Logo               string `json:"logo"`
// 	Description        string `json:"description"`
// 	DescriptionEnglish string `json:"descriptionEnglish"`
// 	DefaultDecimal     int    `json:"defaultDecimal"`
// 	CreateUID          string `json:"createUid"`
// 	CreateTime         int64  `json:"createTime"`
// 	ModifyUID          string `json:"modifyUid"`
// 	ModifyTime         int64  `json:"modifyTime"`
// 	State              int    `json:"state"`
// 	Mark               string `json:"mark"`
// 	TotalNumber        string `json:"totalNumber"`
// 	PublishNumber      string `json:"publishNumber"`
// 	MarketValue        string `json:"marketValue"`
// 	IsLegalCoin        int    `json:"isLegalCoin"`
// 	NeedBlockURL       int    `json:"needBlockUrl"`
// 	BlockChainURL      string `json:"blockChainUrl"`
// 	TradeSearchURL     string `json:"tradeSearchUrl"`
// 	TokenCoinsID       int    `json:"tokenCoinsId"`
// 	IsMining           string `json:"isMining"`
// 	Arithmetic         string `json:"arithmetic"`
// 	Founder            string `json:"founder"`
// 	TeamAddress        string `json:"teamAddress"`
// 	Remark             string `json:"remark"`
// 	TokenName          string `json:"tokenName"`
// 	IsMemo             int    `json:"isMemo"`
// 	WebsiteCurrencyID  string `json:"websiteCurrencyId"`
// 	DrawFlag           int    `json:"drawFlag"`
// 	RechargeFlag       int    `json:"rechargeFlag"`
// 	DrawFee            string `json:"drawFee"`
// 	OnceDrawLimit      int    `json:"onceDrawLimit"`
// 	DailyDrawLimit     int    `json:"dailyDrawLimit"`
// 	TimesFreetrial     string `json:"timesFreetrial"`
// 	HourFreetrial      string `json:"hourFreetrial"`
// 	DayFreetrial       string `json:"dayFreetrial"`
// 	MinFee             string `json:"minFee"`
// 	InConfigTimes      int    `json:"inConfigTimes"`
// 	OutConfigTimes     int    `json:"outConfigTimes"`
// 	MinCash            string `json:"minCash"`
// 	LimitAmount        string `json:"limitAmount"`
// 	ZbExist            bool   `json:"zbExist"`
// }

type CoinsData []struct {
	CurrencyID         string      `json:"currencyId"`
	Name               string      `json:"name"`
	Alias              string      `json:"alias"`
	Logo               string      `json:"logo"`
	Description        string      `json:"description"`
	DescriptionEnglish string      `json:"descriptionEnglish"`
	DefaultDecimal     int         `json:"defaultDecimal"`
	CreateUID          interface{} `json:"createUid"`
	CreateTime         int64       `json:"createTime"`
	ModifyUID          interface{} `json:"modifyUid"`
	ModifyTime         int         `json:"modifyTime"`
	State              int         `json:"state"`
	Mark               string      `json:"mark"`
	TotalNumber        float64     `json:"totalNumber"`
	PublishNumber      float64     `json:"publishNumber"`
	MarketValue        float64     `json:"marketValue"`
	IsLegalCoin        int         `json:"isLegalCoin"`
	NeedBlockURL       int         `json:"needBlockUrl"`
	BlockChainURL      string      `json:"blockChainUrl"`
	TradeSearchURL     interface{} `json:"tradeSearchUrl"`
	TokenCoinsID       int         `json:"tokenCoinsId"`
	IsMining           string      `json:"isMining"`
	Arithmetic         interface{} `json:"arithmetic"`
	Founder            string      `json:"founder"`
	TeamAddress        interface{} `json:"teamAddress"`
	Remark             interface{} `json:"remark"`
	TokenName          string      `json:"tokenName"`
	IsMemo             int         `json:"isMemo"`
	WebsiteCurrencyID  string      `json:"websiteCurrencyId"`
	DrawFlag           int         `json:"drawFlag"`
	RechargeFlag       int         `json:"rechargeFlag"`
	DrawFee            float64     `json:"drawFee"`
	OnceDrawLimit      int         `json:"onceDrawLimit"`
	DailyDrawLimit     int         `json:"dailyDrawLimit"`
	TimesFreetrial     float64     `json:"timesFreetrial"`
	HourFreetrial      float64     `json:"hourFreetrial"`
	DayFreetrial       float64     `json:"dayFreetrial"`
	MinFee             float64     `json:"minFee"`
	InConfigTimes      int         `json:"inConfigTimes"`
	OutConfigTimes     int         `json:"outConfigTimes"`
	MinCash            float64     `json:"minCash"`
	LimitAmount        float64     `json:"limitAmount"`
	ZbExist            bool        `json:"zbExist"`
	Zone               int         `json:"zone"`
}

type PairsData []struct {
	MarketID         string `json:"marketId"`
	WebID            string `json:"webId"`
	ServerID         string `json:"serverId"`
	Name             string `json:"name"`
	LeverType        string `json:"leverType"`
	BuyerCurrencyID  string `json:"buyerCurrencyId"`
	SellerCurrencyID string `json:"sellerCurrencyId"`
	AmountDecimal    int    `json:"amountDecimal"`
	PriceDecimal     int    `json:"priceDecimal"`
	MinAmount        string `json:"minAmount"`
	State            int    `json:"state"`
	OpenTime         int64  `json:"openTime"`
	DefaultFee       string `json:"defaultFee"`
	CreateUID        string `json:"createUid"`
	CreateTime       int    `json:"createTime"`
	ModifyUID        string `json:"modifyUid"`
	ModifyTime       int64  `json:"modifyTime"`
	CombineMarketID  string `json:"combineMarketId"`
	IsCombine        int    `json:"isCombine"`
	IsMining         int    `json:"isMining"`
}

type OrderBook struct {
	Asks      [][]string `json:"asks"`
	Bids      [][]string `json:"bids"`
	Timestamp string     `json:"timestamp"`
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
