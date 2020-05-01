package tokok

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
	API_URL string = "https://www.tokok.com/api/v1"
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
func (e *Tokok) GetCoinsData() error {
	pairsData := PairsData{}

	strRequestUrl := "/exchangeInfo"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &pairsData); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	}

	for _, data := range pairsData {
		base := &coin.Coin{}
		target := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base = coin.GetCoin(data.QuoteAsset)
			if base == nil {
				base = &coin.Coin{}
				base.Code = data.QuoteAsset
				coin.AddCoin(base)
			}
			target = coin.GetCoin(data.BaseAsset)
			if target == nil {
				target = &coin.Coin{}
				target.Code = data.BaseAsset
				coin.AddCoin(target)
			}
		case exchange.JSON_FILE:
			base = e.GetCoinBySymbol(data.QuoteAsset)
			target = e.GetCoinBySymbol(data.BaseAsset)
		}

		if base != nil {
			coinConstraint := e.GetCoinConstraint(base)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       base.ID,
					Coin:         base,
					ExSymbol:     data.QuoteAsset,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = data.QuoteAsset
			}
			e.SetCoinConstraint(coinConstraint)
		}

		if target != nil {
			coinConstraint := e.GetCoinConstraint(target)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       target.ID,
					Coin:         target,
					ExSymbol:     data.BaseAsset,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = data.BaseAsset
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
func (e *Tokok) GetPairsData() error {
	pairsData := PairsData{}

	strRequestUrl := "/exchangeInfo"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	}

	for _, data := range pairsData {
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(data.QuoteAsset)
			target := coin.GetCoin(data.BaseAsset)
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
					ExSymbol:    strings.ToLower(data.Symbol),
					MakerFee:    DEFAULT_MAKER_FEE,
					TakerFee:    DEFAULT_TAKER_FEE,
					LotSize:     math.Pow10(-1 * data.BaseAssetPrecision),
					PriceFilter: math.Pow10(-1 * data.QuoteAssetPrecision),
					Listed:      DEFAULT_LISTED,
				}
			} else {
				pairConstraint.ExSymbol = strings.ToLower(data.Symbol)
				pairConstraint.LotSize = math.Pow10(-1 * data.BaseAssetPrecision)
				pairConstraint.PriceFilter = math.Pow10(-1 * data.QuoteAssetPrecision)
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
func (e *Tokok) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	strRequestUrl := fmt.Sprintf("/depth?symbol=%s", symbol)
	strUrl := API_URL + strRequestUrl

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, nil)
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
			return nil, err
		}
		buydata.Quantity, err = strconv.ParseFloat(bid[1].(string), 64)
		if err != nil {
			return nil, err
		}

		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Asks {
		var selldata exchange.Order

		//Modify according to type and structure
		selldata.Rate, err = strconv.ParseFloat(ask[0].(string), 64)
		if err != nil {
			return nil, err
		}
		selldata.Quantity, err = strconv.ParseFloat(ask[1].(string), 64)
		if err != nil {
			return nil, err
		}

		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

func (e *Tokok) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *Tokok) DoAccountOperation(operation *exchange.AccountOperation) error {
	return nil
}

func (e *Tokok) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := JsonResponse{}
	accountBalance := AccountBalances{}
	strRequest := "/accounts"

	jsonBalanceReturn := e.ApiKeyRequest("POST", make(map[string]interface{}), strRequest)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if !jsonResponse.Result {
		log.Printf("%s UpdateAllBalances Failed: %v", e.GetName(), jsonResponse)
		return
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return
	}

	for _, v := range accountBalance {
		c := e.GetCoinBySymbol(v.CoinCode)
		if c != nil {
			freeamount, err := strconv.ParseFloat(v.HotMoney, 64)
			if err != nil {
				log.Printf("%s UpdateAllBalances Failed: %v", e.GetName(), v.HotMoney)
				return
			}
			balanceMap.Set(c.Code, freeamount)
		}
	}
}

func (e *Tokok) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	return false
}

func (e *Tokok) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := JsonResponse{}
	placeOrder := ""
	strRequest := "/trade"

	price := float64(int(rate/e.GetPriceFilter(pair)+e.GetPriceFilter(pair)/10)) * (e.GetPriceFilter(pair))
	amount := float64(int(quantity/e.GetLotSize(pair)+e.GetLotSize(pair)/10)) * (e.GetLotSize(pair))

	mapParams := make(map[string]interface{})
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["type"] = "2"
	mapParams["entrustPrice"] = price
	mapParams["entrustCount"] = amount
	// mapParams["openTok"] = 0

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	// log.Printf("====Sell return: %+v", jsonPlaceReturn)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if !jsonResponse.Result {
		return nil, fmt.Errorf("%s LimitSell Failed: %v", e.GetName(), jsonPlaceReturn)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      placeOrder,
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Tokok) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := JsonResponse{}
	placeOrder := ""
	strRequest := "/trade"

	price := float64(int(rate/e.GetPriceFilter(pair)+e.GetPriceFilter(pair)/10)) * (e.GetPriceFilter(pair))
	amount := float64(int(quantity/e.GetLotSize(pair)+e.GetLotSize(pair)/10)) * (e.GetLotSize(pair))

	mapParams := make(map[string]interface{})
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["type"] = "1"
	mapParams["entrustPrice"] = price
	mapParams["entrustCount"] = amount
	// mapParams["openTok"] = 0

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	// log.Printf("====Buy return: %+v", jsonPlaceReturn)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if !jsonResponse.Result {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v", e.GetName(), jsonPlaceReturn)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      placeOrder,
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Tokok) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := JsonResponse{}
	orderStatus := OrderStatus{}
	strRequest := "/order/orderInfo"

	mapParams := make(map[string]interface{})
	mapParams["order_id"] = order.OrderID

	jsonOrderStatus := e.ApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if !jsonResponse.Result {
		return fmt.Errorf("%s OrderStatus Failed: %v", e.GetName(), jsonResponse)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order.StatusMessage = jsonOrderStatus
	if orderStatus.Status == 4 {
		order.Status = exchange.Cancelled
	} else if orderStatus.Status == 2 {
		order.Status = exchange.Filled
	} else if orderStatus.Status == 1 || orderStatus.Status == 3 {
		order.Status = exchange.Partial
	} else if orderStatus.Status == 0 || orderStatus.Status == 1 {
		order.Status = exchange.New
	}

	return nil
}

func (e *Tokok) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Tokok) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	cancelOrder := CancelOrder{}
	strRequest := "/cancelEntrust"

	mapParams := make(map[string]interface{})
	mapParams["order_id"] = order.OrderID

	jsonCancelOrder := e.ApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if !cancelOrder.Result {
		return fmt.Errorf("%s CancelOrder Failed: %v", e.GetName(), cancelOrder)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Tokok) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Tokok) ApiKeyRequest(strMethod string, mapParams map[string]interface{}, strRequestPath string) string {
	nonce := fmt.Sprintf("%v", time.Now().Unix())

	request := &http.Request{}
	var strRequestUrl string
	var err error

	if strMethod == "GET" {
		strParams := exchange.Map2UrlQueryInterface(mapParams)
		strRequestUrl = strRequestPath + "?" + strParams
		request, err = http.NewRequest(strMethod, strRequestUrl, nil)
	} else {
		jsonParams := ""
		if nil != mapParams {
			bytesParams, _ := json.Marshal(mapParams)
			jsonParams = string(bytesParams)
		}
		payload := exchange.Map2UrlQueryInterface(mapParams)
		strRequestUrl = API_URL + strRequestPath + "?" + payload
		request, err = http.NewRequest(strMethod, strRequestUrl, strings.NewReader(jsonParams))
		// log.Printf("====jsonParams: %v", jsonParams)
		// log.Printf("====payload: %v", payload)
	}

	if nil != err {
		return err.Error()
	}

	mapParams["ACCESS-KEY"] = e.API_KEY
	mapParams["ACCESS-TIMESTAMP"] = nonce
	createSign := exchange.Map2UrlQueryInterface(mapParams)
	request.Header.Add("ACCESS-SIGN", exchange.ComputeHmac256Base64(createSign, e.API_SECRET))
	request.Header.Add("ACCESS-KEY", e.API_KEY)
	request.Header.Add("ACCESS-TIMESTAMP", nonce)
	request.Header.Add("Content-Type", "application/json")

	// log.Printf("====mapParams: %+v", mapParams)
	// log.Printf("====createSign: %v", createSign)

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

/* func (e *Tokok) ApiKeyRequest(strMethod string, mapParams map[string]interface{}, strRequestPath string) string {
	nonce := fmt.Sprintf("%v", time.Now().Unix())

	request := &http.Request{}
	var strRequestUrl string
	var err error

	if strMethod == "GET" {
		strParams := exchange.Map2UrlQueryInterface(mapParams)
		strRequestUrl = strRequestPath + "?" + strParams
		request, err = http.NewRequest(strMethod, strRequestUrl, nil)
	} else {
		jsonParams := ""
		if nil != mapParams {
			bytesParams, _ := json.Marshal(mapParams)
			jsonParams = string(bytesParams)
		}
		payload := exchange.Map2UrlQueryInterface(mapParams)
		strRequestUrl = API_URL + strRequestPath + "?" + payload
		request, err = http.NewRequest(strMethod, strRequestUrl, nil) //strings.NewReader(payload))
		// log.Printf("====jsonParams: %v", jsonParams)
		log.Printf("====payload: %v", payload)
	}

	if nil != err {
		return err.Error()
	}

	mapParams["ACCESS-KEY"] = e.API_KEY
	mapParams["ACCESS-TIMESTAMP"] = nonce
	createSign := exchange.Map2UrlQueryInterface(mapParams)
	request.Header.Add("ACCESS-SIGN", exchange.ComputeHmac256Base64(createSign, e.API_SECRET))
	request.Header.Add("ACCESS-KEY", e.API_KEY)
	request.Header.Add("ACCESS-TIMESTAMP", nonce)
	request.Header.Add("Content-Type", "application/json")

	log.Printf("====mapParams: %+v", mapParams)
	log.Printf("====createSign: %v", createSign)

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
} */

func (e *Tokok) ApiKeyGET(strRequestPath string, mapParams map[string]string) string {
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

/*
Error code Description
80101 Request frequency too high
80102 This IP is not allowed to access
82101 No account information
82102 Fail to get pair info
83001 Please enable SMS or GOOGLE authentication
83002 Please enable Real-name authentication
83103 Account blocked
83104 Account is prohibited from trading
83005 Trading pair can not be null
83006 Trading pair format error
83007 Incorrect order type
83008 Incorrect order price
83009 Incorrect order amount
83110 Place order failed
83011 Incorrect order ID
83112 Cancel order failed
83013 Order ID can not be null
83114 No order
83015 Required parameters of batch trade can not be null
83016 Exceed maximum order number
84001 Lack of timestamp
84002 Lack of signature
84003 Lack of parameters
84004 Signature error
84005 Invalid signature
84006 Lack of api_key
84007 Invalid API key
*/
