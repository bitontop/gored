package coin

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	cmap "github.com/orcaman/concurrent-map"
)

// modeling structure and functions,
// don't modify unless bugs or new features

type Coin struct { // might also understand as public chain or token
	ID           int
	Code         string `json:"code"`
	Name         string `json:"name"`
	Website      string `json:"website"`
	Explorer     string `json:"explorer"`
	CurrencyType string
	Health       string // the health of the chain
	Blockheigh   int
	Blocktime    int // in seconds
	// Chain      ChainType
}

var LastID int
var coinMap cmap.ConcurrentMap

func Init() {
	if coinMap == nil {
		coinMap = cmap.New()
	}
}

func GenerateCoinID() int {
	return LastID + 1
}

func GetCoinByID(id int) *Coin {
	if tmp, ok := coinMap.Get(fmt.Sprintf("%d", id)); ok {
		return tmp.(*Coin)
	}
	return nil
}

func GetCoinID(code string) int {
	return GetCoin(code).ID
}

func GetCoin(code string) *Coin {
	code = strings.ToUpper(code)
	for _, id := range coinMap.Keys() {
		if tmp, ok := coinMap.Get(id); ok {
			c := tmp.(*Coin)
			if c.Code == code {
				return c
			}
		}
	}
	return nil
}

func GetCoins() []*Coin {
	coins := []*Coin{}
	keySort := []int{}
	for _, key := range coinMap.Keys() {
		id, _ := strconv.Atoi(key)
		keySort = append(keySort, id)
	}
	sort.Ints(keySort)
	for _, key := range keySort {
		coins = append(coins, GetCoinByID(key))
	}
	return coins
}

func AddCoin(coin *Coin) error {
	if coin != nil && coin.Code != "" {
		if coin.ID == 0 {
			coin.ID = GenerateCoinID()
		}
		coin.Code = strings.ToUpper(coin.Code)
		coinMap.Set(fmt.Sprintf("%d", coin.ID), coin)
		LastID = coin.ID
	} else {
		return errors.New("code is not assign yet")
	}
	return nil
}

func DeleteCoin(coin *Coin) {
	coinMap.Remove(fmt.Sprintf("%d", coin.ID))
}
