package coinbene

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"bytes"
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
	API_URL string = "http://api.coinbene.com"
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
func (e *Coinbene) GetCoinsData() {
	pairsData := PairsData{}

	strRequestUrl := "/v1/market/symbol"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &pairsData); err != nil {
		log.Printf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if pairsData.Status != "ok" {
		log.Printf("%s Get Coins Failed: %v", e.GetName(), pairsData.Description)
	}

	for _, data := range pairsData.Symbol {
		base := &coin.Coin{}
		target := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base = coin.GetCoin(data.QuoteAsset)
			if base == nil {
				base = &coin.Coin{}
				base.Code = data.QuoteAsset
				coin.AddCoin(base)
			}
			target = coin.GetCoin(data.BaseAsset)
			if target == nil {
				target = &coin.Coin{}
				target.Code = data.BaseAsset
				coin.AddCoin(target)
			}
		case exchange.JSON_FILE:
			base = e.GetCoinBySymbol(data.QuoteAsset)
			target = e.GetCoinBySymbol(data.BaseAsset)
		}

		if base != nil {
			coinConstraint := &exchange.CoinConstraint{
				CoinID:       base.ID,
				Coin:         base,
				ExSymbol:     data.QuoteAsset,
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
				ExSymbol:     data.BaseAsset,
				TxFee:        DEFAULT_TXFEE,
				Withdraw:     DEFAULT_WITHDRAW,
				Deposit:      DEFAULT_DEPOSIT,
				Confirmation: DEFAULT_CONFIRMATION,
				Listed:       DEFAULT_LISTED,
			}
			e.SetCoinConstraint(coinConstraint)
		}
	}
}

/* GetPairsData - Get Pairs Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func (e *Coinbene) GetPairsData() {
	pairsData := PairsData{}

	strRequestUrl := "/v1/market/symbol"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &pairsData); err != nil {
		log.Printf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if pairsData.Status != "ok" {
		log.Printf("%s Get Pairs Failed: %v", e.GetName(), pairsData.Description)
	}

	for _, data := range pairsData.Symbol {
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(data.QuoteAsset)
			target := coin.GetCoin(data.BaseAsset)
			if base != nil && target != nil {
				p = pair.GetPair(base, target)
			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(data.Ticker)
		}
		if p != nil {
			makerFee, err := strconv.ParseFloat(data.MakerFee, 64)
			if err != nil {
				log.Printf("%s makerFee parse error: %v, %v", e.GetName(), err, data.MakerFee)
				return
			}
			takerFee, err := strconv.ParseFloat(data.TakerFee, 64)
			if err != nil {
				log.Printf("%s takerFee parse error: %v, %v", e.GetName(), err, data.TakerFee)
				return
			}
			lotSize, err := strconv.Atoi(data.LotStepSize)
			if err != nil {
				log.Printf("%s lot size parse error: %v, %v", e.GetName(), err, data.LotStepSize)
				return
			}
			priceSize, err := strconv.Atoi(data.TickSize)
			if err != nil {
				log.Printf("%s price size parse error: %v, %v", e.GetName(), err, data.TickSize)
				return
			}
			pairConstraint := &exchange.PairConstraint{
				PairID:      p.ID,
				Pair:        p,
				ExSymbol:    data.Ticker,
				MakerFee:    makerFee,
				TakerFee:    takerFee,
				LotSize:     math.Pow10(-1 * lotSize),
				PriceFilter: math.Pow10(-1 * priceSize),
				Listed:      DEFAULT_LISTED,
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
func (e *Coinbene) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	strRequestUrl := "/v1/market/orderbook"
	strUrl := API_URL + strRequestUrl

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol

	maker := &exchange.Maker{}
	maker.WorkerIP = exchange.GetExternalIP()
	maker.BeforeTimestamp = float64(time.Now().UnixNano() / 1e6)

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if orderBook.Status != "ok" {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %v", e.GetName(), orderBook.Description)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	for _, bid := range orderBook.Orderbook.Bids {
		var buydata exchange.Order

		buydata.Rate = bid.Price
		buydata.Quantity = bid.Quantity

		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Orderbook.Asks {
		var selldata exchange.Order

		selldata.Rate = ask.Price
		selldata.Quantity = ask.Quantity

		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

/*************** Private API ***************/
func (e *Coinbene) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	accountBalance := AccountBalances{}
	strRequest := "/v1/trade/balance"

	mapParams := make(map[string]string)
	mapParams["account"] = "exchange"

	jsonBalanceReturn := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if accountBalance.Status != "ok" {
		log.Printf("%s UpdateAllBalances Failed: %v", e.GetName(), accountBalance.Description)
		return
	}

	for _, v := range accountBalance.Balance {
		freeAmount, err := strconv.ParseFloat(v.Available, 64)
		if err != nil {
			log.Printf("%s balance parse error: %v, %v", e.GetName(), err, v.Available)
			return
		}
		c := e.GetCoinBySymbol(v.Asset)
		if c != nil {
			balanceMap.Set(c.Code, freeAmount)
		}
	}
}

func (e *Coinbene) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {

	return false
}

func (e *Coinbene) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	placeOrder := PlaceOrder{}
	strRequest := "/v1/trade/order/place"

	mapParams := make(map[string]string)
	mapParams["type"] = "sell-limit"
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["price"] = fmt.Sprintf("%v", rate)
	mapParams["quantity"] = fmt.Sprintf("%v", quantity)

	jsonPlaceReturn := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if placeOrder.Status != "ok" {
		return nil, fmt.Errorf("%s LimitSell Failed: %v", e.GetName(), placeOrder.Description)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      placeOrder.Orderid,
		Rate:         rate,
		Quantity:     quantity,
		Side:         "Sell",
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Coinbene) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	placeOrder := PlaceOrder{}
	strRequest := "/v1/trade/order/place"

	mapParams := make(map[string]string)
	mapParams["type"] = "buy-limit"
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["price"] = fmt.Sprintf("%v", rate)
	mapParams["quantity"] = fmt.Sprintf("%v", quantity)

	jsonPlaceReturn := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if placeOrder.Status != "ok" {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v", e.GetName(), placeOrder.Description)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      placeOrder.Orderid,
		Rate:         rate,
		Quantity:     quantity,
		Side:         "Buy",
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Coinbene) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	orderStatus := OrderStatus{}
	strRequest := "/v1/trade/order/info"

	mapParams := make(map[string]string)
	mapParams["orderid"] = order.OrderID

	jsonOrderStatus := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if orderStatus.Status != "ok" {
		return fmt.Errorf("%s OrderStatus Failed: %v", e.GetName(), orderStatus.Description)
	}

	order.StatusMessage = jsonOrderStatus
	if orderStatus.Order.Orderstatus == "filled" {
		order.Status = exchange.Filled
	} else if orderStatus.Order.Orderstatus == "partialFilled" {
		order.Status = exchange.Partial
	} else if orderStatus.Order.Orderstatus == "canceled" || orderStatus.Order.Orderstatus == "partialCanceled" {
		order.Status = exchange.Canceled
	} else if orderStatus.Order.Orderstatus == "unfilled" {
		order.Status = exchange.New
	} else {
		order.Status = exchange.Other
	}

	return nil
}

func (e *Coinbene) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Coinbene) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	cancelOrder := PlaceOrder{}
	strRequest := "/v1/trade/order/cancel"

	mapParams := make(map[string]string)
	mapParams["orderid"] = order.OrderID

	jsonCancelOrder := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if cancelOrder.Status != "ok" {
		return fmt.Errorf("%s CancelOrder Failed: %v", e.GetName(), cancelOrder.Description)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Coinbene) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Coinbene) ApiKeyPost(strRequestPath string, mapParams map[string]string) string {
	strUrl := API_URL + strRequestPath

	//Signature Request Params
	mapParams["apiid"] = e.API_KEY
	mapParams["secret"] = e.API_SECRET
	mapParams["timestamp"] = fmt.Sprintf("%d", time.Now().UnixNano()/1000000)
	strMessage := strings.ToUpper(exchange.Map2UrlQuery(mapParams))
	mapParams["sign"] = exchange.ComputeMD5(strMessage)
	delete(mapParams, "secret")

	httpClient := &http.Client{}
	bytesParams, _ := json.Marshal(mapParams)

	request, err := http.NewRequest("POST", strUrl, bytes.NewBuffer(bytesParams))
	if nil != err {
		return err.Error()
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	request.Header.Add("Connection", "keep-alive")

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

func (e *Coinbene) ApiKeyGET(strRequestPath string, mapParams map[string]string) string {
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
