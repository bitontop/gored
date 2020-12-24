package exchange

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/pair"

	cmap "github.com/orcaman/concurrent-map"
)

type Exchange interface {
	InitData() error

	/***** Exchange Information *****/
	GetID() int
	GetName() ExchangeName
	GetTradingWebURL(pair *pair.Pair) string
	GetBalance(coin *coin.Coin) float64

	/***** Coin Information *****/
	GetCoinConstraint(coin *coin.Coin) *CoinConstraint
	SetCoinConstraint(coinConstraint *CoinConstraint)
	GetCoins() []*coin.Coin
	GetCoinBySymbol(symbol string) *coin.Coin
	GetSymbolByCoin(coin *coin.Coin) string
	DeleteCoin(coin *coin.Coin)

	/***** Pair Information *****/
	GetPairConstraint(pair *pair.Pair) *PairConstraint
	SetPairConstraint(pairConstraint *PairConstraint)
	GetPairs() []*pair.Pair
	GetPairBySymbol(symbol string) *pair.Pair
	GetSymbolByPair(pair *pair.Pair) string
	HasPair(*pair.Pair) bool
	DeletePair(pair *pair.Pair)

	/***** Public API *****/
	GetCoinsData() error
	GetPairsData() error
	OrderBook(p *pair.Pair) (*Maker, error)

	/***** Private API *****/
	UpdateAllBalances()
	Withdraw(coin *coin.Coin, quantity float64, addr, tag string) bool

	LimitSell(pair *pair.Pair, quantity, rate float64) (*Order, error)
	LimitBuy(pair *pair.Pair, quantity, rate float64) (*Order, error)

	OrderStatus(order *Order) error
	ListOrders() ([]*Order, error)

	CancelOrder(order *Order) error
	CancelAllOrder() error //TODO need to impl cancel all order for exchanges

	/***** Exchange Constraint *****/
	GetConstraintFetchMethod(pair *pair.Pair) *ConstrainFetchMethod
	UpdateConstraint()
	/***** Coin Constraint *****/
	GetTxFee(coin *coin.Coin) float64
	CanWithdraw(coin *coin.Coin) bool
	CanDeposit(coin *coin.Coin) bool
	GetConfirmation(coin *coin.Coin) int
	/***** Pair Constraint *****/
	GetFee(pair *pair.Pair) float64
	GetLotSize(pair *pair.Pair) float64
	GetPriceFilter(pair *pair.Pair) float64

	/***** Ver 2.0  Operation Base Interface *****/
	LoadPublicData(operation *PublicOperation) error
	DoAccountOperation(operation *AccountOperation) error
}

type ExchangeManager struct {
}

var instance *ExchangeManager
var once sync.Once

var exMap cmap.ConcurrentMap
var exIDMap cmap.ConcurrentMap
var supportList = make([]ExchangeName, 0)

func CreateExchangeManager() *ExchangeManager {
	once.Do(func() {
		if instance == nil {
			instance = &ExchangeManager{}

			if exMap == nil {
				exMap = cmap.New()
			}

			if exIDMap == nil {
				exIDMap = cmap.New()
			}

			instance.init()
		}
	})

	return instance
}

func (e *ExchangeManager) init() {
	e.initExchangeNames()
}

func (e *ExchangeManager) Add(exchange Exchange) {
	key := string(exchange.GetName())
	exMap.Set(key, exchange)

	id := fmt.Sprintf("%d", exchange.GetID())
	if exIDMap.Has(id) {
		// log.Fatal("%s ID: %d is exist. Please check.", exchange.GetName(), exchange.GetID())
		log.Printf("%s ID: %d is exist. Please check.", exchange.GetName(), exchange.GetID())
	} else {
		exIDMap.Set(id, exchange)
	}
}

func (e *ExchangeManager) Quantity() int {
	return exMap.Count()
}

func (e *ExchangeManager) Get(name ExchangeName) Exchange {
	if tmp, ok := exMap.Get(string(name)); ok {
		return tmp.(Exchange)
	}
	return nil
}

func (e *ExchangeManager) GetID(name ExchangeName) int {
	eInstance := e.Get(name)
	return eInstance.GetID()
}

func (e *ExchangeManager) GetById(i int) Exchange {
	key := fmt.Sprintf("%d", i)
	if tmp, ok := exIDMap.Get(key); ok {
		return tmp.(Exchange)
	}

	return nil
}

func (e *ExchangeManager) GetStr(name string) Exchange {
	name = strings.ToUpper(strings.TrimSpace(name))
	for _, v := range e.GetSupportExchanges() {
		if string(v) == name {
			return e.Get(v)
		}
	}
	return nil
}

func (e *ExchangeManager) GetSupportExchanges() []ExchangeName {
	return supportList
}

func (e *ExchangeManager) GetExchanges() []Exchange {
	exchanges := []Exchange{}
	idSort := []int{}
	for _, key := range exIDMap.Keys() {
		id, _ := strconv.Atoi(key)
		idSort = append(idSort, id)
	}
	sort.Ints(idSort)
	for _, id := range idSort {
		exchanges = append(exchanges, e.GetById(id))
	}

	return exchanges
}

func (e *ExchangeManager) SubsetPairs(e1, e2 Exchange) []*pair.Pair {
	var pairs []*pair.Pair
	ep1 := e1.GetPairs()

	for _, p := range ep1 {
		if e2.HasPair(p) {
			pairs = append(pairs, p)
		}
	}

	return pairs
}

func (e *ExchangeManager) UpdateExData(conf *Update) {
	switch conf.Method {
	case API_TIGGER:
		break
	case TIME_TIGGER:
		for {
			for _, exName := range conf.ExNames {
				eInstance := e.Get(exName)
				if eInstance != nil {
					if err := eInstance.InitData(); err != nil {
						log.Printf("Updating %s Data is failed.", exName)
					} else {
						log.Printf("%s Data Updated. Coin: %d   Pair: %d", exName, len(eInstance.GetCoins()), len(eInstance.GetPairs()))
					}
				}
			}
			time.Sleep(conf.Time)
		}
	}
}

func (e *ExchangeManager) initExchangeNames() {
	supportList = append(supportList, BINANCE)  // ID = 1
	supportList = append(supportList, BITTREX)  // ID = 2
	supportList = append(supportList, COINEX)   // ID = 3
	supportList = append(supportList, STEX)     // ID = 4
	supportList = append(supportList, BITMEX)   // ID = 5
	supportList = append(supportList, KUCOIN)   // ID = 6
	supportList = append(supportList, BITMAX)   // ID = 7
	supportList = append(supportList, HUOBIOTC) // ID = 8
	supportList = append(supportList, BITSTAMP) // ID = 9
	supportList = append(supportList, OTCBTC)   // ID = 10
	supportList = append(supportList, HUOBI)    // ID = 11
	supportList = append(supportList, BIBOX)    // ID = 12
	supportList = append(supportList, OKEX)     // ID = 13
	supportList = append(supportList, BITZ)     // ID = 14
	supportList = append(supportList, HITBTC)   // ID = 15
	supportList = append(supportList, DRAGONEX) // ID = 16
	supportList = append(supportList, BIGONE)   // ID = 17
	supportList = append(supportList, BITFINEX) // ID = 18
	supportList = append(supportList, GATEIO)   // ID = 19
	supportList = append(supportList, IDEX)     // ID = 20
	supportList = append(supportList, LIQUID)   // ID = 21
	supportList = append(supportList, BITFOREX) // ID = 22
	supportList = append(supportList, TOKOK)    // ID = 23
	supportList = append(supportList, MXC)      // ID = 24
	supportList = append(supportList, BITRUE)   // ID = 25
	supportList = append(supportList, BITATM)   // ID = 26	// not work
	// supportList = append(supportList, TRADESATOSHI) // ID = 27
	supportList = append(supportList, KRAKEN)       // ID = 28
	supportList = append(supportList, POLONIEX)     // ID = 29
	supportList = append(supportList, COINEAL)      // ID = 30
	supportList = append(supportList, TRADEOGRE)    // ID = 31
	supportList = append(supportList, COINBENE)     // ID = 32
	supportList = append(supportList, IBANKDIGITAL) // ID = 33
	supportList = append(supportList, LBANK)        // ID = 34
	// supportList = append(supportList, BINANCEDEX)   // ID = 35
	supportList = append(supportList, BITMART) // ID = 36
	// supportList = append(supportList, GEMINI)    // ID = 37
	supportList = append(supportList, BIKI)      // ID = 38
	supportList = append(supportList, DCOIN)     // ID = 39
	supportList = append(supportList, COINTIGER) // ID = 40
	supportList = append(supportList, BITBAY)    // ID = 41
	supportList = append(supportList, HUOBIDM)   // ID = 42
	supportList = append(supportList, BW)        // ID = 43
	supportList = append(supportList, DERIBIT)   // ID = 44
	supportList = append(supportList, OKEXDM)    // ID = 45
	supportList = append(supportList, GOKO)      // ID = 46
	supportList = append(supportList, BCEX)      // ID = 47
	supportList = append(supportList, DIGIFINEX) // ID = 48
	supportList = append(supportList, LATOKEN)   // ID = 49
	supportList = append(supportList, VIRGOCX)   // ID = 50
	supportList = append(supportList, ABCC)      // ID = 51
	// supportList = append(supportList, BYBIT)     // ID = 52 no orderbook API
	supportList = append(supportList, ZEBITEX)    // ID = 53
	supportList = append(supportList, BITHUMB)    // ID = 54
	supportList = append(supportList, SWITCHEO)   // ID = 55
	supportList = append(supportList, BLOCKTRADE) // ID = 56
	supportList = append(supportList, BKEX)       // ID = 57
	supportList = append(supportList, NEWCAPITAL) // ID = 58
	supportList = append(supportList, COINDEAL)   // ID = 59
	// supportList = append(supportList, HIBITEX)    // ID = 60
	supportList = append(supportList, BGOGO)    // ID = 61
	supportList = append(supportList, FTX)      // ID = 62	orderbook not finished
	supportList = append(supportList, TXBIT)    // ID = 63
	supportList = append(supportList, PROBIT)   // ID = 64
	supportList = append(supportList, BITPIE)   // ID = 65 // api unavailable
	supportList = append(supportList, TAGZ)     // ID = 66
	supportList = append(supportList, IDCM)     // ID = 67
	supportList = append(supportList, HOO)      // ID = 68
	supportList = append(supportList, HOMIEX)   // ID = 69
	supportList = append(supportList, COINBASE) // ID = 70
	// supportList = append(supportList, NICEHASH) // ID = 72
	// supportList = append(supportList, BITBNS) // ID = 73
	supportList = append(supportList, OKSIM) // ID = 74
}
