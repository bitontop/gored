package utils

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
	cmap "github.com/orcaman/concurrent-map"
)

func GetExchangeDataFromJSON(datapath string, exName exchange.ExchangeName) *ExchangeData {
	fileName := fmt.Sprintf("%s/%s.json", datapath, exName)
	exchangeData := &ExchangeData{
		CoinConstraint: cmap.New(),
		PairConstraint: cmap.New(),
	}
	var err error
	data := []byte{}
	if datapath[0:4] == "http" {
		data = []byte(exchange.HttpGetRequest(fileName, nil))
	} else {
		if data, err = ioutil.ReadFile(fileName); err != nil {
			log.Printf("Read %s Failed: %v", fileName, err)
			return nil
		}
	}

	jsonData := &JsonData{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		log.Printf("%s Json Unmarshal Err: %v %s", fileName, err, data)
		return nil
	}

	for _, cc := range jsonData.CoinConstraint {
		key := fmt.Sprintf("%d", cc.CoinID)
		cc.Coin = coin.GetCoinByID(cc.CoinID)
		if cc.Coin != nil {
			exchangeData.CoinConstraint.Set(key, cc)
		}
	}

	for _, pc := range jsonData.PairConstraint {
		key := fmt.Sprintf("%d", pc.PairID)
		pc.Pair = pair.GetPairByID(pc.PairID)
		if pc.Pair != nil {
			exchangeData.PairConstraint.Set(key, pc)
		}
	}

	return exchangeData
}

func GetCommonDataFromJSON(datapath string) {
	fileName := fmt.Sprintf("%s/common.json", datapath)
	var err error
	data := []byte{}
	if datapath[0:4] == "http" {
		data = []byte(exchange.HttpGetRequest(fileName, nil))
	} else {
		if data, err = ioutil.ReadFile(fileName); err != nil {
			log.Printf("Read %s Failed: %v", fileName, err)
			return
		}
	}

	commonData := &CommonData{}
	if err := json.Unmarshal(data, &commonData); err != nil {
		log.Printf("%s Json Unmarshal Err: %v %s", datapath, err, data)
		return
	}

	for _, c := range commonData.Coins {
		coin.AddCoin(c)
	}

	for _, p := range commonData.Pairs {
		pair.SetPair(p.ID, p.Base, p.Target)
	}
}
