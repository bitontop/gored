package exchange

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/pair"
)

type Update struct {
	ExNames []ExchangeName
	Method  UpdateMethod
	Time    time.Duration
}

type Config struct {
	ExName        ExchangeName
	Source        DataSource
	SourceURI     string
	Account_ID    string
	API_KEY       string
	API_SECRET    string
	Two_Factor    string
	Passphrase    string //Memo for bitmart
	TradePassword string
	UserID        string
}

type PairConstraint struct {
	PairID      int
	Pair        *pair.Pair //the code on excahnge with the same chain, eg: BCH, BCC on different exchange, but they are the same chain
	ExID        string
	ExSymbol    string
	MakerFee    float64
	TakerFee    float64
	LotSize     float64 // the decimal place for this coin on exchange for the pairs, eg:  BTC: 0.00001    NEO:1   LTC: 0.001 ETH:0.01
	PriceFilter float64
	Listed      bool
	Issue       string //the issue for the pair if have any problem
}

type CoinConstraint struct {
	CoinID       int
	Coin         *coin.Coin
	ExSymbol     string
	ChainType    ChainType
	TxFee        float64 // the withdraw fee for this exchange
	Withdraw     bool
	Deposit      bool
	Confirmation int
	Listed       bool
	Issue        string //the issue for the chain if have any problem
}

type ConstrainFetchMethod struct {
	PublicAPI   bool
	PrivateAPI  bool
	HealthAPI   bool // get exchange health status from exchange's API directly
	HasWithdraw bool // has withdraw method implemented

	Fee          bool // true only when get Fee from API directly
	LotSize      bool
	PriceFilter  bool
	TxFee        bool
	Withdraw     bool
	Deposit      bool
	Confirmation bool
}

type OrderStatus string

const (
	New       OrderStatus = "New"
	Filled    OrderStatus = "Filled"
	Partial   OrderStatus = "Partial"
	Canceling OrderStatus = "Canceling"
	Canceled  OrderStatus = "Canceled"
	Rejected  OrderStatus = "Rejected"
	Expired   OrderStatus = "Expired"
	Other     OrderStatus = "Other"
)

type Order struct {
	Pair          *pair.Pair
	OrderID       string
	FilledOrders  []int64
	Rate          float64 `bson:"Rate"`
	Quantity      float64 `bson:"Quantity"`
	Side          string
	Status        OrderStatus `json:"status"`
	StatusMessage string
	DealRate      float64
	DealQuantity  float64
	JsonResponse  string

	Canceled     bool
	CancelStatus string
}

type Maker struct {
	WorkerIP        string     `bson:"workerip"`
	WorkerDeadTS    float64    `bson:"workerdeadts"`
	Source          DataSource `bson:"source"`
	BeforeTimestamp float64    `bson:"beforetimestamp"`
	AfterTimestamp  float64    `bson:"aftertimestamp"`
	Timestamp       float64    `bson:"timestamp"`
	Nounce          int        `bson:"nounce"`
	LastUpdateID    int64      `json:"lastUpdateId"`
	Bids            []Order    `json:"bids"`
	Asks            []Order    `json:"asks"`
}

type Margin struct {
	Action   MarginAction
	Pair     *pair.Pair
	Currency *coin.Coin
	Quantity float64
	Order    *MarginOrder
	Balance  *MarginBalance
}

type MarginOrder struct {
	LoanBalance     string `json:"loan-balance"`
	InterestBalance string `json:"interest-balance"`
	InterestRate    string `json:"interest-rate"`
	LoanAmount      string `json:"loan-amount"`
	AccruedAt       int64  `json:"accrued-at"`
	InterestAmount  string `json:"interest-amount"`
	Symbol          string `json:"symbol"`
	Currency        string `json:"currency"`
	ID              int    `json:"id"`
	State           string `json:"state"`
	AccountID       int    `json:"account-id"`
	UserID          int    `json:"user-id"`
	CreatedAt       int64  `json:"created-at"`
}

type MarginBalance struct {
	ID       int    `json:"id"`
	Type     string `json:"type"`
	State    string `json:"state"`
	Symbol   string `json:"symbol"`
	FlPrice  string `json:"fl-price"`
	FlType   string `json:"fl-type"`
	RiskRate string `json:"risk-rate"`
	List     []struct {
		Currency string `json:"currency"`
		Type     string `json:"type"`
		Balance  string `json:"balance"`
	} `json:"list"`
}
