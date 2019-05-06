package exchange

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type ExchangeName string
type DataSource string
type UpdateMethod string

const (
	API_TIGGER  UpdateMethod = "API_TIGGER"
	TIME_TIGGER UpdateMethod = "TIME_TIGGER"

	EXCHANGE_API     DataSource = "EXCHANGE_API"
	MICROSERVICE_API DataSource = "MICROSERVICE_API"
	JSON_FILE        DataSource = "JSON_FILE"
	PSQL             DataSource = "PSQL"

	BINANCE ExchangeName = "BINANCE"
	BITTREX ExchangeName = "BITTREX"
	COINEX ExchangeName = "COINEX"
)

func (e *ExchangeManager) initExchangeNames() {
	supportList = append(supportList, BINANCE)
	supportList = append(supportList, BITTREX)
	supportList = append(supportList, COINEX)
}
