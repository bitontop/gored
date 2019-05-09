package main

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"log"
	"os"
	"time"

	"../gored/exchange/bitmax"
	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/exchange/binance"
	"github.com/bitontop/gored/exchange/bittrex"
	"github.com/bitontop/gored/exchange/coinex"
	"github.com/bitontop/gored/exchange/kucoin"
	"github.com/bitontop/gored/exchange/stex"
	"github.com/bitontop/gored/pair"
	"github.com/bitontop/gored/test/conf"
	"github.com/bitontop/gored/utils"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	exMan := exchange.CreateExchangeManager()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		/* case "export":
		Init(exchange.EXCHANGE_API, "")
		utils.ConvertBaseDataToJson("./data")
		for _, ex := range exMan.GetExchanges() {
			utils.ConvertExchangeDataToJson("./data", ex)
		}
		break */
		case "json":
			Init(exchange.JSON_FILE, "./data")
			for _, ex := range exMan.GetExchanges() {
				for _, coin := range ex.GetCoins() {
					log.Printf("%s Coin %+v", ex.GetName(), coin)
				}
				for _, pair := range ex.GetPairs() {
					log.Printf("%s Pair %+v", ex.GetName(), pair)
				}
			}
			break
		case "renew":
			Init(exchange.JSON_FILE, "./data")
			updateConfig := &exchange.Update{
				ExNames: exMan.GetSupportExchanges(),
				Method:  exchange.TIME_TIGGER,
				Time:    10 * time.Second,
			}
			exMan.UpdateExData(updateConfig)
			break
		}
	}
}

func Init(source exchange.DataSource, sourceURI string) {
	coin.Init()
	pair.Init()
	if source == exchange.JSON_FILE {
		utils.GetCommonDataFromJSON(sourceURI)
	}
	config := &exchange.Config{}
	config.Source = source
	config.SourceURI = sourceURI

	InitBinance(config)
	InitBittrex(config)
	InitCoinex(config)
	InitStex(config)
	InitKucoin(config)
	InitBitmax(config)
}

func InitBinance(config *exchange.Config) {
	conf.Exchange(exchange.BINANCE, config)
	ex := binance.CreateBinance(config)
	log.Printf("Initial [ %12v ] ", ex.GetName())

	exMan := exchange.CreateExchangeManager()
	exMan.Add(ex)
}

func InitBittrex(config *exchange.Config) {
	conf.Exchange(exchange.BITTREX, config)
	ex := bittrex.CreateBittrex(config)
	log.Printf("Initial [ %12v ] ", ex.GetName())

	exMan := exchange.CreateExchangeManager()
	exMan.Add(ex)
}

func InitCoinex(config *exchange.Config) {
	conf.Exchange(exchange.COINEX, config)
	ex := coinex.CreateCoinex(config)
	log.Printf("Initial [ %12v ] ", ex.GetName())

	exMan := exchange.CreateExchangeManager()
	exMan.Add(ex)
}

func InitStex(config *exchange.Config) {
	conf.Exchange(exchange.STEX, config)
	ex := stex.CreateStex(config)
	log.Printf("Initial [ %12v ] ", ex.GetName())

	exMan := exchange.CreateExchangeManager()
	exMan.Add(ex)
}

func InitKucoin(config *exchange.Config) {
	conf.Exchange(exchange.KUCOIN, config)
	ex := kucoin.CreateKucoin(config)
	log.Printf("Initial [ %12v ] ", ex.GetName())

	exMan := exchange.CreateExchangeManager()
	exMan.Add(ex)
}

func InitBitmax(config *exchange.Config) {
	conf.Exchange(exchange.BITMAX, config)
	ex := bitmax.CreateBitmax(config)
	log.Printf("Initial [ %12v ] ", ex.GetName())

	exMan := exchange.CreateExchangeManager()
	exMan.Add(ex)
}
