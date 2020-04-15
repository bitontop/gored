package binance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bitontop/gored/exchange"
)

// contract balance has too many detail
// contract balance doc
// https://binance-docs.github.io/apidocs/futures/cn/#user_data-5
func (e *Binance) doContractAllBalance(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	accountBalances := ContractBalance{}
	strRequest := "/fapi/v1/account"

	// balanceList := []exchange.AssetBalance{}

	jsonAllBalanceReturn := e.ContractApiKeyRequest("GET", make(map[string]string), strRequest)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonAllBalanceReturn
	}

	log.Printf("jsonAllBalanceReturn: %v", jsonAllBalanceReturn)
	if err := json.Unmarshal([]byte(jsonAllBalanceReturn), &accountBalances); err != nil {
		operation.Error = fmt.Errorf("%s ContractAllBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonAllBalanceReturn)
		return operation.Error
	} else if len(accountBalances.Assets) == 0 {
		operation.Error = fmt.Errorf("%s ContractAllBalance Failed: %v", e.GetName(), jsonAllBalanceReturn)
		return operation.Error
	}

	for _, account := range accountBalances.Assets {
		if account.WalletBalance == "0" {
			continue
		}

		total, err := strconv.ParseFloat(account.WalletBalance, 64)
		available, err := strconv.ParseFloat(account.MaxWithdrawAmount, 64)
		frozen := total - available
		if err != nil {
			return fmt.Errorf("%s balance parse fail: %v %+v", e.GetName(), err, account)
		}

		balance := exchange.AssetBalance{
			Coin:             e.GetCoinBySymbol(account.Asset),
			BalanceAvailable: available,
			BalanceFrozen:    frozen,
		}
		operation.BalanceList = append(operation.BalanceList, balance)

	}

	return nil
	// return fmt.Errorf("%s getBalance fail: %v", e.GetName(), jsonBalanceReturn)
}

func (e *Binance) doContractPlaceOrder(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key or passphrase are nil.", e.GetName())
	}

	placeOrder := ContractPlaceOrder{}
	// strRequestUrl := "/fapi/v1/order/test" // test api
	strRequestUrl := "/fapi/v1/order"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(operation.Pair)
	if operation.OrderDirection == exchange.Buy {
		mapParams["side"] = "BUY"
	} else {
		mapParams["side"] = "SELL"
	}
	mapParams["type"] = "LIMIT"
	mapParams["price"] = fmt.Sprintf("%v", operation.Rate)
	mapParams["quantity"] = fmt.Sprintf("%v", operation.Quantity)
	mapParams["timeInForce"] = string(operation.OrderType) //"GTC"
	//  timeInForce:
	// 	GTC - Good Till Cancel 成交为止
	//  IOC - Immediate or Cancel 无法立即成交(吃单)的部分就撤销
	//  FOK - Fill or Kill 无法全部立即成交就撤销
	//  GTX - Good Till Crossing 无法成为挂单方就撤销

	jsonCreatePlaceOrder := e.ContractApiKeyRequest("POST", mapParams, strRequestUrl)
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
		order.Side = "Buy"
	} else if operation.OrderDirection == exchange.Sell {
		order.Side = "Sell"
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

	jsonCreateOrderStatus := e.ContractApiKeyRequest("GET", mapParams, strRequestUrl)
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

	jsonCreateCancelOrder := e.ContractApiKeyRequest("DELETE", mapParams, strRequestUrl)
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

func (e *Binance) ContractApiKeyRequest(strMethod string, mapParams map[string]string, strRequestPath string) string {
	mapParams["recvWindow"] = "500000" //"50000000"
	mapParams["timestamp"] = fmt.Sprintf("%d", time.Now().UTC().UnixNano()/int64(time.Millisecond))
	mapParams["signature"] = exchange.ComputeHmac256NoDecode(exchange.Map2UrlQuery(mapParams), e.API_SECRET)

	payload := exchange.Map2UrlQuery(mapParams)
	strUrl := CONTRACT_URL + strRequestPath + "?" + payload

	// log.Printf("===============server ts: %v, request ts: %v", exchange.HttpGetRequest(CONTRACT_URL+"/fapi/v1/time", nil), mapParams["timestamp"])
	// log.Printf("%v===============strUrl: %v", strMethod, strUrl)
	// log.Printf("================payload: %v", payload)

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
