package mxc

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
	API_URL string = "https://www.mxc.co" //"https://www.mxc.com"
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
func (e *Mxc) GetCoinsData() error {
	jsonResponse := &JsonResponse{}
	pairsData := make(map[string]*PairsData)

	strRequestUrl := "/open/api/v1/data/markets_info"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if jsonResponse.Code != 200 {
		return fmt.Errorf("%s Get Coins Failed: %v", e.GetName(), jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &pairsData); err != nil {
		return fmt.Errorf("%s Get Coins Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for key, _ := range pairsData {
		symbols := strings.Split(key, "_")
		base := &coin.Coin{}
		target := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base = coin.GetCoin(symbols[1])
			if base == nil {
				base = &coin.Coin{}
				base.Code = symbols[1]
				coin.AddCoin(base)
			}
			target = coin.GetCoin(symbols[0])
			if target == nil {
				target = &coin.Coin{}
				target.Code = symbols[0]
				coin.AddCoin(target)
			}
		case exchange.JSON_FILE:
			base = e.GetCoinBySymbol(symbols[1])
			target = e.GetCoinBySymbol(symbols[0])
		}

		if base != nil {
			coinConstraint := e.GetCoinConstraint(base)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       base.ID,
					Coin:         base,
					ExSymbol:     symbols[1],
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = symbols[1]
			}
			e.SetCoinConstraint(coinConstraint)
		}

		if target != nil {
			coinConstraint := e.GetCoinConstraint(target)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       target.ID,
					Coin:         target,
					ExSymbol:     symbols[0],
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = symbols[0]
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
func (e *Mxc) GetPairsData() error {
	jsonResponse := &JsonResponse{}
	pairsData := make(map[string]*PairsData)

	strRequestUrl := "/open/api/v1/data/markets_info"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if jsonResponse.Code != 200 {
		return fmt.Errorf("%s Get Pairs Failed: %v", e.GetName(), jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for key, data := range pairsData {
		symbols := strings.Split(key, "_")
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(symbols[1])
			target := coin.GetCoin(symbols[0])
			if base != nil && target != nil {
				p = pair.GetPair(base, target)
			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(key)
		}
		if p != nil {
			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{
					PairID:      p.ID,
					Pair:        p,
					ExSymbol:    key,
					MakerFee:    data.SellFeeRate,
					TakerFee:    data.BuyFeeRate,
					LotSize:     math.Pow10(-1 * data.QuantityScale),
					PriceFilter: math.Pow10(-1 * data.PriceScale),
					Listed:      DEFAULT_LISTED,
				}
			} else {
				pairConstraint.ExSymbol = key
				pairConstraint.MakerFee = data.SellFeeRate
				pairConstraint.TakerFee = data.BuyFeeRate
				pairConstraint.LotSize = math.Pow10(-1 * data.QuantityScale)
				pairConstraint.PriceFilter = math.Pow10(-1 * data.PriceScale)
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
func (e *Mxc) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	jsonResponse := &JsonResponse{}
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	strRequestUrl := "/open/api/v1/data/depth"
	strUrl := API_URL + strRequestUrl

	mapParams := make(map[string]string)
	mapParams["depth"] = "150"
	mapParams["market"] = symbol

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if jsonResponse.Code != 200 {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %v", e.GetName(), jsonOrderbook)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	var err error
	for _, bid := range orderBook.Bids {
		var buydata exchange.Order

		//Modify according to type and structure
		buydata.Rate, err = strconv.ParseFloat(bid.Price, 64)
		if err != nil {
			return nil, err
		}
		buydata.Quantity, err = strconv.ParseFloat(bid.Quantity, 64)
		if err != nil {
			return nil, err
		}

		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Asks {
		var selldata exchange.Order

		//Modify according to type and structure
		selldata.Rate, err = strconv.ParseFloat(ask.Price, 64)
		if err != nil {
			return nil, err
		}
		selldata.Quantity, err = strconv.ParseFloat(ask.Quantity, 64)
		if err != nil {
			return nil, err
		}

		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

/*************** Public API ***************/

/*************** Private API ***************/
func (e *Mxc) DoAccountOperation(operation *exchange.AccountOperation) error {
	switch operation.Type {
	case exchange.BalanceList:
		return e.getAllBalance(operation)
	case exchange.Balance:
		return e.getBalance(operation)
	}
	return fmt.Errorf("%s Operation type invalid: %s %v", operation.Ex, operation.Wallet, operation.Type)
}

func (e *Mxc) getBalance(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	symbol := e.GetSymbolByCoin(operation.Coin)
	jsonResponse := &JsonResponse{}
	accountBalance := make(map[string]*AccountBalances)
	strRequest := "/open/api/v1/private/account/info"

	jsonBalanceReturn := e.ApiKeyRequest("GET", strRequest, make(map[string]string))
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s UpdateAllBalances Json Unmarshal Err: %v %s", e.GetName(), err, jsonBalanceReturn)
	}
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		return fmt.Errorf("%s UpdateAllBalances Data Unmarshal Err: %v %s", e.GetName(), err, jsonBalanceReturn)
	}

	for key, v := range accountBalance {
		if key != symbol {
			continue
		}

		freeAmount, err := strconv.ParseFloat(v.Available, 64)
		if err != nil {
			return fmt.Errorf("%s balance parse Err: %v %v", e.GetName(), err, v.Available)
		}
		frozen, err := strconv.ParseFloat(v.Frozen, 64)
		if err != nil {
			return fmt.Errorf("%s balance parse Err: %v %v", e.GetName(), err, v.Available)
		}

		operation.BalanceAvailable = freeAmount
		operation.BalanceFrozen = frozen
	}

	return nil
}

func (e *Mxc) getAllBalance(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	accountBalance := make(map[string]*AccountBalances)
	strRequest := "/open/api/v1/private/account/info"

	jsonBalanceReturn := e.ApiKeyRequest("GET", strRequest, make(map[string]string))
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
	}
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		return fmt.Errorf("%s UpdateAllBalances Data Unmarshal Err: %v %s", e.GetName(), err, jsonBalanceReturn)
	}

	for key, v := range accountBalance {
		c := e.GetCoinBySymbol(key)
		if c != nil {
			freeAmount, err := strconv.ParseFloat(v.Available, 64)
			if err != nil {
				return fmt.Errorf("%s balance parse Err: %v %v", e.GetName(), err, v.Available)
			}
			frozen, err := strconv.ParseFloat(v.Frozen, 64)
			if err != nil {
				return fmt.Errorf("%s balance parse Err: %v %v", e.GetName(), err, v.Available)
			}
			b := exchange.AssetBalance{
				Coin:             c,
				BalanceAvailable: freeAmount,
				BalanceFrozen:    frozen,
			}
			operation.BalanceList = append(operation.BalanceList, b)
		}
	}

	return nil
}

func (e *Mxc) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	accountBalance := make(map[string]*AccountBalances)
	strRequest := "/open/api/v1/private/account/info"

	jsonBalanceReturn := e.ApiKeyRequest("GET", strRequest, make(map[string]string))
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} /* else if jsonResponse.Code != 200 {
		log.Printf("%s UpdateAllBalances Failed: %v", e.GetName(), jsonResponse)
		return
	} */
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Data Unmarshal Err: %v %s", e.GetName(), err, jsonBalanceReturn)
		return
	}

	for key, v := range accountBalance {
		c := e.GetCoinBySymbol(key)
		if c != nil {
			freeAmount, err := strconv.ParseFloat(v.Available, 64)
			if err != nil {
				log.Printf("%s balance parse Err: %v %v", e.GetName(), err, v.Available)
				return
			}
			balanceMap.Set(c.Code, freeAmount)
		}
	}
}

func (e *Mxc) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {

	return false
}

func (e *Mxc) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/open/api/v1/private/order"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["market"] = e.GetSymbolByPair(pair)
	mapParams["trade_type"] = "2"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["quantity"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)

	// log.Printf("mapParams: %+v", mapParams)

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != 200 {
		return nil, fmt.Errorf("%s LimitSell Failed: %s", e.GetName(), jsonPlaceReturn)
	}
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Data Unmarshal Err: %v %s", e.GetName(), err, jsonPlaceReturn)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprint(placeOrder.OrderID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Mxc) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/open/api/v1/private/order"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["market"] = e.GetSymbolByPair(pair)
	mapParams["trade_type"] = "1"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["quantity"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)

	// log.Printf("mapParams: %+v", mapParams)

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != 200 {
		return nil, fmt.Errorf("%s LimitBuy Failed: %s", e.GetName(), jsonPlaceReturn)
	}
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Data Unmarshal Err: %v %s", e.GetName(), err, jsonPlaceReturn)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprint(placeOrder.OrderID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Mxc) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	orderStatus := OrderStatus{}
	strRequest := "/open/api/v1/private/order"

	mapParams := make(map[string]string)
	mapParams["market"] = e.GetSymbolByPair(order.Pair)
	mapParams["trade_no"] = order.OrderID

	jsonOrderStatus := e.ApiKeyRequest("GET", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if jsonResponse.Code != 200 {
		return fmt.Errorf("%s OrderStatus Failed: %v", e.GetName(), jsonResponse.Msg)
	}
	if err := json.Unmarshal([]byte(jsonOrderStatus), &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Data Unmarshal Err: %v %s", e.GetName(), err, jsonOrderStatus)
	}

	order.StatusMessage = jsonOrderStatus
	if orderStatus.ID == order.OrderID {
		switch orderStatus.Status {
		case "1":
			order.Status = exchange.New
		case "3":
			order.Status = exchange.Partial
		case "5":
			order.Status = exchange.Canceling
		case "4":
			order.Status = exchange.Cancelled
		case "2":
			order.Status = exchange.Filled
		default:
			order.Status = exchange.Other
		}
	}

	return nil
}

func (e *Mxc) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Mxc) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	strRequest := "/open/api/v1/private/order"

	mapParams := make(map[string]string)
	mapParams["market"] = e.GetSymbolByPair(order.Pair)
	mapParams["trade_no"] = order.OrderID

	jsonCancelOrder := e.ApiKeyRequest("DELETE", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if jsonResponse.Code != 200 {
		return fmt.Errorf("%s CancelOrder Failed: %v", e.GetName(), jsonResponse.Msg)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Mxc) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Mxc) ApiKeyRequest(strMethod string, strRequestPath string, mapParams map[string]string) string {
	strUrl := API_URL + strRequestPath // + "?" + exchange.Map2UrlQuery(mapParams)

	mapParams["api_key"] = e.API_KEY
	mapParams["req_time"] = fmt.Sprintf("%d", time.Now().Unix()) //"1234567890"
	authParams := exchange.Map2UrlQuery(mapParams)
	authParams += "&api_secret=" + e.API_SECRET

	signature := exchange.ComputeMD5(authParams)
	mapParams["sign"] = signature

	strUrl = strUrl + "?" + exchange.Map2UrlQuery(mapParams)

	request, err := http.NewRequest(strMethod, strUrl, nil)
	if nil != err {
		return err.Error()
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
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
