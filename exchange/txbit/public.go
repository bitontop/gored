package txbit

// Contributor 2015-2020 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
	"fmt"
	"log"

	exchange "github.com/bitontop/gored/exchange"
	utils "github.com/bitontop/gored/utils"
)

/*************** PUBLIC  API ***************/
func (e *Txbit) LoadPublicData(operation *exchange.PublicOperation) error {
	switch operation.Type {
	case exchange.TradeHistory:
		return e.doTradeHistory(operation)
	}
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

func (e *Txbit) doTradeHistory(operation *exchange.PublicOperation) error {
	symbol := e.GetSymbolByPair(operation.Pair)
	strRequestUrl := fmt.Sprintf("/public/getmarkethistory?market=%v", symbol)
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

		tradeHistory := TradeHistory{}
		if err := json.Unmarshal(get.ResponseBody, &tradeHistory); err != nil {
			operation.Error = err
			return err
		} else if !tradeHistory.Success {
			operation.Error = err
			return err
			// log.Printf("%+v ", tradeHistory)
		}

		operation.TradeHistory = []*exchange.TradeDetail{}
		for _, trade := range tradeHistory.Result {

			td := &exchange.TradeDetail{
				ID:        fmt.Sprintf("%v", trade.ID),
				Quantity:  trade.Quantity,
				TimeStamp: trade.TimeStamp.UnixNano() / 1e6,
				Rate:      trade.Price,
			}
			if trade.OrderType == "BUY" {
				td.Direction = exchange.Buy
			} else if trade.OrderType == "SELL" {
				td.Direction = exchange.Sell
			}

			operation.TradeHistory = append(operation.TradeHistory, td)
		}
	}

	return nil
}
