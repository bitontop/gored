package coinex

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
)

/*The Base Endpoint URL*/
const (
	API_URL = "https://api.coinex.com"
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
func (e *Coinex) GetCoinsData() {
	pairsData := PairsData{}

	strRequestUrl := "/v1/market/info"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &pairsData); err != nil {
		log.Printf("%s Get Coins Data Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
		return
	}

	for _, data := range pairsData.Data {
		c := &coin.Coin{}
		c2 := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			c = coin.GetCoin(data.TradingName)
			if c == nil {
				log.Printf("====Trading Name:%v", data.TradingName) //==========
				c = &coin.Coin{}
				c.Code = data.TradingName
				//c.Name = data.AssetName
				//c.Website = data.URL
				//c.Explorer = data.BlockURL
				coin.AddCoin(c)
			}
			c2 = coin.GetCoin(data.PricingName)
			if c2 == nil {
				log.Printf("====Pricing Name:%v", data.PricingName) //==========
				c2 = &coin.Coin{}
				c2.Code = data.PricingName
				//c.Name = data.AssetName
				//c.Website = data.URL
				//c.Explorer = data.BlockURL
				coin.AddCoin(c2)
			}
		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(data.TradingName)
		}

		if c != nil {
			coinConstraint := &exchange.CoinConstraint{
				CoinID:   c.ID,
				Coin:     c,
				ExSymbol: data.TradingName,
				//TxFee:        data.TransactionFee,
				//Withdraw:     data.EnableWithdraw,
				//Deposit:      data.EnableCharge,
				//Confirmation: confirmation,
				Listed: true,
			}

			e.SetCoinConstraint(coinConstraint)
		}
	}
}

/* GetPairsData - Get Pairs Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func (e *Coinex) GetPairsData() {
	pairsData := &PairsData{}

	strRequestUrl := "/v1/market/info"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &pairsData); err != nil {
		log.Printf("%s Get Pairs Data Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
		return
	}

	for _, data := range pairsData.Data {
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(data.PricingName)
			target := coin.GetCoin(data.TradingName)
			if base != nil && target != nil {

				p = pair.GetPair(base, target)

			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(data.Symbol)
		}
		if p != nil {
			makerFee, _ := strconv.ParseFloat(data.MakerFeeRate, 64)
			takerFee, _ := strconv.ParseFloat(data.TakerFeeRate, 64)
			pairConstraint := &exchange.PairConstraint{
				PairID:      p.ID,
				Pair:        p,
				ExSymbol:    data.Symbol,
				MakerFee:    makerFee,
				TakerFee:    takerFee,
				LotSize:     DEFAULT_LOT_SIZE,
				PriceFilter: DEFAULT_PRICE_FILTER,
				Listed:      true,
			}
			e.SetPairConstraint(pairConstraint)
		}
	}
}

/*Get Pair Market Depth
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Get Exchange Pair Code ex. symbol := e.GetPairCode(p)
Step 4: Modify API Path(strRequestUrl)
Step 5: Add Params - Depend on API request
Step 6: Convert the response to Standard Maker struct*/
func (e *Coinex) OrderBook(p *pair.Pair) (*exchange.Maker, error) {
	orderBook := CoinexOrderBook{}
	symbol := e.GetSymbolByPair(p)

	mapParams := make(map[string]string)
	mapParams["market"] = symbol
	mapParams["merge"] = "0"

	strRequestUrl := "/v1/market/depth"
	strUrl := API_URL + strRequestUrl

	maker := &exchange.Maker{}
	maker.WorkerIP = exchange.GetExternalIP()
	maker.BeforeTimestamp = float64(time.Now().UnixNano() / 1e6)

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &orderBook); err != nil {
		log.Printf("%s OrderBook json Unmarshal error: %v %v", e.GetName(), err, jsonOrderbook)
		return nil, err
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)

	var err error
	for _, bid := range orderBook.Data.Bids {
		buydata := exchange.Order{}
		buydata.Quantity, err = strconv.ParseFloat(bid[1], 64)
		if err != nil {
			log.Printf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
			return nil, err
		}

		buydata.Rate, err = strconv.ParseFloat(bid[0], 64)
		if err != nil {
			log.Printf("%s OrderBook strconv.ParseFloat Rate error:%v", e.GetName(), err)
			return nil, err
		}
		maker.Bids = append(maker.Bids, buydata)
	}

	for _, ask := range orderBook.Data.Asks {
		selldata := exchange.Order{}
		selldata.Quantity, err = strconv.ParseFloat(ask[1], 64)
		if err != nil {
			log.Printf("%s OrderBook strconv.ParseFloat  Quantity error:%v", e.GetName(), err)
			return nil, err
		}

		selldata.Rate, err = strconv.ParseFloat(ask[0], 64)
		if err != nil {
			log.Printf("%s OrderBook strconv.ParseFloat  Rate error:%v", e.GetName(), err)
			return nil, err
		}
		maker.Asks = append(maker.Asks, selldata)
	}

	return maker, err
}

/*************** Private API ***************/
func (e *Coinex) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	accountBalance := CoinexAccountBalance{}
	strRequest := "/v1/balance/info"

	mapParams := make(map[string]string)
	mapParams["access_id"] = e.API_KEY

	jsonBalanceReturn := e.ApiKeyRequest("GET", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances json Unmarshal error: %v %s", e.GetName(), err, jsonBalanceReturn)
		return
	} else if accountBalance.Code != 0 {
		log.Printf("CoinEX Balance API Err: Code-%d %s", accountBalance.Code, accountBalance.Message)
	} else {
		for key, balance := range accountBalance.Data {
			if err != nil {
				log.Printf("%s UpdateAllBalances err: %+v %v", e.GetName(), balance, err)
				return
			} else {
				c := e.GetCoinBySymbol(key)
				if c != nil {
					balanceMap.Set(c.Code, balance.Available)
				}
			}
		}
	}
}

/* Withdraw(coin *coin.Coin, quantity float64, addr, tag string) */
func (e *Coinex) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {

	return false
}

func (e *Coinex) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	placeOrder := OrderResponse{}
	strRequest := "/v1/order/limit"
	strUrl := API_URL + strRequest

	mapParams := make(map[string]string)
	mapParams["access_id"] = e.API_KEY
	mapParams["market"] = e.GetSymbolByPair(pair)
	mapParams["type"] = "sell"
	mapParams["amount"] = strconv.FormatFloat(quantity, 'E', -1, 64)
	mapParams["price"] = strconv.FormatFloat(rate, 'E', -1, 64)

	jsonPlaceReturn := e.ApiKeyPost(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if placeOrder.Code != 0 {
		return nil, fmt.Errorf("%s LimitSell failed:%v Message:%v", e.GetName(), placeOrder.Code, placeOrder.Message)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%d", placeOrder.Data.ID),
		Rate:         rate,
		Quantity:     quantity,
		Side:         "Sell",
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}
	return order, nil
}

func (e *Coinex) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	placeOrder := OrderResponse{}
	strRequest := "/v1/order/limit"
	strUrl := API_URL + strRequest

	mapParams := make(map[string]string)
	mapParams["access_id"] = e.API_KEY
	mapParams["market"] = e.GetSymbolByPair(pair)
	mapParams["type"] = "buy"
	mapParams["amount"] = strconv.FormatFloat(quantity, 'E', -1, 64)
	mapParams["price"] = strconv.FormatFloat(rate, 'E', -1, 64)

	jsonPlaceReturn := e.ApiKeyPost(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if placeOrder.Code != 0 {
		return nil, fmt.Errorf("%s LimitBuy failed:%v Message:%v", e.GetName(), placeOrder.Code, placeOrder.Message)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%d", placeOrder.Data.ID),
		Rate:         rate,
		Quantity:     quantity,
		Side:         "Buy",
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Coinex) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	orderStatus := OrderResponse{}
	strRequest := "/v1/order/status"
	strUrl := API_URL + strRequest

	mapParams := make(map[string]string)

	mapParams["access_id"] = e.API_KEY
	mapParams["id"] = order.OrderID
	mapParams["market"] = e.GetSymbolByPair(order.Pair)

	jsonOrderStatus := e.ApiKeyRequest("GET", mapParams, strUrl)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if orderStatus.Code != 0 {
		return fmt.Errorf("%s Get OrderStatus Error: %v %s", e.GetName(), orderStatus.Code, orderStatus.Message)
	}

	order.DealRate, _ = strconv.ParseFloat(orderStatus.Data.AvgPrice, 64)
	order.DealQuantity, _ = strconv.ParseFloat(orderStatus.Data.DealAmount, 64)

	if orderStatus.Data.Status == "done" {
		order.Status = exchange.Filled
	} else if orderStatus.Data.Status == "part_deal" {
		order.Status = exchange.Partial
	} else if orderStatus.Data.Status == "not_deal" {
		order.Status = exchange.New
	} else {
		order.Status = exchange.Other
	}

	order.DealRate, _ = strconv.ParseFloat(orderStatus.Data.Price, 64)
	order.DealQuantity, _ = strconv.ParseFloat(orderStatus.Data.Amount, 64)

	return nil
}

func (e *Coinex) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Coinex) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	cancelOrder := OrderResponse{}
	strRequest := "/v1/order/pending"
	strUrl := API_URL + strRequest

	mapParams := make(map[string]string)
	mapParams["access_id"] = e.API_KEY
	mapParams["id"] = order.OrderID
	mapParams["market"] = e.GetSymbolByPair(order.Pair)

	jsonCancelOrder := e.ApiKeyRequest("DELETE", mapParams, strUrl)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if cancelOrder.Code != 0 {
		return fmt.Errorf("%s CancelOrder Error: %v %s", e.GetName(), cancelOrder.Code, cancelOrder.Message)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Coinex) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Coinex) ApiKeyRequest(strMethod string, mapParams map[string]string, strUrl string) string {
	timestamp := time.Now().UnixNano() / 1e6
	mapParams["tonce"] = strconv.FormatInt(timestamp, 10)

	var strRequestUrl string
	if nil == mapParams {
		strRequestUrl = strUrl
	} else {
		strParams := exchange.Map2UrlQuery(mapParams)
		strRequestUrl = strUrl + "?" + strParams
	}

	// 构建Request, 并且按官方要求添加Http Header
	httpClient := &http.Client{}
	request, err := http.NewRequest(strMethod, strRequestUrl, nil)
	if nil != err {
		return err.Error()
	}
	hasher := md5.New()
	signature := fmt.Sprintf("%s&secret_key=%s", exchange.Map2UrlQuery(mapParams), e.API_SECRET)
	hasher.Write([]byte(signature))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("authorization", strings.ToUpper(hex.EncodeToString(hasher.Sum(nil))))
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")

	// 发出请求
	response, err := httpClient.Do(request)
	if nil != err {
		return err.Error()
	}
	defer response.Body.Close()

	// 解析响应内容
	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return err.Error()
	}

	return string(body)

}

func (e *Coinex) ApiKeyPost(strUrl string, mapParams map[string]string) string {
	timestamp := time.Now().UnixNano() / 1e6
	mapParams["tonce"] = strconv.FormatInt(timestamp, 10)

	jsonParams := ""
	if nil != mapParams {
		bytesParams, _ := json.Marshal(mapParams)
		jsonParams = string(bytesParams)
	}

	// 构建Request, 并且按官方要求添加Http Header
	httpClient := &http.Client{}
	request, err := http.NewRequest("POST", strUrl, strings.NewReader(jsonParams))
	if nil != err {
		return err.Error()
	}
	hasher := md5.New()
	signature := fmt.Sprintf("%s&secret_key=%s", exchange.Map2UrlQuery(mapParams), e.API_SECRET)
	hasher.Write([]byte(signature))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("authorization", strings.ToUpper(hex.EncodeToString(hasher.Sum(nil))))
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")

	// 发出请求
	response, err := httpClient.Do(request)
	if nil != err {
		return err.Error()
	}
	defer response.Body.Close()

	// 解析响应内容
	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return err.Error()
	}

	return string(body)
}
