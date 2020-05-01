package bcex

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"

	"strconv"
)

const (
	API_URL string = "http://api.bcex.top"
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
func (e *Bcex) GetCoinsData() error {
	jsonResponse := &JsonResponse{}
	coinsData := make(map[string]*CoinsData)

	strRequestUrl := "/api_market/getTokenPrecision"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if jsonResponse.Code != 0 {
		return fmt.Errorf("%s Get Coins Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Data, &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for key, _ := range coinsData {
		c := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			c = coin.GetCoin(key)
			if c == nil {
				c = &coin.Coin{}
				c.Code = key
				coin.AddCoin(c)
			}
		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(key)
		}

		if c != nil {
			coinConstraint := e.GetCoinConstraint(c)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       c.ID,
					Coin:         c,
					ExSymbol:     key,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = key
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
// this require an API key
func (e *Bcex) GetPairsData() error {
	jsonResponse := &JsonResponse{}
	pairsData := &PairsData{}

	strRequestUrl := "/api_market/getTradeLists"
	strUrl := API_URL + strRequestUrl

	mapParams := make(map[string]string)
	mapParams["api_key"] = e.API_KEY

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if jsonResponse.Code != 0 {
		return fmt.Errorf("%s Get Pairs Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Data, &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, v := range pairsData.Main {
		for _, data := range v {
			p := &pair.Pair{}
			pairSymbol := data.Token + data.Market
			switch e.Source {
			case exchange.EXCHANGE_API:
				base := coin.GetCoin(data.Market)
				target := coin.GetCoin(data.Token)
				if base != nil && target != nil {
					p = pair.GetPair(base, target)
				}
			case exchange.JSON_FILE:
				p = e.GetPairBySymbol(pairSymbol)
			}

			lotsize, err := strconv.Atoi(data.NPrecision)
			//math.Pow10(-1 *lotsize )
			if err != nil {
				log.Printf(" Lot_Size Err: %s\n", err)
			}
			ticksize, err := strconv.Atoi(data.PPrecision)
			//math.Pow10(-1 * ticksize)
			if err != nil {
				log.Printf(" Tick_Size Err: %s\n", err)
			}

			if p != nil {
				pairConstraint := e.GetPairConstraint(p)
				if pairConstraint == nil {
					pairConstraint = &exchange.PairConstraint{
						PairID:      p.ID,
						Pair:        p,
						ExSymbol:    pairSymbol,
						MakerFee:    DEFAULT_MAKER_FEE,
						TakerFee:    DEFAULT_TAKER_FEE,
						LotSize:     math.Pow10(-1 * lotsize),
						PriceFilter: math.Pow10(-1 * ticksize),
						Listed:      true,
					}
				} else {
					pairConstraint.ExSymbol = pairSymbol
					pairConstraint.LotSize = math.Pow10(-1 * lotsize)
					pairConstraint.PriceFilter = math.Pow10(-1 * ticksize)
				}

				e.SetPairConstraint(pairConstraint)
			}
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
func (e *Bcex) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	jsonResponse := &JsonResponse{}
	orderBook := OrderBook{}

	strRequestUrl := "/api_market/market/depth"
	strUrl := API_URL + strRequestUrl

	mapParams := make(map[string]string)
	mapParams["market"] = pair.Base.Code
	mapParams["token"] = pair.Target.Code

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if jsonResponse.Code != 0 {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	var err error
	for _, bid := range orderBook.Bids {
		var buydata exchange.Order
		buydata.Rate, err = strconv.ParseFloat(bid[0], 64)
		if err != nil {
			log.Printf("BCEX Bids OrderBook strconv Rate error:%v", err)
			return nil, err
		}

		buydata.Quantity, err = strconv.ParseFloat(bid[1], 64)
		if err != nil {
			log.Printf("BCEX Bids OrderBook strconv Quantity error:%v", err)
			return nil, err
		}

		maker.Bids = append(maker.Bids, buydata)
	}
	for i := len(orderBook.Asks) - 1; i >= 0; i-- {
		var selldata exchange.Order
		selldata.Rate, err = strconv.ParseFloat(orderBook.Asks[i][0], 64)
		if err != nil {
			log.Printf("BCEX Asks OrderBook strconv Rate error:%v", err)
			return nil, err
		}

		selldata.Quantity, err = strconv.ParseFloat(orderBook.Asks[i][1], 64)
		if err != nil {
			log.Printf("BCEX Asks OrderBook strconv Quantity error:%v", err)
			return nil, err
		}

		maker.Asks = append(maker.Asks, selldata)
	}

	return maker, nil
}

func (e *Bcex) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *Bcex) DoAccountOperation(operation *exchange.AccountOperation) error {
	return nil
}

func (e *Bcex) GetCoinList() []string {
	jsonResponse := &JsonResponse{}
	coinsData := make(map[string]*CoinsData)

	strRequestUrl := "/api_market/getTokenPrecision"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		log.Printf("%s Get Coins List Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
		return nil
	} else if jsonResponse.Code != 0 {
		log.Printf("%s Get Coins List Failed: %v", e.GetName(), jsonResponse.Message)
		return nil
	}
	if err := json.Unmarshal(jsonResponse.Data, &coinsData); err != nil {
		log.Printf("%s Get Coins List Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return nil
	}

	var list []string
	for coinName := range coinsData {
		list = append(list, coinName)
	}

	return list
}

func (e *Bcex) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalances{}
	strRequest := "/api_market/getBalance"

	list := e.GetCoinList()

	for i := 0; i < len(list); i = i + 20 {
		mapParams := make(map[string]interface{})
		mapParams["page"] = "1"
		mapParams["size"] = "20"

		if i+20 > len(list) {
			mapParams["tokens"] = list[i:]
		} else {
			mapParams["tokens"] = list[i : i+20]
		}

		jsonBalanceReturn := e.ApiKeyPost(strRequest, mapParams)
		if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
			log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
			return
		} else if jsonResponse.Code != 0 {
			log.Printf("%s UpdateAllBalances Failed: %v", e.GetName(), jsonResponse.Message)
			return
		}
		if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
			log.Printf("%s UpdateAllBalances Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
			return
		}

		for _, v := range accountBalance.Data {
			c := e.GetCoinBySymbol(v.Token)
			if c != nil {
				freeamount, err := strconv.ParseFloat(v.Usable, 64)
				if err != nil {
					log.Printf("%s parse balance Err: %v %s", e.GetName(), err, v.Usable)
				}
				balanceMap.Set(c.Code, freeamount)
			}
		}
	}
}

func (e *Bcex) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {

	return false
}

func (e *Bcex) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/api_market/placeOrder"

	mapParams := make(map[string]interface{})
	mapParams["market_type"] = "1" //e.GetSymbolByPair(pair)
	mapParams["market"] = e.GetSymbolByCoin(pair.Base)
	mapParams["token"] = e.GetSymbolByCoin(pair.Target)
	mapParams["type"] = "2"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', -1, 64)

	jsonPlaceReturn := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != 0 {
		return nil, fmt.Errorf("%s LimitSell Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%v", placeOrder.ID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Bcex) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/api_market/placeOrder"

	mapParams := make(map[string]interface{})
	mapParams["market_type"] = "1" //e.GetSymbolByPair(pair)
	mapParams["market"] = e.GetSymbolByCoin(pair.Base)
	mapParams["token"] = e.GetSymbolByCoin(pair.Target)
	mapParams["type"] = "1"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', -1, 64)

	jsonPlaceReturn := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != 0 {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%v", placeOrder.ID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Bcex) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	orderStatus := PlaceOrder{}
	strRequest := "/api_market/getOrderByOrderNo"

	mapParams := make(map[string]interface{})
	mapParams["order_no"] = order.OrderID

	jsonOrderStatus := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if jsonResponse.Code != 0 {
		return fmt.Errorf("%s OrderStatus Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order.StatusMessage = jsonOrderStatus
	if orderStatus.Status == 0 {
		order.Status = exchange.Cancelled
	} else if orderStatus.Status == 1 {
		order.Status = exchange.New
	} else if orderStatus.Status == 2 {
		order.Status = exchange.Partial
	} else if orderStatus.Status == 3 {
		order.Status = exchange.Filled
	} else {
		order.Status = exchange.Other
	}

	return nil
}

func (e *Bcex) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Bcex) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	strRequest := "/api_market/cancelOrder"

	mapParams := make(map[string]interface{})
	mapParams["order_nos"] = fmt.Sprintf("[\"%v\"]", order.OrderID)

	jsonCancelOrder := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if jsonResponse.Code != 0 {
		return fmt.Errorf("%s CancelOrder Failed: %v", e.GetName(), jsonResponse.Message)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Bcex) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Bcex) ApiKeyGET(strRequestPath string, mapParams map[string]string) string {
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

func (e *Bcex) ApiKeyPost(strRequestPath string, mapParams map[string]interface{}) string {
	strUrl := API_URL + strRequestPath

	//Signature Request Params
	mapParams["api_key"] = e.API_KEY

	jsonParams := ""
	if nil != mapParams {
		bytesParams, _ := json.Marshal(mapParams)
		jsonParams = string(bytesParams)
	}

	mapParams["sign"] = ComputeSHA1(jsonParams, e.API_SECRET)

	bytesParams, _ := json.Marshal(mapParams)

	// 构建Request, 并且按官方要求添加Http Header
	request, err := http.NewRequest("POST", strUrl, bytes.NewBuffer(bytesParams))
	if nil != err {
		return err.Error()
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	request.Header.Add("Content-Type", "application/json")

	// 发出请求
	httpClient := &http.Client{}
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

func ComputeSHA1(mapParamsJson string, secretKey string) string {
	hasher := sha1.New()
	hasher.Write([]byte(mapParamsJson))

	decodeSecret, rest := pem.Decode([]byte(secretKey))
	if decodeSecret == nil {
		log.Printf("Decode Secret Err: %v", string(rest))
		return ""
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(decodeSecret.Bytes)
	if err != nil {
		log.Printf("Signature Err: %v", err)
		return ""
	}

	decodeSign, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA1, hasher.Sum(nil))
	if err != nil {
		log.Printf("%v", err)
	}

	return url.QueryEscape(base64.StdEncoding.EncodeToString(decodeSign))
}
