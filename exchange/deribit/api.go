package deribit

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
	API_URL = "https://test.deribit.com/api/v2"
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
func (e *Deribit) GetCoinsData() error {
	jsonResponse := &JsonResponse{}
	contractsData := ContractsData{}

	strRequestPath := "/public/get_instruments?currency=BTC&kind=future"
	strUrl := API_URL + strRequestPath

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	}
	if err := json.Unmarshal(jsonResponse.Data, &contractsData); err != nil {
		return fmt.Errorf("%s Get Coins Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, data := range contractsData {
		c := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			c = coin.GetCoin("CQBTC")
			if c == nil && data.InstrumentName == "BTC-27SEP19" {
				c = &coin.Coin{
					Code:     "CQBTC",
					Name:     data.BaseCurrency,
					Explorer: data.SettlementPeriod,
				}
				coin.AddCoin(c)
			}
		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(data.InstrumentName)
		}

		if c != nil {
			coinConstraint := e.GetCoinConstraint(c)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       c.ID,
					Coin:         c,
					ExSymbol:     data.InstrumentName,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       data.IsActive,
				}
			} else {
				coinConstraint.ExSymbol = data.InstrumentName
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
func (e *Deribit) GetPairsData() error {
	jsonResponse := &JsonResponse{}
	contractsData := ContractsData{}

	strRequestPath := "/public/get_instruments?currency=BTC&kind=future"
	strUrl := API_URL + strRequestPath

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	}
	if err := json.Unmarshal(jsonResponse.Data, &contractsData); err != nil {
		return fmt.Errorf("%s Get Pairs Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, data := range contractsData {
		if data.IsActive && data.InstrumentName == "BTC-27SEP19" {
			p := &pair.Pair{}
			switch e.Source {
			case exchange.EXCHANGE_API:
				base := coin.GetCoin("USD") // &coin.Coin{0, "USD", "USD", "", "", "", 0, 0}
				target := coin.GetCoin("CQBTC")
				if base != nil && target != nil {
					p = pair.GetPair(base, target)
				}
			case exchange.JSON_FILE:
				p = e.GetPairBySymbol(data.InstrumentName)
			}
			if p != nil && data.IsActive {
				pairConstraint := e.GetPairConstraint(p)
				if pairConstraint == nil {
					pairConstraint = &exchange.PairConstraint{
						PairID:      p.ID,
						Pair:        p,
						ExSymbol:    data.InstrumentName,
						MakerFee:    DEFAULT_MAKER_FEE,
						TakerFee:    DEFAULT_TAKER_FEE,
						LotSize:     data.MinTradeAmount,
						PriceFilter: data.TickSize,
						Listed:      data.IsActive,
					}
				} else {
					pairConstraint.ExSymbol = data.InstrumentName
					pairConstraint.LotSize = data.MinTradeAmount
					pairConstraint.PriceFilter = data.TickSize
				}
				e.SetPairConstraint(pairConstraint)
			}
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
func (e *Deribit) OrderBook(p *pair.Pair) (*exchange.Maker, error) {
	jsonResponse := &JsonResponse{}
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(p)

	mapParams := make(map[string]string)
	mapParams["instrument_name"] = symbol

	strRequestPath := "/public/get_order_book"
	strUrl := API_URL + strRequestPath

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} /* else if !jsonResponse.Success {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %v", e.GetName(), jsonResponse.Message)
	} */
	if err := json.Unmarshal(jsonResponse.Data, &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)

	var err error
	for _, bid := range orderBook.Bids {
		buydata := exchange.Order{}
		buydata.Quantity = float64(bid[1])
		buydata.Rate = float64(bid[0])
		maker.Bids = append(maker.Bids, buydata)
	}

	for _, ask := range orderBook.Asks {
		selldata := exchange.Order{}
		selldata.Quantity = float64(ask[1])
		selldata.Rate = float64(ask[0])
		maker.Asks = append(maker.Asks, selldata)
	}

	return maker, err
}

func (e *Deribit) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *Deribit) DoAccountOperation(operation *exchange.AccountOperation) error {
	return nil
}

func (e *Deribit) UpdateAllBalances() {
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
	} /* else if !jsonResponse.Success {
		log.Printf("%s UpdateAllBalances Failed: %v", e.GetName(), jsonResponse.Message)
		return
	} */
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
func (e *Deribit) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
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
	} /* else if !jsonResponse.Success {
		log.Printf("%s Withdraw Failed: %v", e.GetName(), jsonResponse.Message)
		return false
	} */
	if err := json.Unmarshal(jsonResponse.Data, &withdraw); err != nil {
		log.Printf("%s Withdraw Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return false
	}

	return true
}

func (e *Deribit) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
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
	} /* else if !jsonResponse.Success {
		return nil, fmt.Errorf("%s LimitSell Failed: %v", e.GetName(), jsonResponse.Message)
	} */
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

func (e *Deribit) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
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
	} /* else if !jsonResponse.Success {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v", e.GetName(), jsonResponse.Message)
	} */
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

func (e *Deribit) OrderStatus(order *exchange.Order) error {
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
	} /* else if !jsonResponse.Success {
		return fmt.Errorf("%s OrderStatus Failed: %v", e.GetName(), jsonResponse.Message)
	} */
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

func (e *Deribit) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Deribit) CancelOrder(order *exchange.Order) error {
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
	} /* else if !jsonResponse.Success {
		return fmt.Errorf("%s CancelOrder Failed: %v", e.GetName(), jsonResponse.Message)
	} */
	if err := json.Unmarshal(jsonResponse.Data, &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Deribit) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Get Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Deribit) ApiKeyGet(strRequestPath string, mapParams map[string]string) string {
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
func (e *Deribit) ApiKeyRequest(strMethod, strRequestPath string, mapParams map[string]string) string {
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
