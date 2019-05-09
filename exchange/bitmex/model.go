package bitmex

import "time"

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Name    string `json:"name"`
	} `json:"error"`
}

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

type PlaceOrder struct {
	OrderID               string      `json:"orderID"`
	ClOrdID               string      `json:"clOrdID"`
	ClOrdLinkID           string      `json:"clOrdLinkID"`
	Account               int         `json:"account"`
	Symbol                string      `json:"symbol"`
	Side                  string      `json:"side"`
	SimpleOrderQty        float64     `json:"simpleOrderQty"`
	OrderQty              interface{} `json:"orderQty"`
	Price                 float64     `json:"price"`
	DisplayQty            interface{} `json:"displayQty"`
	StopPx                interface{} `json:"stopPx"`
	PegOffsetValue        interface{} `json:"pegOffsetValue"`
	PegPriceType          string      `json:"pegPriceType"`
	Currency              string      `json:"currency"`
	SettlCurrency         string      `json:"settlCurrency"`
	OrdType               string      `json:"ordType"`
	TimeInForce           string      `json:"timeInForce"`
	ExecInst              string      `json:"execInst"`
	ContingencyType       string      `json:"contingencyType"`
	ExDestination         string      `json:"exDestination"`
	OrdStatus             string      `json:"ordStatus"`
	Triggered             string      `json:"triggered"`
	WorkingIndicator      bool        `json:"workingIndicator"`
	OrdRejReason          string      `json:"ordRejReason"`
	SimpleLeavesQty       float64     `json:"simpleLeavesQty"`
	LeavesQty             float64     `json:"leavesQty"`
	SimpleCumQty          float64     `json:"simpleCumQty"`
	CumQty                float64     `json:"cumQty"`
	AvgPx                 float64     `json:"avgPx"`
	MultiLegReportingType string      `json:"multiLegReportingType"`
	Text                  string      `json:"text"`
	TransactTime          time.Time   `json:"transactTime"`
	Timestamp             time.Time   `json:"timestamp"`
}

type OrderBook []struct {
	Symbol string  `json:"symbol"`
	ID     int64   `json:"id"`
	Side   string  `json:"side"`
	Size   float64 `json:"size"`
	Price  float64 `json:"price"`
}

type PairsData []struct {
	Symbol                         string      `json:"symbol"`
	RootSymbol                     string      `json:"rootSymbol"`
	State                          string      `json:"state"`
	Typ                            string      `json:"typ"`
	Listing                        time.Time   `json:"listing"`
	Front                          time.Time   `json:"front"`
	Expiry                         time.Time   `json:"expiry"`
	Settle                         time.Time   `json:"settle"`
	RelistInterval                 interface{} `json:"relistInterval"`
	InverseLeg                     string      `json:"inverseLeg"`
	SellLeg                        string      `json:"sellLeg"`
	BuyLeg                         string      `json:"buyLeg"`
	OptionStrikePcnt               interface{} `json:"optionStrikePcnt"`
	OptionStrikeRound              interface{} `json:"optionStrikeRound"`
	OptionStrikePrice              interface{} `json:"optionStrikePrice"`
	OptionMultiplier               interface{} `json:"optionMultiplier"`
	PositionCurrency               string      `json:"positionCurrency"`
	Underlying                     string      `json:"underlying"`
	QuoteCurrency                  string      `json:"quoteCurrency"`
	UnderlyingSymbol               string      `json:"underlyingSymbol"`
	Reference                      string      `json:"reference"`
	ReferenceSymbol                string      `json:"referenceSymbol"`
	CalcInterval                   interface{} `json:"calcInterval"`
	PublishInterval                interface{} `json:"publishInterval"`
	PublishTime                    interface{} `json:"publishTime"`
	MaxOrderQty                    float64     `json:"maxOrderQty"`
	MaxPrice                       float64     `json:"maxPrice"`
	LotSize                        float64     `json:"lotSize"`
	TickSize                       float64     `json:"tickSize"`
	Multiplier                     int         `json:"multiplier"`
	SettlCurrency                  string      `json:"settlCurrency"`
	UnderlyingToPositionMultiplier int         `json:"underlyingToPositionMultiplier"`
	UnderlyingToSettleMultiplier   interface{} `json:"underlyingToSettleMultiplier"`
	QuoteToSettleMultiplier        int         `json:"quoteToSettleMultiplier"`
	IsQuanto                       bool        `json:"isQuanto"`
	IsInverse                      bool        `json:"isInverse"`
	InitMargin                     float64     `json:"initMargin"`
	MaintMargin                    float64     `json:"maintMargin"`
	RiskLimit                      int64       `json:"riskLimit"`
	RiskStep                       int64       `json:"riskStep"`
	Limit                          interface{} `json:"limit"`
	Capped                         bool        `json:"capped"`
	Taxed                          bool        `json:"taxed"`
	Deleverage                     bool        `json:"deleverage"`
	MakerFee                       float64     `json:"makerFee"`
	TakerFee                       float64     `json:"takerFee"`
	SettlementFee                  float64     `json:"settlementFee"`
	InsuranceFee                   int         `json:"insuranceFee"`
	FundingBaseSymbol              string      `json:"fundingBaseSymbol"`
	FundingQuoteSymbol             string      `json:"fundingQuoteSymbol"`
	FundingPremiumSymbol           string      `json:"fundingPremiumSymbol"`
	FundingTimestamp               interface{} `json:"fundingTimestamp"`
	FundingInterval                interface{} `json:"fundingInterval"`
	FundingRate                    interface{} `json:"fundingRate"`
	IndicativeFundingRate          interface{} `json:"indicativeFundingRate"`
	RebalanceTimestamp             interface{} `json:"rebalanceTimestamp"`
	RebalanceInterval              interface{} `json:"rebalanceInterval"`
	OpeningTimestamp               time.Time   `json:"openingTimestamp"`
	ClosingTimestamp               time.Time   `json:"closingTimestamp"`
	SessionInterval                time.Time   `json:"sessionInterval"`
	PrevClosePrice                 float64     `json:"prevClosePrice"`
	LimitDownPrice                 interface{} `json:"limitDownPrice"`
	LimitUpPrice                   interface{} `json:"limitUpPrice"`
	BankruptLimitDownPrice         interface{} `json:"bankruptLimitDownPrice"`
	BankruptLimitUpPrice           interface{} `json:"bankruptLimitUpPrice"`
	PrevTotalVolume                float64     `json:"prevTotalVolume"`
	TotalVolume                    float64     `json:"totalVolume"`
	Volume                         float64     `json:"volume"`
	Volume24H                      float64     `json:"volume24h"`
	PrevTotalTurnover              int64       `json:"prevTotalTurnover"`
	TotalTurnover                  int64       `json:"totalTurnover"`
	Turnover                       int64       `json:"turnover"`
	Turnover24H                    int64       `json:"turnover24h"`
	HomeNotional24H                float64     `json:"homeNotional24h"`
	ForeignNotional24H             float64     `json:"foreignNotional24h"`
	PrevPrice24H                   float64     `json:"prevPrice24h"`
	Vwap                           float64     `json:"vwap"`
	HighPrice                      float64     `json:"highPrice"`
	LowPrice                       float64     `json:"lowPrice"`
	LastPrice                      float64     `json:"lastPrice"`
	LastPriceProtected             float64     `json:"lastPriceProtected"`
	LastTickDirection              string      `json:"lastTickDirection"`
	LastChangePcnt                 float64     `json:"lastChangePcnt"`
	BidPrice                       float64     `json:"bidPrice"`
	MidPrice                       float64     `json:"midPrice"`
	AskPrice                       float64     `json:"askPrice"`
	ImpactBidPrice                 float64     `json:"impactBidPrice"`
	ImpactMidPrice                 float64     `json:"impactMidPrice"`
	ImpactAskPrice                 float64     `json:"impactAskPrice"`
	HasLiquidity                   bool        `json:"hasLiquidity"`
	OpenInterest                   float64     `json:"openInterest"`
	OpenValue                      int64       `json:"openValue"`
	FairMethod                     string      `json:"fairMethod"`
	FairBasisRate                  float64     `json:"fairBasisRate"`
	FairBasis                      float64     `json:"fairBasis"`
	FairPrice                      float64     `json:"fairPrice"`
	MarkMethod                     string      `json:"markMethod"`
	MarkPrice                      float64     `json:"markPrice"`
	IndicativeTaxRate              float64     `json:"indicativeTaxRate"`
	IndicativeSettlePrice          float64     `json:"indicativeSettlePrice"`
	OptionUnderlyingPrice          interface{} `json:"optionUnderlyingPrice"`
	SettledPrice                   interface{} `json:"settledPrice"`
	Timestamp                      time.Time   `json:"timestamp"`
}
