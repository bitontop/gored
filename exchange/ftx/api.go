package ftx

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"

	"strconv"
)

const (
	API_URL string = "https://ftx.com/api"
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
func (e *Ftx) GetCoinsData() error {
	jsonResponse := &JsonResponse{}
	coinsData := CoinsData{}

	strRequestUrl := "/coins"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if !jsonResponse.Success {
		return fmt.Errorf("%s Get Coins Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Result, &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	for _, data := range coinsData {
		c := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			c = coin.GetCoin(data.ID)
			if c == nil {
				c = &coin.Coin{}
				c.Code = data.ID
				coin.AddCoin(c)
			}
		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(data.ID)
		}

		if c != nil {
			coinConstraint := e.GetCoinConstraint(c)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       c.ID,
					Coin:         c,
					ExSymbol:     data.ID,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = data.ID
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
func (e *Ftx) GetPairsData() error { //TODO
	jsonResponse := &JsonResponse{}
	coinsData := CoinsData{}

	strRequestUrl := "/coins"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if !jsonResponse.Success {
		return fmt.Errorf("%s Get Pairs Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Result, &coinsData); err != nil {
		return fmt.Errorf("%s Get Pairs Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	for _, data := range coinsData {
		if data.Underlying == "" {
			continue
		}
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin("USD")
			target := coin.GetCoin(data.ID)
			if base != nil && target != nil {
				p = pair.GetPair(base, target)
			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(data.Name)
		}
		if p != nil {
			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{
					PairID:      p.ID,
					Pair:        p,
					ExSymbol:    fmt.Sprintf("%v/USD", data.ID),
					MakerFee:    DEFAULT_MAKER_FEE,
					TakerFee:    DEFAULT_TAKER_FEE,
					LotSize:     DEFAULT_LOT_SIZE,
					PriceFilter: DEFAULT_PRICE_FILTER,
					Listed:      true,
				}
			} else {
				pairConstraint.ExSymbol = fmt.Sprintf("%v/USD", data.ID)
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

// orderbook TODO
func (e *Ftx) OrderBook(pair *pair.Pair) (*exchange.Maker, error) { // TODO
	jsonResponse := &JsonResponse{}
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	mapParams := make(map[string]string)
	mapParams["market"] = symbol
	mapParams["type"] = "both"

	strRequestUrl := "/v1.1/public/getorderbook"
	strUrl := API_URL + strRequestUrl

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if !jsonResponse.Success {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Result, &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	for _, bid := range orderBook.Buy {
		maker.Bids = append(maker.Bids, bid)
	}
	for _, ask := range orderBook.Sell {
		maker.Asks = append(maker.Asks, ask)
	}
	return maker, nil
}

func (e *Ftx) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/

func (e *Ftx) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalances{}
	strRequest := "/wallet/balances" // "/account"

	jsonBalanceReturn := e.ApiKeyRequest("GET", strRequest, make(map[string]string))
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if !jsonResponse.Success {
		log.Printf("%s UpdateAllBalances Failed: %v", e.GetName(), jsonBalanceReturn)
		return
	}
	if err := json.Unmarshal(jsonResponse.Result, &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
		return
	}

	for _, v := range accountBalance {
		c := e.GetCoinBySymbol(v.Coin)
		if c != nil {
			balanceMap.Set(c.Code, v.Free)
		}
	}
}

// TODO
func (e *Ftx) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil", e.GetName())
		return false
	}

	mapParams := make(map[string]string)
	mapParams["currency"] = e.GetSymbolByCoin(coin)
	mapParams["quantity"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["address"] = addr

	jsonResponse := &JsonResponse{}
	uuid := Uuid{}
	strRequest := "/v1.1/account/withdraw"

	jsonSubmitWithdraw := e.ApiKeyGET(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &jsonResponse); err != nil {
		log.Printf("%s Withdraw Json Unmarshal Err: %v %v", e.GetName(), err, jsonSubmitWithdraw)
		return false
	} else if !jsonResponse.Success {
		log.Printf("%s Withdraw Failed: %v", e.GetName(), jsonResponse.Message)
		return false
	}
	if err := json.Unmarshal(jsonResponse.Result, &uuid); err != nil {
		log.Printf("%s Withdraw Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
		return false
	}
	return true
}

func (e *Ftx) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	mapParams := make(map[string]string)
	mapParams["market"] = e.GetSymbolByPair(pair)
	mapParams["quantity"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["rate"] = strconv.FormatFloat(rate, 'f', -1, 64)

	jsonResponse := &JsonResponse{}
	uuid := Uuid{}
	strRequest := "/orders"

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if !jsonResponse.Success {
		return nil, fmt.Errorf("%s LimitSell Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Result, &uuid); err != nil {
		return nil, fmt.Errorf("%s LimitSell Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      uuid.Id,
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Ftx) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	mapParams := make(map[string]string)
	mapParams["market"] = e.GetSymbolByPair(pair)
	mapParams["quantity"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["rate"] = strconv.FormatFloat(rate, 'f', -1, 64)

	jsonResponse := &JsonResponse{}
	uuid := Uuid{}
	strRequest := "/v1.1/market/buylimit"

	jsonPlaceReturn := e.ApiKeyGET(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if !jsonResponse.Success {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Result, &uuid); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      uuid.Id,
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}
	return order, nil
}

func (e *Ftx) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	mapParams := make(map[string]string)
	mapParams["uuid"] = order.OrderID

	jsonResponse := &JsonResponse{}
	orderStatus := PlaceOrder{}
	strRequest := "/v1.1/account/getorder"

	jsonOrderStatus := e.ApiKeyGET(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if !jsonResponse.Success {
		return fmt.Errorf("%s OrderStatus Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Result, &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	order.StatusMessage = jsonOrderStatus
	if orderStatus.CancelInitiated {
		order.Status = exchange.Canceling
	} else if !orderStatus.IsOpen && orderStatus.QuantityRemaining > 0 {
		order.Status = exchange.Cancelled
	} else if orderStatus.QuantityRemaining == 0 {
		order.Status = exchange.Filled
	} else if orderStatus.QuantityRemaining != orderStatus.Quantity {
		order.Status = exchange.Partial
	} else {
		order.Status = exchange.New
	}

	return nil
}

func (e *Ftx) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Ftx) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	mapParams := make(map[string]string)
	mapParams["uuid"] = order.OrderID

	jsonResponse := &JsonResponse{}
	cancelOrder := PlaceOrder{}
	strRequest := "/v1.1/market/cancel"

	jsonCancelOrder := e.ApiKeyGET(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if !jsonResponse.Success {
		return fmt.Errorf("%s CancelOrder Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Result, &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Ftx) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Ftx) ApiKeyGET(strRequestPath string, mapParams map[string]string) string {
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

func (e *Ftx) ApiKeyRequest(strMethod, strRequestPath string, mapParams map[string]string) string {
	timestamp := fmt.Sprintf("%v", time.Now().UnixNano()/int64(time.Millisecond)) // Millisecond

	var err error
	request := &http.Request{}
	httpClient := &http.Client{}
	postBody := ""
	preSign := ""
	strSignUrl := ""
	strRequestUrl := ""

	if strMethod == "POST" {
		bytesParams, _ := json.Marshal(mapParams)
		postBody = string(bytesParams)
		strSignUrl = strRequestPath // ?
	} else if len(mapParams) != 0 {
		strSignUrl = strRequestPath + "?" + exchange.Map2UrlQuery(mapParams)
	} else {
		strSignUrl = strRequestPath
	}
	strRequestUrl = API_URL + strSignUrl

	// create signature
	preSign = timestamp + strMethod + strSignUrl
	signature := exchange.ComputeHmac256NoDecode(preSign, e.API_SECRET)

	log.Printf("postBody: %v", postBody)
	log.Printf("preSign: %v", preSign)
	log.Printf("strRequestUrl: %v", strRequestUrl)

	// request
	strRequestUrl = strRequestPath // ==========================
	request, err = http.NewRequest(strMethod, strRequestUrl, strings.NewReader(postBody))
	if nil != err {
		return err.Error()
	}

	request.Header.Add("FTX-KEY", e.API_KEY)
	request.Header.Add("FTX-TS", timestamp)
	request.Header.Add("FTX-SIGN", signature)
	// request.Header.Add("Content-Type", "application/json")
	// add FTX-SUBACCOUNT if using subaccount

	log.Printf("key: %v", e.API_KEY)
	log.Printf("secret: %v", e.API_SECRET)

	log.Printf("FTX-KEY: %v", e.API_KEY)
	log.Printf("FTX-TS: %v", timestamp)
	log.Printf("FTX-SIGN: %v", signature)

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
