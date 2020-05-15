package bithumb

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
	API_URL string = "https://global-openapi.bithumb.pro/openapi/v1"
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
func (e *Bithumb) GetCoinsData() error {
	jsonResponse := &JsonResponse{}
	coinsData := CoinsData{}

	strRequestUrl := "/spot/config"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if jsonResponse.Code != "0" {
		return fmt.Errorf("%s Get Coins Failed: %v", e.GetName(), jsonCurrencyReturn)
	}
	if err := json.Unmarshal(jsonResponse.Data, &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, data := range coinsData.CoinConfig {
		c := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			c = coin.GetCoin(data.Name)
			if c == nil {
				c = &coin.Coin{}
				c.Code = data.Name
				c.Name = data.FullName
				coin.AddCoin(c)
			}
		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(data.Name)
		}

		if c != nil {
			txFee, err := strconv.ParseFloat(data.WithdrawFee, 64)
			if err != nil {
				return fmt.Errorf("%s Get Coins txFee parse Err: %v %s", e.GetName(), err, data.WithdrawFee)
			}
			coinConstraint := e.GetCoinConstraint(c)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       c.ID,
					Coin:         c,
					ExSymbol:     data.Name,
					ChainType:    exchange.MAINNET,
					TxFee:        txFee,
					Withdraw:     data.WithdrawStatus == "1",
					Deposit:      data.DepositStatus == "1",
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       true,
				}
			} else {
				coinConstraint.ExSymbol = data.Name
				coinConstraint.TxFee = txFee
				if data.WithdrawStatus == "1" {
					coinConstraint.Withdraw = true
				} else {
					coinConstraint.Withdraw = false
				}
				if data.DepositStatus == "1" {
					coinConstraint.Deposit = true
				} else {
					coinConstraint.Deposit = false
				}
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
func (e *Bithumb) GetPairsData() error {
	jsonResponse := &JsonResponse{}
	pairsData := CoinsData{}

	strRequestUrl := "/spot/config"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if jsonResponse.Code != "0" {
		return fmt.Errorf("%s Get Pairs Failed: %v", e.GetName(), jsonSymbolsReturn)
	}
	if err := json.Unmarshal(jsonResponse.Data, &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, data := range pairsData.SpotConfig {
		symbols := strings.Split(data.Symbol, "-")
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

			priceSize, err := strconv.Atoi(data.Accuracy[0])
			if err != nil {
				return fmt.Errorf("%s Get Pairs rate precision parse Err: %v %s", e.GetName(), err, data.Accuracy[0])
			}
			logSize, err := strconv.Atoi(data.Accuracy[1])
			if err != nil {
				return fmt.Errorf("%s Get Pairs amount precision parse Err: %v %s", e.GetName(), err, data.Accuracy[1])
			}
			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{
					PairID:      p.ID,
					Pair:        p,
					ExSymbol:    data.Symbol,
					MakerFee:    DEFAULT_MAKER_FEE,
					TakerFee:    DEFAULT_TAKER_FEE,
					LotSize:     math.Pow10(-1 * logSize),
					PriceFilter: math.Pow10(-1 * priceSize),
					Listed:      true,
				}
			} else {
				pairConstraint.ExSymbol = data.Symbol
				pairConstraint.LotSize = math.Pow10(-1 * logSize)
				pairConstraint.PriceFilter = math.Pow10(-1 * priceSize)
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
func (e *Bithumb) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	jsonResponse := &JsonResponse{}
	orderBook := OrderBook{}

	strRequestUrl := "/spot/orderBook"
	strUrl := API_URL + strRequestUrl

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(pair)

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if jsonResponse.Code != "0" {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %v", e.GetName(), jsonOrderbook)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	var err error
	for _, bid := range orderBook.B {
		buydata := exchange.Order{}

		buydata.Rate, err = strconv.ParseFloat(bid[0], 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Rate error:%v", e.GetName(), err)
		}
		buydata.Quantity, err = strconv.ParseFloat(bid[1], 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
		}

		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.S {
		selldata := exchange.Order{}

		selldata.Rate, err = strconv.ParseFloat(ask[0], 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Rate error:%v", e.GetName(), err)
		}
		selldata.Quantity, err = strconv.ParseFloat(ask[1], 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
		}

		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

/*************** Public API ***************/

/*************** Private API ***************/
func (e *Bithumb) DoAccountOperation(operation *exchange.AccountOperation) error {
	switch operation.Type {
	case exchange.Transfer:
		return e.transfer(operation)
	case exchange.BalanceList:
		return e.getAllBalance(operation)
	case exchange.Balance:
		return e.getBalance(operation)
	case exchange.Withdraw:
		return e.doWithdraw(operation)
	}
	return fmt.Errorf("%s Operation type invalid: %s %v", operation.Ex, operation.Wallet, operation.Type)
}

func (e *Bithumb) getAllBalance(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalances{}
	strRequest := "/spot/assetList"

	mapParams := make(map[string]string)
	if operation.Wallet == exchange.AssetWallet {
		mapParams["assetType"] = "wallet"
	} else if operation.Wallet == exchange.SpotWallet {
		mapParams["assetType"] = "spot"
	} else {
		return fmt.Errorf("%s getAllBalance unexpected Wallet: %s", e.GetName(), operation.Wallet)
	}

	jsonAllBalanceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonAllBalanceReturn
	}

	if err := json.Unmarshal([]byte(jsonAllBalanceReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s getAllBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonAllBalanceReturn)
		return operation.Error
	} else if jsonResponse.Code != "0" {
		operation.Error = fmt.Errorf("%s getAllBalance failed: %v", e.GetName(), jsonAllBalanceReturn)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		operation.Error = fmt.Errorf("%s getAllBalance Data Unmarshal Err: %v, %s", e.GetName(), err, jsonAllBalanceReturn)
		return operation.Error
	}

	for _, v := range accountBalance {
		available, err := strconv.ParseFloat(v.Count, 64)
		frozen, err := strconv.ParseFloat(v.Frozen, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s getAllBalance parse balance failed: %v, %+v", e.GetName(), err, accountBalance)
			return operation.Error
		}

		if available == 0 && frozen == 0 {
			continue
		}

		b := exchange.AssetBalance{
			Coin:             e.GetCoinBySymbol(v.CoinType),
			BalanceAvailable: available,
			BalanceFrozen:    frozen,
		}
		operation.BalanceList = append(operation.BalanceList, b)
	}

	return nil
}

func (e *Bithumb) getBalance(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalances{}
	strRequest := "/spot/assetList"
	symbol := e.GetSymbolByCoin(operation.Coin)

	mapParams := make(map[string]string)
	if operation.Wallet == exchange.AssetWallet {
		mapParams["assetType"] = "wallet"
	} else if operation.Wallet == exchange.SpotWallet {
		mapParams["assetType"] = "spot"
	} else {
		return fmt.Errorf("%s getAllBalance unexpected Wallet: %s", e.GetName(), operation.Wallet)
	}

	jsonBalanceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonBalanceReturn
	}

	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s getAllBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonBalanceReturn)
		return operation.Error
	} else if jsonResponse.Code != "0" {
		operation.Error = fmt.Errorf("%s getAllBalance failed: %v", e.GetName(), jsonBalanceReturn)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		operation.Error = fmt.Errorf("%s getAllBalance Data Unmarshal Err: %v, %s", e.GetName(), err, jsonBalanceReturn)
		return operation.Error
	}

	for _, v := range accountBalance {
		if v.CoinType == symbol {
			available, err := strconv.ParseFloat(v.Count, 64)
			frozen, err := strconv.ParseFloat(v.Frozen, 64)
			if err != nil {
				operation.Error = fmt.Errorf("%s getAllBalance parse balance failed: %v, %+v", e.GetName(), err, accountBalance)
				return operation.Error
			}

			operation.BalanceFrozen = frozen
			operation.BalanceAvailable = available
			return nil
		}
	}

	operation.Error = fmt.Errorf("%s getBalance get %v account balance fail: %v", e.GetName(), symbol, jsonBalanceReturn)
	return operation.Error
}

// TODO get return structure, withdraw ID
func (e *Bithumb) doWithdraw(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	// withdraw := Withdraw{}
	strRequest := "/withdraw"

	mapParams := make(map[string]string)
	mapParams["coinType"] = e.GetSymbolByCoin(operation.Coin)
	mapParams["address"] = operation.WithdrawAddress
	mapParams["quantity"] = operation.WithdrawAmount
	mapParams["mark"] = operation.WithdrawTag

	jsonWithdrawReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonWithdrawReturn
	}

	if err := json.Unmarshal([]byte(jsonWithdrawReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Json Unmarshal Err: %v, %s", e.GetName(), err, jsonWithdrawReturn)
		return operation.Error
	} else if jsonResponse.Code != "0" {
		operation.Error = fmt.Errorf("%s Withdraw Failed: %v", e.GetName(), jsonWithdrawReturn)
		return operation.Error
	}
	// if err := json.Unmarshal(jsonResponse.Data, &withdraw); err != nil {
	// 	operation.Error = fmt.Errorf("%s Withdraw Data Unmarshal Err: %v, %s", e.GetName(), err, jsonWithdrawReturn)
	// 	return operation.Error
	// }

	// operation.WithdrawID = withdrawResponse.WithdrawalID

	return nil
}

// TODO verify
func (e *Bithumb) transfer(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	// transfer := Transfer{}
	strRequest := "/transfer"

	var from, to string
	if operation.TransferFrom == exchange.SpotWallet {
		from = "SPOT"
	} else if operation.TransferFrom == exchange.AssetWallet {
		from = "WALLET"
	} else {
		return fmt.Errorf("%s Transfer unexpected from type: %s", e.GetName(), operation.TransferFrom)
	}
	if operation.TransferDestination == exchange.SpotWallet {
		to = "SPOT"
	} else if operation.TransferDestination == exchange.AssetWallet {
		to = "WALLET"
	} else {
		return fmt.Errorf("%s Transfer unexpected destination type: %s", e.GetName(), operation.TransferDestination)
	}

	mapParams := make(map[string]string)
	mapParams["coinType"] = e.GetSymbolByCoin(operation.Coin)
	mapParams["quantity"] = operation.TransferAmount
	mapParams["from"] = from
	mapParams["to"] = to

	// mapParams["from"] = "SPOT"
	// mapParams["to"] = "LEVER"

	log.Printf("mapParams: %+v", mapParams) //============

	jsonTransferReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonTransferReturn
	}

	if err := json.Unmarshal([]byte(jsonTransferReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s Transfer Json Unmarshal Err: %v, %s", e.GetName(), err, jsonTransferReturn)
		return operation.Error
	} else if jsonResponse.Code != "0" {
		operation.Error = fmt.Errorf("%s Transfer Failed: %v", e.GetName(), jsonTransferReturn)
		return operation.Error
	}
	// if err := json.Unmarshal(jsonResponse.Data, &transfer); err != nil {
	// 	operation.Error = fmt.Errorf("%s Transfer Data Unmarshal Err: %v, %s", e.GetName(), err, jsonWithdrawReturn)
	// 	return operation.Error
	// }

	return nil
}

func (e *Bithumb) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalances{}
	strRequest := "/spot/assetList"

	mapParams := make(map[string]string)
	mapParams["assetType"] = "spot"

	jsonBalanceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if jsonResponse.Code != "0" {
		log.Printf("%s UpdateAllBalances Failed: %s", e.GetName(), jsonBalanceReturn)
		return
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return
	}

	for _, v := range accountBalance {
		c := e.GetCoinBySymbol(v.CoinType)
		if c != nil {
			freeamount, err := strconv.ParseFloat(v.Count, 64)
			if err == nil {
				balanceMap.Set(c.Code, freeamount)
			}
		}
	}
}

func (e *Bithumb) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	/* if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil", e.GetName())
		return false
	}

	mapParams := make(map[string]string)
	mapParams["currency"] = e.GetSymbolByCoin(coin)
	mapParams["quantity"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["address"] = addr

	jsonResponse := &JsonResponse{}
	uuid := Uuid{}
	strRequest := "/v1.1/account/withdraw"

	jsonSubmitWithdraw := e.ApiKeyGET(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &jsonResponse); err != nil {
		log.Printf("%s Withdraw Json Unmarshal Err: %v %v", e.GetName(), err, jsonSubmitWithdraw)
		return false
	} else if jsonResponse.Code != "0" {
		log.Printf("%s Withdraw Failed: %v", e.GetName(), jsonResponse.Message)
		return false
	}
	if err := json.Unmarshal(jsonResponse.Data, &uuid); err != nil {
		log.Printf("%s Withdraw Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return false
	}
	return true */

	return false
}

func (e *Bithumb) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/spot/placeOrder"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["quantity"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["side"] = "sell"
	mapParams["type"] = "limit"

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != "0" {
		return nil, fmt.Errorf("%s LimitSell Failed: %s", e.GetName(), jsonPlaceReturn)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
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

func (e *Bithumb) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/spot/placeOrder"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["quantity"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["side"] = "buy"
	mapParams["type"] = "limit"

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != "0" {
		return nil, fmt.Errorf("%s LimitBuy Failed: %s", e.GetName(), jsonPlaceReturn)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
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

func (e *Bithumb) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	mapParams := make(map[string]string)
	mapParams["orderId"] = order.OrderID
	mapParams["symbol"] = e.GetSymbolByPair(order.Pair)

	jsonResponse := &JsonResponse{}
	orderStatus := OrderStatus{}
	strRequest := "/spot/singleOrder"

	jsonOrderStatus := e.ApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if jsonResponse.Code != "0" {
		return fmt.Errorf("%s OrderStatus Failed: %s", e.GetName(), jsonOrderStatus)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	} else if string(jsonResponse.Data) == "null" {
		return fmt.Errorf("%s OrderStatus empty return: %s", e.GetName(), jsonOrderStatus)
	}

	tradeNum, err := strconv.ParseFloat(orderStatus.TradedNum, 64)
	if err != nil {
		return fmt.Errorf("%s OrderStatus parse tradeNum fail: %v", e.GetName(), orderStatus.TradedNum)
	}

	order.StatusMessage = jsonOrderStatus
	if tradeNum == 0 && (orderStatus.Status == "send" || orderStatus.Status == "pending") {
		order.Status = exchange.New
	} else if tradeNum > 0 && (orderStatus.Status == "send" || orderStatus.Status == "pending") {
		order.Status = exchange.Partial
	} else if orderStatus.Status == "success" {
		order.Status = exchange.Filled
	} else if orderStatus.Status == "cancel" {
		order.Status = exchange.Cancelled
	} else {
		order.Status = exchange.Other
	}

	return nil
}

func (e *Bithumb) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Bithumb) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	if order == nil {
		return fmt.Errorf("%s CancelOrder Failed, nil order", e.GetName())
	}

	mapParams := make(map[string]string)
	mapParams["orderId"] = order.OrderID
	mapParams["symbol"] = e.GetSymbolByPair(order.Pair)

	jsonResponse := &JsonResponse{}
	// cancelOrder := CancelOrder{}
	strRequest := "/spot/cancelOrder"

	jsonCancelOrder := e.ApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if jsonResponse.Code != "0" {
		return fmt.Errorf("%s CancelOrder Failed: %s", e.GetName(), jsonCancelOrder)
	}
	// if err := json.Unmarshal(jsonResponse.Data, &cancelOrder); err != nil {
	// 	return fmt.Errorf("%s CancelOrder Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	// }

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Bithumb) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Bithumb) ApiKeyGET(strRequestPath string, mapParams map[string]string) string {
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

func (e *Bithumb) ApiKeyRequest(strMethod string, mapParams map[string]string, strRequestPath string) string {
	mapParams["timestamp"] = fmt.Sprintf("%.0d", time.Now().UnixNano()/1e6)
	mapParams["apiKey"] = e.API_KEY
	mapParams["msgNo"] = "1234561284"

	strUrl := API_URL + strRequestPath

	var strParams string
	if nil != mapParams {
		strParams = exchange.Map2UrlQuery(mapParams)
	}

	// log.Printf("strParams: %v", strParams) //===========================

	signature := exchange.ComputeHmac256NoDecode(strParams, e.API_SECRET)
	signMessage := strUrl + "?" + strParams
	mapParams["signature"] = strings.ToLower(signature) // need lower

	jsonParams := ""
	if nil != mapParams {
		bytesParams, _ := json.Marshal(mapParams)
		jsonParams = string(bytesParams)
	}

	request, err := http.NewRequest(strMethod, signMessage, strings.NewReader(jsonParams))
	if nil != err {
		return err.Error()
	}

	request.Header.Add("Content-Type", "application/json")

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
