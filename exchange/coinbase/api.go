package coinbase

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
)

/*API Base Knowledge
Path: API function. Usually after the base endpoint URL
Method:
	Get - Call a URL, API return a response
	Post - Call a URL & send a request, API return a response
Public API:
	It doesn't need authorization/signature , can be called by browser to get response.
	using exchange.HttpGetRequest/exchange.HttpPostRequest
Private API:
	Authorization/Signature is requried. The signature request should look at Exchange API Document.
	using ApiKeyGet/ApiKeyPost
Response:
	Response is a json structure.
	Copy the json to https://transform.now.sh/json-to-go/ convert to go Struct.
	Add the go Struct to model.go

ex. Get /api/v1/depth
Get - Method
/api/v1/depth - Path*/

/*************** Public API ***************/
/*Get Coins Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func (e *Coinbase) GetCoinsData() error {
	return nil
}

/* GetPairsData - Get Pairs Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func (e *Coinbase) GetPairsData() error {
	return nil
}

/*Get Pair Market Depth
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Get Exchange Pair Code ex. symbol := e.GetPairCode(p)
Step 4: Modify API Path(strRequestUrl)
Step 5: Add Params - Depend on API request
Step 6: Convert the response to Standard Maker struct*/
func (e *Coinbase) OrderBook(p *pair.Pair) (*exchange.Maker, error) {
	orderbook := OrderBook{}
	symbol := e.GetSymbolByPair(p)

	mapParams := make(map[string]string)
	mapParams["level"] = "3"

	strRequestUrl := fmt.Sprintf("/products/%s/book", symbol)
	strUrl := API_URL + strRequestUrl

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbookReturn := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbookReturn), &orderbook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbookReturn)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	var err error
	for _, bid := range orderbook.Bids {
		buydata := exchange.Order{}
		buydata.Quantity, err = strconv.ParseFloat(bid[1], 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
		}

		buydata.Rate, err = strconv.ParseFloat(bid[0], 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Rate error:%v", e.GetName(), err)
		}
		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderbook.Asks {
		selldata := exchange.Order{}
		selldata.Quantity, err = strconv.ParseFloat(ask[1], 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
		}

		selldata.Rate, err = strconv.ParseFloat(ask[0], 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Rate error:%v", e.GetName(), err)
		}
		maker.Asks = append(maker.Asks, selldata)
	}

	return maker, nil
}

/*************** Private API ***************/

func (e *Coinbase) DoAccountOperation(operation *exchange.AccountOperation) error {
	switch operation.Type {
	case exchange.Withdraw:
		return e.doWithdraw(operation)
		// case exchange.Transfer:
		// 	return e.transfer(operation)
		// case exchange.BalanceList:
		// 	return e.getAllBalance(operation)
		// case exchange.Balance:
		// 	return e.getBalance(operation)
	}
	return fmt.Errorf("Operation type invalid: %v", operation.Type)
}

func (e *Coinbase) doWithdraw(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	symbol := e.GetSymbolByCoin(operation.Coin)

	if operation.WithdrawTag == "" && (symbol == "XRP" || symbol == "XLM" || symbol == "EOS" || symbol == "ATOM") {
		operation.Error = fmt.Errorf("%s Withdraw Failed, got empty tag: %v, for coin: %v", e.GetName(), operation.WithdrawTag, symbol)
		return operation.Error
	}

	withdraw := WithdrawResponse{}
	strRequestUrl := "/withdrawals/crypto"

	mapParams := make(map[string]interface{})
	mapParams["amount"] = operation.WithdrawAmount
	mapParams["currency"] = symbol
	mapParams["crypto_address"] = operation.WithdrawAddress
	if operation.WithdrawTag != "" {
		mapParams["destination_tag"] = operation.WithdrawTag
	} else {
		mapParams["no_destination_tag"] = true
	}

	jsonCreateWithdraw := e.ApiKeyRequest("POST", mapParams, strRequestUrl)
	if operation.DebugMode {
		operation.RequestURI = strRequestUrl
		operation.CallResponce = jsonCreateWithdraw
	}

	if err := json.Unmarshal([]byte(jsonCreateWithdraw), &withdraw); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Json Unmarshal Err: %v, %s", e.GetName(), err, jsonCreateWithdraw)
		return operation.Error
	} else if withdraw.ID == "" {
		operation.Error = fmt.Errorf("%s Withdraw Failed: %v", e.GetName(), jsonCreateWithdraw)
		return operation.Error
	}

	operation.WithdrawID = withdraw.ID

	return nil
}

func (e *Coinbase) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	accountBalance := AccountBalances{}
	strRequest := "/accounts"

	jsonBalanceReturn := e.ApiKeyRequest("GET", nil, strRequest)
	// log.Printf("%s", jsonBalanceReturn)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances json Unmarshal error: %v %s", e.GetName(), err, jsonBalanceReturn)
		return
	} else {
		for _, balance := range accountBalance {
			freeamount, err := strconv.ParseFloat(balance.Available, 64)
			if err != nil {
				log.Printf("%s UpdateAllBalances err: %+v %v", e.GetName(), balance, err)
				return
			} else {
				c := e.GetCoinBySymbol(balance.Currency)
				if c != nil {
					balanceMap.Set(c.Code, freeamount)
				}
			}
		}
	}
}

/* Withdraw(coin *coin.Coin, quantity float64, addr, tag string) */
func (e *Coinbase) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return false
	}

	withdraw := WithdrawResponse{}
	strRequest := "/withdrawals/crypto"

	mapParams := make(map[string]interface{})
	mapParams["amount"] = quantity
	mapParams["currency"] = e.GetSymbolByCoin(coin)
	mapParams["crypto_address"] = addr
	if tag != "" {
		mapParams["destination_tag"] = tag
	}

	jsonWithdrawReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonWithdrawReturn), &withdraw); err != nil {
		err = fmt.Errorf("%s Withdraw Unmarshal Err: %v %v", e.GetName(), err, jsonWithdrawReturn)
		return false
	}

	return true
}

func (e *Coinbase) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	placeOrder := PlaceOrder{}
	strRequest := "/orders"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]interface{})
	mapParams["product_id"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "sell"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["size"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	order := &exchange.Order{
		Pair:         pair,
		Rate:         rate,
		Quantity:     quantity,
		Side:         "Sell",
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		err = fmt.Errorf("%s LimitSell Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
		order.Error = err
		return nil, err
	} else if placeOrder.Message != "" {
		err = fmt.Errorf("%s LimitSell is failed: %s", e.GetName(), placeOrder.Message)
		order.Error = err
		return nil, err
	}

	order.OrderID = placeOrder.ID

	return order, nil
}

func (e *Coinbase) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	placeOrder := PlaceOrder{}
	strRequest := "/orders"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]interface{})
	mapParams["product_id"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "buy"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["size"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	order := &exchange.Order{
		Pair:         pair,
		Rate:         rate,
		Quantity:     quantity,
		Side:         "Buy",
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		err = fmt.Errorf("%s LimitBuy Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
		order.Error = err
		return nil, err
	} else if placeOrder.Message != "" {
		err = fmt.Errorf("%s LimitBuy is failed: %s", e.GetName(), placeOrder.Message)
		order.Error = err
		return nil, err
	}

	order.OrderID = placeOrder.ID

	return order, nil
}

func (e *Coinbase) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	orderStatus := PlaceOrder{}
	strRequest := fmt.Sprintf("/orders/%s", order.OrderID)
	if order.OrderID == "" {
		return fmt.Errorf("%s Order ID is null.", e.GetName())
	}

	jsonOrderStatus := e.ApiKeyRequest("GET", nil, strRequest)
	if jsonOrderStatus == "The order is cancelled." {
		order.Status = exchange.Cancelled
		return nil
	}

	if err := json.Unmarshal([]byte(jsonOrderStatus), &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	}

	if orderStatus.Status == "done" && orderStatus.DoneReason == "filled" {
		order.Status = exchange.Filled
	} else if orderStatus.Status == "open" || orderStatus.Status == "pending" || orderStatus.Status == "active" {
		order.Status = exchange.New
	} else if orderStatus.Status == "done" {
		order.Status = exchange.Cancelled
	} else {
		order.Status = exchange.Other
	}

	order.JsonResponse = jsonOrderStatus
	order.DealRate, _ = strconv.ParseFloat(orderStatus.ExecutedValue, 64)
	order.DealQuantity, _ = strconv.ParseFloat(orderStatus.Size, 64)

	return nil
}

func (e *Coinbase) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Coinbase) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	cancelOrder := ""
	strRequest := fmt.Sprintf("/orders/%s", order.OrderID)
	if order.OrderID == "" {
		return fmt.Errorf("%s Order ID is null.", e.GetName())
	}

	jsonCancelOrder := e.ApiKeyRequest("DELETE", nil, strRequest)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Coinbase) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Coinbase) ApiKeyRequest(strMethod string, mapParams map[string]interface{}, strRequestPath string) string {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	var bytesParams []byte
	if mapParams != nil {
		bytesParams, _ = json.Marshal(mapParams)
	}

	payload := fmt.Sprintf("%s%s%s%s", timestamp, strMethod, strRequestPath, string(bytesParams))
	signature := exchange.ComputeBase64Hmac256(payload, e.API_SECRET)

	strUrl := API_URL + strRequestPath
	request, err := http.NewRequest(strMethod, strUrl, bytes.NewBuffer(bytesParams))
	if nil != err {
		return err.Error()
	}
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	request.Header.Add("User-Agent", "Go Coinbase Pro Client 1.0")
	request.Header.Add("CB-ACCESS-KEY", e.API_KEY)
	request.Header.Add("CB-ACCESS-SIGN", signature)
	request.Header.Add("CB-ACCESS-TIMESTAMP", timestamp)
	request.Header.Add("CB-ACCESS-PASSPHRASE", e.Passphrase)

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if nil != err {
		return err.Error()
	}
	defer response.Body.Close()

	if response.StatusCode == 404 {
		return "The order is cancelled."
	}

	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return err.Error()
	}

	return string(body)
}
