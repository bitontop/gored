package bitz

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
	API_URL string = "https://apiv2.bit-z.pro"
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
func (e *Bitz) GetCoinsData() {
	jsonResponse := &JsonResponse{}
	pairsData := make(map[string]*PairsData)

	strRequestUrl := "/Market/symbolList"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		log.Printf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if jsonResponse.Status != 200 {
		log.Printf("%s Get Coins Failed: %v %v", e.GetName(), jsonResponse.Status, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &pairsData); err != nil {
		log.Printf("%s Get Coins Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, data := range pairsData {
		base := &coin.Coin{}
		target := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base = coin.GetCoin(data.CoinTo)
			if base == nil {
				base = &coin.Coin{}
				base.Code = data.CoinTo
				coin.AddCoin(base)
			}
			target = coin.GetCoin(data.CoinFrom)
			if target == nil {
				target = &coin.Coin{}
				target.Code = data.CoinFrom
				coin.AddCoin(target)
			}
		case exchange.JSON_FILE:
			base = e.GetCoinBySymbol(data.CoinTo)
			target = e.GetCoinBySymbol(data.CoinFrom)
		}

		if base != nil {
			coinConstraint := &exchange.CoinConstraint{
				CoinID:       base.ID,
				Coin:         base,
				ExSymbol:     data.CoinTo,
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
				ExSymbol:     data.CoinFrom,
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
func (e *Bitz) GetPairsData() {
	jsonResponse := &JsonResponse{}
	pairsData := make(map[string]*PairsData)

	strRequestUrl := "/Market/symbolList"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &jsonResponse); err != nil {
		log.Printf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if jsonResponse.Status != 200 {
		log.Printf("%s Get Pairs Failed: %v %v", e.GetName(), jsonResponse.Status, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &pairsData); err != nil {
		log.Printf("%s Get Pairs Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for symbol := range pairsData {
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(pairsData[symbol].CoinTo)
			target := coin.GetCoin(pairsData[symbol].CoinFrom)
			if base != nil && target != nil {
				p = pair.GetPair(base, target)
			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(pairsData[symbol].Name)
		}
		if p != nil {
			lotSize, _ := strconv.ParseFloat(pairsData[symbol].NumberFloat, 64)
			priceFilter, _ := strconv.ParseFloat(pairsData[symbol].PriceFloat, 64)
			pairConstraint := &exchange.PairConstraint{
				PairID:      p.ID,
				Pair:        p,
				ExSymbol:    pairsData[symbol].Name,
				MakerFee:    DEFAULT_MAKERER_FEE,
				TakerFee:    DEFAULT_TAKER_FEE,
				LotSize:     lotSize,
				PriceFilter: priceFilter,
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
func (e *Bitz) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	jsonResponse := &JsonResponse{}
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	strUrl := fmt.Sprintf("%s/Market/depth?symbol=%s", API_URL, symbol)

	maker := &exchange.Maker{}
	maker.WorkerIP = exchange.GetExternalIP()
	maker.BeforeTimestamp = float64(time.Now().UnixNano() / 1e6)

	jsonOrderbook := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if jsonResponse.Status != 200 {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %v %v", e.GetName(), jsonResponse.Status, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	var err error
	for _, bid := range orderBook.Bids {
		var buydata exchange.Order

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
	for i := len(orderBook.Asks) - 1; i >= 0; i-- {
		var selldata exchange.Order

		//Modify according to type and structure
		selldata.Rate, err = strconv.ParseFloat(orderBook.Asks[i][0], 64)
		if err != nil {
			return nil, err
		}
		selldata.Quantity, err = strconv.ParseFloat(orderBook.Asks[i][1], 64)
		if err != nil {
			return nil, err
		}

		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

/*************** Private API ***************/
func (e *Bitz) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	userInfo := UserInfo{}
	accountBalance := AccountBalances{}
	strRequest := "/Assets/getUserAssets"

	jsonBalanceReturn := e.ApiKeyPOST(make(map[string]string), strRequest)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if jsonResponse.Status != 200 {
		log.Printf("%s UpdateAllBalances Failed: %v %v", e.GetName(), jsonResponse.Status, jsonResponse.Msg)
		return
	}
	if err := json.Unmarshal(jsonResponse.Data, &userInfo); err != nil {
		log.Printf("%s UpdateAllBalances Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return
	}
	if err := json.Unmarshal(userInfo.Info, &accountBalance); err != nil {
		log.Printf("%s Assets are empty: %v %v", e.GetName(), err, userInfo)
		return
	}

	for _, v := range accountBalance {
		c := e.GetCoinBySymbol(v.Name)
		if c != nil {
			balanceFloat, err := strconv.ParseFloat(v.Over, 64)
			if err != nil {
				log.Printf("Bitz balance parse to float64 error: %v", err)
				balanceFloat = 0.0
			}
			balanceMap.Set(c.Code, balanceFloat)
		}
	}
}

func (e *Bitz) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {

	return false
}

func (e *Bitz) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" || e.TradePassword == "" {
		return nil, fmt.Errorf("%s API Key, Secret Key or TradePassword are nil", e.GetName())
	}

	tradePwd := e.TradePassword

	jsonResponse := JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/Trade/addEntrustSheet"

	mapParams := make(map[string]string)
	mapParams["number"] = fmt.Sprintf("%v", quantity)
	mapParams["price"] = fmt.Sprintf("%v", rate)
	mapParams["type"] = "2"
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["tradePwd"] = tradePwd

	jsonPlaceReturn := e.ApiKeyPOST(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Status != 200 {
		return nil, fmt.Errorf("%s LimitSell Failed: %v %v", e.GetName(), jsonResponse.Status, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	orderID := fmt.Sprintf("%v", placeOrder.ID)
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

func (e *Bitz) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" || e.TradePassword == "" {
		return nil, fmt.Errorf("%s API Key, Secret Key or TradePassword are nil", e.GetName())
	}

	tradePwd := e.TradePassword

	jsonResponse := JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/Trade/addEntrustSheet"

	mapParams := make(map[string]string)
	mapParams["number"] = fmt.Sprintf("%v", quantity)
	mapParams["price"] = fmt.Sprintf("%v", rate)
	mapParams["type"] = "1"
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["tradePwd"] = tradePwd

	jsonPlaceReturn := e.ApiKeyPOST(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Status != 200 {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v %v", e.GetName(), jsonResponse.Status, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	orderID := fmt.Sprintf("%v", placeOrder.ID)
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

func (e *Bitz) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	orderStatus := PlaceOrder{}
	strRequest := "/Trade/getEntrustSheetInfo"

	mapParams := make(map[string]string)
	if order.Pair != nil {
		mapParams["entrustSheetId"] = order.OrderID
	} else {
		return fmt.Errorf("Bitz Order Status Pair cannot be null!")
	}

	jsonOrderStatus := e.ApiKeyPOST(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if jsonResponse.Status != 200 {
		return fmt.Errorf("%s OrderStatus Failed: %v %v", e.GetName(), jsonResponse.Status, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order.StatusMessage = jsonOrderStatus
	if orderStatus.Status == 2 {
		order.Status = exchange.Filled
	} else if orderStatus.Status == 1 {
		order.Status = exchange.Partial
	} else if orderStatus.Status == 0 {
		order.Status = exchange.New
	} else if orderStatus.Status == 3 {
		order.Status = exchange.Canceled
	} else {
		order.Status = exchange.Other
	}

	return nil
}

func (e *Bitz) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Bitz) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	cancelOrder := CancelOrder{}
	strRequest := "/Trade/cancelEntrustSheet"

	mapParams := make(map[string]string)
	mapParams["entrustSheetId"] = order.OrderID

	jsonCancelOrder := e.ApiKeyPOST(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if jsonResponse.Status != 200 {
		return fmt.Errorf("%s CancelOrder Failed: %v %v", e.GetName(), jsonResponse.Status, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Bitz) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Bitz) ApiKeyPOST(mapParams map[string]string, strRequestPath string) string {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	strUrl := API_URL + strRequestPath

	mapParams["apiKey"] = e.API_KEY
	mapParams["timeStamp"] = timestamp
	mapParams["nonce"] = timestamp[4:9]

	params := exchange.Map2UrlQuery(mapParams) + e.API_SECRET

	mapParams["sign"] = exchange.ComputeMD5(params)

	request, err := http.NewRequest("POST", strUrl, strings.NewReader(exchange.Map2UrlQuery(mapParams)))
	if err != nil {
		return err.Error()
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Accept", "application/json")

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
