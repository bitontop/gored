package dragonex

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
	API_URL string = "https://openapi.dragonex.io"
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
func (e *Dragonex) GetCoinsData() error {
	jsonResponse := &JsonResponse{}
	coinsData := CoinsData{}

	strRequestUrl := "/api/v1/coin/all/"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if !jsonResponse.Ok {
		return fmt.Errorf("%s Get Coins Failed: %v %v", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
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
				coin.AddCoin(c)
			}
		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(fmt.Sprintf("%v", data.CoinID))
		}

		if c != nil {
			coinConstraint := e.GetCoinConstraint(c)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       c.ID,
					Coin:         c,
					ExSymbol:     fmt.Sprintf("%v", data.CoinID),
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = fmt.Sprintf("%v", data.CoinID)
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
func (e *Dragonex) GetPairsData() error {
	jsonResponse := &JsonResponse{}
	pairsData := PairsData{}

	strRequestUrl := "/api/v1/symbol/all2/"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if !jsonResponse.Ok {
		return fmt.Errorf("%s Get Pairs Failed: %v %v", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, list := range pairsData.List {
		pairStrs := strings.Split(list[1].(string), "_")
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(pairStrs[1])
			target := coin.GetCoin(pairStrs[0])
			if base != nil && target != nil {
				p = pair.GetPair(base, target)
			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(fmt.Sprintf("%0.0f", list[0].(float64)))
		}

		if p != nil {

			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{
					PairID:      p.ID,
					Pair:        p,
					ExID:        strconv.FormatFloat(list[0].(float64), 'f', 0, 64),
					ExSymbol:    list[1].(string),
					MakerFee:    DEFAULT_MAKER_FEE,
					TakerFee:    DEFAULT_TAKER_FEE,
					LotSize:     math.Pow10(-1 * int(list[7].(float64))),
					PriceFilter: math.Pow10(-1 * int(list[5].(float64))),
					Listed:      DEFAULT_LISTED,
				}
			} else {
				pairConstraint.ExID = strconv.FormatFloat(list[0].(float64), 'f', 0, 64)
				pairConstraint.ExSymbol = list[1].(string)
				pairConstraint.LotSize = math.Pow10(-1 * int(list[7].(float64)))
				pairConstraint.PriceFilter = math.Pow10(-1 * int(list[5].(float64)))
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
func (e *Dragonex) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	jsonResponse := &JsonResponse{}
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	strRequestUrl := "/api/v1/market/depth/"
	strUrl := API_URL + strRequestUrl

	mapParams := make(map[string]string)
	mapParams["symbol_id"] = symbol

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if !jsonResponse.Ok {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %v %v", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	var err error
	for _, bid := range orderBook.Buys {
		var buydata exchange.Order

		//Modify according to type and structure
		buydata.Rate, err = strconv.ParseFloat(bid.Price, 64)
		if err != nil {
			return nil, err
		}
		buydata.Quantity, err = strconv.ParseFloat(bid.Volume, 64)
		if err != nil {
			return nil, err
		}

		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Sells {
		var selldata exchange.Order

		//Modify according to type and structure
		selldata.Rate, err = strconv.ParseFloat(ask.Price, 64)
		if err != nil {
			return nil, err
		}
		selldata.Quantity, err = strconv.ParseFloat(ask.Volume, 64)
		if err != nil {
			return nil, err
		}

		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

func (e *Dragonex) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *Dragonex) DoAccountOperation(operation *exchange.AccountOperation) error {
	return nil
}

func (e *Dragonex) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalances{}
	strRequest := "/api/v1/user/own/"

	mapParams := make(map[string]interface{})
	mapParams["access_id"] = e.API_KEY

	jsonBalanceReturn := e.ApiKeyRequest("GET", mapParams, strRequest, false)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if jsonResponse.Code != 1 {
		log.Printf("%s UpdateAllBalances Failed: %v, %s", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
		return
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return
	}

	for _, balance := range accountBalance {
		c := e.GetCoinBySymbol(fmt.Sprintf("%v", balance.CoinID))
		if c != nil {
			freeamount, err := strconv.ParseFloat(balance.Volume, 64)
			if err == nil {
				balanceMap.Set(c.Code, freeamount)
			}
		}
	}

}

func (e *Dragonex) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {

	return false
}

func (e *Dragonex) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/api/v1/order/sell/"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]interface{})
	symbolID, _ := strconv.Atoi(e.GetSymbolByPair(pair))
	mapParams["symbol_id"] = symbolID
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["volume"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequest, false)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != 1 {
		return nil, fmt.Errorf("%s LimitSell Failed: %v, %s", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
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

func (e *Dragonex) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/api/v1/order/buy/"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]interface{})
	symbolID, _ := strconv.Atoi(e.GetSymbolByPair(pair))
	mapParams["symbol_id"] = symbolID
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["volume"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequest, false)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != 1 {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v, %s", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
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

func (e *Dragonex) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	orderStatus := OrderStatus{}

	strRequest := "/api/v1/order/detail/"

	mapParams := make(map[string]interface{})

	symbolID, _ := strconv.Atoi(e.GetSymbolByPair(order.Pair))
	orderID, _ := strconv.Atoi(order.OrderID)
	mapParams["symbol_id"] = symbolID
	mapParams["order_id"] = orderID

	jsonOrderStatus := e.ApiKeyRequest("POST", mapParams, strRequest, false)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if jsonResponse.Code != 1 {
		return fmt.Errorf("%s OrderStatus Failed: %v, %s", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order.StatusMessage = jsonOrderStatus
	tradeVolume, err := strconv.ParseFloat(orderStatus.TradeVolume, 64)
	if err != nil {
		log.Printf("%s trade trade volume parse Failed: %v, %v", e.GetName(), err, orderStatus.TradeVolume)
	}
	volume, err := strconv.ParseFloat(orderStatus.Volume, 64)
	if err != nil {
		log.Printf("%s trade volume parse Failed: %v, %v", e.GetName(), err, orderStatus.Volume)
	}
	if tradeVolume == 0 {
		order.Status = exchange.New
	} else if tradeVolume < volume {
		order.Status = exchange.Partial
	} else if tradeVolume == volume {
		order.Status = exchange.Filled
	} else {
		order.Status = exchange.Other
	}

	return nil
}

func (e *Dragonex) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Dragonex) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	cancelOrder := PlaceOrder{}

	strRequest := "/api/v1/order/cancel/"

	mapParams := make(map[string]interface{})
	symbolID, _ := strconv.Atoi(e.GetSymbolByPair(order.Pair))
	mapParams["symbol_id"] = symbolID
	mapParams["order_id"] = order.OrderID

	jsonCancelOrder := e.ApiKeyRequest("POST", mapParams, strRequest, false)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if jsonResponse.Code != 1 {
		return fmt.Errorf("%s CancelOrder Failed: %v, %s", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Dragonex) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
// token limit 100 times/day, valid 24h
func (e *Dragonex) GetToken() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	token := Token{}
	strRequest := "/api/v1/token/new/"

	jsonTokenReturn := e.ApiKeyRequest("POST", nil, strRequest, true)
	if err := json.Unmarshal([]byte(jsonTokenReturn), &jsonResponse); err != nil {
		log.Printf("%s GetToken Json Unmarshal Err: %v %v", e.GetName(), err, jsonTokenReturn)
		return
	} else if !jsonResponse.Ok {
		log.Printf("%s GetToken Failed: %v, %s", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
		return
	}
	if err := json.Unmarshal(jsonResponse.Data, &token); err != nil {
		log.Printf("%s GetToken Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return
	}

	e.Token = token.Token
}

func (e *Dragonex) ApiKeyRequest(strMethod string, mapParams map[string]interface{}, strRequestPath string, getToken bool) string {

	timestamp := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	strRequestUrl := API_URL + strRequestPath

	updateToken := false
	if time.Now().UTC().YearDay() != e.LastDay || time.Now().UTC().Hour() != e.LastHour {
		updateToken = true
	}

	jsonParams := ""
	if nil != mapParams {
		bytesParams, _ := json.Marshal(mapParams)
		jsonParams = string(bytesParams)
	}

	strMessage := strMethod + "\n" + "\n" + "application/json" + "\n" + timestamp + "\n" + strRequestPath
	signature := exchange.ComputeHmac1(strMessage, e.API_SECRET)

	request, err := http.NewRequest(strMethod, strRequestUrl, strings.NewReader(jsonParams))
	if nil != err {
		return err.Error()
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("auth", e.API_KEY+":"+signature)
	if !getToken {
		if updateToken {
			e.GetToken()
			e.LastDay = time.Now().UTC().YearDay()
			e.LastHour = time.Now().UTC().Hour()
		}
		request.Header.Add("token", e.Token)
	}

	//// request.Header.Add("Content-Sha1", )
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Date", timestamp) //date
	request.Header.Add("CanonicalizedDragonExHeaders", "")

	// 发出请求
	httpClient := &http.Client{}
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
