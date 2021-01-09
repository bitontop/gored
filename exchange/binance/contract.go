package binance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/bitontop/gored/exchange"
)

func (e *Binance) doContractGetOpenOrder(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key or passphrase are nil.", e.GetName())
	}

	openOrders := ContractOpenOrders{}
	strRequest := "/fapi/v1/openOrders"

	mapParams := make(map[string]string)
	if operation.Pair != nil {
		mapParams["symbol"] = e.GetSymbolByPair(operation.Pair)
	}

	jsonGetOpenOrder := e.ContractApiKeyRequest("GET", mapParams, strRequest, operation.TestMode)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		// operation.MapParams = fmt.Sprintf("%+v", mapParams)
		operation.CallResponce = jsonGetOpenOrder
	}

	if err := json.Unmarshal([]byte(jsonGetOpenOrder), &openOrders); err != nil {
		operation.Error = fmt.Errorf("%s doContractGetOpenOrder Json Unmarshal Err: %v, %s", e.GetName(), err, jsonGetOpenOrder)
		return operation.Error
	}

	// store info into orders
	operation.OpenOrders = []*exchange.Order{}
	for _, o := range openOrders {
		rate, err := strconv.ParseFloat(o.Price, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doContractGetOpenOrder parse rate Err: %v, %v", e.GetName(), err, o.Price)
			return operation.Error
		}
		quantity, err := strconv.ParseFloat(o.OrigQty, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doContractGetOpenOrder parse quantity Err: %v, %v", e.GetName(), err, o.OrigQty)
			return operation.Error
		}
		dealQuantity, err := strconv.ParseFloat(o.ExecutedQty, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doContractGetOpenOrder parse dealQuantity Err: %v, %v", e.GetName(), err, o.ExecutedQty)
			return operation.Error
		}

		order := &exchange.Order{
			Pair:         e.GetPairBySymbol(o.Symbol),
			OrderID:      fmt.Sprintf("%v", o.OrderID),
			Rate:         rate,
			Quantity:     quantity,
			DealRate:     rate,
			DealQuantity: dealQuantity,
			Timestamp:    o.UpdateTime,
			// JsonResponse: jsonGetOpenOrder,
		}

		switch o.Side {
		case "BUY":
			order.Direction = exchange.Buy
		case "SELL":
			order.Direction = exchange.Sell
		}

		if o.Status == "CANCELED" {
			order.Status = exchange.Cancelled
		} else if o.Status == "FILLED" {
			order.Status = exchange.Filled
		} else if o.Status == "PARTIALLY_FILLED" {
			order.Status = exchange.Partial
		} else if o.Status == "REJECTED" {
			order.Status = exchange.Rejected
		} else if o.Status == "EXPIRED" {
			order.Status = exchange.Expired
		} else if o.Status == "NEW" {
			order.Status = exchange.New
		} else {
			order.Status = exchange.Other
		}

		operation.OpenOrders = append(operation.OpenOrders, order)
	}

	return nil
}

func (e *Binance) doContractGetOrderHistory(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key or passphrase are nil.", e.GetName())
	}

	allOrders := ContractOpenOrders{}
	strRequest := "/fapi/v1/allOrders"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(operation.Pair)

	jsonGetOpenOrder := e.ContractApiKeyRequest("GET", mapParams, strRequest, operation.TestMode)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		// operation.MapParams = fmt.Sprintf("%+v", mapParams)
		operation.CallResponce = jsonGetOpenOrder
	}

	if err := json.Unmarshal([]byte(jsonGetOpenOrder), &allOrders); err != nil {
		operation.Error = fmt.Errorf("%s doContractGetOrderHistory Json Unmarshal Err: %v, %s", e.GetName(), err, jsonGetOpenOrder)
		return operation.Error
	}

	// store info into orders
	operation.OrderHistory = []*exchange.Order{}
	for _, o := range allOrders {
		rate, err := strconv.ParseFloat(o.Price, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doContractGetOrderHistory parse rate Err: %v, %v", e.GetName(), err, o.Price)
			return operation.Error
		}
		quantity, err := strconv.ParseFloat(o.OrigQty, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doContractGetOrderHistory parse quantity Err: %v, %v", e.GetName(), err, o.OrigQty)
			return operation.Error
		}
		dealQuantity, err := strconv.ParseFloat(o.ExecutedQty, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doContractGetOrderHistory parse dealQuantity Err: %v, %v", e.GetName(), err, o.ExecutedQty)
			return operation.Error
		}

		order := &exchange.Order{
			Pair:         operation.Pair,
			OrderID:      fmt.Sprintf("%v", o.OrderID),
			Rate:         rate,
			Quantity:     quantity,
			DealRate:     rate,
			DealQuantity: dealQuantity,
			Timestamp:    o.UpdateTime,
			// JsonResponse: jsonGetOpenOrder,
		}

		switch o.Side {
		case "BUY":
			order.Direction = exchange.Buy
		case "SELL":
			order.Direction = exchange.Sell
		}

		if o.Status == "CANCELED" {
			order.Status = exchange.Cancelled
		} else if o.Status == "FILLED" {
			order.Status = exchange.Filled
		} else if o.Status == "PARTIALLY_FILLED" {
			// order.Status = exchange.Partial
			continue
		} else if o.Status == "REJECTED" {
			order.Status = exchange.Rejected
		} else if o.Status == "EXPIRED" {
			order.Status = exchange.Expired
		} else if o.Status == "NEW" {
			// order.Status = exchange.New
			continue
		} else {
			order.Status = exchange.Other
		}

		operation.OrderHistory = append(operation.OrderHistory, order)
	}

	return nil
}

func (e *Binance) doContractGetTransferHistory(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key or passphrase are nil.", e.GetName())
	}

	if operation.Coin == nil {
		return fmt.Errorf("%s doContractGetTransferHistory got nil coin", e.GetName())
	}

	transfer := ContractTransferHistory{}
	strRequest := "/sapi/v1/futures/transfer"

	mapParams := make(map[string]string)
	mapParams["asset"] = e.GetSymbolByCoin(operation.Coin)
	mapParams["startTime"] = fmt.Sprintf("%v", operation.TransferStartTime)
	mapParams["size"] = "100"

	jsonTransferOutHistory := e.ApiKeyGet(mapParams, strRequest)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonTransferOutHistory
	}

	if err := json.Unmarshal([]byte(jsonTransferOutHistory), &transfer); err != nil {
		operation.Error = fmt.Errorf("%s doContractGetTransferHistory Json Unmarshal Err: %v, %s", e.GetName(), err, jsonTransferOutHistory)
		return operation.Error
	}

	// store info into orders
	operation.TransferOutHistory = []*exchange.TransferHistory{}
	operation.TransferInHistory = []*exchange.TransferHistory{}
	for _, tx := range transfer.Rows {
		c := e.GetCoinBySymbol(tx.Asset)
		quantity, err := strconv.ParseFloat(tx.Amount, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doContractGetTransferHistory parse quantity Err: %v, %v", e.GetName(), err, tx.Amount)
			return operation.Error
		}

		record := &exchange.TransferHistory{
			Coin:      c,
			Quantity:  quantity,
			TimeStamp: tx.Timestamp,
			StatusMsg: tx.Status,
		}

		switch tx.Type {
		case 1:
			record.Type = exchange.TransferIn
			operation.TransferInHistory = append(operation.TransferInHistory, record)
		case 2:
			record.Type = exchange.TransferOut
			operation.TransferOutHistory = append(operation.TransferOutHistory, record)
		}
	}

	return nil
}

func (e *Binance) doGetPositionInfo(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	positionInfo := PositionInfo{}
	strRequest := "/fapi/v1/positionRisk"

	jsonPositionInfoReturn := e.ContractApiKeyRequest("GET", make(map[string]string), strRequest, operation.TestMode)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonPositionInfoReturn
	}

	// log.Printf("jsonPositionInfoReturn: %v", jsonPositionInfoReturn)
	if err := json.Unmarshal([]byte(jsonPositionInfoReturn), &positionInfo); err != nil {
		operation.Error = fmt.Errorf("%s doGetPositionInfo Json Unmarshal Err: %v, %s", e.GetName(), err, jsonPositionInfoReturn)
		return operation.Error
	}

	return nil
}

func (e *Binance) doGetFutureBalances(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	futureBalances := []*FutureBalance{}
	strRequest := "/fapi/v2/balance"

	jsonAllBalanceReturn := e.ContractApiKeyRequest("GET", make(map[string]string), strRequest, operation.TestMode)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonAllBalanceReturn
	}

	if err := json.Unmarshal([]byte(jsonAllBalanceReturn), &futureBalances); err != nil {
		operation.Error = fmt.Errorf("%s GetFutureBalances Json Unmarshal Err: %v, %s", e.GetName(), err, jsonAllBalanceReturn)
		return operation.Error
	}

	operation.BalanceList = []exchange.AssetBalance{}
	for _, account := range futureBalances {
		total, err := strconv.ParseFloat(account.Balance, 64)
		available, err := strconv.ParseFloat(account.AvailableBalance, 64)
		frozen := total - available
		if err != nil {
			return fmt.Errorf("%s balance parse fail: %v %+v", e.GetName(), err, account)
		}

		balance := exchange.AssetBalance{
			Coin:             e.GetCoinBySymbol(account.Asset),
			Balance:          total,
			BalanceAvailable: available,
			BalanceFrozen:    frozen,
		}
		operation.BalanceList = append(operation.BalanceList, balance)
	}

	return nil
	// return fmt.Errorf("%s getBalance fail: %v", e.GetName(), jsonBalanceReturn)
}

// // contract balance has too many detail
// // contract balance doc
// // https://binance-docs.github.io/apidocs/futures/cn/#user_data-5
func (e *Binance) doContractAllBalance(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	accountBalances := ContractBalance{}
	strRequest := "/fapi/v1/account"

	jsonAllBalanceReturn := e.ContractApiKeyRequest("GET", make(map[string]string), strRequest, operation.TestMode)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonAllBalanceReturn
	}

	// log.Printf("jsonAllBalanceReturn: %v", jsonAllBalanceReturn)
	if err := json.Unmarshal([]byte(jsonAllBalanceReturn), &accountBalances); err != nil {
		operation.Error = fmt.Errorf("%s ContractAllBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonAllBalanceReturn)
		return operation.Error
	}

	operation.BalanceList = []exchange.AssetBalance{}
	for _, account := range accountBalances.Assets {
		total, err := strconv.ParseFloat(account.WalletBalance, 64)
		available, err := strconv.ParseFloat(account.MaxWithdrawAmount, 64)
		frozen := total - available
		if err != nil {
			return fmt.Errorf("%s balance parse fail: %v %+v", e.GetName(), err, account)
		}

		balance := exchange.AssetBalance{
			Coin:             e.GetCoinBySymbol(account.Asset),
			Balance:          total,
			BalanceAvailable: available,
			BalanceFrozen:    frozen,
		}
		operation.BalanceList = append(operation.BalanceList, balance)

	}

	return nil
	// return fmt.Errorf("%s getBalance fail: %v", e.GetName(), jsonBalanceReturn)
}

func (e *Binance) doSetFutureLeverage(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key or passphrase are nil.", e.GetName())
	}

	if operation.Pair == nil {
		return fmt.Errorf("%s SetFutureLeverage Pair is empty", e.GetName())
	}

	if operation.Leverage < 1 || operation.Leverage > 125 {
		return fmt.Errorf("%s SetFutureLeverage is between 1 to 125: %d", e.GetName(), operation.Leverage)
	}

	futureLeverage := FutureLeverage{}
	strRequestUrl := "/fapi/v1/leverage"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(operation.Pair)
	mapParams["leverage"] = fmt.Sprintf("%d", operation.Leverage)

	jsonSetLeverage := e.ContractApiKeyRequest("POST", mapParams, strRequestUrl, operation.TestMode)
	if operation.DebugMode {
		operation.RequestURI = strRequestUrl
		operation.CallResponce = jsonSetLeverage
	}

	if err := json.Unmarshal([]byte(jsonSetLeverage), &futureLeverage); err != nil {
		operation.Error = fmt.Errorf("%s SetFutureLeverage Json Unmarshal Err: %v, %s", e.GetName(), err, jsonSetLeverage)
		return operation.Error
	}

	if futureLeverage.Leverage != operation.Leverage {
		operation.Error = fmt.Errorf("%s SetFutureLeverage Failed: %v", e.GetName(), futureLeverage)
		return operation.Error
	}

	return nil
}

// type: TRADE_LIMIT, TRADE_MARKET, Trade_STOP_LIMIT, Trade_STOP_MARKET
// Stop order need 'StopRate' param
func (e *Binance) doContractPlaceOrder(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key or passphrase are nil.", e.GetName())
	}

	if operation.OrderDirection == "" {
		return fmt.Errorf("%s ContractPlaceOrder empty OrderDirection: %+v", e.GetName(), operation)
	} else if operation.TradeType == "" {
		return fmt.Errorf("%s ContractPlaceOrder empty TradeType: %+v", e.GetName(), operation)
	}

	placeOrder := ContractPlaceOrder{}
	// strRequestUrl := "/fapi/v1/order/test" // test api
	strRequestUrl := "/fapi/v1/order"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(operation.Pair)
	if operation.OrderDirection == exchange.Buy {
		mapParams["side"] = "BUY"
	} else if operation.OrderDirection == exchange.Sell {
		mapParams["side"] = "SELL"
	}
	if operation.TradeType == exchange.Trade_STOP_LIMIT || operation.TradeType == exchange.Trade_STOP_MARKET {
		mapParams["stopPrice"] = fmt.Sprintf("%v", operation.StopRate)
	}
	mapParams["type"] = string(operation.TradeType) // "LIMIT"
	if operation.Rate != 0 {
		mapParams["price"] = fmt.Sprintf("%v", operation.Rate)
	}
	if operation.Quantity != 0 {
		mapParams["quantity"] = fmt.Sprintf("%v", operation.Quantity)
	}
	if operation.OrderType != "" {
		mapParams["timeInForce"] = string(operation.OrderType) //"GTC"
	}
	//  timeInForce:
	// 	GTC - Good Till Cancel 成交为止
	//  IOC - Immediate or Cancel 无法立即成交(吃单)的部分就撤销
	//  FOK - Fill or Kill 无法全部立即成交就撤销
	//  GTX - Good Till Crossing 无法成为挂单方就撤销

	jsonCreatePlaceOrder := e.ContractApiKeyRequest("POST", mapParams, strRequestUrl, operation.TestMode)
	if operation.DebugMode {
		operation.RequestURI = strRequestUrl
		operation.CallResponce = jsonCreatePlaceOrder
	}

	if err := json.Unmarshal([]byte(jsonCreatePlaceOrder), &placeOrder); err != nil {
		operation.Error = fmt.Errorf("%s ContractPlaceOrder Json Unmarshal Err: %v, %s", e.GetName(), err, jsonCreatePlaceOrder)
		return operation.Error
	} else if placeOrder.OrderID == 0 {
		operation.Error = fmt.Errorf("%s ContractPlaceOrder Failed: %v", e.GetName(), jsonCreatePlaceOrder)
		return operation.Error
	}

	order := &exchange.Order{
		Pair:         operation.Pair,
		OrderID:      fmt.Sprintf("%d", placeOrder.OrderID),
		Rate:         operation.Rate,
		Quantity:     operation.Quantity,
		Status:       exchange.New,
		JsonResponse: jsonCreatePlaceOrder,
	}
	if operation.OrderDirection == exchange.Buy {
		order.Direction = exchange.Buy
	} else if operation.OrderDirection == exchange.Sell {
		order.Direction = exchange.Sell
	}

	operation.Order = order

	return nil
}

func (e *Binance) doContractOrderStatus(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key or passphrase are nil.", e.GetName())
	}

	orderStatus := ContractOrderStatus{}
	// strRequestUrl := "/fapi/v1/order/test" // test api
	strRequestUrl := "/fapi/v1/order"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(operation.Pair)
	mapParams["orderId"] = operation.Order.OrderID

	jsonCreateOrderStatus := e.ContractApiKeyRequest("GET", mapParams, strRequestUrl, operation.TestMode)
	if operation.DebugMode {
		operation.RequestURI = strRequestUrl
		operation.CallResponce = jsonCreateOrderStatus
	}

	if err := json.Unmarshal([]byte(jsonCreateOrderStatus), &orderStatus); err != nil {
		operation.Error = fmt.Errorf("%s ContractOrderStatus Json Unmarshal Err: %v, %s", e.GetName(), err, jsonCreateOrderStatus)
		return operation.Error
	} else if orderStatus.OrderID == 0 {
		operation.Error = fmt.Errorf("%s ContractOrderStatus Failed: %v", e.GetName(), jsonCreateOrderStatus)
		return operation.Error
	}

	// order pointer from operation
	order := operation.Order
	if orderStatus.Status == "CANCELED" {
		order.Status = exchange.Cancelled
	} else if orderStatus.Status == "FILLED" {
		order.Status = exchange.Filled
	} else if orderStatus.Status == "PARTIALLY_FILLED" {
		order.Status = exchange.Partial
	} else if orderStatus.Status == "REJECTED" {
		order.Status = exchange.Rejected
	} else if orderStatus.Status == "Expired" {
		order.Status = exchange.Expired
	} else if orderStatus.Status == "NEW" {
		order.Status = exchange.New
	} else {
		order.Status = exchange.Other
	}

	switch orderStatus.Side {
	case "BUY":
		order.Direction = exchange.Buy
	case "SELL":
		order.Direction = exchange.Sell
	}

	order.DealRate, _ = strconv.ParseFloat(orderStatus.Price, 64)
	order.DealQuantity, _ = strconv.ParseFloat(orderStatus.ExecutedQty, 64)

	return nil
}

func (e *Binance) doContractCancelOrder(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key or passphrase are nil.", e.GetName())
	}

	cancelOrder := ContractCancelOrder{}
	// strRequestUrl := "/fapi/v1/order/test" // test api
	strRequestUrl := "/fapi/v1/order"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(operation.Pair)
	mapParams["orderId"] = operation.Order.OrderID

	jsonCreateCancelOrder := e.ContractApiKeyRequest("DELETE", mapParams, strRequestUrl, operation.TestMode)
	if operation.DebugMode {
		operation.RequestURI = strRequestUrl
		operation.CallResponce = jsonCreateCancelOrder
	}

	if err := json.Unmarshal([]byte(jsonCreateCancelOrder), &cancelOrder); err != nil {
		operation.Error = fmt.Errorf("%s ContractCancelOrder Json Unmarshal Err: %v, %s", e.GetName(), err, jsonCreateCancelOrder)
		return operation.Error
	} else if cancelOrder.OrderID == 0 {
		operation.Error = fmt.Errorf("%s ContractCancelOrder Failed: %v", e.GetName(), jsonCreateCancelOrder)
		return operation.Error
	} /* else if cancelOrder.Status != "CANCELED" {
		operation.Error = fmt.Errorf("%s ContractCancelOrder Failed: %v", e.GetName(), jsonCreateCancelOrder)
		return operation.Error
	} */

	operation.Order.Status = exchange.Canceling
	operation.Order.CancelStatus = jsonCreateCancelOrder

	return nil
}

func (e *Binance) ContractApiKeyRequest(strMethod string, mapParams map[string]string, strRequestPath string, testMode bool) string {
	mapParams["recvWindow"] = "500000" //"50000000"
	mapParams["timestamp"] = fmt.Sprintf("%d", time.Now().UTC().UnixNano()/int64(time.Millisecond))
	mapParams["signature"] = exchange.ComputeHmac256NoDecode(exchange.Map2UrlQuery(mapParams), e.API_SECRET)

	payload := exchange.Map2UrlQuery(mapParams)
	strUrl := CONTRACT_URL + strRequestPath + "?" + payload
	if testMode {
		strUrl = CONTRACT_TESTNET_URL + strRequestPath + "?" + payload
	}

	request, err := http.NewRequest(strMethod, strUrl, nil)
	if nil != err {
		return err.Error()
	}
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	request.Header.Add("X-MBX-APIKEY", e.API_KEY)

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if nil != err {
		return err.Error()
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return err.Error()
	}

	return string(body)
}
