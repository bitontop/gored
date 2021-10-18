package ftx

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/bitontop/gored/exchange"
	"github.com/bitontop/gored/utils"
)

/*************** PUBLIC  API ***************/
func (e *Ftx) LoadPublicData(operation *exchange.PublicOperation) error {
	switch operation.Type {
	case exchange.KLine:
		switch operation.Wallet {
		case exchange.SpotWallet:
			return e.doSpotKline(operation)
		}
	case exchange.GetTickerPrice:
		switch operation.Wallet {
		case exchange.SpotWallet:
			return e.doTickerPrice(operation)
		}

	case exchange.GetFutureStats:
		switch operation.Wallet {
		case exchange.SpotWallet: //actually future works the same for FTX
			return e.doGetFutureStats(operation)
		}

	}

	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

func (e *Ftx) doTickerPrice(operation *exchange.PublicOperation) error {
	jsonResponse := &JsonResponse{}
	tickerPrice := PairsData{}

	get := &utils.HttpGet{
		URI:       fmt.Sprintf("%s/api/markets", API_URL),
		Proxy:     operation.Proxy,
		DebugMode: operation.DebugMode,
	}
	if err := utils.HttpGetRequest(get); err != nil {
		operation.Error = err
		return operation.Error
	}

	if operation.DebugMode {
		operation.RequestURI = get.URI
		operation.CallResponce = string(get.ResponseBody)
	}

	jsonTickerPrice := get.ResponseBody
	if err := json.Unmarshal([]byte(jsonTickerPrice), &jsonResponse); err != nil {
		return fmt.Errorf("%s doTickerPrice Json Unmarshal Err: %v %v", e.GetName(), err, jsonTickerPrice)
	} else if !jsonResponse.Success {
		return fmt.Errorf("%s doTickerPrice Failed: %v", e.GetName(), jsonTickerPrice)
	}
	if err := json.Unmarshal(jsonResponse.Result, &tickerPrice); err != nil {
		return fmt.Errorf("%s doTickerPrice Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	operation.TickerPrice = []*exchange.TickerPriceDetail{}
	for _, tp := range tickerPrice {
		if tp.Type != "spot" {
			continue
		}
		p := e.GetPairBySymbol(tp.Name)

		if p == nil {
			if operation.DebugMode {
				log.Printf("doTickerPrice got nil pair for symbol: %v", tp.Name)
			}
			continue
		} else if p.Name == "" {
			continue
		}

		tpd := &exchange.TickerPriceDetail{
			Pair:  p,
			Price: tp.Price,
		}

		operation.TickerPrice = append(operation.TickerPrice, tpd)
	}

	return nil
}

// interval options: 15s, 1min, 5min, 15min, 1hour, 4hour, 1day
func (e *Ftx) doSpotKline(operation *exchange.PublicOperation) error {
	interval := "300"
	if operation.KlineInterval != "" {
		switch operation.KlineInterval {
		case "15s":
			interval = "15"
		case "1min":
			interval = "60"
		case "5min":
			interval = "300"
		case "15min":
			interval = "900"
		case "1hour":
			interval = "3600"
		case "4hour":
			interval = "14400"
		case "1day":
			interval = "86400"
		}
	}

	get := &utils.HttpGet{
		URI: fmt.Sprintf("https://ftx.com/api/markets/%v/candles?resolution=%v&limit=5000", // 1500478320000
			e.GetSymbolByPair(operation.Pair), // BTC/USD
			interval,
		),
		Proxy: operation.Proxy,
	}

	if operation.KlineStartTime != 0 {
		get.URI += fmt.Sprintf("&start_time=%v", operation.KlineStartTime/1000)
	}
	if operation.KlineEndTime != 0 {
		get.URI += fmt.Sprintf("&end_time=%v", operation.KlineEndTime/1000)
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

	jsonResponse := &JsonResponse{}
	rawKline := RawKline{}
	if err := json.Unmarshal([]byte(string(get.ResponseBody)), &jsonResponse); err != nil {
		return fmt.Errorf("%s doSpotKline Json Unmarshal Err: %v %v", e.GetName(), err, string(get.ResponseBody))
	} else if !jsonResponse.Success {
		return fmt.Errorf("%s doSpotKline Failed: %v", e.GetName(), string(get.ResponseBody))
	}
	if err := json.Unmarshal(jsonResponse.Result, &rawKline); err != nil {
		return fmt.Errorf("%s doSpotKline Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
	}

	operation.Kline = []*exchange.KlineDetail{}
	for _, k := range rawKline {

		detail := &exchange.KlineDetail{
			Exchange: e.GetName(),
			Pair:     operation.Pair.Name,
			OpenTime: float64(k.StartTime.Unix()) * 1000,
			Open:     k.Open,
			High:     k.High,
			Low:      k.Low,
			Close:    k.Close,
			Volume:   k.Volume,
		}

		operation.Kline = append(operation.Kline, detail)
	}

	return nil
}

func (e *Ftx) doGetFutureStats(operation *exchange.PublicOperation) error {
	// var str string
	var err error

	jsonResponse := &JsonResponse{}

	get := &utils.HttpGet{
		URI:       fmt.Sprintf("%s/api/futures/%s/stats", API_URL, operation.Pair.Symbol),
		Proxy:     operation.Proxy,
		DebugMode: operation.DebugMode,
	}
	if err = utils.HttpGetRequest(get); err != nil {
		operation.Error = err
		return operation.Error
	}

	if operation.DebugMode {
		operation.RequestURI = get.URI
		operation.CallResponce = string(get.ResponseBody)
	}

	// str = fmt.Sprintf("%s", operation.CallResponce)
	// log.Print(str)

	if err = json.Unmarshal([]byte(operation.CallResponce), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s doGetFutureStats Json Unmarshal Err: %v %v", e.GetName(), err, operation.CallResponce)
		return operation.Error
	} else if !jsonResponse.Success {
		operation.Error = fmt.Errorf("%s doGetFutureStats Failed: %v", e.GetName(), operation.CallResponce)
		return operation.Error
	}

	operation.FutureStats = &exchange.FutureStats{}

	if err := json.Unmarshal(jsonResponse.Result, operation.FutureStats); err != nil {
		operation.Error = fmt.Errorf("%s doGetFutureStats Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
		return operation.Error
	}

	// str = fmt.Sprintf("operation: %#v", operation)
	// log.Print(str)

	return err
}
