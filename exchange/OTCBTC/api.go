package otcbtc

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
)

const (
	API_URL string = "https://bb.otcbtc.com"
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
func (e *Otcbtc) GetCoinsData() {
	errResponse := &ErrorResponse{}
	pairsData := PairsData{}

	strRequestUrl := "/api/v2/markets"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &pairsData); err != nil {
		if err := json.Unmarshal([]byte(jsonCurrencyReturn), &errResponse); err != nil {
			log.Printf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
		} else {
			log.Printf("%s Get Coins Failed: %v", e.GetName(), errResponse.Error.Code)
		}
	}

	for _, data := range pairsData {
		base := &coin.Coin{}
		target := &coin.Coin{}
		baseSymbol := strings.Split(data.Name, "/")[1]
		targetSymbol := strings.Split(data.Name, "/")[0]

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
				TxFee:        DEFAULT_TXFEE,
				Withdraw:     DEFAULT_WITHDRAW,
				Deposit:      DEFAULT_DEPOSIT,
				Confirmation: DEFAULT_CONFIRMATION,
				Listed:       DEFAULT_LISTED,
			}
			e.SetCoinConstraint(coinConstraint)
		}

		if target != nil {
			coinConstraint := &exchange.CoinConstraint{
				CoinID:       target.ID,
				Coin:         target,
				ExSymbol:     targetSymbol,
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
func (e *Otcbtc) GetPairsData() {
	errResponse := &ErrorResponse{}
	pairsData := PairsData{}

	strRequestUrl := "/api/v2/markets"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &pairsData); err != nil {
		if err := json.Unmarshal([]byte(jsonSymbolsReturn), &errResponse); err != nil {
			log.Printf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
		} else {
			log.Printf("%s Get Coins Failed: %v", e.GetName(), errResponse.Error.Code)
		}
	}

	for _, data := range pairsData {
		p := &pair.Pair{}
		baseSymbol := strings.Split(data.Name, "/")[1]
		targetSymbol := strings.Split(data.Name, "/")[0]

		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(baseSymbol)
			target := coin.GetCoin(targetSymbol)
			if base != nil && target != nil {
				p = pair.GetPair(base, target)
			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(data.ID)
		}
		if p != nil {
			pairConstraint := &exchange.PairConstraint{
				PairID:      p.ID,
				Pair:        p,
				ExSymbol:    data.ID,
				MakerFee:    DEFAULT_MAKERER_FEE,
				TakerFee:    DEFAULT_TAKER_FEE,
				LotSize:     data.TradingRule.MinAmount,
				PriceFilter: data.TradingRule.MinPrice,
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
func (e *Otcbtc) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	errResponse := &ErrorResponse{}
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	strRequestUrl := "/api/v2/depth"
	strUrl := API_URL + strRequestUrl

	mapParams := make(map[string]string)
	mapParams["market"] = symbol
	mapParams["limit"] = "100"

	maker := &exchange.Maker{}
	maker.WorkerIP = exchange.GetExternalIP()
	maker.BeforeTimestamp = float64(time.Now().UnixNano() / 1e6)

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &orderBook); err != nil {
		if err := json.Unmarshal([]byte(jsonOrderbook), &errResponse); err != nil {
			log.Printf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
		} else {
			log.Printf("%s Get Coins Failed: %v", e.GetName(), errResponse.Error.Code)
		}
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)

	var err error
	for _, bid := range orderBook.Bids {
		buydata := exchange.Order{}

		//Modify according to type and structure
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

		//Modify according to type and structure
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
	return maker, nil
}

/*************** Private API ***************/
func (e *Otcbtc) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	errResponse := &ErrorResponse{}
	accountBalance := AccountBalance{}
	strRequest := "/api/v2/users/me"

	jsonBalanceReturn := e.ApiKeyGET(make(map[string]string), strRequest)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		if err := json.Unmarshal([]byte(jsonBalanceReturn), &errResponse); err != nil {
			log.Printf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		} else {
			log.Printf("%s Get Coins Failed: %v", e.GetName(), errResponse.Error.Code)
		}
	}

	if accountBalance.OtbFeeEnabled != true {
		return
	}

	for _, account := range accountBalance.Accounts {
		freeAmount, err := strconv.ParseFloat(account.Balance, 64)
		if err != nil {
			log.Printf("Parse freeAmount error: %v", err)
			return
		}
		c := e.GetCoinBySymbol(account.Currency)
		if c != nil {
			balanceMap.Set(c.Code, freeAmount)
		} else {
			c = &coin.Coin{}
			c.Code = strings.ToUpper(account.Currency)
			coin.AddCoin(c)
			balanceMap.Set(c.Code, freeAmount)
		}
	}
}

func (e *Otcbtc) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	return false
}

func (e *Otcbtc) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	errResponse := &ErrorResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/api/v2/orders"

	mapParams := make(map[string]string)
	mapParams["market"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "sell"
	mapParams["volume"] = fmt.Sprintf("%v", quantity)
	mapParams["price"] = fmt.Sprintf("%v", rate)

	jsonPlaceReturn := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		if err := json.Unmarshal([]byte(jsonPlaceReturn), &errResponse); err != nil {
			log.Printf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
		} else {
			log.Printf("%s Get Coins Failed: %v", e.GetName(), errResponse.Error.Code)
		}
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%v", placeOrder.ID),
		Rate:         rate,
		Quantity:     quantity,
		Side:         "Sell",
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Otcbtc) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	errResponse := &ErrorResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/api/v2/orders"

	mapParams := make(map[string]string)
	mapParams["market"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "buy"
	mapParams["volume"] = fmt.Sprintf("%v", quantity)
	mapParams["price"] = fmt.Sprintf("%v", rate)

	jsonPlaceReturn := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		if err := json.Unmarshal([]byte(jsonPlaceReturn), &errResponse); err != nil {
			log.Printf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
		} else {
			log.Printf("%s Get Coins Failed: %v", e.GetName(), errResponse.Error.Code)
		}
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%v", placeOrder.ID),
		Rate:         rate,
		Quantity:     quantity,
		Side:         "Buy",
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Otcbtc) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	errResponse := &ErrorResponse{}
	orderStatus := PlaceOrder{}
	strRequest := "/api/v2/order"

	mapParams := make(map[string]string)
	mapParams["id"] = order.OrderID

	jsonOrderStatus := e.ApiKeyGET(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &orderStatus); err != nil {
		if err := json.Unmarshal([]byte(jsonOrderStatus), &errResponse); err != nil {
			log.Printf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
		} else {
			log.Printf("%s Get Coins Failed: %v", e.GetName(), errResponse.Error.Code)
		}
	}

	order.StatusMessage = jsonOrderStatus
	if orderStatus.State == "cancel" {
		order.Status = exchange.Canceled
	} else if orderStatus.State == "done" {
		order.Status = exchange.Filled
	} else if orderStatus.State == "wait" {
		order.Status = exchange.New
	} else {
		order.Status = exchange.Other
	}

	return nil
}

func (e *Otcbtc) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Otcbtc) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	errResponse := &ErrorResponse{}
	cancelOrder := PlaceOrder{}
	strRequest := "/api/v2/order/delete"

	mapParams := make(map[string]string)
	mapParams["id"] = order.OrderID

	jsonCancelOrder := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &cancelOrder); err != nil {
		if err := json.Unmarshal([]byte(jsonCancelOrder), &errResponse); err != nil {
			log.Printf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
		} else {
			log.Printf("%s Get Coins Failed: %v", e.GetName(), errResponse.Error.Code)
		}
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Otcbtc) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Otcbtc) ApiKeyGET(mapParams map[string]string, strRequestPath string) string {

	timestamp := strconv.FormatInt((time.Now().UTC().UnixNano() / int64(time.Millisecond)), 10)
	strUrl := API_URL + strRequestPath

	payload := "GET|" + strRequestPath + "|" + exchange.Map2UrlQuery(mapParams)

	mapParams["access_key"] = e.API_KEY
	mapParams["nonce"] = timestamp
	mapParams["signature"] = exchange.ComputeHmac256NoDecode(payload, e.API_SECRET)

	return exchange.HttpGetRequest(strUrl, mapParams)
}

func (e *Otcbtc) ApiKeyPost(strRequestPath string, mapParams map[string]string) string {
	httpClient := &http.Client{}
	timestamp := strconv.FormatInt((time.Now().UTC().UnixNano() / int64(time.Millisecond)), 10)
	strUrl := API_URL + strRequestPath

	payload := "POST|" + strRequestPath + "|" + exchange.Map2UrlQuery(mapParams)

	mapParams["access_key"] = e.API_KEY
	mapParams["nonce"] = timestamp
	mapParams["signature"] = exchange.ComputeHmac256NoDecode(payload, e.API_SECRET)

	request, err := http.NewRequest("POST", strUrl, strings.NewReader(exchange.Map2UrlQuery(mapParams)))
	if err != nil {
		return err.Error()
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Accept", "application/json")

	response, err := httpClient.Do(request)
	if err != nil {
		return err.Error()
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err.Error()
	}

	return string(body)
}
