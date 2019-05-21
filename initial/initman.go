package initial

import (
	"sync"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/exchange/bibox"
	"github.com/bitontop/gored/exchange/bigone"
	"github.com/bitontop/gored/exchange/binance"
	"github.com/bitontop/gored/exchange/bitfinex"
	"github.com/bitontop/gored/exchange/bitforex"
	"github.com/bitontop/gored/exchange/bitmax"
	"github.com/bitontop/gored/exchange/bitmex"
	"github.com/bitontop/gored/exchange/bitrue"
	"github.com/bitontop/gored/exchange/bitstamp"
	"github.com/bitontop/gored/exchange/bittrex"
	"github.com/bitontop/gored/exchange/bitz"
	"github.com/bitontop/gored/exchange/coinex"
	"github.com/bitontop/gored/exchange/dragonex"
	"github.com/bitontop/gored/exchange/gateio"
	"github.com/bitontop/gored/exchange/hitbtc"
	"github.com/bitontop/gored/exchange/huobi"
	"github.com/bitontop/gored/exchange/huobiotc"
	"github.com/bitontop/gored/exchange/idex"
	"github.com/bitontop/gored/exchange/kucoin"
	"github.com/bitontop/gored/exchange/liquid"
	"github.com/bitontop/gored/exchange/mxc"
	"github.com/bitontop/gored/exchange/okex"
	"github.com/bitontop/gored/exchange/otcbtc"
	"github.com/bitontop/gored/exchange/stex"
	"github.com/bitontop/gored/exchange/tokok"
)

var instance *InitManager
var once sync.Once

type InitManager struct {
	exMan *exchange.ExchangeManager
}

func CreateInitManager() *InitManager {
	once.Do(func() {
		instance = &InitManager{
			exMan: exchange.CreateExchangeManager(),
		}
	})
	return instance
}

func (e *InitManager) Init(config *exchange.Config) exchange.Exchange {
	switch config.ExName {
	case exchange.BINANCE:
		ex := binance.CreateBinance(config)
		e.exMan.Add(ex)
		return ex

	case exchange.BITTREX:
		ex := bittrex.CreateBittrex(config)
		e.exMan.Add(ex)
		return ex

	case exchange.COINEX:
		ex := coinex.CreateCoinex(config)
		e.exMan.Add(ex)
		return ex

	case exchange.STEX:
		ex := stex.CreateStex(config)
		e.exMan.Add(ex)
		return ex

	case exchange.BITMEX:
		ex := bitmex.CreateBitmex(config)
		e.exMan.Add(ex)
		return ex

	case exchange.KUCOIN:
		ex := kucoin.CreateKucoin(config)
		e.exMan.Add(ex)
		return ex

	case exchange.HUOBIOTC:
		ex := huobiotc.CreateHuobiOTC(config)
		e.exMan.Add(ex)
		return ex

	case exchange.BITMAX:
		ex := bitmax.CreateBitmax(config)
		e.exMan.Add(ex)
		return ex

	case exchange.BITSTAMP:
		ex := bitstamp.CreateBitstamp(config)
		e.exMan.Add(ex)
		return ex

	case exchange.OTCBTC:
		ex := otcbtc.CreateOtcbtc(config)
		e.exMan.Add(ex)
		return ex

	case exchange.HUOBI:
		ex := huobi.CreateHuobi(config)
		e.exMan.Add(ex)
		return ex

	case exchange.BIBOX:
		ex := bibox.CreateBibox(config)
		e.exMan.Add(ex)
		return ex

	case exchange.OKEX:
		ex := okex.CreateOkex(config)
		e.exMan.Add(ex)
		return ex

	case exchange.BITZ:
		ex := bitz.CreateBitz(config)
		e.exMan.Add(ex)
		return ex

	case exchange.HITBTC:
		ex := hitbtc.CreateHitbtc(config)
		e.exMan.Add(ex)
		return ex

	case exchange.DRAGONEX:
		ex := dragonex.CreateDragonex(config)
		e.exMan.Add(ex)
		return ex

	case exchange.BIGONE:
		ex := bigone.CreateBigone(config)
		e.exMan.Add(ex)
		return ex

	case exchange.BITFINEX:
		ex := bitfinex.CreateBitfinex(config)
		e.exMan.Add(ex)
		return ex

	case exchange.GATEIO:
		ex := gateio.CreateGateio(config)
		e.exMan.Add(ex)
		return ex

	case exchange.IDEX:
		ex := idex.CreateIdex(config)
		e.exMan.Add(ex)
		return ex

	case exchange.LIQUID:
		ex := liquid.CreateLiquid(config)
		e.exMan.Add(ex)
		return ex

	case exchange.BITFOREX:
		ex := bitforex.CreateBitforex(config)
		e.exMan.Add(ex)
		return ex

	case exchange.TOKOK:
		ex := tokok.CreateTokok(config)
		e.exMan.Add(ex)
		return ex

	case exchange.MXC:
		ex := mxc.CreateMxc(config)
		e.exMan.Add(ex)
		return ex

	case exchange.BITRUE:
		ex := bitrue.CreateBitrue(config)
		e.exMan.Add(ex)
		return ex

	}
	return nil
}
