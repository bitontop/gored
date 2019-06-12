package conf

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"github.com/bitontop/gored/exchange"
)

func Exchange(name exchange.ExchangeName, config *exchange.Config) {
	config.ExName = name
	switch name {
	case exchange.BINANCE:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.BITTREX:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.COINEX:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.STEX:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.BITMEX:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.KUCOIN:
		config.API_KEY = ""
		config.API_SECRET = ""
		config.Passphrase = ""
		break

	case exchange.BITMAX:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.BITSTAMP:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.OTCBTC:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.HUOBI:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.BIBOX:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.OKEX:
		config.API_KEY = ""
		config.API_SECRET = ""
		config.Passphrase = ""
		config.TradePassword = ""
		break

	case exchange.BITZ:
		config.API_KEY = ""
		config.API_SECRET = ""
		config.TradePassword = ""
		break

	case exchange.HITBTC:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.DRAGONEX:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.BIGONE:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.BITFINEX:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.GATEIO:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.IDEX:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.LIQUID:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.BITFOREX:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.TOKOK:
		config.API_KEY = "5ecac940-b495-48c6-9bd1-dd0e1c6b6b95"
		config.API_SECRET = "14a9a62e-fcaf-49a8-a4a9-dbaf5921bbff"
		break

	case exchange.MXC:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.BITRUE:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.TRADESATOSHI:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.KRAKEN:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.POLONIEX:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.COINEAL:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.TRADEOGRE:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.COINBENE:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.IBANKDIGITAL:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.LBANK:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.BITMART:
		config.API_KEY = "3e43d35bbeefeb881b4de213dc01042f"
		config.API_SECRET = "fa65bc32e69c1b7f2f8965ca7cf5d4fa"
		config.Passphrase = "key3" // key name
		break

	case exchange.BIKI:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.BITATM:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.DCOIN:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.GEMINI:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.COINTIGER:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.BITBAY:
		config.API_KEY = ""
		config.API_SECRET = ""
		break
	}
}
