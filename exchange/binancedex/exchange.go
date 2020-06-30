package binancedex

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/hex"
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
	// "github.com/tendermint/tendermint/crypto/secp256k1"
)

type BinanceDex struct {
	ID      int
	Name    string `bson:"name"`
	Website string `bson:"website"`

	API_KEY    [32]byte
	API_SECRET [32]byte // secp256k1.PrivKeySecp256k1 // gored CI error here

	Source    exchange.DataSource // / exchange API / microservicve api 1 / PSQL
	SourceURI string
}

var pairConstraintMap cmap.ConcurrentMap
var coinConstraintMap cmap.ConcurrentMap
var balanceMap cmap.ConcurrentMap

var instance *BinanceDex
var once sync.Once

/***************************************************/
func CreateBinanceDex(config *exchange.Config) *BinanceDex {
	once.Do(func() {
		instance = &BinanceDex{
			ID:      DEFAULT_ID,
			Name:    "BinanceDex",
			Website: "https://www.binancedex.com/",

			Source:    config.Source,
			SourceURI: config.SourceURI,
		}

		balanceMap = cmap.New()
		coinConstraintMap = cmap.New()
		pairConstraintMap = cmap.New()

		instance.recoveryFromPrivateKey(config.API_SECRET)
		if err := instance.InitData(); err != nil {
			log.Printf("%v", err)
			instance = nil
		}
	})
	return instance
}

func (e *BinanceDex) recoveryFromPrivateKey(privateKey string) error {
	priBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return err
	}

	if len(priBytes) != 32 {
		return fmt.Errorf("Len of Keybytes is not equal to 32 ")
	}

	var keyBytesArray [32]byte
	copy(keyBytesArray[:], priBytes[:32])
	e.API_SECRET = keyBytesArray // secp256k1.PrivKeySecp256k1(keyBytesArray)
	e.API_KEY = keyBytesArray    // e.API_SECRET.PubKey().Address()
	return nil
}

func (e *BinanceDex) InitData() error {
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
func (e *BinanceDex) GetID() int {
	return e.ID
}

func (e *BinanceDex) GetName() exchange.ExchangeName {
	return exchange.BINANCEDEX
}

func (e *BinanceDex) GetTradingWebURL(pair *pair.Pair) string {
	return fmt.Sprintf("https://www.binancedex.com/marker/%s_%s", e.GetSymbolByCoin(pair.Target), e.GetSymbolByCoin(pair.Base))
}

func (e *BinanceDex) GetBalance(coin *coin.Coin) float64 {
	if tmp, ok := balanceMap.Get(coin.Code); ok {
		return tmp.(float64)
	} else {
		return 0.0
	}
}

/*************** Coins on the Exchanges ***************/
func (e *BinanceDex) GetCoinConstraint(coin *coin.Coin) *exchange.CoinConstraint {
	if tmp, ok := coinConstraintMap.Get(fmt.Sprintf("%d", coin.ID)); ok {
		return tmp.(*exchange.CoinConstraint)
	}
	return nil
}

func (e *BinanceDex) SetCoinConstraint(coinConstraint *exchange.CoinConstraint) {
	coinConstraintMap.Set(fmt.Sprintf("%d", coinConstraint.CoinID), coinConstraint)
}

func (e *BinanceDex) GetCoins() []*coin.Coin {
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

func (e *BinanceDex) GetCoinBySymbol(symbol string) *coin.Coin {
	for _, id := range coinConstraintMap.Keys() {
		if tmp, ok := coinConstraintMap.Get(id); ok {
			cc := tmp.(*exchange.CoinConstraint)
			if cc.ExSymbol == symbol {
				return cc.Coin
			}
		} else {
			log.Printf("Get ID %s CoinConstraint Err", id)
		}
	}
	return nil
}

func (e *BinanceDex) GetSymbolByCoin(coin *coin.Coin) string {
	key := fmt.Sprintf("%d", coin.ID)
	if tmp, ok := coinConstraintMap.Get(key); ok {
		cc := tmp.(*exchange.CoinConstraint)
		return cc.ExSymbol
	}
	return ""
}

func (e *BinanceDex) DeleteCoin(coin *coin.Coin) {
	coinConstraintMap.Remove(fmt.Sprintf("%d", coin.ID))
}

/*************** Pairs on the Exchanges ***************/
func (e *BinanceDex) GetPairConstraint(pair *pair.Pair) *exchange.PairConstraint {
	if pair == nil{
		return nil
	}
	if tmp, ok := pairConstraintMap.Get(fmt.Sprintf("%d", pair.ID)); ok {
		return tmp.(*exchange.PairConstraint)
	}
	return nil
}

func (e *BinanceDex) SetPairConstraint(pairConstraint *exchange.PairConstraint) {
	pairConstraintMap.Set(fmt.Sprintf("%d", pairConstraint.PairID), pairConstraint)
}

func (e *BinanceDex) GetPairs() []*pair.Pair {
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

func (e *BinanceDex) GetPairBySymbol(symbol string) *pair.Pair {
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

func (e *BinanceDex) GetSymbolByPair(pair *pair.Pair) string {
	pairConstraint := e.GetPairConstraint(pair)
	if pairConstraint != nil {
		return pairConstraint.ExSymbol
	}
	return ""
}

func (e *BinanceDex) HasPair(pair *pair.Pair) bool {
	return pairConstraintMap.Has(fmt.Sprintf("%d", pair.ID))
}

func (e *BinanceDex) DeletePair(pair *pair.Pair) {
	pairConstraintMap.Remove(fmt.Sprintf("%d", pair.ID))
}

/**************** Exchange Constraint ****************/
func (e *BinanceDex) GetConstraintFetchMethod(pair *pair.Pair) *exchange.ConstrainFetchMethod {
	constrainFetchMethod := &exchange.ConstrainFetchMethod{}
	constrainFetchMethod.PublicAPI = true
	constrainFetchMethod.PrivateAPI = true
	constrainFetchMethod.HealthAPI = true
	constrainFetchMethod.HasWithdraw = false
	constrainFetchMethod.HasTransfer = false
	constrainFetchMethod.Fee = true
	constrainFetchMethod.LotSize = true
	constrainFetchMethod.PriceFilter = true
	constrainFetchMethod.TxFee = true
	constrainFetchMethod.Withdraw = true
	constrainFetchMethod.Deposit = true
	constrainFetchMethod.Confirmation = true
	constrainFetchMethod.ConstrainSource = 1
	constrainFetchMethod.ApiRestrictIP = false
	return constrainFetchMethod
}

func (e *BinanceDex) UpdateConstraint() {
	e.GetCoinsData()
	e.GetPairsData()
}

/**************** Coin Constraint ****************/
func (e *BinanceDex) GetTxFee(coin *coin.Coin) float64 {
	coinConstraint := e.GetCoinConstraint(coin)
	return coinConstraint.TxFee
}

func (e *BinanceDex) CanWithdraw(coin *coin.Coin) bool {
	coinConstraint := e.GetCoinConstraint(coin)
	return coinConstraint.Withdraw
}

func (e *BinanceDex) CanDeposit(coin *coin.Coin) bool {
	coinConstraint := e.GetCoinConstraint(coin)
	return coinConstraint.Deposit
}

func (e *BinanceDex) GetConfirmation(coin *coin.Coin) int {
	coinConstraint := e.GetCoinConstraint(coin)
	return coinConstraint.Confirmation
}

/**************** Pair Constraint ****************/
func (e *BinanceDex) GetFee(pair *pair.Pair) float64 {
	pairConstraint := e.GetPairConstraint(pair)
	return pairConstraint.TakerFee
}

func (e *BinanceDex) GetLotSize(pair *pair.Pair) float64 {
	pairConstraint := e.GetPairConstraint(pair)
	return pairConstraint.LotSize
}

func (e *BinanceDex) GetPriceFilter(pair *pair.Pair) float64 {
	pairConstraint := e.GetPairConstraint(pair)
	return pairConstraint.PriceFilter
}
