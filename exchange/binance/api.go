package binance

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"bytes"
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

/*The Base Endpoint URL*/
const (
	API_URL      = "https://api.binance.com"
	CONTRACT_URL = "https://fapi.binance.com"
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
func (e *Binance) GetCoinsData() error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s Get Coins API Key or Secret Key are nil. Using backupGetCoin", e.GetName())
		return e.backupGetCoinsData()
	}

	coinsData := CoinsData{}
	strUrl := "/sapi/v1/capital/config/getall"

	mapParams := make(map[string]string)

	jsonCurrencyReturn := e.ApiKeyGet(mapParams, strUrl)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	}

	for _, data := range coinsData {
		c := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			c = coin.GetCoin(data.Coin)
			if c == nil {
				c = &coin.Coin{}
				c.Code = data.Coin
				c.Name = data.Name
				coin.AddCoin(c)
			}
		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(data.Coin)
		}

		if c != nil {
			confirmation := DEFAULT_CONFIRMATION
			txFee := DEFAULT_TXFEE
			deposit, withdraw := DEFAULT_DEPOSIT, DEFAULT_WITHDRAW
			for _, netWork := range data.NetworkList {
				if netWork.IsDefault != true {
					continue
				}
				confirmation = netWork.MinConfirm
				tempTX, err := strconv.ParseFloat(netWork.WithdrawFee, 64)
				if err == nil {
					txFee = tempTX
				}
				deposit = netWork.DepositEnable
				withdraw = netWork.WithdrawEnable
			}

			coinConstraint := e.GetCoinConstraint(c)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       c.ID,
					Coin:         c,
					ExSymbol:     data.Coin,
					ChainType:    exchange.MAINNET,
					TxFee:        txFee,
					Withdraw:     withdraw,
					Deposit:      deposit,
					Confirmation: confirmation,
					Listed:       true,
				}
			} else {
				coinConstraint.ExSymbol = data.Coin
				coinConstraint.TxFee = txFee
				coinConstraint.Withdraw = withdraw
				coinConstraint.Deposit = deposit
				coinConstraint.Confirmation = confirmation
			}

			e.SetCoinConstraint(coinConstraint)
		}
	}
	return nil
}

func (e *Binance) backupGetCoinsData() error {
	pairsData := PairsData{}
	strRequestUrl := "/api/v1/exchangeInfo"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
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
					Listed:       true,
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
					Listed:       true,
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
func (e *Binance) GetPairsData() error {
	pairsData := &PairsData{}

	strRequestUrl := "/api/v1/exchangeInfo"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	}

	for _, data := range pairsData.Symbols {
		if data.Status == "TRADING" {
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
				var err error
				lotsize := 0.0
				priceFilter := 0.0
				minTrade := 0.0
				minBase := 0.0
				for _, filter := range data.Filters {
					switch filter.FilterType {
					case "LOT_SIZE":
						lotsize, err = strconv.ParseFloat(filter.StepSize, 64)
						if err != nil {
							log.Printf("%s Lot Size Err: %v", e.GetName(), err)
							lotsize = DEFAULT_LOT_SIZE
						}
						minTrade, err = strconv.ParseFloat(filter.MinQty, 64)
						if err != nil {
							log.Printf("%s minTrade Filter parse Err: %v, %v", e.GetName(), filter.MinQty, err)
							minTrade = 0.0
						}
					case "PRICE_FILTER":
						priceFilter, err = strconv.ParseFloat(filter.TickSize, 64)
						if err != nil {
							log.Printf("%s Price Filter Err: %v", e.GetName(), err)
							priceFilter = DEFAULT_PRICE_FILTER
						}
					case "MIN_NOTIONAL":
						minBase, err = strconv.ParseFloat(filter.MinNotional, 64)
						if err != nil {
							log.Printf("%s minBase Filter parse Err: %v, %v", e.GetName(), filter.MinNotional, err)
							minBase = 0.0
						}
					}
				}
				pairConstraint := e.GetPairConstraint(p)

				if pairConstraint == nil {
					pairConstraint = &exchange.PairConstraint{
						PairID:               p.ID,
						Pair:                 p,
						ExSymbol:             data.Symbol,
						MakerFee:             DEFAULT_MAKER_FEE,
						TakerFee:             DEFAULT_TAKER_FEE,
						LotSize:              lotsize,
						PriceFilter:          priceFilter,
						MinTradeQuantity:     minTrade,
						MinTradeBaseQuantity: minBase,
						Listed:               true,
					}
				} else {
					pairConstraint.ExSymbol = data.Symbol
					pairConstraint.LotSize = lotsize
					pairConstraint.PriceFilter = priceFilter
					pairConstraint.MinTradeQuantity = minTrade
					pairConstraint.MinTradeBaseQuantity = minBase
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
func (e *Binance) OrderBook(p *pair.Pair) (*exchange.Maker, error) {
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(p)

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol
	mapParams["limit"] = "100"

	strRequestUrl := "/api/v1/depth"
	strUrl := API_URL + strRequestUrl

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &orderBook); err != nil {
		return nil, fmt.Errorf("%s OrderBook json Unmarshal error: %v %v", e.GetName(), err, jsonOrderbook)
	} else if orderBook.Code != 0 {
		return nil, fmt.Errorf("%s Get OrderBook failed: %v", e.GetName(), jsonOrderbook)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	maker.LastUpdateID = int64(orderBook.LastUpdateID)

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

	return maker, err
}

/*************** Private API ***************/

func (e *Binance) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	accountBalance := AccountBalances{}
	strRequest := "/api/v3/account"

	jsonBalanceReturn := e.ApiKeyGet(make(map[string]string), strRequest)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances json Unmarshal error: %v %s", e.GetName(), err, jsonBalanceReturn)
		return
	} else {
		for _, balance := range accountBalance.Balances {
			freeamount, err := strconv.ParseFloat(balance.Free, 64)
			if err != nil {
				log.Printf("%s UpdateAllBalances err: %+v %v", e.GetName(), balance, err)
				return
			} else {
				c := e.GetCoinBySymbol(balance.Asset)
				if c != nil {
					balanceMap.Set(c.Code, freeamount)
				}
			}
		}
	}
}

/* Withdraw(coin *coin.Coin, quantity float64, addr, tag string) */
func (e *Binance) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return false
	}

	withdraw := WithdrawResponse{}
	strRequest := "/wapi/v3/withdraw.html"

	mapParams := make(map[string]string)
	mapParams["asset"] = e.GetSymbolByCoin(coin)
	mapParams["address"] = addr
	if tag != "" { //this part is not working yet
		mapParams["addressTag"] = tag
	}
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["timestamp"] = fmt.Sprintf("%d", time.Now().UnixNano()/1e6)

	jsonSubmitWithdraw := e.WApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &withdraw); err != nil {
		log.Printf("%s Withdraw Json Unmarshal Error: %v %v", e.GetName(), err, jsonSubmitWithdraw)
		return false
	}
	if !withdraw.Success {
		log.Printf("%s Withdraw Failed: %s", e.GetName(), withdraw.Msg)
		return false
	}

	return true
}

func (e *Binance) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	placeOrder := PlaceOrder{}
	strRequest := "/api/v3/order"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "SELL"
	mapParams["type"] = "LIMIT"
	mapParams["timeInForce"] = "GTC"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["quantity"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)

	order := &exchange.Order{
		Pair: pair,
		// OrderID:      fmt.Sprintf("%d", placeOrder.OrderID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		err = fmt.Errorf("%s LimitSell Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
		order.Error = err
		return nil, err
	} else if placeOrder.Code != 0 {
		err = fmt.Errorf("%s LimitSell failed:%v Message:%v", e.GetName(), placeOrder.Code, placeOrder.Msg)
		order.Error = err
		return nil, err
	}

	order.OrderID = fmt.Sprintf("%d", placeOrder.OrderID)

	return order, nil
}

func (e *Binance) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	placeOrder := PlaceOrder{}
	strRequest := "/api/v3/order"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "BUY"
	mapParams["type"] = "LIMIT"
	mapParams["timeInForce"] = "GTC"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["quantity"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)

	order := &exchange.Order{
		Pair: pair,
		// OrderID:      fmt.Sprintf("%d", placeOrder.OrderID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		err = fmt.Errorf("%s LimitBuy Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
		order.Error = err
		return nil, err
	} else if placeOrder.Code != 0 {
		err = fmt.Errorf("%s LimitBuy failed:%v Message:%v", e.GetName(), placeOrder.Code, placeOrder.Msg)
		order.Error = err
		return nil, err
	}

	order.OrderID = fmt.Sprintf("%d", placeOrder.OrderID)

	return order, nil
}

func (e *Binance) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	orderStatus := PlaceOrder{}
	strRequest := "/api/v3/order"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(order.Pair)
	mapParams["orderId"] = order.OrderID

	jsonOrderStatus := e.ApiKeyGet(mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if orderStatus.Code != 0 {
		return fmt.Errorf("%s Get OrderStatus Error: %v %s", e.GetName(), orderStatus.Code, orderStatus.Msg)
	}

	if orderStatus.Status == "CANCELED" {
		order.Status = exchange.Cancelled
	} else if orderStatus.Status == "FILLED" {
		order.Status = exchange.Filled
	} else if orderStatus.Status == "PARTIALLY_FILLED" {
		order.Status = exchange.Partial
	} else if orderStatus.Status == "REJECTED" {
		order.Status = exchange.Rejected
	} else if orderStatus.Status == "Expired" {
		order.Status = exchange.Expired
	} else if orderStatus.Status == "NEW" {
		order.Status = exchange.New
	} else {
		order.Status = exchange.Other
	}

	order.DealRate, _ = strconv.ParseFloat(orderStatus.Price, 64)
	order.DealQuantity, _ = strconv.ParseFloat(orderStatus.ExecutedQty, 64)

	return nil
}

func (e *Binance) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Binance) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	cancelOrder := PlaceOrder{}
	strRequest := "/api/v3/order"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(order.Pair)
	mapParams["orderId"] = order.OrderID

	jsonCancelOrder := e.ApiKeyRequest("DELETE", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if cancelOrder.Code != 0 {
		return fmt.Errorf("%s CancelOrder Error: %v %s", e.GetName(), cancelOrder.Code, cancelOrder.Msg)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Binance) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Binance) ApiKeyGet(mapParams map[string]string, strRequestPath string) string {
	mapParams["recvWindow"] = "50000" //"50000000"
	mapParams["timestamp"] = fmt.Sprintf("%d", time.Now().UTC().UnixNano()/int64(time.Millisecond))
	mapParams["signature"] = exchange.ComputeHmac256NoDecode(exchange.Map2UrlQuery(mapParams), e.API_SECRET)

	payload := exchange.Map2UrlQuery(mapParams)
	strUrl := API_URL + strRequestPath + "?" + payload

	request, err := http.NewRequest("GET", strUrl, nil)
	if nil != err {
		return err.Error()
	}
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	request.Header.Add("X-MBX-APIKEY", e.API_KEY)

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
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Binance) ApiKeyRequest(strMethod string, mapParams map[string]string, strRequestPath string) string {
	mapParams["recvWindow"] = "50000" //"50000000"
	mapParams["timestamp"] = fmt.Sprintf("%d", time.Now().UTC().UnixNano()/int64(time.Millisecond))
	mapParams["signature"] = exchange.ComputeHmac256NoDecode(exchange.Map2UrlQuery(mapParams), e.API_SECRET)

	payload := exchange.Map2UrlQuery(mapParams)
	strUrl := API_URL + strRequestPath

	request, err := http.NewRequest(strMethod, strUrl, bytes.NewBuffer([]byte(payload)))
	if nil != err {
		return err.Error()
	}
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	request.Header.Add("X-MBX-APIKEY", e.API_KEY)

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
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Binance) WApiKeyRequest(strMethod string, mapParams map[string]string, strRequestPath string) string {
	mapParams["recvWindow"] = "50000" //"50000000"
	mapParams["timestamp"] = fmt.Sprintf("%d", time.Now().UTC().UnixNano()/int64(time.Millisecond))
	signature := exchange.ComputeHmac256NoDecode(exchange.Map2UrlQuery(mapParams), e.API_SECRET)

	// payload := fmt.Sprintf("%s&signature=%s", exchange.Map2UrlQuery(mapParams), signature)
	// if len(mapParams) == 0 {
	// 	payload = fmt.Sprintf("signature=%s", signature)
	// }

	// strUrl := API_URL + strRequestPath
	// request, err := http.NewRequest(strMethod, strUrl, bytes.NewBuffer([]byte(payload)))
	// mapParams["signature"] = signature

	// signature at end
	strUrl := API_URL + strRequestPath + "?" + fmt.Sprintf("%s&signature=%s", exchange.Map2UrlQuery(mapParams), signature)
	request, err := http.NewRequest(strMethod, strUrl, nil)

	if nil != err {
		return err.Error()
	}

	// log.Printf("=wwww= strUrl: %v", strUrl) // wwwwwwwwwwwwwwwwwwwwwwwwwwwwwww

	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	request.Header.Add("X-MBX-APIKEY", e.API_KEY)

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

func (e *Binance) TestApi() string {
	key := "NWFjY2U2OWQzMjM0NDczMTgwODcwMDNjOWU2M2NhMDE="
	secret := "NWFmNWQ0NGRiMTlhNGMz"
	strUrl := "https://api.btse.com/spot/api/v3.1/user/wallet"
	nonce := fmt.Sprintf("%d", time.Now().UTC().UnixNano()/int64(time.Millisecond))

	mapParams := make(map[string]string)

	log.Printf("====================\nkey: %v, secret: %v, strUrl: %v, nonce: %v, mapParams: %+v", key, secret, strUrl, nonce, mapParams) // ====================

	// signature := exchange.ComputeHmac256NoDecode(exchange.Map2UrlQuery(mapParams), e.API_SECRET)

	// payload := fmt.Sprintf("%s&signature=%s", exchange.Map2UrlQuery(mapParams), signature)
	// if len(mapParams) == 0 {
	// 	payload = fmt.Sprintf("signature=%s", signature)
	// }

	// strUrl := API_URL + strRequestPath
	// request, err := http.NewRequest(strMethod, strUrl, bytes.NewBuffer([]byte(payload)))
	// mapParams["signature"] = signature

	// signature at end
	// strUrl := API_URL + strRequestPath + "?" + fmt.Sprintf("%s&signature=%s", exchange.Map2UrlQuery(mapParams), signature)
	request, err := http.NewRequest("GET", strUrl, nil)

	if nil != err {
		return err.Error()
	}

	// log.Printf("=wwww= strUrl: %v", strUrl) // wwwwwwwwwwwwwwwwwwwwwwwwwwwwwww

	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	request.Header.Add("X-MBX-APIKEY", e.API_KEY)

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
