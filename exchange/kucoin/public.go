package kucoin

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/utils"
)

func (e *Kucoin) LoadPublicData(operation *exchange.PublicOperation) error {
	switch operation.Type {
	case exchange.TradeHistory:
		return e.doTradeHistory(operation)
	case exchange.Orderbook:
		switch operation.Wallet {
		case exchange.SpotWallet:
			return e.doSpotOrderBook(operation)
		}
	case exchange.GetTickerPrice:
		switch operation.Wallet {
		case exchange.SpotWallet:
			return e.doTickerPrice(operation)
		}
	case exchange.KLine:
		switch operation.Wallet {
		case exchange.SpotWallet:
			return e.doSpotKline(operation)
		}

	}
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

// interval options: 1min, 3min, 5min, 15min, 30min, 1hour, 2hour, 4hour, 6hour, 8hour, 12hour, 1day, 1week
func (e *Kucoin) doSpotKline(operation *exchange.PublicOperation) error {
	interval := "5min"
	if operation.KlineInterval != "" {
		switch operation.KlineInterval {
		case "1min":
			interval = "1min"
		case "3min":
			interval = "3min"
		case "5min":
			interval = "5min"
		case "15min":
			interval = "15min"
		case "30min":
			interval = "30min"
		case "1hour":
			interval = "1hour"
		case "2hour":
			interval = "2hour"
		case "4hour":
			interval = "4hour"
		case "6hour":
			interval = "6hour"
		case "8hour":
			interval = "8hour"
		case "12hour":
			interval = "12hour"
		case "1day":
			interval = "1day"
		case "1week":
			interval = "1week"
		}
	}

	baseURL := API_URL
	if e.isSandBox() {
		baseURL = SANDBOX_API_URL
	}
	get := &utils.HttpGet{
		URI: fmt.Sprintf("%s/api/v1/market/candles?symbol=%v&type=%v", baseURL,
			e.GetSymbolByPair(operation.Pair), // ETH-BTC
			interval,
		),
		Proxy:     operation.Proxy,
		DebugMode: operation.DebugMode,
	}
	if operation.KlineStartTime != 0 {
		startTime := operation.KlineStartTime
		for startTime > 9999999999 {
			startTime = startTime / 10
		}
		get.URI += fmt.Sprintf("&startAt=%v", startTime)
	}
	if operation.KlineEndTime != 0 {
		endTime := operation.KlineEndTime
		for endTime > 9999999999 {
			endTime = endTime / 10
		}
		get.URI += fmt.Sprintf("&endAt=%v", endTime)
	}

	if err := utils.HttpGetRequest(get); err != nil {
		operation.Error = err
		return operation.Error
	}
	if operation.DebugMode {
		operation.RequestURI = get.URI
		operation.CallResponce = string(get.ResponseBody)
	}

	rawKline := KLine{}

	jsonKline := get.ResponseBody
	if err := json.Unmarshal([]byte(jsonKline), &rawKline); err != nil {
		return fmt.Errorf("%s Get doSpotKline Json Unmarshal Err: %v %v", e.GetName(), err, jsonKline)
	} else if rawKline.Code != "200000" {
		return fmt.Errorf("%s Get doSpotKline Failed: %v", e.GetName(), jsonKline)
	}

	operation.Kline = []*exchange.KlineDetail{}
	for _, k := range rawKline.Data {
		openTime, err := strconv.ParseFloat(k[0], 64)
		if err != nil {
			log.Printf("%s openTime parse Err: %v %v", e.GetName(), err, k[0])
			operation.Error = err
			return err
		}
		open, err := strconv.ParseFloat(k[1], 64)
		if err != nil {
			log.Printf("%s open parse Err: %v %v", e.GetName(), err, k[1])
			operation.Error = err
			return err
		}
		close, err := strconv.ParseFloat(k[2], 64)
		if err != nil {
			log.Printf("%s close parse Err: %v %v", e.GetName(), err, k[2])
			operation.Error = err
			return err
		}
		high, err := strconv.ParseFloat(k[3], 64)
		if err != nil {
			log.Printf("%s high parse Err: %v %v", e.GetName(), err, k[3])
			operation.Error = err
			return err
		}
		low, err := strconv.ParseFloat(k[4], 64)
		if err != nil {
			log.Printf("%s low parse Err: %v %v", e.GetName(), err, k[4])
			operation.Error = err
			return err
		}
		volume, err := strconv.ParseFloat(k[6], 64) // k[5] Transaction amount, k[6] Transaction volume.
		if err != nil {
			log.Printf("%s volume parse Err: %v %v", e.GetName(), err, k[6])
			operation.Error = err
			return err
		}

		detail := &exchange.KlineDetail{
			Exchange: e.GetName(),
			Pair:     operation.Pair.Name,
			OpenTime: openTime,
			Open:     open,
			High:     high,
			Low:      low,
			Close:    close,
			Volume:   volume,
		}

		operation.Kline = append(operation.Kline, detail)
	}

	return nil
}

func (e *Kucoin) doTickerPrice(operation *exchange.PublicOperation) error {
	jsonResponse := JsonResponse{}
	tickerPrice := TickerPrice{}

	baseURL := API_URL
	if e.isSandBox() {
		baseURL = SANDBOX_API_URL
	}
	get := &utils.HttpGet{
		URI:       fmt.Sprintf("%s/api/v1/market/allTickers", baseURL),
		Proxy:     operation.Proxy,
		DebugMode: operation.DebugMode,
	}
	if err := utils.HttpGetRequest(get); err != nil {
		operation.Error = err
		return operation.Error
	}
	if operation.DebugMode {
		operation.RequestURI = get.URI
		operation.CallResponce = string(get.ResponseBody)
	}

	jsonTickerPrice := get.ResponseBody
	if err := json.Unmarshal([]byte(jsonTickerPrice), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get tickerPrice Json Unmarshal Err: %v %v", e.GetName(), err, jsonTickerPrice)
	} else if jsonResponse.Code != "200000" {
		return fmt.Errorf("%s Get Pairs Failed: %s %v", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &tickerPrice); err != nil {
		return fmt.Errorf("%s Get tickerPrice Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	operation.TickerPrice = []*exchange.TickerPriceDetail{}
	for _, tp := range tickerPrice.Ticker {
		p := e.GetPairBySymbol(tp.Symbol)
		price, err := strconv.ParseFloat(tp.AveragePrice, 64)
		if err != nil {
			log.Printf("%s doTickerPrice parse Err: %v, %v:%v", e.GetName(), err, tp.Symbol, tp.AveragePrice)
			operation.Error = err
			continue
			// return err
		}

		if p == nil {
			if operation.DebugMode {
				log.Printf("doTickerPrice got nil pair for symbol: %v", tp.Symbol)
			}
			continue
		} else if p.Name == "" {
			continue
		}

		tpd := &exchange.TickerPriceDetail{
			Pair:  p,
			Price: price,
		}

		operation.TickerPrice = append(operation.TickerPrice, tpd)
	}

	return nil
}

func (e *Kucoin) doTradeHistory(operation *exchange.PublicOperation) error {
	symbol := e.GetSymbolByPair(operation.Pair)

	baseURL := API_URL
	if e.isSandBox() {
		baseURL = SANDBOX_API_URL
	}
	get := &utils.HttpGet{
		URI: fmt.Sprintf("%s/api/v1/market/histories?symbol=%s", baseURL, symbol),
	}

	err := utils.HttpGetRequest(get)

	if err != nil {
		return err

	} else {
		jsonResponse := JsonResponse{}
		tradeHistory := &TradeHistory{}
		if err := json.Unmarshal(get.ResponseBody, &jsonResponse); err != nil {
			return err
		}

		if err := json.Unmarshal(jsonResponse.Data, &tradeHistory); err != nil {
			return err
		}

		operation.TradeHistory = []*exchange.TradeDetail{}
		for _, d := range *tradeHistory {
			td := &exchange.TradeDetail{}

			td.ID = d.Sequence
			if d.Side == "buy" {
				td.Direction = exchange.Buy
			} else if d.Side == "sell" {
				td.Direction = exchange.Sell
			}

			td.Quantity, err = strconv.ParseFloat(d.Size, 64)
			td.Rate, err = strconv.ParseFloat(d.Price, 64)

			td.TimeStamp = d.Time / 1e6

			operation.TradeHistory = append(operation.TradeHistory, td)
		}
	}

	return nil
}

func (e *Kucoin) doSpotOrderBook(op *exchange.PublicOperation) error {
	jsonResponse := &JsonResponse{}
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(op.Pair)

	maker := &exchange.Maker{
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	baseURL := API_URL
	if e.isSandBox() {
		baseURL = SANDBOX_API_URL
	}
	get := &utils.HttpGet{
		URI:       fmt.Sprintf("%s/api/v1/market/orderbook/level2?symbol=%s", baseURL, symbol),
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
	} else if jsonResponse.Code != "200000" {
		return fmt.Errorf("%s Get Pairs Failed: %s %v", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderBook); err != nil {
		return fmt.Errorf("%s Get Orderbook Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	sequence, err := strconv.Atoi(orderBook.Sequence)
	if err != nil {
		return fmt.Errorf("Kucoin orderbook sequence Atoi err: %v", err)
	}
	maker.LastUpdateID = int64(sequence)
	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)

	for _, bid := range orderBook.Bids {
		buydata := exchange.Order{}

		buydata.Rate, err = strconv.ParseFloat(bid[0], 64)
		if err != nil {
			return fmt.Errorf("%s OrderBook strconv.ParseFloat Rate error:%v", e.GetName(), err)
		}
		buydata.Quantity, err = strconv.ParseFloat(bid[1], 64)
		if err != nil {
			return fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
		}
		maker.Bids = append(maker.Bids, buydata)
	}

	for i := len(orderBook.Asks) - 1; i >= 0; i-- {
		selldata := exchange.Order{}
		selldata.Rate, err = strconv.ParseFloat(orderBook.Asks[i][0], 64)
		if err != nil {
			return fmt.Errorf("%s OrderBook strconv.ParseFloat  Rate error:%v", e.GetName(), err)
		}
		selldata.Quantity, err = strconv.ParseFloat(orderBook.Asks[i][1], 64)
		if err != nil {
			return fmt.Errorf("%s OrderBook strconv.ParseFloat  Quantity error:%v", e.GetName(), err)
		}
		maker.Asks = append(maker.Asks, selldata)
	}

	op.Maker = maker
	return nil
}
