package okexdm

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
	"strconv"
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
)

/*The Base Endpoint URL*/
const (
	API_URL = "https://www.okex.com"
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
Step 3: Modify API Path(strRequestPath)*/
func (e *Okexdm) GetCoinsData() error {
	contractsData := ContractsData{}

	strRequestPath := "/api/futures/v3/instruments/ticker"
	strUrl := API_URL + strRequestPath

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &contractsData); err != nil {
		return fmt.Errorf("%s Get Coins Result Unmarshal Err: %v %s", e.GetName(), err, jsonCurrencyReturn)
	}

	for _, data := range contractsData {
		c := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			c = coin.GetCoin(GetCodeByDate(data.InstrumentID))
			if c == nil {
				c = &coin.Coin{
					Code: GetCodeByDate(data.InstrumentID),
				}
				coin.AddCoin(c)
			}
		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(GetCodeByDate(data.InstrumentID))
		}

		if c != nil {
			coinConstraint := e.GetCoinConstraint(c)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       c.ID,
					Coin:         c,
					ExSymbol:     data.InstrumentID,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = data.InstrumentID
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
func (e *Okexdm) GetPairsData() error {
	contractsData := ContractsData{}

	strRequestPath := "/api/futures/v3/instruments/ticker"
	strUrl := API_URL + strRequestPath

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &contractsData); err != nil {
		return fmt.Errorf("%s Get Pairs Result Unmarshal Err: %v %s", e.GetName(), err, jsonSymbolsReturn)
	}

	for _, data := range contractsData {
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin("USD") // &coin.Coin{000, "USD", "USD", "", "", "", 0, 0}
			target := coin.GetCoin(GetCodeByDate(data.InstrumentID))
			if base != nil && target != nil {
				p = pair.GetPair(base, target)
			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(GetCodeByDate(data.InstrumentID))
		}
		if p != nil {
			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{
					PairID:      p.ID,
					Pair:        p,
					ExSymbol:    data.InstrumentID,
					MakerFee:    DEFAULT_MAKER_FEE,
					TakerFee:    DEFAULT_TAKER_FEE,
					LotSize:     DEFAULT_LOT_SIZE,
					PriceFilter: DEFAULT_PRICE_FILTER,
					Listed:      true,
				}
			} else {
				pairConstraint.ExSymbol = data.InstrumentID
			}
			e.SetPairConstraint(pairConstraint)
		}
	}
	return nil
}

/*Get Pair Market Depth
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Get Exchange Pair Code ex. symbol := e.GetSymbolByPair(p)
Step 4: Modify API Path(strRequestUrl)
Step 5: Add Params - Depend on API request
Step 6: Convert the response to Standard Maker struct*/
func (e *Okexdm) OrderBook(p *pair.Pair) (*exchange.Maker, error) {
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(p)

	mapParams := make(map[string]string)
	mapParams["instrument_id"] = symbol
	mapParams["size"] = "200"

	strRequestPath := "/api/futures/v3/instruments/" + symbol + "/book"
	strUrl := API_URL + strRequestPath

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Result Unmarshal Err: %v %s", e.GetName(), err, jsonOrderbook)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)

	var err error
	for _, bid := range orderBook.Bids {
		buydata := exchange.Order{}
		buydata.Quantity, _ = strconv.ParseFloat(bid[1], 64)
		buydata.Rate, _ = strconv.ParseFloat(bid[0], 64)
		maker.Bids = append(maker.Bids, buydata)
	}

	for _, ask := range orderBook.Asks {
		selldata := exchange.Order{}
		selldata.Quantity, _ = strconv.ParseFloat(ask[1], 64)
		selldata.Rate, _ = strconv.ParseFloat(ask[0], 64)
		maker.Asks = append(maker.Asks, selldata)
	}

	return maker, err
}

func (e *Okexdm) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *Okexdm) DoAccountOperation(operation *exchange.AccountOperation) error {
	return nil
}

func (e *Okexdm) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalances{}

	strRequestPath := "/API Path"

	jsonBalanceReturn := e.ApiKeyGet(strRequestPath, make(map[string]string))
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if !jsonResponse.Success {
		log.Printf("%s UpdateAllBalances Failed: %v", e.GetName(), jsonResponse.Message)
		return
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return
	}

	for _, balance := range accountBalance {
		c := e.GetCoinBySymbol(balance.Asset)
		if c != nil {
			balanceMap.Set(c.Code, balance.Available)
		}
	}
}

/* Withdraw(coin *coin.Coin, quantity float64, addr, tag string) */
func (e *Okexdm) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return false
	}

	jsonResponse := &JsonResponse{}
	withdraw := WithdrawResponse{}
	strRequestPath := "/API Path"

	mapParams := make(map[string]string)
	mapParams["asset"] = e.GetSymbolByCoin(coin)
	mapParams["address"] = addr
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["timestamp"] = fmt.Sprintf("%d", time.Now().UnixNano()/1e6)

	jsonSubmitWithdraw := e.ApiKeyRequest("POST", strRequestPath, mapParams)
	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &jsonResponse); err != nil {
		log.Printf("%s Withdraw Json Unmarshal Err: %v %v", e.GetName(), err, jsonSubmitWithdraw)
		return false
	} else if !jsonResponse.Success {
		log.Printf("%s Withdraw Failed: %v", e.GetName(), jsonResponse.Message)
		return false
	}
	if err := json.Unmarshal(jsonResponse.Data, &withdraw); err != nil {
		log.Printf("%s Withdraw Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return false
	}

	return true
}

func (e *Okexdm) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequestPath := "/API Path"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "SELL"
	mapParams["type"] = "LIMIT"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["quantity"] = strconv.FormatFloat(quantity, 'f', -1, 64)

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequestPath, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if !jsonResponse.Success {
		return nil, fmt.Errorf("%s LimitSell Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      placeOrder.OrderID,
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}
	return order, nil
}

func (e *Okexdm) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequestPath := "/API Path"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "BUY"
	mapParams["type"] = "LIMIT"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["quantity"] = strconv.FormatFloat(quantity, 'f', -1, 64)

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequestPath, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if !jsonResponse.Success {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      placeOrder.OrderID,
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}
	return order, nil
}

func (e *Okexdm) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	orderStatus := PlaceOrder{}
	strRequestPath := "/API Path"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(order.Pair)
	mapParams["orderId"] = order.OrderID

	jsonOrderStatus := e.ApiKeyGet(strRequestPath, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if !jsonResponse.Success {
		return fmt.Errorf("%s OrderStatus Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	if orderStatus.Status == "CANCELED" {
		order.Status = exchange.Cancelled
	} else if orderStatus.Status == "FILLED" {
		order.Status = exchange.Filled
	} else if orderStatus.Status == "PARTIALLY_FILLED" {
		order.Status = exchange.Partial
	} else if orderStatus.Status == "REJECTED" {
		order.Status = exchange.Rejected
	} else if orderStatus.Status == "Expired" {
		order.Status = exchange.Expired
	} else if orderStatus.Status == "NEW" {
		order.Status = exchange.New
	} else {
		order.Status = exchange.Other
	}

	order.DealRate, _ = strconv.ParseFloat(orderStatus.AveragePrice, 64)
	order.DealQuantity, _ = strconv.ParseFloat(orderStatus.ExecutedQty, 64)

	return nil
}

func (e *Okexdm) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Okexdm) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	cancelOrder := PlaceOrder{}
	strRequestPath := "/API Path"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(order.Pair)
	mapParams["orderId"] = order.OrderID

	jsonCancelOrder := e.ApiKeyRequest("DELETE", strRequestPath, mapParams)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if !jsonResponse.Success {
		return fmt.Errorf("%s CancelOrder Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Data, &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Okexdm) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Get Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Okexdm) ApiKeyGet(strRequestPath string, mapParams map[string]string) string {
	mapParams["signature"] = exchange.ComputeHmac256NoDecode(exchange.Map2UrlQuery(mapParams), e.API_SECRET)

	payload := exchange.Map2UrlQuery(mapParams)
	strUrl := API_URL + strRequestPath + "?" + payload

	request, err := http.NewRequest("GET", strUrl, nil)
	if nil != err {
		return err.Error()
	}
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	request.Header.Add("X-MBX-APIKEY", e.API_KEY)

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

/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request*/
func (e *Okexdm) ApiKeyRequest(strMethod, strRequestPath string, mapParams map[string]string) string {
	strUrl := API_URL + strRequestPath

	mapParams["signature"] = exchange.ComputeHmac256NoDecode(exchange.Map2UrlQuery(mapParams), e.API_SECRET)
	jsonParams := ""
	if nil != mapParams {
		bytesParams, _ := json.Marshal(mapParams)
		jsonParams = string(bytesParams)
	}

	request, err := http.NewRequest(strMethod, strUrl, bytes.NewBuffer([]byte(jsonParams)))
	if nil != err {
		return err.Error()
	}
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	request.Header.Add("X-MBX-APIKEY", e.API_KEY)

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

func GetContractByCode(code string) string {
	if code[0:2] == "CW" {
		thisweek := time.Now().UTC().Add(time.Duration(16) * time.Hour)
		for i := 0; i < 7; i++ {
			if thisweek.Add(time.Duration(24*i)*time.Hour).Weekday() != time.Friday {
				continue
			}
			thisweek = thisweek.Add(time.Duration(24*i) * time.Hour)
		}
		return code[2:5] + "-USD-" + thisweek.Format("060102")
	}
	if code[0:2] == "NW" {
		nextweek := time.Now().UTC().Add(time.Duration(16) * time.Hour)
		for i := 7; i < 14; i++ {
			if nextweek.Add(time.Duration(24*i)*time.Hour).Weekday() != time.Friday {
				continue
			}
			nextweek = nextweek.Add(time.Duration(24*i) * time.Hour)
		}
		return code[2:5] + "-USD-" + nextweek.Format("060102")
	}
	if code[0:2] == "CQ" {
		return code[2:5] + "-USD-" + "190927"
	}
	return "CONTRACT DNE"
}

func GetCodeByDate(code string) string {
	thisweek := time.Now().UTC().Add(time.Duration(16) * time.Hour)
	// UTC 8AM OK RESETTLES, + 16 HOURS TO GET TO THE NEXT DAY TO MAKE IT EASY
	for i := 0; i < 7; i++ {
		// SEARCH FOR A FRIDAY FROM THE NEXT 7 DAYS
		if thisweek.Add(time.Duration(24*i)*time.Hour).Weekday() != time.Friday {
			continue
		}
		thisweek = thisweek.Add(time.Duration(24*i) * time.Hour)
	}
	nextweek := time.Now().UTC().Add(time.Duration(16) * time.Hour)
	for i := 7; i < 14; i++ {
		// SEARCH FOR A FRIDAY FROM ONE WEEK LATER FOR A WEEK
		if nextweek.Add(time.Duration(24*i)*time.Hour).Weekday() != time.Friday {
			continue
		}
		nextweek = nextweek.Add(time.Duration(24*i) * time.Hour)
	}
	cw := thisweek.Format("060102")
	nw := nextweek.Format("060102")
	if code[8:14] == cw {
		return "CW" + code[0:3]
	}
	if code[8:14] == nw {
		return "NW" + code[0:3]
	}
	if code[8:14] == "190927" {
		return "CQ" + code[0:3]
	}
	return "CONTRACT DNE"
}

/* thisweek := time.Now().UTC().Add(time.Duration(16) * time.Hour)
for i := 0; i < 7; i++ {
	if thisweek.Add(time.Duration(24*i)*time.Hour).Weekday() != time.Friday {
		continue
	}
	fmt.Println("looped for " + strconv.Itoa(i) + " days")
	thisweek = thisweek.Add(time.Duration(24*i) * time.Hour)
}
fmt.Println("CW IS " + thisweek.Format("060102")) */

/* nextweek := time.Now().UTC().Add(time.Duration(16) * time.Hour)
for i := 7; i < 14; i++ {
	if nextweek.Add(time.Duration(24*i)*time.Hour).Weekday() != time.Friday {
		continue
	}
	fmt.Println("looped for " + strconv.Itoa(i) + " days")
	nextweek = nextweek.Add(time.Duration(24*i) * time.Hour)
}
fmt.Println("NW IS " + nextweek.Format("060102")) */

/* monthcount := time.Now().UTC().Add(time.Duration(16) * time.Hour)
fmt.Println("initial day is " + monthcount.Format("060102"))
for i := 1; i < 15; i++ {
	mm := monthcount.Add(time.Duration(168) * time.Hour).Month()
	if mm == time.April || mm == time.July || mm == time.October || mm == time.January {
		fmt.Println("looped for " + strconv.Itoa(i) + " weeks")
		fmt.Println("exit day is " + monthcount.Format("060102"))
		fmt.Println("exit is the " + strconv.Itoa(int(mm)) + "th month")
		break
	}
	monthcount = monthcount.Add(time.Duration(168) * time.Hour)
	fmt.Println(monthcount.Format("060102"))
}
fmt.Println("==" + monthcount.Format("060102"))

quarterday := monthcount.Add(time.Duration(168) * time.Hour)
for i := -7; i < 0; i++ {
	if quarterday.Add(time.Duration(24*i)*time.Hour).Weekday() != time.Friday {
		continue
	}
	quarterday = quarterday.Add(time.Duration(24*i) * time.Hour)
}
fmt.Println("CQ IS " + quarterday.Format("060102")) */
