package bithumb

// Contributor 2015-2020 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	exchange "github.com/bitontop/gored/exchange"
	utils "github.com/bitontop/gored/utils"
)

/*************** PUBLIC  API ***************/
func (e *Bithumb) LoadPublicData(operation *exchange.PublicOperation) error {
	switch operation.Type {
	case exchange.TradeHistory:
		return e.doTradeHistory(operation)
	}
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

// timestamp only 10 digit precision
func (e *Bithumb) doTradeHistory(operation *exchange.PublicOperation) error {
	symbol := e.GetSymbolByPair(operation.Pair)
	strRequestUrl := fmt.Sprintf("/spot/trades?symbol=%v", symbol)
	strUrl := API_URL + strRequestUrl

	get := &utils.HttpGet{
		URI: strUrl,
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

		tradeHistory := TradeHistory{} //TradeHistory{}
		if err := json.Unmarshal(get.ResponseBody, &tradeHistory); err != nil {
			operation.Error = err
			return err
		} else if tradeHistory.Code != "0" {
			operation.Error = err
			return err
			// log.Printf("%+v ", tradeHistory)
		}

		operation.TradeHistory = []*exchange.TradeDetail{}
		// for _, trade := range tradeHistory.Data {
		for i := len(tradeHistory.Data) - 1; i > 0; i-- {
			trade := tradeHistory.Data[i]
			price, err := strconv.ParseFloat(trade.P, 64)
			if err != nil {
				log.Printf("%s price parse Err: %v %v", e.GetName(), err, trade.P)
				operation.Error = err
				return err
			}
			amount, err := strconv.ParseFloat(trade.V, 64)
			if err != nil {
				log.Printf("%s amount parse Err: %v %v", e.GetName(), err, trade.V)
				operation.Error = err
				return err
			}
			ts, err := strconv.ParseInt(trade.T, 10, 64)
			if err != nil {
				log.Printf("%s ts parse Err: %v %v", e.GetName(), err, trade.P)
				operation.Error = err
				return err
			}

			td := &exchange.TradeDetail{
				ID:        trade.V,
				Quantity:  amount,
				TimeStamp: ts * 1000, //trade.TimeStamp.UnixNano() / 1e6,
				Rate:      price,
			}
			if trade.S == "buy" {
				td.Direction = exchange.Buy
			} else if trade.S == "sell" {
				td.Direction = exchange.Sell
			}

			operation.TradeHistory = append(operation.TradeHistory, td)
		}
	}

	return nil
}
