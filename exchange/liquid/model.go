package liquid

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type PairsData []struct {
	ID                  string      `json:"id"`
	ProductType         string      `json:"product_type"`
	Code                string      `json:"code"`
	Name                string      `json:"name"`
	MarketAsk           float64     `json:"market_ask"`
	MarketBid           float64     `json:"market_bid"`
	Indicator           int         `json:"indicator"`
	Currency            string      `json:"currency"`
	CurrencyPairCode    string      `json:"currency_pair_code"`
	Symbol              string      `json:"symbol"`
	BtcMinimumWithdraw  interface{} `json:"btc_minimum_withdraw"`
	FiatMinimumWithdraw interface{} `json:"fiat_minimum_withdraw"`
	PusherChannel       string      `json:"pusher_channel"`
	TakerFee            float64     `json:"taker_fee"`
	MakerFee            float64     `json:"maker_fee"`
	LowMarketBid        string      `json:"low_market_bid"`
	HighMarketAsk       string      `json:"high_market_ask"`
	Volume24H           string      `json:"volume_24h"`
	LastPrice24H        interface{} `json:"last_price_24h"`
	LastTradedPrice     interface{} `json:"last_traded_price"`
	LastTradedQuantity  interface{} `json:"last_traded_quantity"`
	QuotedCurrency      string      `json:"quoted_currency"`
	BaseCurrency        string      `json:"base_currency"`
	Disabled            bool        `json:"disabled"`
	MarginEnabled       bool        `json:"margin_enabled"`
	CfdEnabled          bool        `json:"cfd_enabled"`
	LastEventTimestamp  interface{} `json:"last_event_timestamp"`
}

type OrderBook struct {
	BuyPriceLevels  [][]string `json:"buy_price_levels"`
	SellPriceLevels [][]string `json:"sell_price_levels"`
}

type AccountBalances []struct {
	Currency string `json:"currency"`
	Balance  string `json:"balance"`
}

type PlaceOrder struct {
	ID                   int         `json:"id"`
	OrderType            string      `json:"order_type"`
	Quantity             string      `json:"quantity"`
	DiscQuantity         string      `json:"disc_quantity"`
	IcebergTotalQuantity string      `json:"iceberg_total_quantity"`
	Side                 string      `json:"side"`
	FilledQuantity       string      `json:"filled_quantity"`
	Price                float64     `json:"price"`
	CreatedAt            int         `json:"created_at"`
	UpdatedAt            int         `json:"updated_at"`
	Status               string      `json:"status"`
	LeverageLevel        int         `json:"leverage_level"`
	SourceExchange       string      `json:"source_exchange"`
	ProductID            int         `json:"product_id"`
	ProductCode          string      `json:"product_code"`
	FundingCurrency      string      `json:"funding_currency"`
	CryptoAccountID      interface{} `json:"crypto_account_id"`
	CurrencyPairCode     string      `json:"currency_pair_code"`
	AveragePrice         float64     `json:"average_price"`
	Target               string      `json:"target"`
	OrderFee             float64     `json:"order_fee"`
	SourceAction         string      `json:"source_action"`
	UnwoundTradeID       interface{} `json:"unwound_trade_id"`
	TradeID              interface{} `json:"trade_id"`
}

type OrderStatus struct {
	ID                   int         `json:"id"`
	OrderType            string      `json:"order_type"`
	Quantity             string      `json:"quantity"`
	DiscQuantity         string      `json:"disc_quantity"`
	IcebergTotalQuantity string      `json:"iceberg_total_quantity"`
	Side                 string      `json:"side"`
	FilledQuantity       string      `json:"filled_quantity"`
	Price                float64     `json:"price"`
	CreatedAt            int         `json:"created_at"`
	UpdatedAt            int         `json:"updated_at"`
	Status               string      `json:"status"`
	LeverageLevel        int         `json:"leverage_level"`
	SourceExchange       string      `json:"source_exchange"`
	ProductID            int         `json:"product_id"`
	ProductCode          string      `json:"product_code"`
	FundingCurrency      string      `json:"funding_currency"`
	CryptoAccountID      interface{} `json:"crypto_account_id"`
	CurrencyPairCode     string      `json:"currency_pair_code"`
	AveragePrice         float64     `json:"average_price"`
	Target               string      `json:"target"`
	OrderFee             float64     `json:"order_fee"`
	SourceAction         string      `json:"source_action"`
	UnwoundTradeID       interface{} `json:"unwound_trade_id"`
	TradeID              interface{} `json:"trade_id"`
}

type CancelOrder struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}
