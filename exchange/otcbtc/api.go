package otcbtc

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
	API_URL string = "https://bb.otcbtc.com"
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
func (e *Otcbtc) GetCoinsData() error {
	errResponse := &ErrorResponse{}
	pairsData := PairsData{}

	strRequestUrl := "/api/v2/markets"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &pairsData); err != nil {
		if err := json.Unmarshal([]byte(jsonCurrencyReturn), &errResponse); err != nil {
			return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
		} else {
			return fmt.Errorf("%s Get Coins Failed: %v", e.GetName(), errResponse.Error.Code)
		}
	}

	for _, data := range pairsData {
		coinStrs := strings.Split(data.TickerID, "_")
		base := &coin.Coin{}
		target := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base = coin.GetCoin(coinStrs[1])
			if base == nil {
				base = &coin.Coin{}
				base.Code = coinStrs[1]
				coin.AddCoin(base)
			}
			target = coin.GetCoin(coinStrs[0])
			if target == nil {
				target = &coin.Coin{}
				target.Code = coinStrs[0]
				coin.AddCoin(target)
			}
		case exchange.JSON_FILE:
			base = e.GetCoinBySymbol(coinStrs[1])
			target = e.GetCoinBySymbol(coinStrs[0])
		}

		if base != nil {
			coinConstraint := e.GetCoinConstraint(base)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       base.ID,
					Coin:         base,
					ExSymbol:     coinStrs[1],
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = coinStrs[1]
			}
			e.SetCoinConstraint(coinConstraint)
		}

		if target != nil {
			coinConstraint := e.GetCoinConstraint(target)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       target.ID,
					Coin:         target,
					ExSymbol:     coinStrs[0],
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = coinStrs[0]
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
func (e *Otcbtc) GetPairsData() error {
	errResponse := &ErrorResponse{}
	pairsData := PairsData{}

	strRequestUrl := "/api/v2/markets"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &pairsData); err != nil {
		if err := json.Unmarshal([]byte(jsonSymbolsReturn), &errResponse); err != nil {
			return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
		} else {
			return fmt.Errorf("%s Get Pairs Failed: %v", e.GetName(), errResponse.Error.Code)
		}
	}

	for _, data := range pairsData {
		pairStrs := strings.Split(data.TickerID, "_")
		p := &pair.Pair{}

		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(pairStrs[1])
			target := coin.GetCoin(pairStrs[0])
			if base != nil && target != nil {
				p = pair.GetPair(base, target)
			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(data.ID)
		}
		if p != nil {
			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{
					PairID:      p.ID,
					Pair:        p,
					ExSymbol:    data.ID,
					MakerFee:    DEFAULT_MAKER_FEE,
					TakerFee:    DEFAULT_TAKER_FEE,
					LotSize:     data.TradingRule.MinAmount,
					PriceFilter: data.TradingRule.MinPrice,
					Listed:      DEFAULT_LISTED,
				}
			} else {
				pairConstraint.ExSymbol = data.ID
				pairConstraint.LotSize = data.TradingRule.MinAmount
				pairConstraint.PriceFilter = data.TradingRule.MinPrice
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
func (e *Otcbtc) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	errResponse := &ErrorResponse{}
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	strRequestUrl := "/api/v2/depth"
	strUrl := API_URL + strRequestUrl

	mapParams := make(map[string]string)
	mapParams["market"] = symbol
	mapParams["limit"] = "100"

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &orderBook); err != nil {
		if err := json.Unmarshal([]byte(jsonOrderbook), &errResponse); err != nil {
			return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
		} else {
			return nil, fmt.Errorf("%s Get Orderbook Failed: %v %v", e.GetName(), errResponse.Error.Code, errResponse.Error.Message)
		}
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)

	var err error
	for _, bid := range orderBook.Bids {
		buydata := exchange.Order{}

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
	for _, ask := range orderBook.Asks {
		selldata := exchange.Order{}

		//Modify according to type and structure
		selldata.Rate, err = strconv.ParseFloat(ask[0], 64)
		if err != nil {
			return nil, err
		}
		selldata.Quantity, err = strconv.ParseFloat(ask[1], 64)
		if err != nil {
			return nil, err
		}

		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

func (e *Otcbtc) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *Otcbtc) DoAccountOperation(operation *exchange.AccountOperation) error {
	return nil
}

func (e *Otcbtc) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	errResponse := &ErrorResponse{}
	accountBalance := AccountBalance{}
	strRequest := "/api/v2/users/me"

	jsonBalanceReturn := e.ApiKeyGET(make(map[string]string), strRequest)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &errResponse); err != nil {
		log.Printf("%s Get Balance Error Response Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if errResponse.Error.Code != 0 {
		log.Printf("%s Get Balance Failed: %v %v", e.GetName(), errResponse.Error.Code, errResponse.Error.Message)
		return
	} else if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		log.Printf("%s Get Balance Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	}

	if !accountBalance.OtbFeeEnabled {
		return
	}

	for _, data := range accountBalance.Accounts {
		freeamount, err := strconv.ParseFloat(data.Balance, 64)
		if err == nil {
			c := e.GetCoinBySymbol(data.Currency)
			if c != nil {
				balanceMap.Set(c.Code, freeamount)
			}
		} else {
			log.Printf("%s %s Get Balance Err: %s\n", e.GetName(), data.Currency, err)
		}
	}
}

func (e *Otcbtc) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	return false
}

func (e *Otcbtc) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	errResponse := &ErrorResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/api/v2/orders"

	mapParams := make(map[string]string)
	mapParams["market"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "sell"
	mapParams["volume"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)

	jsonPlaceReturn := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &errResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Error Response Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if errResponse.Error.Code != 0 {
		return nil, fmt.Errorf("%s LimitSell Failed: %v %v", e.GetName(), errResponse.Error.Code, errResponse.Error.Message)
	} else if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%v", placeOrder.ID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Otcbtc) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	errResponse := &ErrorResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/api/v2/orders"

	mapParams := make(map[string]string)
	mapParams["market"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "buy"
	mapParams["volume"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)

	jsonPlaceReturn := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &errResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Error Response Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if errResponse.Error.Code != 0 {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v %v", e.GetName(), errResponse.Error.Code, errResponse.Error.Message)
	} else if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%v", placeOrder.ID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Otcbtc) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	errResponse := &ErrorResponse{}
	orderStatus := PlaceOrder{}
	strRequest := "/api/v2/order"

	mapParams := make(map[string]string)
	mapParams["id"] = order.OrderID

	jsonOrderStatus := e.ApiKeyGET(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &errResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Error Response Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if errResponse.Error.Code != 0 {
		return fmt.Errorf("%s OrderStatus Failed: %v %v", e.GetName(), errResponse.Error.Code, errResponse.Error.Message)
	} else if err := json.Unmarshal([]byte(jsonOrderStatus), &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	}

	order.StatusMessage = jsonOrderStatus
	if orderStatus.State == "cancel" {
		order.Status = exchange.Cancelled
	} else if orderStatus.State == "done" {
		order.Status = exchange.Filled
	} else if orderStatus.State == "wait" {
		order.Status = exchange.New
	} else {
		order.Status = exchange.Other
	}

	return nil
}

func (e *Otcbtc) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Otcbtc) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	errResponse := &ErrorResponse{}
	cancelOrder := PlaceOrder{}
	strRequest := "/api/v2/order/delete"

	mapParams := make(map[string]string)
	mapParams["id"] = order.OrderID

	jsonCancelOrder := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &errResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Error Response Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if errResponse.Error.Code != 0 {
		return fmt.Errorf("%s OrderStatus Failed: %v %v", e.GetName(), errResponse.Error.Code, errResponse.Error.Message)
	} else if err := json.Unmarshal([]byte(jsonCancelOrder), &cancelOrder); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Otcbtc) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Otcbtc) ApiKeyGET(mapParams map[string]string, strRequestPath string) string {
	strMethod := "GET"
	strUrl := API_URL + strRequestPath

	mapParams["access_key"] = e.API_KEY
	mapParams["nonce"] = fmt.Sprintf("%d", time.Now().UTC().UnixNano()/int64(time.Millisecond))
	payload := fmt.Sprintf("%s|%s|%s", strMethod, strRequestPath, exchange.Map2UrlQuery(mapParams))

	mapParams["signature"] = exchange.ComputeHmac256NoDecode(payload, e.API_SECRET)

	return exchange.HttpGetRequest(strUrl, mapParams)
}

func (e *Otcbtc) ApiKeyPost(strRequestPath string, mapParams map[string]string) string {
	strMethod := "POST"
	strUrl := API_URL + strRequestPath

	mapParams["access_key"] = e.API_KEY
	mapParams["nonce"] = fmt.Sprintf("%d", time.Now().UTC().UnixNano()/int64(time.Millisecond))
	payload := fmt.Sprintf("%s|%s|%s", strMethod, strRequestPath, exchange.Map2UrlQuery(mapParams))

	mapParams["signature"] = exchange.ComputeHmac256NoDecode(payload, e.API_SECRET)

	httpClient := &http.Client{}

	request, err := http.NewRequest("POST", strUrl, strings.NewReader(exchange.Map2UrlQuery(mapParams)))
	if err != nil {
		return err.Error()
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Accept", "application/json")

	response, err := httpClient.Do(request)
	if err != nil {
		return err.Error()
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err.Error()
	}

	return string(body)
}
