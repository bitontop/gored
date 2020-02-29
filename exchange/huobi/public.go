package huobi

// Contributor 2015-2020 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
	"fmt"

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

	}
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

func (e *Huobi) doTradeHistory(operation *exchange.PublicOperation) error {

	get := &utils.HttpGet{
		URI: fmt.Sprintf("https://api.huobi.pro/market/history/trade?symbol=%s&size=%d",
			e.GetSymbolByPair(operation.Pair),
			1000, //TRADE_HISTORY_MAX_LIMIT,
		),
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
