package bkex

// Contributor 2015-2020 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
	"fmt"
	"time"

	exchange "github.com/bitontop/gored/exchange"
	utils "github.com/bitontop/gored/utils"
)

/*************** PUBLIC  API ***************/
func (e *Bkex) LoadPublicData(operation *exchange.PublicOperation) error {
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

func (e *Bkex) doSpotOrderBook(op *exchange.PublicOperation) error {
	jsonResponse := &JsonResponse{}
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(op.Pair)

	mapParams := make(map[string]string)
	mapParams["pair"] = symbol

	maker := &exchange.Maker{
		WorkerIP:        utils.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	/******
		// strRequestUrl := fmt.Sprintf("/public/orderbook/%s", e.GetIDByPair(pair))
		// strUrl := API3_URL + strRequestUrl
		// jsonOrderbook := exchange.HttpGetRequest(strUrl, nil)
	********/
	//!------
	get := &utils.HttpGet{
		URI:       fmt.Sprintf("%s/v1/q/depth?%s", API_URL, utils.Map2UrlQuery(mapParams)),
		Proxy:     op.Proxy,
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
	} else if jsonResponse.Code != 0 {
		return fmt.Errorf("%s Get Orderbook Failed: %s", e.GetName(), jsonOrderbook)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderBook); err != nil {
		return fmt.Errorf("%s Get Orderbook Result Unmarshal Err: %s %s", e.GetName(), err, jsonResponse.Data)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	for _, bid := range orderBook.Bids {
		buydata := exchange.Order{}
		buydata.Quantity = bid.Amt
		buydata.Rate = bid.Price
		maker.Bids = append(maker.Bids, buydata)
	}

	for _, ask := range orderBook.Asks {
		selldata := exchange.Order{}
		selldata.Quantity = ask.Amt
		selldata.Rate = ask.Price
		maker.Asks = append(maker.Asks, selldata)
	}

	op.Maker = maker
	return nil
}

// No trade ID, using timestamp as id
func (e *Bkex) doTradeHistory(operation *exchange.PublicOperation) error {
	symbol := e.GetSymbolByPair(operation.Pair)

	get := &utils.HttpGet{
		URI:       fmt.Sprintf("%s/v1/q/deals?pair=%s", API_URL, symbol),
		Proxy:     operation.Proxy,
		DebugMode: operation.DebugMode,
	}

	err := utils.HttpGetRequest(get)

	if err != nil {
		// log.Printf("%+v", err)
		return err

	} else {
		// log.Printf("%+v  ERR:%+v", string(get.ResponseBody), err)
		tradeHistory := &TradeHistory{}
		if err := json.Unmarshal(get.ResponseBody, &tradeHistory); err != nil {
			return err
		} else {
			// log.Printf("%+v ", tradeHistory)
		}

		// log.Printf("%s", get.ResponseBody)

		operation.TradeHistory = []*exchange.TradeDetail{}
		for i := len(tradeHistory.Data) - 1; i > 0; i-- {
			d2 := tradeHistory.Data[i]
			// d2 := d1.Data[i]
			// log.Printf("d2:%+v", d2)
			td := &exchange.TradeDetail{
				ID:       fmt.Sprintf("%d", d2.CreatedTime),
				Quantity: d2.DealAmount,

				TimeStamp: d2.CreatedTime,
				Rate:      d2.Price,
			}

			if d2.TradeDealDirection == "B" {
				td.Direction = exchange.Buy
			} else if d2.TradeDealDirection == "S" {
				td.Direction = exchange.Sell
			}
			// log.Printf("d2: %+v ", d2)
			// log.Printf("TD: %+v ", td)

			operation.TradeHistory = append(operation.TradeHistory, td)

		}
	}

	return nil

}
