package oksim

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

type JsonResponse struct {
	Code         string          `json:"code"`
	DetailMsg    string          `json:"detailMsg"`
	ErrorCode    string          `json:"error_code"`
	ErrorMessage string          `json:"error_message"`
	Msg          string          `json:"msg"`
	Data         json.RawMessage `json:"data"`
}

type CoinData struct {
	CanDep bool   `json:"canDep"`
	CanWd  bool   `json:"canWd"`
	Ccy    string `json:"ccy"`
	Chain  string `json:"chain"`
	MinWd  int    `json:"minWd"`
	Name   string `json:"name"`
}

type PairData struct {
	Alias     string `json:"alias"`
	BaseCcy   string `json:"baseCcy"`
	Category  string `json:"category"`
	CtMult    string `json:"ctMult"`
	CtType    string `json:"ctType"`
	CtVal     string `json:"ctVal"`
	CtValCcy  string `json:"ctValCcy"`
	ExpTime   string `json:"expTime"`
	InstID    string `json:"instId"`
	InstType  string `json:"instType"`
	Lever     string `json:"lever"`
	ListTime  string `json:"listTime"`
	LotSz     string `json:"lotSz"`
	MinSz     string `json:"minSz"`
	OptType   string `json:"optType"`
	QuoteCcy  string `json:"quoteCcy"`
	SettleCcy string `json:"settleCcy"`
	State     string `json:"state"`
	Stk       string `json:"stk"`
	TickSz    string `json:"tickSz"`
	Uly       string `json:"uly"`
}

type OrderBook struct {
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
	Ts   string     `json:"ts"`
}

type AccountBalances struct {
	AvailBal  string `json:"availBal"`
	Bal       string `json:"bal"`
	Ccy       string `json:"ccy"`
	FrozenBal string `json:"frozenBal"`
}

type PlaceOrder struct {
	ClOrdID string `json:"clOrdId"`
	OrdID   string `json:"ordId"`
	Tag     string `json:"tag"`
	SCode   string `json:"sCode"`
	SMsg    string `json:"sMsg"`
}

type OrderStatus struct {
	InstType    string `json:"instType"`
	InstID      string `json:"instId"`
	Ccy         string `json:"ccy"`
	OrdID       string `json:"ordId"`
	ClOrdID     string `json:"clOrdId"`
	Tag         string `json:"tag"`
	Px          string `json:"px"`
	Sz          string `json:"sz"`
	Pnl         string `json:"pnl"`
	OrdType     string `json:"ordType"`
	Side        string `json:"side"`
	PosSide     string `json:"posSide"`
	TdMode      string `json:"tdMode"`
	AccFillSz   string `json:"accFillSz"`
	FillPx      string `json:"fillPx"`
	TradeID     string `json:"tradeId"`
	FillSz      string `json:"fillSz"`
	FillTime    string `json:"fillTime"`
	State       string `json:"state"`
	AvgPx       string `json:"avgPx"`
	Lever       string `json:"lever"`
	TpTriggerPx string `json:"tpTriggerPx"`
	TpOrdPx     string `json:"tpOrdPx"`
	SlTriggerPx string `json:"slTriggerPx"`
	SlOrdPx     string `json:"slOrdPx"`
	FeeCcy      string `json:"feeCcy"`
	Fee         string `json:"fee"`
	RebateCcy   string `json:"rebateCcy"`
	Rebate      string `json:"rebate"`
	Category    string `json:"category"`
	UTime       string `json:"uTime"`
	CTime       string `json:"cTime"`
}

type Transfer struct {
	TransID string `json:"transId"`
	Ccy     string `json:"ccy"`
	From    string `json:"from"`
	Amt     string `json:"amt"`
	To      string `json:"to"`
}
