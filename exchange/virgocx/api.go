package virgocx

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"

	"strconv"
)

const (
	API_URL string = "https://www.virgocx.ca/api"
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
func (e *Virgocx) GetCoinsData() error {
	jsonResponse := &JsonResponse{}
	pairsData := PairsData{}

	strRequestUrl := "/market/tickers"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if jsonResponse.Code != 0 {
		return fmt.Errorf("%s Get Coins Failed: %v", e.GetName(), jsonCurrencyReturn)
	}
	if err := json.Unmarshal(jsonResponse.Data, &pairsData); err != nil {
		return fmt.Errorf("%s Get Coins Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, data := range pairsData {
		symbols := strings.Split(data.Symbol, "/")
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
					Listed:       true,
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
					Listed:       true,
				}
			} else {
				coinConstraint.ExSymbol = symbols[0]
			}
			e.SetCoinConstraint(coinConstraint)
		}
	}
	// need to wait 1s before using another api
	time.Sleep(time.Second * 1)
	return nil
}

/* GetPairsData - Get Pairs Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func (e *Virgocx) GetPairsData() error {
	jsonResponse := &JsonResponse{}
	pairsData := PairsData{}

	strRequestUrl := "/market/tickers"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if jsonResponse.Code != 0 {
		return fmt.Errorf("%s Get Pairs Failed: %v", e.GetName(), jsonSymbolsReturn)
	}
	if err := json.Unmarshal(jsonResponse.Data, &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, data := range pairsData {
		symbols := strings.Split(data.Symbol, "/")
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(symbols[1])
			target := coin.GetCoin(symbols[0])
			if base != nil && target != nil {
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
					MakerFee:    DEFAULT_MAKER_FEE,
					TakerFee:    DEFAULT_TAKER_FEE,
					LotSize:     DEFAULT_LOT_SIZE,
					PriceFilter: DEFAULT_PRICE_FILTER,
					Listed:      true,
				}
			} else {
				pairConstraint.ExSymbol = data.Symbol
			}
			e.SetPairConstraint(pairConstraint)
		}
	}
	time.Sleep(time.Second * 1)
	return nil
}

/*Get Pair Market Depth
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Get Exchange Pair Code ex. symbol := e.GetPairCode(p)
Step 4: Modify API Path(strRequestUrl)
Step 5: Add Params - Depend on API request
Step 6: Convert the response to Standard Maker struct*/
func (e *Virgocx) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	jsonResponse := &JsonResponse{}
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	strRequestUrl := "/market/depth"
	strUrl := API_URL + strRequestUrl

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if jsonResponse.Code != 0 {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %v", e.GetName(), jsonOrderbook)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	for _, bid := range orderBook.Bids {
		buydata := exchange.Order{}

		buydata.Quantity = bid.Qty
		buydata.Rate = bid.Price

		maker.Bids = append(maker.Bids, buydata)
	}
	for i := len(orderBook.Asks) - 1; i >= 0; i-- {
		var selldata exchange.Order

		selldata.Rate = orderBook.Asks[i].Price
		selldata.Quantity = orderBook.Asks[i].Qty

		maker.Asks = append(maker.Asks, selldata)
	}

	return maker, nil
}

/*************** Public API ***************/

/*************** Private API ***************/

func (e *Virgocx) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalances{}
	strRequest := "/member/accounts"

	mapParams := make(map[string]string)

	jsonBalanceReturn := e.ApiKeyRequest("GET", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if jsonResponse.Code != 0 {
		log.Printf("%s UpdateAllBalances Failed: %v", e.GetName(), jsonBalanceReturn)
		return
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return
	}

	for _, v := range accountBalance {
		c := e.GetCoinBySymbol(v.CoinName)
		if c != nil {
			balanceMap.Set(c.Code, v.Balance)
		}
	}
}

func (e *Virgocx) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil", e.GetName())
		return false
	}

	// mapParams := make(map[string]string)
	// mapParams["currency"] = e.GetSymbolByCoin(coin)
	// mapParams["quantity"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	// mapParams["address"] = addr

	// jsonResponse := &JsonResponse{}
	// uuid := Uuid{}
	// strRequest := "/v1.1/account/withdraw"

	// jsonSubmitWithdraw := e.ApiKeyGET(strRequest, mapParams)
	// if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &jsonResponse); err != nil {
	// 	log.Printf("%s Withdraw Json Unmarshal Err: %v %v", e.GetName(), err, jsonSubmitWithdraw)
	// 	return false
	// } else if jsonResponse.Code != 0 {
	// 	log.Printf("%s Withdraw Failed: %v", e.GetName(), jsonSubmitWithdraw)
	// 	return false
	// }
	// if err := json.Unmarshal(jsonResponse.Data, &uuid); err != nil {
	// 	log.Printf("%s Withdraw Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	// 	return false
	// }
	// return true

	return false
}

func (e *Virgocx) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["qty"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["category"] = "1" // limit
	mapParams["type"] = "2"     // buy: 1, sell: 2

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/member/addOrder"

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != 0 {
		return nil, fmt.Errorf("%s LimitSell Failed: %v", e.GetName(), jsonPlaceReturn)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%v", placeOrder.OrderID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Virgocx) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["qty"] = strconv.FormatFloat(quantity, 'f', -1, 64) // "1.2"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)   // "9000.1"
	mapParams["category"] = "1"                                   // limit
	mapParams["type"] = "1"                                       // buy: 1, sell: 2

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/member/addOrder"

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != 0 {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v", e.GetName(), jsonPlaceReturn)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	} else if placeOrder.OrderID == "0" {

	}
	// log.Printf("==%t,%s, %v", jsonResponse.Data, jsonResponse.Data, string(jsonResponse.Data) == "null")

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%v", placeOrder.OrderID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}
	return order, nil
}

// need test Status code
func (e *Virgocx) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	orderStatus := OrderStatus{}
	strRequest := "/member/queryOrder"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(order.Pair)

	jsonOrderStatus := e.ApiKeyRequest("GET", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if jsonResponse.Code != 0 {
		return fmt.Errorf("%s OrderStatus Failed: %v", e.GetName(), jsonOrderStatus)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	requiredOrder := SingleOrder{}
	found := false
	for _, tempOrder := range orderStatus {
		if fmt.Sprintf("%v", tempOrder.ID) == order.OrderID {
			requiredOrder = tempOrder
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("%s OrderStatus order not found: %v", e.GetName(), jsonOrderStatus)
	}

	// log.Printf("======requiredOrder: %+v", requiredOrder) // ================

	order.StatusMessage = jsonOrderStatus
	if requiredOrder.Status == 0 || requiredOrder.Status == 1 {
		order.Status = exchange.New
	} else if requiredOrder.Status == -1 {
		order.Status = exchange.Cancelled
	} else if requiredOrder.Status == 3 {
		order.Status = exchange.Filled
	} else if requiredOrder.Status == 2 {
		order.Status = exchange.Partial
	} else {
		order.Status = exchange.Other
	}

	return nil
}

func (e *Virgocx) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Virgocx) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	mapParams := make(map[string]string)
	mapParams["id"] = order.OrderID

	cancelOrder := CancelOrder{}
	strRequest := "/member/cancelOrder"

	jsonCancelOrder := e.ApiKeyRequest("GET", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if cancelOrder.Code != 0 || !cancelOrder.Success {
		return fmt.Errorf("%s CancelOrder Failed: %v", e.GetName(), jsonCancelOrder)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Virgocx) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Virgocx) ApiKeyGET(strRequestPath string, mapParams map[string]string) string {
	mapParams["apikey"] = e.API_KEY
	mapParams["nonce"] = fmt.Sprintf("%d", time.Now().UnixNano())

	strUrl := API_URL + strRequestPath + "?" + exchange.Map2UrlQuery(mapParams)

	signature := exchange.ComputeHmac512NoDecode(strUrl, e.API_SECRET)
	httpClient := &http.Client{}

	request, err := http.NewRequest("GET", strUrl, nil)
	if nil != err {
		return err.Error()
	}
	request.Header.Add("Content-Type", "application/json;charset=utf-8")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("apisign", signature)

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

func (e *Virgocx) ApiKeyRequest(strMethod string, strRequestPath string, mapParams map[string]string) string {
	mapParams["apiKey"] = e.API_KEY
	mapParams["apiSecret"] = e.API_SECRET

	// sort params
	keySort := []string{}
	for key := range mapParams {
		keySort = append(keySort, key)
	}
	sort.Strings(keySort)

	preSign := ""
	for _, key := range keySort {
		preSign += mapParams[key]
	}

	// delete secret from params
	delete(mapParams, "apiSecret")

	// try
	signature := exchange.ComputeMD5(preSign)
	mapParams["sign"] = signature

	var strRequestUrl string
	if nil == mapParams {
		strRequestUrl = API_URL + strRequestPath
	} else {

		strParams := exchange.Map2UrlQuery(mapParams)
		strRequestUrl = API_URL + strRequestPath + "?" + strParams
	}

	jsonParams := ""
	if nil != mapParams {
		bytesParams, _ := json.Marshal(mapParams)
		jsonParams = string(bytesParams)
	}

	// 构建Request, 并且按官方要求添加Http Header
	httpClient := &http.Client{}
	request, err := http.NewRequest(strMethod, strRequestUrl, strings.NewReader(jsonParams))
	if nil != err {
		return err.Error()
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("authorization", strings.ToUpper(exchange.ComputeMD5(signature)))
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")

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
