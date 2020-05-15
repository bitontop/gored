package kraken

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
)

const (
	API_URL string = "https://api.kraken.com"
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
func (e *Kraken) GetCoinsData() error {
	jsonResponse := &JsonResponse{}
	coinsData := make(map[string]*CoinsData)

	strRequestUrl := "/0/public/Assets"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if len(jsonResponse.Error) != 0 {
		return fmt.Errorf("%s Get Coins Failed: %v", e.GetName(), jsonResponse.Error)
	}
	if err := json.Unmarshal(jsonResponse.Result, &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	for key, data := range coinsData {
		c := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			c = coin.GetCoin(data.Altname)
			if c == nil {
				c = &coin.Coin{}
				c.Code = data.Altname
				coin.AddCoin(c)
			}
		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(data.Altname)
		}

		if c != nil {
			coinConstraint := e.GetCoinConstraint(c)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       c.ID,
					Coin:         c,
					ExSymbol:     key,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = key
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
func (e *Kraken) GetPairsData() error {
	jsonResponse := &JsonResponse{}
	pairsData := make(map[string]*PairsData)

	strRequestUrl := "/0/public/AssetPairs"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if len(jsonResponse.Error) != 0 {
		return fmt.Errorf("%s Get Pairs Failed: %v", e.GetName(), jsonResponse.Error)
	}
	if err := json.Unmarshal(jsonResponse.Result, &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	for key, data := range pairsData {
		ch := strings.Split(key, ".")
		if len(ch) == 1 {
			p := &pair.Pair{}
			switch e.Source {
			case exchange.EXCHANGE_API:
				base := e.GetCoinBySymbol(data.Quote)
				target := e.GetCoinBySymbol(data.Base)
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
						ExID:        data.Wsname,
						LotSize:     math.Pow10(-1 * data.LotDecimals),
						PriceFilter: math.Pow10(-1 * data.PairDecimals),
						Listed:      DEFAULT_LISTED,
					}
				} else {
					pairConstraint.ExID = data.Wsname
					pairConstraint.ExSymbol = key
					pairConstraint.LotSize = math.Pow10(-1 * data.LotDecimals)
					pairConstraint.PriceFilter = math.Pow10(-1 * data.PairDecimals)
				}
				if len(data.FeesMaker) >= 1 {
					pairConstraint.MakerFee = data.FeesMaker[0][1]
				}
				if len(data.Fees) >= 1 {
					pairConstraint.TakerFee = data.Fees[0][1]
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
func (e *Kraken) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	jsonResponse := &JsonResponse{}
	orderBook := make(map[string]*OrderBook)
	symbol := e.GetSymbolByPair(pair)

	mapParams := make(map[string]string)
	mapParams["pair"] = symbol
	mapParams["count"] = "100"

	strRequestUrl := "/0/public/Depth"
	strUrl := API_URL + strRequestUrl

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if len(jsonResponse.Error) != 0 {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %v", e.GetName(), jsonResponse.Error)
	}
	if err := json.Unmarshal(jsonResponse.Result, &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	var err error
	for _, book := range orderBook {
		for _, bid := range book.Bids {
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
	}
	for _, book := range orderBook {
		for _, ask := range book.Asks {
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
	}
	return maker, nil
}

/*************** Private API ***************/
func (e *Kraken) DoAccountOperation(operation *exchange.AccountOperation) error {
	switch operation.Type {

	// case exchange.Transfer:
	// 	return e.transfer(operation)
	// case exchange.BalanceList:
	// 	return e.getAllBalance(operation)
	// case exchange.Balance:
	// 	return e.getBalance(operation)

	case exchange.Withdraw:
		return e.doWithdraw(operation)

	}
	return fmt.Errorf("%s Operation type invalid: %s %v", operation.Ex, operation.Wallet, operation.Type)
}

func (e *Kraken) doWithdraw(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	if operation.WithdrawTag != "" {
		operation.Error = fmt.Errorf("%s Withdraw Failed, got tag: %v, for coin: %v", e.GetName(), operation.WithdrawTag, operation.Coin.Code)
		return operation.Error
	}

	jsonResponse := &JsonResponse{}
	withdraw := WithdrawResponse{}
	strRequestPath := "/0/private/Withdraw"

	values := url.Values{
		"asset":  {e.GetSymbolByCoin(operation.Coin)},
		"key":    {operation.WithdrawAddress},
		"amount": {operation.WithdrawAmount},
	}

	jsonSubmitWithdraw := e.ApiKeyPost(strRequestPath, values, &WithdrawResponse{})
	if operation.DebugMode {
		operation.RequestURI = strRequestPath
		operation.CallResponce = jsonSubmitWithdraw
	}

	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Json Unmarshal Err: %v, %s", e.GetName(), err, jsonSubmitWithdraw)
		return operation.Error
	} else if len(jsonResponse.Error) != 0 {
		operation.Error = fmt.Errorf("%s Withdraw Failed: %v", e.GetName(), jsonSubmitWithdraw)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Result, &withdraw); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
		return operation.Error
	}

	operation.WithdrawID = withdraw.RefID

	return nil
}

func (e *Kraken) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	accountBalance := make(map[string]string)
	strRequest := "/0/private/Balance"

	jsonBalanceReturn := e.ApiKeyPost(strRequest, url.Values{}, make(map[string]string))
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if len(jsonResponse.Error) != 0 {
		log.Printf("%s UpdateAllBalances Failed: %v", e.GetName(), jsonResponse.Error)
		return
	}
	if err := json.Unmarshal(jsonResponse.Result, &accountBalance); err != nil && jsonResponse.Result != nil {
		log.Printf("%s UpdateAllBalances Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse)
		return
	}

	for symb, balance := range accountBalance {
		c := e.GetCoinBySymbol(symb)
		bal, _ := strconv.ParseFloat(balance, 64)
		if c != nil {
			balanceMap.Set(c.Code, bal)
		}
	}
}

func (e *Kraken) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return false
	}

	jsonResponse := &JsonResponse{}
	withdraw := WithdrawResponse{}
	strRequestPath := "/0/private/Withdraw"

	values := url.Values{
		"asset":  {e.GetSymbolByCoin(coin)},
		"key":    {addr},
		"amount": {strconv.FormatFloat(quantity, 'f', -1, 64)},
	}

	jsonSubmitWithdraw := e.ApiKeyPost(strRequestPath, values, &WithdrawResponse{})
	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &jsonResponse); err != nil {
		log.Printf("%s Withdraw Json Unmarshal Err: %v %v", e.GetName(), err, jsonSubmitWithdraw)
		return false
	} else if len(jsonResponse.Error) != 0 {
		log.Printf("%s Withdraw Failed: %v", e.GetName(), jsonResponse.Error)
		return false
	}
	if err := json.Unmarshal(jsonResponse.Result, &withdraw); err != nil {
		log.Printf("%s Withdraw Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
		return false
	}

	return true
}

func (e *Kraken) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequestPath := "/0/private/AddOrder"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	params := url.Values{
		"pair":      {e.GetSymbolByPair(pair)},
		"type":      {"sell"},
		"ordertype": {"limit"},
		"price":     {strconv.FormatFloat(rate, 'f', priceFilter, 64)},
		"volume":    {strconv.FormatFloat(quantity, 'f', lotSize, 64)},
	}

	jsonPlaceReturn := e.ApiKeyPost(strRequestPath, params, &PlaceOrder{})
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if len(jsonResponse.Error) != 0 {
		return nil, fmt.Errorf("%s LimitSell Failed: %v", e.GetName(), jsonResponse.Error)
	}
	if err := json.Unmarshal(jsonResponse.Result, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      strings.Join(placeOrder.TransactionIds, ""),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}
	return order, nil
}

func (e *Kraken) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequestPath := "/0/private/AddOrder"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	params := url.Values{
		"pair":      {e.GetSymbolByPair(pair)},
		"type":      {"buy"},
		"ordertype": {"limit"},
		"price":     {strconv.FormatFloat(rate, 'f', priceFilter, 64)},
		"volume":    {strconv.FormatFloat(quantity, 'f', lotSize, 64)},
	}

	jsonPlaceReturn := e.ApiKeyPost(strRequestPath, params, &PlaceOrder{})
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if len(jsonResponse.Error) != 0 {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v", e.GetName(), jsonResponse.Error)
	}
	if err := json.Unmarshal(jsonResponse.Result, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      strings.Join(placeOrder.TransactionIds, ""),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}
	return order, nil
}

func (e *Kraken) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	orderStatus := OrderStatus{}
	strRequestPath := "/0/private/QueryOrders"

	params := url.Values{"txid": {order.OrderID}}

	jsonOrderStatus := e.ApiKeyPost(strRequestPath, params, &OrderStatus{})
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if len(jsonResponse.Error) != 0 {
		return fmt.Errorf("%s OrderStatus Failed: %v", e.GetName(), jsonResponse.Error)
	}
	if err := json.Unmarshal(jsonResponse.Result, &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}
	for _, orders := range orderStatus {
		vol, _ := strconv.ParseFloat(orders.Volume, 64)
		if orders.Status == "canceled" {
			order.Status = exchange.Cancelled
		} else if vol == orders.VolumeExecuted {
			order.Status = exchange.Filled
		} else if orders.VolumeExecuted != 0 && orders.Status == "open" {
			order.Status = exchange.Partial
		} else if orders.Status == "open" && orders.VolumeExecuted == 0 {
			order.Status = exchange.New
		} else {
			order.Status = exchange.Other
		}

		order.DealRate = orders.Cost / orders.VolumeExecuted
		order.DealQuantity = orders.VolumeExecuted
	}
	return nil
}

func (e *Kraken) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Kraken) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	cancelOrder := CancelOrder{}
	strRequestPath := "/0/private/CancelOrder"

	params := url.Values{"txid": {order.OrderID}}

	jsonCancelOrder := e.ApiKeyPost(strRequestPath, params, &CancelOrder{})
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if len(jsonResponse.Error) != 0 {
		return fmt.Errorf("%s CancelOrder Failed: %v", e.GetName(), jsonResponse.Error)
	}
	if err := json.Unmarshal(jsonResponse.Result, &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Kraken) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Kraken) ApiKeyPost(strRequestPath string, values url.Values, typ interface{}) string {
	/* if e.Two_Factor != "" {
		mapParams["otp"] = e.Two_Factor
	} */
	strUrl := API_URL + strRequestPath
	httpClient := &http.Client{}
	non := fmt.Sprintf("%d", time.Now().UnixNano())
	values.Set("nonce", non)
	secret, _ := base64.StdEncoding.DecodeString(e.API_SECRET)
	signature := createSignature(strRequestPath, values, secret)

	/* jsonParams := ""
	if nil != mapParams {
		bytesParams, _ := json.Marshal(mapParams)
		jsonParams = string(bytesParams)
	}
	jsonParams = exchange.Map2UrlQuery(mapParams) */

	headers := map[string]string{
		"API-Key":  e.API_KEY,
		"API-Sign": signature,
	}

	request, err := http.NewRequest("POST", strUrl, strings.NewReader(values.Encode()))
	if nil != err {
		return err.Error()
	}

	request.Header.Add("User-Agent", "Kraken GO API Agent (https://github.com/beldur/kraken-go-api-client)")
	for key, value := range headers {
		request.Header.Add(key, value)
	}

	response, err := httpClient.Do(request)
	if nil != err {
		return err.Error()
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return err.Error()
	}

	// Check mime type of response
	mimeType, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil {
		return err.Error()
	}
	if mimeType != "application/json" {
		return err.Error()
	}

	return string(body)
}

//Signature加密
func getSha256(input []byte) []byte {
	sha := sha256.New()
	sha.Write(input)
	return sha.Sum(nil)
}

// getHMacSha512 creates a hmac hash with sha512
func getHMacSha512(message, secret []byte) []byte {
	mac := hmac.New(sha512.New, secret)
	mac.Write(message)
	return mac.Sum(nil)
}

func createSignature(urlPath string, values url.Values, secret []byte) string {
	// See https://www.kraken.com/help/api#general-usage for more information
	shaSum := getSha256([]byte(values.Get("nonce") + values.Encode()))
	macSum := getHMacSha512(append([]byte(urlPath), shaSum...), secret)
	return base64.StdEncoding.EncodeToString(macSum)
}
