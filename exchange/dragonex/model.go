package dragonex

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

type JsonResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Ok   bool            `json:"ok"`
	Data json.RawMessage `json:"data"`
}

type CoinsData []struct {
	CoinID int    `json:"coin_id"`
	Code   string `json:"code"`
}

type PairsData struct {
	Columns []string        `json:"columns"`
	List    [][]interface{} `json:"list"`
}

type OrderBook struct {
	Buys []struct {
		Price  string `json:"price"`
		Volume string `json:"volume"`
	} `json:"buys"`
	Sells []struct {
		Price  string `json:"price"`
		Volume string `json:"volume"`
	} `json:"sells"`
}
