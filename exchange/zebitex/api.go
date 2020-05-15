package zebitex

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
)

/*The Base Endpoint URL*/
const (
	API_URL = "https://zebitex.com"
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
func (e *Zebitex) GetCoinsData() error {
	coinsData := CoinsData{}

	strRequestPath := "/api/v1/orders/tickers"
	strUrl := API_URL + strRequestPath

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v\n", e.GetName(), err, jsonCurrencyReturn)
	}

	for _, data := range coinsData {
		base := &coin.Coin{}
		target := &coin.Coin{}

		switch e.Source {
		case exchange.EXCHANGE_API:
			base = coin.GetCoin(data.QuoteUnit)
			if base == nil {
				base = &coin.Coin{}
				base.Code = data.QuoteUnit
				coin.AddCoin(base)
			}
			target = coin.GetCoin(data.BaseUnit)
			if target == nil {
				target = &coin.Coin{}
				target.Code = data.BaseUnit
				coin.AddCoin(target)
			}
		case exchange.JSON_FILE:
			base = e.GetCoinBySymbol(data.QuoteUnit)
			target = e.GetCoinBySymbol(data.BaseUnit)
		}

		if base != nil {
			coinConstraint := e.GetCoinConstraint(base)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       base.ID,
					Coin:         base,
					ExSymbol:     data.QuoteUnit,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = data.QuoteUnit
			}
			e.SetCoinConstraint(coinConstraint)
		}

		if target != nil {
			coinConstraint := e.GetCoinConstraint(target)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       target.ID,
					Coin:         target,
					ExSymbol:     data.BaseUnit,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = data.BaseUnit
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
func (e *Zebitex) GetPairsData() error {
	pairsData := PairsData{}

	strRequestPath := "/api/v1/orders/tickers"
	strUrl := API_URL + strRequestPath

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v\n", e.GetName(), err, jsonSymbolsReturn)
	}

	for _, data := range pairsData {
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(data.QuoteUnit)
			target := coin.GetCoin(data.BaseUnit)
			if base != nil && target != nil {
				p = pair.GetPair(base, target)
			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(data.Market)
		}

		if p != nil {
			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{
					PairID:      p.ID,
					Pair:        p,
					ExSymbol:    data.Market,
					MakerFee:    data.AskFee,
					TakerFee:    data.BidFee,
					LotSize:     DEFAULT_LOT_SIZE,
					PriceFilter: DEFAULT_PRICE_FILTER,
					Listed:      true,
				}
			} else {
				pairConstraint.ExSymbol = data.Market
				pairConstraint.MakerFee = data.AskFee
				pairConstraint.TakerFee = data.BidFee
			}
			e.SetPairConstraint(pairConstraint)
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
func (e *Zebitex) OrderBook(p *pair.Pair) (*exchange.Maker, error) {
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(p)

	mapParams := make(map[string]string)
	mapParams["market"] = symbol

	strRequestPath := "/api/v1/orders/orderbook"
	strUrl := API_URL + strRequestPath

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v\n", e.GetName(), err, jsonOrderbook)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)

	var err error
	//买入
	for _, bid := range orderBook.Bids {
		buydata := exchange.Order{}
		buydata.Quantity, err = strconv.ParseFloat(bid[1].(string), 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v\n", e.GetName(), err)
		}

		buydata.Rate, err = strconv.ParseFloat(bid[0].(string), 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Rate error:%v\n", e.GetName(), err)
		}

		maker.Bids = append(maker.Bids, buydata)
	}

	//卖出
	for _, ask := range orderBook.Asks {
		selldata := exchange.Order{}
		selldata.Quantity, err = strconv.ParseFloat(ask[1].(string), 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v\n", e.GetName(), err)
		}

		selldata.Rate, err = strconv.ParseFloat(ask[0].(string), 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Rate error:%v\n", e.GetName(), err)
		}

		maker.Asks = append(maker.Asks, selldata)
	}

	return maker, err
}

/*************** Private API ***************/
func (e *Zebitex) DoAccountOperation(operation *exchange.AccountOperation) error {
	switch operation.Type {
	case exchange.BalanceList:
		return e.getAllBalance(operation)
	case exchange.Balance:
		return e.getBalance(operation)
	case exchange.Withdraw: // need test
		return e.doWithdraw(operation)
	}
	return fmt.Errorf("%s Operation type invalid: %s %v", operation.Ex, operation.Wallet, operation.Type)
}

func (e *Zebitex) getAllBalance(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	accountID := AccountBalances{}
	strRequest := "/api/v1/funds"

	jsonAllBalanceReturn := e.ApiKeyGet(strRequest, nil)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonAllBalanceReturn
	}

	// log.Printf("jsonAllBalanceReturn: %v", jsonAllBalanceReturn)
	if err := json.Unmarshal([]byte(jsonAllBalanceReturn), &accountID); err != nil {
		operation.Error = fmt.Errorf("%s getAllBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonAllBalanceReturn)
		return operation.Error
	} else if strings.Contains(jsonAllBalanceReturn, "error") {
		operation.Error = fmt.Errorf("%s getAllBalance Failed: %v", e.GetName(), jsonAllBalanceReturn)
		return operation.Error
	}
	// if err := json.Unmarshal(jsonResponse.Data, &accountID); err != nil {
	// 	operation.Error = fmt.Errorf("%s getAllBalance Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	// 	return operation.Error
	// }

	for _, account := range accountID {
		if account.Balance == "0.00000000" && account.LockedBalance == "0.00000000" {
			continue
		}
		frozen, err := strconv.ParseFloat(account.LockedBalance, 64)
		available, err := strconv.ParseFloat(account.Balance, 64)
		if err != nil {
			return fmt.Errorf("%s balance parse fail: %v %+v", e.GetName(), err, account)
		}

		if frozen+available == 0 {
			continue
		}

		balance := exchange.AssetBalance{
			Coin:             e.GetCoinBySymbol(strings.ToLower(account.Code)),
			BalanceAvailable: available,
			BalanceFrozen:    frozen,
		}
		operation.BalanceList = append(operation.BalanceList, balance)

	}

	return nil
	// return fmt.Errorf("%s getBalance fail: %v", e.GetName(), jsonBalanceReturn)
}

func (e *Zebitex) getBalance(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	accountID := AccountBalances{}
	strRequest := "/api/v1/funds"
	symbol := e.GetSymbolByCoin(operation.Coin)

	jsonBalanceReturn := e.ApiKeyGet(strRequest, nil)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonBalanceReturn
	}

	// log.Printf("jsonBalanceReturn: %v", jsonBalanceReturn)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountID); err != nil {
		operation.Error = fmt.Errorf("%s getBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonBalanceReturn)
		return operation.Error
	} else if strings.Contains(jsonBalanceReturn, "error") {
		operation.Error = fmt.Errorf("%s getBalance Failed: %v", e.GetName(), jsonBalanceReturn)
		return operation.Error
	}

	for _, account := range accountID {
		if strings.ToLower(account.Code) == symbol {
			frozen, err := strconv.ParseFloat(account.LockedBalance, 64)
			allBalance, err := strconv.ParseFloat(account.Balance, 64)
			if err != nil {
				return fmt.Errorf("%s balance parse fail: %v %+v", e.GetName(), err, account)
			}
			operation.BalanceFrozen = frozen
			operation.BalanceAvailable = allBalance - frozen
			return nil
		}
	}

	return fmt.Errorf("%s getBalance fail: %v", e.GetName(), jsonBalanceReturn)
}

// need test
func (e *Zebitex) doWithdraw(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	if operation.WithdrawTag != "" {
		operation.Error = fmt.Errorf("%s Withdraw Failed, got tag: %v, for coin: %v", e.GetName(), operation.WithdrawTag, operation.Coin.Code)
		return operation.Error
	}

	currency := e.GetSymbolByCoin(operation.Coin)

	// create one fund source
	// _, source := e.CreateFundSource(currency, "", addr)
	// select an existing fund source
	sources := e.getFundSource(currency)
	sourceID := ""
	for _, source := range sources {
		if source.Address == operation.WithdrawAddress {
			sourceID = fmt.Sprintf("%v", source.Id)
		}
	}
	if sourceID == "" {
		log.Printf("%s Withdraw fail, addr not found\n", e.GetName())
		return operation.Error
	}

	withdraw := WithdrawResponse{}
	strRequestPath := "/api/v1/withdrawals"

	mapParams := make(map[string]string)
	mapParams["code"] = currency
	mapParams["fund_source_id"] = sourceID // fmt.Sprintf("%d", source.Id)
	mapParams["sum"] = operation.WithdrawAmount
	jsonSubmitWithdraw, code := e.ApiKeyRequest("POST", strRequestPath, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequestPath
		operation.CallResponce = jsonSubmitWithdraw
	}
	if code == 204 {
		return nil
	}

	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &withdraw); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Json Unmarshal Err: %v %v\n", e.GetName(), err, jsonSubmitWithdraw)
	} else {
		operation.Error = fmt.Errorf("%s Withdraw fail: %s\n", e.GetName(), jsonSubmitWithdraw)
	}

	return operation.Error
}

func (e *Zebitex) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	accountBalance := AccountBalances{}
	strRequestPath := "/api/v1/funds"
	jsonBalanceReturn := e.ApiKeyGet(strRequestPath, nil)
	// log.Printf("Balance Return: %v", jsonBalanceReturn) // =====================
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if strings.Contains(jsonBalanceReturn, "error") {
		log.Printf("%s UpdateAllBalances failed: %v\n", e.GetName(), jsonBalanceReturn)
		return
	}

	for _, balance := range accountBalance {
		c := e.GetCoinBySymbol(strings.ToLower(balance.Code))
		if c != nil {
			balanceNum, err := strconv.ParseFloat(balance.Balance, 64)
			if err != nil {
				log.Printf("%s balance parse Err: %v, %v", e.GetName(), err, balance.Balance)
				return
			}
			balanceMap.Set(c.Code, balanceNum)
		}
	}

}

/* FundSources(currency string) */
func (e *Zebitex) getFundSource(currency string) []FundSource {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return nil
	}

	sources := FundSources{}
	strRequestPath := "/api/v1/fund_sources"

	mapParams := make(map[string]string)
	mapParams["currency"] = currency
	jsonSources := e.ApiKeyGet(strRequestPath, mapParams)
	if err := json.Unmarshal([]byte(jsonSources), &sources); err != nil {
		log.Printf("getFundSource %s fundSources Json Unmarshal Err: %v %v\n", e.GetName(), err, jsonSources)
	}

	return sources
}

/* delFundSource(sources []FundSource) */
func (e *Zebitex) delFundSource(sources []FundSource) bool {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return false
	}

	strRequestPath := "/api/v1/fund_sources/selected"
	mapParams := make(map[string]string)

	ids := []int64{}
	for _, source := range sources {
		ids = append(ids, source.Id)
	}

	idsStr, _ := json.Marshal(ids)
	mapParams["ids"] = string(idsStr)
	_, code := e.ApiKeyRequest("DELETE", strRequestPath, mapParams)
	if code == 204 {
		return true
	}

	return false
}

// don't use this for withdraw
/* CreateFundSource(currency, label, addr string) */
func (e *Zebitex) CreateFundSource(currency, label, addr string) (bool, FundSource) {
	source := FundSource{}
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return false, source
	}

	if label == "" {
		label = currency + fmt.Sprintf("%d", rand.Int())
	}

	strRequestPath := "/api/v1/fund_sources"

	mapParams := make(map[string]string)
	mapParams["currency"] = strings.ToLower(currency)
	mapParams["extra"] = label // "0xaC05f7b683b14e5997d288a8C031c5143533F9e3" // kraken eth
	mapParams["uid"] = addr    // "0xaC05f7b683b14e5997d288a8C031c5143533F9e3"

	jsonSource, code := e.ApiKeyRequest("POST", strRequestPath, mapParams)
	log.Printf("============Source json: %v", jsonSource) // ==================
	if err := json.Unmarshal([]byte(jsonSource), &source); err != nil {
		log.Printf("CreateFundSource %s fundSources Json Unmarshal Err: %v %v\n", e.GetName(), err, jsonSource)
	} else if code == 200 {
		log.Printf("============Source: %v", source) // ==================
		return true, source
	}

	return false, source
}

/* Withdraw(coin *coin.Coin, quantity float64, addr, tag string) */
func (e *Zebitex) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return false
	}

	currency := e.GetSymbolByCoin(coin)

	// create one fund source
	// _, source := e.CreateFundSource(currency, "", addr)
	// select an existing fund source
	sources := e.getFundSource(currency)
	sourceID := ""
	for _, source := range sources {
		if source.Address == addr {
			sourceID = fmt.Sprintf("%v", source.Id)
		}
	}
	if sourceID == "" {
		log.Printf("%s Withdraw fail, addr not found\n", e.GetName())
		return false
	}

	withdraw := WithdrawResponse{}
	strRequestPath := "/api/v1/withdrawals"

	mapParams := make(map[string]string)
	mapParams["code"] = currency
	mapParams["fund_source_id"] = sourceID // fmt.Sprintf("%d", source.Id)
	mapParams["sum"] = strconv.FormatFloat(quantity, 'f', 6, 64)
	jsonSubmitWithdraw, code := e.ApiKeyRequest("POST", strRequestPath, mapParams)
	if code == 204 {
		return true
	}

	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &withdraw); err != nil {
		log.Printf("%s Withdraw Json Unmarshal Err: %v %v\n", e.GetName(), err, jsonSubmitWithdraw)
	} else {
		log.Printf("%s Withdraw fail: %s\n", e.GetName(), jsonSubmitWithdraw)
	}

	return false
}

// need to test parameter values
func (e *Zebitex) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.\n", e.GetName())
	}

	placeOrder := PlaceOrder{}
	strRequestPath := "/api/v1/orders"

	mapParams := make(map[string]string)
	mapParams["side"] = "ask"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["market"] = e.GetSymbolByPair(pair)
	mapParams["ordType"] = "limit"

	jsonPlaceReturn, _ := e.ApiKeyRequest("POST", strRequestPath, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v\n", e.GetName(), err, jsonPlaceReturn)
	} else if strings.Contains(jsonPlaceReturn, "error") {
		return nil, fmt.Errorf("%s LimitSell failed: %v\n", e.GetName(), jsonPlaceReturn)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%v", placeOrder.Id),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}
	return order, nil
}

func (e *Zebitex) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.\n", e.GetName())
	}

	placeOrder := PlaceOrder{}
	strRequestPath := "/api/v1/orders"

	mapParams := make(map[string]string)
	mapParams["side"] = "bid"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["market"] = e.GetSymbolByPair(pair)
	mapParams["ordType"] = "limit"

	jsonPlaceReturn, _ := e.ApiKeyRequest("POST", strRequestPath, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v\n", e.GetName(), err, jsonPlaceReturn)
	} else if strings.Contains(jsonPlaceReturn, "error") {
		return nil, fmt.Errorf("%s LimitBuy failed: %v\n", e.GetName(), jsonPlaceReturn)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%v", placeOrder.Id),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}
	return order, nil
}

// Update delay after order cancelled
func (e *Zebitex) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.\n", e.GetName())
	}

	if order == nil {
		return fmt.Errorf("%s OrderStatus failed, nil order input", e.GetName())
	}

	orders := OrdersPage{}
	strRequestPath := "/api/v1/orders/current"

	mapParams := make(map[string]string)
	mapParams["page"] = "1"
	mapParams["per"] = "100"

	jsonOrderStatus := e.ApiKeyGet(strRequestPath, mapParams)
	// log.Printf("Order Return: %v", jsonOrderStatus) // ==========================
	if err := json.Unmarshal([]byte(jsonOrderStatus), &orders); err != nil {
		return fmt.Errorf("%s OrdersPage Json Unmarshal Err: %v %v\n", e.GetName(), err, jsonOrderStatus)
	} else if strings.Contains(jsonOrderStatus, "error") {
		return fmt.Errorf("%s OrderStatus failed: %v\n", e.GetName(), jsonOrderStatus)
	}

	order.StatusMessage = jsonOrderStatus
	for _, orderItem := range orders.Items {
		if string(orderItem.Id) == order.OrderID {
			state := strings.ToUpper(orderItem.State)
			switch state {
			case "wait":
				order.Status = exchange.New
			// not verified
			case "CANCELED":
				order.Status = exchange.Cancelled
			case "FILLED":
				order.Status = exchange.Filled
			case "PARTIALLY_FILLED":
				order.Status = exchange.Partial
			case "REJECTED":
				order.Status = exchange.Rejected
			case "Expired":
				order.Status = exchange.Expired
			case "NEW":
				order.Status = exchange.New
			default:
				order.Status = exchange.Other
			}
			break

			// order.DealRate, _ = strconv.ParseFloat(orderItem.Price, 64)
			// order.DealQuantity, _ = strconv.ParseFloat(orderItem.Filled, 64) // could be 0
		}
	}
	if len(orders.Items) == 0 {
		order.Status = exchange.Cancelled
	}

	return nil
}

func (e *Zebitex) ListOrders() ([]*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.\n", e.GetName())
	}

	orders := OrdersPage{}
	strRequestPath := "/api/v1/orders/day_history"

	mapParams := make(map[string]string)
	mapParams["page"] = "1"
	mapParams["per"] = "100"

	jsonOrderStatus := e.ApiKeyGet(strRequestPath, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &orders); err != nil {
		return nil, fmt.Errorf("%s OrdersPage Json Unmarshal Err: %v %v\n", e.GetName(), err, jsonOrderStatus)
	}

	var res []*exchange.Order
	for _, orderItem := range orders.Items {
		pair := e.GetPairBySymbol(orderItem.Pair)
		rate, _ := strconv.ParseFloat(orderItem.Price, 64)
		quantity, _ := strconv.ParseFloat(orderItem.Amount, 64)

		order := &exchange.Order{
			Pair:         pair,
			OrderID:      string(orderItem.Id),
			Rate:         rate,
			Quantity:     quantity,
			Status:       exchange.New,
			JsonResponse: jsonOrderStatus,
		}

		switch orderItem.Side {
		case "buy":
			order.Direction = exchange.Buy
		case "sell":
			order.Direction = exchange.Sell
		}

		res = append(res, order)
	}

	return res, nil
}

func (e *Zebitex) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.\n", e.GetName())
	}

	mapParams := make(map[string]string)
	mapParams["id"] = order.OrderID

	strRequestPath := fmt.Sprintf("/api/v1/orders/%v/cancel", order.OrderID)
	cont, code := e.ApiKeyRequest("DELETE", strRequestPath, mapParams)
	if code != 204 {
		return fmt.Errorf("%s CancelOrder Failed: %v, %v\n", e.GetName(), code, cont)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = cont

	return nil
}

func (e *Zebitex) CancelAllOrder() error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.\n", e.GetName())
	}

	strRequestPath := "/api/v1/orders/cancel_all"
	_, code := e.ApiKeyRequest("DELETE", strRequestPath, nil)
	if code != 204 {
		return fmt.Errorf("CancelAllOrder Failed: %v\n", code)
	}

	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Get Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Zebitex) ApiKeyGet(strRequestPath string, mapParams map[string]string) string {
	res, _ := e.ApiKeyRequest("GET", strRequestPath, mapParams)
	return res
}

/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request*/
func (e *Zebitex) ApiKeyRequest(strMethod, strRequestPath string, mapParams map[string]string) (string, int) {
	strMethod = strings.ToUpper(strMethod)
	strUrl := API_URL + strRequestPath

	millTime := time.Now().UnixNano() / int64(time.Millisecond)
	var paramStr []byte
	var fields string
	if len(mapParams) > 0 {
		//sort params for sign
		var paramKeys []string
		for k := range mapParams {
			paramKeys = append(paramKeys, k)
		}
		sort.Strings(paramKeys)
		fields = strings.Join(paramKeys, ";")

		paramStr, _ = json.Marshal(mapParams)
	} else {
		fields = ""
		paramStr = []byte("{}")
	}

	payloadStr := fmt.Sprintf("%s|%s|%d|%s", strMethod, strRequestPath, millTime, paramStr)

	//make sign
	sign := exchange.ComputeHmac256NoDecode(payloadStr, e.API_SECRET)
	authStr := fmt.Sprintf("ZEBITEX-HMAC-SHA256 access_key=%s, signature=%s, tonce=%d, signed_params=%s", e.API_KEY, sign, millTime, fields)

	request, err := http.NewRequest(strMethod, strUrl, bytes.NewBuffer(paramStr))
	if nil != err {
		return err.Error(), 0
	}
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	request.Header.Add("Authorization", authStr)

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if nil != err {
		return err.Error(), 0
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return err.Error(), 0
	}

	return string(body), response.StatusCode
}
