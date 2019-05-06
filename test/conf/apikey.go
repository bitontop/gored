package conf

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"../../exchange"
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

		// case exchange.HUOBI:
		// 	config.API_KEY = ""
		// 	config.API_SECRET = ""
		// 	break

	}
}
