package bitatm

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
)

const (
	API_URL string = "https://open.bitatm.com"
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
func (e *BitATM) GetCoinsData() error {
	jsonResponse := &JsonResponse{}
	coinsData := CoinsData{}

	strRequestUrl := "/v1/common/currencies"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		log.Printf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if jsonResponse.Msg != "Ok." {
		log.Printf("%s Get Coins Failed: %v", e.GetName(), jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &coinsData); err != nil {
		log.Printf("%s Get Coins Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, data := range coinsData {
		c := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			c = coin.GetCoin(data.Currencyname)
			if c == nil {
				c = &coin.Coin{
					Code: data.Currencyname,
				}
				coin.AddCoin(c)
			}
		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(data.Currencyname)
		}

		if c != nil {
			coinConstraint := e.GetCoinConstraint(c)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       c.ID,
					Coin:         c,
					ExSymbol:     data.Currencyname,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = data.Currencyname
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
func (e *BitATM) GetPairsData() error {
	jsonResponse := &JsonResponse{}
	pairsData := PairsData{}

	strRequestUrl := "/v1/common/symbols"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &jsonResponse); err != nil {
		log.Printf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if jsonResponse.Msg != "Ok." {
		log.Printf("%s Get Pairs Failed: %v", e.GetName(), jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &pairsData); err != nil {
		log.Printf("%s Get Coins Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, data := range pairsData {
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(data.Basecurrency)
			target := coin.GetCoin(data.Quotecurrency)
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
					LotSize:     DEFAULT_LOT_SIZE,
					PriceFilter: DEFAULT_PRICE_FILTER,
					Listed:      DEFAULT_LISTED,
				}
			} else {
				pairConstraint.ExSymbol = data.Symbol
			}
			e.SetPairConstraint(pairConstraint)
		}
	}
	return nil
}

/*Get Pair Market Depth
Step 1: Change Instance Name    (e *<exchange Instance Name>)
*Step 2: Add Model of API Response
Step 3: Get Exchange Pair Code ex. symbol := e.GetPairCode(p)
Step 4: Modify API Path(strRequestUrl)
*Step 5: Add Params - Depend on API request
*Step 6: Convert the response to Standard Maker struct*/
func (e *BitATM) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	jsonResponse := &JsonResponse{}
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	strRequestUrl := "/market/depth"
	strUrl := API_URL + strRequestUrl + "?Symbol=" + symbol

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if jsonResponse.Code != "Ok." {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %v %s", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	for _, bid := range orderBook.Bids {
		var buydata exchange.Order

		//Modify according to type and structure
		buydata.Rate = bid.Price
		buydata.Quantity = bid.Amount
		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Asks {
		var selldata exchange.Order

		//Modify according to type and structure
		selldata.Rate = ask.Price
		selldata.Quantity = ask.Amount
		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

func (e *BitATM) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *BitATM) DoAccountOperation(operation *exchange.AccountOperation) error {
	return nil
}

func (e *BitATM) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalances{}
	strRequest := "/v1/account/balance"

	jsonBalanceReturn := e.ApiKeyGET(strRequest, make(map[string]interface{}))
	////////log.Printf(jsonBalanceReturn)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if jsonResponse.Code != "200" {
		log.Printf("%s UpdateAllBalances Failed: %v %s", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
		return
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return
	}

	for _, v := range accountBalance {
		c := e.GetCoinBySymbol(v.Currency)
		if c != nil {
			balanceMap.Set(c.Code, v.Balance)
		}
	}
}

func (e *BitATM) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil", e.GetName())
		return false
	}

	mapParams := make(map[string]interface{})
	mapParams["currency"] = e.GetSymbolByCoin(coin)
	mapParams["Amount"] = quantity
	mapParams["address"] = addr

	jsonResponse := &JsonResponse{}
	withdrawal := Withdrawal{}
	strRequest := "/v1/user/withdraw/create"

	jsonSubmitWithdraw := e.ApiKeyPOST(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &jsonResponse); err != nil {
		log.Printf("%s Withdraw Json Unmarshal Err: %v %v", e.GetName(), err, jsonSubmitWithdraw)
		return false
	} else if jsonResponse.Code != "200" {
		log.Printf("%s Withdraw Failed: %v %s", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
		return false
	}
	if err := json.Unmarshal(jsonResponse.Data, &withdrawal); err != nil {
		log.Printf("%s Withdraw Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return false
	}

	return true
}

func (e *BitATM) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	mapParams := make(map[string]interface{})
	mapParams["Symbol"] = e.GetSymbolByPair(pair)
	mapParams["Price"] = rate
	mapParams["Amount"] = quantity
	mapParams["OrderType"] = "Limit"
	mapParams["Direction"] = "Sell"

	jsonResponse := &JsonResponse{}
	strRequest := "/v1/order/create"
	var orderID int64

	jsonPlaceReturn := e.ApiKeyPOST(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != "200" {
		return nil, fmt.Errorf("%s LimitSell Failed: %v %s", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderID); err != nil {
		return nil, fmt.Errorf("%s LimitSell Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%d", orderID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *BitATM) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	mapParams := make(map[string]interface{})
	mapParams["Symbol"] = e.GetSymbolByPair(pair)
	mapParams["Price"] = rate
	mapParams["Amount"] = quantity
	mapParams["OrderType"] = "Limit"
	mapParams["Direction"] = "Buy"

	jsonResponse := &JsonResponse{}
	strRequest := "/v1/order/create"
	var orderID int64

	jsonPlaceReturn := e.ApiKeyPOST(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != "200" {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v %s", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderID); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%d", orderID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *BitATM) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	mapParams := make(map[string]interface{})
	orderID, _ := strconv.Atoi(order.OrderID)
	mapParams["OrderID"] = orderID

	jsonResponse := &JsonResponse{}
	orderStatus := OrderStatus{}
	strRequest := "/v1/order/detail"

	jsonOrderStatus := e.ApiKeyGET(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if jsonResponse.Code != "200" {
		return fmt.Errorf("%s OrderStatus Failed: %v %s", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order.StatusMessage = jsonOrderStatus
	if orderStatus.Orderstatus == 0 {
		order.Status = exchange.New
	} else if orderStatus.Orderstatus == 1 {
		order.Status = exchange.Partial
	} else if orderStatus.Orderstatus == 2 {
		order.Status = exchange.Cancelled
	} else if orderStatus.Orderstatus == 3 {
		order.Status = exchange.Filled
	} else if orderStatus.Orderstatus == 4 {
		order.Status = exchange.Canceling
	} else {
		order.Status = exchange.Other
	}

	return nil
}

func (e *BitATM) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *BitATM) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	mapParams := make(map[string]interface{})
	orderID, _ := strconv.Atoi(order.OrderID)
	mapParams["OrderID"] = orderID

	jsonResponse := &JsonResponse{}
	var cancelID int64
	strRequest := "/v1/order/cancel"

	jsonCancelOrder := e.ApiKeyPOST(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if jsonResponse.Code != "200" {
		return fmt.Errorf("%s CancelOrder Failed: %d %v", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &cancelID); err != nil {
		return fmt.Errorf("%s CancelOrder Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *BitATM) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *BitATM) ApiKeyGET(strRequestPath string, mapParams map[string]interface{}) string {
	mapParams["Accesskey"] = e.API_KEY
	mapParams["Randstr"] = fmt.Sprintf("%d", time.Now().Unix())
	mapParams["Timestamp"] = time.Now().Unix()
	payload := fmt.Sprintf("%s&secretkey=%s", exchange.Map2UrlQueryInterface(mapParams), e.API_SECRET)
	payload = url.QueryEscape(payload)
	payload = strings.ToLower(payload)
	mapParams["Signature"] = exchange.ComputeMD5(payload)

	url := exchange.Map2UrlQueryInterface(mapParams)
	strUrl := API_URL + strRequestPath + "?" + url

	request, err := http.NewRequest("GET", strUrl, nil)
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

func (e *BitATM) ApiKeyPOST(strRequestPath string, mapParams map[string]interface{}) string {
	//mapParamsSig := make(map[string]interface{})
	mapParams["Accesskey"] = e.API_KEY
	mapParams["Randstr"] = fmt.Sprintf("%d", time.Now().Unix())
	mapParams["Timestamp"] = time.Now().Unix()
	payload := fmt.Sprintf("%s&secretkey=%s", exchange.Map2UrlQueryInterface(mapParams), e.API_SECRET)
	payload = url.QueryEscape(payload)
	payload = strings.ToLower(payload)
	mapParams["Signature"] = exchange.ComputeMD5(payload)

	//url := exchange.Map2UrlQueryInterface(mapParams)

	strUrl := API_URL + strRequestPath

	// log.Printf("====payload: %v", payload)
	// log.Printf("====MD5: %v", exchange.ComputeMD5(payload))
	// log.Printf("====url: %v", url)

	jsonParams := ""
	if nil != mapParams {
		bytesParams, _ := json.Marshal(mapParams)
		jsonParams = string(bytesParams)
	}

	request, err := http.NewRequest("POST", strUrl, strings.NewReader(jsonParams))
	if nil != err {
		return err.Error()
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	//request.Header.Add("apisign", signature)

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

	log.Printf("====Post Request: %+v", request)
	log.Printf("====Post Return: %+v", body)

	return string(body)
}
