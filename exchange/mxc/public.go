package mxc

import (
	"fmt"

	"github.com/bitontop/gored/exchange"
)

func (e *Mxc) LoadPublicData(operation *exchange.PublicOperation) error {
	switch operation.Type {
	case exchange.TradeHistory:
		return e.doTradeHistory(operation)

	}
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

// timeStamp parse TODO
func (e *Mxc) doTradeHistory(operation *exchange.PublicOperation) error {
	// symbol := e.GetSymbolByPair(operation.Pair)

	// get := &utils.HttpGet{
	// 	URI: fmt.Sprintf("%s/open/api/v1/data/history?market=%s", API_URL, symbol),
	// }

	// err := utils.HttpGetRequest(get)

	// if err != nil {
	// 	return err

	// } else {
	// 	jsonResponse := JsonResponse{}
	// 	tradeHistory := TradeHistory{}
	// 	if err := json.Unmarshal(get.ResponseBody, &jsonResponse); err != nil {
	// 		return err
	// 	} else if jsonResponse.Code != 200 {
	// 		return err
	// 	}

	// 	if err := json.Unmarshal(jsonResponse.Data, &tradeHistory); err != nil {
	// 		return err
	// 	}

	// 	operation.TradeHistory = []*exchange.TradeDetail{}
	// 	// for _, d := range tradeHistory {
	// 	for i := len(tradeHistory) - 1; i > 0; i-- {
	// 		d := tradeHistory[i]
	// 		td := &exchange.TradeDetail{}

	// 		// td.ID = fmt.Sprintf("%d", d.ID)
	// 		if d.TradeType == "1" {
	// 			td.Direction = exchange.Buy
	// 		} else if d.TradeType == "2" {
	// 			td.Direction = exchange.Sell
	// 		}

	// 		td.Quantity, err = strconv.ParseFloat(d.TradeQuantity, 64)
	// 		td.Rate, err = strconv.ParseFloat(d.TradePrice, 64)

	// 		// TODO
	// 		layout := "2006-03-13 07:21:27.788" //"2006-01-02 15:04:05"
	// 		ts, _ := time.Parse(layout, d.TradeTime)
	// 		td.TimeStamp = ts.UnixNano()
	// 		td.ID = fmt.Sprintf("%v", td.TimeStamp)

	// 		operation.TradeHistory = append(operation.TradeHistory, td)
	// 	}
	// }

	return nil
}
