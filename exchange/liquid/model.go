package liquid

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type CoinsData []struct {
	CurrencyType         string      `json:"currency_type"`
	Currency             string      `json:"currency"`
	Symbol               string      `json:"symbol"`
	AssetsPrecision      int         `json:"assets_precision"`
	QuotingPrecision     int         `json:"quoting_precision"`
	MinimumWithdrawal    float64     `json:"minimum_withdrawal"`
	WithdrawalFee        float64     `json:"withdrawal_fee"`
	MinimumFee           interface{} `json:"minimum_fee"`
	MinimumOrderQuantity interface{} `json:"minimum_order_quantity"`
	DisplayPrecision     int         `json:"display_precision"`
	Depositable          bool        `json:"depositable"`
	Withdrawable         bool        `json:"withdrawable"`
	DiscountFee          float64     `json:"discount_fee"`
	CreditCardFundable   bool        `json:"credit_card_fundable,omitempty"`
	Lendable             bool        `json:"lendable"`
	PositionFundable     bool        `json:"position_fundable"`
}

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
	TakerFee            string      `json:"taker_fee"`
	MakerFee            string      `json:"maker_fee"`
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
	OrderFee             string      `json:"order_fee"`
	SourceAction         string      `json:"source_action"`
	UnwoundTradeID       interface{} `json:"unwound_trade_id"`
	TradeID              interface{} `json:"trade_id"`
}

type CancelOrder struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type Withdraw struct {
	ID            int         `json:"id"`
	Address       string      `json:"address"`
	Amount        string      `json:"amount"`
	State         string      `json:"state"`
	Currency      string      `json:"currency"`
	WithdrawalFee string      `json:"withdrawal_fee"`
	CreatedAt     int         `json:"created_at"`
	UpdatedAt     int         `json:"updated_at"`
	PaymentID     interface{} `json:"payment_id"`
}

type TradeHistory struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	Models      []struct {
		ID        int     `json:"id"`
		Quantity  float64 `json:"quantity"`
		Price     float64 `json:"price"`
		TakerSide string  `json:"taker_side"`
		CreatedAt int64   `json:"created_at"`
	} `json:"models"`
}
