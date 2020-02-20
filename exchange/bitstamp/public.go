package bitstamp

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/utils"
)

func (e *Bitstamp) LoadPublicData(operation *exchange.PublicOperation) error {
	switch operation.Type {

	case exchange.TradeHistory:
		return e.doTradeHistory(operation)

	}
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

func (e *Bitstamp) doTradeHistory(operation *exchange.PublicOperation) error {
	symbol := e.GetSymbolByPair(operation.Pair)

	get := &utils.HttpGet{
		URI: fmt.Sprintf("%s/transactions/%s", API_URL, symbol),
	}

	err := utils.HttpGetRequest(get)

	if err != nil {
		return err

	} else {
		tradeHistory := &TradeHistory{}
		if err := json.Unmarshal(get.ResponseBody, &tradeHistory); err != nil {
			return err
		}

		operation.TradeHistory = []*exchange.TradeDetail{}
		for _, d := range *tradeHistory {
			td := &exchange.TradeDetail{}

			if d.Type == "0" {
				td.Direction = exchange.Buy
			} else if d.Type == "1" {
				td.Direction = exchange.Sell
			}

			td.Quantity, err = strconv.ParseFloat(d.Amount, 64)
			td.Rate, err = strconv.ParseFloat(d.Price, 64)

			t, err := strconv.ParseInt(d.Date, 10, 64)
			if err != nil {
				return err
			}
			td.TimeStamp = t * 1000

			operation.TradeHistory = append(operation.TradeHistory, td)
		}
	}

	return nil
}
