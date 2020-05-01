package dcoin

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

/*The Base Endpoint URL*/
const (
	API_URL = "https://openapi.dcoin.com/open/api"
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
Step 3: Modify API Path(strRequestPath)*/
func (e *Dcoin) GetCoinsData() error {
	jsonResponse := &JsonResponse{}
	coinsData := PairsData{}

	strRequestPath := "/common/symbols"
	strUrl := API_URL + strRequestPath

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if jsonResponse.Code != 0 {
		return fmt.Errorf("%s Get Coins Failed: %v", e.GetName(), jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, data := range coinsData {
		base := &coin.Coin{}
		target := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base = coin.GetCoin(data.BaseCoin)
			if base == nil {
				base = &coin.Coin{}
				base.Code = data.BaseCoin
				coin.AddCoin(base)
			}
			target = coin.GetCoin(data.CountCoin)
			if target == nil {
				target = &coin.Coin{}
				target.Code = data.CountCoin
				coin.AddCoin(target)
			}
		case exchange.JSON_FILE:
			base = e.GetCoinBySymbol(data.BaseCoin)
			target = e.GetCoinBySymbol(data.CountCoin)
		}

		if base != nil {
			coinConstraint := e.GetCoinConstraint(base)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       base.ID,
					Coin:         base,
					ExSymbol:     data.BaseCoin,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = data.BaseCoin
			}
			e.SetCoinConstraint(coinConstraint)
		}

		if target != nil {
			coinConstraint := e.GetCoinConstraint(target)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       target.ID,
					Coin:         target,
					ExSymbol:     data.CountCoin,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = data.CountCoin
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
func (e *Dcoin) GetPairsData() error {
	jsonResponse := &JsonResponse{}
	pairsData := PairsData{}

	strRequestPath := "/common/symbols"
	strUrl := API_URL + strRequestPath

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if jsonResponse.Code != 0 {
		return fmt.Errorf("%s Get Pairs Failed: %v", e.GetName(), jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, data := range pairsData {
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(data.BaseCoin)
			target := coin.GetCoin(data.CountCoin)
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
					LotSize:     math.Pow10(1 * data.AmountPrecision),
					PriceFilter: math.Pow10(1 * data.PricePrecision),
					Listed:      true,
				}
			} else {
				pairConstraint.ExSymbol = data.Symbol
				pairConstraint.LotSize = math.Pow10(1 * data.AmountPrecision)
				pairConstraint.PriceFilter = math.Pow10(1 * data.PricePrecision)
			}
			e.SetPairConstraint(pairConstraint)
		}

	}
	return nil
}

/*Get Pair Market Depth
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Get Exchange Pair Code ex. symbol := e.GetSymbolByPair(p)
Step 4: Modify API Path(strRequestUrl)
Step 5: Add Params - Depend on API request
Step 6: Convert the response to Standard Maker struct*/
func (e *Dcoin) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	jsonResponse := &JsonResponse{}
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol
	mapParams["type"] = "step0"

	strRequestPath := "/market_dept"
	strUrl := API_URL + strRequestPath

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if jsonResponse.Code != 0 {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %v", e.GetName(), jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)

	var err error
	for _, bid := range orderBook.Bids {
		buydata := exchange.Order{}
		buydata.Quantity = bid[1]
		buydata.Rate = bid[0]
		maker.Bids = append(maker.Bids, buydata)
	}

	for _, ask := range orderBook.Asks {
		selldata := exchange.Order{}
		selldata.Quantity = ask[1]
		selldata.Rate = ask[0]
		maker.Asks = append(maker.Asks, selldata)
	}

	return maker, err
}

func (e *Dcoin) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *Dcoin) DoAccountOperation(operation *exchange.AccountOperation) error {
	return nil
}

func (e *Dcoin) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalances{}

	strRequestPath := "/user/account"

	jsonBalanceReturn := e.ApiKeyGet(strRequestPath, make(map[string]interface{}))
	//log.Printf(jsonBalanceReturn)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if jsonResponse.Code != 0 {
		log.Printf("%s UpdateAllBalances Failed: %v", e.GetName(), jsonResponse.Msg)
		return
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return
	}

	for _, balance := range accountBalance.CoinList {
		c := e.GetCoinBySymbol(balance.Coin)
		if c != nil {
			balanceMap.Set(c.Code, balance.Normal)
		}
	}
}

/* Withdraw(coin *coin.Coin, quantity float64, addr, tag string) */
func (e *Dcoin) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	log.Printf("%s Withdraw Not Viable with API.", e.GetName())
	return false
}

func (e *Dcoin) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequestPath := "/create_order"

	mapParams := make(map[string]interface{})
	mapParams["side"] = "SELL"
	mapParams["type"] = 1
	mapParams["volume"] = quantity
	mapParams["price"] = rate
	mapParams["symbol"] = e.GetSymbolByPair(pair)

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequestPath, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != 0 {
		return nil, fmt.Errorf("%s LimitSell Failed: %v", e.GetName(), jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      strconv.Itoa(placeOrder.OrderID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}
	return order, nil
}

func (e *Dcoin) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequestPath := "/create_order"

	mapParams := make(map[string]interface{})
	mapParams["side"] = "BUY"
	mapParams["type"] = 1
	mapParams["volume"] = quantity
	mapParams["price"] = rate
	mapParams["symbol"] = e.GetSymbolByPair(pair)

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequestPath, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != 0 {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v", e.GetName(), jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      strconv.Itoa(placeOrder.OrderID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}
	return order, nil
}

func (e *Dcoin) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	orderStatus := OrderStatus{}
	strRequestPath := "/all_order"

	mapParams := make(map[string]interface{})
	mapParams["symbol"] = e.GetSymbolByPair(order.Pair)
	mapParams["order_id"] = order.OrderID

	jsonOrderStatus := e.ApiKeyGet(strRequestPath, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if jsonResponse.Code != 0 {
		return fmt.Errorf("%s OrderStatus Failed: %v", e.GetName(), jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	if orderStatus.OrderInfo.Status == 4 {
		order.Status = exchange.Cancelled
	} else if orderStatus.OrderInfo.Status == 5 {
		order.Status = exchange.Canceling
	} else if orderStatus.OrderInfo.Status == 2 {
		order.Status = exchange.Filled
	} else if orderStatus.OrderInfo.Status == 3 {
		order.Status = exchange.Partial
	} else if orderStatus.OrderInfo.Status == 6 {
		order.Status = exchange.Expired
	} else if orderStatus.OrderInfo.Status == 1 {
		order.Status = exchange.New
	} else {
		order.Status = exchange.Other
	}

	order.DealRate = orderStatus.OrderInfo.AgePrice
	order.DealQuantity = orderStatus.OrderInfo.DealVolume

	return nil
}

func (e *Dcoin) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Dcoin) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	strRequestPath := "/cancel_order"

	mapParams := make(map[string]interface{})
	orderID, err := strconv.ParseInt(order.OrderID, 10, 64)
	if err != nil {
		return fmt.Errorf("%s CancelOrder orderID parse fail: %v, %v", e.GetName(), err, order.OrderID)
	}
	mapParams["order_id"] = orderID
	mapParams["symbol"] = e.GetSymbolByPair(order.Pair)

	jsonCancelOrder := e.ApiKeyRequest("POST", strRequestPath, mapParams)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if jsonResponse.Code != 0 {
		return fmt.Errorf("%s CancelOrder Failed: %v", e.GetName(), jsonCancelOrder)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Dcoin) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Get Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Dcoin) ApiKeyGet(strRequestPath string, mapParams map[string]interface{}) string {
	mapParams["api_key"] = e.API_KEY
	payload := fmt.Sprintf("%s%s", createPayload(mapParams), e.API_SECRET)
	mapParams["sign"] = exchange.ComputeMD5(payload)

	url := Map2UrlQuery(mapParams)
	strUrl := API_URL + strRequestPath + "?" + url

	request, err := http.NewRequest("GET", strUrl, nil)
	if nil != err {
		return err.Error()
	}
	request.Header.Add("Content-Type", "application/json; charset=utf-8")

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

func Map2UrlQuery(mapParams map[string]interface{}) string {
	var strParams string
	mapSort := []string{}
	for key := range mapParams {
		mapSort = append(mapSort, key)
	}
	sort.Strings(mapSort)

	for _, key := range mapSort {
		strParams += (key + "=" + fmt.Sprintf("%v", mapParams[key]) + "&")
	}

	if 0 < len(strParams) {
		strParams = string([]rune(strParams)[:len(strParams)-1])
	}

	return strParams
}

/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request*/
func (e *Dcoin) ApiKeyRequest(strMethod, strRequestPath string, mapParams map[string]interface{}) string {
	strUrl := API_URL + strRequestPath

	mapParams["api_key"] = e.API_KEY
	payload := fmt.Sprintf("%s%s", createPayload(mapParams), e.API_SECRET)
	mapParams["sign"] = exchange.ComputeMD5(payload)

	strParams := Map2UrlQueryUrl(mapParams)

	request, err := http.NewRequest(strMethod, strUrl, strings.NewReader(strParams))
	if nil != err {
		return err.Error()
	}
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

func Map2UrlQueryUrl(mapParams map[string]interface{}) string {
	var strParams string
	mapSort := []string{}
	for key := range mapParams {
		mapSort = append(mapSort, key)
	}
	sort.Strings(mapSort)

	for _, key := range mapSort {
		strParams += fmt.Sprintf("%s=%v&", key, url.QueryEscape(fmt.Sprintf("%v", mapParams[key])))
	}

	if 0 < len(strParams) {
		strParams = string([]rune(strParams)[:len(strParams)-1])
	}

	return strParams
}

func createPayload(mapParams map[string]interface{}) string {
	var strParams string
	mapSort := []string{}
	for key := range mapParams {
		mapSort = append(mapSort, key)
	}
	sort.Strings(mapSort)

	for _, key := range mapSort {
		strParams += key + fmt.Sprintf("%v", mapParams[key])
	}

	return strParams
}
