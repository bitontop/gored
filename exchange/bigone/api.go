package bigone

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/base64"
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
	API_URL string = "https://big.one/api/v3"
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
func (e *Bigone) GetCoinsData() error {
	jsonResponse := &JsonResponse{}
	pairsData := PairsData{}

	strRequestUrl := "/asset_pairs"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if jsonResponse.Code != 0 {
		return fmt.Errorf("%s Get Coins Failed: %v", e.GetName(), jsonResponse)
	}
	if err := json.Unmarshal(jsonResponse.Data, &pairsData); err != nil {
		return fmt.Errorf("%s Get Coins Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, data := range pairsData {
		base := &coin.Coin{}
		target := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base = coin.GetCoin(data.QuoteAsset.Symbol)
			if base == nil {
				base = &coin.Coin{}
				base.Code = data.QuoteAsset.Symbol
				base.Name = data.QuoteAsset.Name
				coin.AddCoin(base)
			}
			target = coin.GetCoin(data.BaseAsset.Symbol)
			if target == nil {
				target = &coin.Coin{}
				target.Code = data.BaseAsset.Symbol
				target.Name = data.BaseAsset.Name
				coin.AddCoin(target)
			}
		case exchange.JSON_FILE:
			base = e.GetCoinBySymbol(data.QuoteAsset.Symbol)
			target = e.GetCoinBySymbol(data.BaseAsset.Symbol)
		}

		if base != nil {
			coinConstraint := e.GetCoinConstraint(base)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       base.ID,
					Coin:         base,
					ExSymbol:     data.QuoteAsset.Symbol,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = data.QuoteAsset.Symbol
			}
			e.SetCoinConstraint(coinConstraint)
		}

		if target != nil {
			coinConstraint := e.GetCoinConstraint(target)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       target.ID,
					Coin:         target,
					ExSymbol:     data.BaseAsset.Symbol,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       DEFAULT_LISTED,
				}
			} else {
				coinConstraint.ExSymbol = data.BaseAsset.Symbol
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
func (e *Bigone) GetPairsData() error {
	jsonResponse := &JsonResponse{}
	pairsData := PairsData{}

	strRequestUrl := "/asset_pairs"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if false {
		return fmt.Errorf("%s Get Pairs Failed: %v", e.GetName(), jsonResponse)
	}
	if err := json.Unmarshal(jsonResponse.Data, &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	for _, data := range pairsData {
		p := &pair.Pair{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			base := coin.GetCoin(data.QuoteAsset.Symbol)
			target := coin.GetCoin(data.BaseAsset.Symbol)
			if base != nil && target != nil {
				p = pair.GetPair(base, target)
			}
		case exchange.JSON_FILE:
			p = e.GetPairBySymbol(data.Name)
		}
		if p != nil {
			pairConstraint := e.GetPairConstraint(p)
			if pairConstraint == nil {
				pairConstraint = &exchange.PairConstraint{
					PairID:      p.ID,
					Pair:        p,
					ExSymbol:    data.Name,
					MakerFee:    DEFAULT_MAKER_FEE,
					TakerFee:    DEFAULT_TAKER_FEE,
					LotSize:     math.Pow10(-1 * data.BaseScale),
					PriceFilter: math.Pow10(-1 * data.QuoteScale),
					Listed:      DEFAULT_LISTED,
				}
			} else {
				pairConstraint.ExSymbol = data.Name
				pairConstraint.LotSize = math.Pow10(-1 * data.BaseScale)
				pairConstraint.PriceFilter = math.Pow10(-1 * data.QuoteScale)
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
func (e *Bigone) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	jsonResponse := &JsonResponse{}
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(pair)

	strRequestUrl := fmt.Sprintf("/asset_pairs/%v/depth", symbol)
	strUrl := API_URL + strRequestUrl

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	jsonOrderbook := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if jsonResponse.Code != 0 {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %v", e.GetName(), jsonResponse)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
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
		buydata.Quantity, err = strconv.ParseFloat(bid.Quantity, 64)
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
		selldata.Quantity, err = strconv.ParseFloat(ask.Quantity, 64)
		if err != nil {
			return nil, err
		}

		maker.Asks = append(maker.Asks, selldata)
	}
	return maker, nil
}

func (e *Bigone) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *Bigone) DoAccountOperation(operation *exchange.AccountOperation) error {
	return nil
}

func (e *Bigone) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalances{}
	strRequest := "/viewer/accounts"

	jsonBalanceReturn := e.ApiKeyRequest(strRequest, make(map[string]string), "GET")
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if jsonResponse.Code != 0 {
		log.Printf("%s UpdateAllBalances Failed: %v", e.GetName(), jsonResponse)
		return
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return
	}

	for _, v := range accountBalance {
		c := e.GetCoinBySymbol(v.AssetSymbol)
		if c != nil {
			balance, err := strconv.ParseFloat(v.Balance, 64)
			if err == nil {
				balanceMap.Set(c.Code, balance)
			} else {
				log.Printf("%s balance float64 convert err: %v, %v", e.GetName(), err, v.Balance)
				return
			}
		}
	}
}

// read only withdrawal
func (e *Bigone) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("bigone API Key or Secret Key are nil.")
		return false
	}

	jsonResponse := JsonResponse{}
	withdraw := Withdraw{}

	strRequest := "/viewer/withdrawals"

	mapParams := make(map[string]string)
	mapParams["symbol"] = coin.Code
	mapParams["target_address"] = addr
	mapParams["amount"] = fmt.Sprintf("%v", quantity)
	if tag != "" {
		mapParams["memo"] = tag
	}

	jsonWithdrawReturn := e.ApiKeyRequest(strRequest, mapParams, "POST")
	if err := json.Unmarshal([]byte(jsonWithdrawReturn), &jsonResponse); err != nil {
		log.Printf("%s Withdraw Json Unmarshal Err: %v %v", e.GetName(), err, jsonWithdrawReturn)
		return false
	} else if jsonResponse.Code != 0 {
		log.Printf("%s Withdraw Failed: %v", e.GetName(), jsonWithdrawReturn)
		return false
	}
	if err := json.Unmarshal(jsonResponse.Data, &withdraw); err != nil {
		log.Printf("%s Withdraw Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return false
	}

	return true
}

func (e *Bigone) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/viewer/orders"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["asset_pair_name"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "ASK"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)

	jsonPlaceReturn := e.ApiKeyRequest(strRequest, mapParams, "POST")
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitSell Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != 0 {
		return nil, fmt.Errorf("%s LimitSell Failed: %v", e.GetName(), jsonResponse)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitSell Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%v", placeOrder.ID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Sell,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Bigone) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/viewer/orders"

	priceFilter := int(math.Round(math.Log10(e.GetPriceFilter(pair)) * -1))
	lotSize := int(math.Round(math.Log10(e.GetLotSize(pair)) * -1))

	mapParams := make(map[string]string)
	mapParams["asset_pair_name"] = e.GetSymbolByPair(pair)
	mapParams["side"] = "BID"
	mapParams["price"] = strconv.FormatFloat(rate, 'f', priceFilter, 64)
	mapParams["amount"] = strconv.FormatFloat(quantity, 'f', lotSize, 64)

	jsonPlaceReturn := e.ApiKeyRequest(strRequest, mapParams, "POST")
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != 0 {
		return nil, fmt.Errorf("%s LimitBuy Failed: %v", e.GetName(), jsonResponse)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return nil, fmt.Errorf("%s LimitBuy Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         pair,
		OrderID:      fmt.Sprintf("%v", placeOrder.ID),
		Rate:         rate,
		Quantity:     quantity,
		Direction:    exchange.Buy,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	return order, nil
}

func (e *Bigone) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	orderStatus := PlaceOrder{}
	strRequest := fmt.Sprintf("/viewer/orders/%s", order.OrderID)

	jsonOrderStatus := e.ApiKeyRequest(strRequest, make(map[string]string), "GET")
	if err := json.Unmarshal([]byte(jsonOrderStatus), &jsonResponse); err != nil {
		return fmt.Errorf("%s OrderStatus Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderStatus)
	} else if jsonResponse.Code != 0 {
		return fmt.Errorf("%s OrderStatus Failed: %v", e.GetName(), jsonResponse)
	}
	if err := json.Unmarshal(jsonResponse.Data, &orderStatus); err != nil {
		return fmt.Errorf("%s OrderStatus Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order.StatusMessage = jsonOrderStatus
	filledF, err := strconv.ParseFloat(orderStatus.FilledAmount, 64)
	if err != nil {
		log.Printf("Bigone order statuse parse filled amount to float64 error: %v  %v", err, orderStatus.FilledAmount)
		filledF = 0.0
	}

	if filledF > 0 && filledF < order.Quantity {
		order.Status = exchange.Partial
	} else if orderStatus.State == "FILLED" {
		order.Status = exchange.Filled
	} else if orderStatus.State == "CANCELLED" {
		order.Status = exchange.Cancelled
	} else if filledF == 0 {
		order.Status = exchange.New
	} else if orderStatus.State == "PENDING" {
		order.Status = exchange.New
	} else {
		order.Status = exchange.Other
	}

	return nil
}

func (e *Bigone) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Bigone) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	if order == nil {
		return fmt.Errorf("%s CancelOrder Failed, nil Order", e.GetName())
	}

	mapParams := make(map[string]string)
	mapParams["id"] = order.OrderID

	jsonResponse := &JsonResponse{}
	cancelOrder := PlaceOrder{}
	strRequest := fmt.Sprintf("/viewer/orders/%s/cancel", order.OrderID)

	jsonCancelOrder := e.ApiKeyRequest(strRequest, mapParams, "POST")
	if err := json.Unmarshal([]byte(jsonCancelOrder), &jsonResponse); err != nil {
		return fmt.Errorf("%s CancelOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonCancelOrder)
	} else if jsonResponse.Code != 0 {
		return fmt.Errorf("%s CancelOrder Failed: %v", e.GetName(), jsonResponse)
	}
	if err := json.Unmarshal(jsonResponse.Data, &cancelOrder); err != nil {
		return fmt.Errorf("%s CancelOrder Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order.Status = exchange.Canceling
	order.CancelStatus = jsonCancelOrder

	return nil
}

func (e *Bigone) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Bigone) ApiKeyRequest(strRequestPath string, mapParams map[string]string, method string) string {
	nonce := strconv.FormatInt(time.Now().UnixNano() /* +20*1e9 */, 10) //time.Now().UnixNano() + 20*1e9 //strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	strRequestUrl := API_URL + strRequestPath

	//Signature Request Params
	headerParams := make(map[string]interface{})
	headerParams["alg"] = "HS256"
	headerParams["typ"] = "JWT"

	payloadParams := make(map[string]interface{})
	payloadParams["type"] = "OpenAPIV2"
	payloadParams["sub"] = e.API_KEY
	payloadParams["nonce"] = nonce
	// payloadParams["recv_window"] = "100"

	signature := ""

	// encode header & payload:
	header := ""
	bytesHead, err := json.Marshal(headerParams)
	if nil != headerParams {
		header = string(bytesHead)
	}
	header64 := base64Encode([]byte(header))

	payload := ""
	bytesParams, err := json.Marshal(payloadParams)
	if nil != payloadParams {
		payload = string(bytesParams)
	}
	payload64 := base64Encode([]byte(payload))

	// final token
	signature = exchange.ComputeHmac256URL(header64+"."+payload64, e.API_SECRET)
	token := header64 + "." + payload64 + "." + signature

	request := &http.Request{}
	if method == "GET" {
		strUrl := strRequestUrl + "?" + exchange.Map2UrlQuery(mapParams)
		request, err = http.NewRequest(method, strUrl, nil)
		if err != nil {
			return err.Error()
		}
	} else if method == "POST" {
		postBody, err := json.Marshal(mapParams)
		if err != nil {
			log.Printf("postBody marshal err: %v", err)
		}
		request, err = http.NewRequest(method, strRequestUrl, strings.NewReader(string(postBody)))
		if err != nil {
			return err.Error()
		}
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	request.Header.Add("Content-Type", "application/json")

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		log.Printf("123err=%v", err)
		return err.Error()
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err.Error()
	}

	return string(body)
}

func base64Encode(b []byte) string {
	return base64.URLEncoding.EncodeToString(b)
}
