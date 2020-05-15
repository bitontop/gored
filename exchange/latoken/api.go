package latoken

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
	"strings"
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"

	"strconv"
)

const (
	API_URL string = "https://api.latoken.com"
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
func (e *Latoken) GetCoinsData() error {
	coinsData := CoinsData{}

	strRequestUrl := "/api/v1/ExchangeInfo/currencies"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	}

	for _, data := range coinsData {
		c := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			c = coin.GetCoin(data.Symbol)
			if c == nil {
				c = &coin.Coin{}
				c.Code = data.Symbol
				c.Name = data.Name
				coin.AddCoin(c)
			}
		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(data.Symbol)
		}

		if c != nil {
			coinConstraint := e.GetCoinConstraint(c)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       c.ID,
					Coin:         c,
					ExSymbol:     data.Symbol,
					ChainType:    exchange.MAINNET,
					TxFee:        data.Fee,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = data.Symbol
				coinConstraint.TxFee = data.Fee
			}
			// update ChainType
			switch strings.ToUpper(data.Type) {
			case "MAINNET":
				coinConstraint.ChainType = exchange.MAINNET
			case "BEP2":
				coinConstraint.ChainType = exchange.BEP2
			case "ERC20":
				coinConstraint.ChainType = exchange.ERC20
			case "NEP5":
				coinConstraint.ChainType = exchange.NEP5
			case "OMNI":
				coinConstraint.ChainType = exchange.OMNI
			case "TRC20":
				coinConstraint.ChainType = exchange.TRC20
			default:
				coinConstraint.ChainType = exchange.MAINNET
			}
			e.SetCoinConstraint(coinConstraint)
		}
	}
	// time.Sleep(time.Second * 2)
	return nil
}

/* GetPairsData - Get Pairs Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func (e *Latoken) GetPairsData() error {
	pairsData := PairsData{}

	strRequestUrl := "/api/v1/ExchangeInfo/pairs"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	}

	for _, data := range pairsData {
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(data.QuotedCurrency)
			target := coin.GetCoin(data.BaseCurrency)
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
					MakerFee:    data.MakerFee,
					TakerFee:    data.TakerFee,
					LotSize:     math.Pow10(-1 * data.AmountPrecision),
					PriceFilter: math.Pow10(-1 * data.PricePrecision),
					Listed:      true,
				}
			} else {
				pairConstraint.ExSymbol = data.Symbol
				pairConstraint.MakerFee = data.MakerFee
				pairConstraint.TakerFee = data.TakerFee
				pairConstraint.LotSize = math.Pow10(-1 * data.AmountPrecision)
				pairConstraint.PriceFilter = math.Pow10(-1 * data.PricePrecision)
			}
			e.SetPairConstraint(pairConstraint)
		}
	}
	// time.Sleep(time.Second * 2)
	return nil
}

/*Get Pair Market Depth
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Get Exchange Pair Code ex. symbol := e.GetPairCode(p)
Step 4: Modify API Path(strRequestUrl)
Step 5: Add Params - Depend on API request
Step 6: Convert the response to Standard Maker struct*/
func (e *Latoken) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	strRequestUrl := fmt.Sprintf("/api/v1/MarketData/orderBook/%v", symbol)
	strUrl := API_URL + strRequestUrl

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonOrderbook), &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	for _, bid := range orderBook.Bids {
		buydata := exchange.Order{}

		buydata.Quantity = bid.Amount
		buydata.Rate = bid.Price

		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Asks {
		selldata := exchange.Order{}

		selldata.Quantity = ask.Amount
		selldata.Rate = ask.Price

		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

/*************** Public API ***************/
func (e *Latoken) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *Latoken) DoAccountOperation(operation *exchange.AccountOperation) error {
	switch operation.Type {
	// case exchange.Transfer:
	// 	return e.transfer(operation)
	// case exchange.BalanceList:
	// 	return e.getAllBalance(operation)
	// case exchange.Balance:
	// 	return e.getBalance(operation)
	case exchange.Withdraw: // TODO, v2 key
		return e.doWithdraw(operation)
	}
	return fmt.Errorf("%s Operation type invalid: %s %v", operation.Ex, operation.Wallet, operation.Type)
}

func (e *Latoken) doWithdraw(operation *exchange.AccountOperation) error { // TODO
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
	}

	withdrawResponse := WithdrawResponse{}
	strRequest := "/v2/auth/transaction/withdraw"

	mapParams := make(map[string]string)
	mapParams["currency"] = e.GetSymbolByCoin(operation.Coin)
	mapParams["amount"] = operation.WithdrawAmount
	mapParams["destination"] = "4"
	mapParams["to_address"] = operation.WithdrawAddress

	jsonSubmitWithdraw := e.ApiKeyRequest("POST", strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonSubmitWithdraw
	}

	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &withdrawResponse); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Json Unmarshal Err: %v, %s", e.GetName(), err, jsonSubmitWithdraw)
		return operation.Error
	} /* else if !withdrawResponse.Result {
		operation.Error = fmt.Errorf("%s Withdraw Failed: %v", e.GetName(), jsonSubmitWithdraw)
		return operation.Error
	} else if withdrawResponse.WithdrawalID == "" {
		operation.Error = fmt.Errorf("%s Withdraw Failed: %v", e.GetName(), jsonSubmitWithdraw)
		return operation.Error
	} */

	operation.WithdrawID = withdrawResponse.WithdrawalID

	return nil
}

func (e *Latoken) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	accountBalance := AccountBalances{}
	// strRequest := "/api/v1/Account/balances"
	strRequest := "/v2/auth/account"
	// strRequest := "/v2/auth/transaction"

	jsonBalanceReturn := e.ApiKeyRequestV2("GET", strRequest, make(map[string]string))
	log.Printf("balance: %v", jsonBalanceReturn)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	}

	for _, v := range accountBalance {
		c := e.GetCoinBySymbol(v.Symbol)
		if c != nil {
			balanceMap.Set(c.Code, v.Available)
		}
	}
}

func (e *Latoken) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	return false
}

func (e *Latoken) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	placeOrder := PlaceOrder{}
	strRequest := "/api/v1/Order/new"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["Symbol"] = e.GetSymbolByPair(pair)
	mapParams["Side"] = "sell"
	mapParams["Amount"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)
	mapParams["Price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["OrderType"] = "limit"

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
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

func (e *Latoken) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	placeOrder := PlaceOrder{} // TestBuy{} //
	strRequest := "/api/v1/Order/new"
	// strRequest := "/api/v1/Order/test-order" // test buy api

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["Symbol"] = e.GetSymbolByPair(pair)
	mapParams["Side"] = "buy"
	mapParams["Amount"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)
	mapParams["Price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["OrderType"] = "limit"

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
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

func (e *Latoken) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	orderStatus := OrderStatus{}
	strRequest := "/api/v1/Order/get_order"

	mapParams := make(map[string]string)
	mapParams["orderId"] = order.OrderID
	// no timestamp, ?

	jsonOrderStatus := e.ApiKeyRequest("GET", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	}

	order.StatusMessage = jsonOrderStatus
	if orderStatus.OrderStatus == "partiallyFilled" {
		order.Status = exchange.Partial
	} else if orderStatus.ReaminingAmount == 0 {
		order.Status = exchange.Filled
	} else if orderStatus.ExecutedAmount == 0 {
		order.Status = exchange.New
	} else if strings.ToLower(orderStatus.OrderStatus) == "cancelled" { // need to verify this
		order.Status = exchange.Cancelled
	} else {
		order.Status = exchange.Other
	}

	return nil
}

func (e *Latoken) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Latoken) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	cancelOrder := OrderStatus{}
	strRequest := "/api/v1/Order/cancel"

	mapParams := make(map[string]string)
	mapParams["orderId"] = order.OrderID

	jsonCancelOrder := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Latoken) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Latoken) ApiKeyGET(strRequestPath string, mapParams map[string]string) string {
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

func (e *Latoken) ApiKeyRequest(strMethod, strRequestPath string, mapParams map[string]string) string {
	mapParams["timestamp"] = fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))

	signURL := strRequestPath + "?" + exchange.Map2UrlQuery(mapParams)
	strUrl := API_URL + signURL

	signature := exchange.ComputeHmac256NoDecode(signURL, e.API_SECRET)

	request, err := http.NewRequest(strMethod, strUrl, nil)
	if nil != err {
		return err.Error()
	}
	request.Header.Add("Content-Type", "application/json;charset=utf-8")
	request.Header.Add("Accept", "application/json")

	request.Header.Add("X-LA-KEY", e.API_KEY)
	request.Header.Add("X-LA-SIGNATURE", signature)
	request.Header.Add("X-LA-HASHTYPE", "HMAC-SHA256") //HMAC-SHA384, default HMAC-SHA256

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

func (e *Latoken) ApiKeyRequestV2(strMethod, strRequestPath string, mapParams map[string]string) string {
	mapParams["timestamp"] = fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond))

	signURL := strRequestPath + "?" + exchange.Map2UrlQuery(mapParams)
	strUrl := API_URL + signURL

	preSign := strMethod + strRequestPath

	signature := exchange.ComputeHmac256NoDecode(preSign, e.API_SECRET)

	request, err := http.NewRequest(strMethod, strUrl, nil)
	if nil != err {
		return err.Error()
	}
	request.Header.Add("Content-Type", "application/json;charset=utf-8")
	request.Header.Add("Accept", "application/json")

	request.Header.Add("X-LA-APIKEY", e.API_KEY)
	request.Header.Add("X-LA-SIGNATURE", signature)
	request.Header.Add("X-LA-HASHTYPE", "HMAC-SHA256") //HMAC-SHA384, default HMAC-SHA256

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
