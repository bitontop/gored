package kraken

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
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
	API_URL string = "https://api.kraken.com/0"
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
func (e *Kraken) GetCoinsData() error {
	jsonResponse := &JsonResponse{}
	coinsData := make(map[string]*CoinsData)

	strRequestUrl := "/public/Assets"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	} else if len(jsonResponse.Error) != 0 {
		return fmt.Errorf("%s Get Coins Failed: %v", e.GetName(), jsonResponse.Error)
	}
	if err := json.Unmarshal(jsonResponse.Result, &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	for key, data := range coinsData {
		c := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			c = coin.GetCoin(data.Altname)
			if c == nil {
				c = &coin.Coin{}
				c.Code = data.Altname
				coin.AddCoin(c)
			}
		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(data.Altname)
		}

		if c != nil {
			coinConstraint := &exchange.CoinConstraint{
				CoinID:       c.ID,
				Coin:         c,
				ExSymbol:     key,
				ChainType:    exchange.MAINNET,
				TxFee:        DEFAULT_TXFEE,
				Withdraw:     DEFAULT_WITHDRAW,
				Deposit:      DEFAULT_DEPOSIT,
				Confirmation: DEFAULT_CONFIRMATION,
				Listed:       DEFAULT_LISTED,
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
func (e *Kraken) GetPairsData() error {
	jsonResponse := &JsonResponse{}
	pairsData := make(map[string]*PairsData)

	strRequestUrl := "/public/AssetPairs"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	} else if len(jsonResponse.Error) != 0 {
		return fmt.Errorf("%s Get Pairs Failed: %v", e.GetName(), jsonResponse.Error)
	}
	if err := json.Unmarshal(jsonResponse.Result, &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	for key, data := range pairsData {
		ch := strings.Split(key, ".")
		if len(ch) == 1 {
			p := &pair.Pair{}
			switch e.Source {
			case exchange.EXCHANGE_API:
				base := e.GetCoinBySymbol(data.Quote)
				target := e.GetCoinBySymbol(data.Base)
				if base != nil && target != nil {
					p = pair.GetPair(base, target)
				}
			case exchange.JSON_FILE:
				p = e.GetPairBySymbol(key)
			}
			if p != nil {
				pairConstraint := &exchange.PairConstraint{
					PairID:      p.ID,
					Pair:        p,
					ExSymbol:    key,
					LotSize:     math.Pow10(-1 * data.LotDecimals),
					PriceFilter: math.Pow10(-1 * data.PairDecimals),
					Listed:      DEFAULT_LISTED,
				}
				if len(data.FeesMaker) >= 1 {
					pairConstraint.MakerFee = data.FeesMaker[0][1]
				}
				if len(data.Fees) >= 1 {
					pairConstraint.TakerFee = data.Fees[0][1]
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
func (e *Kraken) OrderBook(pair *pair.Pair) (*exchange.Maker, error) {
	jsonResponse := &JsonResponse{}
	orderBook := make(map[string]*OrderBook)
	symbol := e.GetSymbolByPair(pair)

	mapParams := make(map[string]string)
	mapParams["pair"] = symbol
	mapParams["count"] = "100"

	strRequestUrl := "/public/Depth"
	strUrl := API_URL + strRequestUrl

	maker := &exchange.Maker{}
	maker.WorkerIP = exchange.GetExternalIP()
	maker.BeforeTimestamp = float64(time.Now().UnixNano() / 1e6)

	jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
	if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
	} else if len(jsonResponse.Error) != 0 {
		return nil, fmt.Errorf("%s Get Orderbook Failed: %v", e.GetName(), jsonResponse.Error)
	}
	if err := json.Unmarshal(jsonResponse.Result, &orderBook); err != nil {
		return nil, fmt.Errorf("%s Get Orderbook Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	var err error
	for _, book := range orderBook {
		for _, bid := range book.Bids {
			buydata := exchange.Order{}
			buydata.Quantity, err = strconv.ParseFloat(bid[1].(string), 64)
			if err != nil {
				return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
			}

			buydata.Rate, err = strconv.ParseFloat(bid[0].(string), 64)
			if err != nil {
				return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
			}
			maker.Bids = append(maker.Bids, buydata)
		}
	}
	for _, book := range orderBook {
		for _, ask := range book.Asks {
			selldata := exchange.Order{}
			selldata.Quantity, err = strconv.ParseFloat(ask[1].(string), 64)
			if err != nil {
				return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
			}

			selldata.Rate, err = strconv.ParseFloat(ask[0].(string), 64)
			if err != nil {
				return nil, fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
			}
			maker.Asks = append(maker.Asks, selldata)
		}
	}
	return maker, nil
}

/*************** Private API ***************/
func (e *Kraken) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalances{}
	strRequest := "/private/TradeBalance"

	mapParams := make(map[string]string)
	mapParams["asset"] = "xxbt"

	jsonBalanceReturn := e.ApiKeyPost(strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		log.Printf("%s UpdateAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return
	} else if len(jsonResponse.Error) != 0 {
		log.Printf("%s UpdateAllBalances Failed: %v", e.GetName(), jsonResponse.Error)
		return
	}
	if err := json.Unmarshal(jsonResponse.Result, &accountBalance); err != nil {
		log.Printf("%s UpdateAllBalances Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
		return
	}

	//---------------------- TODO
	/* for _, v := range accountBalance {
		c := e.GetCoinBySymbol(v.Currency)
		if c != nil {
			balanceMap.Set(c.Code, v.Available)
		}
	} */
}

func (e *Kraken) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {

	return false
}

func (e *Kraken) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {

	return nil, nil
}

func (e *Kraken) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {

	return nil, nil
}

func (e *Kraken) OrderStatus(order *exchange.Order) error {

	return nil
}

func (e *Kraken) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *Kraken) CancelOrder(order *exchange.Order) error {

	return nil
}

func (e *Kraken) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *Kraken) ApiKeyPost(strRequestPath string, mapParams map[string]string) string {

	//Signature Request Params

	/* if e.Two_Factor != "" {
		mapParams["otp"] = e.Two_Factor
	} */

	strUrl := API_URL + strRequestPath

	httpClient := &http.Client{}

	jsonParams := ""
	if nil != mapParams {
		bytesParams, _ := json.Marshal(mapParams)
		jsonParams = string(bytesParams)
	}
	jsonParams = exchange.Map2UrlQuery(mapParams)

	request, err := http.NewRequest("POST", strUrl, strings.NewReader(jsonParams))
	if nil != err {
		return err.Error()
	}

	mapParams["nonce"] = fmt.Sprintf("%d", time.Now().UnixNano()/1000000) //time.Now().Unix())
	Signature := ComputeHmac512(strRequestPath, mapParams, e.API_SECRET)

	request.Header.Add("User-Agent", "Kraken GO API Agent (https://github.com/beldur/kraken-go-api-client)")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("API-Key", e.API_KEY)
	request.Header.Add("API-Sign", Signature)

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
func ComputeHmac512(strPath string, mapParams map[string]string, strSecret string) string {
	// bytesParams, _ := json.Marshal(mapParams)
	postData := mapParams["nonce"] + exchange.Map2UrlQuery(mapParams)
	b, _ := json.Marshal(postData)
	log.Printf("postData: %v,\n byte : %v,\n jsonB: %v", postData, []byte(postData), b)
	sha := sha256.New()
	sha.Write([]byte(postData))
	shaSum := sha.Sum(nil)

	strMessage := fmt.Sprintf("%s%s", strPath, string(shaSum))
	log.Printf("strMessage: %v", strMessage)
	decodeSecret, _ := base64.StdEncoding.DecodeString(strSecret)

	h := hmac.New(sha512.New, decodeSecret)
	h.Write([]byte(strMessage))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

/* func (e *Kraken) ApiKeyGet(strRequestPath string, mapParams map[string]string) string {
	//strMethod := "POST"
	strMethod := "GET"

	strUrl := API_URL + strRequestPath

	httpClient := &http.Client{}

	// jsonParams := ""
	// if nil != mapParams {
	// 	bytesParams, _ := json.Marshal(mapParams)
	// 	jsonParams = string(bytesParams)
	// }

	// request, err := http.NewRequest(strMethod, strUrl, strings.NewReader(jsonParams))
	strParams := exchange.Map2UrlQuery(mapParams)
	strRequestUrl := strUrl + "?" + strParams
	request, err := http.NewRequest(strMethod, strRequestUrl, nil)
	if nil != err {
		return err.Error()
	}

	//Signature Request Params
	mapParams["nonce"] = fmt.Sprintf("%d", time.Now().UnixNano())
	if e.Two_Factor != "" {
		mapParams["otp"] = e.Two_Factor
	}
	Signature := ComputeHmac512(strRequestPath, mapParams, e.API_SECRET)

	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("API-Key", e.API_KEY)
	request.Header.Add("API-Sign", Signature)

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
} */

/* func (e *Kraken) ApiKeyGET(strRequestPath string, mapParams map[string]string) string {
	mapParams["apikey"] = e.API_KEY
	mapParams["nonce"] = fmt.Sprintf("%d", time.Now().UnixNano())

	strUrl := API_URL + strRequestPath + "?" + exchange.Map2UrlQuery(mapParams)

	signature := exchange.ComputeHmac512NoDecode(strUrl, e.API_SECRET)
	httpClient := &http.Client{}

	request, err := http.NewRequest("GET", strUrl, nil)
	if nil != err {
		return err.Error()
	}
	request.Header.Add("Content-Type", "application/json;charset=utf-8")
	request.Header.Add("Accept", "application/json")
	request.Header.Add("apisign", signature)

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
} */
