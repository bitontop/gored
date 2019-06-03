package gemini

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type ErrorResponse struct {
	// error return
	Result  string `json:"result"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type PairsFile []struct {
	Symbol        string  `json:"symbol"`
	Base          string  `json:"base"`
	Quote         string  `json:"quote"`
	MinOrderSize  float64 `json:"min_order_size"`
	MinOrderIncre float64 `json:"min_order_incre"`
	MinPriceIncre float64 `json:"min_price_incre"`
}

type OrderBook struct {
	Bids []struct {
		Price     string `json:"price"`
		Amount    string `json:"amount"`
		Timestamp string `json:"timestamp"`
	} `json:"bids"`
	Asks []struct {
		Price     string `json:"price"`
		Amount    string `json:"amount"`
		Timestamp string `json:"timestamp"`
	} `json:"asks"`
}

type AccountBalances []struct {
	Type                   string `json:"type"`
	Currency               string `json:"currency"`
	Amount                 string `json:"amount"`
	Available              string `json:"available"`
	AvailableForWithdrawal string `json:"availableForWithdrawal"`
}

/* type AccountBalances []struct {
	Currency      string  `json:"Currency"`
	Balance       float64 `json:"Balance"`
	Available     float64 `json:"Available"`
	Pending       float64 `json:"Pending"`
	CryptoAddress string  `json:"CryptoAddress"`
	Requested     bool    `json:"Requested"`
	Uuid          string  `json:"Uuid"`
} */

/* type OrderBook struct {
	Buy  []exchange.Order `json:"buy"`
	Sell []exchange.Order `json:"sell"`
} */

/* type PairsData []struct {
	MarketCurrency     string      `json:"MarketCurrency"`
	BaseCurrency       string      `json:"BaseCurrency"`
	MarketCurrencyLong string      `json:"MarketCurrencyLong"`
	BaseCurrencyLong   string      `json:"BaseCurrencyLong"`
	MinTradeSize       float64     `json:"MinTradeSize"`
	MarketName         string      `json:"MarketName"`
	IsActive           bool        `json:"IsActive"`
	Created            string      `json:"Created"`
	Notice             interface{} `json:"Notice"`
	IsSponsored        interface{} `json:"IsSponsored"`
	LogoURL            string      `json:"LogoUrl"`
} */

type Uuid struct {
	Id string `json:"uuid"`
}

type PlaceOrder struct {
	AccountId                  string
	OrderUuid                  string `json:"OrderUuid"`
	Exchange                   string `json:"Exchange"`
	Type                       string
	Quantity                   float64 `json:"Quantity"`
	QuantityRemaining          float64 `json:"QuantityRemaining"`
	Limit                      float64 `json:"Limit"`
	Reserved                   float64
	ReserveRemaining           float64
	CommissionReserved         float64
	CommissionReserveRemaining float64
	CommissionPaid             float64
	Price                      float64 `json:"Price"`
	PricePerUnit               float64 `json:"PricePerUnit"`
	Opened                     string
	Closed                     string
	IsOpen                     bool
	Sentinel                   string
	CancelInitiated            bool
	ImmediateOrCancel          bool
	IsConditional              bool
	Condition                  string
	ConditionTarget            string
}

/* type CoinsData []struct {
	Currency        string      `json:"Currency"`
	CurrencyLong    string      `json:"CurrencyLong"`
	MinConfirmation int         `json:"MinConfirmation"`
	TxFee           float64     `json:"TxFee"`
	IsActive        bool        `json:"IsActive"`
	IsRestricted    bool        `json:"IsRestricted"`
	CoinType        string      `json:"CoinType"`
	BaseAddress     string      `json:"BaseAddress"`
	Notice          interface{} `json:"Notice"`
} */
