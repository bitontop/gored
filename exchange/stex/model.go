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
	Notice  string          `json:"notice"`
	Message string          `json:"message"`
	Err     interface{}     `json:"errors"`
}

type JsonResponseV3 struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
	Err     interface{}     `json:"errors"`
}

type CoinsData []struct {
	ID                        int    `json:"id"`
	Code                      string `json:"code"`
	Name                      string `json:"name"`
	Active                    bool   `json:"active"`
	Delisted                  bool   `json:"delisted"`
	Precision                 int    `json:"precision"`
	MinimumWithdrawalAmount   string `json:"minimum_withdrawal_amount"`
	MinimumDepositAmount      string `json:"minimum_deposit_amount"`
	DepositFeeCurrencyID      int    `json:"deposit_fee_currency_id"`
	DepositFeeCurrencyCode    string `json:"deposit_fee_currency_code"`
	DepositFeeConst           string `json:"deposit_fee_const"`
	DepositFeePercent         string `json:"deposit_fee_percent"`
	WithdrawalFeeCurrencyID   int    `json:"withdrawal_fee_currency_id"`
	WithdrawalFeeCurrencyCode string `json:"withdrawal_fee_currency_code"`
	WithdrawalFeeConst        string `json:"withdrawal_fee_const"`
	WithdrawalFeePercent      string `json:"withdrawal_fee_percent"`
	BlockExplorerURL          string `json:"block_explorer_url"`
}

type PairsData []struct {
	ID                int         `json:"id"`
	CurrencyID        int         `json:"currency_id"`
	CurrencyCode      string      `json:"currency_code"`
	CurrencyName      string      `json:"currency_name"`
	MarketCurrencyID  int         `json:"market_currency_id"`
	MarketCode        string      `json:"market_code"`
	MarketName        string      `json:"market_name"`
	MinOrderAmount    string      `json:"min_order_amount"`
	MinBuyPrice       string      `json:"min_buy_price"`
	MinSellPrice      string      `json:"min_sell_price"`
	BuyFeePercent     string      `json:"buy_fee_percent"`
	SellFeePercent    string      `json:"sell_fee_percent"`
	Active            bool        `json:"active"`
	Delisted          bool        `json:"delisted"`
	PairMessage       interface{} `json:"pair_message"`
	CurrencyPrecision int         `json:"currency_precision"`
	MarketPrecision   int         `json:"market_precision"`
	Symbol            string      `json:"symbol"`
	GroupName         string      `json:"group_name"`
	GroupID           int         `json:"group_id"`
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

type UserInfo struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	// UserSessions []struct {
	// 	IP        string `json:"ip"`
	// 	Date      string `json:"date"`
	// 	CreatedAt string `json:"created_at"`
	// 	Active    bool   `json:"active"`
	// } `json:"userSessions"`
	Hash         string            `json:"hash"`
	IntercomHash string            `json:"intercom_hash"`
	Funds        map[string]string `json:"funds"`
	// HoldFunds        map[string]string `json:"hold_funds"`
	// WalletsAddresses map[string]string `json:"wallets_addresses"`
	// PublickKey       map[string]string `json:"publick_key"`
	// AssetsPortfolio  struct {
	// 	PortfolioPrice       int `json:"portfolio_price"`
	// 	FrozenPortfolioPrice int `json:"frozen_portfolio_price"`
	// 	Count                int `json:"count"`
	// 	Assets               []struct {
	// 		WalletsAdress string `json:"wallets_adress"`
	// 		PublickKey    string `json:"publick_key"`
	// 		Funds         string `json:"funds"`
	// 	} `json:"assets"`
	// } `json:"Assets portfolio"`
	OpenOrders int `json:"open_orders"`
	ServerTime int `json:"server_time"`
}

type EmptyAccount struct {
	Success int `json:"success"`
	Data    struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		//UserSessions []interface{} `json:"userSessions"`
		Hash         string        `json:"hash"`
		IntercomHash string        `json:"intercom_hash"`
		Funds        []interface{} `json:"funds"`
		// HoldFunds        []interface{} `json:"hold_funds"`
		// WalletsAddresses []interface{} `json:"wallets_addresses"`
		// AssetsPortfolio struct {
		// 	PortfolioPrice       int           `json:"portfolio_price"`
		// 	FrozenPortfolioPrice int           `json:"frozen_portfolio_price"`
		// 	Count                int           `json:"count"`
		// 	Assets               []interface{} `json:"assets"`
		// } `json:"Assets portfolio"`
		OpenOrders int `json:"open_orders"`
		ServerTime int `json:"server_time"`
	} `json:"data"`
}

type Withdraw struct {
	Code                  string `json:"code"`
	ID                    int    `json:"id"`
	Amount                string `json:"amount"`
	Address               string `json:"address"`
	WithdrawalFee         string `json:"withdrawal_fee"`
	WithdrawalFeeCurrency string `json:"withdrawal_fee_currency"`
	Token                 string `json:"token"`
	Date                  struct {
		Date         string `json:"date"`
		TimezoneType int    `json:"timezone_type"`
		Timezone     string `json:"timezone"`
	} `json:"date"`
	Msg string `json:"msg"`
}

type CancelOrder struct {
	Funds   map[string]string `json:"funds"`
	OrderID string            `json:"order_id"`
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

type TradeDetail struct {
	Funds   map[string]string `json:"funds"`
	OrderID int64             `json:"order_id"`
}

type ActiveOrder map[string]*OrderDetail

type OrderDetail struct {
	Pair           string                    `json:"pair"`
	Type           string                    `json:"type"`
	OriginalAmount string                    `json:"original_amount"`
	BuyAmount      interface{}               `json:"buy_amount"`
	SellAmount     interface{}               `json:"sell_amount"`
	IsYourOrder    int                       `json:"is_your_order"`
	Timestamp      int                       `json:"timestamp"`
	Rates          map[string]*BuySellAmount `json:"rates"`
}

type BuySellAmount struct {
	BuyAmount  interface{} `json:"buy_amount"`
	SellAmount interface{} `json:"sell_amount"`
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

/* type Uuid struct {
	Id string `json:"uuid"`
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

/* type JsonResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
} */
