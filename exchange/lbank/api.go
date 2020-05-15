package lbank

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
)

const (
	API_URL string = "https://www.lbkex.net" //"https://api.lbkex.com"
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
func (e *Lbank) GetCoinsData() error {
	coinsData := CoinsData{}

	strRequestUrl := "/v2/withdrawConfigs.do" //"/v1/withdrawConfigs.do"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if coinsData.Result != "true" {
		return fmt.Errorf("%s Get Coin Failed: %s", e.GetName(), jsonCurrencyReturn)
	}

	for _, data := range coinsData.Data {
		c := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			c = coin.GetCoin(data.AssetCode)
			if c == nil {
				c = &coin.Coin{}
				c.Code = data.AssetCode
				coin.AddCoin(c)
			}
		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(data.AssetCode)
		}

		if c != nil {
			txFee, _ := strconv.ParseFloat(data.Fee, 64)
			coinConstraint := e.GetCoinConstraint(c)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       c.ID,
					Coin:         c,
					ExSymbol:     data.AssetCode,
					ChainType:    exchange.MAINNET,
					TxFee:        txFee,
					Withdraw:     data.CanWithDraw,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = data.AssetCode
				coinConstraint.TxFee = txFee
				coinConstraint.Withdraw = data.CanWithDraw
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
func (e *Lbank) GetPairsData() error {
	pairsData := PairsData{}

	strRequestUrl := "/v1/accuracy.do"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	}

	for _, data := range pairsData {
		coinStr := strings.Split(data.Symbol, "_")
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(coinStr[1])
			target := coin.GetCoin(coinStr[0])
			if base != nil && target != nil {
				p = pair.GetPair(base, target)
			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(data.Symbol)
		}
		if p != nil {
			lotsize, _ := strconv.Atoi(data.QuantityAccuracy)
			ticksize, _ := strconv.Atoi(data.PriceAccuracy)
			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{
					PairID:      p.ID,
					Pair:        p,
					ExSymbol:    data.Symbol,
					MakerFee:    DEFAULT_MAKER_FEE,
					TakerFee:    DEFAULT_TAKER_FEE,
					LotSize:     math.Pow10(lotsize * -1),
					PriceFilter: math.Pow10(ticksize * -1),
					Listed:      DEFAULT_LISTED,
				}
			} else {
				pairConstraint.ExSymbol = data.Symbol
				pairConstraint.LotSize = math.Pow10(lotsize * -1)
				pairConstraint.PriceFilter = math.Pow10(ticksize * -1)
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
func (e *Lbank) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	strRequestUrl := "/v1/depth.do"
	strUrl := API_URL + strRequestUrl

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol
	mapParams["size"] = "60"
	mapParams["merge"] = "0"

	maker := &exchange.Maker{}
	maker.WorkerIP = exchange.GetExternalIP()

	maker.BeforeTimestamp = float64(time.Now().UnixNano() / 1e6)

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	for _, bid := range orderBook.Bids {
		var buydata exchange.Order

		//Modify according to type and structure
		buydata.Rate = bid[0]
		buydata.Quantity = bid[1]
		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Asks {
		var selldata exchange.Order

		//Modify according to type and structure
		selldata.Rate = ask[0]
		selldata.Quantity = ask[1]
		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

func (e *Lbank) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *Lbank) DoAccountOperation(operation *exchange.AccountOperation) error {
	switch operation.Type {

	// case exchange.Transfer:
	// 	return e.transfer(operation)
	// case exchange.BalanceList:
	// 	return e.getAllBalance(operation)
	// case exchange.Balance:
	// 	return e.getBalance(operation)

	case exchange.Withdraw:
		return e.doWithdraw(operation)

	}
	return fmt.Errorf("%s Operation type invalid: %s %v", operation.Ex, operation.Wallet, operation.Type)
}

func (e *Lbank) doWithdraw(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	withdraw := Withdraw{}
	strRequest := "/v1/withdraw.do"

	mapParams := make(map[string]string)

	mapParams["account"] = operation.WithdrawAddress
	mapParams["assetCode"] = e.GetSymbolByCoin(operation.Coin)
	mapParams["amount"] = operation.WithdrawAmount
	if operation.WithdrawTag != "" {
		mapParams["memo"] = operation.WithdrawTag
	}

	jsonWithdrawReturn := e.ApiKeyPost(strRequest, make(map[string]string))
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonWithdrawReturn
	}

	if err := json.Unmarshal([]byte(jsonWithdrawReturn), &withdraw); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Json Unmarshal Err: %v", e.GetName(), err)
		return operation.Error
	} else if withdraw.Result != "true" {
		operation.Error = fmt.Errorf("%s Withdraw Failed: %v", e.GetName(), jsonWithdrawReturn)
		return operation.Error
	}

	operation.WithdrawID = fmt.Sprintf("%v", withdraw.WithdrawID)

	return nil
}

func (e *Lbank) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	accountBalance := AccountBalances{}
	strRequest := "/v1/user_info.do"

	jsonBalanceReturn := e.ApiKeyPost(strRequest, make(map[string]string))
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if accountBalance.Result != "true" {
		log.Printf("%s UpdateAllBalances Failed: %v", e.GetName(), jsonBalanceReturn)
		return
	}

	for key, value := range accountBalance.Info.Free {
		c := e.GetCoinBySymbol(key)
		freeamount, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Printf("%s balance parse error: %v, %v", e.GetName(), err, value)
			return
		}
		if c != nil {
			balanceMap.Set(c.Code, freeamount)
		}
	}
}

func (e *Lbank) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return false
	}

	withdraw := Withdraw{}
	strRequest := "/v1/withdraw.do"

	mapParams := make(map[string]string)

	mapParams["account"] = addr
	mapParams["assetCode"] = e.GetSymbolByCoin(coin)
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	if tag != "" {
		mapParams["memo"] = tag
	}

	jsonWithdrawReturn := e.ApiKeyPost(strRequest, make(map[string]string))
	if err := json.Unmarshal([]byte(jsonWithdrawReturn), &withdraw); err != nil {
		log.Printf("%s Withdraw Json Unmarshal Err: %v %v", e.GetName(), err, jsonWithdrawReturn)
		return false
	} else if withdraw.Result != "true" {
		log.Printf("%s Withdraw Failed: %v", e.GetName(), jsonWithdrawReturn)
		return false
	}

	return true
}

func (e *Lbank) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	placeOrder := PlaceOrder{}
	strRequest := "/v1/create_order.do"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["type"] = "sell"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', -1, 64)

	jsonPlaceReturn := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if placeOrder.Result != "true" {
		return nil, fmt.Errorf("%s LimitSell Failed: %v, %v", e.GetName(), placeOrder.ErrorCode, placeOrder.Result)
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

func (e *Lbank) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	placeOrder := PlaceOrder{}
	strRequest := "/v1/create_order.do"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(pair)
	mapParams["type"] = "buy"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', -1, 64)
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', -1, 64)

	jsonPlaceReturn := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if placeOrder.Result != "true" {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v, %v", e.GetName(), placeOrder.ErrorCode, placeOrder.Result)
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

func (e *Lbank) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	orderStatus := OrderStatus{}
	strRequest := "/v1/orders_info.do"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(order.Pair)
	mapParams["order_id"] = order.OrderID

	jsonOrderStatus := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if orderStatus.Result != "true" {
		return fmt.Errorf("%s OrderStatus Failed: %v, %v", e.GetName(), orderStatus.ErrorCode, orderStatus.Result)
	}

	order.StatusMessage = jsonOrderStatus
	for _, v := range orderStatus.Orders {
		if v.OrderID == order.OrderID {
			if v.Status == -1 {
				order.Status = exchange.Cancelled
			} else if v.Status == 0 {
				order.Status = exchange.New
			} else if v.Status == 1 {
				order.Status = exchange.Partial
			} else if v.Status == 2 {
				order.Status = exchange.Filled
			} else if v.Status == 4 {
				order.Status = exchange.Canceling
			} else {
				order.Status = exchange.Other
			}
		}
	}

	return nil
}

func (e *Lbank) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Lbank) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	cancelOrder := CancelOrders{}
	strRequest := "/v1/cancel_order.do"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(order.Pair)
	mapParams["order_id"] = order.OrderID

	jsonCancelOrder := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if cancelOrder.Result != "true" {
		return fmt.Errorf("%s CancelOrder Failed: %v", e.GetName(), cancelOrder.ErrorCode, cancelOrder.Result)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Lbank) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Lbank) ApiKeyPost(strRequestPath string, mapParams map[string]string) string {
	//Signature Request Params
	mapParams["api_key"] = e.API_KEY
	mapParams["sign"] = ComputeMD5(mapParams, e.API_SECRET)

	httpClient := &http.Client{}
	payload := exchange.Map2UrlQuery(mapParams)
	strUrl := fmt.Sprintf("%s%s?%s", API_URL, strRequestPath, payload)

	// Parameters do not post on body.
	request, err := http.NewRequest("POST", strUrl, nil)
	if nil != err {
		return err.Error()
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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

//Signature加密
func ComputeMD5(mapParams map[string]string, strSecret string) string {
	var secret string

	/* The Secret Key is PKCS8 RSA Private Key.
	   Header and Footer are necessary.
	   Each Line cannot be over 64 bits.*/
	if len(strSecret) > 32 {
		ch := strings.Split(strSecret, "\n")
		if ch[0] != "-----BEGIN RSA PRIVATE KEY-----" {
			secret = "-----BEGIN RSA PRIVATE KEY-----\n"
			for i := 0; i <= len(strSecret); i = i + 64 {
				if len(strSecret) > i+63 {
					secret += fmt.Sprintf("%s\n", strSecret[i:i+64])
				} else {
					secret += fmt.Sprintf("%s\n", strSecret[i:len(strSecret)])
				}
			}
			secret += "-----END RSA PRIVATE KEY-----"
		}
	}
	// log.Printf("%s", secret)

	// Decoding String Private Key by PEM
	decodeSecret, rest := pem.Decode([]byte(secret))
	if decodeSecret == nil {
		log.Printf("Decode Secret Err: %v", string(rest))
		return ""
	}

	//Parse PKCS8 Private Key
	privateKey, err := x509.ParsePKCS8PrivateKey(decodeSecret.Bytes)
	if err != nil {
		log.Printf("Signature Err: %v", err)
		return ""
	}

	// MD5 Encrypt Parameters
	strMessage := exchange.Map2UrlQuery(mapParams)
	hasher := md5.New()
	hasher.Write([]byte(strMessage))

	// RSA SHA256 Sign Parameters
	h := sha256.New()
	md5Message := strings.ToUpper(hex.EncodeToString(hasher.Sum(nil)))
	h.Write([]byte(md5Message))
	// Convert PKCS8 to RSA Private Key for signature
	decodeSign, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA256, h.Sum(nil))
	if err != nil {
		return ""
	}

	// Signature converts to URL Decode.
	return url.QueryEscape(base64.StdEncoding.EncodeToString(decodeSign))
}
