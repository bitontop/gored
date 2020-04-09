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
	ExpireTS      int64
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
	HasTransfer bool // has transfer method

	Fee             bool // true only when get Fee from API directly
	LotSize         bool
	PriceFilter     bool
	TxFee           bool
	Withdraw        bool
	Deposit         bool
	Confirmation    bool
	ConstrainSource int // 1）API   2)WEB 3）Manual
	ApiRestrictIP   bool
}

type OrderStatus string

const (
	New       OrderStatus = "New"
	Filled    OrderStatus = "Filled"
	Partial   OrderStatus = "Partial"
	Canceling OrderStatus = "Canceling"
	Cancelled OrderStatus = "Cancelled"
	Rejected  OrderStatus = "Rejected"
	Expired   OrderStatus = "Expired"
	Other     OrderStatus = "Other"
)

type Order struct {
	EX            ExchangeName
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
	Timestamp     int64 `bson:"timestamp"`
	JsonResponse  string

	Canceled     bool
	CancelStatus string
	Error        error
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
	Action        MarginAction
	Pair          *pair.Pair
	Currency      *coin.Coin
	Rate          float64
	Quantity      float64
	TransferID    int
	Order         *Order
	MarginOrder   *MarginOrder
	MarginBalance *MarginBalance
}

type MarginOrder struct {
	ID              int
	LoanAmount      float64
	LoanBalance     float64
	InterestRate    float64
	InterestAmount  float64
	InterestBalance float64
	State           string
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

//Account Operation Data Modeling

type OperationType string

const (
	Withdraw    OperationType = "Withdraw"
	Transfer    OperationType = "Transfer"   // transfer  between inneral wallet
	Balance     OperationType = "Balance"    // balance(s) of different accounts
	BalanceList OperationType = "BalanceAll" // balance(s) of different accounts

	//Public Query
	GetCoin OperationType = "GetCoin"
	GetPair OperationType = "GetPair"

	TradeHistory  OperationType = "TradeHistory"
	Orderbook     OperationType = "Orderbook"
	CoinChainType OperationType = "CoinChainType"

	//Trade (Private Action)
	PlaceOrder     OperationType = "PlaceOrder"
	CancelOrder    OperationType = "CancelOrder"
	GetOrderStatus OperationType = "GetOrderStatus"

	//User (Private Action)
	GetOpenOrder         OperationType = "GetOpenOrder"    // New and Partial Orders
	GetOrderHistory      OperationType = "GetOrderHistory" // All Orders other than open orders
	GetDepositHistory    OperationType = "GetDepositHistory"
	GetWithdrawalHistory OperationType = "GetWithdrawalHistory"
	GetDepositAddress    OperationType = "GetDepositAddress" // Get address for one coin
)

type WalletType string

const (
	AssetWallet   WalletType = "AssetWallet"
	SpotWallet    WalletType = "SpotWallet"
	FiatOTCWallet WalletType = "FiatOTCWallet"
	MarginWallet  WalletType = "MarginWallet"
	// ##### New - Contract
	ContractWallet WalletType = "ContractWallet"
)

type AccountOperation struct {
	ID int `json:"id"` //dummy at this moment for

	Type OperationType `json:"type"`
	Ex   ExchangeName  `json:"exchange_name"`

	//#Transfer,Balance,Withdraw
	Coin *coin.Coin `json:"transfer_coin"` //BOT standard symbol, not the symbol on exchange

	//specific operations
	// #Transfer
	TransferFrom        WalletType `json:"transfer_from"`
	TransferDestination WalletType `json:"transfer_dest"`
	TransferAmount      string     `json:"transfer_amount"`

	// #Withdraw
	WithdrawAddress string `json:"withdraw_address"`
	WithdrawTag     string `json:"withdraw_tag"`
	WithdrawAmount  string `json:"withdraw_amount"` //here using string instead of float64
	WithdrawID      string `json:"withdraw_id"`

	// #Balance
	Wallet WalletType `json:"wallet"` // Contract/Spot operation. Default spot if empty

	//#Single Balance
	BalanceAvailable float64 `json:"balance_available"` //the fund able to do trading
	BalanceFrozen    float64 `json:"balance_frozen"`    // the fund in order or frozen can't do trading         the total amount of fund should be   BalanceAvailable + BalanceFrozen

	//#Balance Listed
	//Coin = nil
	BalanceList []AssetBalance `json:"balance_list"`

	// #OpenOrder, OrderHistory
	OpenOrders   []*Order
	OrderHistory []*Order

	// #GetWithdrawal/DepositHistory
	WithdrawalHistory []*WDHistory
	DepositHistory    []*WDHistory

	// #GetDepositAddress
	// Input: Coin. Get addresses for mainnet and erc20.
	DepositAddresses map[ChainType]*DepositAddr // key: chainType

	//#Debug
	DebugMode  bool   `json:"debug_mode"`
	RequestURI string `json:"request_uri"`

	// MapParams    string `json:"map_params"`
	CallResponce string `json:"call_responce"`
	Error        error  `json:"error"`

	// ##### New Changes - Contract
	// OperationType  WalletType `json:"operation_type"`  // replace by walelt!!!
	Pair           *pair.Pair `json:"pair"`
	Rate           float64
	Quantity       float64
	Order          *Order
	OrderDirection TradeDirection
}

type PublicOperation struct {
	ID int `json:"id"` //dummy at this moment for

	Type OperationType `json:"type"`
	EX   ExchangeName  `json:"exchange_name"`

	Coin          *coin.Coin     `json:"op_coin"` //BOT standard symbol, not the symbol on exchange
	Pair          *pair.Pair     `json:"op_pair"`
	Maker         *Maker         `json:"maker"`
	TradeHistory  []*TradeDetail `json:"history"`
	CoinChainType []ChainType    `json:"chain_type"`

	//#Debug
	DebugMode    bool   `json:"debug mode"`
	RequestURI   string `json:"request_uri"`
	CallResponce string `json:"call_responce"`

	//#Network
	Proxy   string `json:"proxy"`
	Timeout int    `json:"timeout"`

	Error error `json:"error"`

	// ##### New Changes - Contract
	Wallet WalletType `json:"wallet"` // Contract/Spot operation. Default spot if empty
}

type WDHistory struct {
	ID        string     `json:"id"`
	Coin      *coin.Coin `json:"wd_history_coin"`
	Quantity  float64    `json:"quantity"`
	Tag       string     `json:"tag"`
	Address   string     `json:"address"`
	TxHash    string     `json:"txhash"`
	ChainType ChainType  `json:"chain_type"`
	Status    string     `json:"status"`
	TimeStamp int64      `json:"timestamp"`
}

type DepositAddr struct {
	Coin    *coin.Coin `json:"wd_history_coin"`
	Address string     `json:"address"`
	Tag     string     `json:"tag"`
	Chain   ChainType  `json:"chain_type"`
}

type AssetBalance struct {
	Coin             *coin.Coin `json:"balance_coin"`
	BalanceAvailable float64    `json:"balance_available"` //the fund able to do trading
	BalanceFrozen    float64    `json:"balance_frozen"`    // the fund in order or frozen can't do trading         the total amount of fund should be   BalanceAvailable + BalanceFrozen

}

type TradeDirection string

const (
	Buy  TradeDirection = "b"
	Sell TradeDirection = "s"
)

type TradeDetail struct {
	ID        string         `json:"id"`
	Quantity  float64        `json:"quantity"`  //amount 	/ Qty
	TimeStamp int64          `json:"timestamp"` //TS ts
	Rate      float64        `json:"rate"`      //Price
	Direction TradeDirection `json:"direction"` //Buy or Sell  /'b' 's'
	BestMatch bool           `json:"best_match"`
}

type ChainTypeRequest struct {
	Exchange  string   `json:"exchange, omitempty"`
	CoinID    int      `json:"coin_id, omitempty"`
	ExSymbol  string   `json:"ex_symbol, omitempty"`
	ChainType []string `json:"chain_type, omitempty"`
	CTSource  string   `json:"ct_source, omitempty"`
}
