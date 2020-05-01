package probit

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
	"sort"
	"strconv"
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
)

/*The Base Endpoint URL*/
const (
	API_URL = "https://api.probit.com/api/exchange/v1"
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
func (e *Probit) GetCoinsData() error {
	coinsData := CoinsData{}

	strRequestUrl := "/currency"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	}

	for _, data := range coinsData.Data {
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

		txFee, err := strconv.ParseFloat(data.WithdrawalFee, 64)
		if err != nil {
			log.Printf("%s Get Coins txfee parse Err: %v, %v", e.GetName(), err, data.WithdrawalFee)
		}

		if c != nil {
			coinConstraint := e.GetCoinConstraint(c)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       c.ID,
					Coin:         c,
					ExSymbol:     data.ID,
					ChainType:    exchange.MAINNET,
					TxFee:        txFee,
					Withdraw:     !data.WithdrawalSuspended,
					Deposit:      !data.DepositSuspended,
					Confirmation: data.MinConfirmationCount,
					Listed:       true,
				}
			} else {
				coinConstraint.ExSymbol = data.ID
				coinConstraint.TxFee = txFee
				coinConstraint.Withdraw = !data.WithdrawalSuspended
				coinConstraint.Deposit = !data.DepositSuspended
				coinConstraint.Confirmation = data.MinConfirmationCount
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
func (e *Probit) GetPairsData() error {
	pairsData := &PairsData{}

	strRequestUrl := "/market"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	}

	for _, data := range pairsData.Data {
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(data.QuoteCurrencyID)
			target := coin.GetCoin(data.BaseCurrencyID)
			if base != nil && target != nil {
				p = pair.GetPair(base, target)
			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(data.ID)
		}

		if p != nil && !data.Closed {
			makerFee, _ := strconv.ParseFloat(data.MakerFeeRate, 64)
			takerFee, _ := strconv.ParseFloat(data.TakerFeeRate, 64)
			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{
					PairID:      p.ID,
					Pair:        p,
					ExSymbol:    data.ID,
					MakerFee:    makerFee,
					TakerFee:    takerFee,
					LotSize:     math.Pow10(-1 * data.QuantityPrecision),
					PriceFilter: math.Pow10(-1 * data.CostPrecision),
					Listed:      !data.Closed,
				}
			} else {
				pairConstraint.ExSymbol = data.ID
				pairConstraint.MakerFee = makerFee
				pairConstraint.TakerFee = takerFee
				pairConstraint.LotSize = math.Pow10(-1 * data.QuantityPrecision)
				pairConstraint.PriceFilter = math.Pow10(-1 * data.CostPrecision)
				pairConstraint.Listed = !data.Closed
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
func (e *Probit) OrderBook(p *pair.Pair) (*exchange.Maker, error) {
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(p)

	strRequestUrl := "/order_book"
	strUrl := API_URL + strRequestUrl

	mapParams := make(map[string]string)
	mapParams["market_id"] = symbol

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &orderBook); err != nil {
		return nil, fmt.Errorf("%s OrderBook json Unmarshal error: %v %v", e.GetName(), err, jsonOrderbook)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)

	var err error
	for _, order := range orderBook.Data {
		if order.Side == "buy" {
			buydata := exchange.Order{}
			buydata.Quantity, err = strconv.ParseFloat(order.Quantity, 64)
			if err != nil {
				return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
			}

			buydata.Rate, err = strconv.ParseFloat(order.Price, 64)
			if err != nil {
				return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Rate error:%v", e.GetName(), err)
			}
			maker.Bids = append(maker.Bids, buydata)
		} else if order.Side == "sell" {
			selldata := exchange.Order{}
			selldata.Quantity, err = strconv.ParseFloat(order.Quantity, 64)
			if err != nil {
				return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
			}

			selldata.Rate, err = strconv.ParseFloat(order.Price, 64)
			if err != nil {
				return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Rate error:%v", e.GetName(), err)
			}
			maker.Asks = append(maker.Asks, selldata)
		}
	}

	// sort maker's price
	Sort(maker.Bids[:], func(i, j int) bool {
		return maker.Bids[i].Rate > maker.Bids[j].Rate
	})
	Sort(maker.Asks[:], func(i, j int) bool {
		return maker.Asks[i].Rate < maker.Asks[j].Rate
	})

	return maker, err
}

func Sort(slice interface{}, less func(i, j int) bool) {
	sort.Slice(slice, less)
}

func (e *Probit) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *Probit) DoAccountOperation(operation *exchange.AccountOperation) error {
	return nil
}

func (e *Probit) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	accountBalance := AccountBalances{}
	strRequest := "/api/v3/account"

	jsonBalanceReturn := e.ApiKeyGet(make(map[string]string), strRequest)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances json Unmarshal error: %v %s", e.GetName(), err, jsonBalanceReturn)
		return
	} else {
		for _, balance := range accountBalance.Balances {
			freeamount, err := strconv.ParseFloat(balance.Free, 64)
			if err != nil {
				log.Printf("%s UpdateAllBalances err: %+v %v", e.GetName(), balance, err)
				return
			} else {
				c := e.GetCoinBySymbol(balance.Asset)
				if c != nil {
					balanceMap.Set(c.Code, freeamount)
				}
			}
		}
	}
}

/* Withdraw(coin *coin.Coin, quantity float64, addr, tag string) */
func (e *Probit) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return false
	}

	withdraw := WithdrawResponse{}
	strRequest := "/wapi/v3/withdraw.html"

	mapParams := make(map[string]string)
	mapParams["asset"] = e.GetSymbolByCoin(coin)
	mapParams["address"] = addr
	if tag != "" { //this part is not working yet
		mapParams["addressTag"] = tag
	}
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["timestamp"] = fmt.Sprintf("%d", time.Now().UnixNano()/1e6)

	jsonSubmitWithdraw := e.WApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &withdraw); err != nil {
		log.Printf("%s Withdraw Json Unmarshal Error: %v %v", e.GetName(), err, jsonSubmitWithdraw)
		return false
	}
	if !withdraw.Success {
		log.Printf("%s Withdraw Failed: %s", e.GetName(), withdraw.Msg)
		return false
	}

	return true
}

func (e *Probit) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	placeOrder := PlaceOrder{}
	strRequest := "/api/v3/order"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "SELL"
	mapParams["type"] = "LIMIT"
	mapParams["timeInForce"] = "GTC"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["quantity"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if placeOrder.Code != 0 {
		return nil, fmt.Errorf("%s LimitSell failed:%v Message:%v", e.GetName(), placeOrder.Code, placeOrder.Msg)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%d", placeOrder.OrderID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}
	return order, nil
}

func (e *Probit) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	placeOrder := PlaceOrder{}
	strRequest := "/api/v3/order"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "BUY"
	mapParams["type"] = "LIMIT"
	mapParams["timeInForce"] = "GTC"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["quantity"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if placeOrder.Code != 0 {
		return nil, fmt.Errorf("%s LimitBuy failed:%v Message:%v", e.GetName(), placeOrder.Code, placeOrder.Msg)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%d", placeOrder.OrderID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Probit) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	orderStatus := PlaceOrder{}
	strRequest := "/api/v3/order"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(order.Pair)
	mapParams["orderId"] = order.OrderID

	jsonOrderStatus := e.ApiKeyGet(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if orderStatus.Code != 0 {
		return fmt.Errorf("%s Get OrderStatus Error: %v %s", e.GetName(), orderStatus.Code, orderStatus.Msg)
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

	order.DealRate, _ = strconv.ParseFloat(orderStatus.Price, 64)
	order.DealQuantity, _ = strconv.ParseFloat(orderStatus.ExecutedQty, 64)

	return nil
}

func (e *Probit) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Probit) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	cancelOrder := PlaceOrder{}
	strRequest := "/api/v3/order"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(order.Pair)
	mapParams["orderId"] = order.OrderID

	jsonCancelOrder := e.ApiKeyRequest("DELETE", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if cancelOrder.Code != 0 {
		return fmt.Errorf("%s CancelOrder Error: %v %s", e.GetName(), cancelOrder.Code, cancelOrder.Msg)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Probit) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Probit) ApiKeyGet(mapParams map[string]string, strRequestPath string) string {
	mapParams["recvWindow"] = "50000" //"50000000"
	mapParams["timestamp"] = fmt.Sprintf("%d", time.Now().UTC().UnixNano()/int64(time.Millisecond))
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
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Probit) ApiKeyRequest(strMethod string, mapParams map[string]string, strRequestPath string) string {
	mapParams["recvWindow"] = "50000" //"50000000"
	mapParams["timestamp"] = fmt.Sprintf("%d", time.Now().UTC().UnixNano()/int64(time.Millisecond))
	mapParams["signature"] = exchange.ComputeHmac256NoDecode(exchange.Map2UrlQuery(mapParams), e.API_SECRET)

	payload := exchange.Map2UrlQuery(mapParams)
	strUrl := API_URL + strRequestPath

	request, err := http.NewRequest(strMethod, strUrl, bytes.NewBuffer([]byte(payload)))
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
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Probit) WApiKeyRequest(strMethod string, mapParams map[string]string, strRequestPath string) string {
	mapParams["recvWindow"] = "50000" //"50000000"
	mapParams["timestamp"] = fmt.Sprintf("%d", time.Now().UTC().UnixNano()/int64(time.Millisecond))
	signature := exchange.ComputeHmac256NoDecode(exchange.Map2UrlQuery(mapParams), e.API_SECRET)

	payload := fmt.Sprintf("%s&signature=%s", exchange.Map2UrlQuery(mapParams), signature)
	strUrl := API_URL + strRequestPath

	request, err := http.NewRequest(strMethod, strUrl, bytes.NewBuffer([]byte(payload)))
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
