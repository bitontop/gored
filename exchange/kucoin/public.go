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

	}
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

func (e *Kucoin) doTickerPrice(operation *exchange.PublicOperation) error {
	jsonResponse := JsonResponse{}
	tickerPrice := TickerPrice{}

	get := &utils.HttpGet{
		URI:       fmt.Sprintf("%s/api/v1/market/allTickers", API_URL), // TODO sandbox
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
			log.Printf("%s doTickerPrice parse Err: %v %v", e.GetName(), err, tp.AveragePrice)
			operation.Error = err
			return err
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

	get := &utils.HttpGet{
		URI: fmt.Sprintf("%s/api/v1/market/histories?symbol=%s", API_URL, symbol),
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

	get := &utils.HttpGet{
		URI:       fmt.Sprintf("%s/api/v1/market/orderbook/level2?symbol=%s", API_URL, symbol),
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
