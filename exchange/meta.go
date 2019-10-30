package exchange

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type ChainType string
type DataSource string
type ExchangeName string
type MarginAction string
type UpdateMethod string

const (
	MAINNET ChainType = "MAINNET"
	BEP2    ChainType = "BEP2"
	ERC20   ChainType = "ERC20"
	NEP5    ChainType = "NEP5"
	OMNI    ChainType = "OMNI"
	TRC20   ChainType = "TRC20"

	EXCHANGE_API     DataSource = "EXCHANGE_API"
	WEBSOCKET        DataSource = "WEBSOCKET"
	MICROSERVICE_API DataSource = "MICROSERVICE_API"
	JSON_FILE        DataSource = "JSON_FILE"
	PSQL             DataSource = "PSQL"

	TRANSFER_IN  MarginAction = "TRANSFER_IN"
	TRANSFER_OUT MarginAction = "TRANSFER_OUT"
	LOAN_REQUEST MarginAction = "LOAN_REQUEST"
	LOAN_REPAY   MarginAction = "LOAN_REPAY"
	ORDER_STATUS MarginAction = "ORDER_STATUS"
	BALANCE      MarginAction = "BALANCE"
	LIMIT_BUY    MarginAction = "LIMIT_BUY"
	LIMIT_SELL   MarginAction = "LIMIT_SELL"
	MARKET_BUY   MarginAction = "MARKET_BUY"
	MARKET_SELL  MarginAction = "MARKET_SELL"

	API_TIGGER  UpdateMethod = "API_TIGGER"
	TIME_TIGGER UpdateMethod = "TIME_TIGGER"

	ABCC         ExchangeName = "ABCC"
	BCEX         ExchangeName = "BCEX"
	BELFRIES     ExchangeName = "BELFRIES"
	BGOGO        ExchangeName = "BGOGO"
	BIBOX        ExchangeName = "BIBOX"
	BIGONE       ExchangeName = "BIGONE"
	BIKI         ExchangeName = "BIKI"
	BINANCE      ExchangeName = "BINANCE"
	BINANCEDEX   ExchangeName = "BINANCEDEX"
	BITATM       ExchangeName = "BITATM"
	BITBAY       ExchangeName = "BITBAY"
	BITCOIN      ExchangeName = "BITCOIN"
	BITFINEX     ExchangeName = "BITFINEX"
	BITFOREX     ExchangeName = "BITFOREX"
	BITHUMB      ExchangeName = "BITHUMB"
	BITMART      ExchangeName = "BITMART"
	BITMAX       ExchangeName = "BITMAX"
	BITMEX       ExchangeName = "BITMEX"
	BITSTAMP     ExchangeName = "BITSTAMP"
	BITTREX      ExchangeName = "BITTREX"
	BITLISH      ExchangeName = "BITLISH"
	BITRUE       ExchangeName = "BITRUE"
	BITZ         ExchangeName = "BITZ"
	BKEX         ExchangeName = "BKEX"
	BTSE         ExchangeName = "BTSE"
	BYBIT        ExchangeName = "BYBIT"
	BW           ExchangeName = "BW"
	BLANK        ExchangeName = "BLANK"
	BLEUTRADE    ExchangeName = "BLEUTRADE"
	BLOCKTRADE   ExchangeName = "BLOCKTRADE"
	COINALL      ExchangeName = "COINALL"
	COINMEX      ExchangeName = "COINMEX"
	COINBASE     ExchangeName = "COINBASE"
	COINBENE     ExchangeName = "COINBENE"
	COINEAL      ExchangeName = "COINEAL"
	COINEX       ExchangeName = "COINEX"
	COINSUPER    ExchangeName = "COINSUPER"
	COINTIGER    ExchangeName = "COINTIGER"
	COINDEAL     ExchangeName = "COINDEAL"
	CRYPTOPIA    ExchangeName = "CRYPTOPIA"
	DCOIN        ExchangeName = "DCOIN"
	DERIBIT      ExchangeName = "DERIBIT"
	DIGIFINEX    ExchangeName = "DIGIFINEX"
	DRAGONEX     ExchangeName = "DRAGONEX"
	EXMO         ExchangeName = "EXMO"
	EXX          ExchangeName = "EXX"
	FATBTC       ExchangeName = "FATBTC"
	FCOIN        ExchangeName = "FCOIN"
	FTX          ExchangeName = "FTX"
	GATEIO       ExchangeName = "GATEIO"
	GEMINI       ExchangeName = "GEMINI"
	GOKO         ExchangeName = "GOKO"
	GRAVIEX      ExchangeName = "GRAVIEX"
	HITBTC       ExchangeName = "HITBTC"
	HIBITEX      ExchangeName = "HIBITEX"
	HOTBIT       ExchangeName = "HOTBIT"
	HUOBI        ExchangeName = "HUOBI"
	HUOBIDM      ExchangeName = "HUOBIDM"
	HUOBIOTC     ExchangeName = "HUOBIOTC"
	IBANKDIGITAL ExchangeName = "IBANKDIGITAL"
	IDAX         ExchangeName = "IDAX"
	IDEX         ExchangeName = "IDEX"
	KRAKEN       ExchangeName = "KRAKEN"
	KUCOIN       ExchangeName = "KUCOIN"
	LATOKEN      ExchangeName = "LATOKEN"
	LBANK        ExchangeName = "LBANK"
	LIQUID       ExchangeName = "LIQUID"
	LIVECOIN     ExchangeName = "LIVECOIN"
	MXC          ExchangeName = "MXC"
	NEWCAPITAL   ExchangeName = "NEWCAPITAL"
	OKEX         ExchangeName = "OKEX"
	OKEXDM       ExchangeName = "OKEXDM"
	OTCBTC       ExchangeName = "OTCBTC"
	P2PB2B       ExchangeName = "P2PB2B"
	POLONIEX     ExchangeName = "POLONIEX"
	PROBIT       ExchangeName = "PROBIT"
	RIGHTBTC     ExchangeName = "RIGHTBTC"
	STEX         ExchangeName = "STEX"
	SWITCHEO     ExchangeName = "SWITCHEO"
	TOKOK        ExchangeName = "TOKOK"
	TIDEX        ExchangeName = "TIDEX"
	TOPBTC       ExchangeName = "TOPBTC"
	TRADEOGRE    ExchangeName = "TRADEOGRE"
	TRADESATOSHI ExchangeName = "TRADESATOSHI"
	TXBIT        ExchangeName = "TXBIT"
	UEX          ExchangeName = "UEX"
	VIRGOCX      ExchangeName = "VIRGOCX"
	ZBEX         ExchangeName = "ZBEX"
	ZEBITEX      ExchangeName = "ZEBITEX"
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
	supportList = append(supportList, HUOBIDM)   // ID = 42
	supportList = append(supportList, BW)        // ID = 43
	supportList = append(supportList, DERIBIT)   // ID = 44
	supportList = append(supportList, OKEXDM)    // ID = 45
	supportList = append(supportList, GOKO)      // ID = 46
	supportList = append(supportList, BCEX)      // ID = 47
	supportList = append(supportList, DIGIFINEX) // ID = 48
	supportList = append(supportList, LATOKEN)   // ID = 49
	supportList = append(supportList, VIRGOCX)   // ID = 50
	supportList = append(supportList, ABCC)      // ID = 51
	// supportList = append(supportList, BYBIT)     // ID = 52 no orderbook API
	supportList = append(supportList, ZEBITEX)    // ID = 53
	supportList = append(supportList, BITHUMB)    // ID = 54
	supportList = append(supportList, SWITCHEO)   // ID = 55
	supportList = append(supportList, BLOCKTRADE) // ID = 56
	supportList = append(supportList, BKEX)       // ID = 57
	supportList = append(supportList, NEWCAPITAL) // ID = 58
	supportList = append(supportList, COINDEAL)   // ID = 59
	supportList = append(supportList, HIBITEX)    // ID = 60
	supportList = append(supportList, BGOGO)      // ID = 61
	// supportList = append(supportList, FTX)        // ID = 62	orderbook not finished
	supportList = append(supportList, TXBIT)  // ID = 63
	supportList = append(supportList, PROBIT) // ID = 64
}
