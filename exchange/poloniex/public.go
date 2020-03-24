package poloniex

// Contributor 2015-2020 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	exchange "github.com/bitontop/gored/exchange"
	utils "github.com/bitontop/gored/utils"
)

/*************** PUBLIC  API ***************/
func (e *Poloniex) LoadPublicData(operation *exchange.PublicOperation) error {
	switch operation.Type {
	case exchange.TradeHistory:
		return e.doTradeHistory(operation)
	case exchange.Orderbook:
		switch operation.Wallet {
		case exchange.SpotWallet:
			return e.doSpotOrderBook(operation)
		}
	}
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

// timestamp 10 digit precision
func (e *Poloniex) doTradeHistory(operation *exchange.PublicOperation) error {
	symbol := e.GetSymbolByPair(operation.Pair)
	strRequestUrl := fmt.Sprintf("/public?command=returnTradeHistory&currencyPair=%v", symbol)
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
		} else if len(tradeHistory) == 0 {
			operation.Error = fmt.Errorf("Got Empty Trade History")
			return fmt.Errorf("Got Empty Trade History")
			// log.Printf("%+v ", tradeHistory)
		}

		operation.TradeHistory = []*exchange.TradeDetail{}
		// for _, trade := range tradeHistory {
		for i := len(tradeHistory) - 1; i > 0; i-- {
			trade := tradeHistory[i]
			price, err := strconv.ParseFloat(trade.Rate, 64)
			if err != nil {
				log.Printf("%s price parse Err: %v %v", e.GetName(), err, trade.Rate)
				operation.Error = err
				return err
			}
			amount, err := strconv.ParseFloat(trade.Amount, 64)
			if err != nil {
				log.Printf("%s amount parse Err: %v %v", e.GetName(), err, trade.Amount)
				operation.Error = err
				return err
			}

			layout := "2006-01-02 15:04:05"
			ts, _ := time.Parse(layout, trade.Date)

			td := &exchange.TradeDetail{
				ID:        fmt.Sprintf("%v", trade.TradeID),
				Quantity:  amount,
				TimeStamp: ts.Unix() * 1000,
				Rate:      price,
			}
			if trade.Type == "buy" {
				td.Direction = exchange.Buy
			} else if trade.Type == "sell" {
				td.Direction = exchange.Sell
			}

			operation.TradeHistory = append(operation.TradeHistory, td)
		}
	}

	return nil
}

func (e *Poloniex) doSpotOrderBook(op *exchange.PublicOperation) error {
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(op.Pair)

	mapParams := make(map[string]string)
	mapParams["command"] = "returnOrderBook"
	mapParams["currencyPair"] = symbol

	maker := &exchange.Maker{
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	get := &utils.HttpGet{
		URI:       fmt.Sprintf("%s/public?%s", API_URL, exchange.Map2UrlQuery(mapParams)),
		Proxy:     op.Proxy,
		DebugMode: op.DebugMode,
	}
	if err := utils.HttpGetRequest(get); err != nil {
		op.Error = err
		return op.Error
	}

	jsonOrderbook := get.ResponseBody
	if err := json.Unmarshal([]byte(jsonOrderbook), &orderBook); err != nil {
		return fmt.Errorf("%s OrderBook json Unmarshal error: %v %v", e.GetName(), err, jsonOrderbook)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)

	var err error
	for _, bid := range orderBook.Bids {
		var buydata exchange.Order

		//Modify according to type and structure
		buydata.Rate, err = strconv.ParseFloat(bid[0].(string), 64)
		if err != nil {
			return fmt.Errorf("Poloniex Bids Rate ParseFloat error:%v", err)
		}
		buydata.Quantity = bid[1].(float64)

		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Asks {
		var selldata exchange.Order

		//Modify according to type and structure
		selldata.Rate, err = strconv.ParseFloat(ask[0].(string), 64)
		if err != nil {
			return fmt.Errorf("Poloniex Asks Rate ParseFloat error:%v", err)
		}
		selldata.Quantity = ask[1].(float64)

		maker.Asks = append(maker.Asks, selldata)
	}
	maker.LastUpdateID = orderBook.Seq

	op.Maker = maker
	return nil
}
