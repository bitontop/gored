package ftx

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"

	"strconv"
)

const (
	API_URL string = "https://ftx.com"
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

	strRequestUrl := "/api/coins"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if !jsonResponse.Success {
		return fmt.Errorf("%s Get Coins Failed: %v", e.GetName(), jsonCurrencyReturn)
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
func (e *Ftx) GetPairsData() error {
	jsonResponse := &JsonResponse{}
	pairsData := PairsData{}

	strRequestUrl := "/api/markets"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if !jsonResponse.Success {
		return fmt.Errorf("%s Get Pairs Failed: %v", e.GetName(), jsonSymbolsReturn)
	}
	if err := json.Unmarshal(jsonResponse.Result, &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	for _, data := range pairsData {
		if data.Type != "spot" {
			continue
		}
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(data.QuoteCurrency)
			target := coin.GetCoin(data.BaseCurrency)
			if base != nil && target != nil {
				p = pair.GetPair(base, target)
			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(data.Name)
		}
		if p != nil {
			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{ // no minTradeBaseQuantity
					PairID:           p.ID,
					Pair:             p,
					ExSymbol:         data.Name,
					MakerFee:         DEFAULT_MAKER_FEE,
					TakerFee:         DEFAULT_TAKER_FEE,
					LotSize:          data.SizeIncrement,
					PriceFilter:      data.PriceIncrement,
					MinTradeQuantity: data.MinProvideSize,
					Listed:           true,
				}
			} else {
				pairConstraint.ExSymbol = data.Name
				pairConstraint.LotSize = data.SizeIncrement
				pairConstraint.PriceFilter = data.PriceIncrement
				pairConstraint.MinTradeQuantity = data.MinProvideSize
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

// orderbook
func (e *Ftx) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	jsonResponse := &JsonResponse{}
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	if symbol == "" {
		symbol = pair.Symbol
	}

	// log.Printf("FTX symbol:%s", symbol)

	strRequestUrl := fmt.Sprintf("/api/markets/%v/orderbook?depth=%v", symbol, 100)
	strUrl := API_URL + strRequestUrl

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if !jsonResponse.Success {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %v", e.GetName(), jsonOrderbook)
	}
	if err := json.Unmarshal(jsonResponse.Result, &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
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
	return maker, nil
}

/*************** Private API ***************/

func (e *Ftx) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalances{}
	strRequest := "/api/wallet/balances" // "/account"

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
			// str := fmt.Sprintf("%s=%.8f", c.Code, v.Free)
			// log.Printf(str)
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

	symbol := e.GetSymbolByPair(pair)
	if symbol == "" {
		symbol = pair.Symbol
	}

	mapParams := make(map[string]string)
	mapParams["market"] = symbol // future "BTC-PERP", spot "ALTHEDGE/USD"
	mapParams["size"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["side"] = "sell"
	mapParams["type"] = "limit"
	// mapParams["reduceOnly"] = false

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/api/orders"

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if !jsonResponse.Success {
		return nil, fmt.Errorf("%s LimitSell Failed: %s", e.GetName(), jsonPlaceReturn)
	}
	if err := json.Unmarshal(jsonResponse.Result, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%v", placeOrder.ID),
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

	symbol := e.GetSymbolByPair(pair)
	if symbol == "" {
		symbol = pair.Symbol
	}

	// log.Printf("FTX LimitBuy symbol:%s", symbol)

	mapParams := make(map[string]string)
	mapParams["market"] = symbol // future "BTC-PERP", spot "ALTHEDGE/USD"
	mapParams["size"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["side"] = "buy"
	mapParams["type"] = "limit"
	// mapParams["reduceOnly"] = false

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/api/orders"

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if !jsonResponse.Success {
		return nil, fmt.Errorf("%s LimitBuy Failed: %s", e.GetName(), jsonPlaceReturn)
	}
	if err := json.Unmarshal(jsonResponse.Result, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%v", placeOrder.ID),
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

	jsonResponse := &JsonResponse{}
	orderStatus := OrderStatus{}
	strRequest := fmt.Sprintf("/api/orders/%v", order.OrderID)

	mapParams := make(map[string]string)

	jsonOrderStatus := e.ApiKeyRequest("GET", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if !jsonResponse.Success {
		return fmt.Errorf("%s OrderStatus Failed: %v", e.GetName(), jsonOrderStatus)
	}
	if err := json.Unmarshal(jsonResponse.Result, &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	order.Rate = orderStatus.Price
	order.Quantity = orderStatus.Size
	order.DealRate = orderStatus.AvgFillPrice
	order.DealQuantity = orderStatus.FilledSize

	order.StatusMessage = jsonOrderStatus
	if (orderStatus.Status == "new" || orderStatus.Status == "open") && order.DealQuantity == 0 {
		order.Status = exchange.New
	} else if (orderStatus.Status == "new" || orderStatus.Status == "open") && order.DealQuantity < order.Quantity {
		order.Status = exchange.Partial
	} else if order.DealQuantity == order.Quantity {
		order.Status = exchange.Filled
	} else if orderStatus.Status == "closed" || orderStatus.Status == "cancelled" {
		order.Status = exchange.Cancelled
		order.Canceled = true
	} else {
		order.Status = exchange.Other
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

	jsonResponse := &JsonResponse{}
	cancelOrderMsg := ""
	strRequest := fmt.Sprintf("/api/orders/%v", order.OrderID)

	mapParams := make(map[string]string)

	jsonCancelOrder := e.ApiKeyRequest("DELETE", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if !jsonResponse.Success {
		return fmt.Errorf("%s CancelOrder Failed: %v", e.GetName(), jsonCancelOrder)
	}
	if err := json.Unmarshal(jsonResponse.Result, &cancelOrderMsg); err != nil {
		return fmt.Errorf("%s CancelOrder Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	log.Printf("%s CancelOrder: %v", e.GetName(), cancelOrderMsg)

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
	if strMethod == "POST" {
		preSign += postBody
	}
	signature := exchange.ComputeHmac256NoDecode(preSign, e.API_SECRET)

	// log.Printf("postBody: %v", postBody)
	// log.Printf("preSign: %v", preSign)
	// log.Printf("strRequestUrl: %v", strRequestUrl)

	// request
	// strRequestUrl = strRequestPath // ==========================
	// request, err = http.NewRequest(strMethod, strRequestUrl, strings.NewReader(postBody))
	request, err = http.NewRequest(strMethod, strRequestUrl, bytes.NewBuffer([]byte(postBody)))
	if nil != err {
		return err.Error()
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Add("FTX-KEY", e.API_KEY)
	request.Header.Add("FTX-TS", timestamp)
	request.Header.Add("FTX-SIGN", signature)
	if e.Sub_Account_Name != "" {
		request.Header.Add("FTX-SUBACCOUNT", e.Sub_Account_Name)
	}
	// add FTX-SUBACCOUNT if using subaccount

	response, err := httpClient.Do(request)
	if nil != err {
		return err.Error()
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return err.Error()
	}

	// log.Printf("JSON: %v", string(body))
	return string(body)
}
