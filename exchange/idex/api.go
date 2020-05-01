package idex

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
	API_URL string = "https://api.idex.market"
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
func (e *Idex) GetCoinsData() error {
	coinsData := CoinsData{}

	strRequestUrl := "/returnCurrencies"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	}

	for symbol, data := range coinsData {
		c := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			c = coin.GetCoin(symbol)
			if c == nil {
				c = &coin.Coin{}
				c.Code = symbol
				c.Name = data.Name
				coin.AddCoin(c)
			}
		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(symbol)
		}

		if c != nil {
			coinConstraint := e.GetCoinConstraint(c)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       c.ID,
					Coin:         c,
					ExSymbol:     data.Address,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
				}
			} else {
				coinConstraint.ExSymbol = data.Address
			}

			e.SetCoinConstraint(coinConstraint)
			coinDecimals.Set(c.Code, data.Decimals)
		}
	}
	return nil
}

/* GetPairsData - Get Pairs Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func (e *Idex) GetPairsData() error {
	pairsData := PairsData{}

	strRequestUrl := "/return24Volume"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	}

	for symbol := range pairsData {
		pairStrs := strings.Split(symbol, "_")
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			if len(pairStrs) >= 2 {
				base := coin.GetCoin(pairStrs[0])
				target := coin.GetCoin(pairStrs[1])
				if base != nil && target != nil {
					p = pair.GetPair(base, target)
				}
			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(symbol)
		}

		if p != nil {
			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{
					PairID:      p.ID,
					Pair:        p,
					ExSymbol:    symbol,
					MakerFee:    DEFAULT_MAKER_FEE,
					TakerFee:    DEFAULT_TAKER_FEE,
					LotSize:     DEFAULT_LOT_SIZE,
					PriceFilter: DEFAULT_PRICE_FILTER,
					Listed:      true,
				}
			} else {
				pairConstraint.ExSymbol = symbol
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
func (e *Idex) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	jsonResponse := &JsonResponse{}
	orderBook := OrderBook{}

	strRequestUrl := "/returnOrderBook"
	strUrl := API_URL + strRequestUrl

	mapParams := make(map[string]string)
	mapParams["market"] = e.GetSymbolByPair(pair)
	mapParams["count"] = "100"

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
		buydata.Quantity, err = strconv.ParseFloat(bid.Amount, 64)
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
		selldata.Quantity, err = strconv.ParseFloat(ask.Amount, 64)
		if err != nil {
			return nil, err
		}
		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

func (e *Idex) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *Idex) DoAccountOperation(operation *exchange.AccountOperation) error {
	return nil
}

func (e *Idex) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}
}

func (e *Idex) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	return false
}

func (e *Idex) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/order"

	baseDecimal := 0
	if tmp, ok := coinDecimals.Get(pair.Base.Code); ok {
		baseDecimal = tmp.(int)
	}
	targetDecimal := 0
	if tmp, ok := coinDecimals.Get(pair.Target.Code); ok {
		targetDecimal = tmp.(int)
	}

	mapParams := make(map[string]interface{})
	mapParams["tokenBuy"] = e.GetSymbolByCoin(pair.Base)
	mapParams["amountBuy"] = fmt.Sprintf("%0.0f", quantity*rate*math.Pow10(baseDecimal))
	mapParams["tokenSell"] = e.GetSymbolByCoin(pair.Target)
	mapParams["amountSell"] = fmt.Sprintf("%0.0f", quantity*math.Pow10(targetDecimal))
	mapParams["expires"] = EXPIRES

	jsonPlaceReturn := e.ApiKeyPOST(strRequest, mapParams)
	log.Printf("jsonPlaceReturn: %s", jsonPlaceReturn)
	if err := json.Unmarshal(jsonResponse.Result, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Result Unmarshal Err: %v %s", e.GetName(), err, jsonPlaceReturn)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%v", placeOrder.OrderNumber),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Idex) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/order"

	baseDecimal := 0
	if tmp, ok := coinDecimals.Get(pair.Base.Code); ok {
		baseDecimal = tmp.(int)
	}
	targetDecimal := 0
	if tmp, ok := coinDecimals.Get(pair.Target.Code); ok {
		targetDecimal = tmp.(int)
	}

	mapParams := make(map[string]interface{})
	mapParams["tokenBuy"] = e.GetSymbolByCoin(pair.Target)
	mapParams["amountBuy"] = fmt.Sprintf("%0.0f", quantity*math.Pow10(baseDecimal))
	mapParams["tokenSell"] = e.GetSymbolByCoin(pair.Base)
	mapParams["amountSell"] = fmt.Sprintf("%0.0f", rate*quantity*math.Pow10(targetDecimal))
	mapParams["expires"] = EXPIRES

	jsonPlaceReturn := e.ApiKeyPOST(strRequest, mapParams)
	log.Printf("jsonPlaceReturn: %s", jsonPlaceReturn)
	if err := json.Unmarshal(jsonResponse.Result, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Result Unmarshal Err: %v %s", e.GetName(), err, jsonPlaceReturn)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%v", placeOrder.OrderNumber),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Idex) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	orderStatus := OrderStatus{}
	strRequest := "/orderpending"

	mapParams := make(map[string]interface{})
	mapParams["cmd"] = "orderpending/order"

	orderID, err := strconv.Atoi(order.OrderID)
	if err != nil {
		return fmt.Errorf("convert id from string to int error :%v", err)
	}

	body := make(map[string]interface{})
	body["id"] = orderID

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

func (e *Idex) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Idex) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	cancelOrder := CancelOrder{}
	strRequest := "/orderpending"

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

func (e *Idex) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Idex) ApiKeyPOST(strRequestPath string, mapParams map[string]interface{}) string {
	strRequestUrl := API_URL + strRequestPath

	mapParams["address"] = e.API_KEY
	mapParams["nonce"] = fmt.Sprintf("%d", time.Now().Unix())
	// signature := SoliditySha3(mapParams, e.API_SECRET)
	// log.Printf("signature: %d %s", len(signature), signature)
	// log.Printf("v: %d %d", signature[64], int(signature[64]))
	// mapParams["v"] = int(signature[64]) + 27
	// mapParams["r"] = signature[:32]
	// mapParams["s"] = signature[32:64]

	jsonParams := ""
	if nil != mapParams {
		bytesParams, _ := json.Marshal(mapParams)
		jsonParams = string(bytesParams)
	}

	log.Printf("jsonParams: %s", jsonParams)

	request, err := http.NewRequest("POST", strRequestUrl, strings.NewReader(jsonParams))
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

/* func SoliditySha3(mapParams map[string]interface{}, strSecret string) string {
	hash := solsha3.SoliditySHA3(
		solsha3.Address(CONTRACT_ADDRESS),
		solsha3.Address(mapParams["tokenBuy"]),
		solsha3.Uint256(mapParams["amountBuy"]),
		solsha3.Address(mapParams["tokenSell"]),
		solsha3.Uint256(mapParams["amountSell"]),
		solsha3.Uint256(mapParams["expires"]),
		solsha3.Uint256(mapParams["nonce"]),
		solsha3.Address(mapParams["address"]),
	)

	// Sign Hash Message by Private Key to get V S R message

	return fmt.Sprintln(hex.EncodeToString(hash))
} */
