package bitmax

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
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
	API_URL string = "https://bitmax.io"
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
func (e *Bitmax) GetCoinsData() error {
	coinsData := CoinsData{}

	strRequestUrl := "/api/v1/assets"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %s", e.GetName(), err, jsonCurrencyReturn)
	}

	for _, data := range coinsData {
		c := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			c = coin.GetCoin(data.AssetCode)
			if c == nil {
				c = &coin.Coin{}
				c.Code = data.AssetCode
				c.Name = data.AssetName
				coin.AddCoin(c)
			}
		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(data.AssetCode)
		}

		if c != nil {
			deposit := true
			withdraw := true
			if data.Status == "NotTrading" || data.Status == "NoTransaction" {
				deposit = false
				withdraw = false
			} else if data.Status == "NoWithdraw" {
				withdraw = false
			} else if data.Status == "NoDeposit" {
				deposit = false
			}
			coinConstraint := e.GetCoinConstraint(c)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       c.ID,
					Coin:         c,
					ExSymbol:     data.AssetCode,
					ChainType:    exchange.MAINNET,
					TxFee:        data.WithdrawalFee,
					Withdraw:     withdraw,
					Deposit:      deposit,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = data.AssetCode
				coinConstraint.TxFee = data.WithdrawalFee
				coinConstraint.Withdraw = withdraw
				coinConstraint.Deposit = deposit
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
func (e *Bitmax) GetPairsData() error {
	pairsData := PairsData{}

	strRequestUrl := "/api/v1/products"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %s", e.GetName(), err, jsonSymbolsReturn)
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
		if p != nil && data.Status != "NoTrading" {
			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{
					PairID:      p.ID,
					Pair:        p,
					ExSymbol:    data.Symbol,
					MakerFee:    DEFAULT_MAKER_FEE,
					TakerFee:    DEFAULT_TAKER_FEE,
					LotSize:     math.Pow10(-1 * data.QtyScale),
					PriceFilter: math.Pow10(-1 * data.PriceScale),
					Listed:      data.Status != "NoTrading",
				}
			} else {
				pairConstraint.ExSymbol = data.Symbol
				pairConstraint.LotSize = math.Pow10(-1 * data.QtyScale)
				pairConstraint.PriceFilter = math.Pow10(-1 * data.PriceScale)
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
func (e *Bitmax) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	strRequestUrl := "/api/v1/depth"
	strUrl := API_URL + strRequestUrl

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol
	mapParams["n"] = "100"

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)

	var err error
	for _, bid := range orderBook.Bids {
		buydata := exchange.Order{}
		buydata.Rate, err = strconv.ParseFloat(bid[0], 64)
		if err != nil {
			return nil, err
		}
		buydata.Quantity, err = strconv.ParseFloat(bid[1], 64)
		if err != nil {
			return nil, err
		}
		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Asks {
		selldata := exchange.Order{}
		selldata.Rate, err = strconv.ParseFloat(ask[0], 64)
		if err != nil {
			return nil, err
		}
		selldata.Quantity, err = strconv.ParseFloat(ask[1], 64)
		if err != nil {
			return nil, err
		}
		maker.Asks = append(maker.Asks, selldata)
	}
	maker.LastUpdateID = int64(orderBook.Seqnum)

	return maker, nil
}

func (e *Bitmax) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *Bitmax) DoAccountOperation(operation *exchange.AccountOperation) error {
	return nil
}

func (e *Bitmax) AccountGroup() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	strUrl := "/api/v1/user/info"
	account := AccountGroup{}

	jsonAccountGroup := e.ApiKeyGet(nil, strUrl, "user/info")
	err := json.Unmarshal([]byte(jsonAccountGroup), &account)
	if err != nil {
		log.Printf("%s get Account Group jsonUnmarshal error :%v", e.GetName(), err)
	}

	e.Account_Group = fmt.Sprintf("%v", account.AccountGroup) //set the AccountGroup
}

func (e *Bitmax) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	accountBalance := AccountBalances{}
	jsonResponse := JsonResponse{}
	strRequest := fmt.Sprintf("/%v/api/v1/balance", e.Account_Group)

	jsonBalanceReturn := e.ApiKeyGet(nil, strRequest, "balance")
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if jsonResponse.Code != 0 {
		log.Printf("%s UpdateAllBalances Failed: %v %v", e.GetName(), jsonResponse.Code, jsonResponse.Message)
		return
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return
	}

	for _, freeBalance := range accountBalance {
		freeAmount, err := strconv.ParseFloat(freeBalance.AvailableAmount, 64)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		c := e.GetCoinBySymbol(freeBalance.AssetCode)
		if c != nil {
			balanceMap.Set(c.Code, freeAmount)
		}

	}

}

func (e *Bitmax) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil", e.GetName())
		return false
	}

	// API 1.2 be deprecated
	/* withdraw := Withdrawal{}
	strRequest := fmt.Sprintf("/%v/api/v1/withdraw", e.Account_Group)

	mapParams := make(map[string]string)
	mapParams["requestId"] = fmt.Sprintf("%v%v", time.Now().UTC().UnixNano()/1000000, time.Now().UTC().UnixNano())
	mapParams["assetCode"] = e.GetSymbolByCoin(coin)
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["address"] = addr

	jsonSubmitWithdraw := e.ApiKeyRequest(mapParams, "POST", strRequest, "withdraw")
	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &withdraw); err != nil {
		log.Printf("%s Withdraw Json Unmarshal Err: %v %v", e.GetName(), err, jsonSubmitWithdraw)
		return false
	} else if withdraw.Status != "success" {
		log.Printf("%s Withdraw Failed: %v", e.GetName(), withdraw.Msg)
		return false
	}
	return true */
	return false
}

func (e *Bitmax) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequestUrl := fmt.Sprintf("/%v/api/v1/order", e.Account_Group)

	mapParams := make(map[string]string)
	mapParams["coid"] = fmt.Sprintf("%v%v", time.Now().UTC().UnixNano(), time.Now().UTC().UnixNano()/1000000)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["orderPrice"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["orderQty"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["orderType"] = "limit"
	mapParams["side"] = "sell"

	jsonPlaceReturn := e.ApiKeyRequest(mapParams, "POST", strRequestUrl, "order")
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != 0 {
		return nil, fmt.Errorf("%s LimitSell Failed: %v %+v", e.GetName(), jsonResponse.Code, jsonResponse)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	} else if !placeOrder.Success {
		return nil, fmt.Errorf("%s LimitSell failed :%v", e.GetName(), placeOrder.Success)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      placeOrder.Coid,
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Bitmax) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequestUrl := fmt.Sprintf("/%v/api/v1/order", e.Account_Group)

	mapParams := make(map[string]string)
	mapParams["coid"] = fmt.Sprintf("%v%v", time.Now().UTC().UnixNano(), time.Now().UTC().UnixNano()/1000000)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["orderPrice"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["orderQty"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["orderType"] = "limit"
	mapParams["side"] = "buy"

	jsonPlaceReturn := e.ApiKeyRequest(mapParams, "POST", strRequestUrl, "order")
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != 0 {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v %v", e.GetName(), jsonResponse.Code, jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	} else if !placeOrder.Success {
		return nil, fmt.Errorf("%s LimitBuy failed :%v", e.GetName(), placeOrder.Success)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      placeOrder.Coid,
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Bitmax) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	orderStatus := OrderStatus{}
	strRequestUrl := fmt.Sprintf("/%v/api/v1/order/%v", e.Account_Group, order.OrderID)

	jsonOrderStatus := e.ApiKeyGet(nil, strRequestUrl, "order")
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if jsonResponse.Code != 0 {
		return fmt.Errorf("%s OrderStatus Failed: %v %v", e.GetName(), jsonResponse.Code, jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order.StatusMessage = jsonOrderStatus
	if orderStatus.Status == "New" {
		order.Status = exchange.New
	} else if orderStatus.Status == "PartiallyFilled" {
		order.Status = exchange.Partial
	} else if orderStatus.Status == "Filled" {
		order.Status = exchange.Filled
	} else if orderStatus.Status == "Canceled" {
		order.Status = exchange.Cancelled
	} else {
		order.Status = exchange.Other
	}

	return nil
}

func (e *Bitmax) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Bitmax) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	cancelOrder := PlaceOrder{}
	strRequestUrl := fmt.Sprintf("/%v/api/v1/order", e.Account_Group)

	mapParams := make(map[string]string)
	mapParams["coid"] = fmt.Sprintf("%v%v", time.Now().UTC().UnixNano()/1000000, time.Now().UTC().UnixNano())
	if order != nil {
		mapParams["origCoid"] = order.OrderID
	} else {
		log.Printf("%s CancelOrder, nil order", e.GetName())
		mapParams["origCoid"] = ""
	}
	mapParams["symbol"] = e.GetSymbolByPair(order.Pair)

	jsonCancelOrder := e.ApiKeyRequest(mapParams, "DELETE", strRequestUrl, "order")
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if jsonResponse.Code != 0 {
		return fmt.Errorf("%s CancelOrder Failed: %v %v", e.GetName(), jsonResponse.Code, jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Data, &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	if cancelOrder.Action == "cancel" {
		order.Status = exchange.Canceling
		order.CancelStatus = jsonCancelOrder
	}

	return nil
}

func (e *Bitmax) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Bitmax) ApiKeyGet(mapParams map[string]string, strRequestPath, path string) string {
	timestamp := fmt.Sprintf("%v", time.Now().UTC().UnixNano()/1000000)
	strUrl := API_URL + strRequestPath

	request, err := http.NewRequest("GET", strUrl, nil)
	if nil != err {
		return err.Error()
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("x-auth-key", e.API_KEY)
	request.Header.Add("x-auth-signature", ComputeHmac256EncodingTwice(CreatePayload(timestamp, path, ""), e.API_SECRET))
	request.Header.Add("x-auth-timestamp", timestamp)

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

func (e *Bitmax) ApiKeyRequest(mapParams map[string]string, strMethod, strRequestPath, path string) string {
	timestamp := fmt.Sprintf("%v", time.Now().UTC().UnixNano()/1000000)
	strUrl := API_URL + strRequestPath

	mapParams["time"] = timestamp

	jsonParams := ""
	if nil != mapParams {
		bytesParams, _ := json.Marshal(mapParams)
		jsonParams = string(bytesParams)
	}

	httpClient := &http.Client{}

	request, err := http.NewRequest(strMethod, strUrl, strings.NewReader(jsonParams))
	if nil != err {
		return err.Error()
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("x-auth-key", e.API_KEY)
	request.Header.Add("x-auth-timestamp", timestamp)
	request.Header.Add("x-auth-signature", ComputeHmac256EncodingTwice(CreatePayload(timestamp, path, mapParams["coid"]), e.API_SECRET))
	request.Header.Add("x-auth-coid", mapParams["coid"])

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

func CreatePayload(nonce string, path string, coid string) string {
	if coid != "" {
		return fmt.Sprintf("%v+%v+%v", nonce, path, coid)
	} else {
		return fmt.Sprintf("%v+%v", nonce, path)
	}
}

func ComputeHmac256EncodingTwice(strMessage string, strSecret string) string {
	key, _ := base64.StdEncoding.DecodeString(strSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(strMessage))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
