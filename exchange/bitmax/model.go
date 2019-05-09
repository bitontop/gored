package bitmax

import "encoding/json"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type JsonResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Status  string          `json:"status"`
	Email   string          `json:"email"`
	Data    json.RawMessage `json:"data"`
}

type BitmaxCoins []struct {
	AssetCode        string  `json:"assetCode"`
	AssetName        string  `json:"assetName"`
	PrecisionScale   int     `json:"precisionScale"`
	NativeScale      int     `json:"nativeScale"`
	WithdrawalFee    float64 `json:"withdrawalFee"`
	MinWithdrawalAmt float64 `json:"minWithdrawalAmt"`
	Status           string  `json:"status"`
}

type BitmaxPairs []struct {
	Symbol        string `json:"symbol"`
	BaseAsset     string `json:"baseAsset"`
	QuoteAsset    string `json:"quoteAsset"`
	PriceScale    int    `json:"priceScale"`
	QtyScale      int    `json:"qtyScale"`
	NotionalScale int    `json:"notionalScale"`
	MinQty        string `json:"minQty"`
	MaxQty        string `json:"maxQty"`
	MinNotional   string `json:"minNotional"`
	MaxNotional   string `json:"maxNotional"`
	Status        string `json:"status"`
	MiningStatus  string `json:"miningStatus"`
}

type BitmaxOrderBook struct {
	M    string     `json:"m"`
	S    string     `json:"s"`
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
}

type BitmaxAccountGroup struct {
	AccountGroup int `json:"accountGroup"`
}

type BitmaxBalance []struct {
	AssetCode       string `json:"assetCode"`
	AssetName       string `json:"assetName"`
	TotalAmount     string `json:"totalAmount"`
	AvailableAmount string `json:"availableAmount"`
	InOrderAmount   string `json:"inOrderAmount"`
}

type Withdrawal struct {
	Msg    string `json:"msg"`
	Status string `json:"status"`
}

type BitmaxWithdraw struct {
	RequestID string `json:"requestId"`
	Time      int64  `json:"time"`
	AssetCode string `json:"assetCode"`
	Amount    string `json:"amount"`
	Address   struct {
		Address string `json:"address"`
	} `json:"address"`
}

type BitmaxOrder struct {
	Coid    string `json:"coid"`
	Action  string `json:"action"`
	Success bool   `json:"success"`
}

type BitmaxOrderStatus struct {
	Time       int64  `json:"time"`
	Coid       string `json:"coid"`
	Symbol     string `json:"symbol"`
	BaseAsset  string `json:"baseAsset"`
	QuoteAsset string `json:"quoteAsset"`
	Side       string `json:"side"`
	OrderPrice string `json:"orderPrice"`
	StopPrice  string `json:"stopPrice"`
	OrderQty   string `json:"orderQty"`
	Filled     string `json:"filled"`
	Fee        string `json:"fee"`
	FeeAsset   string `json:"feeAsset"`
	Status     string `json:"status"`
}

type BitmaxCancel struct {
	Coid    string `json:"coid"`
	Action  string `json:"action"`
	Success bool   `json:"success"`
}

/* type PlaceOrder struct {
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
} */

/* type Uuid struct {
	Id string `json:"uuid"`
} */

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
