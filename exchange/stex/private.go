package stex

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/bitontop/gored/exchange"
)

func (e *Stex) DoAccountOperation(operation *exchange.AccountOperation) error {
	switch operation.Type {

	// case exchange.Transfer:
	// 	return e.transfer(operation)
	// case exchange.BalanceList:
	// 	return e.getAllBalance(operation)
	// case exchange.Balance:
	// 	return e.getBalance(operation)
	case exchange.GetOpenOrder:
		if operation.Wallet == exchange.SpotWallet {
			return e.doGetOpenOrder(operation)
		}

	case exchange.Withdraw:
		return e.doWithdraw(operation)

	}
	return fmt.Errorf("%s Operation type invalid: %s %v", operation.Ex, operation.Wallet, operation.Type)
}

func (e *Stex) doGetOpenOrder(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	jsonResponse := &JsonResponseV3{}
	openOrders := OpenOrders{}

	strRequestUrl := "/trading/orders"
	if operation.Pair != nil {
		strRequestUrl += fmt.Sprintf("/%s", e.GetIDByPair(operation.Pair))
	}

	jsonOpenOrders := e.ApiKeyGet(nil, strRequestUrl)
	if operation.DebugMode {
		operation.RequestURI = strRequestUrl
		operation.CallResponce = jsonOpenOrders
	}

	if err := json.Unmarshal([]byte(jsonOpenOrders), &jsonResponse); err != nil {
		return fmt.Errorf("%s doGetOpenOrder Json Unmarshal Err: %v %s", e.GetName(), err, jsonOpenOrders)
	} else if !jsonResponse.Success {
		return fmt.Errorf("%s doGetOpenOrder Failed: %v %s", e.GetName(), jsonResponse.Message, jsonOpenOrders)
	}

	if err := json.Unmarshal(jsonResponse.Data, &openOrders); err != nil {
		return fmt.Errorf("%s doGetOpenOrder Unmarshal Error: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	// store info into orders
	operation.OpenOrders = []*exchange.Order{}
	for _, o := range openOrders {
		rate, err := strconv.ParseFloat(o.Price, 64)
		if err != nil {
			log.Printf("%s parse rate Err: %v %s", e.GetName(), err, o.Price)
		}
		quantity, err := strconv.ParseFloat(o.InitialAmount, 64)
		if err != nil {
			log.Printf("%s parse quantity Err: %v %s", e.GetName(), err, o.InitialAmount)
		}
		dealQuantity, err := strconv.ParseFloat(o.ProcessedAmount, 64)
		if err != nil {
			log.Printf("%s parse dealQuantity Err: %v %s", e.GetName(), err, o.ProcessedAmount)
		}
		ts, err := strconv.ParseInt(o.Timestamp, 10, 64)
		if err != nil {
			log.Printf("%s parse timestamp Err: %v %s", e.GetName(), err, o.Timestamp)
		}

		order := &exchange.Order{
			Pair:         e.GetPairBySymbol(o.CurrencyPairName),
			OrderID:      fmt.Sprintf("%v", o.ID),
			Rate:         rate,
			Quantity:     quantity,
			DealRate:     rate,
			DealQuantity: dealQuantity,
			Timestamp:    ts,
			// JsonResponse: jsonGetOpenOrder,
		}

		switch o.Type {
		case "BUY":
			order.Direction = exchange.Buy
		case "SELL":
			order.Direction = exchange.Sell
		}

		if o.Status == "PROCESSING" {
			order.Status = exchange.New
		} else if o.Status == "PENDING" {
			order.Status = exchange.Cancelled
		} else if o.Status == "FINISHED" {
			order.Status = exchange.Filled
		} else if o.Status == "PARTIAL" {
			order.Status = exchange.Partial
		} else if o.Status == "CANCELLED" {
			order.Status = exchange.Cancelled
		} else {
			log.Printf("%v OpenOrder %v unknown type: %v", e.GetName(), order.OrderID, o.Status)
			order.Status = exchange.Other
		}

		operation.OpenOrders = append(operation.OpenOrders, order)
	}

	return nil
}
