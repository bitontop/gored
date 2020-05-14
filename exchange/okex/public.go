package okex

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/utils"
)

func (e *Okex) LoadPublicData(operation *exchange.PublicOperation) error {
	switch operation.Type {

	case exchange.TradeHistory:
		return e.doTradeHistory(operation)
	case exchange.Orderbook:
		if operation.Wallet == exchange.SpotWallet {
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

func (e *Okex) doTickerPrice(operation *exchange.PublicOperation) error {
	tickerPrice := TickerPrice{}

	get := &utils.HttpGet{
		URI:       fmt.Sprintf("%s/api/spot/v3/instruments/ticker", API_URL),
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
	if err := json.Unmarshal([]byte(jsonTickerPrice), &tickerPrice); err != nil {
		operation.Error = fmt.Errorf("%s doTickerPrice json Unmarshal error: %v %v", e.GetName(), err, string(jsonTickerPrice))
		return operation.Error
	} else if len(tickerPrice) == 0 {
		operation.Error = fmt.Errorf("%s doTickerPrice failed, got empty ticker: %v", e.GetName(), string(jsonTickerPrice))
		return operation.Error
	}

	operation.TickerPrice = []*exchange.TickerPriceDetail{}
	for _, tp := range tickerPrice {
		p := e.GetPairBySymbol(tp.InstrumentID)
		if p == nil {
			if operation.DebugMode {
				log.Printf("doTickerPrice got nil pair for symbol: %v", tp.InstrumentID)
			}
			continue
		} else if p.Name == "" {
			continue
		}

		bid, err := strconv.ParseFloat(tp.BestBid, 64)
		if err != nil {
			log.Printf("%s doTickerPrice parse Err: %v %v", e.GetName(), err, tp.BestBid)
			operation.Error = err
			return err
		}
		ask, err := strconv.ParseFloat(tp.BestAsk, 64)
		if err != nil {
			log.Printf("%s doTickerPrice parse Err: %v %v", e.GetName(), err, tp.BestAsk)
			operation.Error = err
			return err
		}

		tpd := &exchange.TickerPriceDetail{
			Pair:  p,
			Price: (bid + ask) / 2,
		}

		operation.TickerPrice = append(operation.TickerPrice, tpd)
	}

	return nil
}

func (e *Okex) doSpotOrderBook(op *exchange.PublicOperation) error {
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(op.Pair)

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
		URI:       fmt.Sprintf("%s/api/spot/v3/instruments/%s/book", API_URL, symbol),
		Proxy:     op.Proxy,
		DebugMode: op.DebugMode,
	}
	if err := utils.HttpGetRequest(get); err != nil {
		op.Error = err
		return op.Error
	}

	jsonOrderbook := get.ResponseBody
	//! ###

	if err := json.Unmarshal([]byte(jsonOrderbook), &orderBook); err != nil {
		return fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %s %s", e.GetName(), err, jsonOrderbook)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	for _, bid := range orderBook.Bids {
		var buydata exchange.Order
		buydata.Rate, _ = strconv.ParseFloat(bid[0], 64)
		buydata.Quantity, _ = strconv.ParseFloat(bid[1], 64)

		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Asks {
		var selldata exchange.Order
		selldata.Rate, _ = strconv.ParseFloat(ask[0], 64)
		selldata.Quantity, _ = strconv.ParseFloat(ask[1], 64)

		maker.Asks = append(maker.Asks, selldata)
	}

	op.Maker = maker
	return nil
}

func (e *Okex) doTradeHistory(operation *exchange.PublicOperation) error {
	symbol := e.GetSymbolByPair(operation.Pair)

	get := &utils.HttpGet{
		URI: fmt.Sprintf("%s/api/spot/v3/instruments/%s/trades", API_URL, symbol),
	}

	err := utils.HttpGetRequest(get)

	if err != nil {
		return err

	} else {
		tradeHistory := TradeHistory{}
		if err := json.Unmarshal(get.ResponseBody, &tradeHistory); err != nil {
			return err
		}

		operation.TradeHistory = []*exchange.TradeDetail{}
		for i := len(tradeHistory) - 1; i > 0; i-- {
			// for _, d := range *tradeHistory {
			d := tradeHistory[i]
			td := &exchange.TradeDetail{}

			td.ID = d.TradeID
			if d.Side == "buy" {
				td.Direction = exchange.Buy
			} else if d.Side == "sell" {
				td.Direction = exchange.Sell
			}

			td.Quantity, err = strconv.ParseFloat(d.Size, 64)
			td.Rate, err = strconv.ParseFloat(d.Price, 64)

			td.TimeStamp = d.Timestamp.UnixNano() / 1e6

			operation.TradeHistory = append(operation.TradeHistory, td)
		}
	}

	return nil
}
