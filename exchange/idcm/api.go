package idcm

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"

	"strconv"
)

const (
	API_URL string = "https://api.idcs.io:8323/api/v1/RealTimeQuote"
	// "https://api.idcs.io:8323/api/v1/RealTimeQuote"
	// "https://api.IDCM.cc:8323/api/v1"
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
func (e *Idcm) GetCoinsData() error {
	// jsonResponse := &JsonResponse{}
	pairsData := PairsData{}

	strRequestUrl := "/GetRealTimeQuotes" // "/getticker"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &pairsData); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if pairsData.StatusCode != "200" {
		return fmt.Errorf("%s Get Coins Failed: %s", e.GetName(), jsonCurrencyReturn)
	}
	// if err := json.Unmarshal(jsonResponse.Data, &pairsData); err != nil {
	// 	return fmt.Errorf("%s Get Coins Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	// }

	for _, data := range pairsData.Data {
		baseSymbol := strings.Split(data.TradePairCode, "/")[1]
		targetSymbol := strings.Split(data.TradePairCode, "/")[0]
		base := &coin.Coin{}
		target := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base = coin.GetCoin(baseSymbol)
			if base == nil {
				base = &coin.Coin{}
				base.Code = baseSymbol
				coin.AddCoin(base)
			}
			target = coin.GetCoin(targetSymbol)
			if target == nil {
				target = &coin.Coin{}
				target.Code = targetSymbol
				coin.AddCoin(target)
			}
		case exchange.JSON_FILE:
			base = e.GetCoinBySymbol(baseSymbol)
			target = e.GetCoinBySymbol(targetSymbol)
		}

		if base != nil {
			coinConstraint := &exchange.CoinConstraint{
				CoinID:       base.ID,
				Coin:         base,
				ExSymbol:     baseSymbol,
				ChainType:    exchange.MAINNET,
				TxFee:        DEFAULT_TXFEE,
				Withdraw:     DEFAULT_WITHDRAW,
				Deposit:      DEFAULT_DEPOSIT,
				Confirmation: DEFAULT_CONFIRMATION,
				Listed:       true,
			}
			e.SetCoinConstraint(coinConstraint)
		}

		if target != nil {
			coinConstraint := &exchange.CoinConstraint{
				CoinID:       target.ID,
				Coin:         target,
				ExSymbol:     targetSymbol,
				ChainType:    exchange.MAINNET,
				TxFee:        DEFAULT_TXFEE,
				Withdraw:     DEFAULT_WITHDRAW,
				Deposit:      DEFAULT_DEPOSIT,
				Confirmation: DEFAULT_CONFIRMATION,
				Listed:       true,
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
func (e *Idcm) GetPairsData() error {
	// jsonResponse := &JsonResponse{}
	pairsData := PairsData{}

	strRequestUrl := "/GetRealTimeQuotes" // "/getticker"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if pairsData.StatusCode != "200" {
		return fmt.Errorf("%s Get Pairs Failed: %s", e.GetName(), jsonSymbolsReturn)
	}
	// if err := json.Unmarshal(jsonResponse.Data, &pairsData); err != nil {
	// 	return fmt.Errorf("%s Get Pairs Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	// }

	for _, data := range pairsData.Data {
		baseSymbol := strings.Split(data.TradePairCode, "/")[1]
		targetSymbol := strings.Split(data.TradePairCode, "/")[0]
		exSymbol := fmt.Sprintf("%s-%s", targetSymbol, baseSymbol)
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(baseSymbol)
			target := coin.GetCoin(targetSymbol)
			if base != nil && target != nil {
				p = pair.GetPair(base, target)
			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(exSymbol)
		}
		if p != nil && data.Open != 0 {
			pairConstraint := &exchange.PairConstraint{
				PairID:      p.ID,
				Pair:        p,
				ExSymbol:    exSymbol,
				MakerFee:    DEFAULT_MAKER_FEE,
				TakerFee:    DEFAULT_TAKER_FEE,
				LotSize:     math.Pow10(-1 * data.QuantityDigit),
				PriceFilter: math.Pow10(-1 * data.PriceDigit),
				Listed:      true,
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
func (e *Idcm) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	jsonResponse := &JsonResponse{}
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol

	strRequestUrl := "/getdepth"
	strUrl := API_URL + strRequestUrl

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if jsonResponse.Code != "200" {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %v", e.GetName(), jsonOrderbook)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)

	for _, bid := range orderBook.Bids {
		buydata := exchange.Order{}

		buydata.Rate = bid.Price
		buydata.Quantity = bid.Amount

		maker.Bids = append(maker.Bids, buydata)
	}

	for _, ask := range orderBook.Asks {
		selldata := exchange.Order{}

		selldata.Rate = ask.Price
		selldata.Quantity = ask.Amount

		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

/*************** Private API ***************/
func (e *Idcm) DoAccoutOperation(operation *exchange.AccountOperation) error {
	switch operation.Type {

	// case exchange.Transfer:
	// 	return e.transfer(operation)
	// case exchange.BalanceList:
	// 	return e.getAllBalance(operation)
	// case exchange.Balance:
	// 	return e.getBalance(operation)

	case exchange.Withdraw:
		return e.doWithdraw(operation)

	}
	return fmt.Errorf("Operation type invalid: %v", operation.Type)
}

func (e *Idcm) doWithdraw(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	mapParams := make(map[string]string)
	mapParams["currency"] = e.GetSymbolByCoin(operation.Coin)
	mapParams["quantity"] = operation.WithdrawAmount
	mapParams["address"] = operation.WithdrawAddress

	jsonResponse := &JsonResponse{}
	uuid := Uuid{}
	strRequest := "/v1.1/account/withdraw"

	jsonSubmitWithdraw := e.ApiKeyGET(strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.MapParams = fmt.Sprintf("%+v", mapParams)
		operation.CallResponce = jsonSubmitWithdraw
	}

	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Json Unmarshal Err: %v, %s", e.GetName(), err, jsonSubmitWithdraw)
		return operation.Error
	} else if jsonResponse.Code != "200" {
		operation.Error = fmt.Errorf("%s Withdraw Failed: %v", e.GetName(), jsonSubmitWithdraw)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &uuid); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	operation.WithdrawID = uuid.Id

	return nil
}

func (e *Idcm) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalances{}
	strRequest := "/v1.1/account/getbalances"

	jsonBalanceReturn := e.ApiKeyGET(strRequest, make(map[string]string))
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if jsonResponse.Code != "200" {
		log.Printf("%s UpdateAllBalances Failed: %v", e.GetName(), jsonBalanceReturn)
		return
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return
	}

	for _, v := range accountBalance {
		c := e.GetCoinBySymbol(v.Currency)
		if c != nil {
			balanceMap.Set(c.Code, v.Available)
		}
	}
}

func (e *Idcm) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
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
	} else if jsonResponse.Code != "200" {
		log.Printf("%s Withdraw Failed: %v", e.GetName(), jsonSubmitWithdraw)
		return false
	}
	if err := json.Unmarshal(jsonResponse.Data, &uuid); err != nil {
		log.Printf("%s Withdraw Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return false
	}
	return true
}

func (e *Idcm) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	mapParams := make(map[string]string)
	mapParams["market"] = e.GetSymbolByPair(pair)
	mapParams["quantity"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["rate"] = strconv.FormatFloat(rate, 'f', -1, 64)

	jsonResponse := &JsonResponse{}
	uuid := Uuid{}
	strRequest := "/v1.1/market/selllimit"

	jsonPlaceReturn := e.ApiKeyGET(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != "200" {
		return nil, fmt.Errorf("%s LimitSell Failed: %v", e.GetName(), jsonPlaceReturn)
	}
	if err := json.Unmarshal(jsonResponse.Data, &uuid); err != nil {
		return nil, fmt.Errorf("%s LimitSell Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      uuid.Id,
		Rate:         rate,
		Quantity:     quantity,
		Side:         "Sell",
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Idcm) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
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
	} else if jsonResponse.Code != "200" {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v", e.GetName(), jsonPlaceReturn)
	}
	if err := json.Unmarshal(jsonResponse.Data, &uuid); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      uuid.Id,
		Rate:         rate,
		Quantity:     quantity,
		Side:         "Buy",
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}
	return order, nil
}

func (e *Idcm) OrderStatus(order *exchange.Order) error {
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
	} else if jsonResponse.Code != "200" {
		return fmt.Errorf("%s OrderStatus Failed: %v", e.GetName(), jsonOrderStatus)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
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

func (e *Idcm) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Idcm) CancelOrder(order *exchange.Order) error {
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
	} else if jsonResponse.Code != "200" {
		return fmt.Errorf("%s CancelOrder Failed: %v", e.GetName(), jsonCancelOrder)
	}
	if err := json.Unmarshal(jsonResponse.Data, &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Idcm) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Idcm) ApiKeyGET(strRequestPath string, mapParams map[string]string) string {
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

func (e *Idcm) ApiKeyRequest(strMethod string, strRequestPath string, mapParams map[string]string) string {
	// mapParams["timestamp"] = fmt.Sprintf("%d", time.Now().UnixNano())
	// log.Printf("timestamp: %v", mapParams["timestamp"]) //==============
	// fmt.Sprintf("%d", time.Now().UnixNano())
	// fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))

	// mapParams["signature"] = signature

	strUrl := API_URL + strRequestPath /* + "?" + exchange.Map2UrlQuery(mapParams) */

	jsonParams := ""
	if nil != mapParams {
		bytesParams, _ := json.Marshal(mapParams)
		jsonParams = string(bytesParams)
	}

	// jsonParams = "1" // *************************************************

	// signature := ComputeHmac384NoDecode(exchange.Map2UrlQuery(mapParams), e.API_SECRET)
	signature := ComputeHmac384NoDecode(jsonParams, e.API_SECRET)

	httpClient := &http.Client{}

	// log.Printf("url: %v, body: %v", strUrl, jsonParams) // ========================

	request, err := http.NewRequest(strMethod, strUrl, strings.NewReader(jsonParams))
	if nil != err {
		return err.Error()
	}

	request.Header.Add("X-IDCM-APIKEY", e.API_KEY)
	request.Header.Add("X-IDCM-SIGNATURE", signature)
	request.Header.Add("X-IDCM-INPUT", jsonParams) // input: mapParams to json
	request.Header.Add("Accept", "text/html, application/xhtml+xml, */*")
	request.Header.Add("Content-Type", "application/json;charset=utf-8")
	// method -> post

	// request.Header.Add("Content-Type", "application/json;charset=utf-8")
	// request.Header.Add("Accept", "application/json")
	// request.Header.Add("apisign", signature)

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

func ComputeHmac384NoDecode(strMessage string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha512.New384, key)
	h.Write([]byte(strMessage))

	// return hex.EncodeToString(h.Sum(nil))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
