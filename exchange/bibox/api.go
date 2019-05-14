package bibox

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
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
	API_URL string = "https://api.bibox.com/v1"
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
func (e *Bibox) GetCoinsData() {
	jsonResponse := &JsonResponse{}
	pairsData := PairsData{}

	strRequestUrl := "/mdata?cmd=pairList"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		log.Printf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if jsonResponse.Error != (Error{}) {
		log.Printf("%s Get Coins Failed: %v", e.GetName(), jsonResponse.Error)
	}
	if err := json.Unmarshal(jsonResponse.Result, &pairsData); err != nil {
		log.Printf("%s Get Coins Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	for _, data := range pairsData {
		pairStrs := strings.Split(data.Pair, "_")

		base := &coin.Coin{}
		target := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base = coin.GetCoin(pairStrs[1])
			if base == nil {
				base = &coin.Coin{}
				base.Code = pairStrs[1]
				coin.AddCoin(base)
			}
			target = coin.GetCoin(pairStrs[0])
			if target == nil {
				target = &coin.Coin{}
				target.Code = pairStrs[0]
				coin.AddCoin(target)
			}
		case exchange.JSON_FILE:
			base = e.GetCoinBySymbol(pairStrs[1])
			target = e.GetCoinBySymbol(pairStrs[0])
		}

		if base != nil {
			coinConstraint := &exchange.CoinConstraint{
				CoinID:       base.ID,
				Coin:         base,
				ExSymbol:     pairStrs[1],
				TxFee:        DEFAULT_TXFEE,
				Withdraw:     DEFAULT_WITHDRAW,
				Deposit:      DEFAULT_DEPOSIT,
				Confirmation: DEFAULT_CONFIRMATION,
				Listed:       DEFAULT_LISTED,
			}
			e.SetCoinConstraint(coinConstraint)
		}
		if target != nil {
			coinConstraint := &exchange.CoinConstraint{
				CoinID:       target.ID,
				Coin:         target,
				ExSymbol:     pairStrs[0],
				TxFee:        DEFAULT_TXFEE,
				Withdraw:     DEFAULT_WITHDRAW,
				Deposit:      DEFAULT_DEPOSIT,
				Confirmation: DEFAULT_CONFIRMATION,
				Listed:       DEFAULT_LISTED,
			}
			e.SetCoinConstraint(coinConstraint)
		}
	}
}

/* GetPairsData - Get Pairs Information (If API provide)
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Add Model of API Response
Step 3: Modify API Path(strRequestUrl)*/
func (e *Bibox) GetPairsData() {
	jsonResponse := &JsonResponse{}
	pairsData := PairsData{}

	strRequestUrl := "/mdata?cmd=pairList"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &jsonResponse); err != nil {
		log.Printf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if jsonResponse.Error != (Error{}) {
		log.Printf("%s Get Pairs Failed: %v", e.GetName(), jsonResponse.Error)
	}
	if err := json.Unmarshal(jsonResponse.Result, &pairsData); err != nil {
		log.Printf("%s Get Pairs Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
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
			pairConstraint := &exchange.PairConstraint{
				PairID:      p.ID,
				Pair:        p,
				ExSymbol:    data.Pair,
				MakerFee:    DEFAULT_MAKERER_FEE,
				TakerFee:    DEFAULT_TAKER_FEE,
				LotSize:     DEFAULT_LOT_SIZE,
				PriceFilter: DEFAULT_PRICE_FILTER,
				Listed:      true,
			}
			e.SetPairConstraint(pairConstraint)
		}
	}
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
	symbol := e.GetSymbolByPair(pair)

	strRequestUrl := fmt.Sprintf("/mdata?cmd=depth&pair=%s", symbol)
	strUrl := API_URL + strRequestUrl

	maker := &exchange.Maker{}
	maker.WorkerIP = exchange.GetExternalIP()
	maker.BeforeTimestamp = float64(time.Now().UnixNano() / 1e6)

	jsonOrderbook := exchange.HttpGetRequest(strUrl, nil)
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

		//Modify according to type and structure
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

/*************** Private API ***************/
func (e *Bibox) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalances{}
	strRequest := "/transfer"

	params := &Asset{
		Cmd: "transfer/assets",
		Body: &AssetDetail{
			Select: 1,
		},
	}
	params.Cmd = "transfer/assets"
	assetDetail := &AssetDetail{}
	assetDetail.Select = 1
	params.Body = assetDetail

	bytes, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	payload := string(bytes)

	jsonBalanceReturn := e.ApiKeyPOST(strRequest, payload)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if jsonResponse.Error != (Error{}) {
		log.Printf("%s UpdateAllBalances Failed: %v", e.GetName(), jsonResponse.Error)
		return
	}
	if err := json.Unmarshal(jsonResponse.Result, &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
		return
	}

	for _, v := range accountBalance[0].AssetsList {
		c := e.GetCoinBySymbol(v.CoinSymbol)
		if c != nil {
			balanceMap.Set(c.Code, v.Balance)
		}
	}
}

func (e *Bibox) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {

	return false
}

func (e *Bibox) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	orderReturn := PlaceOrder{}
	strRequest := "/orderpending"

	params := &OrderParam{
		Cmd:   "orderpending/trade",
		Index: 12345,
		BodyDetails: &OrderParamDetails{
			AccountType: 0,
			Amount:      quantity,
			OrderSide:   2,
			OrderType:   2,
			Pair:        e.GetSymbolByPair(pair),
			Price:       rate,
		},
	}

	bytes, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	payload := string(bytes)

	jsonPlaceReturn := e.ApiKeyPOST(strRequest, payload)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Error != (Error{}) {
		return nil, fmt.Errorf("%s LimitSell Failed: %v", e.GetName(), jsonResponse.Error)
	}
	if err := json.Unmarshal(jsonResponse.Result, &orderReturn); err != nil {
		return nil, fmt.Errorf("%s LimitSell Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%v", orderReturn[0].Result),
		Rate:         rate,
		Quantity:     quantity,
		Side:         "Sell",
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
	orderReturn := PlaceOrder{}
	strRequest := "/orderpending"

	params := &OrderParam{
		Cmd:   "orderpending/trade",
		Index: 12345,
		BodyDetails: &OrderParamDetails{
			AccountType: 0,
			Amount:      quantity,
			OrderSide:   1,
			OrderType:   2,
			Pair:        e.GetSymbolByPair(pair),
			Price:       rate,
		},
	}

	bytes, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	payload := string(bytes)

	jsonPlaceReturn := e.ApiKeyPOST(strRequest, payload)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Error != (Error{}) {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v", e.GetName(), jsonResponse.Error)
	}
	if err := json.Unmarshal(jsonResponse.Result, &orderReturn); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%v", orderReturn[0].Result),
		Rate:         rate,
		Quantity:     quantity,
		Side:         "Buy",
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
	strRequest := "/orderpending"

	orderID, err := strconv.Atoi(order.OrderID)
	if err != nil {
		return fmt.Errorf("convert id from string to int error :%v", err)
	}
	params := &StatusParam{
		Cmd: "orderpending/order",
		BodyDetails: &StatusParamDetails{
			Id: orderID,
		},
	}

	bytes, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	payload := string(bytes)

	jsonOrderStatus := e.ApiKeyPOST(strRequest, payload)
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if jsonResponse.Error != (Error{}) {
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
		order.Status = exchange.Canceled
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
	strRequest := "/orderpending"

	ordersId, err := strconv.Atoi(order.OrderID)
	if err != nil {
		return fmt.Errorf("convert order Id from string %v to int err :%v", order.OrderID, err)
	}
	params := &CancelParam{
		Cmd:   "orderpending/cancelTrade",
		Index: 12345,
		BodyDetails: &CancelParamDetails{
			OrdersId: ordersId,
		},
	}

	bytes, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}
	payload := string(bytes)

	jsonCancelOrder := e.ApiKeyPOST(strRequest, payload)
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if jsonResponse.Error != (Error{}) {
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
func (e *Bibox) ApiKeyPOST(strRequestPath string, payload string) string {

	strRequestUrl := API_URL + strRequestPath

	params := make(map[string]string)
	params["cmds"] = "[" + payload + "]"
	params["apikey"] = e.API_KEY
	sign := ComputeHmacMd5(params["cmds"], e.API_SECRET)
	params["sign"] = sign

	request, err := http.NewRequest("POST", strRequestUrl, strings.NewReader(exchange.Map2UrlQuery(params)))
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

func ComputeHmacMd5(strMessage string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(md5.New, key)
	h.Write([]byte(strMessage))

	return hex.EncodeToString(h.Sum([]byte("")))
}
