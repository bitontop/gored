package huobi

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

type JsonResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Status  string          `json:"status"`
	Data    json.RawMessage `json:"data"`
	Tick    json.RawMessage `json:"tick"`
	ErrCode string          `json:"err-code"`
	ErrMsg  string          `json:"err-msg"`
}

type CoinsData []struct {
	Currency string `json:"currency"`
	Chains   []struct {
		Chain                  string `json:"chain"`
		NumOfConfirmations     int    `json:"numOfConfirmations"`
		NumOfFastConfirmations int    `json:"numOfFastConfirmations"`
		DepositStatus          string `json:"depositStatus"`
		MinDepositAmt          string `json:"minDepositAmt"`
		WithdrawStatus         string `json:"withdrawStatus"`
		MinWithdrawAmt         string `json:"minWithdrawAmt"`
		WithdrawPrecision      int    `json:"withdrawPrecision"`
		MaxWithdrawAmt         string `json:"maxWithdrawAmt"`
		WithdrawQuotaPerDay    string `json:"withdrawQuotaPerDay"`
		WithdrawQuotaPerYear   string `json:"withdrawQuotaPerYear"`
		WithdrawQuotaTotal     string `json:"withdrawQuotaTotal"`
		WithdrawFeeType        string `json:"withdrawFeeType"`
		TransactFeeWithdraw    string `json:"transactFeeWithdraw"`
		MinTransactFeeWithdraw string `json:"minTransactFeeWithdraw"`
	} `json:"chains"`
	InstStatus string `json:"instStatus"`
}

// doesn't work anymore
// type CoinsData []struct {
// 	CurrencyAddrWithTag     bool          `json:"currency-addr-with-tag"`
// 	FastConfirms            int           `json:"fast-confirms"`
// 	SafeConfirms            int           `json:"safe-confirms"`
// 	CurrencyType            string        `json:"currency-type"`
// 	SupportSites            []interface{} `json:"support-sites"`
// 	OtcEnable               int           `json:"otc-enable"`
// 	CountryDisabled         bool          `json:"country-disabled"`
// 	Tags                    interface{}   `json:"tags"`
// 	DepositEnabled          bool          `json:"deposit-enabled"`
// 	WithdrawEnabled         bool          `json:"withdraw-enabled"`
// 	WhiteEnabled            bool          `json:"white-enabled"`
// 	WithdrawPrecision       int           `json:"withdraw-precision"`
// 	CurrencyPartition       string        `json:"currency-partition"`
// 	QuoteCurrency           bool          `json:"quote-currency"`
// 	WithdrawMinAmount       string        `json:"withdraw-min-amount"`
// 	Weight                  int           `json:"weight"`
// 	Visible                 bool          `json:"visible"`
// 	ShowPrecision           string        `json:"show-precision"`
// 	DepositMinAmount        string        `json:"deposit-min-amount"`
// 	VisibleAssetsTimestamp  int64         `json:"visible-assets-timestamp"`
// 	DepositEnableTimestamp  int64         `json:"deposit-enable-timestamp"`
// 	WithdrawEnableTimestamp int64         `json:"withdraw-enable-timestamp"`
// 	Name                    string        `json:"name"`
// 	State                   string        `json:"state"`
// 	DisplayName             string        `json:"display-name"`
// 	DepositDesc             string        `json:"deposit-desc"`
// 	WithdrawDesc            string        `json:"withdraw-desc"`
// 	SuspendVisibleDesc      string        `json:"suspend-visible-desc"`
// 	SuspendDepositDesc      string        `json:"suspend-deposit-desc"`
// 	SuspendWithdrawDesc     string        `json:"suspend-withdraw-desc"`
// 	CurrencyAddrOneoff      bool          `json:"currency-addr-oneoff,omitempty"`
// 	Blockchains             string        `json:"blockchains,omitempty"`
// }

type PairsData []struct {
	BaseCurrency             string  `json:"base-currency"`
	QuoteCurrency            string  `json:"quote-currency"`
	PricePrecision           int     `json:"price-precision"`
	AmountPrecision          int     `json:"amount-precision"`
	SymbolPartition          string  `json:"symbol-partition"`
	Symbol                   string  `json:"symbol"`
	State                    string  `json:"state"`
	ValuePrecision           int     `json:"value-precision"`
	MinOrderAmt              float64 `json:"min-order-amt"`
	MaxOrderAmt              float64 `json:"max-order-amt"`
	MinOrderValue            float64 `json:"min-order-value"`
	LeverageRatio            float64 `json:"leverage-ratio,omitempty"`
	SuperMarginLeverageRatio float64 `json:"super-margin-leverage-ratio,omitempty"`
	FundingLeverageRatio     float64 `json:"funding-leverage-ratio,omitempty"`
}

type OrderBook struct {
	Bids    [][]float64 `json:"bids"`
	Asks    [][]float64 `json:"asks"`
	ID      interface{} `json:"id"`
	Ts      int64       `json:"ts"`
	Version int64       `json:"version"`
}

type KLines []struct {
	ID     int64   `json:"id"`
	Amount float64 `json:"amount"`
	Count  int     `json:"count"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	Low    float64 `json:"low"`
	High   float64 `json:"high"`
	Vol    float64 `json:"vol"`
}

type AccountsReturn []struct {
	ID      int64  `json:"id"`
	Type    string `json:"type"`
	State   string `json:"state"`
	SubType string `json:"sub-type"`
}

type AccountBalances struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	State string `json:"state"`
	List  []struct {
		Currency string `json:"currency"`
		Type     string `json:"type"`
		Balance  string `json:"balance"`
		Address  string `json:"address"`
	} `json:"list"`
}

type OrderStatus struct {
	ID               int    `json:"id"`
	Symbol           string `json:"symbol"`
	AccountID        int    `json:"account-id"`
	Amount           string `json:"amount"`
	Price            string `json:"price"`
	CreatedAt        int64  `json:"created-at"`
	Type             string `json:"type"`
	FieldAmount      string `json:"field-amount"`
	FieldCashAmount  string `json:"field-cash-amount"`
	FieldFees        string `json:"field-fees"`
	FilledAmount     string `json:"filled-amount"`
	FilledCashAmount string `json:"filled-cash-amount"`
	FilledFees       string `json:"filled-fees"`
	FinishedAt       int64  `json:"finished-at"`
	UserID           int    `json:"user-id"`
	Source           string `json:"source"`
	State            string `json:"state"`
	CanceledAt       int    `json:"canceled-at"`
	Exchange         string `json:"exchange"`
	Batch            string `json:"batch"`
}

type TradeHistory struct {
	Status string `json:"status"`
	Ch     string `json:"ch"`
	Ts     uint   `json:"ts"`
	Data   []struct {
		ID   uint `json:"id"`
		Ts   uint `json:"ts"`
		Data []struct {
			Amount  float64 `json:"amount"`
			TradeID uint    `json:"trade-id"`
			Ts      int64   `json:"ts"`
			// ID        uint64  `json:"id"`
			Price     float64 `json:"price"`
			Direction string  `json:"direction"`
		} `json:"data"`
	} `json:"data"`
}

type DWHistory struct {
	ID         int     `json:"id"`
	Type       string  `json:"type"`
	Currency   string  `json:"currency"`
	TxHash     string  `json:"tx-hash"`
	Amount     float64 `json:"amount"`
	Address    string  `json:"address"`
	AddressTag string  `json:"address-tag"`
	Fee        int     `json:"fee"`
	State      string  `json:"state"`
	CreatedAt  int64   `json:"created-at"`
	UpdatedAt  int64   `json:"updated-at"`
}

type DepositAddress struct {
	Currency   string `json:"currency"`
	Address    string `json:"address"`
	AddressTag string `json:"addressTag"`
	Chain      string `json:"chain"`
}

type SubAccountBalances []struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
	List []struct {
		Currency string `json:"currency"`
		Type     string `json:"type"`
		Balance  string `json:"balance"`
	} `json:"list"`
}

type SubAccountList []struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Subtype string `json:"subtype"` // sub account type
	State   string `json:"state"`
}

type SubAllAccountBalances []struct {
	Currency string `json:"currency"`
	Type     string `json:"type"`
	Balance  string `json:"balance"`
}

type TickerPrice []struct {
	Symbol  string  `json:"symbol"`
	Open    float64 `json:"open"`
	High    float64 `json:"high"`
	Low     float64 `json:"low"`
	Close   float64 `json:"close"`
	Amount  float64 `json:"amount"`
	Vol     float64 `json:"vol"`
	Count   int     `json:"count"`
	Bid     float64 `json:"bid"`
	BidSize float64 `json:"bidSize"`
	Ask     float64 `json:"ask"`
	AskSize float64 `json:"askSize"`
}

type TransferHistory []struct {
	AccountID    int     `json:"accountId"`
	Currency     string  `json:"currency"`
	TransactAmt  float64 `json:"transactAmt"`
	TransactType string  `json:"transactType"`
	TransferType string  `json:"transferType"`
	TransactID   int     `json:"transactId"`
	TransactTime int64   `json:"transactTime"`
	Transferer   int     `json:"transferer"`
	Transferee   int     `json:"transferee"`
}

type SubTransfer struct {
	Status  string `json:"status"`
	Data    int    `json:"data"`
	ErrCode string `json:"err-code"`
	ErrMsg  string `json:"err-msg"`
}
