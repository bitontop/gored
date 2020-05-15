package idcm

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"crypto/hmac"
	"crypto/sha256"
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
	API_URL_PUB string = "https://api.idcs.io:8323/api/v1/RealTimeQuote"
	API_URL     string = "https://api.IDCM.cc:8323/api/v1"
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
	strUrl := API_URL_PUB + strRequestUrl

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
			coinConstraint := e.GetCoinConstraint(base)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
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
			} else {
				coinConstraint.ExSymbol = baseSymbol
			}
			e.SetCoinConstraint(coinConstraint)
		}

		if target != nil {
			coinConstraint := e.GetCoinConstraint(target)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
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
			} else {
				coinConstraint.ExSymbol = targetSymbol
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
	strUrl := API_URL_PUB + strRequestUrl

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
			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{
					PairID:      p.ID,
					Pair:        p,
					ExSymbol:    exSymbol,
					MakerFee:    DEFAULT_MAKER_FEE,
					TakerFee:    DEFAULT_TAKER_FEE,
					LotSize:     math.Pow10(-1 * data.QuantityDigit),
					PriceFilter: math.Pow10(-1 * data.PriceDigit),
					Listed:      true,
				}
			} else {
				pairConstraint.ExSymbol = exSymbol
				pairConstraint.LotSize = math.Pow10(-1 * data.QuantityDigit)
				pairConstraint.PriceFilter = math.Pow10(-1 * data.PriceDigit)
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
	strUrl := API_URL_PUB + strRequestUrl

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

/*************** Public API ***************/
func (e *Idcm) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *Idcm) DoAccountOperation(operation *exchange.AccountOperation) error {
	switch operation.Type {

	// case exchange.Transfer:
	// 	return e.transfer(operation)
	case exchange.BalanceList:
		return e.getAllBalance(operation)
	case exchange.Balance:
		return e.getBalance(operation)

	case exchange.Withdraw:
		return e.doWithdraw(operation)

	}
	return fmt.Errorf("%s Operation type invalid: %s %v", operation.Ex, operation.Wallet, operation.Type)
}

func (e *Idcm) getAllBalance(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	accountBalances := AccountBalances{}
	strRequest := "/getuserinfo"

	mapParams := make(map[string]string)

	jsonAllBalanceReturn := e.ApiKeyRequest("POST", strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonAllBalanceReturn
	}

	// log.Printf("jsonAllBalanceReturn: %v", jsonAllBalanceReturn)
	if err := json.Unmarshal([]byte(jsonAllBalanceReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s getAllBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonAllBalanceReturn)
		return operation.Error
	} else if jsonResponse.Code != "200" {
		operation.Error = fmt.Errorf("%s getAllBalance Failed: %v", e.GetName(), jsonAllBalanceReturn)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountBalances); err != nil {
		operation.Error = fmt.Errorf("%s getAllBalance Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	for _, account := range accountBalances {
		if account.Free+account.Freezed == 0 {
			continue
		}

		balance := exchange.AssetBalance{
			Coin:             e.GetCoinBySymbol(account.Code),
			BalanceAvailable: account.Free,
			BalanceFrozen:    account.Freezed,
		}
		operation.BalanceList = append(operation.BalanceList, balance)

	}

	return nil
}

func (e *Idcm) getBalance(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	accountBalances := AccountBalances{}
	strRequest := "/getuserinfo"
	symbol := e.GetSymbolByCoin(operation.Coin)

	mapParams := make(map[string]string)

	jsonBalanceReturn := e.ApiKeyRequest("POST", strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonBalanceReturn
	}

	// log.Printf("jsonBalanceReturn: %v", jsonBalanceReturn)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s getBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonBalanceReturn)
		return operation.Error
	} else if jsonResponse.Code != "200" {
		operation.Error = fmt.Errorf("%s getBalance Failed: %v", e.GetName(), jsonBalanceReturn)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountBalances); err != nil {
		operation.Error = fmt.Errorf("%s getBalance Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	for _, account := range accountBalances {
		if account.Code == symbol {
			operation.BalanceFrozen = account.Freezed
			operation.BalanceAvailable = account.Free
			return nil
		}
	}

	return fmt.Errorf("%s getBalance fail: %v", e.GetName(), jsonBalanceReturn)
}

func (e *Idcm) doWithdraw(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	if operation.WithdrawTag != "" {
		operation.Error = fmt.Errorf("%s Withdraw Failed, got tag: %v, for coin: %v", e.GetName(), operation.WithdrawTag, operation.Coin.Code)
		return operation.Error
	}

	mapParams := make(map[string]string)
	mapParams["Symbol"] = e.GetSymbolByCoin(operation.Coin)
	mapParams["Amount"] = operation.WithdrawAmount
	mapParams["Address"] = operation.WithdrawAddress

	jsonResponse := &JsonResponse{}
	withdraw := Withdraw{}
	strRequest := "/withdraw"

	jsonSubmitWithdraw := e.ApiKeyRequest("POST", strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonSubmitWithdraw
	}

	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Json Unmarshal Err: %v, %s", e.GetName(), err, jsonSubmitWithdraw)
		return operation.Error
	} else if jsonResponse.Code != "200" {
		operation.Error = fmt.Errorf("%s Withdraw Failed: %v", e.GetName(), jsonSubmitWithdraw)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &withdraw); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	operation.WithdrawID = withdraw.WithdrawID

	return nil
}

func (e *Idcm) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalances{}
	strRequest := "/getuserinfo"

	jsonBalanceReturn := e.ApiKeyRequest("POST", strRequest, make(map[string]string))
	// log.Printf("================jsonBalanceReturn: %v", jsonBalanceReturn) // =========================
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
		c := e.GetCoinBySymbol(v.Code)
		if c != nil {
			balanceMap.Set(c.Code, v.Free)
		}
	}
}

func (e *Idcm) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil", e.GetName())
		return false
	}

	// mapParams := make(map[string]string)
	// mapParams["currency"] = e.GetSymbolByCoin(coin)
	// mapParams["quantity"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	// mapParams["address"] = addr

	// jsonResponse := &JsonResponse{}
	// uuid := Uuid{}
	// strRequest := "/v1.1/account/withdraw"

	// jsonSubmitWithdraw := e.ApiKeyGET(strRequest, mapParams)
	// if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &jsonResponse); err != nil {
	// 	log.Printf("%s Withdraw Json Unmarshal Err: %v %v", e.GetName(), err, jsonSubmitWithdraw)
	// 	return false
	// } else if jsonResponse.Code != "200" {
	// 	log.Printf("%s Withdraw Failed: %v", e.GetName(), jsonSubmitWithdraw)
	// 	return false
	// }
	// if err := json.Unmarshal(jsonResponse.Data, &uuid); err != nil {
	// 	log.Printf("%s Withdraw Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	// 	return false
	// }
	// return true
	return false
}

func (e *Idcm) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/trade"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["Symbol"] = e.GetSymbolByPair(pair)
	mapParams["Size"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)
	mapParams["Price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["Type"] = "1" // limit
	mapParams["Side"] = "1" // 0 buy, 1 sell

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != "200" {
		return nil, fmt.Errorf("%s LimitSell Failed: %v", e.GetName(), jsonPlaceReturn)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      placeOrder.Orderid,
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Idcm) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/trade"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["Symbol"] = e.GetSymbolByPair(pair)
	mapParams["Size"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)
	mapParams["Price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["Type"] = "1" // limit
	mapParams["Side"] = "0" // 0 buy, 1 sell

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != "200" {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v", e.GetName(), jsonPlaceReturn)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      placeOrder.Orderid,
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Idcm) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	orderStatus := OrderStatus{}
	strRequest := "/getorderinfo"

	mapParams := make(map[string]string)
	mapParams["Symbol"] = e.GetSymbolByPair(order.Pair)
	mapParams["OrderID"] = order.OrderID

	jsonOrderStatus := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if jsonResponse.Code != "200" {
		return fmt.Errorf("%s OrderStatus Failed: %v", e.GetName(), jsonOrderStatus)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	if len(orderStatus) == 0 {
		return fmt.Errorf("%s OrderStatus Order not found, : %v", e.GetName(), jsonOrderStatus)
	}

	status := orderStatus[0]
	order.StatusMessage = jsonOrderStatus
	if status.Status == 1 {
		order.Status = exchange.Partial
	} else if status.Status == 2 {
		order.Status = exchange.Filled
	} else if status.Status == 0 {
		order.Status = exchange.New
	} else if status.Status == -2 {
		order.Status = exchange.Cancelled
	} else {
		order.Status = exchange.Other
	}

	return nil
}

func (e *Idcm) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

// will return true even if the order doesn't exist
func (e *Idcm) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	mapParams := make(map[string]string)
	mapParams["uuid"] = order.OrderID

	jsonResponse := &JsonResponse{}
	var result bool
	strRequest := "/cancel_order"

	jsonCancelOrder := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if jsonResponse.Code != "200" {
		return fmt.Errorf("%s CancelOrder Failed: %v", e.GetName(), jsonCancelOrder)
	}
	if err := json.Unmarshal(jsonResponse.Data, &result); err != nil {
		return fmt.Errorf("%s CancelOrder Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	} else if !result {
		return fmt.Errorf("%s CancelOrder Failed: %v", e.GetName(), jsonCancelOrder)
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

	strUrl := API_URL + strRequestPath /* + "?" + exchange.Map2UrlQuery(mapParams) */

	jsonParams := ""
	if nil != mapParams {
		bytesParams, _ := json.Marshal(mapParams)
		jsonParams = string(bytesParams)
	}

	signature := ComputeHmac384NoDecode(jsonParams, e.API_SECRET)

	httpClient := &http.Client{}

	// log.Printf("url: %v, body: %v, signature: %v", strUrl, jsonParams, signature) // ========================

	request, err := http.NewRequest(strMethod, strUrl, strings.NewReader(jsonParams))
	if nil != err {
		return err.Error()
	}

	request.Header.Add("X-IDCM-APIKEY", e.API_KEY)
	request.Header.Add("X-IDCM-SIGNATURE", signature)
	request.Header.Add("X-IDCM-INPUT", jsonParams) // input: mapParams to json
	// request.Header.Add("Accept", "text/html, application/xhtml+xml, */*")
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

func ComputeHmac256Base64(strMessage string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(strMessage))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
