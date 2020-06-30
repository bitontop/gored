package txbit

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"sync"

	cmap "github.com/orcaman/concurrent-map"

	"github.com/bitontop/gored/coin"
	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/pair"
	"github.com/bitontop/gored/utils"
)

type Txbit struct {
	ID      int
	Name    string `bson:"name"`
	Website string `bson:"website"`

	API_KEY    string
	API_SECRET string

	Source    exchange.DataSource // / exchange API / microservicve api 1 / PSQL
	SourceURI string
}

var pairConstraintMap cmap.ConcurrentMap
var coinConstraintMap cmap.ConcurrentMap
var balanceMap cmap.ConcurrentMap

var instance *Txbit
var once sync.Once

/***************************************************/
func CreateTxbit(config *exchange.Config) *Txbit {
	once.Do(func() {
		instance = &Txbit{
			ID:      DEFAULT_ID,
			Name:    "Txbit",
			Website: "https://txbit.io/",

			API_KEY:    config.API_KEY,
			API_SECRET: config.API_SECRET,
			Source:     config.Source,
			SourceURI:  config.SourceURI,
		}

		balanceMap = cmap.New()
		coinConstraintMap = cmap.New()
		pairConstraintMap = cmap.New()

		if err := instance.InitData(); err != nil {
			log.Printf("%v", err)
			instance = nil
		}
	})
	return instance
}

func (e *Txbit) InitData() error {
	switch e.Source {
	case exchange.EXCHANGE_API:
		if err := e.GetCoinsData(); err != nil {
			return err
		}
		if err := e.GetPairsData(); err != nil {
			return err
		}
		break
	case exchange.MICROSERVICE_API:
		break
	case exchange.JSON_FILE:
		exchangeData := utils.GetExchangeDataFromJSON(e.SourceURI, e.GetName())
		if exchangeData == nil {
			return fmt.Errorf("%s Initial Data Error.", e.GetName())
		} else {
			coinConstraintMap = exchangeData.CoinConstraint
			pairConstraintMap = exchangeData.PairConstraint
		}
		break
	case exchange.PSQL:
	default:
		return fmt.Errorf("%s Initial Coin: There is not selected data source.", e.GetName())
	}
	return nil
}

/**************** Exchange Information ****************/
func (e *Txbit) GetID() int {
	return e.ID
}

func (e *Txbit) GetName() exchange.ExchangeName {
	return exchange.TXBIT
}

func (e *Txbit) GetBalance(coin *coin.Coin) float64 {
	if tmp, ok := balanceMap.Get(coin.Code); ok {
		return tmp.(float64)
	} else {
		return 0.0
	}
}

func (e *Txbit) GetTradingWebURL(pair *pair.Pair) string {
	return fmt.Sprintf("https://txbit.io/Trade/%s", e.GetSymbolByPair(pair))
}

/*************** Coins on the Exchanges ***************/
func (e *Txbit) GetCoinConstraint(coin *coin.Coin) *exchange.CoinConstraint {
	if tmp, ok := coinConstraintMap.Get(fmt.Sprintf("%d", coin.ID)); ok {
		return tmp.(*exchange.CoinConstraint)
	}
	return nil
}

func (e *Txbit) SetCoinConstraint(coinConstraint *exchange.CoinConstraint) {
	coinConstraintMap.Set(fmt.Sprintf("%d", coinConstraint.CoinID), coinConstraint)
}

func (e *Txbit) GetCoins() []*coin.Coin {
	coinList := []*coin.Coin{}
	keySort := []int{}
	for _, key := range coinConstraintMap.Keys() {
		id, _ := strconv.Atoi(key)
		keySort = append(keySort, id)
	}
	sort.Ints(keySort)
	for _, key := range keySort {
		c := coin.GetCoinByID(key)
		if c != nil {
			coinList = append(coinList, c)
		}
	}
	return coinList
}

func (e *Txbit) GetSymbolByCoin(coin *coin.Coin) string {
	key := fmt.Sprintf("%d", coin.ID)
	if tmp, ok := coinConstraintMap.Get(key); ok {
		cc := tmp.(*exchange.CoinConstraint)
		return cc.ExSymbol
	}
	return ""
}

func (e *Txbit) GetCoinBySymbol(symbol string) *coin.Coin {
	for _, id := range coinConstraintMap.Keys() {
		if tmp, ok := coinConstraintMap.Get(id); ok {
			cc := tmp.(*exchange.CoinConstraint)
			if cc.ExSymbol == symbol {
				return cc.Coin
			}
		}
	}
	return nil
}

func (e *Txbit) DeleteCoin(coin *coin.Coin) {
	coinConstraintMap.Remove(fmt.Sprintf("%d", coin.ID))
}

/*************** Pairs on the Exchanges ***************/
func (e *Txbit) GetPairConstraint(pair *pair.Pair) *exchange.PairConstraint {
	if pair == nil{
		return nil
	}
	if tmp, ok := pairConstraintMap.Get(fmt.Sprintf("%d", pair.ID)); ok {
		return tmp.(*exchange.PairConstraint)
	}
	return nil
}

func (e *Txbit) SetPairConstraint(pairConstraint *exchange.PairConstraint) {
	pairConstraintMap.Set(fmt.Sprintf("%d", pairConstraint.PairID), pairConstraint)
}

func (e *Txbit) GetPairs() []*pair.Pair {
	pairList := []*pair.Pair{}
	keySort := []int{}
	for _, key := range pairConstraintMap.Keys() {
		id, _ := strconv.Atoi(key)
		keySort = append(keySort, id)
	}
	sort.Ints(keySort)
	for _, key := range keySort {
		p := pair.GetPairByID(key)
		if p != nil {
			pairList = append(pairList, p)
		}
	}
	return pairList
}

func (e *Txbit) GetPairBySymbol(symbol string) *pair.Pair {
	for _, id := range pairConstraintMap.Keys() {
		if tmp, ok := pairConstraintMap.Get(id); ok {
			pc := tmp.(*exchange.PairConstraint)
			if pc.ExSymbol == symbol {
				return pc.Pair
			}
		}
	}
	return nil
}

func (e *Txbit) GetSymbolByPair(pair *pair.Pair) string {
	pairConstraint := e.GetPairConstraint(pair)
	if pairConstraint != nil {
		return pairConstraint.ExSymbol
	}
	return ""
}

func (e *Txbit) HasPair(pair *pair.Pair) bool {
	return pairConstraintMap.Has(fmt.Sprintf("%d", pair.ID))
}

func (e *Txbit) DeletePair(pair *pair.Pair) {
	pairConstraintMap.Remove(fmt.Sprintf("%d", pair.ID))
}

/**************** Exchange Constraint ****************/
func (e *Txbit) GetConstraintFetchMethod(pair *pair.Pair) *exchange.ConstrainFetchMethod {
	constrainFetchMethod := &exchange.ConstrainFetchMethod{}
	constrainFetchMethod.PublicAPI = true
	constrainFetchMethod.PrivateAPI = true
	constrainFetchMethod.HealthAPI = true
	constrainFetchMethod.HasWithdraw = true
	constrainFetchMethod.HasTransfer = false
	constrainFetchMethod.Fee = false
	constrainFetchMethod.LotSize = false
	constrainFetchMethod.PriceFilter = false
	constrainFetchMethod.TxFee = true
	constrainFetchMethod.Withdraw = true
	constrainFetchMethod.Deposit = true
	constrainFetchMethod.Confirmation = false
	constrainFetchMethod.ConstrainSource = 1
	constrainFetchMethod.ApiRestrictIP = false
	return constrainFetchMethod
}

func (e *Txbit) UpdateConstraint() {
	e.GetCoinsData()
	e.GetPairsData()
}

/**************** Coin Constraint ****************/
func (e *Txbit) GetTxFee(coin *coin.Coin) float64 {
	coinConstraint := e.GetCoinConstraint(coin)
	if coinConstraint == nil {
		return 0.0
	}
	return coinConstraint.TxFee
}

func (e *Txbit) CanWithdraw(coin *coin.Coin) bool {
	coinConstraint := e.GetCoinConstraint(coin)
	if coinConstraint == nil {
		return false
	}
	return coinConstraint.Withdraw
}

func (e *Txbit) CanDeposit(coin *coin.Coin) bool {
	coinConstraint := e.GetCoinConstraint(coin)
	if coinConstraint == nil {
		return false
	}
	return coinConstraint.Deposit
}

func (e *Txbit) GetConfirmation(coin *coin.Coin) int {
	coinConstraint := e.GetCoinConstraint(coin)
	if coinConstraint == nil {
		return 0
	}
	return coinConstraint.Confirmation
}

/**************** Pair Constraint ****************/
func (e *Txbit) GetFee(pair *pair.Pair) float64 {
	pairConstraint := e.GetPairConstraint(pair)
	if pairConstraint == nil {
		return 0.0
	}
	return pairConstraint.TakerFee
}

func (e *Txbit) GetLotSize(pair *pair.Pair) float64 {
	pairConstraint := e.GetPairConstraint(pair)
	if pairConstraint == nil {
		return 0.0
	}
	return pairConstraint.LotSize
}

func (e *Txbit) GetPriceFilter(pair *pair.Pair) float64 {
	pairConstraint := e.GetPairConstraint(pair)
	if pairConstraint == nil {
		return 0.0
	}
	return pairConstraint.PriceFilter
}
