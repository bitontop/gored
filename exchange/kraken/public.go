package kraken

// Contributor 2015-2020 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"fmt"

	exchange "github.com/bitontop/gored/exchange"
)

/*************** PUBLIC  API ***************/
func (e *Kraken) LoadPublicData(operation *exchange.PublicOperation) error {
	switch operation.Type {
	case exchange.TradeHistory:
		return e.doTradeHistory(operation)
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
