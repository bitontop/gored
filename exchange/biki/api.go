package biki

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
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
)

const (
	API_URL string = "https://api.biki.com"
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
func (e *Biki) GetCoinsData() error {
	jsonResponse := &JsonResponse{}
	pairsData := PairsData{}

	strRequestUrl := "/open/api/common/symbols"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if jsonResponse.Code != "0" {
		return fmt.Errorf("%s Get Coins Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Result, &pairsData); err != nil {
		return fmt.Errorf("%s Get Coins Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	for _, data := range pairsData {
		base := &coin.Coin{}
		target := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base = coin.GetCoin(data.CountCoin)
			if base == nil {
				base = &coin.Coin{}
				base.Code = data.CountCoin
				coin.AddCoin(base)
			}
			target = coin.GetCoin(data.BaseCoin)
			if target == nil {
				target = &coin.Coin{}
				target.Code = data.BaseCoin
				coin.AddCoin(target)
			}
		case exchange.JSON_FILE:
			base = e.GetCoinBySymbol(data.CountCoin)
			target = e.GetCoinBySymbol(data.BaseCoin)
		}

		if base != nil {
			coinConstraint := e.GetCoinConstraint(base)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       base.ID,
					Coin:         base,
					ExSymbol:     strings.ToLower(data.CountCoin),
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = strings.ToLower(data.CountCoin)
			}
			e.SetCoinConstraint(coinConstraint)
		}

		if target != nil {
			coinConstraint := e.GetCoinConstraint(target)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       target.ID,
					Coin:         target,
					ExSymbol:     strings.ToLower(data.BaseCoin),
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = strings.ToLower(data.BaseCoin)
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
func (e *Biki) GetPairsData() error {
	jsonResponse := &JsonResponse{}
	pairsData := PairsData{}

	strRequestUrl := "/open/api/common/symbols"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if jsonResponse.Code != "0" {
		return fmt.Errorf("%s Get Pairs Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Result, &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	for _, data := range pairsData {
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(data.CountCoin)
			target := coin.GetCoin(data.BaseCoin)
			if base != nil && target != nil {
				p = pair.GetPair(base, target)
			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(data.Symbol)
		}
		if p != nil {
			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{
					PairID:      p.ID,
					Pair:        p,
					ExSymbol:    data.Symbol,
					MakerFee:    DEFAULT_MAKER_FEE,
					TakerFee:    DEFAULT_TAKER_FEE,
					LotSize:     math.Pow10(-1 * data.AmountPrecision),
					PriceFilter: math.Pow10(-1 * data.PricePrecision),
					Listed:      DEFAULT_LISTED,
				}
			} else {
				pairConstraint.ExSymbol = data.Symbol
				pairConstraint.LotSize = math.Pow10(-1 * data.AmountPrecision)
				pairConstraint.PriceFilter = math.Pow10(-1 * data.PricePrecision)
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
func (e *Biki) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	jsonResponse := &JsonResponse{}
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	strRequestUrl := "/open/api/market_dept"
	strUrl := API_URL + strRequestUrl

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol
	mapParams["type"] = "step0"

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if jsonResponse.Code != "0" {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Result, &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	var err error
	for _, bid := range orderBook.Tick.Bids {
		var buydata exchange.Order

		//Modify according to type and structure
		buydata.Rate = bid[0]
		if err != nil {
			return nil, err
		}
		buydata.Quantity = bid[1]
		if err != nil {
			return nil, err
		}

		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Tick.Asks {
		var selldata exchange.Order

		//Modify according to type and structure
		selldata.Rate = ask[0]
		if err != nil {
			return nil, err
		}
		selldata.Quantity = ask[1]
		if err != nil {
			return nil, err
		}

		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

func (e *Biki) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *Biki) DoAccountOperation(operation *exchange.AccountOperation) error {
	return nil
}

func (e *Biki) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalances{}
	strRequest := "/open/api/user/account"

	jsonBalanceReturn := e.ApiKeyGet(strRequest, make(map[string]string))
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if jsonResponse.Code != "0" {
		log.Printf("%s UpdateAllBalances Failed: %v", e.GetName(), jsonResponse.Message)
		return
	}
	if err := json.Unmarshal(jsonResponse.Result, &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
		return
	}

	for _, v := range accountBalance.CoinList {
		c := e.GetCoinBySymbol(v.Coin)
		if c != nil {
			freeamount, err := strconv.ParseFloat(v.Normal, 64)
			if err != nil {
				log.Printf("%s Balance parse error: %v, %v", e.GetName(), err, v.Normal)
				return
			}
			balanceMap.Set(c.Code, freeamount)
		}
	}
}

func (e *Biki) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {

	return false
}

func (e *Biki) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/open/api/create_order"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "SELL"
	mapParams["type"] = "1"
	mapParams["volume"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)

	jsonPlaceReturn := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != "0" {
		return nil, fmt.Errorf("%s LimitSell Failed: %v", e.GetName(), jsonPlaceReturn)
	}
	if err := json.Unmarshal(jsonResponse.Result, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%v", placeOrder.OrderID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Biki) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/open/api/create_order"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "BUY"
	mapParams["type"] = "1"
	mapParams["volume"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)

	jsonPlaceReturn := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != "0" {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v", e.GetName(), jsonPlaceReturn)
	}
	if err := json.Unmarshal(jsonResponse.Result, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%v", placeOrder.OrderID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Biki) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	orderStatus := OrderStatus{}
	strRequest := "/open/api/order_info"

	mapParams := make(map[string]string)
	mapParams["order_id"] = order.OrderID
	mapParams["symbol"] = e.GetSymbolByPair(order.Pair)

	jsonOrderStatus := e.ApiKeyGet(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if jsonResponse.Code != "0" {
		return fmt.Errorf("%s OrderStatus Failed: %v", e.GetName(), jsonResponse)
	}
	if err := json.Unmarshal(jsonResponse.Result, &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	order.StatusMessage = jsonOrderStatus
	if orderStatus.OrderInfo.Status == 0 {
		order.Status = exchange.New
	} else if orderStatus.OrderInfo.Status == 1 {
		order.Status = exchange.New
	} else if orderStatus.OrderInfo.Status == 2 {
		order.Status = exchange.Filled
	} else if orderStatus.OrderInfo.Status == 3 {
		order.Status = exchange.Partial
	} else if orderStatus.OrderInfo.Status == 4 {
		order.Status = exchange.Cancelled
	} else if orderStatus.OrderInfo.Status == 5 {
		order.Status = exchange.Canceling
	} else if orderStatus.OrderInfo.Status == 6 {
		order.Status = exchange.Expired
	} else {
		order.Status = exchange.Other
	}

	return nil
}

func (e *Biki) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Biki) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	strRequest := "/open/api/cancel_order"

	mapParams := make(map[string]string)
	mapParams["order_id"] = order.OrderID
	mapParams["symbol"] = e.GetSymbolByPair(order.Pair)

	jsonCancelOrder := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if jsonResponse.Code != "0" {
		return fmt.Errorf("%s CancelOrder Failed: %v", e.GetName(), jsonResponse.Message)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Biki) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Biki) ApiKeyGet(strRequestPath string, mapParams map[string]string) string {

	strURL := API_URL + strRequestPath
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)

	mapParams["api_key"] = e.API_KEY
	mapParams["time"] = timeStamp

	stringData := MapSortByKey(mapParams) + e.API_SECRET
	sign := exchange.ComputeMD5(stringData)

	mapParams["sign"] = sign

	return exchange.HttpGetRequest(strURL, mapParams)
}

func (e *Biki) ApiKeyPost(strRequestPath string, mapParams map[string]string) string {

	//create url and http client
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	strURL := API_URL + strRequestPath
	httpClient := &http.Client{}
	postValues := url.Values{}

	mapParams["api_key"] = e.API_KEY
	mapParams["time"] = timeStamp

	stringData := MapSortByKey(mapParams) + e.API_SECRET
	hash := exchange.ComputeMD5(stringData)

	mapParams["sign"] = hash

	//stores mapParam value in postvalue
	for key, value := range mapParams {
		postValues.Set(key, value)
	}

	//encode into json format
	queryString := postValues.Encode()
	request, err := http.NewRequest("POST", strURL, strings.NewReader(queryString))
	if err != nil {
		log.Println(err)
	}

	//add the correct header
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	//make request
	postRespond, err := httpClient.Do(request)
	if err != nil {
		log.Println(err)
	}

	defer postRespond.Body.Close()
	body, err := ioutil.ReadAll(postRespond.Body)
	if nil != err {
		return err.Error()
	}

	//return request body
	return string(body)
}

func MapSortByKey(mapValue map[string]string) string {
	keys := make([]string, 0, len(mapValue))
	for key := range mapValue {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	mapParams := ""
	for _, key := range keys {
		mapParams += (key + url.QueryEscape(mapValue[key]))
	}
	//mapParams = mapParams[:len(mapParams)-1]
	//log.Printf("payload here, %v", mapParams)
	return mapParams
}
