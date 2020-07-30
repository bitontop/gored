package okex

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
	API_URL string = "https://www.okex.com"
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
func (e *Okex) GetCoinsData() error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
	}

	coinsData := CoinsData{}

	strRequestUrl := "/api/account/v3/currencies"
	jsonCurrencyReturn := e.ApiKeyRequest("GET", nil, strRequestUrl)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Result Unmarshal Err: %v %s", e.GetName(), err, []byte(jsonCurrencyReturn))
	}

	for _, data := range coinsData {
		c := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			c = coin.GetCoin(data.Currency) //data.Currency)
			if c == nil {
				c = &coin.Coin{}
				c.Code = data.Currency //data.Currency
				c.Name = data.Name
				coin.AddCoin(c)
			}

		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(data.Currency) //data.Currency)
		}

		if c != nil {
			coinConstraint := e.GetCoinConstraint(c)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       c.ID,
					Coin:         c,
					ExSymbol:     data.Currency,
					ChainType:    exchange.MAINNET,
					Confirmation: DEFAULT_CONFIRMATION,
				}
			} else {
				coinConstraint.ExSymbol = data.Currency
			}

			if data.CanDeposit == "1" {
				coinConstraint.Deposit = true
			} else {
				coinConstraint.Deposit = false
			}

			if data.CanWithdraw == "1" {
				coinConstraint.Withdraw = true
			} else {
				coinConstraint.Withdraw = false
			}

			e.SetCoinConstraint(coinConstraint)
		}
	}
	return e.WithdrawFee()
}

func (e *Okex) WithdrawFee() error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
	}

	withdrawFee := WithdrawFee{}
	strRequest := "/api/account/v3/withdrawal/fee"

	jsonWithdrawFee := e.ApiKeyRequest("GET", nil, strRequest)
	if err := json.Unmarshal([]byte(jsonWithdrawFee), &withdrawFee); err != nil {
		return fmt.Errorf("%s WithdrawFee Unmarshal Err: %v %v", e.GetName(), err, jsonWithdrawFee)
	}

	for _, data := range withdrawFee {
		c := e.GetCoinBySymbol(data.Currency)
		if c != nil {
			coinConstraint := e.GetCoinConstraint(c)
			if data.MinFee != "" {
				minFee, err := strconv.ParseFloat(data.MinFee, 64)
				if err != nil {
					return fmt.Errorf("%s minFee conver to float64 err: %v %+v", e.GetName(), err, data)
				} else {
					coinConstraint.TxFee = minFee
					coinConstraint.Listed = true
				}
			} else {
				coinConstraint.Listed = false
			}
		}

	}
	return nil
}

/* GetPairsData - Get Pairs Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func (e *Okex) GetPairsData() error {
	//jsonResponse := &JsonResponse{}
	pairsData := PairsData{}

	strRequestUrl := "/api/spot/v3/instruments"
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
			p = e.GetPairBySymbol(data.InstrumentID)
		}

		lotSize, err := strconv.ParseFloat(data.SizeIncrement, 64)
		if err != nil {
			return fmt.Errorf("%s Convert lotSize to Float64 Err: %v %v", e.GetName(), err, data.SizeIncrement)
		}
		priceFilter, err := strconv.ParseFloat(data.TickSize, 64)
		if err != nil {
			return fmt.Errorf("%s Convert lotSize to Float64 Err: %v %v", e.GetName(), err, data.TickSize)
		}
		minTrade, err := strconv.ParseFloat(data.MinSize, 64)
		if err != nil {
			return fmt.Errorf("%s Convert minTrade to Float64 Err: %v %v", e.GetName(), err, data.MinSize)
		}

		if p != nil {
			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{ // no minBaseQuantity
					PairID:           p.ID,
					Pair:             p,
					ExSymbol:         data.InstrumentID,
					MakerFee:         DEFAULT_MAKER_FEE,
					TakerFee:         DEFAULT_TAKER_FEE,
					LotSize:          lotSize,
					PriceFilter:      priceFilter,
					MinTradeQuantity: minTrade,
					Listed:           DEFAULT_LISTED,
				}
			} else {
				pairConstraint.ExSymbol = data.InstrumentID
				pairConstraint.LotSize = lotSize
				pairConstraint.PriceFilter = priceFilter
				pairConstraint.MinTradeQuantity = minTrade
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
func (e *Okex) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	strRequestUrl := fmt.Sprintf("/api/spot/v3/instruments/%s/book", symbol)
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
		var buydata exchange.Order
		buydata.Rate, _ = strconv.ParseFloat(bid[0], 64)
		buydata.Quantity, _ = strconv.ParseFloat(bid[1], 64)

		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Asks {
		var selldata exchange.Order
		selldata.Rate, _ = strconv.ParseFloat(ask[0], 64)
		selldata.Quantity, _ = strconv.ParseFloat(ask[1], 64)

		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

/*************** Private API ***************/

func (e *Okex) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		log.Printf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
		return
	}

	accountBalance := AccountBalances{}
	strRequest := "/api/spot/v3/accounts"

	jsonBalanceReturn := e.ApiKeyRequest("GET", nil, strRequest)
	// log.Printf("jsonBalanceReturn: %v", jsonBalanceReturn)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		errorJson := ErrorMsg{}
		if err := json.Unmarshal([]byte(jsonBalanceReturn), &errorJson); err != nil {
			log.Printf("%s UpdateAllBalances Err: Code: %v Msg: %v", e.GetName(), errorJson.Code, errorJson.Msg)
		} else {
			log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		}
		return
	} else if len(accountBalance) == 0 {
		log.Printf("%s UpdateAllBalances Failed: %v", e.GetName(), jsonBalanceReturn)
		return
	}

	for _, v := range accountBalance {
		c := e.GetCoinBySymbol(v.Currency)
		if c != nil {
			balanceAvailable, err := strconv.ParseFloat(v.Available, 64)
			if err != nil {
				log.Printf("%s available balance conver to float64 err : %v", e.GetName, err)
				balanceAvailable = 0.0
			}
			balanceMap.Set(c.Code, balanceAvailable)
		}
	}
}

func (e *Okex) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		log.Printf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
		return false
	}

	withdrawResponse := WithdrawResponse{}
	strRequest := "/api/account/v3/withdrawal"

	mapParams := make(map[string]interface{})
	mapParams["currency"] = e.GetSymbolByCoin(coin)
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["destination"] = "4"
	mapParams["to_address"] = addr
	mapParams["trade_pwd"] = e.TradePassword
	mapParams["fee"] = e.GetTxFee(coin)

	jsonSubmitWithdraw := e.ApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &withdrawResponse); err != nil {
		log.Printf("%s Withdraw Json Unmarshal Err: %v %v", e.GetName(), err, jsonSubmitWithdraw)
		return false
	} else if !withdrawResponse.Result {
		log.Printf("%s Withdraw Failed: %v %v", e.GetName(), withdrawResponse.Code, withdrawResponse.Message)
		return false
	}

	return true
}

func (e *Okex) Transfer(coin *coin.Coin, quantity float64, from, to int) bool {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		log.Printf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
		return false
	}

	// transfer := Transfer{}
	// strRequest := "/api/account/v3/transfer"

	// mapParams := make(map[string]interface{})
	// mapParams["currency"] = e.GetSymbolByCoin(coin)
	// mapParams["amount"] = quantity
	// mapParams["from"] = from
	// mapParams["to"] = to

	// jsonTransfer := e.ApiKeyRequest("POST", mapParams, strRequest)

	// if err := json.Unmarshal([]byte(jsonTransfer), &transfer); err != nil {
	// 	log.Printf("%s Transfer Unmarshal Err: %v %v", e.GetName, err, jsonTransfer)
	// 	return false
	// } else if !transfer.Result {
	// 	log.Printf("%s Transfer Failed: %v %v", e.GetName, transfer.Code, transfer.Message)
	// 	return false
	// }

	return true
}

func (e *Okex) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return nil, fmt.Errorf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
	}

	placeOrder := PlaceOrder{}
	strRequest := "/api/spot/v3/orders"

	mapParams := make(map[string]interface{})
	mapParams["side"] = "sell"
	mapParams["instrument_id"] = e.GetSymbolByPair(pair)
	mapParams["type"] = "limit"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["size"] = strconv.FormatFloat(quantity, 'f', -1, 64)

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if !placeOrder.Result {
		return nil, fmt.Errorf("%s LimitSell Failed: %v", e.GetName(), jsonPlaceReturn)
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

func (e *Okex) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return nil, fmt.Errorf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
	}

	placeOrder := PlaceOrder{}
	strRequest := "/api/spot/v3/orders"

	mapParams := make(map[string]interface{})
	mapParams["side"] = "buy"
	mapParams["instrument_id"] = e.GetSymbolByPair(pair)
	mapParams["type"] = "limit"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["size"] = strconv.FormatFloat(quantity, 'f', -1, 64)

	jsonPlaceReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if !placeOrder.Result {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v", e.GetName(), jsonPlaceReturn)
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

func (e *Okex) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
	}

	orderStatus := OrderStatus{}
	strRequest := fmt.Sprintf("/api/spot/v3/orders/%s", order.OrderID)

	mapParams := make(map[string]string)
	mapParams["instrument_id"] = e.GetSymbolByPair(order.Pair)

	strRequest += fmt.Sprintf("?%s", exchange.Map2UrlQuery(mapParams))

	jsonOrderStatus := e.ApiKeyRequest("GET", nil, strRequest)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	}

	order.StatusMessage = jsonOrderStatus
	if orderStatus.Status == "open" {
		order.Status = exchange.New
	} else if orderStatus.Status == "part_filled" {
		order.Status = exchange.Partial
	} else if orderStatus.Status == "filled" {
		order.Status = exchange.Filled
	} else if orderStatus.Status == "canceling" {
		order.Status = exchange.Canceling
	} else if orderStatus.Status == "cancelled" {
		order.Status = exchange.Cancelled
	} else if orderStatus.Status == "failure" {
		order.Status = exchange.Rejected
	} else {
		order.Status = exchange.Other
	}

	return nil
}

func (e *Okex) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Okex) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
	}

	cancelOrder := PlaceOrder{}
	strRequest := fmt.Sprintf("/api/spot/v3/cancel_orders/%s", order.OrderID)

	mapParams := make(map[string]interface{})
	mapParams["instrument_id"] = e.GetSymbolByPair(order.Pair)

	jsonCancelOrder := e.ApiKeyRequest("POST", mapParams, strRequest)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if !cancelOrder.Result {
		return fmt.Errorf("%s CancelOrder Failed: %v", e.GetName(), jsonCancelOrder)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Okex) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Okex) ApiKeyRequest(method string, mapParams map[string]interface{}, strRequestPath string) string {
	TimeStamp := IsoTime()

	jsonParams := ""
	var bytesParams []byte
	if len(mapParams) != 0 {
		bytesParams, _ = json.Marshal(mapParams)
		jsonParams = string(bytesParams)
	}

	var strMessage string
	if method == "GET" {
		if len(mapParams) > 0 {
			strRequestPath += "?" + exchange.Map2UrlQueryInterface(mapParams)
		}
		strMessage = TimeStamp + method + strRequestPath
		if len(mapParams) > 0 {
			strMessage += jsonParams
		}
	} else {
		strMessage = TimeStamp + method + strRequestPath + jsonParams
	}

	// log.Printf("===================strMessage: %v", strMessage)

	signature := exchange.ComputeHmac256Base64(strMessage, e.API_SECRET)
	strUrl := API_URL + strRequestPath

	httpClient := &http.Client{}
	request, err := http.NewRequest(method, strUrl, bytes.NewReader(bytesParams))
	if nil != err {
		return err.Error()
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json; charset=UTF-8")
	request.Header.Add("OK-ACCESS-KEY", e.API_KEY)
	request.Header.Add("OK-ACCESS-SIGN", signature)
	request.Header.Add("OK-ACCESS-TIMESTAMP", TimeStamp)
	request.Header.Add("OK-ACCESS-PASSPHRASE", e.Passphrase)

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

func IsoTime() string {
	utcTime := time.Now().UTC()
	iso := utcTime.String()
	isoBytes := []byte(iso)
	iso = string(isoBytes[:10]) + "T" + string(isoBytes[11:23]) + "Z"
	return iso
}
