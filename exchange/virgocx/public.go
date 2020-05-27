package virgocx

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/utils"
)

func (e *Virgocx) LoadPublicData(operation *exchange.PublicOperation) error {
	switch operation.Type {
	case exchange.KLine:
		return e.doKline(operation)
	}
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

// interval: 1, 5, 10, 30, 60, 240, 1440, 7200, 10080, 43200
func (e *Virgocx) doKline(operation *exchange.PublicOperation) error {
	interval := "5"
	if operation.KlineInterval != "" {
		interval = operation.KlineInterval
	}

	get := &utils.HttpGet{
		URI: fmt.Sprintf("https://www.virgocx.ca/api/market/history/kline?symbol=%v&period=%v",
			e.GetSymbolByPair(operation.Pair), // BTC/CAD
			interval,
		),
		Proxy: operation.Proxy,
	}

	err := utils.HttpGetRequest(get)

	if err != nil {
		log.Printf("%+v", err)
		operation.Error = err
		return err

	}

	if operation.DebugMode {
		operation.RequestURI = get.URI
		operation.CallResponce = string(get.ResponseBody)
	}

	jsonResponse := JsonResponse{}
	rawKline := RawKline{}

	if err := json.Unmarshal([]byte(get.ResponseBody), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s doKline Json Unmarshal Err: %v %v", e.GetName(), err, string(get.ResponseBody))
		return operation.Error
	} else if jsonResponse.Code != 0 {
		operation.Error = fmt.Errorf("%s doKline Failed: %v", e.GetName(), string(get.ResponseBody))
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &rawKline); err != nil {
		operation.Error = fmt.Errorf("%s doKline Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	operation.Kline = []*exchange.KlineDetail{}
	for _, k := range rawKline {
		detail := &exchange.KlineDetail{
			OpenTime: float64(k.CreateTime),
			Open:     k.Open,
			High:     k.High,
			Low:      k.Low,
			Close:    k.Close,
			Volume:   k.Qty,
		}

		operation.Kline = append(operation.Kline, detail)
	}

	return nil
}
