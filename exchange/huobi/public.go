package huobi

import (
	"fmt"

	"github.com/bitontop/gored/exchange"
)

// Contributionr 2015-2020 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

/*************** PUBLIC  API ***************/
func (e *Huobi) LoadPublicData(operation *exchange.AccountOperation) error {
	switch operation.Type {

	case exchange.TradeHistory:
		return e.doTradeHistory(operation)

	}
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

func (e *Huobi) doTradeHistory(operation *exchange.AccountOperation) error {

	// jsonResponse := &JsonResponse{}
	// orderBook := OrderBook{}
	// symbol := e.GetSymbolByPair(pair)

	// strRequestUrl := "/market/depth"
	// strUrl := API_URL + strRequestUrl

	// mapParams := make(map[string]string)
	// mapParams["symbol"] = symbol
	// mapParams["type"] = "step0"

	// maker := &exchange.Maker{
	// 	WorkerIP:        exchange.GetExternalIP(),
	// 	Source:          exchange.EXCHANGE_API,
	// 	BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	// }

	// jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	// if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
	// 	return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	// } else if jsonResponse.Status != "ok" {
	// 	return nil, fmt.Errorf("%s Get Orderbook Failed: %v", e.GetName(), jsonOrderbook)
	// }
	// if err := json.Unmarshal(jsonResponse.Tick, &orderBook); err != nil {
	// 	return nil, fmt.Errorf("%s Get Orderbook Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Tick)
	// }

	// if len(orderBook.Bids) == 0 || len(orderBook.Asks) == 0 {
	// 	return nil, fmt.Errorf("%s Get Orderbook Failed, Empty Orderbook: %v", e.GetName(), jsonOrderbook)
	// }

	// maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	// for _, bid := range orderBook.Bids {
	// 	var buydata exchange.Order

	// 	buydata.Rate = bid[0]
	// 	buydata.Quantity = bid[1]

	// 	maker.Bids = append(maker.Bids, buydata)
	// }
	// for _, ask := range orderBook.Asks {
	// 	var selldata exchange.Order

	// 	selldata.Rate = ask[0]
	// 	selldata.Quantity = ask[1]

	// 	maker.Asks = append(maker.Asks, selldata)
	// }
	// maker.LastUpdateID = orderBook.Version

	// return maker,nil

	return nil

}

// func (e *Huobi) doWithdraw(operation *exchange.AccountOperation) error {
// 	if e.API_KEY == "" || e.API_SECRET == "" {
// 		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
// 	}

// 	jsonResponse := &JsonResponse{}
// 	var withdrawID int64
// 	strRequest := "/v1/dw/withdraw/api/create"

// 	mapParams := make(map[string]string)
// 	mapParams["address"] = operation.WithdrawAddress
// 	mapParams["amount"] = operation.WithdrawAmount
// 	mapParams["currency"] = e.GetSymbolByCoin(operation.Coin)
// 	// mapParams["fee"] = strconv.FormatFloat(e.GetTxFee(operation.Coin), 'f', -1, 64) // Required parameter
// 	if operation.WithdrawTag != "" {
// 		mapParams["tag"] = operation.WithdrawTag
// 	}

// 	jsonWithdraw := e.ApiKeyRequest("POST", mapParams, strRequest)
// 	if operation.DebugMode {
// 		operation.RequestURI = strRequest
// 		operation.MapParams = fmt.Sprintf("%+v", mapParams)
// 		operation.CallResponce = jsonWithdraw
// 	}

// 	if err := json.Unmarshal([]byte(jsonWithdraw), &jsonResponse); err != nil {
// 		operation.Error = fmt.Errorf("%s Withdraw Json Unmarshal Err: %v, %s", e.GetName(), err, jsonWithdraw)
// 		return operation.Error
// 	} else if jsonResponse.Status != "ok" {
// 		operation.Error = fmt.Errorf("%s Withdraw Failed: %v", e.GetName(), jsonWithdraw)
// 		return operation.Error
// 	}
// 	if err := json.Unmarshal(jsonResponse.Data, &withdrawID); err != nil {
// 		operation.Error = fmt.Errorf("%s Withdraw Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
// 		return operation.Error
// 	}

// 	operation.WithdrawID = fmt.Sprintf("%v", withdrawID)

// 	return nil
// }
