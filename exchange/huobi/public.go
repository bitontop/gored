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
