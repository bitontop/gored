package binance

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
func (e *Binance) LoadPublicData(operation *exchange.PublicOperation) error {
	switch operation.Type {
	case exchange.TradeHistory:
		return e.doTradeHistory(operation)
	case exchange.Orderbook:
		switch operation.Wallet {
		case exchange.SpotWallet:
			return e.doSpotOrderBook(operation)
		case exchange.ContractWallet:
			return e.doContractOrderBook(operation)
		}
	case exchange.KLine:
		switch operation.Wallet {
		case exchange.ContractWallet:
			return e.doContractKline(operation)
		case exchange.SpotWallet:
			return e.doSpotKline(operation)
		}
	}
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

// interval options: 1m, 3m, 5m, 15m, 30m, 1h, 2h, 4h, 6h, 8h, 12h, 1d, 3d, 1w, 1M
func (e *Binance) doSpotKline(operation *exchange.PublicOperation) error {
	interval := "5m"
	if operation.KlineInterval != "" {
		interval = operation.KlineInterval
	}

	get := &utils.HttpGet{
		URI: fmt.Sprintf("https://api.binance.com/api/v3/klines?symbol=%v&interval=%v&limit=1000", // 1500478320000
			e.GetSymbolByPair(operation.Pair), // BTCUSDT
			interval,
		),
		Proxy: operation.Proxy,
	}

	if operation.KlineStartTime != 0 {
		get.URI += fmt.Sprintf("&startTime=%v", operation.KlineStartTime)
	}

	err := utils.HttpGetRequest(get)

	if err != nil {
		log.Printf("%+v", err)
		operation.Error = err
		return err

	}

	if operation.DebugMode {
		operation.RequestURI = get.URI
		operation.CallResponce = string(get.ResponseBody)
	}

	var rawKline [][]interface{}
	if err := json.Unmarshal(get.ResponseBody, &rawKline); err != nil {
		operation.Error = fmt.Errorf("%s doContractKline Json Unmarshal Err: %v %v", e.GetName(), err, string(get.ResponseBody))
		return operation.Error
	}

	operation.Kline = []*exchange.KlineDetail{}
	for _, k := range rawKline {
		open, err := strconv.ParseFloat(k[1].(string), 64)
		if err != nil {
			log.Printf("%s open parse Err: %v %v", e.GetName(), err, k[1])
			operation.Error = err
			return err
		}
		high, err := strconv.ParseFloat(k[2].(string), 64)
		if err != nil {
			log.Printf("%s high parse Err: %v %v", e.GetName(), err, k[2])
			operation.Error = err
			return err
		}
		low, err := strconv.ParseFloat(k[3].(string), 64)
		if err != nil {
			log.Printf("%s low parse Err: %v %v", e.GetName(), err, k[3])
			operation.Error = err
			return err
		}
		close, err := strconv.ParseFloat(k[4].(string), 64)
		if err != nil {
			log.Printf("%s close parse Err: %v %v", e.GetName(), err, k[4])
			operation.Error = err
			return err
		}
		volume, err := strconv.ParseFloat(k[5].(string), 64)
		if err != nil {
			log.Printf("%s volume parse Err: %v %v", e.GetName(), err, k[5])
			operation.Error = err
			return err
		}
		quoteAssetVolume, err := strconv.ParseFloat(k[7].(string), 64)
		if err != nil {
			log.Printf("%s quoteAssetVolume parse Err: %v %v", e.GetName(), err, k[7])
			operation.Error = err
			return err
		}
		takerBuyBaseVolume, err := strconv.ParseFloat(k[9].(string), 64)
		if err != nil {
			log.Printf("%s takerBuyBaseVolume parse Err: %v %v", e.GetName(), err, k[9])
			operation.Error = err
			return err
		}
		takerBuyQuoteVolume, err := strconv.ParseFloat(k[10].(string), 64)
		if err != nil {
			log.Printf("%s takerBuyQuoteVolume parse Err: %v %v", e.GetName(), err, k[10])
			operation.Error = err
			return err
		}

		detail := &exchange.KlineDetail{
			OpenTime:            k[0].(float64),
			Open:                open,
			High:                high,
			Low:                 low,
			Close:               close,
			Volume:              volume,
			CloseTime:           k[6].(float64),
			QuoteAssetVolume:    quoteAssetVolume,
			TradesCount:         k[8].(float64),
			TakerBuyBaseVolume:  takerBuyBaseVolume,
			TakerBuyQuoteVolume: takerBuyQuoteVolume,
		}

		operation.Kline = append(operation.Kline, detail)
	}

	return nil
}

// interval options: 1m, 3m, 5m, 15m, 30m, 1h, 2h, 4h, 6h, 8h, 12h, 1d, 3d, 1w, 1M
func (e *Binance) doContractKline(operation *exchange.PublicOperation) error {
	interval := "5m"
	if operation.KlineInterval != "" {
		interval = operation.KlineInterval
	}

	get := &utils.HttpGet{
		URI: fmt.Sprintf("https://fapi.binance.com/fapi/v1/klines?symbol=%v&interval=%v&limit=1500",
			e.GetSymbolByPair(operation.Pair), // BTCUSDT
			interval,
		),
		Proxy: operation.Proxy,
	}

	if operation.KlineStartTime != 0 {
		get.URI += fmt.Sprintf("&startTime=%v", operation.KlineStartTime)
	}

	err := utils.HttpGetRequest(get)

	if err != nil {
		log.Printf("%+v", err)
		operation.Error = err
		return err

	}

	if operation.DebugMode {
		operation.RequestURI = get.URI
		operation.CallResponce = string(get.ResponseBody)
	}

	var rawKline [][]interface{}
	if err := json.Unmarshal(get.ResponseBody, &rawKline); err != nil {
		operation.Error = fmt.Errorf("%s doContractKline Json Unmarshal Err: %v %v", e.GetName(), err, string(get.ResponseBody))
		return operation.Error
	}

	operation.Kline = []*exchange.KlineDetail{}
	for _, k := range rawKline {
		open, err := strconv.ParseFloat(k[1].(string), 64)
		if err != nil {
			log.Printf("%s open parse Err: %v %v", e.GetName(), err, k[1])
			operation.Error = err
			return err
		}
		high, err := strconv.ParseFloat(k[2].(string), 64)
		if err != nil {
			log.Printf("%s high parse Err: %v %v", e.GetName(), err, k[2])
			operation.Error = err
			return err
		}
		low, err := strconv.ParseFloat(k[3].(string), 64)
		if err != nil {
			log.Printf("%s low parse Err: %v %v", e.GetName(), err, k[3])
			operation.Error = err
			return err
		}
		close, err := strconv.ParseFloat(k[4].(string), 64)
		if err != nil {
			log.Printf("%s close parse Err: %v %v", e.GetName(), err, k[4])
			operation.Error = err
			return err
		}
		volume, err := strconv.ParseFloat(k[5].(string), 64)
		if err != nil {
			log.Printf("%s volume parse Err: %v %v", e.GetName(), err, k[5])
			operation.Error = err
			return err
		}
		quoteAssetVolume, err := strconv.ParseFloat(k[7].(string), 64)
		if err != nil {
			log.Printf("%s quoteAssetVolume parse Err: %v %v", e.GetName(), err, k[7])
			operation.Error = err
			return err
		}
		takerBuyBaseVolume, err := strconv.ParseFloat(k[9].(string), 64)
		if err != nil {
			log.Printf("%s takerBuyBaseVolume parse Err: %v %v", e.GetName(), err, k[9])
			operation.Error = err
			return err
		}
		takerBuyQuoteVolume, err := strconv.ParseFloat(k[10].(string), 64)
		if err != nil {
			log.Printf("%s takerBuyQuoteVolume parse Err: %v %v", e.GetName(), err, k[10])
			operation.Error = err
			return err
		}

		detail := &exchange.KlineDetail{
			OpenTime:            k[0].(float64),
			Open:                open,
			High:                high,
			Low:                 low,
			Close:               close,
			Volume:              volume,
			CloseTime:           k[6].(float64),
			QuoteAssetVolume:    quoteAssetVolume,
			TradesCount:         k[8].(float64),
			TakerBuyBaseVolume:  takerBuyBaseVolume,
			TakerBuyQuoteVolume: takerBuyQuoteVolume,
		}

		operation.Kline = append(operation.Kline, detail)
	}

	return nil
}

func (e *Binance) doTradeHistory(operation *exchange.PublicOperation) error {

	get := &utils.HttpGet{
		URI: fmt.Sprintf("https://api.binance.com/api/v3/trades?symbol=%s&limit=%d",
			e.GetSymbolByPair(operation.Pair),
			1000, //TRADE_HISTORY_MAX_LIMIT,
		),
		Proxy: operation.Proxy,
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
				ID:        fmt.Sprintf("%v", trade.ID),
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

func (e *Binance) doContractOrderBook(operation *exchange.PublicOperation) error {
	orderbook := ContractOrderBook{}
	symbol := e.GetSymbolByPair(operation.Pair)

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol

	strRequestUrl := "/fapi/v1/depth"
	strUrl := CONTRACT_URL + strRequestUrl

	operation.Maker = &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbookReturn := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbookReturn), &orderbook); err != nil {
		return fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbookReturn)
	}

	operation.Maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	var err error
	for _, bid := range orderbook.Bids {
		buydata := exchange.Order{}
		buydata.Quantity, err = strconv.ParseFloat(bid[1], 64)
		if err != nil {
			return fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
		}

		buydata.Rate, err = strconv.ParseFloat(bid[0], 64)
		if err != nil {
			return fmt.Errorf("%s OrderBook strconv.ParseFloat Rate error:%v", e.GetName(), err)
		}
		operation.Maker.Bids = append(operation.Maker.Bids, buydata)
	}
	for _, ask := range orderbook.Asks {
		selldata := exchange.Order{}
		selldata.Quantity, err = strconv.ParseFloat(ask[1], 64)
		if err != nil {
			return fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
		}

		selldata.Rate, err = strconv.ParseFloat(ask[0], 64)
		if err != nil {
			return fmt.Errorf("%s OrderBook strconv.ParseFloat Rate error:%v", e.GetName(), err)
		}
		operation.Maker.Asks = append(operation.Maker.Asks, selldata)
	}
	return nil
}

func (e *Binance) doSpotOrderBook(op *exchange.PublicOperation) error {
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(op.Pair)

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol
	mapParams["limit"] = "100"

	maker := &exchange.Maker{
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	get := &utils.HttpGet{
		URI:       fmt.Sprintf("%s/api/v1/depth?%s", API_URL, exchange.Map2UrlQuery(mapParams)),
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
	maker.LastUpdateID = int64(orderBook.LastUpdateID)

	var err error
	for _, bid := range orderBook.Bids {
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

	for _, ask := range orderBook.Asks {
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
	op.Maker = maker
	return nil
}
