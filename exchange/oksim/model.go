package oksim

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

type JsonResponse struct {
	Code string          `json:"code"`
	Msg  string          `json:"message"`
	Data json.RawMessage `json:"data"`
}

type PairData struct {
	InstType  string `json:"instType"`
	InstID    string `json:"instId"`
	Uly       string `json:"uly"`
	Category  string `json:"category"`
	BaseCcy   string `json:"baseCcy"`
	QuoteCcy  string `json:"quoteCcy"`
	SettleCcy string `json:"settleCcy"`
	CtVal     string `json:"ctVal"`
	CtMult    string `json:"ctMult"`
	CtValCcy  string `json:"ctValCcy"`
	OptType   string `json:"optType"`
	Stk       string `json:"stk"`
	ListTime  string `json:"listTime"`
	ExpTime   string `json:"expTime"`
	Lever     string `json:"lever"`
	TickSz    string `json:"tickSz"`
	LotSz     string `json:"lotSz"`
	MinSz     string `json:"minSz"`
	CtType    string `json:"ctType"`
	Alias     string `json:"alias"`
	State     string `json:"state"`
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
