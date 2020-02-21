package bittrex

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/utils"
)

func (e *Bittrex) LoadPublicData(operation *exchange.PublicOperation) error {
	switch operation.Type {

	case exchange.TradeHistory:
		return e.doTradeHistory(operation)

	}
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

func (e *Bittrex) doTradeHistory(operation *exchange.PublicOperation) error {
	symbol := e.GetSymbolByPair(operation.Pair)

	get := &utils.HttpGet{
		URI: fmt.Sprintf("%s/v1.1/public/getmarkethistory?market=%s", API_URL, symbol),
	}

	err := utils.HttpGetRequest(get)

	if err != nil {
		return err

	} else {
		jsonResponse := JsonResponse{}
		tradeHistory := &TradeHistory{}
		if err := json.Unmarshal(get.ResponseBody, &jsonResponse); err != nil {
			return err
		}

		if err := json.Unmarshal(jsonResponse.Result, &tradeHistory); err != nil {
			return err
		}

		operation.TradeHistory = []*exchange.TradeDetail{}
		for _, d := range *tradeHistory {
			td := &exchange.TradeDetail{}

			td.ID = fmt.Sprintf("%d", d.ID)
			if d.OrderType == "BUY" {
				td.Direction = exchange.Buy
			} else if d.OrderType == "SELL" {
				td.Direction = exchange.Sell
			}

			td.Quantity = d.Quantity
			td.Rate = d.Price

			layout := "2006-01-02T15:04:05.00"
			t, err := time.Parse(layout, d.TimeStamp)
			if err != nil {
				// log.Printf("%+v", err)
				layout = "2006-01-02T15:04:05.0"
				t, err = time.Parse(layout, d.TimeStamp)
				if err != nil {
					// log.Printf("%+v", err)
					layout = "2006-01-02T15:04:05"
					t, err = time.Parse(layout, d.TimeStamp)
				}
			}

			td.TimeStamp = t.UnixNano() / 1e6

			operation.TradeHistory = append(operation.TradeHistory, td)
		}
	}

	return nil
}
