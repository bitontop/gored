package homiex

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/utils"
)

func (e *Homiex) LoadPublicData(operation *exchange.PublicOperation) error {
	switch operation.Type {
	case exchange.TradeHistory:
		return e.doTradeHistory(operation)

	}
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

// using timeStamp as tradeID
func (e *Homiex) doTradeHistory(operation *exchange.PublicOperation) error {
	symbol := e.GetSymbolByPair(operation.Pair)

	get := &utils.HttpGet{
		URI: fmt.Sprintf("%s/openapi/quote/v1/trades?symbol=%s&limit=%v", API_URL, symbol, 1000),
	}

	err := utils.HttpGetRequest(get)

	if err != nil {
		return err

	} else {
		tradeHistory := TradeHistory{}
		if err := json.Unmarshal(get.ResponseBody, &tradeHistory); err != nil {
			return err
		} else if len(tradeHistory) == 0 {
			return fmt.Errorf("%v TradeHistory got 0 record, %v", e.GetName(), string(get.ResponseBody))
		}

		operation.TradeHistory = []*exchange.TradeDetail{}
		for _, trade := range tradeHistory {
			td := &exchange.TradeDetail{}

			td.ID = fmt.Sprintf("%d", trade.Time)
			if trade.IsBuyerMaker {
				td.Direction = exchange.Buy
			} else {
				td.Direction = exchange.Sell
			}

			td.Quantity, err = strconv.ParseFloat(trade.Qty, 64)
			td.Rate, err = strconv.ParseFloat(trade.Price, 64)
			if err != nil {
				return fmt.Errorf("%v TradeHistory err: %v", e.GetName(), err)
			}

			td.TimeStamp = trade.Time

			operation.TradeHistory = append(operation.TradeHistory, td)
		}
	}

	return nil
}
