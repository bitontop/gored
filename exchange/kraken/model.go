package kraken

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
	"time"
)

type JsonResponse struct {
	Error  []interface{}   `json:"error"`
	Result json.RawMessage `json:"result"`
}

type CoinsData struct {
	Aclass          string `json:"aclass"`
	Altname         string `json:"altname"`
	Decimals        int    `json:"decimals"`
	DisplayDecimals int    `json:"display_decimals"`
}

type PairsData struct {
	Altname           string      `json:"altname"`
	Wsname            string      `json:"wsname"`
	AclassBase        string      `json:"aclass_base"`
	Base              string      `json:"base"`
	AclassQuote       string      `json:"aclass_quote"`
	Quote             string      `json:"quote"`
	Lot               string      `json:"lot"`
	PairDecimals      int         `json:"pair_decimals"`
	LotDecimals       int         `json:"lot_decimals"`
	LotMultiplier     int         `json:"lot_multiplier"`
	LeverageBuy       []int       `json:"leverage_buy"`
	LeverageSell      []int       `json:"leverage_sell"`
	Fees              [][]float64 `json:"fees"`
	FeesMaker         [][]float64 `json:"fees_maker"`
	FeeVolumeCurrency string      `json:"fee_volume_currency"`
	MarginCall        int         `json:"margin_call"`
	MarginStop        int         `json:"margin_stop"`
}

type OrderBook struct {
	Asks [][]interface{} `json:"asks"`
	Bids [][]interface{} `json:"bids"`
}

type WithdrawResponse struct {
	RefID string `json:"refid"`
}

type PlaceOrder struct {
	Description    OrderDescription `json:"descr"`
	TransactionIds []string         `json:"txid"`
}

type OrderDescription struct {
	AssetPair      string `json:"pair"`
	Close          string `json:"close"`
	Leverage       string `json:"leverage"`
	Order          string `json:"order"`
	OrderType      string `json:"ordertype"`
	PrimaryPrice   string `json:"price"`
	SecondaryPrice string `json:"price2"`
	Type           string `json:"type"`
}

type OrderStatus map[string]Order

type Order struct {
	TransactionID  string           `json:"-"`
	ReferenceID    string           `json:"refid"`
	UserRef        int              `json:"userref"`
	Status         string           `json:"status"`
	OpenTime       float64          `json:"opentm"`
	StartTime      float64          `json:"starttm"`
	ExpireTime     float64          `json:"expiretm"`
	Description    OrderDescription `json:"descr"`
	Volume         string           `json:"vol"`
	VolumeExecuted float64          `json:"vol_exec,string"`
	Cost           float64          `json:"cost,string"`
	Fee            float64          `json:"fee,string"`
	Price          float64          `json:"price,string"`
	StopPrice      float64          `json:"stopprice.string"`
	LimitPrice     float64          `json:"limitprice,string"`
	Misc           string           `json:"misc"`
	OrderFlags     string           `json:"oflags"`
	CloseTime      float64          `json:"closetm"`
	Reason         string           `json:"reason"`
}

type CancelOrder struct {
	Count   int  `json:"count"`
	Pending bool `json:"pending"`
}

type Trade struct {
	Symbol    string    `json:"symbol"`
	Price     string    `json:"price"`
	Volume    string    `json:"volume"`
	Direction string    `json:"direction"`
	TradeTime time.Time `json:"tradeTime"`
}

/* type AccountBalances struct {
	ADA  float64 `json:"ADA,string"`
	ATOM float64 `json:"ATOM,string"`
	BCH  float64 `json:"BCH,string"`
	BSV  float64 `json:"BSV,string"`
	DASH float64 `json:"DASH,string"`
	EOS  float64 `json:"EOS,string"`
	GNO  float64 `json:"GNO,string"`
	KFEE float64 `json:"ADA,string"`
	QTUM float64 `json:"QTUM,string"`
	USDT float64 `json:"USDT,string"`
	XDAO float64 `json:"XDAO,string"`
	XETC float64 `json:"XETC,string"`
	XETH float64 `json:"XETH,string"`
	XICN float64 `json:"XICN,string"`
	XLTC float64 `json:"XLTC,string"`
	XMLN float64 `json:"XMLN,string"`
	XNMC float64 `json:"XNMC,string"`
	XREP float64 `json:"XREP,string"`
	XTZ  float64 `json:"XTZ,string"`
	XXBT float64 `json:"XXBT,string"`
	XXDG float64 `json:"XXDG,string"`
	XXLM float64 `json:"XXLM,string"`
	XXMR float64 `json:"XXMR,string"`
	XXRP float64 `json:"XXRP,string"`
	XXVN float64 `json:"XXVN,string"`
	XZEC float64 `json:"XZEC,string"`
	ZCAD float64 `json:"ZCAD,string"`
	ZEUR float64 `json:"ZEUR,string"`
	ZGBP float64 `json:"ZGBP,string"`
	ZJPY float64 `json:"ZJPY,string"`
	ZKRW float64 `json:"ZKRW,string"`
	ZUSD float64 `json:"ZUSD,string"`
} */
