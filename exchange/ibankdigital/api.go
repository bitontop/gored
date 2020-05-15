package ibankdigital

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
	API_URL string = "https://www.ibankex.io/api"
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
func (e *Ibankdigital) GetCoinsData() error {
	jsonResponse := &JsonResponse{}
	var coinsData []string

	strRequestUrl := "/v1/common/currencys"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if jsonResponse.Status != "ok" {
		return fmt.Errorf("%s Get Coins Failed: %v", e.GetName(), jsonResponse)
	}
	if err := json.Unmarshal(jsonResponse.Data, &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, data := range coinsData {
		c := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			c = coin.GetCoin(data)
			if c == nil {
				c = &coin.Coin{}
				c.Code = data
				coin.AddCoin(c)
			}
		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(data)
		}

		if c != nil {
			coinConstraint := e.GetCoinConstraint(c)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       c.ID,
					Coin:         c,
					ExSymbol:     data,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = data
			}
			e.SetCoinConstraint(coinConstraint)
		}
	}
	e.updateCoinInfo()

	return nil
}

/* GetPairsData - Get Pairs Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func (e *Ibankdigital) GetPairsData() error {
	jsonResponse := &JsonResponse{}
	pairsData := PairsData{}

	strRequestUrl := "/v1/common/symbols"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if jsonResponse.Status != "ok" {
		return fmt.Errorf("%s Get Pairs Failed: %v", e.GetName(), jsonResponse)
	}
	if err := json.Unmarshal(jsonResponse.Data, &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
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
			p = e.GetPairBySymbol(data.Symbol)
		}
		if p != nil && data.State == "online" {
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
					Listed:      data.State == "online",
				}
			} else {
				pairConstraint.ExSymbol = data.Symbol
				pairConstraint.LotSize = math.Pow10(-1 * data.AmountPrecision)
				pairConstraint.PriceFilter = math.Pow10(-1 * data.PricePrecision)
				pairConstraint.Listed = data.State == "online"
			}
			e.SetPairConstraint(pairConstraint)
		}
	}
	return nil
}

func (e *Ibankdigital) updateCoinInfo() error {
	jsonResponse := &HuobiJsonResponse{}
	coinsData := HuobiCoins{}

	strUrl := "https://www.huobi.com/-/x/pro/v1/settings/currencys?r=sqyeinryv8&language=en-US"

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if jsonResponse.Status != "ok" {
		return fmt.Errorf("%s Get Coins Failed: %v", e.GetName(), jsonResponse)
	}
	if err := json.Unmarshal(jsonResponse.Data, &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, data := range coinsData {
		c := &coin.Coin{}
		c = coin.GetCoin(strings.ToLower(data.DisplayName))

		if c != nil {
			coinConstraint := &exchange.CoinConstraint{
				CoinID:       c.ID,
				Coin:         c,
				ExSymbol:     data.Name,
				ChainType:    exchange.MAINNET,
				TxFee:        DEFAULT_TXFEE,
				Withdraw:     data.WithdrawEnabled,
				Deposit:      data.DepositEnabled,
				Confirmation: data.SafeConfirms,
				Listed:       data.State == "online",
			}
			e.SetCoinConstraint(coinConstraint)
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
func (e *Ibankdigital) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	strRequestUrl := "/market/depth"
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
	if err := json.Unmarshal([]byte(jsonOrderbook), &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if orderBook.Status != "ok" {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %v", e.GetName(), jsonOrderbook)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	for _, bid := range orderBook.Tick.Bids {
		var buydata exchange.Order

		//Modify according to type and structure
		buydata.Rate = bid[0]
		buydata.Quantity = bid[1]
		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Tick.Asks {
		var selldata exchange.Order

		//Modify according to type and structure
		selldata.Rate = ask[0]
		selldata.Quantity = ask[1]
		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

func (e *Ibankdigital) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *Ibankdigital) DoAccountOperation(operation *exchange.AccountOperation) error {
	switch operation.Type {

	// case exchange.Transfer:
	// 	return e.transfer(operation)
	// case exchange.BalanceList:
	// 	return e.getAllBalance(operation)
	// case exchange.Balance:
	// 	return e.getBalance(operation)

	// case exchange.Withdraw:
	// 	return e.doWithdraw(operation)

	}
	return fmt.Errorf("%s Operation type invalid: %s %v", operation.Ex, operation.Wallet, operation.Type)
}

// not a valid api, disappeared from api document.
func (e *Ibankdigital) doWithdraw(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := JsonResponse{}
	strRequest := "v1/dw/withdraw/api/create"

	mapParams := make(map[string]string)
	mapParams["address"] = operation.WithdrawAddress
	mapParams["amount"] = operation.WithdrawAmount
	mapParams["currency"] = e.GetSymbolByCoin(operation.Coin)

	jsonSubmitWithdraw := e.ApiKeyPost(strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonSubmitWithdraw
	}

	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Json Unmarshal Err: %v", e.GetName(), err)
		return operation.Error
	} else if jsonResponse.Status != "ok" {
		operation.Error = fmt.Errorf("%s Withdraw Failed: %v", e.GetName(), jsonSubmitWithdraw)
		return operation.Error
	}

	// api document not found
	// operation.WithdrawID =

	return nil
}

func (e *Ibankdigital) GetAccounts() { //doesn't work well, always got err-msg of signature not valid
	jsonResponse := JsonResponse{}
	accountId := AccountID{}

	strRequest := "/v1/account/accounts"

	jsonAccountsReturn := e.ApiKeyGet(strRequest, make(map[string]string))
	if err := json.Unmarshal([]byte(jsonAccountsReturn), &jsonResponse); err != nil {
		log.Printf("%s GetAccounts Json Unmarshal Err: %v %v", e.GetName(), err, jsonAccountsReturn)
	} else if jsonResponse.Status != "ok" {
		log.Printf("%s GetAccounts Failed: %v", e.GetName(), jsonResponse)
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountId); err != nil {
		log.Printf("%s GetAccounts Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	if len(accountId) >= 1 {
		e.Account_ID = fmt.Sprintf("%v", accountId[0].ID)
	} else {
		e.Account_ID = "error"
	}

}

func (e *Ibankdigital) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalances{}
	strRequest := fmt.Sprintf("/v1/account/accounts/%s/balance", e.Account_ID)

	mapParams := make(map[string]string)
	mapParams["account-id"] = e.Account_ID

	jsonBalanceReturn := e.ApiKeyGet(strRequest, mapParams)
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

	for _, balance := range accountBalance.List {
		if balance.Type == "trade" {
			freeAmount, err := strconv.ParseFloat(balance.Balance, 64)
			if err != nil {
				log.Printf("%s balance parse error: %v, %v", e.GetName(), err, balance.Balance)
				return
			}
			c := e.GetCoinBySymbol(balance.Currency)
			if c != nil {
				balanceMap.Set(c.Code, freeAmount)
			}
		}
	}
}

func (e *Ibankdigital) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil", e.GetName())
		return false
	}

	jsonResponse := JsonResponse{}
	strRequest := "v1/dw/withdraw/api/create"

	mapParams := make(map[string]string)
	mapParams["address"] = addr
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["currency"] = e.GetSymbolByCoin(coin)

	jsonSubmitWithdraw := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &jsonResponse); err != nil {
		log.Printf("%s Withdraw Json Unmarshal Err: %v %v", e.GetName(), err, jsonSubmitWithdraw)
		return false
	} else if jsonResponse.Status != "ok" {
		log.Printf("%s Withdraw Failed: %v", e.GetName(), jsonResponse)
		return false
	}

	return true
}

func (e *Ibankdigital) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := ""
	strRequest := "/v1/order/orders/place"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["account-id"] = e.Account_ID
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["type"] = "sell-limit"

	jsonPlaceReturn := e.ApiKeyPost(strRequest, mapParams)
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
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Ibankdigital) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := ""
	strRequest := "/v1/order/orders/place"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["account-id"] = e.Account_ID
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["type"] = "buy-limit"

	jsonPlaceReturn := e.ApiKeyPost(strRequest, mapParams)
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
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Ibankdigital) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	orderStatus := OrderStatus{}
	strRequest := fmt.Sprintf("/v1/order/orders/%s", order.OrderID)

	mapParams := make(map[string]string)
	mapParams["order-id"] = order.OrderID

	jsonOrderStatus := e.ApiKeyGet(strRequest, mapParams)
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
		order.Status = exchange.Cancelled
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

func (e *Ibankdigital) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Ibankdigital) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	cancelOrder := ""
	strRequest := fmt.Sprintf("/v1/order/orders/%s/submitcancel", order.OrderID)

	mapParams := make(map[string]string)
	mapParams["order-id"] = order.OrderID

	jsonCancelOrder := e.ApiKeyPost(strRequest, mapParams)
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

func (e *Ibankdigital) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Ibankdigital) ApiKeyGet(strRequestPath string, mapParams map[string]string) string {
	strMethod := "GET"
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05")

	mapParams["AccessKeyId"] = e.API_KEY
	mapParams["SignatureMethod"] = "HmacSHA256"
	mapParams["SignatureVersion"] = "2"
	mapParams["Timestamp"] = timestamp

	hostName := "www.ibankex.io"
	mapParams["Signature"] = CreateSign(mapParams, strMethod, hostName, strRequestPath, e.API_SECRET)
	strUrl := API_URL + strRequestPath
	httpClient := &http.Client{}

	var strRequestUrl string
	if nil == mapParams {
		strRequestUrl = strUrl
	} else {
		strParams := exchange.Map2UrlQueryUrl(mapParams)
		strRequestUrl = strUrl + "?" + strParams
	}

	// 构建Request, 并且按官方要求添加Http Header
	request, err := http.NewRequest("GET", strRequestUrl, nil)
	if nil != err {
		return err.Error()
	}
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

func (e *Ibankdigital) ApiKeyPost(strRequestPath string, mapParams map[string]string) string {
	strMethod := "POST"
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05")

	mapParams2Sign := make(map[string]string)
	mapParams2Sign["AccessKeyId"] = e.API_KEY
	mapParams2Sign["SignatureMethod"] = "HmacSHA256"
	mapParams2Sign["SignatureVersion"] = "2"
	mapParams2Sign["Timestamp"] = timestamp

	hostName := "www.ibankex.io"
	mapParams2Sign["Signature"] = CreateSign(mapParams2Sign, strMethod, hostName, strRequestPath, e.API_SECRET)
	strUrl := API_URL + strRequestPath + "?" + exchange.Map2UrlQueryUrl(mapParams2Sign)

	return exchange.HttpPostRequest(strUrl, mapParams)
}

func CreateSign(mapParams map[string]string, strMethod, strHostUrl, strRequestPath, strSecretKey string) string {
	sortedParams := exchange.Map2UrlQueryUrl(mapParams) //将数据根据ASCII进行排序
	strPayload := strMethod + "\n" + strHostUrl + "\n" + strRequestPath + "\n" + sortedParams

	return exchange.ComputeHmac256Base64(strPayload, strSecretKey)
}
