package huobi

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
func (e *Huobi) LoadPublicData(operation *exchange.PublicOperation) error {
	switch operation.Type {

	case exchange.TradeHistory:
		return e.doTradeHistory(operation)
	case exchange.CoinChainType:
		return e.getCoinChainType(operation)
	case exchange.Orderbook:
		if operation.Wallet == exchange.SpotWallet {
			return e.doSpotOrderBook(operation)
		}
	}
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

func (e *Huobi) doSpotOrderBook(op *exchange.PublicOperation) error {
	jsonResponse := &JsonResponse{}
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(op.Pair)

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol
	mapParams["type"] = "step0"

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
		URI:       fmt.Sprintf("%s/market/depth?%s", API_URL, utils.Map2UrlQuery(mapParams)),
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
	} else if jsonResponse.Status != "ok" {
		return fmt.Errorf("%s Get Orderbook Failed: %s", e.GetName(), jsonOrderbook)
	}
	if err := json.Unmarshal(jsonResponse.Tick, &orderBook); err != nil {
		return fmt.Errorf("%s Get Orderbook Result Unmarshal Err: %s %s", e.GetName(), err, jsonResponse.Tick)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	for _, bid := range orderBook.Bids {
		var buydata exchange.Order

		buydata.Rate = bid[0]
		buydata.Quantity = bid[1]

		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Asks {
		var selldata exchange.Order

		selldata.Rate = ask[0]
		selldata.Quantity = ask[1]

		maker.Asks = append(maker.Asks, selldata)
	}

	op.Maker = maker
	return nil
}

func (e *Huobi) doTradeHistory(operation *exchange.PublicOperation) error {

	get := &utils.HttpGet{
		URI: fmt.Sprintf("https://api.huobi.pro/market/history/trade?symbol=%s&size=%d",
			e.GetSymbolByPair(operation.Pair),
			1000, //TRADE_HISTORY_MAX_LIMIT,
		),
		Proxy: operation.Proxy,
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
			for _, d2 := range tradeHistory.Data[i].Data {
				// d2 := d1.Data[i]
				// log.Printf("d2:%+v", d2)
				td := &exchange.TradeDetail{
					ID:       fmt.Sprintf("%d", d2.TradeID),
					Quantity: d2.Amount,

					TimeStamp: d2.Ts,
					Rate:      d2.Price,
				}

				if d2.Direction == "buy" {
					td.Direction = exchange.Buy
				} else if d2.Direction == "sell" {
					td.Direction = exchange.Sell
				}
				// log.Printf("d2: %+v ", d2)
				// log.Printf("TD: %+v ", td)

				operation.TradeHistory = append(operation.TradeHistory, td)
			}
		}
	}

	return nil

}

func (e *Huobi) getCoinChainType(operation *exchange.PublicOperation) error {
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
