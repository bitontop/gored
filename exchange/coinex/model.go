package coinex

import "encoding/json"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type AccountBalances struct {
	Available string `json:"available"`
	Frozen    string `json:"frozen"`
}

type PlaceOrder struct {
	Amount       string `json:"amount"`
	AvgPrice     string `json:"avg_price"`
	CreateTime   int    `json:"create_time"`
	DealAmount   string `json:"deal_amount"`
	DealFee      string `json:"deal_fee"`
	DealMoney    string `json:"deal_money"`
	ID           int    `json:"id"`
	Left         string `json:"left"`
	MakerFeeRate string `json:"maker_fee_rate"`
	Market       string `json:"market"`
	OrderType    string `json:"order_type"`
	Price        string `json:"price"`
	Status       string `json:"status"`
	TakerFeeRate string `json:"taker_fee_rate"`
	Type         string `json:"type"`
}

type OrderBook struct {
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
	Last string     `json:"last"`
}

type PairsData struct {
	Symbol         string `json:"name"`
	MinAmount      string `json:"min_amount"`
	MakerFeeRate   string `json:"maker_fee_rate"`
	TakerFeeRate   string `json:"taker_fee_rate"`
	PricingName    string `json:"pricing_name"`
	PricingDecimal int    `json:"pricing_decimal"`
	TradingName    string `json:"trading_name"`
	TradingDecimal int    `json:"trading_decimal"`
}

type Withdraw struct {
	ActualAmount   string `json:"actual_amount"`
	Amount         string `json:"amount"`
	CoinAddress    string `json:"coin_address"`
	CoinType       string `json:"coin_type"`
	CoinWithdrawID int    `json:"coin_withdraw_id"`
	Confirmations  int    `json:"confirmations"`
	CreateTime     int    `json:"create_time"`
	Status         string `json:"status"`
	TxFee          string `json:"tx_fee"`
	TxID           string `json:"tx_id"`
}
