package initial

import (
	"sync"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/exchange/bibox"
	"github.com/bitontop/gored/exchange/bigone"
	"github.com/bitontop/gored/exchange/binance"
	"github.com/bitontop/gored/exchange/bitfinex"
	"github.com/bitontop/gored/exchange/bitmax"
	"github.com/bitontop/gored/exchange/bitmex"
	"github.com/bitontop/gored/exchange/bittrex"
	"github.com/bitontop/gored/exchange/bitz"
	"github.com/bitontop/gored/exchange/coinex"
	"github.com/bitontop/gored/exchange/dragonex"
	"github.com/bitontop/gored/exchange/gateio"
	"github.com/bitontop/gored/exchange/hitbtc"
	"github.com/bitontop/gored/exchange/huobi"
	"github.com/bitontop/gored/exchange/huobiotc"
	"github.com/bitontop/gored/exchange/kucoin"
	"github.com/bitontop/gored/exchange/okex"
	"github.com/bitontop/gored/exchange/otcbtc"
	"github.com/bitontop/gored/exchange/stex"

	"github.com/bitontop/gored/exchange/bitstamp"
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

	case exchange.GATEIO:
		ex := gateio.CreateGateio(config)
		e.exMan.Add(ex)
		return ex

		// case exchange.TRADESATOSHI:
		// 	ex := tradesatoshi.CreateTradeSatoshi(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.BITFINEX:
		// 	ex := bitfinex.CreateBitfinex(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.OKEX:
		// 	ex := okex.CreateOkex(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.HITBTC:
		// 	ex := hitbtc.CreateHitbtc(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.CRYPTOPIA:
		// 	ex := cryptopia.CreateCryptopia(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.ZBEX:
		// 	ex := zbex.CreateZBex(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.POLONIEX:
		// 	ex := poloniex.CreatePoloniex(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.GATEIO:
		// 	ex := gateio.CreateGateIo(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.LBANK:
		// 	ex := lbank.CreateLbank(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.COINBENE:
		// 	ex := coinbene.CreateCoinbene(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.KRAKEN:
		// 	ex := kraken.CreateKraken(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.BCEX:
		// 	ex := bcex.CreateBcex(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.IDAX:
		// 	ex := idax.CreateIdax(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.BITZ:
		// 	ex := bitz.CreateBitz(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.TRADEOGRE:
		// 	ex := tradeogre.CreateTradeogre(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.DIGIFINEX:
		// 	ex := digifinex.CreateDigifinex(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.EXX:
		// 	ex := exx.CreateExx(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.TOPBTC:
		// 	ex := topbtc.CreateTopbtc(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.BIGONE:
		// 	ex := bigone.CreateBigone(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.FATBTC:
		// 	ex := fatbtc.CreateFatbtc(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.BITMART:
		// 	ex := bitmart.CreateBitmart(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.BIBOX:
		// 	ex := bibox.CreateBibox(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.BITFOREX:
		// 	ex := bitforex.CreateBitforex(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.TOKOK:
		// 	ex := tokok.CreateTokok(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.LIQUID:
		// 	ex := liquid.CreateLiquid(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.COINEAL:
		// 	ex := coineal.CreateCoineal(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.COINMEX:
		// 	ex := coinmex.CreateCoinmex(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.IBANKDIGITAL:
		// 	ex := ibankdigital.CreateIbankdigital(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.HOTBIT:
		// 	ex := hotbit.CreateHotbit(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.UEX:
		// 	ex := uex.CreateUex(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.BLEUTRADE:
		// 	ex := bleutrade.CreateBleutrade(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.DRAGONEX:
		// 	ex := dragonex.CreateDragonex(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.BGOGO:
		// 	ex := bgogo.CreateBgogo(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.IDEX:
		// 	ex := idex.CreateIdex(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.BITRUE:
		// 	ex := bitrue.CreateBitrue(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.MXC:
		// 	ex := mxc.CreateMxc(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.GRAVIEX:
		// 	ex := graviex.CreateGraviex(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.P2PB2B:
		// 	ex := p2pb2b.CreateP2pb2b(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.COINBASE:
		// 	ex := coinbase.CreateCoinbase(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.LIVECOIN:
		// 	ex := livecoin.CreateLivecoin(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.BIKI:
		// 	ex := biki.CreateBiki(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.TIDEX:
		// 	ex := tidex.CreateTidex(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.BITLISH:
		// 	ex := bitlish.CreateBitlish(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.EXMO:
		// 	ex := exmo.CreateExmo(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.RIGHTBTC:
		// 	ex := rightbtc.CreateRightbtc(config)
		// 	e.exMan.Add(ex)
		// 	return ex

		// case exchange.FCOIN:
		// 	ex := fcoin.CreateFcoin(config)
		// 	e.exMan.Add(ex)
		// 	return ex
	case exchange.BITFINEX:
		ex := bitfinex.CreateBitfinex(config)
		e.exMan.Add(ex)
		return ex
	}
	return nil
}
