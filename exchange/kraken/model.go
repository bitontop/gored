package kraken

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
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
