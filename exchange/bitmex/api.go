package bitmex

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

/*The Base Endpoint URL*/
const (
	API_URL = "https://www.bitmex.com/api/v1"
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
func (e *Bitmex) GetCoinsData() error {
	coinsData := PairsData{}

	strRequestUrl := "/instrument/active"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	}

	for _, data := range coinsData {
		base := &coin.Coin{}
		target := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base = coin.GetCoin(data.QuoteCurrency)
			if base == nil {
				base = &coin.Coin{}
				base.Code = data.QuoteCurrency
				coin.AddCoin(base)
			}

			target = coin.GetCoin(data.RootSymbol)
			if target == nil {
				target = &coin.Coin{}
				target.Code = data.RootSymbol
				coin.AddCoin(target)
			}
		case exchange.JSON_FILE:
			base = e.GetCoinBySymbol(data.QuoteCurrency)
			target = e.GetCoinBySymbol(data.RootSymbol)
		}

		if base != nil {
			coinConstraint := e.GetCoinConstraint(base)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       base.ID,
					Coin:         base,
					ExSymbol:     data.QuoteCurrency,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       true,
				}
			} else {
				coinConstraint.ExSymbol = data.QuoteCurrency
			}

			e.SetCoinConstraint(coinConstraint)
		}

		if target != nil {
			coinConstraint := e.GetCoinConstraint(target)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       target.ID,
					Coin:         target,
					ExSymbol:     data.RootSymbol,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       true,
				}
			} else {
				coinConstraint.ExSymbol = data.RootSymbol
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
func (e *Bitmex) GetPairsData() error {
	pairsData := &PairsData{}

	strRequestUrl := "/instrument/active"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	}

	for _, data := range *pairsData {
		if (data.QuoteCurrency == "USD" && data.Typ == "FFWCSX") || data.QuoteCurrency != "USD" {
			p := &pair.Pair{}
			switch e.Source {
			case exchange.EXCHANGE_API:
				base := coin.GetCoin(data.QuoteCurrency)
				target := coin.GetCoin(data.RootSymbol)
				if base != nil && target != nil && base != target {
					p = pair.GetPair(base, target)
				}
			case exchange.JSON_FILE:
				p = e.GetPairBySymbol(data.Symbol)
			}

			if p != nil {
				pairConstraint := e.GetPairConstraint(p)
				if pairConstraint == nil {
					pairConstraint = &exchange.PairConstraint{
						PairID:      p.ID,
						Pair:        p,
						ExSymbol:    data.Symbol,
						MakerFee:    data.MakerFee,
						TakerFee:    data.TakerFee,
						LotSize:     data.LotSize,
						PriceFilter: data.TickSize,
						Listed:      true,
					}
				} else {
					pairConstraint.ExSymbol = data.Symbol
					pairConstraint.MakerFee = data.MakerFee
					pairConstraint.TakerFee = data.TakerFee
					pairConstraint.LotSize = data.LotSize
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
Step 3: Get Exchange Pair Code ex. symbol := e.GetPairCode(p)
Step 4: Modify API Path(strRequestUrl)
Step 5: Add Params - Depend on API request
Step 6: Convert the response to Standard Maker struct*/
func (e *Bitmex) OrderBook(p *pair.Pair) (*exchange.Maker, error) {
	errResponse := ErrorResponse{}
	orderBook := OrderBook{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(p)
	mapParams["depth"] = "0"

	strRequestUrl := "/orderBook/L2"
	strUrl := API_URL + strRequestUrl

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &orderBook); err != nil {
		if err := json.Unmarshal([]byte(jsonOrderbook), &errResponse); err != nil {
			return nil, fmt.Errorf("%s OrderBook Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
		} else {
			return nil, fmt.Errorf("%s Get OrderBook Failed: %v %v", e.GetName(), errResponse.Error.Name, errResponse.Error.Message)
		}
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)

	for _, bid := range orderBook {
		if bid.Side == "Buy" {
			buydata := exchange.Order{}

			buydata.Rate = bid.Price
			buydata.Quantity = bid.Size / bid.Price

			maker.Bids = append(maker.Bids, buydata)
		}
	}
	for i := len(orderBook) - 1; i >= 0; i-- {
		if orderBook[i].Side == "Sell" {
			selldata := exchange.Order{}

			selldata.Rate = orderBook[i].Price
			selldata.Quantity = orderBook[i].Size / orderBook[i].Price

			maker.Asks = append(maker.Asks, selldata)
		}
	}

	return maker, nil
}

func (e *Bitmex) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *Bitmex) DoAccountOperation(operation *exchange.AccountOperation) error {
	return nil
}

func (e *Bitmex) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

}

/* Withdraw(coin *coin.Coin, quantity float64, addr, tag string) */
func (e *Bitmex) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return false
	}
	return false
}

func (e *Bitmex) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	errResponse := ErrorResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/api/v1/order"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "Sell"
	mapParams["simpleOrderQty"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)

	jsonPlaceReturn := e.ApiKeyPost(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		if err := json.Unmarshal([]byte(jsonPlaceReturn), &errResponse); err != nil {
			return nil, fmt.Errorf("%s LimitSell Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
		} else {
			return nil, fmt.Errorf("%s Place LimitSell Failed: %v %v", e.GetName(), errResponse.Error.Name, errResponse.Error.Message)
		}
	} else {
		order := &exchange.Order{
			Pair:         pair,
			Direction:    exchange.Sell,
			OrderID:      placeOrder.OrderID,
			Rate:         rate,
			Quantity:     quantity,
			Status:       exchange.New,
			JsonResponse: jsonPlaceReturn,
		}
		return order, nil
	}
}

func (e *Bitmex) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	errResponse := ErrorResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/api/v1/order"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "Buy"
	mapParams["simpleOrderQty"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)

	jsonPlaceReturn := e.ApiKeyPost(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		if err := json.Unmarshal([]byte(jsonPlaceReturn), &errResponse); err != nil {
			return nil, fmt.Errorf("%s LimitSell Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
		} else {
			return nil, fmt.Errorf("%s Place LimitSell Failed: %v %v", e.GetName(), errResponse.Error.Name, errResponse.Error.Message)
		}
	} else {
		order := &exchange.Order{
			Pair:         pair,
			Direction:    exchange.Buy,
			OrderID:      placeOrder.OrderID,
			Rate:         rate,
			Quantity:     quantity,
			Status:       exchange.New,
			JsonResponse: jsonPlaceReturn,
		}
		return order, nil
	}
}

func (e *Bitmex) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	errResponse := ErrorResponse{}
	orderStatus := []PlaceOrder{}
	strRequest := "/api/v1/order"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(order.Pair)

	jsonOrderStatus := e.ApiKeyGet(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &orderStatus); err != nil {
		if err := json.Unmarshal([]byte(jsonOrderStatus), &errResponse); err != nil {
			return fmt.Errorf("%s OrderStatus Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
		} else {
			return fmt.Errorf("%s Get OrderStatus Failed: %v %v", e.GetName(), errResponse.Error.Name, errResponse.Error.Message)
		}
	} else {
		for _, orderStatus := range orderStatus {
			if orderStatus.OrderID == order.OrderID {
				if orderStatus.OrdStatus == "Filled" {
					order.Status = exchange.Filled
				} else if orderStatus.OrdStatus == "Canceled" {
					order.Status = exchange.Cancelled
				} else if orderStatus.OrdStatus == "Partial" {
					order.Status = exchange.Partial
				} else {
					order.Status = exchange.Other
				}
				order.DealRate = orderStatus.AvgPx
				order.DealQuantity = orderStatus.SimpleOrderQty - orderStatus.SimpleLeavesQty

				return nil
			}
		}
	}

	return fmt.Errorf("%s Could not find Order: %v", e.GetName(), order.OrderID)
}

func (e *Bitmex) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Bitmex) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	return nil
}

func (e *Bitmex) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: GET and Signature is required  --reference Binance
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Bitmex) ApiKeyGet(mapParams map[string]string, strRequestPath string) string {
	strMethod := "GET"
	timestamp := time.Now().Unix() + 5

	var strRequestUrl string
	if nil == mapParams {
		strRequestUrl = strRequestPath
	} else {
		strParams := exchange.Map2UrlQuery(mapParams)
		strRequestUrl = strRequestPath + "?" + strParams
	}

	strPayload := fmt.Sprintf("%s%s%d", strMethod, strRequestUrl, timestamp)

	mapParams2Sign := make(map[string]string)
	mapParams2Sign["api-expires"] = strconv.FormatInt(timestamp, 10)
	mapParams2Sign["api-key"] = e.API_KEY
	mapParams2Sign["api-signature"] = exchange.ComputeHmac256Base64(strPayload, e.API_SECRET)

	strUrl := API_URL + strRequestUrl

	httpClient := &http.Client{}

	// 构建Request, 并且按官方要求添加Http Header
	request, err := http.NewRequest(strMethod, strUrl, nil)
	if nil != err {
		return err.Error()
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	request.Header.Set("api-expires", mapParams2Sign["api-expires"])
	request.Header.Set("api-key", mapParams2Sign["api-key"])
	request.Header.Set("api-signature", mapParams2Sign["api-signature"])

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

/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Bitmex) ApiKeyPost(mapParams map[string]string, strRequestPath string) string {
	strMethod := "POST"
	timestamp := time.Now().Unix() + 5

	jsonParams := ""
	if nil != mapParams {
		bytesParams, _ := json.Marshal(mapParams)
		jsonParams = string(bytesParams)
	}
	strUrl := API_URL + strRequestPath
	httpClient := &http.Client{}

	// 构建Request, 并且按官方要求添加Http Header
	request, err := http.NewRequest(strMethod, strUrl, strings.NewReader(jsonParams))
	if nil != err {
		return err.Error()
	}
	// log.Printf("Request: %s", request)
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("api-expires", strconv.FormatInt(timestamp, 10))
	request.Header.Add("api-key", e.API_KEY)
	strPayload := fmt.Sprintf("%s%s%d%s", strMethod, strRequestPath, timestamp, jsonParams)
	request.Header.Add("api-signature", exchange.ComputeHmac256Base64(strPayload, e.API_SECRET))

	// 发出请求
	response, err := httpClient.Do(request)
	// log.Printf("Response: %s", response)
	if nil != err {
		return err.Error()
	}
	defer response.Body.Close()

	// 解析响应内容
	body, err := ioutil.ReadAll(response.Body)
	// log.Printf("Body: %s", body)
	if nil != err {
		return err.Error()
	}

	return string(body)
}
