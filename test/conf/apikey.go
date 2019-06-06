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
		config.API_KEY = ""
		config.API_SECRET = ""
		config.Passphrase = "" // key name
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
		config.API_KEY = "1dbi4tegcy2oe3qfte6c84bcrywsm8f2"
		config.API_SECRET = "1p3k53wtebqa129wht01sr6raq0q6ugb"
		break

	case exchange.GEMINI:
		config.API_KEY = ""
		config.API_SECRET = ""
		break

	case exchange.COINTIGER:
		config.API_KEY = "fe6d48e4-cf08-4566-b217-3ca8e8c4014c"
		config.API_SECRET = "NDRhYWViNmE3YzY2OGMzM2FhMjUyZjUxODdjNDBjNTE4NTczOTEyMzcxYmQ1MDg2OGUyZDI4ZGU2MGZhZWMzOA=="
		break
	}
}
