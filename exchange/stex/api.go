package stex

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
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
)

const (
	API_URL  string = "https://app.stex.com/api2"
	API3_URL string = "https://api3.stex.com"
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
func (e *Stex) GetCoinsData() error {
	jsonResponse := &JsonResponseV3{}
	coinsData := CoinsData{}

	strRequestUrl := "/public/currencies"
	strUrl := API3_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if !jsonResponse.Success {
		return fmt.Errorf("%s Get Coins Failed: %v", e.GetName(), jsonResponse.Message)
	}

	if err := json.Unmarshal(jsonResponse.Data, &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, data := range coinsData {
		c := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			c = coin.GetCoin(data.Code)
			if c == nil {
				c = &coin.Coin{}
				c.Code = data.Code
				c.Name = data.Name
				c.Explorer = data.BlockExplorerURL
				coin.AddCoin(c)
			}
		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(data.Code)
		}

		if c != nil {
			txFee, _ := strconv.ParseFloat(data.WithdrawalFeeConst, 64)
			coinConstraint := e.GetCoinConstraint(c)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       c.ID,
					Coin:         c,
					ExSymbol:     fmt.Sprintf("%d", data.ID),
					ChainType:    exchange.MAINNET,
					TxFee:        txFee,
					Withdraw:     data.Active,
					Deposit:      data.Active,
					Confirmation: data.MinimumTxConfirmations,
					Listed:       !(data.Delisted),
				}
			} else {
				coinConstraint.ExSymbol = fmt.Sprintf("%d", data.ID)
				coinConstraint.TxFee = txFee
				coinConstraint.Withdraw = data.Active
				coinConstraint.Deposit = data.Active
				coinConstraint.Confirmation = data.MinimumTxConfirmations
				coinConstraint.Listed = !data.Delisted
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
func (e *Stex) GetPairsData() error {
	jsonResponse := JsonResponseV3{}
	pairsData := PairsData{}

	strRequestUrl := "/public/currency_pairs/list/ALL"
	strUrl := API3_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if !jsonResponse.Success {
		return fmt.Errorf("%s Get Coins Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Data, &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, data := range pairsData {
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(data.MarketCode)
			target := coin.GetCoin(data.CurrencyCode)
			if base != nil && target != nil {
				p = pair.GetPair(base, target)
			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(data.Symbol)
		}
		if p != nil && !data.Delisted {
			makerFee, _ := strconv.ParseFloat(data.BuyFeePercent, 64)
			takerFee, _ := strconv.ParseFloat(data.SellFeePercent, 64)
			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{
					PairID:      p.ID,
					Pair:        p,
					ExID:        fmt.Sprintf("%d", data.ID),
					ExSymbol:    data.Symbol,
					MakerFee:    makerFee / 100,
					TakerFee:    takerFee / 100,
					LotSize:     math.Pow10(data.CurrencyPrecision * -1),
					PriceFilter: math.Pow10(data.MarketPrecision * -1),
					Listed:      !data.Delisted,
				}
			} else {
				pairConstraint.ExID = fmt.Sprintf("%d", data.ID)
				pairConstraint.ExSymbol = data.Symbol
				pairConstraint.MakerFee = makerFee / 100
				pairConstraint.TakerFee = takerFee / 100
				pairConstraint.LotSize = math.Pow10(data.CurrencyPrecision * -1)
				pairConstraint.PriceFilter = math.Pow10(data.MarketPrecision * -1)
				pairConstraint.Listed = !data.Delisted
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
func (e *Stex) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	// return e.webpageOrderBook(pair) // using webpage orderbook
	jsonResponse := JsonResponseV3{}
	orderBook := OrderBook{}

	strRequestUrl := fmt.Sprintf("/public/orderbook/%s", e.GetIDByPair(pair))
	strUrl := API3_URL + strRequestUrl

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		log.Printf("%s Get Orderbook Err: %v, %s, Using webpage Orderbook...", e.GetName(), err, jsonOrderbook)
		return e.webpageOrderBook(pair) //nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if !jsonResponse.Success {
		log.Printf("%s Get Orderbook fail: %s, Using webpage Orderbook...", e.GetName(), jsonOrderbook)
		return e.webpageOrderBook(pair) //nil, fmt.Errorf("Get Orderbook Failed: %v", jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderBook); err != nil {
		log.Printf("%s Get Orderbook Err: %v, %s, Using webpage Orderbook...", e.GetName(), err, jsonOrderbook)
		return e.webpageOrderBook(pair) //nil, fmt.Errorf("%s Get Orderbook Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	var err error
	for _, bid := range orderBook.Bid {
		var buydata exchange.Order

		//Modify according to type and structure
		buydata.Rate, err = strconv.ParseFloat(bid.Price, 64)
		if err != nil {
			return e.webpageOrderBook(pair) //nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Rate error:%v", e.GetName(), err)
		}
		buydata.Quantity, err = strconv.ParseFloat(bid.Amount, 64)
		if err != nil {
			return e.webpageOrderBook(pair) //nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
		}

		maker.Bids = append(maker.Bids, buydata)
	}

	for _, ask := range orderBook.Ask {
		var selldata exchange.Order

		//Modify according to type and structure
		selldata.Rate, err = strconv.ParseFloat(ask.Price, 64)
		if err != nil {
			return e.webpageOrderBook(pair) //nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Rate error:%v", e.GetName(), err)
		}
		selldata.Quantity, err = strconv.ParseFloat(ask.Amount, 64)
		if err != nil {
			return e.webpageOrderBook(pair) //nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
		}

		maker.Asks = append(maker.Asks, selldata)
	}

	// // test two orderbook get the same value
	// webMaker, _ := e.webpageOrderBook(pair)
	// for i, _ := range webMaker.Bids {
	// 	if maker.Bids[i].Rate != webMaker.Bids[i].Rate {
	// 		log.Printf("Rate not match here: %v->%v", maker.Bids[i].Rate, webMaker.Bids[i].Rate)
	// 	}
	// 	if maker.Bids[i].Quantity != webMaker.Bids[i].Quantity {
	// 		log.Printf("Quantity not match here: %v->%v", maker.Bids[i].Quantity, webMaker.Bids[i].Quantity)
	// 	}
	// }
	// for i, _ := range webMaker.Asks {
	// 	if maker.Asks[i].Rate != webMaker.Asks[i].Rate {
	// 		log.Printf("Rate not match here: %v->%v", maker.Asks[i].Rate, webMaker.Asks[i].Rate)
	// 	}
	// 	if maker.Asks[i].Quantity != webMaker.Asks[i].Quantity {
	// 		log.Printf("Quantity not match here: %v->%v", maker.Asks[i].Quantity, webMaker.Asks[i].Quantity)
	// 	}
	// }

	return maker, nil
}

// OrderBook from webpage
func (e *Stex) webpageOrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	orderBookBuy := WebOrderBook{}
	orderBookSell := WebOrderBook{}

	// strRequestUrl := fmt.Sprintf("/public/orderbook/%v", e.GetIDByPair(pair))
	strUrlBuy := fmt.Sprintf("https://app.stex.com/en/basic-trade/buy-glass/%v", e.GetIDByPair(pair))
	strUrlSell := fmt.Sprintf("https://app.stex.com/en/basic-trade/sell-glass/%v", e.GetIDByPair(pair))

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbookBuy := exchange.HttpGetRequest(strUrlBuy, nil)
	if err := json.Unmarshal([]byte(jsonOrderbookBuy), &orderBookBuy); err != nil {
		return nil, fmt.Errorf("%s Get WebOrderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbookBuy)
	} else if len(orderBookBuy) == 0 {
		return nil, fmt.Errorf("Got empty WebOrderbook: %v", jsonOrderbookBuy)
	}

	jsonOrderbookSell := exchange.HttpGetRequest(strUrlSell, nil)
	if err := json.Unmarshal([]byte(jsonOrderbookSell), &orderBookSell); err != nil {
		return nil, fmt.Errorf("%s Get WebOrderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbookSell)
	} else if len(orderBookSell) == 0 {
		return nil, fmt.Errorf("Got empty WebOrderbook: %v", jsonOrderbookSell)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	var err error
	for _, bid := range orderBookBuy {
		var buydata exchange.Order

		buydata.Rate, err = strconv.ParseFloat(bid.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("%s WebOrderbook strconv.ParseFloat Rate error:%v", e.GetName(), err)
		}
		buydata.Quantity, err = strconv.ParseFloat(bid.Amount, 64)
		if err != nil {
			return nil, fmt.Errorf("%s WebOrderbook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
		}

		maker.Bids = append(maker.Bids, buydata)
	}

	for _, ask := range orderBookSell {
		var selldata exchange.Order

		selldata.Rate, err = strconv.ParseFloat(ask.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("%s WebOrderbook strconv.ParseFloat Rate error:%v", e.GetName(), err)
		}
		selldata.Quantity, err = strconv.ParseFloat(ask.Amount, 64)
		if err != nil {
			return nil, fmt.Errorf("%s WebOrderbook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
		}

		maker.Asks = append(maker.Asks, selldata)
	}

	return maker, nil
}

/*************** Private API ***************/

func (e *Stex) doWithdraw(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	if operation.WithdrawTag != "" {
		operation.Error = fmt.Errorf("%s Withdraw Failed, got tag: %v, for coin: %v", e.GetName(), operation.WithdrawTag, operation.Coin.Code)
		return operation.Error
	}

	jsonResponse := JsonResponseV3{}
	withdraw := WithdrawResult{}
	strRequestUrl := "/profile/withdraw"

	mapParams := make(map[string]interface{})
	mapParams["currency_id "] = e.GetSymbolByCoin(operation.Coin)
	mapParams["address"] = operation.WithdrawAddress // "0x37E0Fc27C6cDB5035B2a3d0682B4E7C05A4e6C46"
	mapParams["amount"] = operation.WithdrawAmount

	// log.Printf("mapParams: %+v", mapParams)

	jsonCreateWithdraw := e.ApiKeyRequest("POST", mapParams, strRequestUrl)
	if operation.DebugMode {
		operation.RequestURI = strRequestUrl
		operation.CallResponce = jsonCreateWithdraw
	}

	if err := json.Unmarshal([]byte(jsonCreateWithdraw), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Json Unmarshal Err: %v, %s", e.GetName(), err, jsonCreateWithdraw)
		return operation.Error
	} else if !jsonResponse.Success {
		operation.Error = fmt.Errorf("%s Withdraw Failed: %v", e.GetName(), jsonCreateWithdraw)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &withdraw); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	operation.WithdrawID = fmt.Sprintf("%v", withdraw.ID)

	return nil
}

func (e *Stex) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponseV3{}
	accountBalance := WalletDetails{}

	strRequestUrl := "/profile/wallets"

	jsonBalanceReturn := e.ApiKeyGet(make(map[string]string), strRequestUrl)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %s", e.GetName(), err, jsonBalanceReturn)
		return
	} else if !jsonResponse.Success {
		log.Printf("%s UpdateAllBalances Failed: %v %s", e.GetName(), jsonResponse.Message, jsonBalanceReturn)
		return
	}

	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return
	}

	for _, data := range accountBalance {
		c := e.GetCoinBySymbol(fmt.Sprintf("%d", data.CurrencyID))
		if c != nil {
			balance, err := strconv.ParseFloat(data.Balance, 64)
			if err != nil {
				log.Printf("Parse stex balance error: %v", err)
			}
			balanceMap.Set(c.Code, balance)
		}
	}
}

func (e *Stex) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil", e.GetName())
		return false
	}

	jsonResponse := JsonResponseV3{}
	withdraw := WithdrawResult{}
	strRequestUrl := "/profile/withdraw"

	mapParams := make(map[string]interface{})
	mapParams["currency_id "] = e.GetSymbolByCoin(coin)
	mapParams["address"] = addr
	mapParams["amount"] = quantity

	jsonCreateWithdraw := e.ApiKeyRequest("POST", mapParams, strRequestUrl)
	if err := json.Unmarshal([]byte(jsonCreateWithdraw), &jsonResponse); err != nil {
		log.Printf("%s Withdraw Json Unmarshal Err: %v %s", e.GetName(), err, jsonCreateWithdraw)
		return false
	} else if !jsonResponse.Success {
		log.Printf("%s Withdraw Failed: %v %v", e.GetName(), jsonResponse.Message, jsonCreateWithdraw)
		return false
	}
	if err := json.Unmarshal(jsonResponse.Data, &withdraw); err != nil {
		log.Printf("%s Withdraw Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return false
	}

	return true
}

func (e *Stex) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	placeOrder := PlaceOrder{}
	jsonResponse := JsonResponseV3{}

	strRequestUrl := fmt.Sprintf("/trading/orders/%s", e.GetIDByPair(pair))

	mapParams := make(map[string]interface{})
	mapParams["type"] = "SELL"
	mapParams["amount"] = quantity
	mapParams["price"] = rate

	jsonPlaceOrder := e.ApiKeyRequest("POST", mapParams, strRequestUrl)
	if err := json.Unmarshal([]byte(jsonPlaceOrder), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %s", e.GetName(), err, jsonPlaceOrder)
	} else if !jsonResponse.Success {
		return nil, fmt.Errorf("%s LimitSell Failed: %v %s", e.GetName(), jsonResponse.Message, jsonPlaceOrder)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%d", placeOrder.ID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceOrder,
	}

	return order, nil
}

func (e *Stex) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	placeOrder := PlaceOrder{}
	jsonResponse := JsonResponseV3{}

	strRequestUrl := fmt.Sprintf("/trading/orders/%s", e.GetIDByPair(pair))

	mapParams := make(map[string]interface{})
	mapParams["type"] = "BUY"
	mapParams["amount"] = quantity
	mapParams["price"] = rate

	jsonPlaceOrder := e.ApiKeyRequest("POST", mapParams, strRequestUrl)
	if err := json.Unmarshal([]byte(jsonPlaceOrder), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %s", e.GetName(), err, jsonPlaceOrder)
	} else if !jsonResponse.Success {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v %s", e.GetName(), jsonResponse.Message, jsonPlaceOrder)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%d", placeOrder.ID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceOrder,
	}

	return order, nil
}

func (e *Stex) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponseV3{}
	placeOrder := PlaceOrder{}

	strRequestUrl := fmt.Sprintf("/trading/order/%s", order.OrderID)

	jsonOrderStatus := e.ApiKeyGet(nil, strRequestUrl)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Order Status Json Unmarshal Err: %v %s", e.GetName(), err, jsonOrderStatus)
	} else if !jsonResponse.Success {
		return fmt.Errorf("%s Get Order Status Failed: %v %s", e.GetName(), jsonResponse.Message, jsonOrderStatus)
	}

	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return fmt.Errorf("%s Get Order Status Unmarshal Error: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	if placeOrder.Status == "PROCESSING" {
		order.Status = exchange.Other
	} else if placeOrder.Status == "PENDING" {
		order.Status = exchange.Cancelled
	} else if placeOrder.Status == "FINISHED" {
		order.Status = exchange.Filled
	} else if placeOrder.Status == "PARTIAL" {
		order.Status = exchange.Partial
	} else if placeOrder.Status == "CANCELLED" {
		order.Status = exchange.Cancelled
	}

	return nil
}

func (e *Stex) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Stex) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponseV3{}
	cancelOrder := CancelOrder{}

	strRequestUrl := fmt.Sprintf("/trading/order/%s", order.OrderID)

	jsonCancelOrder := e.ApiKeyRequest("DELETE", nil, strRequestUrl)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %s", e.GetName(), err, jsonCancelOrder)
	} else if !jsonResponse.Success {
		return fmt.Errorf("%s CancelOrder Failed: %v %s", e.GetName(), jsonResponse.Message, jsonCancelOrder)
	}
	if err := json.Unmarshal(jsonResponse.Data, &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Stex) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Stex) ApiKeyGet(mapParams map[string]string, strRequestPath string) string {
	if mapParams != nil {
		strRequestPath = fmt.Sprintf("%s?%s", strRequestPath, exchange.Map2UrlQuery(mapParams))
	}

	strUrl := API3_URL + strRequestPath

	request, err := http.NewRequest("GET", strUrl, nil)
	if nil != err {
		return err.Error()
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json; charset=UTF-8")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", e.API_SECRET))

	// log.Printf("request: %+v", request)

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if nil != err {
		return err.Error()
	}
	defer response.Body.Close()

	// log.Printf("response: %+v", response)

	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return err.Error()
	}

	return string(body)
}

func (e *Stex) ApiKeyRequest(method string, mapParams map[string]interface{}, strRequestPath string) string {
	httpClient := &http.Client{}

	var bytesParams []byte
	if mapParams != nil {
		bytesParams, _ = json.Marshal(mapParams)
	}

	strUrl := API3_URL + strRequestPath

	request, err := http.NewRequest(method, strUrl, bytes.NewReader(bytesParams))
	if nil != err {
		return err.Error()
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json; charset=UTF-8")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", e.API_SECRET))

	// 发出请求
	response, err := httpClient.Do(request)
	if err != nil {
		log.Printf("Stex Request error: %v", err)
		return err.Error()
	}
	defer response.Body.Close()

	// 解析响应内容
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err.Error()
	}

	return string(body)
}
