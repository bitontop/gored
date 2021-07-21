package pair

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	coin "github.com/bitontop/gored/coin"
	cmap "github.com/orcaman/concurrent-map"
)

//Pair is the common name pairs across diff excahnges
type Pair struct {
	ID     int
	Name   string
	Base   *coin.Coin
	Target *coin.Coin

	Symbol string //hard fix for missing constrian
}

var pairMap cmap.ConcurrentMap

func Init() {
	if pairMap == nil {
		pairMap = cmap.New()
	}
}

func GeneratePairID() int {
	if pairMap.Count() > 0 {
		keySort := []int{}
		for _, key := range pairMap.Keys() {
			if id, err := strconv.Atoi(key); err == nil {
				keySort = append(keySort, id)
			}
		}
		sort.Ints(keySort)
		return keySort[len(keySort)-1] + 1
	} else {
		return 1
	}
}

func SetPair(id int, base, target *coin.Coin) *Pair {
	if base != nil && target != nil {
		key := GetKey(base, target)
		p := &Pair{id, key, base, target}
		pairMap.Set(fmt.Sprintf("%d", id), p)
		return p
	} else {
		return nil
	}
}

func GetPairID(name string) int {
	return GetPairByKey(name).ID
}

func GetPair(base, target *coin.Coin) *Pair {
	for _, id := range pairMap.Keys() {
		if tmp, ok := pairMap.Get(id); ok {
			p := tmp.(*Pair)
			if p.Base.ID == base.ID && p.Target.ID == target.ID {
				return p
			}
		}
	}

	return SetPair(GeneratePairID(), base, target)
}

func GetString(pair *Pair) string {
	key := GetKey(pair.Base, pair.Target)
	return key
}

func GetKey(base, target *coin.Coin) string {
	key := ""
	if base != nil && target != nil {
		key = (base.Code + coin.SEPARATOR + target.Code)
	}
	return key
}

func GetPairByKey(key string) *Pair {
	key = strings.ToUpper(key)
	for _, id := range pairMap.Keys() {
		if tmp, ok := pairMap.Get(id); ok {
			p := tmp.(*Pair)
			if p.Name == key {
				return p
			}
		}
	}
	return nil
}

func GetPairByID(id int) *Pair {
	if tmp, ok := pairMap.Get(fmt.Sprintf("%d", id)); ok {
		return tmp.(*Pair)
	}
	return nil
}

func GetPairs() []*Pair {
	pairs := []*Pair{}
	keySort := []int{}
	for _, key := range pairMap.Keys() {
		if id, err := strconv.Atoi(key); err == nil {
			keySort = append(keySort, id)
		}
	}
	sort.Ints(keySort)
	for _, key := range keySort {
		pairs = append(pairs, GetPairByID(key))
	}

	return pairs
}

func DeletePair(pair *Pair) {
	pairMap.Remove(fmt.Sprintf("%d", pair.ID))
}
