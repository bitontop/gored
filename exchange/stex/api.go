package stex

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
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
				coin.AddCoin(c)
			}
		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(data.Code)
		}

		if c != nil {
			txFee, _ := strconv.ParseFloat(data.WithdrawalFeeConst, 64)
			coinConstraint := &exchange.CoinConstraint{
				CoinID:       c.ID,
				Coin:         c,
				ExSymbol:     data.Code,
				ChainType:    exchange.MAINNET,
				TxFee:        txFee,
				Withdraw:     data.Active,
				Deposit:      data.Active,
				Confirmation: DEFAULT_CONFIRMATION,
				Listed:       true,
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
		if p != nil {
			pairConstraint := &exchange.PairConstraint{
				PairID:      p.ID,
				Pair:        p,
				ExID:        fmt.Sprintf("%d", data.ID),
				ExSymbol:    data.Symbol,
				MakerFee:    DEFAULT_MAKER_FEE,
				TakerFee:    DEFAULT_TAKER_FEE,
				LotSize:     math.Pow10(data.CurrencyPrecision * -1),
				PriceFilter: math.Pow10(data.MarketPrecision * -1),
				Listed:      true,
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
	jsonResponse := JsonResponseV3{}
	orderBook := OrderBook{}

	strRequestUrl := fmt.Sprintf("/public/orderbook/%v", e.GetIDByPair(pair))
	strUrl := API3_URL + strRequestUrl

	maker := &exchange.Maker{}
	maker.WorkerIP = exchange.GetExternalIP()
	maker.BeforeTimestamp = float64(time.Now().UnixNano() / 1e6)

	jsonOrderbook := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if !jsonResponse.Success {
		return nil, fmt.Errorf("Get Orderbook Failed: %v", jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	var err error
	for _, bid := range orderBook.Bid {
		var buydata exchange.Order

		//Modify according to type and structure
		buydata.Rate, err = strconv.ParseFloat(bid.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
		}
		buydata.Quantity, err = strconv.ParseFloat(bid.Amount, 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
		}

		maker.Bids = append(maker.Bids, buydata)
	}

	for _, ask := range orderBook.Ask {
		var selldata exchange.Order

		//Modify according to type and structure
		selldata.Rate, err = strconv.ParseFloat(ask.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
		}
		selldata.Quantity, err = strconv.ParseFloat(ask.Amount, 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
		}

		maker.Asks = append(maker.Asks, selldata)
	}

	return maker, nil
}

/*************** Private API ***************/
func (e *Stex) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalances{}

	mapParams := make(map[string]string)
	mapParams["method"] = "GetInfo"

	jsonBalanceReturn := e.ApiKeyPost(mapParams)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if jsonResponse.Success != 1 {
		log.Printf("%s UpdateAllBalances Failed: %v %v", e.GetName(), jsonResponse.Error, jsonResponse.Message)
		return
	}

	if strings.Contains(jsonBalanceReturn, "\"funds\":[]") {
		EmptyAccount := EmptyAccount{}
		if err := json.Unmarshal(jsonResponse.Data, &EmptyAccount); err != nil {
			log.Printf("%s Get Balance Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
			return
		} else {
			log.Printf("%s Get Balance: This account has no balance", e.GetName())
			return
		}
	}

	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return
	}

	for coinName, _ := range accountBalance.Funds {
		c := e.GetCoinBySymbol(coinName)
		if c != nil {
			Fundf, err := strconv.ParseFloat(accountBalance.Funds[coinName], 64)
			if err != nil {
				log.Printf("Parse stex balance error: %v", err)
			}
			balanceMap.Set(c.Code, Fundf)
		}
	}
}

func (e *Stex) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil", e.GetName())
		return false
	}

	withdraw := Withdraw{}

	mapParams := make(map[string]string)
	mapParams["method"] = "Withdraw"
	mapParams["currency"] = e.GetSymbolByCoin(coin)
	mapParams["address"] = addr
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', -1, 64)

	jsonResponse := &JsonResponse{}
	jsonSubmitWithdraw := e.ApiKeyPost(mapParams)

	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &jsonResponse); err != nil {
		log.Printf("%s Withdraw Json Unmarshal Err: %v %v", e.GetName(), err, jsonSubmitWithdraw)
		return false
	} else if jsonResponse.Success != 1 {
		log.Printf("%s Withdraw Failed: %v %v", e.GetName(), jsonResponse.Error, jsonResponse.Message)
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

	placeOrder := TradeDetail{}
	jsonResponse := JsonResponse{}

	mapParams := make(map[string]string)
	mapParams["method"] = "Trade"
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["rate"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["type"] = "SELL"
	mapParams["pair"] = e.GetSymbolByPair(pair)

	jsonPlaceReturn := e.ApiKeyPost(mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Success != 1 {
		return nil, fmt.Errorf("%s LimitSell Failed: %v %v", e.GetName(), jsonResponse.Error, jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	orderID := fmt.Sprintf("%v", placeOrder.OrderID)
	order := &exchange.Order{
		Pair:         pair,
		OrderID:      orderID,
		Rate:         rate,
		Quantity:     quantity,
		Side:         "Sell",
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Stex) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	placeOrder := TradeDetail{}
	jsonResponse := JsonResponse{}

	mapParams := make(map[string]string)
	mapParams["method"] = "Trade"
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["rate"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["type"] = "BUY"
	mapParams["pair"] = e.GetSymbolByPair(pair)

	jsonPlaceReturn := e.ApiKeyPost(mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Success != 1 {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v %v", e.GetName(), jsonResponse.Error, jsonResponse.Message)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	orderID := fmt.Sprintf("%v", placeOrder.OrderID)
	order := &exchange.Order{
		Pair:         pair,
		OrderID:      orderID,
		Rate:         rate,
		Quantity:     quantity,
		Side:         "Buy",
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Stex) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	orderDetail := ActiveOrder{}

	mapParams := make(map[string]string)
	if order.Pair == nil {
		return fmt.Errorf("Stex Order Status Pair cannot be nill!")
	}

	for i := 1; i <= 4; i++ {
		mapParams["method"] = "TradeHistory"
		mapParams["pair"] = e.GetSymbolByPair(order.Pair)
		mapParams["status"] = fmt.Sprintf("%v", i)
		mapParams["from_id"] = order.OrderID
		mapParams["end_id"] = order.OrderID
		mapParams["owner"] = "ALL"

		jsonOrderStatus := e.ApiKeyPost(mapParams)
		if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
			return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
		} else if jsonResponse.Success != 1 {
			return fmt.Errorf("%s OrderStatus Failed: %v %v", e.GetName(), jsonResponse.Error, jsonResponse.Message)
		}

		if err := json.Unmarshal(jsonResponse.Data, &orderDetail); err != nil && i == 4 {
			return fmt.Errorf("%s OrderStatus order does not exist: %v %s", e.GetName(), err, jsonOrderStatus)
		} else if err == nil {
			order.Side = orderDetail[order.OrderID].Type
			var orderAmount, dealAmount float64
			if order.Side == "sell" {
				orderAmount, _ = strconv.ParseFloat(fmt.Sprintf("%v", orderDetail[order.OrderID].SellAmount), 64)
				dealAmount = orderAmount / order.Rate
			} else {
				orderAmount, _ = strconv.ParseFloat(fmt.Sprintf("%v", orderDetail[order.OrderID].BuyAmount), 64)
				dealAmount = orderAmount / order.Rate
			}

			if dealAmount == 0 {
				order.Status = exchange.New
			} else if dealAmount > 0 && dealAmount < order.Quantity {
				order.Status = exchange.Partial
			}

			if i == 4 {
				order.Status = exchange.Canceled
			} else if i == 3 {
				order.Status = exchange.Filled
			}

			order.StatusMessage = jsonOrderStatus
		}
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

	jsonResponse := &JsonResponse{}
	cancelOrder := CancelOrder{}

	mapParams := make(map[string]string)
	mapParams["order_id"] = order.OrderID
	mapParams["method"] = "CancelOrder"

	jsonCancelOrder := e.ApiKeyPost(mapParams)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if jsonResponse.Success != 1 {
		return fmt.Errorf("%s CancelOrder Failed: %v %v", e.GetName(), jsonResponse.Error, jsonResponse.Message)
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
func (e *Stex) ApiKeyPost(mapParams map[string]string) string {
	httpClient := &http.Client{}

	mapParams["nonce"] = fmt.Sprintf("%d", time.Now().UnixNano())
	payload := exchange.Map2UrlQuery(mapParams)

	request, err := http.NewRequest("POST", API_URL, strings.NewReader(payload))
	if err != nil {
		return err.Error()
	}

	sig := exchange.ComputeHmac512NoDecode(payload, e.API_SECRET)
	request.Header.Add("Key", e.API_KEY)
	request.Header.Add("Sign", sig)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Accept", "application/json")

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
