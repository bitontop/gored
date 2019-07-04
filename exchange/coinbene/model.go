package coinbene

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type PairsData struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Timestamp   int64  `json:"timestamp"`
	Symbol      []struct {
		Ticker      string `json:"ticker"`
		BaseAsset   string `json:"baseAsset"`
		QuoteAsset  string `json:"quoteAsset"`
		TakerFee    string `json:"takerFee"`
		MakerFee    string `json:"makerFee"`
		TickSize    string `json:"tickSize"`
		LotStepSize string `json:"lotStepSize"`
		MinQuantity string `json:"minQuantity"`
	} `json:"symbol"`
}

type OrderBook struct {
	Orderbook struct {
		Asks OrderBookDetail `json:"asks"`
		Bids OrderBookDetail `json:"bids"`
	} `json:"orderbook"`
	Status      string `json:"status"`
	Description string `json:"description"`
	Symbol      string `json:"symbol"`
	Timestamp   int64  `json:"timestamp"`
}

type OrderBookDetail []struct {
	Quantity float64 `json:"quantity"`
	Price    float64 `json:"price"`
}

type AccountBalances struct {
	Account string `json:"account"`
	Balance []struct {
		Asset     string `json:"asset"`
		Available string `json:"available"`
		Reserved  string `json:"reserved"`
		Total     string `json:"total"`
	} `json:"balance"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Timestamp   int64  `json:"timestamp"`
}

type PlaceOrder struct {
	Orderid     string `json:"orderid"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Timestamp   int64  `json:"timestamp"`
}

type OrderStatus struct {
	Order struct {
		Createtime     int64  `json:"createtime"`
		Filledamount   string `json:"filledamount"`
		Filledquantity string `json:"filledquantity"`
		Orderid        string `json:"orderid"`
		Orderquantity  string `json:"orderquantity"`
		Orderstatus    string `json:"orderstatus"`
		Price          string `json:"price"`
		Symbol         string `json:"symbol"`
		Type           string `json:"type"`
	} `json:"order"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Timestamp   int64  `json:"timestamp"`
}

type Withdraw struct {
	Status     string `json:"status"`
	Timestamp  int64  `json:"timestamp"`
	WithdrawID int    `json:"withdrawId"`
}
