package stex

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

type JsonResponse struct {
	Success int             `json:"success"`
	Data    json.RawMessage `json:"data"`
	Error   string          `json:"error"`
	Notice  string          `json:"notice"`
	Message string          `json:"message"`
}

type JsonResponseV3 struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
	Err     interface{}     `json:"errors"`
}

type CoinsData []struct {
	ID                        int         `json:"id"`
	Code                      string      `json:"code"`
	Name                      string      `json:"name"`
	Active                    bool        `json:"active"`
	Delisted                  bool        `json:"delisted"`
	Precision                 int         `json:"precision"`
	MinimumTxConfirmations    int         `json:"minimum_tx_confirmations"`
	MinimumWithdrawalAmount   string      `json:"minimum_withdrawal_amount"`
	MinimumDepositAmount      string      `json:"minimum_deposit_amount"`
	DepositFeeCurrencyID      int         `json:"deposit_fee_currency_id"`
	DepositFeeCurrencyCode    string      `json:"deposit_fee_currency_code"`
	DepositFeeConst           string      `json:"deposit_fee_const"`
	DepositFeePercent         string      `json:"deposit_fee_percent"`
	WithdrawalFeeCurrencyID   int         `json:"withdrawal_fee_currency_id"`
	WithdrawalFeeCurrencyCode string      `json:"withdrawal_fee_currency_code"`
	WithdrawalFeeConst        string      `json:"withdrawal_fee_const"`
	WithdrawalFeePercent      string      `json:"withdrawal_fee_percent"`
	BlockExplorerURL          string      `json:"block_explorer_url"`
	ProtocolSpecificSettings  interface{} `json:"protocol_specific_settings"`
}

type PairsData []struct {
	ID                int    `json:"id"`
	CurrencyID        int    `json:"currency_id"`
	CurrencyCode      string `json:"currency_code"`
	CurrencyName      string `json:"currency_name"`
	MarketCurrencyID  int    `json:"market_currency_id"`
	MarketCode        string `json:"market_code"`
	MarketName        string `json:"market_name"`
	MinOrderAmount    string `json:"min_order_amount"`
	MinBuyPrice       string `json:"min_buy_price"`
	MinSellPrice      string `json:"min_sell_price"`
	BuyFeePercent     string `json:"buy_fee_percent"`
	SellFeePercent    string `json:"sell_fee_percent"`
	Active            bool   `json:"active"`
	Delisted          bool   `json:"delisted"`
	PairMessage       string `json:"pair_message"`
	CurrencyPrecision int    `json:"currency_precision"`
	MarketPrecision   int    `json:"market_precision"`
	Symbol            string `json:"symbol"`
	GroupName         string `json:"group_name"`
	GroupID           int    `json:"group_id"`
}

type OrderBook struct {
	Ask []struct {
		CurrencyPairID   int     `json:"currency_pair_id"`
		Amount           string  `json:"amount"`
		Price            string  `json:"price"`
		Amount2          string  `json:"amount2"`
		Count            int     `json:"count"`
		CumulativeAmount float64 `json:"cumulative_amount"`
	} `json:"ask"`
	Bid []struct {
		CurrencyPairID   int     `json:"currency_pair_id"`
		Amount           string  `json:"amount"`
		Price            string  `json:"price"`
		Amount2          string  `json:"amount2"`
		Count            int     `json:"count"`
		CumulativeAmount float64 `json:"cumulative_amount"`
	} `json:"bid"`
}

type WebOrderBook []struct {
	CurrencyPairID   int     `json:"currency_pair_id"`
	Amount           string  `json:"amount"`
	Price            string  `json:"price"`
	Amount2          string  `json:"amount2"`
	Count            int     `json:"count"`
	CumulativeAmount float64 `json:"cumulative_amount"`
}

type AccountBalances struct {
	Email        string            `json:"email"`
	Username     string            `json:"username"`
	Hash         string            `json:"hash"`
	IntercomHash string            `json:"intercom_hash"`
	Funds        map[string]string `json:"funds"`
	OpenOrders   int               `json:"open_orders"`
	ServerTime   int               `json:"server_time"`
}

type EmptyAccount struct {
	Success int `json:"success"`
	Data    struct {
		Email        string        `json:"email"`
		Username     string        `json:"username"`
		Hash         string        `json:"hash"`
		IntercomHash string        `json:"intercom_hash"`
		Funds        []interface{} `json:"funds"`
		OpenOrders   int           `json:"open_orders"`
		ServerTime   int           `json:"server_time"`
	} `json:"data"`
}

type WithdrawResult struct {
	ID                int         `json:"id"`
	Amount            string      `json:"amount"`
	CurrencyID        int         `json:"currency_id"`
	CurrencyCode      string      `json:"currency_code"`
	Fee               string      `json:"fee"`
	FeeCurrencyID     int         `json:"fee_currency_id"`
	FeeCurrencyCode   string      `json:"fee_currency_code"`
	Status            string      `json:"status"`
	CreatedAt         string      `json:"created_at"`
	CreatedTs         string      `json:"created_ts"`
	UpdatedAt         string      `json:"updated_at"`
	UpdatedTs         string      `json:"updated_ts"`
	Txid              interface{} `json:"txid"`
	WithdrawalAddress struct {
		Address                        string `json:"address"`
		AddressName                    string `json:"address_name"`
		AdditionalAddressParameter     string `json:"additional_address_parameter"`
		AdditionalAddressParameterName string `json:"additional_address_parameter_name"`
	} `json:"withdrawal_address"`
}

type CancelOrder struct {
	PutIntoProcessingQueue    []*PlaceOrder `json:"put_into_processing_queue"`
	NotPutIntoProcessingQueue []interface{} `json:"not_put_into_processing_queue"`
	Message                   string        `json:"message"`
}

type PlaceOrder struct {
	ID              int         `json:"id"`
	CurrencyPairID  int         `json:"currency_pair_id"`
	Price           string      `json:"price"`
	TriggerPrice    float64     `json:"trigger_price"`
	InitialAmount   string      `json:"initial_amount"`
	ProcessedAmount string      `json:"processed_amount"`
	Type            string      `json:"type"`
	OriginalType    string      `json:"original_type"`
	Created         string      `json:"created"`
	Timestamp       interface{} `json:"timestamp"`
	Status          string      `json:"status"`
}

type WalletDetails []struct {
	ID              int    `json:"id"`
	CurrencyID      int    `json:"currency_id"`
	Delisted        bool   `json:"delisted"`
	Disabled        bool   `json:"disabled"`
	DisableDeposits bool   `json:"disable_deposits"`
	CurrencyCode    string `json:"currency_code"`
	Balance         string `json:"balance"`
	FrozenBalance   string `json:"frozen_balance"`
	BonusBalance    string `json:"bonus_balance"`
}

type TradeHistory []struct {
	ID        int    `json:"id"`
	Price     string `json:"price"`
	Amount    string `json:"amount"`
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
}

type OpenOrders []struct {
	ID               int     `json:"id"`
	CurrencyPairID   int     `json:"currency_pair_id"`
	CurrencyPairName string  `json:"currency_pair_name"`
	Price            string  `json:"price"`
	TriggerPrice     float64 `json:"trigger_price"`
	InitialAmount    string  `json:"initial_amount"`
	ProcessedAmount  string  `json:"processed_amount"`
	Type             string  `json:"type"`
	OriginalType     string  `json:"original_type"`
	Created          string  `json:"created"`
	Timestamp        string  `json:"timestamp"`
	Status           string  `json:"status"`
}
