package stex

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/utils"
)

func (e *Stex) LoadPublicData(operation *exchange.PublicOperation) error {
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

func (e *Stex) doSpotOrderBook(op *exchange.PublicOperation) error {

	jsonResponse := JsonResponseV3{}
	orderBook := OrderBook{}
	maker := &exchange.Maker{
		// WorkerIP:        exchange.GetExternalIP(),
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
		URI:       fmt.Sprintf("%s/public/orderbook/%s", API3_URL, e.GetIDByPair(op.Pair)),
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
		log.Printf("%s Get Orderbook Err: %v, %s, Using webpage Orderbook...", e.GetName(), err, jsonOrderbook)
		op.Maker, op.Error = e.doWebOrderBook(op)
		return op.Error

	} else if !jsonResponse.Success {
		log.Printf("%s Get Orderbook fail: %s, Using webpage Orderbook...", e.GetName(), jsonOrderbook)
		op.Maker, op.Error = e.doWebOrderBook(op)
		return op.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderBook); err != nil {
		log.Printf("%s Get Orderbook Err: %v, %s, Using webpage Orderbook...", e.GetName(), err, jsonOrderbook)
		op.Maker, op.Error = e.doWebOrderBook(op)
		return op.Error
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	var err error
	for _, bid := range orderBook.Bid {
		var buydata exchange.Order

		//Modify according to type and structure
		buydata.Rate, err = strconv.ParseFloat(bid.Price, 64)
		if err != nil {
			op.Maker, op.Error = e.doWebOrderBook(op)
			return op.Error
		}
		buydata.Quantity, err = strconv.ParseFloat(bid.Amount, 64)
		if err != nil {
			op.Maker, op.Error = e.doWebOrderBook(op)
			return op.Error
		}

		maker.Bids = append(maker.Bids, buydata)
	}

	for _, ask := range orderBook.Ask {
		var selldata exchange.Order

		//Modify according to type and structure
		selldata.Rate, err = strconv.ParseFloat(ask.Price, 64)
		if err != nil {
			op.Maker, op.Error = e.doWebOrderBook(op)
			return op.Error
		}
		selldata.Quantity, err = strconv.ParseFloat(ask.Amount, 64)
		if err != nil {
			op.Maker, op.Error = e.doWebOrderBook(op)
			return op.Error
		}

		maker.Asks = append(maker.Asks, selldata)
	}
	op.Maker = maker
	return nil
}

// OrderBook from webpage
func (e *Stex) doWebOrderBook(op *exchange.PublicOperation) (*exchange.Maker, error) { //(pair *pair.Pair) (*exchange.Maker, error) {
	pair := op.Pair

	orderBookBuy := WebOrderBook{}
	orderBookSell := WebOrderBook{}

	// strRequestUrl := fmt.Sprintf("/public/orderbook/%v", e.GetIDByPair(pair))
	strUrlBuy := fmt.Sprintf("https://app.stex.com/en/basic-trade/buy-glass/%v", e.GetIDByPair(pair))
	strUrlSell := fmt.Sprintf("https://app.stex.com/en/basic-trade/sell-glass/%v", e.GetIDByPair(pair))

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	getBuy := &utils.HttpGet{
		URI:       strUrlBuy,
		Proxy:     op.Proxy,
		DebugMode: op.DebugMode,
	}
	if err := utils.HttpGetRequest(getBuy); err != nil {
		op.Error = err
		return nil, op.Error
	}

	getSell := &utils.HttpGet{
		URI:       strUrlSell,
		Proxy:     op.Proxy,
		DebugMode: op.DebugMode,
	}
	if err := utils.HttpGetRequest(getSell); err != nil {
		op.Error = err
		return nil, op.Error
	}

	// jsonOrderbookBuy := exchange.HttpGetRequest(strUrlBuy, nil)
	jsonOrderbookBuy := getBuy.ResponseBody
	if err := json.Unmarshal([]byte(jsonOrderbookBuy), &orderBookBuy); err != nil {
		return nil, fmt.Errorf("%s Get WebOrderbook Json Unmarshal Err: %v %s", e.GetName(), err, jsonOrderbookBuy)
	} else if len(orderBookBuy) == 0 {
		return nil, fmt.Errorf("Got empty WebOrderbook: %s", jsonOrderbookBuy)
	}

	// jsonOrderbookSell := exchange.HttpGetRequest(strUrlSell, nil)
	jsonOrderbookSell := getSell.ResponseBody
	if err := json.Unmarshal([]byte(jsonOrderbookSell), &orderBookSell); err != nil {
		return nil, fmt.Errorf("%s Get WebOrderbook Json Unmarshal Err: %v %s", e.GetName(), err, jsonOrderbookSell)
	} else if len(orderBookSell) == 0 {
		return nil, fmt.Errorf("Got empty WebOrderbook: %s", jsonOrderbookSell)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	var err error
	for _, bid := range orderBookBuy {
		var buydata exchange.Order

		buydata.Rate, err = strconv.ParseFloat(bid.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("%s WebOrderbook strconv.ParseFloat Rate error:%v", e.GetName(), err)
		}
		buydata.Quantity, err = strconv.ParseFloat(bid.Amount, 64)
		if err != nil {
			return nil, fmt.Errorf("%s WebOrderbook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
		}

		maker.Bids = append(maker.Bids, buydata)
	}

	for _, ask := range orderBookSell {
		var selldata exchange.Order

		selldata.Rate, err = strconv.ParseFloat(ask.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("%s WebOrderbook strconv.ParseFloat Rate error:%s", e.GetName(), err)
		}
		selldata.Quantity, err = strconv.ParseFloat(ask.Amount, 64)
		if err != nil {
			return nil, fmt.Errorf("%s WebOrderbook strconv.ParseFloat Quantity error:%s", e.GetName(), err)
		}

		maker.Asks = append(maker.Asks, selldata)
	}

	return maker, nil
}

func (e *Stex) doTradeHistory(operation *exchange.PublicOperation) error {
	get := &utils.HttpGet{
		URI:   fmt.Sprintf("%s/public/trades/%s", API3_URL, e.GetIDByPair(operation.Pair)),
		Proxy: operation.Proxy,
	}

	err := utils.HttpGetRequest(get)

	if err != nil {
		return err

	} else {
		jsonResponse := JsonResponseV3{}
		tradeHistory := TradeHistory{}
		if err := json.Unmarshal(get.ResponseBody, &jsonResponse); err != nil {
			return err
		}

		if err := json.Unmarshal(jsonResponse.Data, &tradeHistory); err != nil {
			return err
		}

		operation.TradeHistory = []*exchange.TradeDetail{}
		for i := len(tradeHistory) - 1; i > 0; i-- {
			d := tradeHistory[i]
			td := &exchange.TradeDetail{}

			td.ID = fmt.Sprintf("%d", d.ID)
			if d.Type == "BUY" {
				td.Direction = exchange.Buy
			} else if d.Type == "SELL" {
				td.Direction = exchange.Sell
			}

			td.Quantity, err = strconv.ParseFloat(d.Amount, 64)
			td.Rate, err = strconv.ParseFloat(d.Price, 64)

			t, err := strconv.ParseInt(d.Timestamp, 10, 64)
			if err != nil {
				return err
			}
			td.TimeStamp = t * 1000

			operation.TradeHistory = append(operation.TradeHistory, td)
		}
	}

	return nil
}
