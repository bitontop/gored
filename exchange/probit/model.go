package probit

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type CoinsData struct {
	Data []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		DisplayName struct {
			KoKr string `json:"ko-kr"`
		} `json:"display_name"`
		Platform             string `json:"platform"`
		Precision            int    `json:"precision"`
		DisplayPrecision     int    `json:"display_precision"`
		MinConfirmationCount int    `json:"min_confirmation_count"`
		MinWithdrawalAmount  string `json:"min_withdrawal_amount"`
		WithdrawalFee        string `json:"withdrawal_fee"`
		DepositSuspended     bool   `json:"deposit_suspended"`
		WithdrawalSuspended  bool   `json:"withdrawal_suspended"`
		InternalPrecision    int    `json:"internal_precision"`
		ShowInUI             bool   `json:"show_in_ui"`
		SuspendedReason      string `json:"suspended_reason"`
		MinDepositAmount     string `json:"min_deposit_amount"`
	} `json:"data"`
}

type PairsData struct {
	Data []struct {
		ID                string `json:"id"`
		BaseCurrencyID    string `json:"base_currency_id"`
		QuoteCurrencyID   string `json:"quote_currency_id"`
		MinPrice          string `json:"min_price"`
		MaxPrice          string `json:"max_price"`
		PriceIncrement    string `json:"price_increment"`
		MinQuantity       string `json:"min_quantity"`
		MaxQuantity       string `json:"max_quantity"`
		QuantityPrecision int    `json:"quantity_precision"`
		MinCost           string `json:"min_cost"`
		MaxCost           string `json:"max_cost"`
		CostPrecision     int    `json:"cost_precision"`
		TakerFeeRate      string `json:"taker_fee_rate"`
		MakerFeeRate      string `json:"maker_fee_rate"`
		ShowInUI          bool   `json:"show_in_ui"`
		Closed            bool   `json:"closed"`
	} `json:"data"`
}

type OrderBook struct {
	Data []struct {
		Side     string `json:"side"`
		Price    string `json:"price"`
		Quantity string `json:"quantity"`
	} `json:"data"`
}

// Old

type AccountBalances struct {
	MakerCommission  int  `json:"makerCommission"`
	TakerCommission  int  `json:"takerCommission"`
	BuyerCommission  int  `json:"buyerCommission"`
	SellerCommission int  `json:"sellerCommission"`
	CanTrade         bool `json:"canTrade"`
	CanWithdraw      bool `json:"canWithdraw"`
	CanDeposit       bool `json:"canDeposit"`
	Balances         []struct {
		Asset  string `json:"asset"`
		Free   string `json:"free"`
		Locked string `json:"locked"`
	} `json:"balances"`
}

type WithdrawResponse struct {
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

type PlaceOrder struct {
	Symbol        string `json:"symbol"`
	OrderID       int    `json:"orderId"`
	ClientOrderID string `json:"clientOrderId"`
	TransactTime  int64  `json:"transactTime"`
	Price         string `json:"price"`
	OrigQty       string `json:"origQty"`
	ExecutedQty   string `json:"executedQty"`
	Status        string `json:"status"`
	TimeInForce   string `json:"timeInForce"`
	Type          string `json:"type"`
	Side          string `json:"side"`
	StopPrice     string `json:"stopPrice"`
	IcebergQty    string `json:"icebergQty"`
	Time          int64  `json:"time"`
	Code          int    `json:"code"`
	Msg           string `json:"msg"`
}
