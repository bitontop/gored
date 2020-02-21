package hitbtc

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/bitontop/gored/exchange"
)

func (e *Hitbtc) LoadPublicData(operation *exchange.PublicOperation) error {
	switch operation.Type {

	case exchange.TradeHistory:
		return e.doTradeHistory(operation)

	}
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

// Authorization required
func (e *Hitbtc) doTradeHistory(operation *exchange.PublicOperation) error {
	strRequest := "/api/2/history/trades"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(operation.Pair)

	jsonTradeHistoryReturn := e.ApiKeyRequest("POST", mapParams, strRequest)

	tradeHistory := &TradeHistory{}
	if err := json.Unmarshal([]byte(jsonTradeHistoryReturn), &tradeHistory); err != nil {
		return err
	}

	operation.TradeHistory = []*exchange.TradeDetail{}
	for _, d := range *tradeHistory {
		td := &exchange.TradeDetail{}

		td.ID = fmt.Sprintf("%d", d.ID)
		if d.Side == "buy" {
			td.Direction = exchange.Buy
		} else if d.Side == "sell" {
			td.Direction = exchange.Sell
		}

		td.Quantity, _ = strconv.ParseFloat(d.Quantity, 64)
		td.Rate, _ = strconv.ParseFloat(d.Price, 64)

		td.TimeStamp = d.Timestamp.UnixNano() / 1e6

		operation.TradeHistory = append(operation.TradeHistory, td)
	}

	return nil
}
