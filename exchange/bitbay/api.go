package bitbay

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"

	"strconv"
)

const (
	API_URL string = "https://api.bitbay.net/rest"
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
func (e *Bitbay) GetCoinsData() error {
	pairsData := PairsData{}

	strRequestUrl := "/trading/ticker"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &pairsData); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if pairsData.Status != "Ok" {
		return fmt.Errorf("%s Get Coins Failed: %v", e.GetName(), jsonCurrencyReturn)
	}

	for _, data := range pairsData.Pairs {
		base := &coin.Coin{}
		target := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base = coin.GetCoin(data.Market.Second.Currency)
			if base == nil {
				base = &coin.Coin{}
				base.Code = data.Market.Second.Currency
				coin.AddCoin(base)
			}
			target = coin.GetCoin(data.Market.First.Currency)
			if target == nil {
				target = &coin.Coin{}
				target.Code = data.Market.First.Currency
				coin.AddCoin(target)
			}
		case exchange.JSON_FILE:
			base = e.GetCoinBySymbol(data.Market.Second.Currency)
			target = e.GetCoinBySymbol(data.Market.First.Currency)
		}

		if base != nil {
			coinConstraint := e.GetCoinConstraint(base)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       base.ID,
					Coin:         base,
					ExSymbol:     data.Market.Second.Currency,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = data.Market.Second.Currency
			}
			e.SetCoinConstraint(coinConstraint)
		}

		if target != nil {
			coinConstraint := e.GetCoinConstraint(target)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       target.ID,
					Coin:         target,
					ExSymbol:     data.Market.First.Currency,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = data.Market.First.Currency
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
func (e *Bitbay) GetPairsData() error {
	pairsData := PairsData{}

	strRequestUrl := "/trading/ticker"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if pairsData.Status != "Ok" {
		return fmt.Errorf("%s Get Pairs Failed: %v", e.GetName(), jsonSymbolsReturn)
	}

	for _, data := range pairsData.Pairs {
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(data.Market.Second.Currency)
			target := coin.GetCoin(data.Market.First.Currency)
			if base != nil && target != nil {

				p = pair.GetPair(base, target)

			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(data.Market.Code)
		}
		if p != nil {
			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{
					PairID:      p.ID,
					Pair:        p,
					ExSymbol:    data.Market.Code,
					MakerFee:    DEFAULT_MAKER_FEE,
					TakerFee:    DEFAULT_TAKER_FEE,
					LotSize:     math.Pow10(-1 * data.Market.First.Scale),
					PriceFilter: math.Pow10(-1 * data.Market.Second.Scale),
					Listed:      DEFAULT_LISTED,
				}
			} else {
				pairConstraint.ExSymbol = data.Market.Code
				pairConstraint.LotSize = math.Pow10(-1 * data.Market.First.Scale)
				pairConstraint.PriceFilter = math.Pow10(-1 * data.Market.Second.Scale)
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
func (e *Bitbay) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	strRequestUrl := fmt.Sprintf("/trading/orderbook/%s", symbol)
	strUrl := API_URL + strRequestUrl

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonOrderbook), &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if orderBook.Status != "Ok" {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %v", e.GetName(), jsonOrderbook)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	var err error
	for _, bid := range orderBook.Buy {
		buydata := exchange.Order{}
		buydata.Quantity, err = strconv.ParseFloat(bid.Ca, 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
		}

		buydata.Rate, err = strconv.ParseFloat(bid.Ra, 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Rate error:%v", e.GetName(), err)
		}
		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Sell {
		selldata := exchange.Order{}
		selldata.Quantity, err = strconv.ParseFloat(ask.Ca, 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
		}

		selldata.Rate, err = strconv.ParseFloat(ask.Ra, 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Rate error:%v", e.GetName(), err)
		}
		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

func (e *Bitbay) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *Bitbay) DoAccountOperation(operation *exchange.AccountOperation) error {
	return nil
}

func (e *Bitbay) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	accountBalance := AccountBalances{}
	strRequest := "/balances/BITBAY/balance"

	jsonBalanceReturn := e.ApiKeyGET(strRequest, make(map[string]interface{}))
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if accountBalance.Status != "Ok" {
		log.Printf("%s UpdateAllBalances Failed: %v, %s", e.GetName(), jsonBalanceReturn, accountBalance.Errors)
		return
	}

	for _, v := range accountBalance.Balances {
		c := e.GetCoinBySymbol(v.Currency)
		if c != nil {
			balanceMap.Set(c.Code, v.AvailableFunds)
		}
	}
}

func (e *Bitbay) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {

	return false
}

func (e *Bitbay) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	placeOrder := PlaceOrder{}
	strRequest := fmt.Sprintf("/trading/offer/", e.GetSymbolByPair(pair))

	mapParams := make(map[string]interface{})
	price := float64(int(rate/e.GetPriceFilter(pair)+e.GetPriceFilter(pair)/10)) * (e.GetPriceFilter(pair))
	amount := float64(int(quantity/e.GetLotSize(pair)+e.GetLotSize(pair)/10)) * (e.GetLotSize(pair))
	mapParams["rate"] = price
	mapParams["amount"] = amount
	mapParams["offerType"] = "sell"
	mapParams["mode"] = "limit"

	jsonPlaceReturn := e.ApiKeyGET(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if placeOrder.Status != "Ok" {
		return nil, fmt.Errorf("%s LimitSell Failed: %v", e.GetName(), jsonPlaceReturn)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      placeOrder.OfferID,
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Bitbay) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	placeOrder := PlaceOrder{}
	strRequest := fmt.Sprintf("/trading/offer/", e.GetSymbolByPair(pair))

	mapParams := make(map[string]interface{})
	price := float64(int(rate/e.GetPriceFilter(pair)+e.GetPriceFilter(pair)/10)) * (e.GetPriceFilter(pair))
	amount := float64(int(quantity/e.GetLotSize(pair)+e.GetLotSize(pair)/10)) * (e.GetLotSize(pair))
	mapParams["rate"] = price
	mapParams["amount"] = amount
	mapParams["offerType"] = "buy"
	mapParams["mode"] = "limit"

	jsonPlaceReturn := e.ApiKeyGET(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if placeOrder.Status != "Ok" {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v", e.GetName(), jsonPlaceReturn)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      placeOrder.OfferID,
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Bitbay) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	orderStatus := OrderStatus{}
	strRequest := fmt.Sprintf("/trading/offer/", e.GetSymbolByPair(order.Pair))

	jsonOrderStatus := e.ApiKeyGET(strRequest, make(map[string]interface{}))
	if err := json.Unmarshal([]byte(jsonOrderStatus), &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if orderStatus.Status != "Ok" {
		return fmt.Errorf("%s OrderStatus Failed: %v", e.GetName(), jsonOrderStatus)
	}

	order.StatusMessage = jsonOrderStatus
	currentAmount, err := strconv.ParseFloat(orderStatus.Items[0].CurrentAmount, 64)
	startAmount, err := strconv.ParseFloat(orderStatus.Items[0].StartAmount, 64)
	if err != nil {
		return fmt.Errorf("%s OrderStatus amount parse Failed: %v, %v, %v", e.GetName(), err, orderStatus.Items[0].CurrentAmount, orderStatus.Items[0].StartAmount)
	}
	if currentAmount == startAmount {
		order.Status = exchange.New
	} else if currentAmount == 0.0 {
		order.Status = exchange.Filled
	} else if currentAmount < startAmount {
		order.Status = exchange.Partial
	}

	return nil
}

func (e *Bitbay) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Bitbay) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	mapParams := make(map[string]interface{})
	mapParams["uuid"] = order.OrderID

	side := ""
	if order.Direction == exchange.Buy {
		side = "BUY"
	} else {
		side = "SELL"
	}

	cancelOrder := CancelOrder{}
	strRequest := fmt.Sprintf("/trading/offer/%s/%s/%s/%s", e.GetSymbolByPair(order.Pair), order.OrderID, side, order.Rate)

	jsonCancelOrder := e.ApiKeyGET(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if cancelOrder.Status != "Ok" {
		return fmt.Errorf("%s CancelOrder Failed: %v", e.GetName(), jsonCancelOrder)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Bitbay) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
// ------------------        TODO
func (e *Bitbay) ApiKeyGET(strRequestPath string, mapParams map[string]interface{}) string {
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())

	strUrl := API_URL + strRequestPath + "?" + exchange.Map2UrlQueryInterface(mapParams)

	// signature := exchange.ComputeHmac512NoDecode(strUrl, e.API_SECRET)

	// mapParams["API-Key"] = e.API_KEY
	// mapParams["Request-Timestamp"] = timestamp

	/* jsonParams := ""
	if nil != mapParams {
		bytesParams, _ := json.Marshal(mapParams)
		jsonParams = string(bytesParams)
	} */

	key := []byte(e.API_SECRET)
	h := hmac.New(sha512.New, key)
	// h.Write([]byte(e.API_KEY))
	// h.Write([]byte(timestamp))
	h.Write([]byte(exchange.Map2UrlQueryInterface(mapParams)))
	signature := hex.EncodeToString(h.Sum(nil))

	request, err := http.NewRequest("GET", strUrl, nil) //strings.NewReader(jsonParams)
	if nil != err {
		return err.Error()
	}

	request.Header.Add("Content-Type", "application/json;charset=utf-8")
	request.Header.Add("API-Key", e.API_KEY)
	request.Header.Add("API-Hash", signature)
	request.Header.Add("operation-id", timestamp)
	request.Header.Add("Request-Timestamp", timestamp)

	//request.Header.Add("Accept", "application/json")
	//request.Header.Add("apisign", signature)

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
