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

type SubAccountBalances map[string](map[string]*AccountBalances)

type TickerPrice struct {
	Date   int64                   `json:"date"`
	Ticker map[string]TickerDetail `json:"ticker"`
}

type TickerDetail struct {
	Buy        string `json:"buy"`
	BuyAmount  string `json:"buy_amount"`
	Open       string `json:"open"`
	High       string `json:"high"`
	Last       string `json:"last"`
	Low        string `json:"low"`
	Sell       string `json:"sell"`
	SellAmount string `json:"sell_amount"`
	Vol        string `json:"vol"`
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
	Time int64      `json:"time"`
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

type TradeHistory []struct {
	ID     int    `json:"id"`
	Type   string `json:"type"`
	Price  string `json:"price"`
	Amount string `json:"amount"`
	Date   int    `json:"date"`
	DateMs int64  `json:"date_ms"`
}

type OpenOrders struct {
	Count    int `json:"count"`
	CurrPage int `json:"curr_page"`
	Data     []struct {
		Amount       string `json:"amount"`
		AssetFee     string `json:"asset_fee"`
		AvgPrice     string `json:"avg_price"`
		ClientID     string `json:"client_id"`
		CreateTime   int64  `json:"create_time"`
		DealAmount   string `json:"deal_amount"`
		DealFee      string `json:"deal_fee"`
		DealMoney    string `json:"deal_money"`
		FeeAsset     string `json:"fee_asset"`
		FeeDiscount  string `json:"fee_discount"`
		ID           int64  `json:"id"`
		Left         string `json:"left"`
		MakerFeeRate string `json:"maker_fee_rate"`
		Market       string `json:"market"`
		OrderType    string `json:"order_type"`
		Price        string `json:"price"`
		Status       string `json:"status"`
		TakerFeeRate string `json:"taker_fee_rate"`
		Type         string `json:"type"`
	} `json:"data"`
	HasNext bool `json:"has_next"`
}

type WithdrawHistory []struct {
	ActualAmount   string `json:"actual_amount"`
	Amount         string `json:"amount"`
	CoinAddress    string `json:"coin_address"`
	CoinType       string `json:"coin_type"`
	CoinWithdrawID int    `json:"coin_withdraw_id"`
	Confirmations  int    `json:"confirmations"`
	CreateTime     int64  `json:"create_time"`
	Status         string `json:"status"`
	TxFee          string `json:"tx_fee"`
	TxID           string `json:"tx_id"`
}

type DepositHistory []struct {
	ActualAmount        string `json:"actual_amount"`
	ActualAmountDisplay string `json:"actual_amount_display"`
	AddExplorer         string `json:"add_explorer"`
	Amount              string `json:"amount"`
	AmountDisplay       string `json:"amount_display"`
	CoinAddress         string `json:"coin_address"`
	CoinAddressDisplay  string `json:"coin_address_display"`
	CoinDepositID       int    `json:"coin_deposit_id"`
	CoinType            string `json:"coin_type"`
	Confirmations       int    `json:"confirmations"`
	CreateTime          int64  `json:"create_time"`
	Explorer            string `json:"explorer"`
	Remark              string `json:"remark"`
	SmartContractName   string `json:"smart_contract_name"`
	Status              string `json:"status"`
	StatusDisplay       string `json:"status_display"`
	TransferMethod      string `json:"transfer_method"`
	TxID                string `json:"tx_id"`
	TxIDDisplay         string `json:"tx_id_display"`
}

type TransferHistory struct {
	CurrPage int `json:"curr_page"`
	Data     []struct {
		Time         int64  `json:"time"`
		Amount       string `json:"amount"`
		CoinType     string `json:"coin_type"`
		TransferFrom string `json:"transfer_from"`
		TransferTo   string `json:"transfer_to"`
		Status       string `json:"status"`
	} `json:"data"`
	HasNext   bool `json:"has_next"`
	PerPage   int  `json:"per_page"`
	Total     int  `json:"total"`
	TotalPage int  `json:"total_page"`
}
