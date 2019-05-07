package exchange

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/pair"

	cmap "github.com/orcaman/concurrent-map"
)

type Exchange interface {
	InitData()

	/***** Exchange Information *****/
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
	GetCoinsData()
	GetPairsData()
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
}

type ExchangeManager struct {
}

var instance *ExchangeManager
var once sync.Once

var exMap cmap.ConcurrentMap
var supportList cmap.ConcurrentMap

func CreateExchangeManager() *ExchangeManager {
	once.Do(func() {
		instance = &ExchangeManager{}
		instance.init()

		if exMap == nil {
			exMap = cmap.New()
		}

		if supportList == nil {
			supportList = cmap.New()
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
	var id int
	for _, key := range supportList.Keys() {
		if tmp, ok := supportList.Get(key); ok {
			if name == tmp.(ExchangeName) {
				id, _ = strconv.Atoi(key)
				break
			}
		}
	}

	return id
}

func (e *ExchangeManager) GetById(i int) Exchange {
	for _, ex := range e.GetExchanges() {
		if ex.GetID() == i {
			return ex
		}
	}
	return nil
}

func (e *ExchangeManager) GetStr(name string) Exchange {
	for _, v := range e.GetSupportExchanges() {
		if string(v) == name {
			return e.Get(v)
		}
	}
	return nil
}

func (e *ExchangeManager) GetSupportExchanges() []ExchangeName {
	list := []ExchangeName{}
	idSort := []int{}
	for _, key := range supportList.Keys() {
		id, _ := strconv.Atoi(key)
		idSort = append(idSort, id)
	}
	sort.Ints(idSort)
	for _, id := range idSort {
		key := fmt.Sprintf("%d", id)
		if tmp, ok := supportList.Get(key); ok {
			list = append(list, tmp.(ExchangeName))
		}
	}

	return list
}

func (e *ExchangeManager) GetExchanges() []Exchange {
	exchanges := []Exchange{}
	for _, key := range exMap.Keys() {
		exchanges = append(exchanges, e.GetStr(key))
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
				eInstance.InitData()
				log.Printf("%s Data Updated. Coin: %d   Pair: %d", exName, len(eInstance.GetCoins()), len(eInstance.GetPairs()))
			}
			time.Sleep(conf.Time)
		}
	}
}
