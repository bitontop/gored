package cointiger

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
	"net/url"
	"sort"
	"strconv"
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
)

/*The Base Endpoint URL*/
const (
	API_URL = "https://api.cointiger.com/exchange/trading"
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
Step 3: Modify API Path(strRequestPath)*/
func (e *Cointiger) GetCoinsData() error {
	jsonResponse := &JsonResponse{}
	pairsData := PairsData{}

	strRequestPath := "/api/v2/currencys/v2"
	strUrl := API_URL + strRequestPath

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if jsonResponse.Msg != "suc" {
		return fmt.Errorf("%s Get Coins Failed: %s", e.GetName(), jsonCurrencyReturn)
	}
	if err := json.Unmarshal(jsonResponse.Data, &pairsData); err != nil {
		return fmt.Errorf("%s Get Coins Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, data := range pairsData {
		for _, details := range data {
			base := &coin.Coin{}
			target := &coin.Coin{}
			switch e.Source {
			case exchange.EXCHANGE_API:
				base = coin.GetCoin(details.QuoteCurrency)
				if base == nil {
					base = &coin.Coin{}
					base.Code = details.QuoteCurrency
					coin.AddCoin(base)
				}
				target = coin.GetCoin(details.BaseCurrency)
				if target == nil {
					target = &coin.Coin{}
					target.Code = details.BaseCurrency
					coin.AddCoin(target)
				}
			case exchange.JSON_FILE:
				base = e.GetCoinBySymbol(details.QuoteCurrency)
				target = e.GetCoinBySymbol(details.BaseCurrency)
			}

			// data contain target's txfee only
			if base != nil {
				if e.GetCoinConstraint(base) == nil {
					coinConstraint := e.GetCoinConstraint(base)
					if coinConstraint == nil {
						coinConstraint = &exchange.CoinConstraint{
							CoinID:       base.ID,
							Coin:         base,
							ExSymbol:     details.QuoteCurrency,
							ChainType:    exchange.MAINNET,
							TxFee:        DEFAULT_TXFEE,
							Withdraw:     DEFAULT_WITHDRAW,
							Deposit:      DEFAULT_DEPOSIT,
							Confirmation: DEFAULT_CONFIRMATION,
							Listed:       DEFAULT_LISTED,
						}
					} else {
						coinConstraint.ExSymbol = details.QuoteCurrency
					}
					e.SetCoinConstraint(coinConstraint)
				}
			}

			if target != nil {
				coinConstraint := e.GetCoinConstraint(target)
				if coinConstraint == nil {
					coinConstraint = &exchange.CoinConstraint{
						CoinID:       target.ID,
						Coin:         target,
						ExSymbol:     details.BaseCurrency,
						ChainType:    exchange.MAINNET,
						TxFee:        details.WithdrawFeeMin,
						Withdraw:     DEFAULT_WITHDRAW,
						Deposit:      DEFAULT_DEPOSIT,
						Confirmation: DEFAULT_CONFIRMATION,
						Listed:       DEFAULT_LISTED,
					}
				} else {
					coinConstraint.ExSymbol = details.BaseCurrency
				}
				e.SetCoinConstraint(coinConstraint)
			}
		}
	}
	return nil
}

/* GetPairsData - Get Pairs Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func (e *Cointiger) GetPairsData() error {
	jsonResponse := &JsonResponse{}
	pairsData := PairsData{}

	strRequestPath := "/api/v2/currencys"
	strUrl := API_URL + strRequestPath

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if jsonResponse.Msg != "suc" {
		return fmt.Errorf("%s Get Pairs Failed: %s", e.GetName(), jsonSymbolsReturn)
	}
	if err := json.Unmarshal(jsonResponse.Data, &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, pairs := range pairsData {
		for _, data := range pairs {
			p := &pair.Pair{}
			switch e.Source {
			case exchange.EXCHANGE_API:
				base := coin.GetCoin(data.QuoteCurrency)
				target := coin.GetCoin(data.BaseCurrency)
				if base != nil && target != nil {
					p = pair.GetPair(base, target)
				}
			case exchange.JSON_FILE:
				p = e.GetPairBySymbol(data.BaseCurrency + data.QuoteCurrency)
			}
			if p != nil {
				// ↓ orderbook precision, not trading precision.
				priceFilter, err := strconv.ParseFloat(data.DepthSelect.Step0, 64)
				if err != nil {
					log.Printf("%s priceFilter parse error: %v, %v", e.GetBalance, err, data.DepthSelect.Step0)
				}
				pairConstraint := e.GetPairConstraint(p)
				if pairConstraint == nil {
					pairConstraint = &exchange.PairConstraint{
						PairID:      p.ID,
						Pair:        p,
						ExSymbol:    data.BaseCurrency + data.QuoteCurrency,
						MakerFee:    DEFAULT_MAKER_FEE,
						TakerFee:    DEFAULT_TAKER_FEE,
						LotSize:     math.Pow10(-1 * data.AmountPrecision),
						PriceFilter: priceFilter, //math.Pow10(-1 * data.PricePrecision),  // orderbook precision
						Listed:      DEFAULT_LISTED,
					}
				} else {
					pairConstraint.ExSymbol = data.BaseCurrency + data.QuoteCurrency
					pairConstraint.LotSize = math.Pow10(-1 * data.AmountPrecision)
					pairConstraint.PriceFilter = priceFilter
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
Step 3: Get Exchange Pair Code ex. symbol := e.GetSymbolByPair(p)
Step 4: Modify API Path(strRequestUrl)
Step 5: Add Params - Depend on API request
Step 6: Convert the response to Standard Maker struct*/
func (e *Cointiger) OrderBook(p *pair.Pair) (*exchange.Maker, error) {
	jsonResponse := &JsonResponse{}
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(p)

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol
	mapParams["type"] = "step0"

	strRequestPath := "/api/market/depth"
	strUrl := API_URL + strRequestPath

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if jsonResponse.Code != "0" {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %s", e.GetName(), jsonOrderbook)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)

	var err error
	for _, bid := range orderBook.DepthData.Tick.Buys {
		buydata := exchange.Order{}
		buydata.Quantity, err = strconv.ParseFloat(fmt.Sprintf("%v", bid[1]), 64)
		if err != nil {
			return nil, err
		}
		buydata.Rate, err = strconv.ParseFloat(fmt.Sprintf("%v", bid[0]), 64)
		if err != nil {
			return nil, err
		}
		maker.Bids = append(maker.Bids, buydata)
	}

	for _, ask := range orderBook.DepthData.Tick.Asks {
		selldata := exchange.Order{}
		selldata.Quantity, err = strconv.ParseFloat(fmt.Sprintf("%v", ask[1]), 64)
		if err != nil {
			return nil, err
		}
		selldata.Rate, err = strconv.ParseFloat(fmt.Sprintf("%v", ask[0]), 64)
		if err != nil {
			return nil, err
		}
		maker.Asks = append(maker.Asks, selldata)
	}

	return maker, err
}

func (e *Cointiger) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *Cointiger) DoAccountOperation(operation *exchange.AccountOperation) error {
	return nil
}

func (e *Cointiger) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalances{}

	strRequestPath := "/api/user/balance"

	jsonBalanceReturn := e.ApiKeyGet(strRequestPath, make(map[string]interface{}))
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if jsonResponse.Msg != "suc" {
		log.Printf("%s UpdateAllBalances Failed: %s", e.GetName(), jsonBalanceReturn)
		return
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return
	}

	for _, balance := range accountBalance {
		c := e.GetCoinBySymbol(balance.Coin)
		floatBalance, err := strconv.ParseFloat(fmt.Sprintf("%v", balance.Normal), 64)
		if err != nil {
			return
		}
		if c != nil {
			balanceMap.Set(c.Code, floatBalance)
		}
	}
}

/* Withdraw(coin *coin.Coin, quantity float64, addr, tag string) */
func (e *Cointiger) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	log.Printf("%s Withdraw Not Viable with API.", e.GetName())
	return false
}

/*
btcusdt tradepair - price filter minimum integer(交易对价格精度最小为个位数)
*/
func (e *Cointiger) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequestPath := "/order"

	mapParams := make(map[string]interface{})
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["volume"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["side"] = "SELL"
	mapParams["type"] = "1"
	mapParams["time"] = strconv.FormatInt(time.Now().UTC().UnixNano(), 10)[:13]

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequestPath, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Msg != "suc" {
		return nil, fmt.Errorf("%s LimitSell Failed: %s", e.GetName(), jsonPlaceReturn)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      strconv.Itoa(placeOrder.OrderID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}
	return order, nil
}

/*
btcusdt tradepair - price filter minimum integer(交易对价格精度最小为个位数)
*/
func (e *Cointiger) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequestPath := "/order"

	mapParams := make(map[string]interface{})
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["volume"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["side"] = "BUY"
	mapParams["type"] = "1"
	mapParams["time"] = strconv.FormatInt(time.Now().UTC().UnixNano(), 10)[:13]

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequestPath, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Msg != "suc" {
		return nil, fmt.Errorf("%s LimitBuy Failed: %s", e.GetName(), jsonPlaceReturn)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      strconv.Itoa(placeOrder.OrderID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}
	return order, nil
}

func (e *Cointiger) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	orderStatus := OrderStatus{}
	strRequestPath := "/api/v2/order/details"

	mapParams := make(map[string]interface{})
	mapParams["symbol"] = e.GetSymbolByPair(order.Pair)
	mapParams["order_id"] = order.OrderID

	jsonOrderStatus := e.ApiKeyGet(strRequestPath, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if jsonResponse.Msg != "suc" {
		return fmt.Errorf("%s OrderStatus Failed: %s", e.GetName(), jsonOrderStatus)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	if orderStatus.Status == 4 {
		order.Status = exchange.Cancelled
	} else if orderStatus.Status == 2 {
		order.Status = exchange.Filled
	} else if orderStatus.Status == 3 {
		order.Status = exchange.Partial
	} else if orderStatus.Status == 1 {
		order.Status = exchange.New
	} else {
		order.Status = exchange.Other
	}

	order.DealRate, _ = strconv.ParseFloat(orderStatus.Price, 64)
	order.DealQuantity, _ = strconv.ParseFloat(orderStatus.DealMoney, 64)

	return nil
}

func (e *Cointiger) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Cointiger) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	cancelOrder := CancelOrder{}
	strRequestPath := "/order/batch_cancel"

	symbol := e.GetSymbolByPair(order.Pair)
	orderidlist := make(map[string][]string)
	orderidlist[symbol] = []string{order.OrderID}

	bytes, _ := json.Marshal(orderidlist)

	mapParams := make(map[string]interface{})
	mapParams["orderIdList"] = string(bytes)
	mapParams["time"] = strconv.FormatInt(time.Now().UTC().UnixNano(), 10)[:13]

	jsonCancelOrder := e.ApiKeyRequest("POST", strRequestPath, mapParams)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if jsonResponse.Msg != "suc" {
		return fmt.Errorf("%s CancelOrder Failed: %s", e.GetName(), jsonCancelOrder)
	}
	if err := json.Unmarshal(jsonResponse.Data, &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Cointiger) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Get Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/

/*
sign err sometimes due to timestamp lag
*/
func (e *Cointiger) ApiKeyGet(strRequestPath string, mapParams map[string]interface{}) string {
	mapParams["time"] = strconv.FormatInt(time.Now().UTC().UnixNano(), 10)[:13]
	payload := createPayload(mapParams)
	mapParams["sign"] = exchange.ComputeHmac512NoDecode(payload+e.API_SECRET, e.API_SECRET)
	mapParams["api_key"] = e.API_KEY

	url := exchange.Map2UrlQueryInterface(mapParams)
	strUrl := API_URL + strRequestPath + "?" + url

	request, err := http.NewRequest("GET", strUrl, nil)
	if nil != err {
		return err.Error()
	}
	request.Header.Add("Language", "en_US")
	request.Header.Add("User-Agent", "Mozilla/5.0(Macintosh;U;IntelMacOSX10_6_8;en-us)AppleWebKit/534.50(KHTML,likeGecko)Version/5.1Safari/534.50")
	request.Header.Add("Referer", "https://api.cointiger.com")

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

/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request*/

/*
sign err sometimes due to timestamp lag
*/
func (e *Cointiger) ApiKeyRequest(strMethod, strRequestPath string, mapParams map[string]interface{}) string {

	payload := createPayload(mapParams)
	Params := make(map[string]interface{})
	Params["sign"] = exchange.ComputeHmac512NoDecode(payload+e.API_SECRET, e.API_SECRET)
	Params["api_key"] = e.API_KEY
	Params["time"] = strconv.FormatInt(time.Now().UTC().UnixNano(), 10)[:13]

	urlParams := exchange.Map2UrlQueryInterface(Params)
	strUrl := API_URL + "/api/v2" + strRequestPath + "?" + urlParams

	query := exchange.Map2UrlQueryInterface(mapParams)
	values, err := url.ParseQuery(query)

	request, err := http.PostForm(strUrl, values)
	if nil != err {
		return err.Error()
	}
	request.Header.Add("Language", "en_US")
	request.Header.Add("User-Agent", "Mozilla/5.0(Macintosh;U;IntelMacOSX10_6_8;en-us)AppleWebKit/534.50(KHTML,likeGecko)Version/5.1Safari/534.50")
	request.Header.Add("Referer", "https://api.cointiger.com")

	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return ""
	}
	return string(body)
}

func createPayload(mapParams map[string]interface{}) string {
	var strParams string
	mapSort := []string{}
	for key := range mapParams {
		mapSort = append(mapSort, key)
	}
	sort.Strings(mapSort)

	for _, key := range mapSort {
		strParams += key + fmt.Sprintf("%v", mapParams[key])
	}

	return strParams
}
