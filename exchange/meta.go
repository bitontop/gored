package exchange

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type ChainType string
type DataSource string
type ExchangeName string
type MarginAction string
type UpdateMethod string

// from goredmergin
type ContractAction string
type ContractTransDir string
type OffSetType string
type OrderPriceType string

const (
	MAINNET ChainType = "MAINNET"
	BEP2    ChainType = "BEP2" //TODO //
	ERC20   ChainType = "ERC20"
	NEP5    ChainType = "NEP5" //NEO //
	OMNI    ChainType = "OMNI"
	TRC20   ChainType = "TRC20"
	BNB     ChainType = "BNB" //
	CET     ChainType = "CET" //
	NXT     ChainType = "NXT" //
	OTHER   ChainType = "OTHER"

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

	// ************ from goredmergin ************
	CONTRACT_MARKET_BUY   ContractAction = "CONTRACT_MARKET_BUY"
	CONTRACT_MARKET_SELL  ContractAction = "CONTRACT_MARKET_SELL"
	CONTRACT_LIMIT_BUY    ContractAction = "CONTRACT_LIMIT_BUY"
	CONTRACT_LIMIT_SELL   ContractAction = "CONTRACT_LIMIT_SELL"
	GET_ADDR              ContractAction = "GET_ADDR"
	LIQUIDATION           ContractAction = "LIQUIDATION"
	CONTRACT_TRANSFER     ContractAction = "CONTRACT_TRANSFER"
	CONTRACT_ORDER_STATUS ContractAction = "CONTRACT_ORDER_STATUS"
	CONTRACT_BALANCE      ContractAction = "CONTRACT_BALANCE"

	OPEN  OffSetType = "open"
	CLOSE OffSetType = "close"

	// Huobi, Bitmex Contract
	LIMIT          OrderPriceType = "limit"
	OPTIMAL_5_FOK  OrderPriceType = "optimal_5_fok"
	OPTIMAL_5      OrderPriceType = "optimal_5"
	OPTIMAL_10     OrderPriceType = "optimal_10"
	OPTIMAL_20     OrderPriceType = "optimal_20"
	OPTIMAL_20_FOK OrderPriceType = "optimal_20_fok"
	BBO            OrderPriceType = "opponent"
	BBO_FOK        OrderPriceType = "opponent_fok"

	FUTURE_TO_SPOT ContractTransDir = "futures-to-pro"
	SPOT_TO_FUTURE ContractTransDir = "pro-to-futures"

	// Binance Contract
	GTC OrderPriceType = "GTC"
	IOC OrderPriceType = "IOC"
	FOK OrderPriceType = "FOK"
	GTX OrderPriceType = "GTX"

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
	BITBNS       ExchangeName = "BITBNS"
	BITCOIN      ExchangeName = "BITCOIN"
	BITFINEX     ExchangeName = "BITFINEX"
	BITFOREX     ExchangeName = "BITFOREX"
	BITHUMB      ExchangeName = "BITHUMB"
	BITMART      ExchangeName = "BITMART"
	BITMAX       ExchangeName = "BITMAX"
	BITMEX       ExchangeName = "BITMEX"
	BITPIE       ExchangeName = "BITPIE"
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
	HOMIEX       ExchangeName = "HOMIEX"
	HOO          ExchangeName = "HOO"
	HOTBIT       ExchangeName = "HOTBIT"
	HUOBI        ExchangeName = "HUOBI"
	HUOBIDM      ExchangeName = "HUOBIDM"
	HUOBIOTC     ExchangeName = "HUOBIOTC"
	IBANKDIGITAL ExchangeName = "IBANKDIGITAL"
	IDAX         ExchangeName = "IDAX"
	IDEX         ExchangeName = "IDEX"
	IDCM         ExchangeName = "IDCM"
	KRAKEN       ExchangeName = "KRAKEN"
	KUCOIN       ExchangeName = "KUCOIN"
	LATOKEN      ExchangeName = "LATOKEN"
	LBANK        ExchangeName = "LBANK"
	LIQUID       ExchangeName = "LIQUID"
	LIVECOIN     ExchangeName = "LIVECOIN"
	MXC          ExchangeName = "MXC"
	NEWCAPITAL   ExchangeName = "NEWCAPITAL"
	NICEHASH     ExchangeName = "NICEHASH"
	OKEX         ExchangeName = "OKEX"
	OKEXDM       ExchangeName = "OKEXDM"
	OKSIM        ExchangeName = "OKSIM"
	OTCBTC       ExchangeName = "OTCBTC"
	P2PB2B       ExchangeName = "P2PB2B"
	POLONIEX     ExchangeName = "POLONIEX"
	PROBIT       ExchangeName = "PROBIT"
	RIGHTBTC     ExchangeName = "RIGHTBTC"
	STEX         ExchangeName = "STEX"
	SWITCHEO     ExchangeName = "SWITCHEO"
	TAGZ         ExchangeName = "TAGZ"
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

type TradeDirection string
type OrderType string //deprecated  after 2020
const (
	Buy  TradeDirection = "b"
	Sell TradeDirection = "s"
	BUY  OrderType      = "BUY"  //deprecated  after 2020
	SELL OrderType      = "SELL" //deprecated  after 2020
)

type OrderTradeType string

const (
	TRADE_LIMIT  OrderTradeType = "LIMIT"
	TRADE_MARKET OrderTradeType = "MARKET"
	// Stop order need 'StopRate' param
	Trade_STOP_LIMIT  OrderTradeType = "STOP"
	Trade_STOP_MARKET OrderTradeType = "STOP_MARKET"
	//TODO Types below need more params,
	Trade_TAKE_PROFIT          OrderTradeType = "TAKE_PROFIT"
	Trade_TAKE_PROFIT_MARKET   OrderTradeType = "TAKE_PROFIT_MARKET"
	Trade_TRAILING_STOP_MARKET OrderTradeType = "TRAILING_STOP_MARKET"
)
