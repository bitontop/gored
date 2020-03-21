package txbit

// Contributor 2015-2020 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	exchange "github.com/bitontop/gored/exchange"
	utils "github.com/bitontop/gored/utils"
)

/*************** PUBLIC  API ***************/
func (e *Txbit) LoadPublicData(operation *exchange.PublicOperation) error {
	switch operation.Type {
	case exchange.TradeHistory:
		return e.doTradeHistory(operation)

	case exchange.Orderbook:
		if operation.Wallet == exchange.SpotWallet {
			return e.doSpotOrderBook(operation)
		}

	}
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

func (e *Txbit) doSpotOrderBook(op *exchange.PublicOperation) error {
	pair := op.Pair

	// }func (e *Txbit) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {

	jsonResponse := &JsonResponse{}
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	mapParams := make(map[string]string)
	mapParams["market"] = symbol
	mapParams["type"] = "both"

	maker := &exchange.Maker{
		WorkerIP:        utils.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	// strRequestUrl := "/public/getorderbook"
	// strUrl := API_URL + strRequestUrl
	// jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)

	//!------
	get := &utils.HttpGet{
		URI:       fmt.Sprintf("%s/public/getorderbook?%s", API_URL, utils.Map2UrlQuery(mapParams)),
		Proxy:     op.Proxy,
		Timeout:   op.Timeout,
		DebugMode: op.DebugMode,
	}
	if err := utils.HttpGetRequest(get); err != nil {
		op.Error = err
		return op.Error
	}

	jsonOrderbook := get.ResponseBody
	//! ###

	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %s %s", e.GetName(), err, jsonOrderbook)
	} else if !jsonResponse.Success {
		return fmt.Errorf("%s Get Orderbook Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Result, &orderBook); err != nil {
		return fmt.Errorf("%s Get Orderbook Result Unmarshal Err: %s %s", e.GetName(), err, jsonResponse.Result)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	for _, bid := range orderBook.Buy {
		var buydata exchange.Order
		buydata.Rate = bid.Rate
		buydata.Quantity = bid.Quantity

		maker.Bids = append(maker.Bids, buydata)
	}
	for i := len(orderBook.Sell) - 1; i >= 0; i-- {
		var selldata exchange.Order

		selldata.Rate = orderBook.Sell[i].Rate
		selldata.Quantity = orderBook.Sell[i].Quantity

		maker.Asks = append(maker.Asks, selldata)
	}

	op.Maker = maker
	return nil
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
		// for _, trade := range tradeHistory.Result {
		for i := len(tradeHistory.Result) - 1; i > 0; i-- {
			trade := tradeHistory.Result[i]

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
