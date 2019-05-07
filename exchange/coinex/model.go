package coinex

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type CoinexAccountBalance struct {
	Code    int                           `json:"code"`
	Data    map[string]*CoinexCoinBalance `json:"data"`
	Message string                        `json:"message"`
}

type CoinexCoinBalance struct {
	Available string `json:"available"`
	Frozen    string `json:"frozen"`
}

/* type AccountBalances struct {
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
} */

type OrderResponse struct {
	Code int `json:"code"`
	Data struct {
		Amount       string `json:"amount"`
		AvgPrice     string `json:"avg_price"`
		CreateTime   int    `json:"create_time"`
		DealAmount   string `json:"deal_amount"`
		DealFee      string `json:"deal_fee"`
		DealMoney    string `json:"deal_money"`
		ID           int    `json:"id"`
		Left         string `json:"left"`
		MakerFeeRate string `json:"maker_fee_rate"`
		Market       string `json:"market"`
		OrderType    string `json:"order_type"`
		Price        string `json:"price"`
		Status       string `json:"status"`
		TakerFeeRate string `json:"taker_fee_rate"`
		Type         string `json:"type"`
	} `json:"data"`
	Message string `json:"message"`
}

/* type PlaceOrder struct {
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
} */

type CoinexOrderBook struct {
	Code int `json:"code"`
	Data struct {
		Asks [][]string `json:"asks"`
		Bids [][]string `json:"bids"`
		Last string     `json:"last"`
	} `json:"data"`
	Message string `json:"message"`
}

/* type OrderBook struct {
	LastUpdateID int             `json:"lastUpdateId"`
	Bids         [][]interface{} `json:"bids"`
	Asks         [][]interface{} `json:"asks"`
} */

/* type PairsData struct {
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
} */

type PairsData struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]*PairDetail `json:"data"`
}

type PairDetail struct {
	Symbol         string `json:"name"`
	MinAmount      string `json:"min_amount"`
	MakerFeeRate   string `json:"maker_fee_rate"`
	TakerFeeRate   string `json:"taker_fee_rate"`
	PricingName    string `json:"pricing_name"`
	PricingDecimal int    `json:"pricing_decimal"`
	TradingName    string `json:"trading_name"`
	TradingDecimal int    `json:"trading_decimal"`
}

/* type CoinsData []struct {
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
} */
