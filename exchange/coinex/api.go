package coinex

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"../../coin"
	"../../exchange"
	"../../pair"
)

/*The Base Endpoint URL*/
const (
	API_URL = "https://api.coinegg.im"
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

func getCoinFromCoinMarket() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/info", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := url.Values{}
	q.Add("start", "1")
	q.Add("limit", "5000")
	q.Add("convert", "USD")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", "b54bcf4d-1bca-4e8e-9a24-22ff2c3d462c")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		os.Exit(1)
	}
	fmt.Println(resp.Status)
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBody))
	return string(respBody)
}

/*************** Public API ***************/
/*Get Coins Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func (e *Coinex) GetCoinsData() {
	coinsData := CoinsData{}

	strUrl := "" //"https://www.binance.com/assetWithdraw/getAllAsset.html"

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &coinsData); err != nil {
		log.Printf("%s Get Coins Data Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
		return
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
				c.Website = data.URL
				c.Explorer = data.BlockURL
				coin.AddCoin(c)
			}
		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(data.AssetCode)
		}

		if c != nil {
			confirmation, _ := strconv.Atoi(data.ConfirmTimes)
			coinConstraint := &exchange.CoinConstraint{
				CoinID:       c.ID,
				Coin:         c,
				ExSymbol:     data.AssetCode,
				TxFee:        data.TransactionFee,
				Withdraw:     data.EnableWithdraw,
				Deposit:      data.EnableCharge,
				Confirmation: confirmation,
				Listed:       true,
			}

			e.SetCoinConstraint(coinConstraint)
		}
	}
}

/* GetPairsData - Get Pairs Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func (e *Coinex) GetPairsData() {
	pairsData := &PairsData{}

	strRequestUrl := "/api/v1/exchangeInfo"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &pairsData); err != nil {
		log.Printf("%s Get Pairs Data Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
		return
	}

	for _, data := range pairsData.Symbols {
		if data.Status == "TRADING" {
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
			if p != nil {
				var err error
				lotsize := 0.0
				priceFilter := 0.0
				for _, filter := range data.Filters {
					switch filter.FilterType {
					case "LOT_SIZE":
						lotsize, err = strconv.ParseFloat(filter.StepSize, 64)
						if err != nil {
							log.Printf("%s Lot Size Err: %v", e.GetName(), err)
							lotsize = DEFAULT_LOT_SIZE
						}
					case "PRICE_FILTER":
						priceFilter, err = strconv.ParseFloat(filter.TickSize, 64)
						if err != nil {
							log.Printf("%s Price Filter Err: %v", e.GetName(), err)
							priceFilter = DEFAULT_PRICE_FILTER
						}
					}
				}
				pairConstraint := &exchange.PairConstraint{
					PairID:      p.ID,
					Pair:        p,
					ExSymbol:    data.Symbol,
					MakerFee:    DEFAULT_MAKERER_FEE,
					TakerFee:    DEFAULT_TAKER_FEE,
					LotSize:     lotsize,
					PriceFilter: priceFilter,
					Listed:      true,
				}
				e.SetPairConstraint(pairConstraint)
			}
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
func (e *Coinex) OrderBook(p *pair.Pair) (*exchange.Maker, error) {
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(p)

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol
	mapParams["limit"] = "100"

	strRequestUrl := "/api/v1/depth"
	strUrl := API_URL + strRequestUrl

	maker := &exchange.Maker{}
	maker.WorkerIP = exchange.GetExternalIP()
	maker.BeforeTimestamp = float64(time.Now().UnixNano() / 1e6)

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &orderBook); err != nil {
		log.Printf("%s OrderBook json Unmarshal error: %v %v", e.GetName(), err, jsonOrderbook)
		return nil, err
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	maker.LastUpdateID = orderBook.LastUpdateID

	var err error
	for _, bid := range orderBook.Bids {
		buydata := exchange.Order{}
		buydata.Quantity, err = strconv.ParseFloat(bid[1].(string), 64)
		if err != nil {
			log.Printf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
			return nil, err
		}

		buydata.Rate, err = strconv.ParseFloat(bid[0].(string), 64)
		if err != nil {
			log.Printf("%s OrderBook strconv.ParseFloat Rate error:%v", e.GetName(), err)
			return nil, err
		}
		maker.Bids = append(maker.Bids, buydata)
	}

	for _, ask := range orderBook.Asks {
		selldata := exchange.Order{}
		selldata.Quantity, err = strconv.ParseFloat(ask[1].(string), 64)
		if err != nil {
			log.Printf("%s OrderBook strconv.ParseFloat  Quantity error:%v", e.GetName(), err)
			return nil, err
		}

		selldata.Rate, err = strconv.ParseFloat(ask[0].(string), 64)
		if err != nil {
			log.Printf("%s OrderBook strconv.ParseFloat  Rate error:%v", e.GetName(), err)
			return nil, err
		}
		maker.Asks = append(maker.Asks, selldata)
	}

	return maker, err
}

/*************** Private API ***************/
func (e *Coinex) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	accountBalance := AccountBalances{}
	strRequest := "/api/v3/account"

	jsonBalanceReturn := e.ApiKeyRequest("GET", make(map[string]string), strRequest)
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
func (e *Coinex) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
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
	mapParams["amount"] = fmt.Sprintf("%f", quantity)
	mapParams["timestamp"] = fmt.Sprintf("%d", time.Now().UnixNano()/1e6)

	jsonSubmitWithdraw := e.ApiKeyRequest("POST", mapParams, strRequest)
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

func (e *Coinex) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	placeOrder := PlaceOrder{}
	strRequest := "/api/v3/order"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "SELL"
	mapParams["type"] = "LIMIT"
	mapParams["timeInForce"] = "GTC"
	mapParams["price"] = fmt.Sprintf("%f", rate)
	mapParams["quantity"] = fmt.Sprintf("%f", quantity)

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
		Side:         "Sell",
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}
	return order, nil
}

func (e *Coinex) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	placeOrder := PlaceOrder{}
	strRequest := "/api/v3/order"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "BUY"
	mapParams["type"] = "LIMIT"
	mapParams["timeInForce"] = "GTC"
	mapParams["price"] = fmt.Sprintf("%f", rate)
	mapParams["quantity"] = fmt.Sprintf("%f", quantity)

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
		Side:         "Buy",
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Coinex) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	orderStatus := PlaceOrder{}
	strRequest := "/api/v3/order"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(order.Pair)
	mapParams["orderId"] = order.OrderID

	jsonOrderStatus := e.ApiKeyRequest("GET", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if orderStatus.Code != 0 {
		return fmt.Errorf("%s Get OrderStatus Error: %v %s", e.GetName(), orderStatus.Code, orderStatus.Msg)
	}

	if orderStatus.Status == "CANCELED" {
		order.Status = exchange.Canceled
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

func (e *Coinex) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Coinex) CancelOrder(order *exchange.Order) error {
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

func (e *Coinex) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Coinex) ApiKeyRequest(strMethod string, mapParams map[string]string, strRequestPath string) string {
	mapParams["recvWindow"] = fmt.Sprintf("50000000")
	mapParams["timestamp"] = fmt.Sprintf("%d", time.Now().UTC().UnixNano()/int64(time.Millisecond))

	payload := exchange.Map2UrlQuery(mapParams)
	mapParams["signature"] = exchange.ComputeHmac256(payload, e.API_SECRET)
	strUrl := API_URL + strRequestPath + "?" + exchange.Map2UrlQuery(mapParams)

	httpClient := &http.Client{}

	request, err := http.NewRequest(http.MethodGet, strUrl, nil)
	if nil != err {
		return err.Error()
	}
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	request.Header.Add("X-MBX-APIKEY", e.API_KEY)

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
