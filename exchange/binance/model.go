package binance

import "time"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type AccountBalances struct {
	MakerCommission  int  `json:"makerCommission"`
	TakerCommission  int  `json:"takerCommission"`
	BuyerCommission  int  `json:"buyerCommission"`
	SellerCommission int  `json:"sellerCommission"`
	CanTrade         bool `json:"canTrade"`
	CanWithdraw      bool `json:"canWithdraw"`
	CanDeposit       bool `json:"canDeposit"`
	Balances         []struct {
		Asset  string `json:"asset"`
		Free   string `json:"free"`
		Locked string `json:"locked"`
	} `json:"balances"`
}

type WithdrawResponse struct {
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
	ID      string `json:"id"`
}

type PlaceOrder struct {
	Symbol        string `json:"symbol"`
	OrderID       int    `json:"orderId"`
	ClientOrderID string `json:"clientOrderId"`
	TransactTime  int64  `json:"transactTime"`
	Price         string `json:"price"`
	OrigQty       string `json:"origQty"`
	ExecutedQty   string `json:"executedQty"`
	Status        string `json:"status"`
	TimeInForce   string `json:"timeInForce"`
	Type          string `json:"type"`
	Side          string `json:"side"`
	StopPrice     string `json:"stopPrice"`
	IcebergQty    string `json:"icebergQty"`
	Time          int64  `json:"time"`
	Code          int    `json:"code"`
	Msg           string `json:"msg"`
}

type OrderBook struct {
	LastUpdateID int             `json:"lastUpdateId"`
	Bids         [][]interface{} `json:"bids"`
	Asks         [][]interface{} `json:"asks"`
}

type ContractOrderBook struct {
	LastUpdateID int64      `json:"lastUpdateId"`
	E            int64      `json:"E"`
	T            int64      `json:"T"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

type PairsData struct {
	Timezone   string `json:"timezone"`
	ServerTime int64  `json:"serverTime"`
	RateLimits []struct {
		RateLimitType string `json:"rateLimitType"`
		Interval      string `json:"interval"`
		Limit         int    `json:"limit"`
	} `json:"rateLimits"`
	ExchangeFilters []interface{} `json:"exchangeFilters"`
	Symbols         []struct {
		Symbol             string   `json:"symbol"`
		Status             string   `json:"status"`
		BaseAsset          string   `json:"baseAsset"`
		BaseAssetPrecision int      `json:"baseAssetPrecision"`
		QuoteAsset         string   `json:"quoteAsset"`
		QuotePrecision     int      `json:"quotePrecision"`
		OrderTypes         []string `json:"orderTypes"`
		IcebergAllowed     bool     `json:"icebergAllowed"`
		Filters            []struct {
			FilterType       string `json:"filterType"`
			MinPrice         string `json:"minPrice,omitempty"`
			MaxPrice         string `json:"maxPrice,omitempty"`
			TickSize         string `json:"tickSize,omitempty"`
			MinQty           string `json:"minQty,omitempty"`
			MaxQty           string `json:"maxQty,omitempty"`
			StepSize         string `json:"stepSize,omitempty"`
			MinNotional      string `json:"minNotional,omitempty"`
			Limit            int    `json:"limit,omitempty"`
			MaxNumAlgoOrders int    `json:"maxNumAlgoOrders,omitempty"`
		} `json:"filters"`
	} `json:"symbols"`
}

type CoinsData []struct {
	ID                      string      `json:"id"`
	AssetCode               string      `json:"assetCode"`
	AssetName               string      `json:"assetName"`
	Unit                    string      `json:"unit"`
	TransactionFee          float64     `json:"transactionFee"`
	CommissionRate          float64     `json:"commissionRate"`
	FreeAuditWithdrawAmt    float64     `json:"freeAuditWithdrawAmt"`
	FreeUserChargeAmount    float64     `json:"freeUserChargeAmount"`
	MinProductWithdraw      string      `json:"minProductWithdraw"`
	WithdrawIntegerMultiple string      `json:"withdrawIntegerMultiple"`
	ConfirmTimes            string      `json:"confirmTimes"`
	ChargeLockConfirmTimes  interface{} `json:"chargeLockConfirmTimes"`
	CreateTime              interface{} `json:"createTime"`
	Test                    int         `json:"test"`
	URL                     string      `json:"url"`
	AddressURL              string      `json:"addressUrl"`
	BlockURL                string      `json:"blockUrl"`
	EnableCharge            bool        `json:"enableCharge"`
	EnableWithdraw          bool        `json:"enableWithdraw"`
	RegEx                   string      `json:"regEx"`
	RegExTag                string      `json:"regExTag"`
	Gas                     float64     `json:"gas"`
	ParentCode              string      `json:"parentCode"`
	IsLegalMoney            bool        `json:"isLegalMoney"`
	ReconciliationAmount    float64     `json:"reconciliationAmount"`
	SeqNum                  string      `json:"seqNum"`
	ChineseName             string      `json:"chineseName"`
	CnLink                  string      `json:"cnLink"`
	EnLink                  string      `json:"enLink"`
	LogoURL                 string      `json:"logoUrl"`
	FullLogoURL             string      `json:"fullLogoUrl"`
	ForceStatus             bool        `json:"forceStatus"`
	ResetAddressStatus      bool        `json:"resetAddressStatus"`
	ChargeDescCn            interface{} `json:"chargeDescCn"`
	ChargeDescEn            interface{} `json:"chargeDescEn"`
	AssetLabel              interface{} `json:"assetLabel"`
	SameAddress             bool        `json:"sameAddress"`
	DepositTipStatus        bool        `json:"depositTipStatus"`
	DynamicFeeStatus        bool        `json:"dynamicFeeStatus"`
	DepositTipEn            interface{} `json:"depositTipEn"`
	DepositTipCn            interface{} `json:"depositTipCn"`
	AssetLabelEn            interface{} `json:"assetLabelEn"`
	SupportMarket           interface{} `json:"supportMarket"`
	FeeReferenceAsset       string      `json:"feeReferenceAsset"`
	FeeRate                 float64     `json:"feeRate"`
	FeeDigit                int         `json:"feeDigit"`
	AssetDigit              int         `json:"assetDigit"`
	LegalMoney              bool        `json:"legalMoney"`
}

type TradeHistory []struct {
	ID           int    `json:"id"`
	Price        string `json:"price"`
	Qty          string `json:"qty"`
	QuoteQty     string `json:"quoteQty"`
	Time         int64  `json:"time"`
	IsBuyerMaker bool   `json:"isBuyerMaker"`
	IsBestMatch  bool   `json:"isBestMatch"`
}

type ContractPlaceOrder struct {
	ClientOrderID string `json:"clientOrderId"`
	CumQuote      string `json:"cumQuote"`
	ExecutedQty   string `json:"executedQty"`
	OrderID       int    `json:"orderId"`
	OrigQty       string `json:"origQty"`
	Price         string `json:"price"`
	ReduceOnly    bool   `json:"reduceOnly"`
	Side          string `json:"side"`
	Status        string `json:"status"`
	StopPrice     string `json:"stopPrice"`
	Symbol        string `json:"symbol"`
	TimeInForce   string `json:"timeInForce"`
	Type          string `json:"type"`
	ActivatePrice string `json:"activatePrice"`
	PriceRate     string `json:"priceRate"`
	UpdateTime    int64  `json:"updateTime"`
	WorkingType   string `json:"workingType"`
}

type ContractOrderStatus struct {
	AvgPrice      string `json:"avgPrice"`
	ClientOrderID string `json:"clientOrderId"`
	CumQuote      string `json:"cumQuote"`
	ExecutedQty   string `json:"executedQty"`
	OrderID       int    `json:"orderId"`
	OrigQty       string `json:"origQty"`
	OrigType      string `json:"origType"`
	Price         string `json:"price"`
	ReduceOnly    bool   `json:"reduceOnly"`
	Side          string `json:"side"`
	Status        string `json:"status"`
	StopPrice     string `json:"stopPrice"`
	Symbol        string `json:"symbol"`
	Time          int64  `json:"time"`
	TimeInForce   string `json:"timeInForce"`
	Type          string `json:"type"`
	ActivatePrice string `json:"activatePrice"`
	PriceRate     string `json:"priceRate"`
	UpdateTime    int64  `json:"updateTime"`
	WorkingType   string `json:"workingType"`
}

type ContractCancelOrder struct {
	ClientOrderID string `json:"clientOrderId"`
	CumQty        string `json:"cumQty"`
	CumQuote      string `json:"cumQuote"`
	ExecutedQty   string `json:"executedQty"`
	OrderID       int    `json:"orderId"`
	OrigQty       string `json:"origQty"`
	Price         string `json:"price"`
	ReduceOnly    bool   `json:"reduceOnly"`
	Side          string `json:"side"`
	Status        string `json:"status"`
	StopPrice     string `json:"stopPrice"`
	Symbol        string `json:"symbol"`
	TimeInForce   string `json:"timeInForce"`
	OrigType      string `json:"origType"`
	Type          string `json:"type"`
	ActivatePrice string `json:"activatePrice"`
	PriceRate     string `json:"priceRate"`
	UpdateTime    int64  `json:"updateTime"`
	WorkingType   string `json:"workingType"`
}

type ContractBalance struct {
	Assets []struct {
		Asset                  string `json:"asset"`
		InitialMargin          string `json:"initialMargin"`
		MaintMargin            string `json:"maintMargin"`
		MarginBalance          string `json:"marginBalance"`
		MaxWithdrawAmount      string `json:"maxWithdrawAmount"`
		OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
		PositionInitialMargin  string `json:"positionInitialMargin"`
		UnrealizedProfit       string `json:"unrealizedProfit"`
		WalletBalance          string `json:"walletBalance"`
	} `json:"assets"`
	CanDeposit        bool   `json:"canDeposit"`
	CanTrade          bool   `json:"canTrade"`
	CanWithdraw       bool   `json:"canWithdraw"`
	FeeTier           int    `json:"feeTier"`
	MaxWithdrawAmount string `json:"maxWithdrawAmount"`
	Positions         []struct {
		Isolated               bool   `json:"isolated"`
		Leverage               string `json:"leverage"`
		InitialMargin          string `json:"initialMargin"`
		MaintMargin            string `json:"maintMargin"`
		OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
		PositionInitialMargin  string `json:"positionInitialMargin"`
		Symbol                 string `json:"symbol"`
		UnrealizedProfit       string `json:"unrealizedProfit"`
	} `json:"positions"`
	TotalInitialMargin          string `json:"totalInitialMargin"`
	TotalMaintMargin            string `json:"totalMaintMargin"`
	TotalMarginBalance          string `json:"totalMarginBalance"`
	TotalOpenOrderInitialMargin string `json:"totalOpenOrderInitialMargin"`
	TotalPositionInitialMargin  string `json:"totalPositionInitialMargin"`
	TotalUnrealizedProfit       string `json:"totalUnrealizedProfit"`
	TotalWalletBalance          string `json:"totalWalletBalance"`
	UpdateTime                  int    `json:"updateTime"`
}

// private operation
type OpenOrders []struct {
	Symbol              string `json:"symbol"`
	OrderID             int    `json:"orderId"`
	OrderListID         int    `json:"orderListId"`
	ClientOrderID       string `json:"clientOrderId"`
	Price               string `json:"price"`
	OrigQty             string `json:"origQty"`
	ExecutedQty         string `json:"executedQty"`
	CummulativeQuoteQty string `json:"cummulativeQuoteQty"`
	Status              string `json:"status"`
	TimeInForce         string `json:"timeInForce"`
	Type                string `json:"type"`
	Side                string `json:"side"`
	StopPrice           string `json:"stopPrice"`
	IcebergQty          string `json:"icebergQty"`
	Time                int64  `json:"time"`
	UpdateTime          int64  `json:"updateTime"`
	IsWorking           bool   `json:"isWorking"`
	OrigQuoteOrderQty   string `json:"origQuoteOrderQty"`
}

type CloseOrders []struct {
	Symbol          string `json:"symbol"`
	ID              int    `json:"id"`
	OrderID         int    `json:"orderId"`
	OrderListID     int    `json:"orderListId"`
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	QuoteQty        string `json:"quoteQty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
	Time            int64  `json:"time"`
	IsBuyer         bool   `json:"isBuyer"`
	IsMaker         bool   `json:"isMaker"`
	IsBestMatch     bool   `json:"isBestMatch"`
}

type WithdrawHistory []struct {
	Address         string    `json:"address"`
	Amount          string    `json:"amount"`
	ApplyTime       time.Time `json:"applyTime"`
	Coin            string    `json:"coin"`
	ID              string    `json:"id"`
	WithdrawOrderID string    `json:"withdrawOrderId,omitempty"`
	Network         string    `json:"network"`
	TransferType    int       `json:"transferType"`
	Status          int       `json:"status"`
	TxID            string    `json:"txId"`
}

type DepositHistory []struct {
	Address    string `json:"address"`
	AddressTag string `json:"addressTag"`
	Amount     string `json:"amount"`
	Coin       string `json:"coin"`
	InsertTime int64  `json:"insertTime"`
	Network    string `json:"network"`
	Status     int    `json:"status"`
	TxID       string `json:"txId"`
}

type DepositAddress struct {
	Code    int    `json:"code"`
	Address string `json:"address"`
	Coin    string `json:"coin"`
	Tag     string `json:"tag"`
	URL     string `json:"url"`
}
