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
	PairID               int
	Pair                 *pair.Pair //the code on excahnge with the same chain, eg: BCH, BCC on different exchange, but they are the same chain
	ExID                 string
	ExSymbol             string
	MakerFee             float64
	TakerFee             float64
	LotSize              float64 // the decimal place for this coin on exchange for the pairs, eg:  BTC: 0.00001    NEO:1   LTC: 0.001 ETH:0.01
	PriceFilter          float64
	MinTradeQuantity     float64 // minimum trade target quantity
	MinTradeBaseQuantity float64 // minimum trade base quantity
	Listed               bool
	Issue                string //the issue for the pair if have any problem
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
	Rate          float64   `bson:"Rate"`
	Quantity      float64   `bson:"Quantity"`
	Side          OrderType //TODO  SIDE ==> Direction, depreated after all changed. All gored, goredws fixed.
	Direction     TradeDirection
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
	Withdraw           OperationType = "Withdraw"
	Transfer           OperationType = "Transfer"           // transfer  between inneral wallet
	SubAccountTransfer OperationType = "SubAccountTransfer" // transfer  between main/sub account
	Balance            OperationType = "Balance"            // balance(s) of different accounts
	BalanceList        OperationType = "BalanceAll"         // balance(s) of different accounts
	SubBalanceList     OperationType = "SubBalanceList"     // balance(s) of subaccount
	SubAllBalanceList  OperationType = "SubAllBalanceList"  // balance(s) of all subaccounts
	GetSubAccountList  OperationType = "GetSubAccountList"  // get sub accounts list
	GetPositionInfo    OperationType = "GetPositionInfo"    // position information for Contract

	//Public Query
	GetCoin OperationType = "GetCoin"
	GetPair OperationType = "GetPair"

	TradeHistory   OperationType = "TradeHistory"
	Orderbook      OperationType = "Orderbook"
	CoinChainType  OperationType = "CoinChainType"
	KLine          OperationType = "KLine"
	GetTickerPrice OperationType = "GetTickerPrice"

	//Trade (Private Action)
	PlaceOrder     OperationType = "PlaceOrder"
	CancelOrder    OperationType = "CancelOrder"
	GetOrderStatus OperationType = "GetOrderStatus"

	//User (Private Action)
	GetOpenOrder         OperationType = "GetOpenOrder"    // New and Partial Orders
	GetOrderHistory      OperationType = "GetOrderHistory" // All Orders other than open orders
	GetDepositHistory    OperationType = "GetDepositHistory"
	GetWithdrawalHistory OperationType = "GetWithdrawalHistory"
	GetTransferHistory   OperationType = "GetTransferHistory"
	GetDepositAddress    OperationType = "GetDepositAddress" // Get address for one coin

	// Contract
	GetFutureBalance  OperationType = "GetFutureBalance"
	SetFutureLeverage OperationType = "SetFutureLeverage"
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

	Type     OperationType `json:"type"`
	Ex       ExchangeName  `json:"exchange_name"`
	Sandbox  bool          `json:"sandbox"`
	TestMode bool          `json:"test_mode"`

	//#Transfer,TransferHistory,Balance,Withdraw
	Coin *coin.Coin `json:"transfer_coin"` //BOT standard symbol, not the symbol on exchange

	//specific operations
	// #Inner Transfer
	TransferFrom        WalletType `json:"transfer_from"`
	TransferDestination WalletType `json:"transfer_dest"`
	TransferAmount      string     `json:"transfer_amount"`
	TransferStartTime   int64      `json:"transfer_start_time"`
	TransferEndTime     int64      `json:"transfer_end_time"`

	// #SubAccount Transfer
	SubTransferFrom   string `json:"sub_transfer_from"`
	SubTransferTo     string `json:"sub_transfer_to"`
	SubTransferAmount string `json:"sub_transfer_amount"`

	// start/end timestamp
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`

	// #Withdraw
	WithdrawAddress string `json:"withdraw_address"`
	WithdrawTag     string `json:"withdraw_tag"`
	WithdrawAmount  string `json:"withdraw_amount"` //here using string instead of float64
	WithdrawID      string `json:"withdraw_id"`
	WithdrawChain   string `json:"withdraw_chain"`

	// #Balance
	Wallet       WalletType `json:"wallet"`         // Contract/Spot operation. Default spot if empty
	SubAccountID string     `json:"sub_account_id"` // Sub account id. eg. Sub account email

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

	// #GetSubAccountList
	SubAccountList []*SubAccountInfo

	// #Sub Account Transfer History
	SubUserName        string // coinex only
	TransferInHistory  []*TransferHistory
	TransferOutHistory []*TransferHistory

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
	StopRate       float64 // for STOP_LIMIT, STOP_MARKET
	Quantity       float64
	Order          *Order
	OrderType      OrderPriceType // eg. FOK
	TradeType      OrderTradeType // eg. TRADE_LIMIT
	OrderDirection TradeDirection //TradeDirection
	Leverage       int
}

type PublicOperation struct {
	ID int `json:"id"` //dummy at this moment for

	Type     OperationType `json:"type"`
	EX       ExchangeName  `json:"exchange_name"`
	TestMode bool          `json:"test_mode"`

	Coin           *coin.Coin           `json:"op_coin"` //BOT standard symbol, not the symbol on exchange
	Pair           *pair.Pair           `json:"op_pair"`
	Maker          *Maker               `json:"maker"`
	TradeHistory   []*TradeDetail       `json:"history"`
	CoinChainType  []ChainType          `json:"chain_type"`
	KlineInterval  string               `json:"kline_interval"`
	KlineStartTime int64                `json:"kline_start_time"`
	KlineEndTime   int64                `json:"kline_end_time"`
	Kline          []*KlineDetail       `json:"kline"`
	TickerPrice    []*TickerPriceDetail `json:"ticker_price"`

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

type TransferType string

const (
	TransferIn  TransferType = "TransferIn"
	TransferOut TransferType = "TransferOut"
)

type TransferHistory struct {
	ID        string       `json:"id"`
	Coin      *coin.Coin   `json:"transfer_history_coin"`
	Type      TransferType `json:"type"`
	Quantity  float64      `json:"quantity"`
	TimeStamp int64        `json:"timestamp"`
	StatusMsg string       `json: status_msg`
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
	Balance          float64    `json:"balance"`
	BalanceAvailable float64    `json:"balance_available"` //the fund able to do trading
	BalanceFrozen    float64    `json:"balance_frozen"`    // the fund in order or frozen can't do trading         the total amount of fund should be   BalanceAvailable + BalanceFrozen

}

type SubAccountInfo struct {
	ID          string     `json:"id"` // account ID, email, etc.
	Status      string     `json:"status"`
	Activated   bool       `json:"activated"`
	AccountType WalletType `json:"account_type"`
	TimeStamp   int64      `json:"timestamp"`
}

type KlineDetail struct {
	ID                  int          `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Exchange            ExchangeName `json:"exchange"`
	Pair                string       `json:"pair"`
	OpenTime            float64      `json:"open_time"`
	Open                float64      `json:"open"`
	High                float64      `json:"high"`
	Low                 float64      `json:"low"`
	Close               float64      `json:"close"`
	Volume              float64      `json:"volume"`
	CloseTime           float64      `json:"close_time"`
	QuoteAssetVolume    float64      `json:"quote_asset_volume"`
	TradesCount         float64      `json:"trades_count"`
	TakerBuyBaseVolume  float64      `json:"taker_buy_base_volume"`
	TakerBuyQuoteVolume float64      `json:"taker_buy_quote_volume"`
}

type TickerPriceDetail struct {
	Pair  *pair.Pair `json:"pair"`
	Price float64    `json:"price"`
}

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
