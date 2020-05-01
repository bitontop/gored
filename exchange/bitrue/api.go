package bitrue

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
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
)

const (
	API_URL string = "https://www.bitrue.com"
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
func (e *Bitrue) GetCoinsData() error {
	pairsData := PairsData{}
	jsonResponse := JsonResponse{}

	strRequestUrl := "/api/v1/exchangeInfo"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if jsonResponse.Code != 0 {
		return fmt.Errorf("%s Get Coins Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &pairsData); err != nil {
		return fmt.Errorf("%s Get Coins Result Unmarshal Err: %v %s", e.GetName(), err, jsonCurrencyReturn)
	}

	for _, data := range pairsData.Symbols {
		base := &coin.Coin{}
		target := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base = coin.GetCoin(data.QuoteAsset)
			if base == nil {
				base = &coin.Coin{}
				base.Code = data.QuoteAsset
				coin.AddCoin(base)
			}
			target = coin.GetCoin(data.BaseAsset)
			if target == nil {
				target = &coin.Coin{}
				target.Code = data.BaseAsset
				coin.AddCoin(target)
			}
		case exchange.JSON_FILE:
			base = e.GetCoinBySymbol(data.QuoteAsset)
			target = e.GetCoinBySymbol(data.BaseAsset)
		}

		if base != nil {
			coinConstraint := e.GetCoinConstraint(base)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       base.ID,
					Coin:         base,
					ExSymbol:     data.QuoteAsset,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = data.QuoteAsset
			}
			e.SetCoinConstraint(coinConstraint)
		}

		if target != nil {
			coinConstraint := e.GetCoinConstraint(target)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       target.ID,
					Coin:         target,
					ExSymbol:     data.BaseAsset,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = data.BaseAsset
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
func (e *Bitrue) GetPairsData() error {
	pairsData := PairsData{}
	jsonResponse := JsonResponse{}

	strRequestUrl := "/api/v1/exchangeInfo"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if jsonResponse.Code != 0 {
		return fmt.Errorf("%s Get Pairs Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Result Unmarshal Err: %v %s", e.GetName(), err, jsonSymbolsReturn)
	}

	for _, data := range pairsData.Symbols {
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
			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{
					PairID:      p.ID,
					Pair:        p,
					ExSymbol:    data.Symbol,
					MakerFee:    DEFAULT_MAKER_FEE,
					TakerFee:    DEFAULT_TAKER_FEE,
					LotSize:     math.Pow10(-1 * data.BaseAssetPrecision),
					PriceFilter: math.Pow10(-1 * data.QuotePrecision),
					Listed:      DEFAULT_LISTED,
				}
			} else {
				pairConstraint.ExSymbol = data.Symbol
				pairConstraint.LotSize = math.Pow10(-1 * data.BaseAssetPrecision)
				pairConstraint.PriceFilter = math.Pow10(-1 * data.QuotePrecision)
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
func (e *Bitrue) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)
	jsonResponse := JsonResponse{}

	strRequestUrl := "/api/v1/depth"
	strUrl := API_URL + strRequestUrl

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol
	mapParams["limit"] = "0"

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if jsonResponse.Code != 0 {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal([]byte(jsonOrderbook), &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Data Unmarshal Err: %v %s", e.GetName(), err, jsonOrderbook)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	var err error
	for _, bid := range orderBook.Bids {
		buydata := exchange.Order{}
		buydata.Quantity, err = strconv.ParseFloat(bid[1].(string), 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
		}

		buydata.Rate, err = strconv.ParseFloat(bid[0].(string), 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Rate error:%v", e.GetName(), err)
		}
		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Asks {
		selldata := exchange.Order{}
		selldata.Quantity, err = strconv.ParseFloat(ask[1].(string), 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
		}

		selldata.Rate, err = strconv.ParseFloat(ask[0].(string), 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Rate error:%v", e.GetName(), err)
		}
		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

func (e *Bitrue) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *Bitrue) DoAccountOperation(operation *exchange.AccountOperation) error {
	return nil
}

func (e *Bitrue) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalances{}
	strRequest := "/api/v1/account"

	jsonBalanceReturn := e.ApiKeyRequest("GET", make(map[string]string), strRequest)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if jsonResponse.Code != 0 {
		log.Printf("%s UpdateAllBalances Failed: %v", e.GetName(), jsonResponse.Message)
		return
	}
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Data Unmarshal Err: %v %s", e.GetName(), err, jsonBalanceReturn)
		return
	}

	for _, v := range accountBalance.Balances {
		c := e.GetCoinBySymbol(v.Asset)
		if c != nil {
			freeAmount, err := strconv.ParseFloat(v.Free, 64)
			if err != nil {
				log.Printf("%s balance parse Err: %v %v", e.GetName(), err, v.Free)
				return
			}
			balanceMap.Set(c.Code, freeAmount)
		}
	}
}

func (e *Bitrue) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {

	return false
}

func (e *Bitrue) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/api/v1/order"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "SELL"
	mapParams["type"] = "LIMIT"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["quantity"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != 0 {
		return nil, fmt.Errorf("%s LimitSell Failed: %v", e.GetName(), jsonResponse.Message)
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

func (e *Bitrue) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/api/v1/order"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "BUY"
	mapParams["type"] = "LIMIT"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["quantity"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != 0 {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v", e.GetName(), jsonResponse.Message)
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

func (e *Bitrue) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	orderStatus := OrderStatus{}
	strRequest := "/api/v1/order"

	mapParams := make(map[string]string)
	mapParams["method"] = "GET"
	mapParams["symbol"] = e.GetSymbolByPair(order.Pair)
	mapParams["orderId"] = order.OrderID

	jsonOrderStatus := e.ApiKeyRequest("GET", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if jsonResponse.Code != 0 {
		return fmt.Errorf("%s OrderStatus Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal([]byte(jsonOrderStatus), &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Data Unmarshal Err: %v %s", e.GetName(), err, jsonOrderStatus)
	}

	order.StatusMessage = jsonOrderStatus
	if orderStatus.OrderID == order.OrderID {
		switch orderStatus.Status {
		case "NEW":
			order.Status = exchange.New
		case "PARTIALLY_FILLED":
			order.Status = exchange.Partial
		case "REJECTED":
			order.Status = exchange.Rejected
		case "PENDING_CANCEL":
			order.Status = exchange.Canceling
		case "CANCELED":
			order.Status = exchange.Cancelled
		case "FILLED":
			order.Status = exchange.Filled
		case "EXPIRED":
			order.Status = exchange.Expired
		default:
			order.Status = exchange.Other
		}
	}

	return nil
}

func (e *Bitrue) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Bitrue) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	cancelOrder := PlaceOrder{}
	strRequest := "/api/v1/order"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(order.Pair)

	jsonCancelOrder := e.ApiKeyRequest("DELETE", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if jsonResponse.Code != 0 {
		return fmt.Errorf("%s CancelOrder Failed: %v", e.GetName(), jsonResponse.Message)
	}
	if err := json.Unmarshal([]byte(jsonCancelOrder), &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Data Unmarshal Err: %v %s", e.GetName(), err, jsonCancelOrder)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Bitrue) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Bitrue) ApiKeyRequest(strMethod string, mapParams map[string]string, strRequestPath string) string {
	mapParams["timestamp"] = fmt.Sprintf("%.0d", time.Now().UnixNano()/1e6)

	strUrl := API_URL + strRequestPath

	var strParams string
	if nil != mapParams {
		strParams = exchange.Map2UrlQuery(mapParams)
	}

	signature := exchange.ComputeHmac256NoDecode(strParams, e.API_SECRET)
	signMessage := strUrl + "?" + strParams + "&signature=" + signature

	request, err := http.NewRequest(strMethod, signMessage, nil)
	if nil != err {
		return err.Error()
	}

	request.Header.Add("X-MBX-APIKEY", e.API_KEY)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

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

func (e *Bitrue) ApiKeyGET(strRequestPath string, mapParams map[string]string) string {
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
