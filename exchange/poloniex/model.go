package poloniex

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type CoinsData struct {
	ID             int         `json:"id"`
	Name           string      `json:"name"`
	HumanType      string      `json:"humanType"`
	CurrencyType   string      `json:"currencyType"`
	TxFee          string      `json:"txFee"`
	MinConf        int         `json:"minConf"`
	DepositAddress interface{} `json:"depositAddress"`
	Disabled       int         `json:"disabled"`
	Delisted       int         `json:"delisted"`
	Frozen         int         `json:"frozen"`
	IsGeofenced    int         `json:"isGeofenced"`
}

type PairsData struct {
	ID            int    `json:"id"`
	Last          string `json:"last"`
	LowestAsk     string `json:"lowestAsk"`
	HighestBid    string `json:"highestBid"`
	PercentChange string `json:"percentChange"`
	BaseVolume    string `json:"baseVolume"`
	QuoteVolume   string `json:"quoteVolume"`
	IsFrozen      string `json:"isFrozen"`
	High24Hr      string `json:"high24hr"`
	Low24Hr       string `json:"low24hr"`
}

type OrderBook struct {
	Asks     [][]interface{} `json:"asks"`
	Bids     [][]interface{} `json:"bids"`
	IsFrozen string          `json:"isFrozen"`
	Seq      int64           `json:"seq"`
}

type Kline []struct {
	Date            int     `json:"date"`
	High            float64 `json:"high"`
	Low             float64 `json:"low"`
	Open            float64 `json:"open"`
	Close           float64 `json:"close"`
	Volume          float64 `json:"volume"`
	QuoteVolume     float64 `json:"quoteVolume"`
	WeightedAverage float64 `json:"weightedAverage"`
}

type Withdraw struct {
	Response string `json:"response"`
}

type PlaceOrder struct {
	OrderNumber     string `json:"orderNumber"`
	ResultingTrades []struct {
		Amount  string `json:"amount"`
		Date    string `json:"date"`
		Rate    string `json:"rate"`
		Total   string `json:"total"`
		TradeID string `json:"tradeID"`
		Type    string `json:"type"`
	} `json:"resultingTrades"`
}

type OrderStatus struct {
	Result  map[string]*OrderDetail `json:"result"`
	Success int                     `json:"success"`
}

type OrderDetail struct {
	Status         string `json:"status"`
	Rate           string `json:"rate"`
	Amount         string `json:"amount"`
	CurrencyPair   string `json:"currencyPair"`
	Date           string `json:"date"`
	Total          string `json:"total"`
	Type           string `json:"type"`
	StartingAmount string `json:"startingAmount"`
}

type CancelOrder struct {
	Success       int    `json:"success"`
	Amount        string `json:"amount"`
	ClientOrderID string `json:"clientOrderId"`
	Message       string `json:"message"`
}

type TradeHistory []struct {
	GlobalTradeID int    `json:"globalTradeID"`
	TradeID       int    `json:"tradeID"`
	Date          string `json:"date"`
	Type          string `json:"type"`
	Rate          string `json:"rate"`
	Amount        string `json:"amount"`
	Total         string `json:"total"`
	OrderNumber   int64  `json:"orderNumber"`
}

// type CancelOrder struct {
// 	Success int    `json:"success"`
// 	Amount  string `json:"amount"`
// 	Message string `json:"message"`
// }
