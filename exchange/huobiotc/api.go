package huobiotc

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
)

/*The Base Endpoint URL*/
const (
	API_URL = "https://otc-api.eiijo.cn/v1/data"
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
func (e *HuobiOTC) GetCoinsData() error {
	currency := make(map[string]string)
	currency["CNY"] = "1"
	currency["BTC"] = "1"
	currency["USDT"] = "2"

	for symbol, id := range currency {
		c := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			c = coin.GetCoin(symbol)
			if c == nil {
				c = &coin.Coin{}
				c.Code = symbol
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
					ExSymbol:     id,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: DEFAULT_CONFIRMATION,
					Listed:       true,
				}
			} else {
				coinConstraint.ExSymbol = id
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
func (e *HuobiOTC) GetPairsData() error {
	currency := make(map[string]string)
	currency["CNY"] = "1"
	coinID := make(map[string]string)
	coinID["BTC"] = "1"
	coinID["USDT"] = "2"

	for currency, _ := range currency {
		for c, _ := range coinID {
			p := &pair.Pair{}
			switch e.Source {
			case exchange.EXCHANGE_API:
				base := coin.GetCoin(currency)
				target := coin.GetCoin(c)
				if base != nil && target != nil {
					p = pair.GetPair(base, target)
				}
			case exchange.JSON_FILE:
				p = e.GetPairBySymbol(fmt.Sprintf("%s_%s", currency, c))
			}

			if p != nil {
				pairConstraint := e.GetPairConstraint(p)
				if pairConstraint == nil {
					pairConstraint = &exchange.PairConstraint{
						PairID:      p.ID,
						Pair:        p,
						ExSymbol:    fmt.Sprintf("%s_%s", currency, c),
						MakerFee:    DEFAULT_MAKER_FEE,
						TakerFee:    DEFAULT_TAKER_FEE,
						LotSize:     DEFAULT_LOT_SIZE,
						PriceFilter: DEFAULT_PRICE_FILTER,
						Listed:      true,
					}
				} else {
					pairConstraint.ExSymbol = fmt.Sprintf("%s_%s", currency, c)
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
func (e *HuobiOTC) OrderBook(p *pair.Pair) (*exchange.Maker, error) {
	jsonResponse := &JsonResponse{}
	orderBook := OrderBook{}

	strRequestUrl := "/trade-market"
	strUrl := API_URL + strRequestUrl

	mapParams := make(map[string]string)
	mapParams["country"] = "0"   //all countries
	mapParams["payMethod"] = "0" //all methods
	mapParams["currency"] = e.GetSymbolByCoin(p.Base)
	mapParams["coinId"] = e.GetSymbolByCoin(p.Target)
	mapParams["blockType"] = "general"
	mapParams["online"] = "1"

	maker := &exchange.Maker{
		WorkerIP:        exchange.GetExternalIP(),
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	for side := 0; side < 2; side++ {
		currPage := 1
		if side == 0 {
			mapParams["tradeType"] = "sell"
		} else {
			mapParams["tradeType"] = "buy"
		}

		for {
			mapParams["currPage"] = fmt.Sprintf("%v", currPage)

			jsonOrderbook := exchange.HttpGetRequest(strUrl, mapParams)
			if err := json.Unmarshal([]byte(jsonOrderbook), &jsonResponse); err != nil {
				return nil, fmt.Errorf("%s Get Orderbook Json Unmarshal Err: %v %v", e.GetName(), err, jsonOrderbook)
			} else if !jsonResponse.Success {
				return nil, fmt.Errorf("%s Get Pairs Failed: %d %v", e.GetName(), jsonResponse.Code, jsonResponse.Message)
			}
			if err := json.Unmarshal(jsonResponse.Data, &orderBook); err != nil {
				return nil, fmt.Errorf("%s Get Orderbook Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
			}

			maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)

			if side == 0 {
				for _, ask := range orderBook {
					selldata := exchange.Order{}
					selldata.Quantity = ask.TradeCount
					selldata.Rate = ask.Price
					maker.Asks = append(maker.Asks, selldata)
				}
			} else if side == 1 {
				for _, bid := range orderBook {
					buydata := exchange.Order{}
					buydata.Quantity = bid.TradeCount
					buydata.Rate = bid.Price
					maker.Bids = append(maker.Bids, buydata)
				}
			}

			if jsonResponse.CurrPage == jsonResponse.TotalPage {
				break
			} else {
				currPage = jsonResponse.CurrPage + 1
			}
		}
	}

	return maker, nil
}

func (e *HuobiOTC) LoadPublicData(operation *exchange.PublicOperation) error {
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

/*************** Private API ***************/
func (e *HuobiOTC) DoAccountOperation(operation *exchange.AccountOperation) error {
	return nil
}

func (e *HuobiOTC) UpdateAllBalances() {
	if e.API_KEY == "" || e.API_SECRET == "" {
		log.Printf("%s API Key or Secret Key are nil.", e.GetName())
		return
	}
}

/* Withdraw(coin *coin.Coin, quantity float64, addr, tag string) */
func (e *HuobiOTC) Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool {
	return false
}

func (e *HuobiOTC) LimitSell(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	return nil, nil
}

func (e *HuobiOTC) LimitBuy(pair *pair.Pair, quantity, rate float64) (*exchange.Order, error) {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return nil, fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	return nil, nil
}

func (e *HuobiOTC) OrderStatus(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	return nil
}

func (e *HuobiOTC) ListOrders() ([]*exchange.Order, error) {
	return nil, nil
}

func (e *HuobiOTC) CancelOrder(order *exchange.Order) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	return nil
}

func (e *HuobiOTC) CancelAllOrder() error {
	return nil
}

/*************** Signature Http Request ***************/
/*Method: API Request and Signature is required
Step 1: Change Instance Name    (e *<exchange Instance Name>)
Step 2: Create mapParams Depend on API Signature request
Step 3: Add HttpGetRequest below strUrl if API has different requests*/
func (e *HuobiOTC) ApiKeyPost(strRequestPath string, mapParams map[string]string) string {
	return ""
}
