package bitmart

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
	API_URL string = "https://openapi.bitmart.com"
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
func (e *Bitmart) GetCoinsData() error {
	coinsData := CoinsData{}

	strRequestUrl := "/v2/currencies"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	}

	for _, data := range coinsData {
		c := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			c = coin.GetCoin(data.ID)
			if c == nil {
				c = &coin.Coin{}
				c.Code = data.ID
				c.Name = data.Name
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
					Withdraw:     data.WithdrawEnabled,
					Deposit:      data.DepositEnabled,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = data.ID
				coinConstraint.Withdraw = data.WithdrawEnabled
				coinConstraint.Deposit = data.DepositEnabled
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
func (e *Bitmart) GetPairsData() error {
	pairsData := PairsData{}

	strRequestUrl := "/v2/symbols_details"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	}

	for _, data := range pairsData {
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(data.QuoteCurrency)
			target := coin.GetCoin(data.BaseCurrency)
			if base != nil && target != nil {

				p = pair.GetPair(base, target)

			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(data.ID)
		}
		if p != nil {
			lotSize, err := strconv.ParseFloat(data.QuoteIncrement, 64)
			if err != nil {
				return fmt.Errorf("%s lotSize parse error: %v, %v", e.GetName(), err, data.QuoteIncrement)
			}
			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{
					PairID:      p.ID,
					Pair:        p,
					ExSymbol:    data.ID,
					MakerFee:    DEFAULT_MAKER_FEE,
					TakerFee:    DEFAULT_TAKER_FEE,
					LotSize:     lotSize,
					PriceFilter: math.Pow10(-1 * data.PriceMaxPrecision),
					Listed:      DEFAULT_LISTED,
				}
			} else {
				pairConstraint.ExSymbol = data.ID
				pairConstraint.LotSize = lotSize
				pairConstraint.PriceFilter = math.Pow10(-1 * data.PriceMaxPrecision)
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
func (e *Bitmart) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	strRequestUrl := fmt.Sprintf("/v2/symbols/%s/orders", symbol)
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
	for _, bid := range orderBook.Buys {
		var buydata exchange.Order

		//Modify according to type and structure
		buydata.Rate, err = strconv.ParseFloat(bid.Price, 64)
		if err != nil {
			return nil, err
		}
		buydata.Quantity, err = strconv.ParseFloat(bid.Amount, 64)
		if err != nil {
			return nil, err
		}

		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Sells {
		var selldata exchange.Order

		//Modify according to type and structure
		selldata.Rate, err = strconv.ParseFloat(ask.Price, 64)
		if err != nil {
			return nil, err
		}
		selldata.Quantity, err = strconv.ParseFloat(ask.Amount, 64)
		if err != nil {
			return nil, err
		}

		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

func (e *Bitmart) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *Bitmart) DoAccountOperation(operation *exchange.AccountOperation) error {
	return nil
}

func (e *Bitmart) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	accountBalance := AccountBalances{}
	strRequest := "/v2/wallet"

	jsonBalanceReturn := e.ApiKeyRequest("GET", strRequest, nil)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	}

	for _, balance := range accountBalance {
		c := e.GetCoinBySymbol(balance.ID)
		if c != nil {
			freeAmount, err := strconv.ParseFloat(balance.Available, 64)
			if err != nil {
				log.Printf("%s balance parse Err: %v %v", e.GetName(), err, balance.Available)
				return
			}
			balanceMap.Set(c.Code, freeAmount)
		}

	}

}

func (e *Bitmart) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {

	return false
}

func (e *Bitmart) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	placeOrder := PlaceOrder{}
	strRequest := "/v2/orders"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "sell"
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if placeOrder.Message != "" {
		return nil, fmt.Errorf("%s LimitSell Failed: %v", e.GetName(), jsonPlaceReturn)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%v", placeOrder.EntrustID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Bitmart) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	placeOrder := PlaceOrder{}
	strRequest := "/v2/orders"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "buy"
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if placeOrder.Message != "" {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v", e.GetName(), jsonPlaceReturn)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%v", placeOrder.EntrustID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Bitmart) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	orderStatus := OrderStatus{}
	strRequest := fmt.Sprintf("/v2/orders/%s", order.OrderID)

	mapParams := make(map[string]string)
	mapParams["entrust_id"] = order.OrderID

	jsonOrderStatus := e.ApiKeyRequest("GET", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	}

	order.StatusMessage = jsonOrderStatus
	if orderStatus.Status == 4 {
		order.Status = exchange.Cancelled
	} else if orderStatus.Status == 0 {
		order.Status = exchange.Other
	} else if orderStatus.Status == 3 {
		order.Status = exchange.Filled
	} else if orderStatus.Status == 5 || orderStatus.Status == 2 {
		order.Status = exchange.Partial
	} else if orderStatus.Status == 1 {
		order.Status = exchange.New
	} else {
		order.Status = exchange.Other
	}

	return nil
}

func (e *Bitmart) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Bitmart) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	strRequest := fmt.Sprintf("/v2/orders/%s", order.OrderID)

	mapParams := make(map[string]string)
	mapParams["entrust_id"] = order.OrderID

	jsonCancelOrder := e.ApiKeyRequest("DELETE", strRequest, mapParams)
	if jsonCancelOrder != "{}" {
		return fmt.Errorf("%s CancelOrder Failed: %v", e.GetName(), jsonCancelOrder)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Bitmart) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Bitmart) ApiKeyRequest(strMethod string, strRequestPath string, mapParams map[string]string) string {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	strUrl := API_URL + strRequestPath

	request := &http.Request{}
	var err error

	if strMethod == "POST" || strMethod == "DELETE" {
		jsonParams := ""
		if nil != mapParams {
			bytesParams, _ := json.Marshal(mapParams)
			jsonParams = string(bytesParams)
		}
		request, err = http.NewRequest(strMethod, strUrl, strings.NewReader(jsonParams))
	} else if strMethod == "GET" {
		if mapParams != nil {
			strParams := exchange.Map2UrlQuery(mapParams)
			strUrl = strUrl + "?" + strParams
		}
		request, err = http.NewRequest(strMethod, strUrl, nil)
	}

	if nil != err {
		return err.Error()
	}

	if strMethod != "GET" {
		request.Header.Add("X-BM-SIGNATURE", exchange.ComputeHmac256NoDecode(exchange.Map2UrlQueryUrl(mapParams), e.API_SECRET))
	}

	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-BM-TIMESTAMP", fmt.Sprintf("%v", timestamp))
	request.Header.Add("X-BM-AUTHORIZATION", "Bearer "+e.GetToken(e.API_KEY, e.API_SECRET, e.Passphrase))

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

func (e *Bitmart) GetToken(key string, secret string, memo string) string {
	token := exchange.ComputeHmac256NoDecode((key + ":" + secret + ":" + memo), secret)

	mapParams := make(map[string]string)
	mapParams["grant_type"] = "client_credentials"
	mapParams["client_id"] = key
	mapParams["client_secret"] = token

	accessToken := AccessToken{}
	strRequest := "https://openapi.bitmart.com/v2/authentication"

	jsonBitmart := e.TokenReq(strRequest, mapParams)
	err := json.Unmarshal([]byte(jsonBitmart), &accessToken)
	if err != nil {
		log.Printf("Create AccessToken json unmarshal error : %v", jsonBitmart)
	}
	token = accessToken.AccessToken

	return token
}

func (e *Bitmart) TokenReq(resource string, mapParams map[string]string) string {

	req, err := http.NewRequest("POST", resource, strings.NewReader(exchange.Map2UrlQuery(mapParams)))
	if err != nil {
		return err.Error()
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	// 发出请求
	httpClient := &http.Client{}
	response, err := httpClient.Do(req)
	if err != nil {
		log.Printf("err=%v", err)
		return err.Error()
	}

	defer response.Body.Close()

	// 解析响应内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err.Error()
	}

	return string(body)
}
