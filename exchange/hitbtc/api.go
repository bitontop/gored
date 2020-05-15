package hitbtc

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
)

const (
	API_URL string = "https://api.hitbtc.com"
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
func (e *Hitbtc) GetCoinsData() error {
	coinsData := CoinsData{}

	strRequestUrl := "/api/2/public/currency"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	}

	for _, data := range coinsData {
		c := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			c = coin.GetCoin(data.ID)
			if c == nil {
				c = &coin.Coin{}
				c.Code = data.ID
				c.Name = data.FullName
				coin.AddCoin(c)
			}
		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(data.ID)
		}

		if c != nil {
			txFee, _ := strconv.ParseFloat(data.PayoutFee, 64)
			coinConstraint := e.GetCoinConstraint(c)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       c.ID,
					Coin:         c,
					ExSymbol:     data.ID,
					ChainType:    exchange.MAINNET,
					TxFee:        txFee,
					Withdraw:     data.PayoutEnabled,
					Deposit:      data.PayinEnabled,
					Confirmation: data.PayinConfirmations,
					Listed:       !data.Delisted,
				}
			} else {
				coinConstraint.ExSymbol = data.ID
				coinConstraint.TxFee = txFee
				coinConstraint.Withdraw = data.PayoutEnabled
				coinConstraint.Deposit = data.PayinEnabled
				coinConstraint.Confirmation = data.PayinConfirmations
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
func (e *Hitbtc) GetPairsData() error {
	pairsData := PairsData{}

	strRequestUrl := "/api/2/public/symbol"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	}

	for _, data := range pairsData {
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(data.QuoteCurrency)
			target := coin.GetCoin(data.BaseCurrency)
			if base != nil && target != nil {
				p = pair.GetPair(base, target)
			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(data.ID)
		}
		if p != nil {
			lotSize, err := strconv.ParseFloat(data.QuantityIncrement, 64)
			if err != nil {
				log.Printf("%s Lot Size Err: %s", e.GetName(), err)
			}
			priceFilter, err := strconv.ParseFloat(data.TickSize, 64)
			if err != nil {
				log.Printf("%s Price Filter Err: %s", e.GetName(), err)
			}
			makerFee, err := strconv.ParseFloat(data.ProvideLiquidityRate, 64)
			if err != nil {
				log.Printf("%s Maker Fee Err: %s", e.GetName(), err)
			}
			takerFee, err := strconv.ParseFloat(data.TakeLiquidityRate, 64)
			if err != nil {
				log.Printf("%s Taker Fee Err: %s", e.GetName(), err)
			}
			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{
					PairID:      p.ID,
					Pair:        p,
					ExSymbol:    data.ID,
					MakerFee:    makerFee,
					TakerFee:    takerFee,
					LotSize:     lotSize,
					PriceFilter: priceFilter,
					Listed:      DEFAULT_LISTED,
				}
			} else {
				pairConstraint.ExSymbol = data.ID
				pairConstraint.MakerFee = makerFee
				pairConstraint.TakerFee = takerFee
				pairConstraint.LotSize = lotSize
				pairConstraint.PriceFilter = priceFilter
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
func (e *Hitbtc) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	strRequestUrl := fmt.Sprintf("/api/2/public/orderbook/%s", symbol)
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
	var err error
	for _, bid := range orderBook.Bid {
		var buydata exchange.Order

		buydata.Rate, err = strconv.ParseFloat(bid.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("%s Bids Rate ParseFloat error: %v", e.GetName(), err)
		}
		buydata.Quantity, err = strconv.ParseFloat(bid.Size, 64)
		if err != nil {
			return nil, fmt.Errorf("%s Bids Quantity ParseFloat error: %v", e.GetName(), err)
		}

		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Ask {
		var selldata exchange.Order

		selldata.Rate, err = strconv.ParseFloat(ask.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("%s Asks Rate ParseFloat error: %v", e.GetName(), err)
		}
		selldata.Quantity, err = strconv.ParseFloat(ask.Size, 64)
		if err != nil {
			return nil, fmt.Errorf("%s Asks Quantity ParseFloat error: %v", e.GetName(), err)
		}

		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

/*************** Private API ***************/
func (e *Hitbtc) DoAccountOperation(operation *exchange.AccountOperation) error {
	switch operation.Type {
	case exchange.BalanceList:
		return e.getAllBalance(operation)
	case exchange.Balance:
		return e.getBalance(operation)
	case exchange.Withdraw:
		return e.doWithdraw(operation)
	}
	return fmt.Errorf("%s Operation type invalid: %s %v", operation.Ex, operation.Wallet, operation.Type)
}

func (e *Hitbtc) getAllBalance(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	accountBalance := AccountBalances{}
	errResponse := ErrResponse{}
	strRequest := "/api/2/trading/balance"

	jsonAllBalanceReturn := e.ApiKeyRequest("GET", make(map[string]string), strRequest)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonAllBalanceReturn
	}

	json.Unmarshal([]byte(jsonAllBalanceReturn), &errResponse)
	// log.Printf("=========ALLBALANCE ERR: %v", errResponse) // ==============
	if err := json.Unmarshal([]byte(jsonAllBalanceReturn), &accountBalance); err != nil {
		operation.Error = fmt.Errorf("%s getAllBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonAllBalanceReturn)
		return operation.Error
	} else if errResponse.Error.Code != 0 {
		operation.Error = fmt.Errorf("%s getAllBalance failed: %v", e.GetName(), jsonAllBalanceReturn)
		return operation.Error
	}
	// log.Printf("=========ALLBALANCE: %v", jsonAllBalanceReturn) // ==============

	for _, v := range accountBalance {
		available, err := strconv.ParseFloat(v.Available, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s getAllBalance parse balance Err: %v, %s", e.GetName(), err, accountBalance)
			return operation.Error
		}
		frozen, err := strconv.ParseFloat(v.Available, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s getAllBalance parse balance Err: %v, %s", e.GetName(), err, accountBalance)
			return operation.Error
		}

		if available == 0 && frozen == 0 {
			continue
		}

		b := exchange.AssetBalance{
			Coin:             e.GetCoinBySymbol(v.Currency),
			BalanceAvailable: available,
			BalanceFrozen:    frozen,
		}
		operation.BalanceList = append(operation.BalanceList, b)
	}

	return nil
}

func (e *Hitbtc) getBalance(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	accountBalance := AccountBalances{}
	errResponse := ErrResponse{}
	strRequest := "/api/2/trading/balance"
	symbol := e.GetSymbolByCoin(operation.Coin)

	jsonBalanceReturn := e.ApiKeyRequest("GET", make(map[string]string), strRequest)
	// log.Printf("RETURN : %v", jsonBalanceReturn)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonBalanceReturn
	}

	json.Unmarshal([]byte(jsonBalanceReturn), &errResponse)
	// log.Printf("=========BALANCE ERR: %v, %v", errResponse, errResponse.Error.Code) // ==============
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		operation.Error = fmt.Errorf("%s getBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonBalanceReturn)
		return operation.Error
	} else if errResponse.Error.Code != 0 {
		operation.Error = fmt.Errorf("%s getBalance failed: %v", e.GetName(), jsonBalanceReturn)
		return operation.Error
	}

	for _, v := range accountBalance {
		if v.Currency != symbol {
			// log.Printf("v.Currency,symbol: [%v-%v]", v.Currency, symbol)
			continue
		}

		available, err := strconv.ParseFloat(v.Available, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s getBalance parse balance Err: %v, %s", e.GetName(), err, accountBalance)
			return operation.Error
		}
		frozen, err := strconv.ParseFloat(v.Available, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s getBalance parse balance Err: %v, %s", e.GetName(), err, accountBalance)
			return operation.Error
		}

		// if available == 0 && frozen == 0 {
		// 	continue
		// }

		operation.BalanceFrozen = frozen
		operation.BalanceAvailable = available
		return nil
	}

	operation.Error = fmt.Errorf("%s getBalance get %v account balance fail: %v", e.GetName(), symbol, jsonBalanceReturn)
	return operation.Error
}

func (e *Hitbtc) doWithdraw(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	if operation.WithdrawTag != "" {
		operation.Error = fmt.Errorf("%s Withdraw Failed, got tag: %v, for coin: %v", e.GetName(), operation.WithdrawTag, operation.Coin.Code)
		return operation.Error
	}

	withdraw := Withdraw{}
	errResponse := ErrResponse{}
	strRequest := "/api/2/account/crypto/withdraw"

	mapParams := make(map[string]string)
	mapParams["currency"] = e.GetSymbolByCoin(operation.Coin)
	mapParams["amount"] = operation.WithdrawAmount
	mapParams["address"] = operation.WithdrawAddress

	jsonWithdrawReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonWithdrawReturn
	}

	json.Unmarshal([]byte(jsonWithdrawReturn), &errResponse)
	if err := json.Unmarshal([]byte(jsonWithdrawReturn), &withdraw); err != nil {
		operation.Error = fmt.Errorf("%s withdraw Json Unmarshal Err: %v, %s", e.GetName(), err, jsonWithdrawReturn)
		return operation.Error
	} else if errResponse.Error.Code != 0 {
		operation.Error = fmt.Errorf("%s withdraw failed: %v", e.GetName(), jsonWithdrawReturn)
		return operation.Error
	} else if withdraw.ID == "" {
		operation.Error = fmt.Errorf("%s withdraw failed: %v", e.GetName(), jsonWithdrawReturn)
		return operation.Error
	}

	operation.WithdrawID = withdraw.ID

	return nil
}

func (e *Hitbtc) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	accountBalance := AccountBalances{}
	errResponse := ErrResponse{}
	strRequest := "/api/2/trading/balance"

	jsonBalanceReturn := e.ApiKeyRequest("GET", make(map[string]string), strRequest)
	json.Unmarshal([]byte(jsonBalanceReturn), &errResponse)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if errResponse.Error.Code != 0 {
		log.Printf("%s UpdateAllBalances Failed: %v %v", e.GetName(), errResponse.Error.Code, errResponse.Error.Message)
		return
	}

	for _, v := range accountBalance {
		freeamount, err := strconv.ParseFloat(v.Available, 64)
		if err == nil {
			c := e.GetCoinBySymbol(v.Currency)
			if c != nil {
				balanceMap.Set(c.Code, freeamount)
			}
		} else {
			log.Printf("%s %s Get Balance Err: %s\n", e.GetName(), v.Currency, err)
		}
	}
}

func (e *Hitbtc) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	return false
}

func (e *Hitbtc) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	placeOrder := PlaceOrder{}
	errResponse := ErrResponse{}
	strRequest := "/api/2/order"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "sell"
	mapParams["type"] = "limit"
	mapParams["quantity"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	json.Unmarshal([]byte(jsonPlaceReturn), &errResponse)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if errResponse.Error.Code != 0 {
		return nil, fmt.Errorf("%s LimitSell Failed: %v", e.GetName(), errResponse)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%d", placeOrder.ID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Hitbtc) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	placeOrder := PlaceOrder{}
	errResponse := ErrResponse{}
	strRequest := "/api/2/order"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "buy"
	mapParams["type"] = "limit"
	mapParams["quantity"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	json.Unmarshal([]byte(jsonPlaceReturn), &errResponse)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if errResponse.Error.Code != 0 {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v", e.GetName(), errResponse)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%d", placeOrder.ID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Hitbtc) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	errResponse := &ErrResponse{}
	orderStatus := PlaceOrder{}
	strRequest := fmt.Sprintf("/api/2/order/%s", order.OrderID)

	jsonOrderStatus := e.ApiKeyRequest("GET", make(map[string]string), strRequest)
	json.Unmarshal([]byte(jsonOrderStatus), &errResponse)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if errResponse.Error.Code != 0 {
		return fmt.Errorf("%s OrderStatus Failed: %v", e.GetName(), errResponse)
	}

	order.StatusMessage = jsonOrderStatus
	var err error
	if err == nil && errResponse.Error.Code == 0 {
		if orderStatus.Status == "Filled" {
			order.Status = exchange.Filled
		} else if orderStatus.Status == "partiallyFilled" {
			order.Status = exchange.Partial
		} else if orderStatus.Status == "canceled" {
			order.Status = exchange.Cancelled
		} else if orderStatus.Status == "expired" {
			order.Status = exchange.Expired
		} else if orderStatus.Status == "suspended" {
			order.Status = exchange.Other
		}

		order.DealRate = order.Rate
		order.DealQuantity, _ = strconv.ParseFloat(orderStatus.CumQuantity, 64)
	}

	return nil
}

func (e *Hitbtc) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Hitbtc) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	errResponse := &ErrResponse{}
	cancelOrder := []*PlaceOrder{}
	strRequest := "/api/2/order"

	jsonCancelOrder := e.ApiKeyRequest("DELETE", make(map[string]string), strRequest)
	json.Unmarshal([]byte(jsonCancelOrder), &errResponse)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if errResponse.Error.Code != 0 {
		return fmt.Errorf("%s CancelOrder Failed: %v", e.GetName(), errResponse)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Hitbtc) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Hitbtc) ApiKeyRequest(strMethod string, mapParams map[string]string, strRequestPath string) string {
	strUrl := API_URL + strRequestPath

	var bytesParams []byte
	if mapParams != nil {
		bytesParams, _ = json.Marshal(mapParams)
	}

	request, err := http.NewRequest(strMethod, strUrl, bytes.NewReader(bytesParams))
	if nil != err {
		return err.Error()
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	request.Header.Add("Accept", "application/json")
	request.SetBasicAuth(e.API_KEY, e.API_SECRET)

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
