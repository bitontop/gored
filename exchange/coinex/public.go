package coinex

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/utils"
)

func (e *Coinex) LoadPublicData(operation *exchange.PublicOperation) error {
	switch operation.Type {
	case exchange.TradeHistory:
		return e.doTradeHistory(operation)
	case exchange.CoinChainType:
		return e.getCoinChainType(operation)
	case exchange.Orderbook:
		if operation.Wallet == exchange.SpotWallet {
			return e.doSpotOrderBook(operation)
		}
	case exchange.KLine:
		switch operation.Wallet {
		case exchange.SpotWallet:
			return e.doSpotKline(operation)
		}
	case exchange.GetTickerPrice:
		switch operation.Wallet {
		case exchange.SpotWallet:
			return e.doTickerPrice(operation)
		}
	}
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

// interval options: 1min, 3min, 5min, 15min, 30min, 1hour, 2hour, 4hour, 6hour, 12hour, 1day, 3day, 1week
func (e *Coinex) doSpotKline(operation *exchange.PublicOperation) error {
	interval := "5min"
	if operation.KlineInterval != "" {
		interval = operation.KlineInterval
	}

	get := &utils.HttpGet{
		URI: fmt.Sprintf("https://api.coinex.com/v1/market/kline?market=%v&type=%v&limit=1000", // ETHBTC
			e.GetSymbolByPair(operation.Pair), // BTCUSDT
			interval,
		),
		Proxy: operation.Proxy,
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

	jsonKLine := get.ResponseBody
	jsonResponse := &JsonResponse{}
	var rawKline [][]interface{}

	if err := json.Unmarshal([]byte(jsonKLine), &jsonResponse); err != nil {
		return fmt.Errorf("%s doSpotKline Json Unmarshal Err: %s %s", e.GetName(), err, jsonKLine)
	} else if jsonResponse.Code != 0 {
		return fmt.Errorf("%s doSpotKline Failed: %s", e.GetName(), jsonKLine)
	}
	if err := json.Unmarshal(jsonResponse.Data, &rawKline); err != nil {
		return fmt.Errorf("%s doSpotKline Result Unmarshal Err: %s %s", e.GetName(), err, jsonResponse.Data)
	}

	operation.Kline = []*exchange.KlineDetail{}
	for _, k := range rawKline {
		open, err := strconv.ParseFloat(k[1].(string), 64)
		if err != nil {
			log.Printf("%s open parse Err: %v %v", e.GetName(), err, k[1])
			operation.Error = err
			return err
		}
		close, err := strconv.ParseFloat(k[2].(string), 64)
		if err != nil {
			log.Printf("%s close parse Err: %v %v", e.GetName(), err, k[2])
			operation.Error = err
			return err
		}
		high, err := strconv.ParseFloat(k[3].(string), 64)
		if err != nil {
			log.Printf("%s high parse Err: %v %v", e.GetName(), err, k[3])
			operation.Error = err
			return err
		}
		low, err := strconv.ParseFloat(k[4].(string), 64)
		if err != nil {
			log.Printf("%s low parse Err: %v %v", e.GetName(), err, k[4])
			operation.Error = err
			return err
		}
		volume, err := strconv.ParseFloat(k[5].(string), 64)
		if err != nil {
			log.Printf("%s volume parse Err: %v %v", e.GetName(), err, k[5])
			operation.Error = err
			return err
		}

		detail := &exchange.KlineDetail{
			Exchange: e.GetName(),
			Pair:     operation.Pair.Name,
			OpenTime: k[0].(float64) * 1000,
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

func (e *Coinex) doTickerPrice(operation *exchange.PublicOperation) error {
	jsonResponse := &JsonResponse{}
	tickerPrice := TickerPrice{}

	get := &utils.HttpGet{
		URI:       fmt.Sprintf("%s/v1/market/ticker/all", API_URL),
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
		return fmt.Errorf("%s doTickerPrice Json Unmarshal Err: %s %s", e.GetName(), err, jsonTickerPrice)
	} else if jsonResponse.Code != 0 {
		return fmt.Errorf("%s doTickerPrice Failed: %s", e.GetName(), jsonTickerPrice)
	}
	if err := json.Unmarshal(jsonResponse.Data, &tickerPrice); err != nil {
		return fmt.Errorf("%s doTickerPrice Result Unmarshal Err: %s %s", e.GetName(), err, jsonResponse.Data)
	}

	operation.TickerPrice = []*exchange.TickerPriceDetail{}
	for symbol, tp := range tickerPrice.Ticker {
		p := e.GetPairBySymbol(symbol)
		price, err := strconv.ParseFloat(tp.Last, 64)
		if err != nil {
			log.Printf("%s doTickerPrice parse Err: %v %v", e.GetName(), err, tp.Last)
			operation.Error = err
			return err
		}

		if p == nil {
			if operation.DebugMode {
				log.Printf("doTickerPrice got nil pair for symbol: %v", symbol)
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

func (e *Coinex) doSpotOrderBook(op *exchange.PublicOperation) error {
	jsonResponse := &JsonResponse{}
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(op.Pair)

	mapParams := make(map[string]string)
	mapParams["market"] = symbol
	mapParams["merge"] = "0"

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
		URI:       fmt.Sprintf("%s/v1/market/depth?%s", API_URL, utils.Map2UrlQuery(mapParams)),
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
	var err error
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

	for _, ask := range orderBook.Asks {
		selldata := exchange.Order{}

		selldata.Rate, err = strconv.ParseFloat(ask[0], 64)
		if err != nil {
			return fmt.Errorf("%s OrderBook strconv.ParseFloat Rate error:%v", e.GetName(), err)
		}
		selldata.Quantity, err = strconv.ParseFloat(ask[1], 64)
		if err != nil {
			return fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
		}

		maker.Asks = append(maker.Asks, selldata)
	}
	maker.LastUpdateID = orderBook.Time

	op.Maker = maker
	return nil
}

func (e *Coinex) doTradeHistory(operation *exchange.PublicOperation) error {
	symbol := e.GetSymbolByPair(operation.Pair)

	get := &utils.HttpGet{
		URI:   fmt.Sprintf("%s/v1/market/deals?market=%s", API_URL, symbol),
		Proxy: operation.Proxy,
	}

	err := utils.HttpGetRequest(get)

	if err != nil {
		return err

	} else {
		jsonResponse := JsonResponse{}
		tradeHistory := TradeHistory{}
		if err := json.Unmarshal(get.ResponseBody, &jsonResponse); err != nil {
			return err
		}

		if err := json.Unmarshal(jsonResponse.Data, &tradeHistory); err != nil {
			return err
		}

		operation.TradeHistory = []*exchange.TradeDetail{}
		// for _, d := range tradeHistory {
		for i := len(tradeHistory) - 1; i > 0; i-- {
			d := tradeHistory[i]
			td := &exchange.TradeDetail{}

			td.ID = fmt.Sprintf("%d", d.ID)
			if d.Type == "buy" {
				td.Direction = exchange.Buy
			} else if d.Type == "sell" {
				td.Direction = exchange.Sell
			}

			td.Quantity, err = strconv.ParseFloat(d.Amount, 64)
			td.Rate, err = strconv.ParseFloat(d.Price, 64)

			td.TimeStamp = d.DateMs

			operation.TradeHistory = append(operation.TradeHistory, td)
		}
	}

	return nil
}

func (e *Coinex) getCoinChainType(operation *exchange.PublicOperation) error {
	operation.CoinChainType = []exchange.ChainType{}
	request := &exchange.ChainTypeRequest{
		Exchange: string(operation.EX),
		CoinID:   operation.Coin.ID,
	}

	byteJson, err := json.Marshal(request)
	post := &utils.HttpPost{
		URI:         "http://127.0.0.1:52020/getchaintype",
		RequestBody: byteJson,
	}

	err = utils.HttpPostRequest(post)
	if err != nil {
		return err

	} else {
		chainType := []*exchange.ChainTypeRequest{}
		if err := json.Unmarshal(post.ResponseBody, &chainType); err != nil {
			return err
		}

		for _, data := range chainType {
			for _, ct := range data.ChainType {
				switch ct {
				case "MAINNET":
					operation.CoinChainType = append(operation.CoinChainType, exchange.MAINNET)
				case "BEP2":
					operation.CoinChainType = append(operation.CoinChainType, exchange.BEP2)
				case "ERC20":
					operation.CoinChainType = append(operation.CoinChainType, exchange.ERC20)
				case "NEP5":
					operation.CoinChainType = append(operation.CoinChainType, exchange.NEP5)
				case "OMNI":
					operation.CoinChainType = append(operation.CoinChainType, exchange.OMNI)
				case "TRC20":
					operation.CoinChainType = append(operation.CoinChainType, exchange.TRC20)
				default:
					operation.CoinChainType = append(operation.CoinChainType, exchange.OTHER)
				}
			}
		}
	}

	return nil
}
