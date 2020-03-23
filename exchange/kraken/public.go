package kraken

// Contributor 2015-2020 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	exchange "github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/utils"
)

/*************** PUBLIC  API ***************/
func (e *Kraken) LoadPublicData(operation *exchange.PublicOperation) error {
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

// no tradeID
func (e *Kraken) doTradeHistory(operation *exchange.PublicOperation) error {
	// symbol := e.GetSymbolByPair(operation.Pair)
	// strRequestUrl := fmt.Sprintf("/spot/trades?symbol=%v", symbol)
	// strUrl := API_URL + strRequestUrl

	// get := &utils.HttpGet{
	// 	URI: strUrl,
	// }

	// err := utils.HttpGetRequest(get)

	// if err != nil {
	// 	log.Printf("%+v", err)
	// 	operation.Error = err
	// 	return err

	// } else {
	// 	// log.Printf("%+v  ERR:%+v", string(get.ResponseBody), err) // ======================
	// 	if operation.DebugMode {
	// 		operation.RequestURI = get.URI
	// 		operation.CallResponce = string(get.ResponseBody)
	// 	}

	// 	tradeHistory := [][]Trade{} //TradeHistory{}
	// 	if err := json.Unmarshal(get.ResponseBody, &tradeHistory); err != nil {
	// 		operation.Error = err
	// 		return err
	// 	} else if len(tradeHistory) == 0 {
	// 		operation.Error = fmt.Errorf("Got Empty Trade History")
	// 		return fmt.Errorf("Got Empty Trade History")
	// 		// log.Printf("%+v ", tradeHistory)
	// 	}

	// 	operation.TradeHistory = []*exchange.TradeDetail{}
	// 	for _, trade := range tradeHistory {
	// 		price, err := strconv.ParseFloat(trade[0].Price, 64)
	// 		if err != nil {
	// 			log.Printf("%s price parse Err: %v %v", e.GetName(), err, trade[0].Price)
	// 			operation.Error = err
	// 			return err
	// 		}
	// 		amount, err := strconv.ParseFloat(trade[0].Volume, 64)
	// 		if err != nil {
	// 			log.Printf("%s amount parse Err: %v %v", e.GetName(), err, trade[0].Volume)
	// 			operation.Error = err
	// 			return err
	// 		}

	// 		td := &exchange.TradeDetail{
	// 			ID:        trade.V,
	// 			Quantity:  amount,
	// 			TimeStamp: trade.TimeStamp.UnixNano() / 1e6,
	// 			Rate:      price,
	// 		}
	// 		if trade.S == "buy" {
	// 			td.Direction = exchange.Buy
	// 		} else if trade.S == "sell" {
	// 			td.Direction = exchange.Sell
	// 		}

	// 		operation.TradeHistory = append(operation.TradeHistory, td)
	// 	}
	// }

	return nil
}

func (e *Kraken) doSpotOrderBook(op *exchange.PublicOperation) error {
	jsonResponse := &JsonResponse{}
	orderBook := make(map[string]*OrderBook)
	symbol := e.GetSymbolByPair(op.Pair)

	mapParams := make(map[string]string)
	mapParams["pair"] = symbol
	mapParams["count"] = "100"

	maker := &exchange.Maker{
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	get := &utils.HttpGet{
		URI:       fmt.Sprintf("%s/0/public/Depth?%s", API_URL, exchange.Map2UrlQuery(mapParams)),
		Proxy:     op.Proxy,
		DebugMode: op.DebugMode,
	}
	if err := utils.HttpGetRequest(get); err != nil {
		op.Error = err
		return op.Error
	}

	jsonOrderbook := get.ResponseBody
	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if len(jsonResponse.Error) != 0 {
		return fmt.Errorf("%s Get Orderbook Failed: %v", e.GetName(), jsonResponse.Error)
	}
	if err := json.Unmarshal(jsonResponse.Result, &orderBook); err != nil {
		return fmt.Errorf("%s Get Orderbook Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)

	var err error
	for _, book := range orderBook {
		for _, bid := range book.Bids {
			buydata := exchange.Order{}
			buydata.Quantity, err = strconv.ParseFloat(bid[1].(string), 64)
			if err != nil {
				return fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
			}

			buydata.Rate, err = strconv.ParseFloat(bid[0].(string), 64)
			if err != nil {
				return fmt.Errorf("%s OrderBook strconv.ParseFloat Rate error:%v", e.GetName(), err)
			}
			maker.Bids = append(maker.Bids, buydata)
		}

		for _, ask := range book.Asks {
			selldata := exchange.Order{}
			selldata.Quantity, err = strconv.ParseFloat(ask[1].(string), 64)
			if err != nil {
				return fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
			}

			selldata.Rate, err = strconv.ParseFloat(ask[0].(string), 64)
			if err != nil {
				return fmt.Errorf("%s OrderBook strconv.ParseFloat Rate error:%v", e.GetName(), err)
			}
			maker.Asks = append(maker.Asks, selldata)
		}
	}

	op.Maker = maker
	return nil
}
