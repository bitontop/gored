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
	ID         int
	Code       string `json:"code"`
	Name       string `json:"name"`
	Website    string `json:"website"`
	Explorer   string `json:"explorer"`
	Health     string // the health of the chain
	Blockheigh int
	Blocktime  int // in seconds
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
	coin := GetCoin(code)
	if coin == nil {
		return -1 // return -1 if coin not found
	}
	return coin.ID
}

func GetCoin(code string) *Coin {
	code = strings.TrimSpace(strings.ToUpper(code)) //trim for psql space

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

func ExistID(id int) bool {
	if tmp, ok := coinMap.Get(fmt.Sprintf("%d", id)); ok {
		c := tmp.(*Coin)
		if c.Code != "" {
			// log.Printf("Exist id, coin: %+v", c)
			return true
		}
	}
	return false
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
		} else if ExistID(coin.ID) {
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
