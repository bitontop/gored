package coinbase

// Contributor 2015-2020 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bitontop/gored/coin"
	exchange "github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
	utils "github.com/bitontop/gored/utils"
)

/*The Base Endpoint URL*/
const (
	API_URL = "https://api.pro.coinbase.com"
)

/*************** PUBLIC  API ***************/
func (e *Coinbase) LoadPublicData(operation *exchange.PublicOperation) error {
	switch operation.Type {
	case exchange.GetCoin:
		return e.doGetCoin(operation)
	case exchange.GetPair:
		return e.doGetPair(operation)
	case exchange.Orderbook:
		switch operation.Wallet {
		case exchange.SpotWallet:
			return e.doSpotOrderBook(operation)
		}
	case exchange.TradeHistory:
		return e.doTradeHistory(operation)

	}
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

func (e *Coinbase) doGetCoin(operation *exchange.PublicOperation) error {
	coinsData := CoinsData{}

	strRequestUrl := "/currencies"
	strUrl := API_URL + strRequestUrl

	jsonCurrencyReturn := exchange.HttpGetRequest(strUrl, nil)
	// log.Printf("jsonCurrencyReturn: %v", jsonCurrencyReturn) // ==========
	if err := json.Unmarshal([]byte(jsonCurrencyReturn), &coinsData); err != nil {
		return fmt.Errorf("%s Get Coins Json Unmarshal Err: %v %v", e.GetName(), err, jsonCurrencyReturn)
	}

	for _, data := range coinsData {
		c := &coin.Coin{}
		switch e.Source {
		case exchange.EXCHANGE_API:
			c = coin.GetCoin(data.ID)
			if c == nil {
				c = &coin.Coin{}
				c.Code = data.ID
				c.Name = data.Name
				coin.AddCoin(c)
			}
		case exchange.JSON_FILE:
			c = e.GetCoinBySymbol(data.ID)
		}

		if c != nil {
			coinConstraint := e.GetCoinConstraint(c)
			if coinConstraint == nil {
				coinConstraint = &exchange.CoinConstraint{
					CoinID:       c.ID,
					Coin:         c,
					ExSymbol:     data.ID,
					ChainType:    exchange.MAINNET,
					TxFee:        DEFAULT_TXFEE,
					Withdraw:     DEFAULT_WITHDRAW,
					Deposit:      DEFAULT_DEPOSIT,
					Confirmation: data.Details.NetworkConfirmations,
					Listed:       data.Status == "online",
				}
			} else {
				coinConstraint.ExSymbol = data.ID
				coinConstraint.Confirmation = data.Details.NetworkConfirmations
				if data.Status == "online" {
					coinConstraint.Listed = true
				} else {
					coinConstraint.Listed = false
				}
			}

			e.SetCoinConstraint(coinConstraint)
		}
	}
	return nil
}

// precision doesn't match
func (e *Coinbase) doGetPair(operation *exchange.PublicOperation) error {
	pairsData := PairsData{}

	strRequestUrl := "/products"
	strUrl := API_URL + strRequestUrl

	jsonSymbolsReturn := exchange.HttpGetRequest(strUrl, nil)
	// log.Printf("jsonSymbolsReturn: %v", jsonSymbolsReturn) // ==========
	if err := json.Unmarshal([]byte(jsonSymbolsReturn), &pairsData); err != nil {
		return fmt.Errorf("%s Get Pairs Json Unmarshal Err: %v %v", e.GetName(), err, jsonSymbolsReturn)
	}

	for _, data := range pairsData {
		if data.Status == "online" {
			p := &pair.Pair{}
			switch e.Source {
			case exchange.EXCHANGE_API:
				base := coin.GetCoin(data.QuoteCurrency)
				target := coin.GetCoin(data.BaseCurrency)
				if base != nil && target != nil {
					p = pair.GetPair(base, target)
				}
			case exchange.JSON_FILE:
				p = e.GetPairBySymbol(data.ID)
			}
			if p != nil {
				var err error
				lotsize := 0.0
				priceFilter := 0.0

				lotsize, err = strconv.ParseFloat(data.BaseIncrement, 64)
				if err != nil {
					log.Printf("%s Lot Size Err: %v", e.GetName(), err)
					lotsize = DEFAULT_LOT_SIZE
				}
				priceFilter, err = strconv.ParseFloat(data.QuoteIncrement, 64)
				if err != nil {
					log.Printf("%s Price Filter Err: %v", e.GetName(), err)
					priceFilter = DEFAULT_PRICE_FILTER
				}

				pairConstraint := e.GetPairConstraint(p)
				if pairConstraint == nil {
					pairConstraint = &exchange.PairConstraint{
						PairID:      p.ID,
						Pair:        p,
						ExSymbol:    data.ID,
						MakerFee:    DEFAULT_MAKER_FEE,
						TakerFee:    DEFAULT_TAKER_FEE,
						LotSize:     lotsize,
						PriceFilter: priceFilter,
						Listed:      true,
					}
				} else {
					pairConstraint.ExSymbol = data.ID
					pairConstraint.LotSize = lotsize
					pairConstraint.PriceFilter = priceFilter
				}
				e.SetPairConstraint(pairConstraint)
			}
		}
	}
	return nil
}

// precision doesn't match
func (e *Coinbase) doSpotOrderBook(op *exchange.PublicOperation) error {
	orderBook := OrderBook{}
	symbol := e.GetSymbolByPair(op.Pair)

	op.Maker = &exchange.Maker{
		Source:          exchange.EXCHANGE_API,
		BeforeTimestamp: float64(time.Now().UnixNano() / 1e6),
	}

	get := &utils.HttpGet{
		URI:       fmt.Sprintf("%s/products/%s/book?level=3", API_URL, symbol),
		Proxy:     op.Proxy,
		DebugMode: op.DebugMode,
	}
	if err := utils.HttpGetRequest(get); err != nil {
		op.Error = err
		return op.Error
	}

	jsonOrderbook := get.ResponseBody
	if err := json.Unmarshal([]byte(jsonOrderbook), &orderBook); err != nil {
		return fmt.Errorf("%s OrderBook json Unmarshal error: %v %v", e.GetName(), err, jsonOrderbook)
	}

	op.Maker.AfterTimestamp = float64(time.Now().UnixNano() / 1e6)
	var err error
	for _, bid := range orderBook.Bids {
		buydata := exchange.Order{}
		buydata.Quantity, err = strconv.ParseFloat(bid[1], 64)
		if err != nil {
			return fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
		}

		buydata.Rate, err = strconv.ParseFloat(bid[0], 64)
		if err != nil {
			return fmt.Errorf("%s OrderBook strconv.ParseFloat Rate error:%v", e.GetName(), err)
		}
		op.Maker.Bids = append(op.Maker.Bids, buydata)
	}
	for _, ask := range orderBook.Asks {
		selldata := exchange.Order{}
		selldata.Quantity, err = strconv.ParseFloat(ask[1], 64)
		if err != nil {
			return fmt.Errorf("%s OrderBook strconv.ParseFloat Quantity error:%v", e.GetName(), err)
		}

		selldata.Rate, err = strconv.ParseFloat(ask[0], 64)
		if err != nil {
			return fmt.Errorf("%s OrderBook strconv.ParseFloat Rate error:%v", e.GetName(), err)
		}
		op.Maker.Asks = append(op.Maker.Asks, selldata)
	}

	return nil
}

func (e *Coinbase) doTradeHistory(operation *exchange.PublicOperation) error {
	symbol := e.GetSymbolByPair(operation.Pair)
	strRequestUrl := fmt.Sprintf("/products/%v/trades", symbol)
	strUrl := API_URL + strRequestUrl

	get := &utils.HttpGet{
		URI:   strUrl,
		Proxy: operation.Proxy,
	}

	err := utils.HttpGetRequest(get)

	if err != nil {
		log.Printf("%+v", err)
		operation.Error = err
		return err

	} else {
		// log.Printf("%+v  ERR:%+v", string(get.ResponseBody), err) // ======================
		if operation.DebugMode {
			operation.RequestURI = get.URI
			operation.CallResponce = string(get.ResponseBody)
		}

		tradeHistory := TradeHistory{}
		if err := json.Unmarshal(get.ResponseBody, &tradeHistory); err != nil {
			operation.Error = err
			return err
		} else {
			// log.Printf("%+v ", tradeHistory)
		}

		operation.TradeHistory = []*exchange.TradeDetail{}
		for _, trade := range tradeHistory {
			price, err := strconv.ParseFloat(trade.Price, 64)
			if err != nil {
				log.Printf("%s price parse Err: %v %v", e.GetName(), err, trade.Price)
				operation.Error = err
				return err
			}
			amount, err := strconv.ParseFloat(trade.Size, 64)
			if err != nil {
				log.Printf("%s amount parse Err: %v %v", e.GetName(), err, trade.Size)
				operation.Error = err
				return err
			}

			td := &exchange.TradeDetail{
				ID:        fmt.Sprintf("%v", trade.TradeID),
				Quantity:  amount,
				TimeStamp: trade.Time.UnixNano() / 1e6,
				Rate:      price,
			}
			if trade.Side == "buy" {
				td.Direction = exchange.Buy
			} else if trade.Side == "sell" {
				td.Direction = exchange.Sell
			}

			operation.TradeHistory = append(operation.TradeHistory, td)
		}
	}

	return nil
}
