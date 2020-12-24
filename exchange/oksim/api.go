package oksim

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

const (
	API_URL string = "https://www.okex.com"
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
func (e *Oksim) GetCoinsData() error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}

	mapParams := make(map[string]interface{})
	mapParams["instType"] = "SPOT"

	strRequestUrl := "/api/v5/public/instruments"
	jsonCurrencyReturn := e.ApiKeyRequest("GET", mapParams, strRequestUrl)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Coins Result Unmarshal Err: %v %s", e.GetName(), err, []byte(jsonCurrencyReturn))
	} else if jsonResponse.Code != "0" {
		return fmt.Errorf("%s Coins Failed: %v", e.GetName(), jsonResponse.Msg)
	}

	coinsData := []*PairData{}
	if err := json.Unmarshal(jsonResponse.Data, &coinsData); err != nil {
		return fmt.Errorf("%s Coins Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, data := range coinsData {
		base := &coin.Coin{}
		target := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base = coin.GetCoin(data.CtValCcy)
			if base == nil {
				base = &coin.Coin{}
				base.Code = data.CtValCcy
				coin.AddCoin(base)
			}
			target = coin.GetCoin(data.SettleCcy)
			if target == nil {
				target = &coin.Coin{}
				target.Code = data.SettleCcy
				coin.AddCoin(target)
			}

		case exchange.JSON_FILE:
			base = e.GetCoinBySymbol(data.CtValCcy) //data.Currency)
			target = e.GetCoinBySymbol(data.SettleCcy)
		}

		if base != nil {
			coinConstraint := e.GetCoinConstraint(base)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       base.ID,
					Coin:         base,
					ExSymbol:     data.CtValCcy,
					ChainType:    exchange.MAINNET,
					Confirmation: DEFAULT_CONFIRMATION,
				}
			} else {
				coinConstraint.ExSymbol = data.CtValCcy
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
func (e *Oksim) GetPairsData() error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}

	mapParams := make(map[string]interface{})
	mapParams["instType"] = "SPOT"

	strRequestUrl := "/api/v5/public/instruments"
	jsonSymbolsReturn := e.ApiKeyRequest("GET", mapParams, strRequestUrl)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Pairs Result Unmarshal Err: %v %s", e.GetName(), err, []byte(jsonSymbolsReturn))
	} else if jsonResponse.Code != "0" {
		return fmt.Errorf("%s Get Pairs Failed: %v", e.GetName(), jsonResponse.Msg)
	}

	pairsData := []*PairData{}
	if err := json.Unmarshal(jsonResponse.Data, &pairsData); err != nil {
		return fmt.Errorf("%s Pairs Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, data := range pairsData {
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(data.CtValCcy)
			target := coin.GetCoin(data.SettleCcy)
			if base != nil && target != nil {
				p = pair.GetPair(base, target)
			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(data.InstID)
		}

		lotSize, err := strconv.ParseFloat(data.LotSz, 64)
		if err != nil {
			return fmt.Errorf("%s Convert lotSize to Float64 Err: %v %v", e.GetName(), err, data.LotSz)
		}
		priceFilter, err := strconv.ParseFloat(data.TickSz, 64)
		if err != nil {
			return fmt.Errorf("%s Convert lotSize to Float64 Err: %v %v", e.GetName(), err, data.TickSz)
		}
		minTrade, err := strconv.ParseFloat(data.MinSz, 64)
		if err != nil {
			return fmt.Errorf("%s Convert minTrade to Float64 Err: %v %v", e.GetName(), err, data.MinSz)
		}

		if p != nil {
			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{ // no minBaseQuantity
					PairID:           p.ID,
					Pair:             p,
					ExSymbol:         data.InstID,
					MakerFee:         DEFAULT_MAKER_FEE,
					TakerFee:         DEFAULT_TAKER_FEE,
					LotSize:          lotSize,
					PriceFilter:      priceFilter,
					MinTradeQuantity: minTrade,
					Listed:           DEFAULT_LISTED,
				}
			} else {
				pairConstraint.ExSymbol = data.InstID
				pairConstraint.LotSize = lotSize
				pairConstraint.PriceFilter = priceFilter
				pairConstraint.MinTradeQuantity = minTrade
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
func (e *Oksim) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return nil, fmt.Errorf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
	}

	symbol := e.GetSymbolByPair(pair)

	jsonResponse := &JsonResponse{}

	mapParams := make(map[string]interface{})
	mapParams["instId"] = symbol

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	strRequestUrl := "/api/v5/market/books"
	jsonOrderbook := e.ApiKeyRequest("GET", mapParams, strRequestUrl)
	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s Get OrderBook Result Unmarshal Err: %v %s", e.GetName(), err, []byte(jsonOrderbook))
	} else if jsonResponse.Code != "0" {
		return nil, fmt.Errorf("%s Get OrderBook Failed: %v", e.GetName(), jsonResponse.Msg)
	}

	orderBooks := []*OrderBook{}
	if err := json.Unmarshal(jsonResponse.Data, &orderBooks); err != nil {
		return nil, fmt.Errorf("%s OrderBook Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	for _, orderBook := range orderBooks {
		for _, bid := range orderBook.Bids {
			var buydata exchange.Order
			buydata.Rate, _ = strconv.ParseFloat(bid[0], 64)
			buydata.Quantity, _ = strconv.ParseFloat(bid[1], 64)

			maker.Bids = append(maker.Bids, buydata)
		}
		for _, ask := range orderBook.Asks {
			var selldata exchange.Order
			selldata.Rate, _ = strconv.ParseFloat(ask[0], 64)
			selldata.Quantity, _ = strconv.ParseFloat(ask[1], 64)

			maker.Asks = append(maker.Asks, selldata)
		}
	}
	return maker, nil
}

/*************** Private API ***************/

func (e *Oksim) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		log.Printf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}

	strRequest := "/api/v5/asset/balances"
	jsonBalanceReturn := e.ApiKeyRequest("GET", nil, strRequest)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if jsonResponse.Code != "0" {
		log.Printf("%s UpdateAllBalances Err: Code: %v Msg: %v", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
		return
	}

	accountBalance := []*AccountBalances{}
	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return
	} else if len(accountBalance) == 0 {
		log.Printf("%s UpdateAllBalances Failed: %v", e.GetName(), jsonBalanceReturn)
		return
	}

	for _, v := range accountBalance {
		c := e.GetCoinBySymbol(v.Ccy)
		if c != nil {
			balanceAvailable, err := strconv.ParseFloat(v.AvailBal, 64)
			if err != nil {
				log.Printf("%s available balance conver to float64 err : %v", e.GetName, err)
				balanceAvailable = 0.0
			}
			balanceMap.Set(c.Code, balanceAvailable)
		}
	}
}

func (e *Oksim) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		log.Printf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
		return false
	}

	return false
}

func (e *Oksim) Transfer(coin *coin.Coin, quantity float64, from, to int) bool {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		log.Printf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
		return false
	}

	return true
}

func (e *Oksim) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return nil, fmt.Errorf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	strRequest := "/api/v5/trade/order"

	mapParams := make(map[string]interface{})
	mapParams["instId"] = e.GetSymbolByPair(pair)
	mapParams["tdMode"] = "cash"
	mapParams["side"] = "sell"
	mapParams["ordType"] = "limit"
	mapParams["px"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["sz"] = strconv.FormatFloat(quantity, 'f', -1, 64)

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != "0" {
		return nil, fmt.Errorf("%s LimitSell Failed: Code: %v Msg: %v", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
	}

	placeOrder := []*PlaceOrder{}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	} else if len(placeOrder) == 0 {
		return nil, fmt.Errorf("%s LimitSell Failed: %v", e.GetName(), jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      placeOrder[0].OrdID,
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Oksim) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return nil, fmt.Errorf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	strRequest := "/api/v5/trade/order"

	mapParams := make(map[string]interface{})
	mapParams["instId"] = e.GetSymbolByPair(pair)
	mapParams["tdMode"] = "cash"
	mapParams["side"] = "buy"
	mapParams["ordType"] = "limit"
	mapParams["px"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["sz"] = strconv.FormatFloat(quantity, 'f', -1, 64)

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != "0" {
		return nil, fmt.Errorf("%s LimitBuy Failed: Code: %v Msg: %v", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
	}

	placeOrder := []*PlaceOrder{}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	} else if len(placeOrder) == 0 {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v", e.GetName(), jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      placeOrder[0].OrdID,
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Oksim) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	strRequest := "/api/v5/trade/order"

	mapParams := make(map[string]interface{})
	mapParams["instId"] = e.GetSymbolByPair(order.Pair)
	mapParams["ordId"] = order.OrderID

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != "0" {
		return fmt.Errorf("%s OrderStatus Failed: Code: %v Msg: %v", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
	}

	orderStatus := []*OrderStatus{}
	if err := json.Unmarshal(jsonResponse.Data, &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	} else if len(orderStatus) == 0 {
		return fmt.Errorf("%s OrderStatus Failed: %v", e.GetName(), jsonResponse.Data)
	}

	order.DealRate, _ = strconv.ParseFloat(orderStatus[0].AccFillSz, 64)
	order.DealQuantity, _ = strconv.ParseFloat(orderStatus[0].AvgPx, 64)

	if orderStatus[0].State == "live" {
		order.Status = exchange.New
	} else if orderStatus[0].State == "partially_filled" {
		order.Status = exchange.Partial
	} else if orderStatus[0].State == "filled" {
		order.Status = exchange.Filled
	} else if orderStatus[0].State == "canceled" {
		order.Status = exchange.Cancelled
	} else {
		order.Status = exchange.Other
	}

	return nil
}

func (e *Oksim) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Oksim) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	strRequest := "/api/v5/trade/cancel-order"

	mapParams := make(map[string]interface{})
	mapParams["instId"] = e.GetSymbolByPair(order.Pair)
	mapParams["ordId"] = order.OrderID

	jsonCancelOrder := e.ApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if jsonResponse.Code != "0" {
		return fmt.Errorf("%s OrderStatus Failed: Code: %v Msg: %v", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
	}

	cancelOrder := []*PlaceOrder{}
	if err := json.Unmarshal(jsonResponse.Data, &cancelOrder); err != nil {
		return fmt.Errorf("%s OrderStatus Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	} else if len(cancelOrder) == 0 {
		return fmt.Errorf("%s OrderStatus Failed: %v", e.GetName(), jsonResponse.Data)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Oksim) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Oksim) ApiKeyRequest(method string, mapParams map[string]interface{}, strRequestPath string) string {
	TimeStamp := IsoTime()

	jsonParams := ""
	var bytesParams []byte
	if len(mapParams) != 0 {
		bytesParams, _ = json.Marshal(mapParams)
		jsonParams = string(bytesParams)
	}

	var strMessage string
	if method == "GET" {
		if len(mapParams) > 0 {
			strRequestPath += "?" + exchange.Map2UrlQueryInterface(mapParams)
		}
		strMessage = TimeStamp + method + strRequestPath
		if len(mapParams) > 0 {
			strMessage += jsonParams
		}
	} else {
		strMessage = TimeStamp + method + strRequestPath + jsonParams
	}

	// log.Printf("===================strMessage: %v", strMessage)

	signature := exchange.ComputeHmac256Base64(strMessage, e.API_SECRET)
	strUrl := API_URL + strRequestPath

	httpClient := &http.Client{}
	request, err := http.NewRequest(method, strUrl, bytes.NewReader(bytesParams))
	if nil != err {
		return err.Error()
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json; charset=UTF-8")
	request.Header.Add("OK-ACCESS-KEY", e.API_KEY)
	request.Header.Add("OK-ACCESS-SIGN", signature)
	request.Header.Add("OK-ACCESS-TIMESTAMP", TimeStamp)
	request.Header.Add("OK-ACCESS-PASSPHRASE", e.Passphrase)
	request.Header.Add("x-simulated-trading", "1")

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

func IsoTime() string {
	utcTime := time.Now().UTC()
	iso := utcTime.String()
	isoBytes := []byte(iso)
	iso = string(isoBytes[:10]) + "T" + string(isoBytes[11:23]) + "Z"
	return iso
}
