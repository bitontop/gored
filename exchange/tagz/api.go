package tagz

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

/*The Base Endpoint URL*/
const (
	API_URL = "https://trade.tagz.com"
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
func (e *Tagz) GetCoinsData() error {
	pairsData := PairsData{}

	strRequestUrl := "/marketdata/market/symbols"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &pairsData); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if len(pairsData) == 0 {
		return fmt.Errorf("%s Get Coins Failed: %s", e.GetName(), jsonCurrencyReturn)
	}

	for _, data := range pairsData {
		baseSymbol := strings.Split(data, "_")[1]
		targetSymbol := strings.Split(data, "_")[0]
		base := &coin.Coin{}
		target := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base = coin.GetCoin(baseSymbol)
			if base == nil {
				base = &coin.Coin{}
				base.Code = baseSymbol
				coin.AddCoin(base)
			}
			target = coin.GetCoin(targetSymbol)
			if target == nil {
				target = &coin.Coin{}
				target.Code = targetSymbol
				coin.AddCoin(target)
			}
		case exchange.JSON_FILE:
			base = e.GetCoinBySymbol(baseSymbol)
			target = e.GetCoinBySymbol(targetSymbol)
		}

		if base != nil {
			coinConstraint := e.GetCoinConstraint(base)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       base.ID,
					Coin:         base,
					ExSymbol:     baseSymbol,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       true,
				}
			} else {
				coinConstraint.ExSymbol = baseSymbol
			}
			e.SetCoinConstraint(coinConstraint)
		}

		if target != nil {
			coinConstraint := e.GetCoinConstraint(target)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       target.ID,
					Coin:         target,
					ExSymbol:     targetSymbol,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       true,
				}
			} else {
				coinConstraint.ExSymbol = targetSymbol
			}
			e.SetCoinConstraint(coinConstraint)
		}
	}

	// limit 1 per 1s
	time.Sleep(time.Second * 2)

	return nil
}

/* GetPairsData - Get Pairs Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func (e *Tagz) GetPairsData() error {
	pairsData := PairsData{}

	strRequestUrl := "/marketdata/market/symbols"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if len(pairsData) == 0 {
		return fmt.Errorf("%s Get Pairs Failed: %s", e.GetName(), jsonSymbolsReturn)
	}

	for _, data := range pairsData {
		baseSymbol := strings.Split(data, "_")[1]
		targetSymbol := strings.Split(data, "_")[0]
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(baseSymbol)
			target := coin.GetCoin(targetSymbol)
			if base != nil && target != nil {
				p = pair.GetPair(base, target)
			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(data)
		}

		if p != nil {
			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{
					PairID:      p.ID,
					Pair:        p,
					ExSymbol:    data,
					MakerFee:    DEFAULT_MAKER_FEE,
					TakerFee:    DEFAULT_TAKER_FEE,
					LotSize:     DEFAULT_LOT_SIZE,
					PriceFilter: DEFAULT_PRICE_FILTER,
					Listed:      DEFAULT_LISTED,
				}
			} else {
				pairConstraint.ExSymbol = data
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
func (e *Tagz) OrderBook(p *pair.Pair) (*exchange.Maker, error) {
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(p)

	strRequestUrl := fmt.Sprintf("/marketdata/instruments/%v/depth", symbol)
	strUrl := API_URL + strRequestUrl

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if len(orderBook.Bids) == 0 && len(orderBook.Asks) == 0 {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %s", e.GetName(), jsonOrderbook)
	}

	maker.LastUpdateID = int64(orderBook.Version)
	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)

	// for _, bid := range orderBook.Bids {
	// 	buydata := exchange.Order{}

	// 	buydata.Rate = bid.Price
	// 	buydata.Quantity = bid.Amount

	// 	maker.Bids = append(maker.Bids, buydata)
	// }

	for i := len(orderBook.Bids) - 1; i >= 0; i-- {
		buydata := exchange.Order{}

		buydata.Rate = orderBook.Bids[i].Price
		buydata.Quantity = orderBook.Bids[i].Amount

		maker.Bids = append(maker.Bids, buydata)
	}

	for _, ask := range orderBook.Asks {
		selldata := exchange.Order{}

		selldata.Rate = ask.Price
		selldata.Quantity = ask.Amount

		maker.Asks = append(maker.Asks, selldata)
	}

	// for i := len(orderBook.Asks) - 1; i >= 0; i-- {
	// 	selldata := exchange.Order{}

	// 	selldata.Rate = strconv.ParseFloat(orderBook.Asks[i][0], 64)
	// 	selldata.Quantity, err = strconv.ParseFloat(orderBook.Asks[i][1], 64)

	// 	maker.Asks = append(maker.Asks, selldata)
	// }

	return maker, nil
}

/*************** Public API ***************/
func (e *Tagz) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *Tagz) DoAccountOperation(operation *exchange.AccountOperation) error {
	switch operation.Type {
	// case exchange.Transfer:
	// 	return e.transfer(operation)
	// case exchange.BalanceList:
	// 	return e.getAllBalance(operation)
	// case exchange.Balance:
	// 	return e.getBalance(operation)
	// case exchange.Withdraw:
	// 	return e.doWithdraw(operation)
	}
	return fmt.Errorf("%s Operation type invalid: %s %v", operation.Ex, operation.Wallet, operation.Type)
}

func (e *Tagz) doWithdraw(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("Tagz API Key or Secret Key or passphrase are nil.")
	}

	// need to use inner transfer before withdraw

	jsonResponse := JsonResponse{}
	withdraw := Withdraw{}
	strRequestUrl := "/api/v1/withdrawals"

	mapParams := make(map[string]string)
	mapParams["currency"] = e.GetSymbolByCoin(operation.Coin)
	mapParams["address"] = operation.WithdrawAddress
	mapParams["amount"] = operation.WithdrawAmount

	jsonCreateWithdraw := e.ApiKeyRequest("POST", strRequestUrl, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequestUrl
		operation.CallResponce = jsonCreateWithdraw
	}

	if err := json.Unmarshal([]byte(jsonCreateWithdraw), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Json Unmarshal Err: %v, %s", e.GetName(), err, jsonCreateWithdraw)
		return operation.Error
	} else if jsonResponse.Code != "200000" {
		operation.Error = fmt.Errorf("%s Withdraw Failed: %v", e.GetName(), jsonCreateWithdraw)
		return operation.Error
	}

	if err := json.Unmarshal(jsonResponse.Data, &withdraw); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	operation.WithdrawID = withdraw.WithdrawalID

	return nil
}

func (e *Tagz) transfer(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
	}

	innerTrans := InnerTrans{}
	strRequestUrl := "/api/v2/accounts/inner-transfer"

	mapParams := make(map[string]string)
	mapParams["clientOid"] = "12345"
	mapParams["currency"] = e.GetSymbolByCoin(operation.Coin)
	mapParams["amount"] = operation.TransferAmount
	switch operation.TransferFrom {
	case exchange.AssetWallet:
		mapParams["from"] = "main"
	case exchange.SpotWallet:
		mapParams["from"] = "trade"
	}
	switch operation.TransferDestination {
	case exchange.AssetWallet:
		mapParams["to"] = "main"
	case exchange.SpotWallet:
		mapParams["to"] = "trade"
	}

	jsonTransferReturn := e.ApiKeyRequest("POST", strRequestUrl, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequestUrl
		operation.CallResponce = jsonTransferReturn
	}

	// log.Printf("return: %v", jsonTransferReturn)
	if err := json.Unmarshal([]byte(jsonTransferReturn), &innerTrans); err != nil {
		operation.Error = fmt.Errorf("%s InnerTrans Json Unmarshal Err: %v, %s", e.GetName(), err, jsonTransferReturn)
		return operation.Error
	} else if innerTrans.Code != "200000" {
		operation.Error = fmt.Errorf("%s InnerTrans Failed: %v", e.GetName(), jsonTransferReturn)
		return operation.Error
	} else if innerTrans.Msg != "" {
		operation.Error = fmt.Errorf("%s InnerTrans Failed: %v", e.GetName(), jsonTransferReturn)
		return operation.Error
	}

	log.Printf("InnerTrans response %v", jsonTransferReturn)

	return nil
}

func (e *Tagz) getAllBalance(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	accountID := AccountID{}
	strRequest := "/api/v1/accounts"
	accountType := ""
	// balanceList := []exchange.AssetBalance{}

	mapParams := make(map[string]string)
	if operation.Wallet == exchange.AssetWallet {
		mapParams["type"] = "main" // "trade"
		accountType = "main"
	} else if operation.Wallet == exchange.SpotWallet {
		mapParams["type"] = "trade"
		accountType = "trade"
	}

	jsonAllBalanceReturn := e.ApiKeyRequest("GET", strRequest, nil)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonAllBalanceReturn
	}

	// log.Printf("jsonAllBalanceReturn: %v", jsonAllBalanceReturn)
	if err := json.Unmarshal([]byte(jsonAllBalanceReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s getAllBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonAllBalanceReturn)
		return operation.Error
	} else if jsonResponse.Code != "200000" {
		operation.Error = fmt.Errorf("%s Withdraw Failed: %v", e.GetName(), jsonAllBalanceReturn)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountID); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	for _, account := range accountID {
		if account.Type == accountType {
			frozen, err := strconv.ParseFloat(account.Holds, 64)
			avaliable, err := strconv.ParseFloat(account.Available, 64)
			if err != nil {
				return fmt.Errorf("%s balance parse fail: %v %+v", e.GetName(), err, account)
			}

			balance := exchange.AssetBalance{
				Coin:             e.GetCoinBySymbol(account.Currency),
				BalanceAvailable: avaliable,
				BalanceFrozen:    frozen,
			}
			operation.BalanceList = append(operation.BalanceList, balance)
		}
	}

	return nil
	// return fmt.Errorf("%s getBalance fail: %v", e.GetName(), jsonBalanceReturn)
}

func (e *Tagz) getBalance(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	accountID := AccountID{}
	strRequest := "/api/v1/accounts"
	accountType := ""

	mapParams := make(map[string]string)
	mapParams["currency"] = e.GetSymbolByCoin(operation.Coin)
	if operation.Wallet == exchange.AssetWallet {
		mapParams["type"] = "main" // "trade"
		accountType = "main"
	} else if operation.Wallet == exchange.SpotWallet {
		mapParams["type"] = "trade"
		accountType = "trade"
	}

	jsonBalanceReturn := e.ApiKeyRequest("GET", strRequest, nil)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonBalanceReturn
	}

	// log.Printf("jsonBalanceReturn: %v", jsonBalanceReturn)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s getBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonBalanceReturn)
		return operation.Error
	} else if jsonResponse.Code != "200000" {
		operation.Error = fmt.Errorf("%s getBalance Failed: %v", e.GetName(), jsonBalanceReturn)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountID); err != nil {
		operation.Error = fmt.Errorf("%s getBalance Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	for _, account := range accountID {
		if account.Type == accountType && account.Currency == mapParams["currency"] {
			frozen, err := strconv.ParseFloat(account.Holds, 64)
			avaliable, err := strconv.ParseFloat(account.Available, 64)
			if err != nil {
				return fmt.Errorf("%s balance parse fail: %v %+v", e.GetName(), err, account)
			}
			operation.BalanceFrozen = frozen
			operation.BalanceAvailable = avaliable
			return nil
		}
	}

	return fmt.Errorf("%s getBalance fail: %v", e.GetName(), jsonBalanceReturn)
}

func (e *Tagz) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalance{}
	strRequest := "/api/v1/accounts"

	mapParams := make(map[string]string)
	mapParams["type"] = "trade"

	jsonBalanceReturn := e.ApiKeyRequest("GET", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if jsonResponse.Code != "200000" {
		log.Printf("%s UpdateAllBalances Failed: %s %v", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
		return
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return
	}

	for _, balance := range accountBalance {
		c := e.GetCoinBySymbol(balance.Currency)
		if c != nil {
			freeamount, err := strconv.ParseFloat(balance.Available, 64)
			if err == nil {
				balanceMap.Set(c.Code, freeamount)
			}
		}
	}
}

/* Withdraw(coin *coin.Coin, quantity float64, addr, tag string) */
func (e *Tagz) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		log.Printf("Tagz API Key or Secret Key or passphrase are nil.")
		return false
	}

	// need to use inner transfer before withdraw
	// e.InnerTrans(quantity, coin, "trade", "main", fmt.Sprintf("%v", time.Now().UnixNano()/int64(time.Millisecond)))

	jsonResponse := JsonResponse{}
	withdraw := Withdraw{}
	strRequestUrl := "/api/v1/withdrawals"

	mapParams := make(map[string]string)
	mapParams["currency"] = e.GetSymbolByCoin(coin)
	mapParams["address"] = addr
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', -1, 64)

	jsonCreateWithdraw := e.ApiKeyRequest("POST", strRequestUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonCreateWithdraw), &jsonResponse); err != nil {
		log.Printf("%s Withdraw Json Unmarshal Err: %v %v", e.GetName(), err, jsonCreateWithdraw)
		return false
	} else if jsonResponse.Code != "200000" {
		log.Printf("%s Withdraw Failed: %s %v", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
		return false
	}

	if err := json.Unmarshal(jsonResponse.Data, &withdraw); err != nil {
		log.Printf("%s Withdraw Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return false
	}

	log.Printf("the withdraw state response %v and the withdraw id: %v", jsonCreateWithdraw, withdraw.WithdrawalID)
	return true
}

func (e *Tagz) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := OrderDetail{}
	strRequest := "/api/v1/orders"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["clientOid"] = fmt.Sprintf("%v", time.Now().UnixNano()) //Unique order id selected by you to identify your order
	mapParams["side"] = "sell"
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["type"] = "limit"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["size"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != "200000" {
		return nil, fmt.Errorf("%s LimitSell Failed: %s %v", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
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

func (e *Tagz) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := OrderDetail{}
	strRequest := "/api/v1/orders"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["clientOid"] = fmt.Sprintf("%v", time.Now().UnixNano()) //Unique order id selected by you to identify your order
	mapParams["side"] = "buy"
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["type"] = "limit"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["size"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != "200000" {
		return nil, fmt.Errorf("%s LimitBuy Failed: %s %v", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
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

func (e *Tagz) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	orderStatus := OrderStatus{}
	strRequest := fmt.Sprintf("/api/v1/orders/%s", order.OrderID)

	jsonOrderStatus := e.ApiKeyRequest("GET", strRequest, nil)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if jsonResponse.Code != "200000" {
		return fmt.Errorf("%s OrderStatus Failed: %s %v", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	switch orderStatus.OpType {
	case "DEAL":
		dealSize, _ := strconv.ParseFloat(orderStatus.DealSize, 64)
		if dealSize == order.Quantity {
			order.Status = exchange.Filled
		} else if dealSize > 0 && dealSize < order.Quantity {
			order.Status = exchange.Partial
		} else if math.Abs(dealSize-0.0) < 0.00000000001 {
			order.Status = exchange.New
		} else {
			order.Status = exchange.Other
		}
	case "CANCEL":
		order.Status = exchange.Cancelled
	default:
		order.Status = exchange.Other
	}

	order.DealRate, _ = strconv.ParseFloat(orderStatus.DealFunds, 64)
	order.DealQuantity, _ = strconv.ParseFloat(orderStatus.DealSize, 64)

	return nil
}

func (e *Tagz) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Tagz) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	cancelOrder := CancelOrder{}
	strRequest := fmt.Sprintf("/api/v1/orders/%s", order.OrderID)

	jsonCancelOrder := e.ApiKeyRequest("DELETE", strRequest, nil)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if jsonResponse.Code != "200000" {
		return fmt.Errorf("%s CancelOrder Failed: %s %v", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
	}
	if err := json.Unmarshal(jsonResponse.Data, &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Tagz) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Tagz) ApiKeyRequest(strMethod, strRequestPath string, mapParams map[string]string) string {
	nonce := time.Now().UnixNano() / int64(time.Millisecond) //Millisecond无误
	strRequestUrl := API_URL + strRequestPath

	httpClient := &http.Client{}
	var err error
	request := &http.Request{}
	signature := fmt.Sprintf("%v", nonce) + strMethod + strRequestPath
	jsonParams := ""

	if strMethod == "GET" || strMethod == "DELETE" {
		if nil != mapParams {
			payload := exchange.Map2UrlQuery(mapParams)
			strRequestUrl = API_URL + strRequestPath + "?" + payload
			signature = signature + "?" + payload
		}
		request, err = http.NewRequest(strMethod, strRequestUrl, nil)
	} else {
		if nil != mapParams {
			bytesParams, _ := json.Marshal(mapParams)
			jsonParams = string(bytesParams)
			signature = signature + jsonParams
		}
		request, err = http.NewRequest(strMethod, strRequestUrl, strings.NewReader(jsonParams))
	}

	if nil != err {
		return err.Error()
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("KC-API-KEY", e.API_KEY)
	request.Header.Add("KC-API-SIGN", exchange.ComputeHmac256Base64(signature, e.API_SECRET))
	request.Header.Add("KC-API-TIMESTAMP", fmt.Sprintf("%v", nonce))
	request.Header.Add("KC-API-PASSPHRASE", e.Passphrase)

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
