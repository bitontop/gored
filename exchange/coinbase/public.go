package coinbase

// Contributor 2015-2020 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/bitontop/gored/coin"
	exchange "github.com/bitontop/gored/exchange"
	utils "github.com/bitontop/gored/utils"
)

/*************** PUBLIC  API ***************/
func (e *Coinbase) LoadPublicData(operation *exchange.PublicOperation) error {
	switch operation.Type {
	case exchange.GetCoin:
		return e.doGetCoin(operation)
	case exchange.GetPair:
		// return e.doGetPair
	case exchange.TradeHistory:
		return e.doTradeHistory(operation)

	}
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

func (e *Coinbase) doGetCoin(operation *exchange.PublicOperation) error {
	coinsData := CoinsData{}

	strUrl := "https://www.binance.com/assetWithdraw/getAllAsset.html"

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	}

	for _, data := range coinsData {
		c := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			c = coin.GetCoin(data.AssetCode)
			if c == nil {
				c = &coin.Coin{}
				c.Code = data.AssetCode
				c.Name = data.AssetName
				c.Website = data.URL
				c.Explorer = data.BlockURL
				coin.AddCoin(c)
			}
		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(data.AssetCode)
		}

		if c != nil {
			confirmation, _ := strconv.Atoi(data.ConfirmTimes)
			coinConstraint := &exchange.CoinConstraint{
				CoinID:       c.ID,
				Coin:         c,
				ExSymbol:     data.AssetCode,
				ChainType:    exchange.MAINNET,
				TxFee:        data.TransactionFee,
				Withdraw:     data.EnableWithdraw,
				Deposit:      data.EnableCharge,
				Confirmation: confirmation,
				Listed:       true,
			}

			e.SetCoinConstraint(coinConstraint)
		}
	}
	return nil
}

func (e *Coinbase) doTradeHistory(operation *exchange.PublicOperation) error {

	get := &utils.HttpGet{
		URI: fmt.Sprintf("https://api.binance.com/api/v3/trades?symbol=%s&limit=%d",
			e.GetSymbolByPair(operation.Pair),
			1000, //TRADE_HISTORY_MAX_LIMIT,
		),
	}

	err := utils.HttpGetRequest(get)

	if err != nil {
		log.Printf("%+v", err)
		operation.Error = err
		return err

	} else {
		// log.Printf("%+v  ERR:%+v", string(get.ResponseBody), err) // ======================
		if operation.DebugMode {
			operation.RequestURI = get.URI
			operation.CallResponce = string(get.ResponseBody)
		}

		tradeHistory := TradeHistory{}
		if err := json.Unmarshal(get.ResponseBody, &tradeHistory); err != nil {
			operation.Error = err
			return err
		} else {
			// log.Printf("%+v ", tradeHistory)
		}

		operation.TradeHistory = []*exchange.TradeDetail{}
		for _, trade := range tradeHistory {
			price, err := strconv.ParseFloat(trade.Price, 64)
			if err != nil {
				log.Printf("%s price parse Err: %v %v", e.GetName(), err, trade.Price)
				operation.Error = err
				return err
			}
			amount, err := strconv.ParseFloat(trade.Qty, 64)
			if err != nil {
				log.Printf("%s amount parse Err: %v %v", e.GetName(), err, trade.Qty)
				operation.Error = err
				return err
			}

			td := &exchange.TradeDetail{
				Quantity:  amount,
				TimeStamp: trade.Time,
				Rate:      price,
			}
			if trade.IsBuyerMaker {
				td.Direction = exchange.Buy
			} else if !trade.IsBuyerMaker {
				td.Direction = exchange.Sell
			}

			operation.TradeHistory = append(operation.TradeHistory, td)
		}
	}

	return nil
}
