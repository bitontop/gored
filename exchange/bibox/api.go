package bibox

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
	API_URL string = "https://api.bibox.com"
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
func (e *Bibox) GetCoinsData() error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	coinsData := CoinsData{}
	strRequestUrl := "/v1/transfer"

	mapParams := make(map[string]interface{})
	mapParams["cmd"] = "transfer/coinConfig"

	body := make(map[string]interface{})
	mapParams["body"] = body

	jsonCurrencyReturn := e.ApiKeyPOST(strRequestUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if jsonResponse.Error != (Error{}) {
		return fmt.Errorf("%s Get Coins Failed: %v", e.GetName(), jsonResponse.Error)
	}
	if err := json.Unmarshal(jsonResponse.Result, &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	for _, result := range coinsData {
		for _, data := range result.Result {
			c := &coin.Coin{}
			switch e.Source {
			case exchange.EXCHANGE_API:
				c = coin.GetCoin(data.CoinSymbol)
				if c == nil {
					c = &coin.Coin{}
					c.Code = data.CoinSymbol
					coin.AddCoin(c)
				}
			case exchange.JSON_FILE:
				c = e.GetCoinBySymbol(data.CoinSymbol)
			}

			if c != nil {
				coinConstraint := e.GetCoinConstraint(c)
				if coinConstraint == nil {
					coinConstraint = &exchange.CoinConstraint{
						CoinID:       c.ID,
						Coin:         c,
						ExSymbol:     data.CoinSymbol,
						ChainType:    exchange.MAINNET,
						TxFee:        data.WithdrawFee,
						Confirmation: DEFAULT_CONFIRMATION,
					}
				} else {
					coinConstraint.ExSymbol = data.CoinSymbol
					coinConstraint.TxFee = data.WithdrawFee
				}

				if data.EnableDeposit == 1 {
					coinConstraint.Deposit = true
				} else {
					coinConstraint.Deposit = false
				}

				if data.EnableWithdraw == 1 {
					coinConstraint.Withdraw = true
				} else {
					coinConstraint.Withdraw = false
				}

				if data.IsActive == 1 {
					coinConstraint.Listed = true
				} else {
					coinConstraint.Listed = false
				}

				e.SetCoinConstraint(coinConstraint)
			}
		}
	}
	return nil
}

/* GetPairsData - Get Pairs Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func (e *Bibox) GetPairsData() error {
	jsonResponse := &JsonResponse{}
	pairsData := PairsData{}

	strRequestUrl := "/v1/mdata"
	strUrl := API_URL + strRequestUrl

	mapParams := make(map[string]string)
	mapParams["cmd"] = "pairList"

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if jsonResponse.Error != (Error{}) {
		return fmt.Errorf("%s Get Pairs Failed: %v", e.GetName(), jsonResponse.Error)
	}
	if err := json.Unmarshal(jsonResponse.Result, &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	for _, data := range pairsData {
		pairStrs := strings.Split(data.Pair, "_")
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(pairStrs[1])
			target := coin.GetCoin(pairStrs[0])
			if base != nil && target != nil {
				p = pair.GetPair(base, target)
			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(data.Pair)
		}

		if p != nil {
			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{
					PairID:      p.ID,
					Pair:        p,
					ExSymbol:    data.Pair,
					MakerFee:    DEFAULT_MAKER_FEE,
					TakerFee:    DEFAULT_TAKER_FEE,
					LotSize:     DEFAULT_LOT_SIZE,
					PriceFilter: DEFAULT_PRICE_FILTER,
					Listed:      true,
				}
			} else {
				pairConstraint.ExSymbol = data.Pair
			}

			if pairStrs[1] == "USDT" || pairStrs[1] == "DAI" {
				pairConstraint.PriceFilter = USDT_PRICE_FILTER
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
func (e *Bibox) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	jsonResponse := &JsonResponse{}
	orderBook := OrderBook{}

	strRequestUrl := "/v1/mdata"
	strUrl := API_URL + strRequestUrl

	mapParams := make(map[string]string)
	mapParams["cmd"] = "depth"
	mapParams["pair"] = e.GetSymbolByPair(pair)

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if jsonResponse.Error != (Error{}) {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %v", e.GetName(), jsonResponse.Error)
	}
	if err := json.Unmarshal(jsonResponse.Result, &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)

	var err error
	for _, bid := range orderBook.Bids {
		var buydata exchange.Order

		//Modify according to type and structure
		buydata.Rate, err = strconv.ParseFloat(bid.Price, 64)
		if err != nil {
			return nil, err
		}
		buydata.Quantity, err = strconv.ParseFloat(bid.Volume, 64)
		if err != nil {
			return nil, err
		}

		maker.Bids = append(maker.Bids, buydata)
	}
	for _, ask := range orderBook.Asks {
		var selldata exchange.Order

		selldata.Rate, err = strconv.ParseFloat(ask.Price, 64)
		if err != nil {
			return nil, err
		}
		selldata.Quantity, err = strconv.ParseFloat(ask.Volume, 64)
		if err != nil {
			return nil, err
		}
		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

func (e *Bibox) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *Bibox) DoAccountOperation(operation *exchange.AccountOperation) error {

	switch operation.Type {

	case exchange.Transfer:
		return e.transfer(operation)
	// case exchange.BalanceList:
	// 	return e.getAllBalance(operation)
	// case exchange.Balance:
	// 	return e.getBalance(operation)

	case exchange.Withdraw:
		return e.doWithdraw(operation)

	}
	return fmt.Errorf("%s Operation type invalid: %s %v", operation.Ex, operation.Wallet, operation.Type)
}

// direction: 0钱包转币币; 1币币转钱包
func (e *Bibox) transfer(operation *exchange.AccountOperation) error { //(coin *coin.Coin, quantity float64, direction int) bool {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	innerTrans := InnerTrans{}
	strRequest := "/v2/assets/transfer/spot"

	mapParams := make(map[string]interface{})

	mapParams["symbol"] = e.GetSymbolByCoin(operation.Coin)
	mapParams["amount"] = operation.TransferAmount
	switch operation.TransferFrom {
	case exchange.AssetWallet:
		mapParams["type"] = 0
	case exchange.SpotWallet:
		mapParams["type"] = 1
	default:
		return fmt.Errorf("%s transfer type not supported: %v", e.GetName(), operation.TransferFrom)
	}

	jsonInnerReturn := e.ApiKeyPOSTInner(strRequest, mapParams) //ApiKeyPOST
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonInnerReturn
	}

	if err := json.Unmarshal([]byte(jsonInnerReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s Transfer Json Unmarshal Err: %v, %v", e.GetName(), err, jsonInnerReturn)
		return operation.Error
	} else if jsonResponse.Error.Code != "" {
		operation.Error = fmt.Errorf("%s Transfer Failed: %v", e.GetName(), jsonInnerReturn)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Result, &innerTrans); err != nil {
		operation.Error = fmt.Errorf("%s Transfer Result Unmarshal Err: %v %s", e.GetName(), err, jsonInnerReturn)
		return operation.Error
	}

	return nil
}

// googleAuth: google验证码
// tag： the remark of withdraw address
// memo: memo is required for some tokens, such as EOS
// if it's new address, need to provide trade_pwd and google code. Could do this on webpage.
func (e *Bibox) doWithdraw(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("bibox API Key or Secret Key are nil.")
	}

	jsonResponse := JsonResponse{}
	withdraw := Withdraw{}

	strRequestUrl := "/v1/transfer"

	mapParams := make(map[string]interface{})
	mapParams["cmd"] = "transfer/transferOut"

	body := make(map[string]interface{})
	body["coin_symbol"] = e.GetSymbolByCoin(operation.Coin)
	body["amount"] = operation.WithdrawAmount
	body["addr"] = operation.WithdrawAddress
	if operation.WithdrawTag != "" {
		body["addr_remark"] = operation.WithdrawTag
	}
	// body["totp_code"] = 123456     //googleAuth, need if new addr
	// body["trade_pwd"] = "tradePWD" //tradePWD, need if new addr
	// body["memo"] = "" // memo is required for some tokens, such as EOS

	mapParams["body"] = body

	jsonWithdraw := e.ApiKeyPOST(strRequestUrl, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequestUrl
		operation.CallResponce = jsonWithdraw
	}

	if err := json.Unmarshal([]byte(jsonWithdraw), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Json Unmarshal Err: %v, %s", e.GetName(), err, jsonWithdraw)
		return operation.Error
	} else if jsonResponse.Error.Code != "" {
		operation.Error = fmt.Errorf("%s Withdraw Failed: %v", e.GetName(), jsonWithdraw)
		return operation.Error
	}
	if err := json.Unmarshal([]byte(jsonWithdraw), &withdraw); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Result Unmarshal Err: %v %s", e.GetName(), err, jsonWithdraw)
		return operation.Error
	}

	operation.WithdrawID = fmt.Sprintf("%v", withdraw.Result)

	return nil
}

func (e *Bibox) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	// test Inner Transfer, withdraw
	// e.InnerTrans(e.GetCoinBySymbol("ETH"), 0.01, 0)
	//e.Withdraw(e.GetCoinBySymbol("ETH"), 100, "addr", "tag") //, 123456, "tradePWD", "")

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalances{}
	strRequest := "/v1/transfer"

	mapParams := make(map[string]interface{})
	mapParams["cmd"] = "transfer/assets"

	body := make(map[string]interface{})
	body["select"] = 1

	mapParams["body"] = body

	jsonBalanceReturn := e.ApiKeyPOST(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if jsonResponse.Error.Code != "" {
		log.Printf("%s UpdateAllBalances Failed: %v", e.GetName(), jsonResponse.Error)
		return
	}
	if err := json.Unmarshal(jsonResponse.Result, &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
		return
	}

	for _, result := range accountBalance {
		for _, v := range result.Result.AssetsList {
			freeamount, err := strconv.ParseFloat(v.Balance, 64)
			if err == nil {
				c := e.GetCoinBySymbol(v.CoinSymbol)
				if c != nil {
					balanceMap.Set(c.Code, freeamount)
				}
			} else {
				log.Printf("%s %s Get Balance Err: %s\n", e.GetName(), v.CoinSymbol, err)
			}
		}
	}

}

// googleAuth: google验证码
// tag： 提现地址备注
// memo： 提现标签(can be "", not required)
// need to update interface to use more params
func (e *Bibox) Withdraw(coin *coin.Coin, quantity float64, addr, tag string /* , googleAuth int, tradePWD, memo string */) bool {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("bibox API Key or Secret Key are nil.")
		return false
	}

	jsonResponse := JsonResponse{}
	withdraw := Withdraw{}

	strRequestUrl := "/v1/transfer"

	mapParams := make(map[string]interface{})
	mapParams["cmd"] = "transfer/transferOut"

	body := make(map[string]interface{})
	body["coin_symbol"] = e.GetSymbolByCoin(coin)
	body["amount"] = quantity
	body["totp_code"] = 123456     //googleAuth
	body["trade_pwd"] = "tradePWD" //tradePWD
	body["addr"] = addr
	body["addr_remark"] = tag
	body["memo"] = "" //memo

	mapParams["body"] = body

	jsonWithdraw := e.ApiKeyPOST(strRequestUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonWithdraw), &jsonResponse); err != nil {
		log.Printf("%s Withdraw Json Unmarshal Err: %v %v", e.GetName(), err, jsonWithdraw)
		return false
	} else if jsonResponse.Error.Code != "" {
		log.Printf("%s Withdraw Failed: %v", e.GetName(), jsonWithdraw)
		return false
	}
	if err := json.Unmarshal([]byte(jsonWithdraw), &withdraw); err != nil {
		log.Printf("%s Withdraw Result Unmarshal Err: %v %s", e.GetName(), err, jsonWithdraw)
		return false
	}

	return true
}

func (e *Bibox) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/v1/orderpending"

	mapParams := make(map[string]interface{})
	mapParams["cmd"] = "orderpending/trade"

	body := make(map[string]interface{})
	body["pair"] = e.GetSymbolByPair(pair)
	body["account_type"] = 0
	body["order_type"] = 2
	body["order_side"] = 2
	body["price"] = strconv.FormatFloat(rate, 'f', -1, 64)
	body["amount"] = strconv.FormatFloat(quantity, 'f', -1, 64)

	mapParams["body"] = body

	jsonPlaceReturn := e.ApiKeyPOST(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Error.Code != "" {
		return nil, fmt.Errorf("%s LimitSell Failed: %v", e.GetName(), jsonResponse.Error)
	}
	if err := json.Unmarshal(jsonResponse.Result, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      placeOrder[0].Result,
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Bibox) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/v1/orderpending"

	mapParams := make(map[string]interface{})
	mapParams["cmd"] = "orderpending/trade"

	body := make(map[string]interface{})
	body["pair"] = e.GetSymbolByPair(pair)
	body["account_type"] = 0
	body["order_type"] = 2
	body["order_side"] = 1
	body["price"] = strconv.FormatFloat(rate, 'f', -1, 64)
	body["amount"] = strconv.FormatFloat(quantity, 'f', -1, 64)

	mapParams["body"] = body

	jsonPlaceReturn := e.ApiKeyPOST(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Error.Code != "" {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v", e.GetName(), jsonResponse.Error)
	}
	if err := json.Unmarshal(jsonResponse.Result, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      placeOrder[0].Result,
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Bibox) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	orderStatus := OrderStatus{}
	strRequest := "/v1/orderpending"

	mapParams := make(map[string]interface{})
	mapParams["cmd"] = "orderpending/order"

	orderID, err := strconv.Atoi(order.OrderID)
	if err != nil {
		return fmt.Errorf("convert id from string to int error :%v", err)
	}

	body := make(map[string]interface{})
	body["id"] = orderID
	body["account_type"] = 0

	mapParams["body"] = body

	jsonOrderStatus := e.ApiKeyPOST(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if jsonResponse.Error.Code != "" {
		return fmt.Errorf("%s OrderStatus Failed: %v", e.GetName(), jsonResponse.Error)
	}
	if err := json.Unmarshal(jsonResponse.Result, &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	order.StatusMessage = jsonOrderStatus
	if orderStatus[0].Result.Status == 1 {
		order.Status = exchange.New
	} else if orderStatus[0].Result.Status == 2 {
		order.Status = exchange.Partial
	} else if orderStatus[0].Result.Status == 3 {
		order.Status = exchange.Filled
	} else if orderStatus[0].Result.Status == 4 {
		order.Status = exchange.Canceling
	} else if orderStatus[0].Result.Status == 5 {
		order.Status = exchange.Cancelled
	} else if orderStatus[0].Result.Status == 6 {
		order.Status = exchange.Canceling

	}

	return nil
}

func (e *Bibox) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Bibox) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	cancelOrder := CancelOrder{}
	strRequest := "/v1/orderpending"

	mapParams := make(map[string]interface{})
	mapParams["cmd"] = "orderpending/cancelTrade"

	orderID, err := strconv.Atoi(order.OrderID)
	if err != nil {
		return fmt.Errorf("convert id from string to int error :%v", err)
	}

	body := make(map[string]interface{})
	body["orders_id"] = orderID

	mapParams["body"] = body

	jsonCancelOrder := e.ApiKeyPOST(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if jsonResponse.Error.Code != "" {
		return fmt.Errorf("%s CancelOrder Failed: %v", e.GetName(), jsonResponse.Error)
	}
	if err := json.Unmarshal(jsonResponse.Result, &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	} else if cancelOrder[0].Result != "撤销中" {
		return fmt.Errorf("%s Cancel Order error :%v", e.GetName(), cancelOrder[0].Result)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Bibox) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Bibox) ApiKeyPOST(strRequestPath string, mapParams map[string]interface{}) string {
	strRequestUrl := API_URL + strRequestPath

	jsonParams := ""
	if nil != mapParams {
		bytesParams, _ := json.Marshal(mapParams)
		jsonParams = "[" + string(bytesParams) + "]"
	}

	Params := make(map[string]string)
	Params["cmds"] = jsonParams
	Params["apikey"] = e.API_KEY
	Params["sign"] = exchange.ComputeHmacMd5(jsonParams, e.API_SECRET)

	// log.Printf("Params: %+v", Params)

	request, err := http.NewRequest("POST", strRequestUrl, strings.NewReader(exchange.Map2UrlQuery(Params)))
	if err != nil {
		return err.Error()
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Accept", "application/json")

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

func (e *Bibox) ApiKeyPOSTInner(strRequestPath string, mapParams map[string]interface{}) string {
	strRequestUrl := API_URL + strRequestPath

	jsonParams := ""
	if nil != mapParams {
		bytesParams, _ := json.Marshal(mapParams)
		jsonParams = /* "[" +  */ string(bytesParams) /* + "]" */
	}

	Params := make(map[string]string)
	Params["body"] = jsonParams //"{\"type\":0,\"amount\":\"100\",\"symbol\":\"USDT\"}"
	Params["apikey"] = e.API_KEY
	Params["sign"] = exchange.ComputeHmacMd5(jsonParams, e.API_SECRET)

	// log.Printf(`Params["body"]: %v`, Params["body"])
	//log.Printf("====Post mapParams: %+v", Params)

	request, err := http.NewRequest("POST", strRequestUrl, strings.NewReader(exchange.Map2UrlQuery(Params)))
	if err != nil {
		return err.Error()
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Accept", "application/json")

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
