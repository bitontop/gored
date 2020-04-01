package huobi

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/bitontop/gored/exchange"
)

func (e *Huobi) DoAccoutOperation(operation *exchange.AccountOperation) error {
	switch operation.Type {

	// case exchange.Transfer:
	// 	return e.transfer(operation)
	// case exchange.BalanceList:
	// 	return e.getAllBalance(operation)
	// case exchange.Balance:
	// 	return e.getBalance(operation)

	case exchange.Withdraw:
		return e.doWithdraw(operation)
	case exchange.GetOpenOrder:
		return e.getOpenOrder(operation)
	case exchange.GetOrderHistory:
		return e.getOrderHistory(operation)
	case exchange.GetDepositAddress:
		return e.getDepositAddress(operation)
	case exchange.GetDepositHistory:
		return e.getDepositHistory(operation)
	case exchange.GetWithdrawalHistory:
		return e.getWithdrawalHistory(operation)
	}
	return fmt.Errorf("Operation type invalid: %v", operation.Type)
}

func (e *Huobi) doWithdraw(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	var withdrawID int64
	strRequest := "/v1/dw/withdraw/api/create"

	mapParams := make(map[string]string)
	mapParams["address"] = operation.WithdrawAddress
	mapParams["amount"] = operation.WithdrawAmount
	mapParams["currency"] = e.GetSymbolByCoin(operation.Coin)
	// mapParams["fee"] = strconv.FormatFloat(e.GetTxFee(operation.Coin), 'f', -1, 64) // Required parameter
	if operation.WithdrawTag != "" {
		mapParams["tag"] = operation.WithdrawTag
	}

	jsonWithdraw := e.ApiKeyRequest("POST", mapParams, strRequest)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonWithdraw
	}

	if err := json.Unmarshal([]byte(jsonWithdraw), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Json Unmarshal Err: %v, %s", e.GetName(), err, jsonWithdraw)
		return operation.Error
	} else if jsonResponse.Status != "ok" {
		operation.Error = fmt.Errorf("%s Withdraw Failed: %v", e.GetName(), jsonWithdraw)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &withdrawID); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	operation.WithdrawID = fmt.Sprintf("%v", withdrawID)

	return nil
}

func (e *Huobi) getOpenOrder(op *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	if e.Account_ID == "" {
		e.Account_ID = e.GetAccounts()
		if e.Account_ID == "" {
			return fmt.Errorf("%s Get AccountID Err", e.GetName())
		}
	}

	jsonResponse := &JsonResponse{}
	openOrders := []*OrderStatus{}
	strRequest := "/v1/order/openOrders"

	mapParams := make(map[string]string)
	mapParams["account-id"] = e.Account_ID
	mapParams["symbol"] = e.GetSymbolByPair(op.Pair)

	jsonOrders := e.ApiKeyRequest("GET", mapParams, strRequest)
	if op.DebugMode {
		op.RequestURI = strRequest
		op.CallResponce = jsonOrders
	}

	if err := json.Unmarshal([]byte(jsonOrders), &jsonResponse); err != nil {
		op.Error = fmt.Errorf("%s Get OpenOrders Json Unmarshal Err: %v, %s", e.GetName(), err, jsonOrders)
		return op.Error
	} else if jsonResponse.Status != "ok" {
		op.Error = fmt.Errorf("%s Withdraw Failed: %v", e.GetName(), jsonOrders)
		return op.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &openOrders); err != nil {
		op.Error = fmt.Errorf("%s Withdraw Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return op.Error
	}

	result := []*exchange.Order{}
	for _, data := range openOrders {
		order := &exchange.Order{
			Pair:    op.Pair,
			OrderID: fmt.Sprintf("%d", data.ID),
			Side:    data.Type,
		}

		order.Quantity, _ = strconv.ParseFloat(data.Amount, 64)
		order.Rate, _ = strconv.ParseFloat(data.Price, 64)

		if data.State == "canceled" {
			order.Status = exchange.Cancelled
		} else if data.State == "filled" {
			order.Status = exchange.Filled
		} else if data.State == "partial-filled" || data.State == "partial-canceled" {
			order.Status = exchange.Partial
		} else if data.State == "submitting" || data.State == "submitted" {
			order.Status = exchange.New
		} else {
			order.Status = exchange.Other
		}

		if data.FilledAmount != "" && data.FilledCashAmount != "" {
			dealQ, _ := strconv.ParseFloat(data.FilledAmount, 64)
			totalP, _ := strconv.ParseFloat(data.FilledCashAmount, 64)
			order.DealQuantity = dealQ
			if dealQ > 0 {
				order.DealRate = totalP / dealQ
			}
		} else {
			dealQ, _ := strconv.ParseFloat(data.FieldAmount, 64)
			totalP, _ := strconv.ParseFloat(data.FieldCashAmount, 64)
			order.DealQuantity = dealQ
			if dealQ > 0 {
				order.DealRate = totalP / dealQ
			}
		}
		result = append(result, order)
	}
	op.OpenOrders = result

	return nil
}

func (e *Huobi) getOrderHistory(op *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	orders := []*OrderStatus{}
	strRequest := "/v1/order/orders"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(op.Pair)

	jsonOrders := e.ApiKeyRequest("GET", mapParams, strRequest)
	if op.DebugMode {
		op.RequestURI = strRequest
		op.CallResponce = jsonOrders
	}

	if err := json.Unmarshal([]byte(jsonOrders), &jsonResponse); err != nil {
		op.Error = fmt.Errorf("%s Get Order History Json Unmarshal Err: %v, %s", e.GetName(), err, jsonOrders)
		return op.Error
	} else if jsonResponse.Status != "ok" {
		op.Error = fmt.Errorf("%s Get Order History Failed: %v", e.GetName(), jsonOrders)
		return op.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &orders); err != nil {
		op.Error = fmt.Errorf("%s Get Order History Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return op.Error
	}

	result := []*exchange.Order{}
	for _, data := range orders {
		order := &exchange.Order{
			Pair:    op.Pair,
			OrderID: fmt.Sprintf("%d", data.ID),
			Side:    data.Type,
		}

		order.Quantity, _ = strconv.ParseFloat(data.Amount, 64)
		order.Rate, _ = strconv.ParseFloat(data.Price, 64)

		if data.State == "canceled" {
			order.Status = exchange.Cancelled
		} else if data.State == "filled" {
			order.Status = exchange.Filled
		} else if data.State == "partial-filled" || data.State == "partial-canceled" {
			order.Status = exchange.Partial
		} else if data.State == "submitting" || data.State == "submitted" {
			order.Status = exchange.New
		} else {
			order.Status = exchange.Other
		}

		if data.FilledAmount != "" && data.FilledCashAmount != "" {
			dealQ, _ := strconv.ParseFloat(data.FilledAmount, 64)
			totalP, _ := strconv.ParseFloat(data.FilledCashAmount, 64)
			order.DealQuantity = dealQ
			if dealQ > 0 {
				order.DealRate = totalP / dealQ
			}
		} else {
			dealQ, _ := strconv.ParseFloat(data.FieldAmount, 64)
			totalP, _ := strconv.ParseFloat(data.FieldCashAmount, 64)
			order.DealQuantity = dealQ
			if dealQ > 0 {
				order.DealRate = totalP / dealQ
			}
		}
		result = append(result, order)
	}
	op.OrderHistory = result

	return nil
}

func (e *Huobi) getDepositAddress(op *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	address := []*DepositAddress{}
	strRequest := "/v2/account/deposit/address"

	mapParams := make(map[string]string)
	mapParams["currency"] = e.GetSymbolByCoin(op.Coin)

	jsonDepositAddress := e.ApiKeyRequest("GET", mapParams, strRequest)
	if op.DebugMode {
		op.RequestURI = strRequest
		op.CallResponce = jsonDepositAddress
	}

	if err := json.Unmarshal([]byte(jsonDepositAddress), &jsonResponse); err != nil {
		op.Error = fmt.Errorf("%s Get Deposit Address Json Unmarshal Err: %v, %s", e.GetName(), err, jsonDepositAddress)
		return op.Error
	} else if jsonResponse.Status != "ok" {
		op.Error = fmt.Errorf("%s Get Deposit Address Failed: %v", e.GetName(), jsonDepositAddress)
		return op.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &address); err != nil {
		op.Error = fmt.Errorf("%s Get Deposit Address Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return op.Error
	}

	for _, data := range address {
		addr := &exchange.DepositAddr{
			Coin:    op.Coin,
			Address: data.Address,
			Tag:     data.AddressTag,
			Chain:   exchange.MAINNET,
		}
		op.DepositAddresses[addr.Chain] = addr
	}

	return nil
}

func (e *Huobi) getDepositHistory(op *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	histories := []*DWHistory{}
	strRequest := "/v1/query/deposit-withdraw"

	mapParams := make(map[string]string)
	mapParams["type"] = "deposit"

	jsonDWHistory := e.ApiKeyRequest("GET", mapParams, strRequest)
	if op.DebugMode {
		op.RequestURI = strRequest
		op.CallResponce = jsonDWHistory
	}

	if err := json.Unmarshal([]byte(jsonDWHistory), &jsonResponse); err != nil {
		op.Error = fmt.Errorf("%s Deposit History Json Unmarshal Err: %v, %s", e.GetName(), err, jsonDWHistory)
		return op.Error
	} else if jsonResponse.Status != "ok" {
		op.Error = fmt.Errorf("%s Deposit History Failed: %v", e.GetName(), jsonDWHistory)
		return op.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &histories); err != nil {
		op.Error = fmt.Errorf("%s Deposit History Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return op.Error
	}

	result := []*exchange.WDHistory{}
	for _, data := range histories {
		history := &exchange.WDHistory{
			ID:       fmt.Sprintf("%d", data.ID),
			Coin:     e.GetCoinBySymbol(data.Currency),
			Quantity: data.Amount,
			Tag:      data.AddressTag,
			Address:  data.Address,
			TxHash:   data.TxHash,
		}
		result = append(result, history)
	}
	op.WithdrawalHistory = result

	return nil
}

func (e *Huobi) getWithdrawalHistory(op *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	histories := []*DWHistory{}
	strRequest := "/v1/query/deposit-withdraw"

	mapParams := make(map[string]string)
	mapParams["type"] = "withdraw"

	jsonDWHistory := e.ApiKeyRequest("GET", mapParams, strRequest)
	if op.DebugMode {
		op.RequestURI = strRequest
		op.CallResponce = jsonDWHistory
	}

	if err := json.Unmarshal([]byte(jsonDWHistory), &jsonResponse); err != nil {
		op.Error = fmt.Errorf("%s Withdraw History Json Unmarshal Err: %v, %s", e.GetName(), err, jsonDWHistory)
		return op.Error
	} else if jsonResponse.Status != "ok" {
		op.Error = fmt.Errorf("%s Withdraw History Failed: %v", e.GetName(), jsonDWHistory)
		return op.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &histories); err != nil {
		op.Error = fmt.Errorf("%s Withdraw History Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return op.Error
	}

	result := []*exchange.WDHistory{}
	for _, data := range histories {
		history := &exchange.WDHistory{
			ID:       fmt.Sprintf("%d", data.ID),
			Coin:     e.GetCoinBySymbol(data.Currency),
			Quantity: data.Amount,
			Tag:      data.AddressTag,
			Address:  data.Address,
			TxHash:   data.TxHash,
		}
		result = append(result, history)
	}
	op.WithdrawalHistory = result

	return nil
}
