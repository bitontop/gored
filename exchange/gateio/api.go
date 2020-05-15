package gateio

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"crypto/hmac"
	"crypto/sha512"
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
	API_URL     string = "https://data.gate.io"
	Private_URL string = "https://api.gateio.io"
	// Private_URL string = "https://api.gateio.life"
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
func (e *Gateio) GetCoinsData() error {
	coinsData := CoinsData{}

	strRequestUrl := "/api2/1/marketlist"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if coinsData.Result != "true" {
		return fmt.Errorf("%s Get Coins Failed: %+v", e.GetName(), coinsData)
	}

	for _, data := range coinsData.Data {
		baseSymbol := strings.TrimSpace(data.CurrSuffix)
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
			target = coin.GetCoin(data.Symbol)
			if target == nil {
				target = &coin.Coin{}
				target.Code = data.Symbol
				target.Name = data.Name
				coin.AddCoin(target)
			}
		case exchange.JSON_FILE:
			base = e.GetCoinBySymbol(baseSymbol)
			target = e.GetCoinBySymbol(data.Symbol)
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
					Listed:       DEFAULT_LISTED,
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
					ExSymbol:     data.Symbol,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = data.Symbol
			}
			e.SetCoinConstraint(coinConstraint)
		}
	}

	return e.SetConstraint()
}

func (e *Gateio) SetConstraint() error {
	coinsConstrain := CoinsConstrain{}

	strRequestUrl := "/api2/1/coininfo"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &coinsConstrain); err != nil {
		return fmt.Errorf("%s Get Coins' Constraint Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if coinsConstrain.Result != "true" {
		return fmt.Errorf("%s Get Coins' Constraint Failed: %+v", e.GetName(), coinsConstrain)
	}

	for _, coins := range coinsConstrain.Coins {
		for symbol, data := range coins {
			c := e.GetCoinBySymbol(symbol)
			if c != nil {
				constrain := e.GetCoinConstraint(c)
				constrain.Deposit = data.DepositDisabled == 0
				constrain.Withdraw = data.WithdrawDisabled == 0
				constrain.Listed = data.Delisted == 0
			}
		}
	}
	return nil
}

/* GetPairsData - Get Pairs Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func (e *Gateio) GetPairsData() error {
	pairsData := PairsData{}

	strRequestUrl := "/api2/1/marketinfo"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if pairsData.Result != "true" {
		return fmt.Errorf("%s Get Pairs Failed: %+v", e.GetName(), pairsData)
	}

	for _, pairs := range pairsData.Pairs {
		for symbol, data := range pairs {
			p := &pair.Pair{}
			switch e.Source {
			case exchange.EXCHANGE_API:
				pairSymbol := strings.Split(symbol, "_")
				base := coin.GetCoin(pairSymbol[1])
				target := coin.GetCoin(pairSymbol[0])
				if base != nil && target != nil {

					p = pair.GetPair(base, target)

				}
			case exchange.JSON_FILE:
				p = e.GetPairBySymbol(symbol)
			}

			priceFilter := math.Pow10(data.DecimalPlaces * -1)
			if p != nil {
				pairConstraint := e.GetPairConstraint(p)
				if pairConstraint == nil {
					pairConstraint = &exchange.PairConstraint{
						PairID:      p.ID,
						Pair:        p,
						ExSymbol:    symbol,
						MakerFee:    data.Fee / 100,
						TakerFee:    data.Fee / 100,
						LotSize:     DEFAULT_LOT_SIZE,
						PriceFilter: priceFilter,
						Listed:      DEFAULT_LISTED,
					}
				} else {
					pairConstraint.ExSymbol = symbol
					pairConstraint.MakerFee = data.Fee / 100
					pairConstraint.TakerFee = data.Fee / 100
					pairConstraint.PriceFilter = priceFilter
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
Step 3: Get Exchange Pair Code ex. symbol := e.GetPairCode(p)
Step 4: Modify API Path(strRequestUrl)
Step 5: Add Params - Depend on API request
Step 6: Convert the response to Standard Maker struct*/
func (e *Gateio) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	strRequestUrl := fmt.Sprintf("/api2/1/orderBook/%s", symbol)
	strUrl := API_URL + strRequestUrl

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonOrderbook), &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if orderBook.Result != "true" {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %+v", e.GetName(), orderBook)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	var err error
	for _, bid := range orderBook.Bids {
		var buydata exchange.Order

		//Modify according to type and structure
		buydata.Rate, err = strconv.ParseFloat(bid[0], 64)
		if err != nil {
			return nil, fmt.Errorf("GateIo Bids Rate ParseFloat error:%v", err)
		}
		buydata.Quantity, err = strconv.ParseFloat(bid[1], 64)
		if err != nil {
			return nil, fmt.Errorf("GateIo Bids Quantity ParseFloat error:%v", err)
		}

		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Asks {
		var selldata exchange.Order

		//Modify according to type and structure
		selldata.Rate, err = strconv.ParseFloat(ask[0], 64)
		if err != nil {
			return nil, fmt.Errorf("GateIo Asks Rate ParseFloat error:%v", err)
		}
		selldata.Quantity, err = strconv.ParseFloat(ask[1], 64)
		if err != nil {
			return nil, fmt.Errorf("GateIo Asks Quantity ParseFloat error:%v", err)
		}

		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

/*************** Private API ***************/
func (e *Gateio) DoAccountOperation(operation *exchange.AccountOperation) error {
	switch operation.Type {
	// case exchange.BalanceList:
	// 	return e.getAllBalance(operation)
	// case exchange.Balance:
	// 	return e.getBalance(operation)
	case exchange.Withdraw:
		return e.doWithdraw(operation)
	}
	return fmt.Errorf("%s Operation type invalid: %s %v", operation.Ex, operation.Wallet, operation.Type)
}

// need to add address to address book, or set TOTP or add phone number
func (e *Gateio) doWithdraw(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
	}

	withdrawResponse := WithdrawResponse{}
	strRequest := "/api2/1/private/withdraw" //"https://api.gateio.life/api2/1/private/withdraw"

	mapParams := make(map[string]string)
	mapParams["currency"] = e.GetSymbolByCoin(operation.Coin)
	mapParams["amount"] = operation.WithdrawAmount
	mapParams["address"] = operation.WithdrawAddress
	if operation.WithdrawTag != "" {
		mapParams["address"] += fmt.Sprintf(" %v", operation.WithdrawTag)
	}

	log.Printf("mapParams: %v", mapParams) //=================

	jsonSubmitWithdraw := e.ApiKeyPost(strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonSubmitWithdraw
	}

	log.Printf("Withdraw: %v", jsonSubmitWithdraw) // ===========================

	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &withdrawResponse); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Json Unmarshal Err: %v, %s", e.GetName(), err, jsonSubmitWithdraw)
		return operation.Error
	} else if withdrawResponse.Result != "true" || withdrawResponse.Message != "Success" {
		operation.Error = fmt.Errorf("%s Withdraw Failed: %v", e.GetName(), jsonSubmitWithdraw)
		return operation.Error
	} /* else if withdrawResponse.WithdrawalID == "" {
		operation.Error = fmt.Errorf("%s Withdraw Failed: %v", e.GetName(), jsonSubmitWithdraw)
		return operation.Error
	} */

	// operation.WithdrawID = withdrawResponse.WithdrawalID

	return nil
}

func (e *Gateio) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	accountBalance := AccountBalances{}
	freeBalance := make(map[string]string)
	strRequest := "/api2/1/private/balances"

	jsonBalanceReturn := e.ApiKeyPost(strRequest, make(map[string]string))
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if accountBalance.Result != "true" {
		log.Printf("%s UpdateAllBalances Failed: %+v", e.GetName(), accountBalance)
		return
	}
	if string(accountBalance.Available) != "[]" {
		if err := json.Unmarshal(accountBalance.Available, &freeBalance); err != nil {
			log.Printf("%s UpdateAllBalances Data Unmarshal Err: %v %v", e.GetName(), err, accountBalance.Available)
			return
		}
	} else {
		return
	}

	for key, data := range freeBalance {
		freeAmount, err := strconv.ParseFloat(data, 64)
		if err != nil {
			log.Printf("%s freeAmount parse Failed: %v", e.GetName(), data)
			return
		}
		c := e.GetCoinBySymbol(key)
		if c != nil {
			balanceMap.Set(c.Code, freeAmount)
		}
	}
}

func (e *Gateio) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {

	return false
}

func (e *Gateio) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	placeOrder := PlaceOrder{}
	strRequest := "/api2/1/private/sell"

	mapParams := make(map[string]string)
	mapParams["currencyPair"] = e.GetSymbolByPair(pair)
	mapParams["rate"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["orderType"] = "" //ioc: immediate order cancel

	jsonPlaceReturn := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if placeOrder.Result != "true" {
		return nil, fmt.Errorf("%s LimitSell Failed: %+v", e.GetName(), placeOrder)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%d", placeOrder.OrderNumber),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Gateio) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	placeOrder := PlaceOrder{}
	strRequest := "/api2/1/private/buy"

	mapParams := make(map[string]string)
	mapParams["currencyPair"] = e.GetSymbolByPair(pair)
	mapParams["rate"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["orderType"] = "" //ioc: immediate order cancel

	jsonPlaceReturn := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if placeOrder.Result != "true" {
		return nil, fmt.Errorf("%s LimitBuy Failed: %+v", e.GetName(), placeOrder)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%d", placeOrder.OrderNumber),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Gateio) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	orderStatus := OrderStatus{}
	strRequest := "/api2/1/private/getOrder"

	mapParams := make(map[string]string)
	mapParams["orderNumber"] = order.OrderID
	mapParams["currencyPair"] = e.GetSymbolByPair(order.Pair)

	jsonOrderStatus := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if orderStatus.Result != "true" {
		return fmt.Errorf("%s OrderStatus Failed: %+v", e.GetName(), orderStatus)
	}

	order.StatusMessage = jsonOrderStatus
	if orderStatus.Order.Status == "cancelled" {
		order.Status = exchange.Cancelled
	} else if orderStatus.Order.Status == "filled" {
		order.Status = exchange.Filled
	} else if orderStatus.Order.Status == "partial-filled" || orderStatus.Order.Status == "partial-canceled" {
		order.Status = exchange.Partial
	} else if orderStatus.Order.Status == "submitting" || orderStatus.Order.Status == "submitted" {
		order.Status = exchange.New
	} else {
		order.Status = exchange.Other
	}

	return nil
}

func (e *Gateio) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Gateio) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	cancelOrder := CancelOrder{}
	strRequest := "/api2/1/private/cancelOrder"

	mapParams := make(map[string]string)
	mapParams["orderNumber"] = order.OrderID
	mapParams["currencyPair"] = e.GetSymbolByPair(order.Pair)

	jsonCancelOrder := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if !cancelOrder.Result {
		return fmt.Errorf("%s CancelOrder Failed: %v", e.GetName(), cancelOrder.Message)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Gateio) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Gateio) ApiKeyPost(strRequestPath string, mapParams map[string]string) string {
	strMethod := "POST"

	payload := exchange.Map2UrlQuery(mapParams)
	Signature := ComputeHmac512(payload, e.API_SECRET)

	strUrl := Private_URL + strRequestPath

	httpClient := &http.Client{}

	request, err := http.NewRequest(strMethod, strUrl, strings.NewReader(payload))
	if nil != err {
		return err.Error()
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("key", e.API_KEY)
	request.Header.Set("sign", Signature)

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

func (e *Gateio) WithdrawKeyRequest(strRequestPath string, mapParams map[string]string) string {
	strMethod := "POST"

	payload := exchange.Map2UrlQuery(mapParams)
	Signature := ComputeHmac512(payload, e.API_SECRET)

	strUrl := /* Private_URL + */ strRequestPath

	httpClient := &http.Client{}

	request, err := http.NewRequest(strMethod, strUrl, strings.NewReader(payload))
	if nil != err {
		return err.Error()
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("key", e.API_KEY)
	request.Header.Set("sign", Signature)

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

func ComputeHmac512(strMessage string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha512.New, key)
	h.Write([]byte(strMessage))

	return fmt.Sprintf("%x", h.Sum(nil))
}
