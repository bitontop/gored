package gemini

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"

	"strconv"
)

const (
	API_URL string = "https://api.sandbox.gemini.com"
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
func (e *Gemini) GetCoinsData() error {
	file, err := ioutil.ReadFile("../exchange/gemini/constraint.json")
	if err != nil {
		log.Printf("%s getCoin read file err: %v", e.GetName, err)
	}

	pairsFile := PairsFile{}
	if err := json.Unmarshal(file, &pairsFile); err != nil {
		log.Printf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, file)
	}

	for _, data := range pairsFile {
		base := &coin.Coin{}
		target := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base = coin.GetCoin(data.Base)
			if base == nil {
				base = &coin.Coin{}
				base.Code = data.Base
				coin.AddCoin(base)
			}
			target = coin.GetCoin(data.Quote)
			if target == nil {
				target = &coin.Coin{}
				target.Code = data.Quote
				coin.AddCoin(target)
			}
		case exchange.JSON_FILE:
			base = e.GetCoinBySymbol(data.Base)
			target = e.GetCoinBySymbol(data.Quote)
		}

		if base != nil {
			coinConstraint := e.GetCoinConstraint(base)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       base.ID,
					Coin:         base,
					ExSymbol:     data.Base,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = data.Base
			}
			e.SetCoinConstraint(coinConstraint)
		}

		if target != nil {
			coinConstraint := e.GetCoinConstraint(target)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       target.ID,
					Coin:         target,
					ExSymbol:     data.Quote,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = data.Quote
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
func (e *Gemini) GetPairsData() error {
	file, err := ioutil.ReadFile("../exchange/gemini/constraint.json")
	if err != nil {
		log.Printf("%s getPair read file err: %v", e.GetName, err)
	}

	pairsFile := PairsFile{}
	if err := json.Unmarshal(file, &pairsFile); err != nil {
		log.Printf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, file)
	}

	for _, data := range pairsFile {
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(data.Base)
			target := coin.GetCoin(data.Quote)
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
					MakerFee:    DEFAULT_MAKER_FEE,
					TakerFee:    DEFAULT_TAKER_FEE,
					LotSize:     data.MinOrderIncre,
					PriceFilter: data.MinPriceIncre,
					Listed:      DEFAULT_LISTED,
				}
			} else {
				pairConstraint.ExSymbol = data.Symbol
				pairConstraint.LotSize = data.MinOrderIncre
				pairConstraint.PriceFilter = data.MinPriceIncre
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
func (e *Gemini) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	errResponse := ErrorResponse{}
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	strRequestUrl := fmt.Sprintf("/v1/book/%s", symbol)
	strUrl := API_URL + strRequestUrl

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonOrderbook), &errResponse); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if errResponse.Result == "error" {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %v", e.GetName(), errResponse)
	}

	if err := json.Unmarshal([]byte(jsonOrderbook), &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	var err error
	for _, bid := range orderBook.Bids {
		buydata := exchange.Order{}
		buydata.Quantity, err = strconv.ParseFloat(bid.Amount, 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
		}

		buydata.Rate, err = strconv.ParseFloat(bid.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Rate error:%v", e.GetName(), err)
		}
		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Asks {
		selldata := exchange.Order{}
		selldata.Quantity, err = strconv.ParseFloat(ask.Amount, 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
		}

		selldata.Rate, err = strconv.ParseFloat(ask.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Rate error:%v", e.GetName(), err)
		}
		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

func (e *Gemini) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *Gemini) DoAccountOperation(operation *exchange.AccountOperation) error {
	return nil
}

func (e *Gemini) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}
	accountBalance := AccountBalances{}
	strRequest := "/v1/balances"

	mapParams := make(map[string]interface{})
	mapParams["request"] = strRequest

	jsonBalanceReturn := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Result Unmarshal Err: %v %s", e.GetName(), err, jsonBalanceReturn)
		return
	}

	for _, v := range accountBalance {
		c := e.GetCoinBySymbol(v.Currency)
		if c != nil {
			freeAmount, err := strconv.ParseFloat(v.Amount, 64)
			if err != nil {
				log.Printf("%s balance parse Err: %v %v", e.GetName(), err, v.Amount)
				return
			}
			balanceMap.Set(c.Code, freeAmount)
		}
	}
}

/*
txHash is Only shown for ETH and GUSD withdrawals.
withdrawalID and message are Only shown for BTC, ZEC, LTC and BCH withdrawals.
*/
func (e *Gemini) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil", e.GetName())
		return false
	}
	withdrawal := Withdrawal{}
	strRequest := "/v1/withdraw" + "/" + strings.ToLower(coin.Code)

	mapParams := make(map[string]interface{})
	mapParams["request"] = strRequest
	mapParams["address"] = addr
	mapParams["ammount"] = strconv.FormatFloat(quantity, 'f', -1, 64)

	jsonSubmitWithdraw := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &withdrawal); err != nil {
		log.Printf("%s Withdraw Result Unmarshal Err: %v %s", e.GetName(), err, jsonSubmitWithdraw)
		return false
	}
	return true
}

func (e *Gemini) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}
	sellorder := PlaceOrder{}
	strRequest := "/v1/order/new"

	mapParams := make(map[string]interface{})
	mapParams["request"] = strRequest
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["side"] = "sell"
	mapParams["type"] = "exchange limit"

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &sellorder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Result Unmarshal Err: %v %s", e.GetName(), err, jsonPlaceReturn)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      sellorder.OrderID,
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Gemini) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}
	buyorder := PlaceOrder{}
	strRequest := "/v1/order/new"

	mapParams := make(map[string]interface{})
	mapParams["request"] = strRequest
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["side"] = "buy"
	mapParams["type"] = "exchange limit"

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &buyorder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Result Unmarshal Err: %v %s", e.GetName(), err, jsonPlaceReturn)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      buyorder.OrderID,
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}
	return order, nil
}

func (e *Gemini) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}
	orderStatus := PlaceOrder{}
	strRequest := "/v1/order/status"

	id, _ := strconv.ParseInt(order.OrderID, 0, 0)

	mapParams := make(map[string]interface{})
	mapParams["request"] = strRequest
	mapParams["order_id"] = id

	jsonOrderStatus := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Result Unmarshal Err: %v %s", e.GetName(), err, jsonOrderStatus)
	}

	order.StatusMessage = jsonOrderStatus
	if orderStatus.IsCancelled {
		order.Status = exchange.Cancelled
	} else if !orderStatus.IsLive && orderStatus.RemainingAmount != "0" {
		order.Status = exchange.Canceling
	} else if orderStatus.RemainingAmount == "0" {
		order.Status = exchange.Filled
	} else if orderStatus.IsLive && orderStatus.ExecutedAmount != "0" {
		order.Status = exchange.Partial
	} else {
		order.Status = exchange.New
	}

	return nil
}

func (e *Gemini) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Gemini) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}
	cancelOrder := PlaceOrder{}
	strRequest := "/v1/order/cancel"

	mapParams := make(map[string]interface{})
	mapParams["request"] = strRequest
	mapParams["order_id"] = order.OrderID

	jsonCancelOrder := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Result Unmarshal Err: %v %s", e.GetName(), err, jsonCancelOrder)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Gemini) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Gemini) ApiKeyRequest(strMethod string, strRequestPath string, mapParams map[string]interface{}) string {
	mapParams["nonce"] = /* fmt.Sprintf("%d",  */ time.Now().UnixNano() //)

	strUrl := API_URL + strRequestPath

	var bytesParams []byte
	if nil != mapParams {
		bytesParams, _ = json.Marshal(mapParams)
	}

	b64 := base64.StdEncoding.EncodeToString(bytesParams)
	signature := ComputeHmac384NoDecode(b64, e.API_SECRET)

	request, err := http.NewRequest("POST", strUrl, bytes.NewBuffer(bytesParams))
	if nil != err {
		return err.Error()
	}

	request.Header.Add("Content-Length", "0")
	request.Header.Add("Content-Type", "text/plain")
	request.Header.Add("X-GEMINI-APIKEY", e.API_KEY)
	request.Header.Add("X-GEMINI-PAYLOAD", b64)
	request.Header.Add("X-GEMINI-SIGNATURE", signature)
	request.Header.Add("Cache-Control", "no-cache")

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

func ComputeHmac384NoDecode(strMessage string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha512.New384, key)
	h.Write([]byte(strMessage))

	return hex.EncodeToString(h.Sum(nil))
}
