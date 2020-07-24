package poloniex

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
)

const (
	API_URL string = "https://poloniex.com"
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
func (e *Poloniex) GetCoinsData() error {
	coinsData := make(map[string]*CoinsData)

	strRequestUrl := "/public?command=returnCurrencies"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	}

	for key, data := range coinsData {
		if data.Delisted == 1 {
			continue
		}
		c := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			c = coin.GetCoin(key)
			if c == nil {
				c = &coin.Coin{}
				c.Code = key
				c.Name = data.Name
				coin.AddCoin(c)
			}
		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(key)
		}

		if c != nil {
			txFee, err := strconv.ParseFloat(data.TxFee, 64)
			if err != nil {
				return fmt.Errorf("%s txFee parse error: %v", e.GetName, err)
			}
			coinConstraint := e.GetCoinConstraint(c)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       c.ID,
					Coin:         c,
					ExSymbol:     key,
					ChainType:    exchange.MAINNET,
					TxFee:        txFee,
					Withdraw:     data.Disabled == 0,
					Deposit:      data.Disabled == 0,
					Confirmation: data.MinConf,
					Listed:       data.Delisted == 0,
				}
			} else {
				coinConstraint.ExSymbol = key
				coinConstraint.TxFee = txFee
				coinConstraint.Withdraw = data.Disabled == 0
				coinConstraint.Deposit = data.Disabled == 0
				coinConstraint.Confirmation = data.MinConf
				coinConstraint.Listed = data.Delisted == 0
			}
			e.SetCoinConstraint(coinConstraint)
		}
	}
	return nil
}

/* GetPairsData - Get Pairs Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func (e *Poloniex) GetPairsData() error {
	pairsData := make(map[string]*PairsData)

	strRequestUrl := "/public?command=returnTicker"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	}

	for key, data := range pairsData {
		coinStrs := strings.Split(key, "_")
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(coinStrs[0])
			target := coin.GetCoin(coinStrs[1])
			if base != nil && target != nil {
				p = pair.GetPair(base, target)
			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(key)
		}
		if p != nil {
			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{
					PairID:      p.ID,
					Pair:        p,
					ExID:        fmt.Sprintf("%v", data.ID),
					ExSymbol:    key,
					MakerFee:    DEFAULT_MAKER_FEE,
					TakerFee:    DEFAULT_TAKER_FEE,
					LotSize:     DEFAULT_LOT_SIZE,
					PriceFilter: DEFAULT_PRICE_FILTER,
					Listed:      DEFAULT_LISTED,
				}
				switch p.Base.Code {
				case "USDT", "USD":
					pairConstraint.MinTradeBaseQuantity = 10
				case "BTC":
					pairConstraint.MinTradeBaseQuantity = 0.001
				case "ETH":
					pairConstraint.MinTradeBaseQuantity = 0.05
				}
			} else {
				pairConstraint.ExSymbol = key
			}
			e.SetPairConstraint(pairConstraint)
		}
	}
	return nil
}

/*Get Pair Market Depth
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Get Exchange Pair Code ex. symbol := e.GetPairCode(p)
Step 4: Modify API Path(strRequestUrl)
Step 5: Add Params - Depend on API request
Step 6: Convert the response to Standard Maker struct*/
func (e *Poloniex) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	strRequestUrl := "/public" //?command=returnOrderBook&currencyPair=BTC_ETH&depth=10"
	strUrl := API_URL + strRequestUrl

	mapParams := make(map[string]string)
	mapParams["command"] = "returnOrderBook"
	mapParams["currencyPair"] = symbol

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	var err error
	for _, bid := range orderBook.Bids {
		var buydata exchange.Order

		//Modify according to type and structure
		buydata.Rate, err = strconv.ParseFloat(bid[0].(string), 64)
		if err != nil {
			return nil, fmt.Errorf("Poloniex Bids Rate ParseFloat error:%v", err)
		}
		buydata.Quantity = bid[1].(float64)

		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Asks {
		var selldata exchange.Order

		//Modify according to type and structure
		selldata.Rate, err = strconv.ParseFloat(ask[0].(string), 64)
		if err != nil {
			return nil, fmt.Errorf("Poloniex Asks Rate ParseFloat error:%v", err)
		}
		selldata.Quantity = ask[1].(float64)

		maker.Asks = append(maker.Asks, selldata)
	}
	maker.LastUpdateID = orderBook.Seq
	return maker, nil
}

/*************** Private API ***************/
func (e *Poloniex) DoAccountOperation(operation *exchange.AccountOperation) error {
	switch operation.Type {

	// case exchange.Transfer:
	// 	return e.transfer(operation)
	// case exchange.BalanceList:
	// 	return e.getAllBalance(operation)
	// case exchange.Balance:
	// 	return e.getBalance(operation)

	case exchange.Withdraw:
		return e.doWithdraw(operation)

	}
	return fmt.Errorf("%s Operation type invalid: %s %v", operation.Ex, operation.Wallet, operation.Type)
}

func (e *Poloniex) doWithdraw(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	if operation.WithdrawTag != "" {
		operation.Error = fmt.Errorf("%s Withdraw Failed, got tag: %v, for coin: %v", e.GetName(), operation.WithdrawTag, operation.Coin.Code)
		return operation.Error
	}

	withdraw := Withdraw{}
	strRequest := "/tradingApi"

	mapParams := make(map[string]string)
	mapParams["command"] = "withdraw"
	mapParams["currency"] = e.GetSymbolByCoin(operation.Coin)
	mapParams["amount"] = operation.WithdrawAmount
	mapParams["address"] = operation.WithdrawAddress

	jsonSubmitWithdraw := e.ApiKeyPost(strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonSubmitWithdraw
	}

	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &withdraw); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Json Unmarshal Err: %v, %s", e.GetName(), err, jsonSubmitWithdraw)
		return operation.Error
	} else if withdraw.Response == "" {
		operation.Error = fmt.Errorf("%s Withdraw Failed: %v", e.GetName(), jsonSubmitWithdraw)
		return operation.Error
	}

	// WithdrawID not provided, sample response: { response: 'Withdrew 2.0 ETH.' }
	// operation.WithdrawID = withdraw.Response

	return nil
}

func (e *Poloniex) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	accountBalance := make(map[string]string)
	strRequest := "/tradingApi"
	mapParams := make(map[string]string)
	mapParams["command"] = "returnBalances"

	jsonBalanceReturn := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if accountBalance == nil {
		log.Printf("%s UpdateAllBalances Failed: %v", e.GetName(), jsonBalanceReturn)
		return
	}

	for key, v := range accountBalance {
		freeAmount, err := strconv.ParseFloat(v, 64)
		if err != nil {
			log.Printf("%s balance parse error: %v", e.GetName(), err)
			return
		}

		c := e.GetCoinBySymbol(key)
		if c != nil {
			balanceMap.Set(c.Code, freeAmount)
		}
	}
}

func (e *Poloniex) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil", e.GetName())
		return false
	}

	withdraw := Withdraw{}
	strRequest := "/tradingApi"

	mapParams := make(map[string]string)
	mapParams["command"] = "withdraw"
	mapParams["currency"] = e.GetSymbolByCoin(coin)
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["address"] = addr

	jsonSubmitWithdraw := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &withdraw); err != nil {
		log.Printf("%s Withdraw Json Unmarshal Err: %v %v", e.GetName(), err, jsonSubmitWithdraw)
		return false
	} else if withdraw.Response == "" {
		log.Printf("%s Withdraw Failed: %v", e.GetName(), jsonSubmitWithdraw)
		return false
	}

	return true
}

func (e *Poloniex) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	placeOrder := PlaceOrder{}
	strRequest := "/tradingApi"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["command"] = "sell"
	mapParams["currencyPair"] = e.GetSymbolByPair(pair)
	mapParams["rate"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)

	jsonPlaceReturn := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if placeOrder.OrderNumber == "" {
		return nil, fmt.Errorf("%s LimitSell Failed: %v", e.GetName(), jsonPlaceReturn)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      placeOrder.OrderNumber,
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Poloniex) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	placeOrder := PlaceOrder{}
	strRequest := "/tradingApi"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["command"] = "buy"
	mapParams["currencyPair"] = e.GetSymbolByPair(pair)
	mapParams["rate"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)

	jsonPlaceReturn := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if placeOrder.OrderNumber == "" {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v", e.GetName(), jsonPlaceReturn)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      placeOrder.OrderNumber,
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Poloniex) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	orderStatus := OrderStatus{}
	strRequest := "/tradingApi"

	mapParams := make(map[string]string)
	mapParams["command"] = "returnOrderStatus"
	mapParams["orderNumber"] = order.OrderID

	jsonOrderStatus := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if orderStatus.Success != 1 {
		return fmt.Errorf("%s OrderStatus Failed: %v", e.GetName(), jsonOrderStatus)
	}

	order.StatusMessage = jsonOrderStatus
	for key, data := range orderStatus.Result {
		if key == order.OrderID {
			if data.Status == "Partially filled" {
				order.Status = exchange.Partial
			} else if data.Status == "Open" {
				order.Status = exchange.New
			} else {
				order.Status = exchange.Other
			}
		}
	}

	rate, err := strconv.ParseFloat(orderStatus.Result[order.OrderID].Rate, 64)
	if err != nil {
		log.Printf("%s OrderStatus parse error: %v", e.GetName(), err)
		return err
	}
	quantity, err := strconv.ParseFloat(orderStatus.Result[order.OrderID].StartingAmount, 64)
	if err != nil {
		log.Printf("%s OrderStatus parse error: %v", e.GetName(), err)
		return err
	}
	quantityLeft, err := strconv.ParseFloat(orderStatus.Result[order.OrderID].Amount, 64)
	if err != nil {
		log.Printf("%s OrderStatus parse error: %v", e.GetName(), err)
		return err
	}

	order.Pair = e.GetPairBySymbol(orderStatus.Result[order.OrderID].CurrencyPair)
	order.Rate = rate
	order.Quantity = quantity
	order.DealRate = rate
	order.DealQuantity = order.Quantity - quantityLeft
	if orderStatus.Result[order.OrderID].Type == "buy" {
		order.Direction = exchange.Buy
	} else if orderStatus.Result[order.OrderID].Type == "sell" {
		order.Direction = exchange.Sell
	}

	return nil
}

func (e *Poloniex) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Poloniex) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	mapParams := make(map[string]string)
	mapParams["command"] = "cancelOrder"
	mapParams["orderNumber"] = order.OrderID

	cancelOrder := CancelOrder{}
	strRequest := "/tradingApi"

	jsonCancelOrder := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if cancelOrder.Success != 1 {
		return fmt.Errorf("%s CancelOrder Failed: %s", e.GetName(), jsonCancelOrder)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Poloniex) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Poloniex) ApiKeyPost(strRequestPath string, mapParams map[string]string) string {
	strMethod := "POST"

	mapParams["nonce"] = fmt.Sprintf("%d", time.Now().UnixNano())

	payload := exchange.Map2UrlQuery(mapParams)
	Signature := exchange.ComputeHmac512NoDecode(payload, e.API_SECRET)

	strUrl := API_URL + strRequestPath

	request, err := http.NewRequest(strMethod, strUrl, strings.NewReader(payload))
	if nil != err {
		return err.Error()
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Key", e.API_KEY)
	request.Header.Set("Sign", Signature)

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

func (e *Poloniex) ApiKeyGET(strRequestPath string, mapParams map[string]string) string {
	mapParams["apikey"] = e.API_KEY
	mapParams["nonce"] = fmt.Sprintf("%d", time.Now().UnixNano())

	strUrl := API_URL + strRequestPath + "?" + exchange.Map2UrlQuery(mapParams)

	signature := exchange.ComputeHmac512NoDecode(strUrl, e.API_SECRET)
	httpClient := &http.Client{}

	request, err := http.NewRequest("GET", strUrl, nil)
	if nil != err {
		return err.Error()
	}
	request.Header.Add("Content-Type", "application/json;charset=utf-8")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("apisign", signature)

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
