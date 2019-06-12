package exchange

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type ExchangeName string
type DataSource string
type ChainType string
type UpdateMethod string

const (
	API_TIGGER  UpdateMethod = "API_TIGGER"
	TIME_TIGGER UpdateMethod = "TIME_TIGGER"

	EXCHANGE_API     DataSource = "EXCHANGE_API"
	MICROSERVICE_API DataSource = "MICROSERVICE_API"
	JSON_FILE        DataSource = "JSON_FILE"
	PSQL             DataSource = "PSQL"

	MAINNET ChainType = "MAINNET"
	BEP2    ChainType = "BEP2"
	ERC20   ChainType = "ERC20"
	NEP5    ChainType = "NEP5"
	OMNI    ChainType = "OMNI"
	TRC20   ChainType = "TRC20"

	BCEX         ExchangeName = "BCEX"
	BGOGO        ExchangeName = "BGOGO"
	BIBOX        ExchangeName = "BIBOX"
	BIGONE       ExchangeName = "BIGONE"
	BIKI         ExchangeName = "BIKI"
	BINANCE      ExchangeName = "BINANCE"
	BINANCEDEX   ExchangeName = "BINANCEDEX"
	BITATM       ExchangeName = "BITATM"
	BITBAY       ExchangeName = "BITBAY"
	BITFINEX     ExchangeName = "BITFINEX"
	BITFOREX     ExchangeName = "BITFOREX"
	BITMART      ExchangeName = "BITMART"
	BITMAX       ExchangeName = "BITMAX"
	BITMEX       ExchangeName = "BITMEX"
	BITSTAMP     ExchangeName = "BITSTAMP"
	BITTREX      ExchangeName = "BITTREX"
	BITLISH      ExchangeName = "BITLISH"
	BITRUE       ExchangeName = "BITRUE"
	BITZ         ExchangeName = "BITZ"
	BLANK        ExchangeName = "BLANK"
	BLEUTRADE    ExchangeName = "BLEUTRADE"
	COINALL      ExchangeName = "COINALL"
	COINMEX      ExchangeName = "COINMEX"
	COINBASE     ExchangeName = "COINBASE"
	COINBENE     ExchangeName = "COINBENE"
	COINEAL      ExchangeName = "COINEAL"
	COINEX       ExchangeName = "COINEX"
	COINSUPER    ExchangeName = "COINSUPER"
	COINTIGER    ExchangeName = "COINTIGER"
	CRYPTOPIA    ExchangeName = "CRYPTOPIA"
	DCOIN        ExchangeName = "DCOIN"
	DIGIFINEX    ExchangeName = "DIGIFINEX"
	DRAGONEX     ExchangeName = "DRAGONEX"
	EXMO         ExchangeName = "EXMO"
	EXX          ExchangeName = "EXX"
	FATBTC       ExchangeName = "FATBTC"
	FCOIN        ExchangeName = "FCOIN"
	GATEIO       ExchangeName = "GATEIO"
	GEMINI       ExchangeName = "GEMINI"
	GRAVIEX      ExchangeName = "GRAVIEX"
	HITBTC       ExchangeName = "HITBTC"
	HOTBIT       ExchangeName = "HOTBIT"
	HUOBI        ExchangeName = "HUOBI"
	HUOBIDM      ExchangeName = "HUOBIDM"
	HUOBIOTC     ExchangeName = "HUOBIOTC"
	IBANKDIGITAL ExchangeName = "IBANKDIGITAL"
	IDAX         ExchangeName = "IDAX"
	IDEX         ExchangeName = "IDEX"
	KRAKEN       ExchangeName = "KRAKEN"
	KUCOIN       ExchangeName = "KUCOIN"
	LBANK        ExchangeName = "LBANK"
	LIQUID       ExchangeName = "LIQUID"
	LIVECOIN     ExchangeName = "LIVECOIN"
	MXC          ExchangeName = "MXC"
	OKEX         ExchangeName = "OKEX"
	OTCBTC       ExchangeName = "OTCBTC"
	P2PB2B       ExchangeName = "P2PB2B"
	POLONIEX     ExchangeName = "POLONIEX"
	RIGHTBTC     ExchangeName = "RIGHTBTC"
	STEX         ExchangeName = "STEX"
	TOKOK        ExchangeName = "TOKOK"
	TIDEX        ExchangeName = "TIDEX"
	TOPBTC       ExchangeName = "TOPBTC"
	TRADEOGRE    ExchangeName = "TRADEOGRE"
	TRADESATOSHI ExchangeName = "TRADESATOSHI"
	UEX          ExchangeName = "UEX"
	ZBEX         ExchangeName = "ZBEX"
)

func (e *ExchangeManager) initExchangeNames() {
	supportList = append(supportList, BINANCE)      // ID = 1
	supportList = append(supportList, BITTREX)      // ID = 2
	supportList = append(supportList, COINEX)       // ID = 3
	supportList = append(supportList, STEX)         // ID = 4
	supportList = append(supportList, BITMEX)       // ID = 5
	supportList = append(supportList, KUCOIN)       // ID = 6
	supportList = append(supportList, BITMAX)       // ID = 7
	supportList = append(supportList, HUOBIOTC)     // ID = 8
	supportList = append(supportList, BITSTAMP)     // ID = 9
	supportList = append(supportList, OTCBTC)       // ID = 10
	supportList = append(supportList, HUOBI)        // ID = 11
	supportList = append(supportList, BIBOX)        // ID = 12
	supportList = append(supportList, OKEX)         // ID = 13
	supportList = append(supportList, BITZ)         // ID = 14
	supportList = append(supportList, HITBTC)       // ID = 15
	supportList = append(supportList, DRAGONEX)     // ID = 16
	supportList = append(supportList, BIGONE)       // ID = 17
	supportList = append(supportList, BITFINEX)     // ID = 18
	supportList = append(supportList, GATEIO)       // ID = 19
	supportList = append(supportList, IDEX)         // ID = 20
	supportList = append(supportList, LIQUID)       // ID = 21
	supportList = append(supportList, BITFOREX)     // ID = 22
	supportList = append(supportList, TOKOK)        // ID = 23
	supportList = append(supportList, MXC)          // ID = 24
	supportList = append(supportList, BITRUE)       // ID = 25
	supportList = append(supportList, BITATM)       // ID = 26
	supportList = append(supportList, TRADESATOSHI) // ID = 27
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
	supportList = append(supportList, HUOBIDM)   // ID = 99
}
