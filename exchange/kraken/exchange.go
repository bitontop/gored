package kraken

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

type Kraken struct {
	ID      int
	Name    string `bson:"name"`
	Website string `bson:"website"`

	API_KEY    string
	API_SECRET string
	Two_Factor string

	Source    exchange.DataSource // / exchange API / microservicve api 1 / PSQL
	SourceURI string
}

var pairConstraintMap cmap.ConcurrentMap
var coinConstraintMap cmap.ConcurrentMap
var balanceMap cmap.ConcurrentMap

var instance *Kraken
var once sync.Once

/***************************************************/
func CreateKraken(config *exchange.Config) *Kraken {
	once.Do(func() {
		instance = &Kraken{
			ID:      DEFAULT_ID,
			Name:    "Kraken",
			Website: "https://www.kraken.com/",

			API_KEY:    config.API_KEY,
			API_SECRET: config.API_SECRET,
			Two_Factor: config.Two_Factor,
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

func (e *Kraken) InitData() error {
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
func (e *Kraken) GetID() int {
	return e.ID
}

func (e *Kraken) GetName() exchange.ExchangeName {
	return exchange.KRAKEN
}

func (e *Kraken) GetBalance(coin *coin.Coin) float64 {
	if tmp, ok := balanceMap.Get(coin.Code); ok {
		return tmp.(float64)
	} else {
		return 0.0
	}
}

func (e *Kraken) GetTradingWebURL(pair *pair.Pair) string {
	return fmt.Sprintf("https://trade.kraken.com/markets/kraken/%s/%s", pair.Target.Code, pair.Base.Code)
}

/*************** Coins on the Exchanges ***************/
func (e *Kraken) GetCoinConstraint(coin *coin.Coin) *exchange.CoinConstraint {
	if tmp, ok := coinConstraintMap.Get(fmt.Sprintf("%d", coin.ID)); ok {
		return tmp.(*exchange.CoinConstraint)
	}
	return nil
}

func (e *Kraken) SetCoinConstraint(coinConstraint *exchange.CoinConstraint) {
	coinConstraintMap.Set(fmt.Sprintf("%d", coinConstraint.CoinID), coinConstraint)
}

func (e *Kraken) GetCoins() []*coin.Coin {
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

func (e *Kraken) GetSymbolByCoin(coin *coin.Coin) string {
	key := fmt.Sprintf("%d", coin.ID)
	if tmp, ok := coinConstraintMap.Get(key); ok {
		cc := tmp.(*exchange.CoinConstraint)
		return cc.ExSymbol
	}
	return ""
}

func (e *Kraken) GetCoinBySymbol(symbol string) *coin.Coin {
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

func (e *Kraken) DeleteCoin(coin *coin.Coin) {
	coinConstraintMap.Remove(fmt.Sprintf("%d", coin.ID))
}

/*************** Pairs on the Exchanges ***************/
func (e *Kraken) GetPairConstraint(pair *pair.Pair) *exchange.PairConstraint {
	if pair == nil{
		return nil
	}
	if tmp, ok := pairConstraintMap.Get(fmt.Sprintf("%d", pair.ID)); ok {
		return tmp.(*exchange.PairConstraint)
	}
	return nil
}

func (e *Kraken) SetPairConstraint(pairConstraint *exchange.PairConstraint) {
	pairConstraintMap.Set(fmt.Sprintf("%d", pairConstraint.PairID), pairConstraint)
}

func (e *Kraken) GetPairs() []*pair.Pair {
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

func (e *Kraken) GetPairBySymbol(symbol string) *pair.Pair {
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

func (e *Kraken) GetSymbolByPair(pair *pair.Pair) string {
	pairConstraint := e.GetPairConstraint(pair)
	if pairConstraint != nil {
		return pairConstraint.ExSymbol
	}
	return ""
}

func (e *Kraken) HasPair(pair *pair.Pair) bool {
	return pairConstraintMap.Has(fmt.Sprintf("%d", pair.ID))
}

func (e *Kraken) DeletePair(pair *pair.Pair) {
	pairConstraintMap.Remove(fmt.Sprintf("%d", pair.ID))
}

/**************** Exchange Constraint ****************/
func (e *Kraken) GetConstraintFetchMethod(pair *pair.Pair) *exchange.ConstrainFetchMethod {
	constrainFetchMethod := &exchange.ConstrainFetchMethod{}
	constrainFetchMethod.PublicAPI = true
	constrainFetchMethod.PrivateAPI = true
	constrainFetchMethod.HealthAPI = true
	constrainFetchMethod.HasWithdraw = true
	constrainFetchMethod.HasTransfer = false
	constrainFetchMethod.Fee = true
	constrainFetchMethod.LotSize = true
	constrainFetchMethod.PriceFilter = true
	constrainFetchMethod.TxFee = false
	constrainFetchMethod.Withdraw = false
	constrainFetchMethod.Deposit = false
	constrainFetchMethod.Confirmation = false
	constrainFetchMethod.ConstrainSource = 1
	constrainFetchMethod.ApiRestrictIP = false
	return constrainFetchMethod
}

func (e *Kraken) UpdateConstraint() {
	e.GetCoinsData()
	e.GetPairsData()
}

/**************** Coin Constraint ****************/
func (e *Kraken) GetTxFee(coin *coin.Coin) float64 {
	coinConstraint := e.GetCoinConstraint(coin)
	if coinConstraint == nil {
		return 0.0
	}
	return coinConstraint.TxFee
}

func (e *Kraken) CanWithdraw(coin *coin.Coin) bool {
	coinConstraint := e.GetCoinConstraint(coin)
	if coinConstraint == nil {
		return false
	}
	return coinConstraint.Withdraw
}

func (e *Kraken) CanDeposit(coin *coin.Coin) bool {
	coinConstraint := e.GetCoinConstraint(coin)
	if coinConstraint == nil {
		return false
	}
	return coinConstraint.Deposit
}

func (e *Kraken) GetConfirmation(coin *coin.Coin) int {
	coinConstraint := e.GetCoinConstraint(coin)
	if coinConstraint == nil {
		return 0
	}
	return coinConstraint.Confirmation
}

/**************** Pair Constraint ****************/
func (e *Kraken) GetFee(pair *pair.Pair) float64 {
	pairConstraint := e.GetPairConstraint(pair)
	if pairConstraint == nil {
		return 0.0
	}
	return pairConstraint.TakerFee
}

func (e *Kraken) GetLotSize(pair *pair.Pair) float64 {
	pairConstraint := e.GetPairConstraint(pair)
	if pairConstraint == nil {
		return 0.0
	}
	return pairConstraint.LotSize
}

func (e *Kraken) GetPriceFilter(pair *pair.Pair) float64 {
	pairConstraint := e.GetPairConstraint(pair)
	if pairConstraint == nil {
		return 0.0
	}
	return pairConstraint.PriceFilter
}
