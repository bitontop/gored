package bitz

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
	// API_URL string = "https://apiv2.bit-z.pro" // lookup apiv2.bit-z.pro: no such host
	API_URL string = "https://apiv2.bitz.com"
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
func (e *Bitz) GetCoinsData() error {
	jsonResponse := &JsonResponse{}
	coinsData := make(map[string]interface{})

	strRequestUrl := "/Market/coinRate"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if jsonResponse.Status != 200 {
		return fmt.Errorf("%s Get Coins Failed: %v %v", e.GetName(), jsonResponse.Status, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for coinName, _ := range coinsData {
		c := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			c = coin.GetCoin(coinName)
			if c == nil {
				c = &coin.Coin{}
				c.Code = coinName
				coin.AddCoin(c)
			}
		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(coinName)
		}

		if c != nil {
			coinConstraint := e.GetCoinConstraint(c)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       c.ID,
					Coin:         c,
					ExSymbol:     coinName,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = coinName
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
func (e *Bitz) GetPairsData() error {
	jsonResponse := &JsonResponse{}
	pairsData := make(map[string]*PairsData)

	strRequestUrl := "/Market/symbolList"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if jsonResponse.Status != 200 {
		return fmt.Errorf("%s Get Pairs Failed: %v %v", e.GetName(), jsonResponse.Status, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, data := range pairsData {
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(data.CoinTo)
			target := coin.GetCoin(data.CoinFrom)
			if base != nil && target != nil {
				p = pair.GetPair(base, target)
			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(data.Name)
		}
		if p != nil {
			lotSize, _ := strconv.Atoi(data.NumberFloat)
			priceFilter, _ := strconv.Atoi(data.PriceFloat)
			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{
					PairID:      p.ID,
					Pair:        p,
					ExSymbol:    data.Name,
					MakerFee:    DEFAULT_MAKER_FEE,
					TakerFee:    DEFAULT_TAKER_FEE,
					LotSize:     math.Pow10(lotSize * -1),
					PriceFilter: math.Pow10(priceFilter * -1),
					Listed:      DEFAULT_LISTED,
				}
			} else {
				pairConstraint.ExSymbol = data.Name
				pairConstraint.LotSize = math.Pow10(lotSize * -1)
				pairConstraint.PriceFilter = math.Pow10(priceFilter * -1)
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
func (e *Bitz) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	jsonResponse := &JsonResponse{}
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol

	strRequestUrl := "/Market/depth"
	strUrl := API_URL + strRequestUrl

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if jsonResponse.Status != 200 {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %v %v", e.GetName(), jsonResponse.Status, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	var err error
	for _, bid := range orderBook.Bids {
		var buydata exchange.Order

		//Modify according to type and structure
		buydata.Rate, err = strconv.ParseFloat(bid[0], 64)
		if err != nil {
			return nil, err
		}
		buydata.Quantity, err = strconv.ParseFloat(bid[1], 64)
		if err != nil {
			return nil, err
		}

		maker.Bids = append(maker.Bids, buydata)
	}
	for i := len(orderBook.Asks) - 1; i >= 0; i-- {
		var selldata exchange.Order

		//Modify according to type and structure
		selldata.Rate, err = strconv.ParseFloat(orderBook.Asks[i][0], 64)
		if err != nil {
			return nil, err
		}
		selldata.Quantity, err = strconv.ParseFloat(orderBook.Asks[i][1], 64)
		if err != nil {
			return nil, err
		}

		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

/*************** Private API ***************/
func (e *Bitz) DoAccountOperation(operation *exchange.AccountOperation) error {
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

func (e *Bitz) doWithdraw(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.TradePassword == "" {
		return fmt.Errorf("%s API Key, Secret Key or TradePassword are nil", e.GetName())
	}

	if operation.WithdrawTag != "" {
		operation.Error = fmt.Errorf("%s Withdraw Failed, got tag: %v, for coin: %v", e.GetName(), operation.WithdrawTag, operation.Coin.Code)
		return operation.Error
	}

	jsonResponse := JsonResponse{}
	withdraw := Withdraw{}
	strRequest := "/Trade/coinOut"

	mapParams := make(map[string]string)
	mapParams["coin"] = e.GetSymbolByCoin(operation.Coin) /* strconv.FormatFloat(quantity, 'f', -1, 64) */ //===========
	mapParams["number"] = operation.WithdrawAmount
	mapParams["address"] = operation.WithdrawAddress

	jsonWithdrawReturn := e.ApiKeyPOST(mapParams, strRequest)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonWithdrawReturn
	}

	if err := json.Unmarshal([]byte(jsonWithdrawReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Json Unmarshal Err: %v, %s", e.GetName(), err, jsonWithdrawReturn)
		return operation.Error
	} else if jsonResponse.Status != 200 {
		operation.Error = fmt.Errorf("%s Withdraw Failed: %v", e.GetName(), jsonWithdrawReturn)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &withdraw); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	} else if withdraw.ID == 0 {
		operation.Error = fmt.Errorf("%s Withdraw Failed, withdraw ID = 0: %v", e.GetName(), jsonWithdrawReturn)
		return operation.Error
	}

	operation.WithdrawID = fmt.Sprintf("%v", withdraw.ID)

	return nil
}

func (e *Bitz) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	userInfo := UserInfo{}
	accountBalance := AccountBalances{}
	strRequest := "/Assets/getUserAssets"

	jsonBalanceReturn := e.ApiKeyPOST(make(map[string]string), strRequest)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if jsonResponse.Status != 200 {
		log.Printf("%s UpdateAllBalances Failed: %s", e.GetName(), jsonBalanceReturn)
		return
	}
	if err := json.Unmarshal(jsonResponse.Data, &userInfo); err != nil {
		log.Printf("%s UpdateAllBalances Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return
	}
	if err := json.Unmarshal(userInfo.Info, &accountBalance); err != nil {
		log.Printf("%s Assets are empty: %v %v", e.GetName(), err, userInfo)
		return
	}

	for _, v := range accountBalance {
		c := e.GetCoinBySymbol(v.Name)
		if c != nil {
			freeamount, err := strconv.ParseFloat(v.Over, 64)
			if err != nil {
				log.Printf("%s balance Convert to float64 Error: %v %v", e.GetName(), err, v.Over)
			} else {
				balanceMap.Set(c.Code, freeamount)
			}
		}
	}
}

func (e *Bitz) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" || e.TradePassword == "" {
		log.Printf("%s API Key, Secret Key or TradePassword are nil", e.GetName())
		return false
	}

	jsonResponse := JsonResponse{}
	withdraw := Withdraw{}
	strRequest := "/Trade/coinOut"

	mapParams := make(map[string]string)
	mapParams["coin"] = e.GetSymbolByCoin(coin) /* strconv.FormatFloat(quantity, 'f', -1, 64) */ //===========
	mapParams["number"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["address"] = addr

	jsonWithdrawReturn := e.ApiKeyPOST(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonWithdrawReturn), &jsonResponse); err != nil {
		log.Printf("%s Withdraw Json Unmarshal Err: %v %v", e.GetName(), err, jsonWithdrawReturn)
		return false
	} else if jsonResponse.Status != 200 {
		log.Printf("%s Withdraw Failed: %v %+v", e.GetName(), jsonResponse.Status, jsonResponse)
		return false
	}
	if err := json.Unmarshal(jsonResponse.Data, &withdraw); err != nil {
		log.Printf("%s Withdraw Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return false
	}

	return true
}

func (e *Bitz) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" || e.TradePassword == "" {
		return nil, fmt.Errorf("%s API Key, Secret Key or TradePassword are nil", e.GetName())
	}

	jsonResponse := JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/Trade/addEntrustSheet"

	mapParams := make(map[string]string)
	mapParams["number"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["type"] = "2"
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["tradePwd"] = e.TradePassword

	jsonPlaceReturn := e.ApiKeyPOST(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Status != 200 {
		return nil, fmt.Errorf("%s LimitSell Failed: %s", e.GetName(), jsonPlaceReturn)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	orderID := fmt.Sprintf("%v", placeOrder.ID)
	order := &exchange.Order{
		Pair:         pair,
		OrderID:      orderID,
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Bitz) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" || e.TradePassword == "" {
		return nil, fmt.Errorf("%s API Key, Secret Key or TradePassword are nil", e.GetName())
	}

	jsonResponse := JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/Trade/addEntrustSheet"

	mapParams := make(map[string]string)
	mapParams["number"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["type"] = "1"
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["tradePwd"] = e.TradePassword

	jsonPlaceReturn := e.ApiKeyPOST(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Status != 200 {
		return nil, fmt.Errorf("%s LimitBuy Failed: %s", e.GetName(), jsonPlaceReturn)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	orderID := fmt.Sprintf("%v", placeOrder.ID)
	order := &exchange.Order{
		Pair:         pair,
		OrderID:      orderID,
		Rate:         rate,
		Quantity:     quantity,
		Direction:         exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Bitz) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	orderStatus := OrderDetails{}
	strRequest := "/Trade/getEntrustSheetInfo"

	mapParams := make(map[string]string)
	if order.Pair != nil {
		mapParams["entrustSheetId"] = order.OrderID
	} else {
		return fmt.Errorf("%s Order Status Pair cannot be null!", e.GetName())
	}

	jsonOrderStatus := e.ApiKeyPOST(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if jsonResponse.Status != 200 {
		return fmt.Errorf("%s OrderStatus Failed: %s", e.GetName(), jsonOrderStatus)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order.StatusMessage = jsonOrderStatus
	if orderStatus.Status == 2 {
		order.Status = exchange.Filled
	} else if orderStatus.Status == 1 {
		order.Status = exchange.Partial
	} else if orderStatus.Status == 0 {
		order.Status = exchange.New
	} else if orderStatus.Status == 3 {
		order.Status = exchange.Cancelled
	} else {
		order.Status = exchange.Other
	}

	return nil
}

func (e *Bitz) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Bitz) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	cancelOrder := CancelOrder{}
	strRequest := "/Trade/cancelEntrustSheet"

	mapParams := make(map[string]string)
	mapParams["entrustSheetId"] = order.OrderID

	jsonCancelOrder := e.ApiKeyPOST(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if jsonResponse.Status != 200 {
		return fmt.Errorf("%s CancelOrder Failed: %s", e.GetName(), jsonCancelOrder)
	}
	if err := json.Unmarshal(jsonResponse.Data, &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Bitz) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Bitz) ApiKeyPOST(mapParams map[string]string, strRequestPath string) string {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	strUrl := API_URL + strRequestPath

	mapParams["apiKey"] = e.API_KEY
	mapParams["timeStamp"] = timestamp
	mapParams["nonce"] = timestamp[4:]

	params := exchange.Map2UrlQuery(mapParams) + e.API_SECRET
	mapParams["sign"] = exchange.ComputeMD5(params)

	request, err := http.NewRequest("POST", strUrl, strings.NewReader(exchange.Map2UrlQuery(mapParams)))
	if err != nil {
		return err.Error()
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Accept", "application/json")

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
