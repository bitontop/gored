package huobi

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
	API_URL string = "https://api.huobi.pro"
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
func (e *Huobi) GetCoinsData() {
	jsonResponse := &JsonResponse{}
	coinsData := CoinsData{}

	//strRequestUrl := "/v1/common/currencys"
	strUrl := "https://www.huobi.com/-/x/pro/v1/settings/currencys?r=sqyeinryv8&language=en-US"

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		log.Printf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if jsonResponse.Status != "ok" {
		log.Printf("%s Get Coins Failed: %v", e.GetName(), jsonResponse)
	}
	if err := json.Unmarshal(jsonResponse.Data, &coinsData); err != nil {
		log.Printf("%s Get Coins Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, data := range coinsData {
		c := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			c = coin.GetCoin(data.DisplayName)
			if c == nil {
				c = &coin.Coin{}
				c.Code = data.DisplayName
				c.CurrencyType = data.CurrencyType
				coin.AddCoin(c)
			} else if c.CurrencyType == "" {
				c.CurrencyType = data.CurrencyType
			}
		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(data.DisplayName)
		}

		if c != nil {
			coinConstraint := &exchange.CoinConstraint{
				CoinID:       c.ID,
				Coin:         c,
				ExSymbol:     data.Name,
				TxFee:        DEFAULT_TXFEE,
				Withdraw:     data.WithdrawEnabled,
				Deposit:      data.DepositEnabled,
				Confirmation: data.SafeConfirms,
				Listed:       true,
			}
			e.SetCoinConstraint(coinConstraint)
		}
	}
}

/* GetPairsData - Get Pairs Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func (e *Huobi) GetPairsData() {
	jsonResponse := &JsonResponse{}
	pairsData := PairsData{}

	strRequestUrl := "/v1/common/symbols"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &jsonResponse); err != nil {
		log.Printf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if jsonResponse.Status != "ok" {
		log.Printf("%s Get Pairs Failed: %v", e.GetName(), jsonResponse)
	}
	if err := json.Unmarshal(jsonResponse.Data, &pairsData); err != nil {
		log.Printf("%s Get Pairs Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, data := range pairsData {
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(strings.ToUpper(data.QuoteCurrency))
			target := coin.GetCoin(strings.ToUpper(data.BaseCurrency))
			if base != nil && target != nil {
				p = pair.GetPair(base, target)
			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(data.Symbol)
		}
		if p != nil {
			pairConstraint := &exchange.PairConstraint{
				PairID:      p.ID,
				Pair:        p,
				ExSymbol:    data.Symbol,
				MakerFee:    DEFAULT_MAKERER_FEE,
				TakerFee:    DEFAULT_TAKER_FEE,
				LotSize:     math.Pow10(data.AmountPrecision * -1),
				PriceFilter: math.Pow10(data.PricePrecision * -1),
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
func (e *Huobi) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	jsonResponse := &JsonResponse{}
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	strRequestUrl := "/market/depth"
	strUrl := API_URL + strRequestUrl

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol
	mapParams["type"] = "step0"

	maker := &exchange.Maker{}
	maker.WorkerIP = exchange.GetExternalIP()
	maker.BeforeTimestamp = float64(time.Now().UnixNano() / 1e6)

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if jsonResponse.Status != "ok" {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %v", e.GetName(), jsonResponse)
	}
	if err := json.Unmarshal(jsonResponse.Tick, &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Tick)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	for _, bid := range orderBook.Bids {
		var buydata exchange.Order

		buydata.Rate = bid[0]
		buydata.Quantity = bid[1]

		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Asks {
		var selldata exchange.Order

		selldata.Rate = ask[0]
		selldata.Quantity = ask[1]

		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

/*************** Private API ***************/
func (e *Huobi) GetAccounts() string {
	jsonResponse := &JsonResponse{}
	accountsReturn := AccountsReturn{}

	strRequest := "/v1/account/accounts"

	jsonAccountsReturn := e.ApiKeyRequest("GET", make(map[string]string), strRequest)
	if err := json.Unmarshal([]byte(jsonAccountsReturn), &jsonResponse); err != nil {
		log.Printf("%s Get AccountID Json Unmarshal Err: %v %v", e.GetName(), err, jsonAccountsReturn)
		return ""
	} else if jsonResponse.Status != "ok" {
		log.Printf("%s Get AccountID Failed: %v", e.GetName(), jsonResponse)
		return ""
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountsReturn); err != nil {
		log.Printf("%s Get AccountID Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return ""
	}

	accountID := strconv.FormatInt(accountsReturn[0].ID, 10)
	return accountID
}

func (e *Huobi) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	if e.Account_ID == "" {
		e.Account_ID = e.GetAccounts()
		if e.Account_ID == "" {
			return
		}
	}

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalances{}
	strRequest := fmt.Sprintf("/v1/account/accounts/%s/balance", e.Account_ID)

	jsonBalanceReturn := e.ApiKeyRequest("GET", make(map[string]string), strRequest)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if jsonResponse.Status != "ok" {
		log.Printf("%s UpdateAllBalances Failed: %v", e.GetName(), jsonResponse)
		return
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return
	}

	for _, v := range accountBalance.List {
		if v.Type == "trade" {
			freeamount, err := strconv.ParseFloat(v.Balance, 64)
			if err == nil {
				c := e.GetCoinBySymbol(v.Currency)
				if c != nil {
					balanceMap.Set(c.Code, freeamount)
				}
			} else {
				log.Printf("%s %s Get Balance Err: %s\n", e.GetName(), v.Currency, err)
			}
		}
	}
}

func (e *Huobi) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {

	return false
}

func (e *Huobi) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	if e.Account_ID == "" {
		e.Account_ID = e.GetAccounts()
		if e.Account_ID == "" {
			return nil, fmt.Errorf("%s Get AccountID Err", e.GetName())
		}
	}

	jsonResponse := &JsonResponse{}
	placeOrder := ""
	strRequest := "/v1/order/orders/place"

	priceFilter := e.GetPriceFilter(pair)
	lotSize := e.GetLotSize(pair)
	price := float64(int(rate/priceFilter+priceFilter/10)) * (priceFilter)
	qty := float64(int(quantity/lotSize+lotSize/10)) * (lotSize)

	log.Printf("%s LimitSell Rate: %s   Quantity: %s", e.GetName(), fmt.Sprintf("%f", price), fmt.Sprintf("%f", qty))

	mapParams := make(map[string]string)
	mapParams["account-id"] = e.Account_ID
	mapParams["amount"] = fmt.Sprintf("%f", qty)
	mapParams["price"] = fmt.Sprintf("%f", price)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["type"] = "sell-limit"

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Status != "ok" {
		return nil, fmt.Errorf("%s LimitSell Failed: %v", e.GetName(), jsonResponse)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      placeOrder,
		Rate:         rate,
		Quantity:     quantity,
		Side:         "Sell",
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Huobi) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	if e.Account_ID == "" {
		e.Account_ID = e.GetAccounts()
		if e.Account_ID == "" {
			return nil, fmt.Errorf("%s Get AccountID Err", e.GetName())
		}
	}

	jsonResponse := &JsonResponse{}
	placeOrder := ""
	strRequest := "/v1/order/orders/place"

	priceFilter := e.GetPriceFilter(pair)
	lotSize := e.GetLotSize(pair)
	price := float64(int(rate/priceFilter+priceFilter/10)) * (priceFilter)
	qty := float64(int(quantity/lotSize+lotSize/10)) * (lotSize)

	mapParams := make(map[string]string)
	mapParams["account-id"] = e.Account_ID
	mapParams["amount"] = fmt.Sprintf("%f", qty)
	mapParams["price"] = fmt.Sprintf("%f", price)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["type"] = "buy-limit"

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Status != "ok" {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v", e.GetName(), jsonResponse)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      placeOrder,
		Rate:         rate,
		Quantity:     quantity,
		Side:         "Buy",
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Huobi) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	if e.Account_ID == "" {
		e.Account_ID = e.GetAccounts()
		if e.Account_ID == "" {
			return fmt.Errorf("%s Get AccountID Err", e.GetName())
		}
	}

	mapParams := make(map[string]string)
	mapParams["uuid"] = order.OrderID

	jsonResponse := &JsonResponse{}
	orderStatus := OrderStatus{}
	strRequest := fmt.Sprintf("/v1/order/orders/%s", order.OrderID)

	jsonOrderStatus := e.ApiKeyRequest("GET", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if jsonResponse.Status != "ok" {
		return fmt.Errorf("%s OrderStatus Failed: %v", e.GetName(), jsonResponse)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order.StatusMessage = jsonOrderStatus
	if orderStatus.State == "canceled" {
		order.Status = exchange.Canceled
	} else if orderStatus.State == "filled" {
		order.Status = exchange.Filled
	} else if orderStatus.State == "partial-filled" || orderStatus.State == "partial-canceled" {
		order.Status = exchange.Partial
	} else if orderStatus.State == "submitting" || orderStatus.State == "submitted" {
		order.Status = exchange.New
	} else {
		order.Status = exchange.Other
	}

	return nil
}

func (e *Huobi) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Huobi) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	if e.Account_ID == "" {
		e.Account_ID = e.GetAccounts()
		if e.Account_ID == "" {
			return fmt.Errorf("%s Get AccountID Err", e.GetName())
		}
	}

	jsonResponse := &JsonResponse{}
	cancelOrder := ""
	strRequest := fmt.Sprintf("/v1/order/orders/%s/submitcancel", order.OrderID)

	jsonCancelOrder := e.ApiKeyRequest("POST", make(map[string]string), strRequest)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if jsonResponse.Status != "ok" {
		return fmt.Errorf("%s CancelOrder Failed: %v", e.GetName(), jsonResponse)
	}
	if err := json.Unmarshal(jsonResponse.Data, &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Huobi) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Huobi) ApiKeyRequest(strMethod string, mapParams map[string]string, strRequestPath string) string {

	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05")
	strUrl := API_URL + strRequestPath

	mapParams["AccessKeyId"] = e.API_KEY
	mapParams["SignatureMethod"] = "HmacSHA256"
	mapParams["SignatureVersion"] = "2"
	mapParams["Timestamp"] = timestamp

	hostName := "api.huobi.pro"
	mapParams["Signature"] = CreateSign(mapParams, strMethod, hostName, strRequestPath, e.API_SECRET)

	var strRequestUrl string
	strParams := MapSortByKey(mapParams)
	strRequestUrl = strUrl + "?" + strParams

	if strMethod == "POST" {
		return exchange.HttpPostRequest(strRequestUrl, mapParams)
	}

	// 构建Request, 并且按官方要求添加Http Header
	request, err := http.NewRequest("GET", strRequestUrl, nil)
	if nil != err {
		return err.Error()
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")

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

func CreateSign(mapParams map[string]string, strMethod, strHostUrl, strRequestPath, strSecretKey string) string {
	sortedParams := MapSortByKey(mapParams) //将数据根据ASCII进行排序
	strPayload := strMethod + "\n" + strHostUrl + "\n" + strRequestPath + "\n" + sortedParams

	return exchange.ComputeHmac256Base64(strPayload, strSecretKey)
}

func MapSortByKey(mapValue map[string]string) string {
	keys := make([]string, 0, len(mapValue))
	for key := range mapValue {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	mapParams := ""
	for _, key := range keys {
		mapParams += (key + "=" + url.QueryEscape(mapValue[key]) + "&")
	}
	mapParams = mapParams[:len(mapParams)-1]
	return mapParams
}
