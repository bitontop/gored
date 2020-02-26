package bitfinex

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/utils"
)

func (e *Bitfinex) LoadPublicData(operation *exchange.PublicOperation) error {
	switch operation.Type {

	case exchange.TradeHistory:
		return e.doTradeHistory(operation)

	}
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

// Symbol that is not same as API.
func (e *Bitfinex) doTradeHistory(operation *exchange.PublicOperation) error {
	symbol := e.GetSymbolByPair(operation.Pair)
	symbol = fmt.Sprintf("t%v", strings.ToUpper(symbol))
	// log.Printf("Symbol: %s", symbol)

	get := &utils.HttpGet{
		URI: fmt.Sprintf("%s/v2/trades/%s/hist", API_URL, symbol),
	}

	err := utils.HttpGetRequest(get)

	if err != nil {
		return err

	} else {
		tradeHistory := TradeHistory{}
		if err := json.Unmarshal(get.ResponseBody, &tradeHistory); err != nil {
			return err
		}

		operation.TradeHistory = []*exchange.TradeDetail{}
		// for _, d := range tradeHistory {
		for i := len(tradeHistory) - 1; i > 0; i-- {
			d := tradeHistory[i]
			td := &exchange.TradeDetail{}

			td.ID = fmt.Sprintf("%.0f", d[0])
			if d[2] > 0 {
				td.Direction = exchange.Buy
				td.Quantity = d[2]
			} else if d[2] < 0 {
				td.Direction = exchange.Sell
				td.Quantity = d[2] * -1
			}

			td.Rate = d[3]

			td.TimeStamp = int64(d[1])

			operation.TradeHistory = append(operation.TradeHistory, td)
		}
	}

	return nil
}
