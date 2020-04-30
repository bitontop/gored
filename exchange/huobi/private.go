package huobi

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/bitontop/gored/exchange"
)

func (e *Huobi) DoAccountOperation(operation *exchange.AccountOperation) error {
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
		if operation.Wallet == exchange.SpotWallet {
			return e.getOpenOrder(operation)
		}
	case exchange.GetOrderHistory:
		if operation.Wallet == exchange.SpotWallet {
			return e.getOrderHistory(operation)
		}
	case exchange.GetDepositAddress:
		if operation.Wallet == exchange.SpotWallet {
			return e.getDepositAddress(operation)
		}
	case exchange.GetDepositHistory:
		if operation.Wallet == exchange.SpotWallet {
			return e.getDepositHistory(operation)
		}
	case exchange.GetWithdrawalHistory:
		if operation.Wallet == exchange.SpotWallet {
			return e.getWithdrawalHistory(operation)
		}
	case exchange.SubBalanceList:
		if operation.Wallet == exchange.SpotWallet {
			return e.doSubBalance(operation)
		}
	case exchange.GetSubAccountList: // all type accounts
		// if operation.Wallet == exchange.SpotWallet {
		return e.doSubAccountList(operation)
		// }
	case exchange.SubAllBalanceList:
		if operation.Wallet == exchange.SpotWallet {
			return e.doSubAllBalance(operation)
		}
	}
	return fmt.Errorf("%s Operation type invalid: %s %v", operation.Ex, operation.Wallet, operation.Type)
}

// 查询当前用户的"所有"账户 ID 及其相关信息
// sub account if 'Subtype' != "" （仅对逐仓杠杆账户有效）
func (e *Huobi) doSubAccountList(operation *exchange.AccountOperation) error { //TODO, test with sub account
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	accountList := SubAccountList{}
	strRequest := "/v1/account/accounts"

	mapParams := make(map[string]string)

	jsonSubAccountReturn := e.ApiKeyRequest("GET", mapParams, strRequest)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonSubAccountReturn
	}

	if err := json.Unmarshal([]byte(jsonSubAccountReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s doSubAllBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonSubAccountReturn)
		return operation.Error
	} else if jsonResponse.Status != "ok" {
		operation.Error = fmt.Errorf("%s doSubAllBalance failed: %v", e.GetName(), jsonSubAccountReturn)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountList); err != nil {
		operation.Error = fmt.Errorf("%s doSubAllBalance Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	operation.SubAccountList = []*exchange.SubAccountInfo{}
	for _, account := range accountList {
		var accountType exchange.WalletType
		if account.Type == "spot" {
			accountType = exchange.SpotWallet
		} else if account.Type == "margin" || account.Type == "super-margin" {
			accountType = exchange.MarginWallet
		} else if account.Type == "otc" {
			accountType = exchange.FiatOTCWallet
		}

		a := &exchange.SubAccountInfo{
			ID:          fmt.Sprintf("%v", account.ID),
			Status:      account.State,
			AccountType: accountType,
			Activated:   true,
			// TimeStamp:   account.CreateTime,
		}
		operation.SubAccountList = append(operation.SubAccountList, a)
	}

	return nil
}

func (e *Huobi) doSubBalance(operation *exchange.AccountOperation) error { //TODO, test with sub account
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	accountBalance := SubAccountBalances{}
	strRequest := fmt.Sprintf("/v1/account/accounts/%v", operation.SubAccountID)

	mapParams := make(map[string]string)
	mapParams["sub-uid"] = operation.SubAccountID

	jsonBalanceReturn := e.ApiKeyRequest("GET", mapParams, strRequest)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonBalanceReturn
	}

	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s doSubBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonBalanceReturn)
		return operation.Error
	} else if jsonResponse.Status != "ok" {
		operation.Error = fmt.Errorf("%s doSubBalance failed: %v", e.GetName(), jsonBalanceReturn)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		operation.Error = fmt.Errorf("%s doSubBalance Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	operation.BalanceList = []exchange.AssetBalance{}
	if len(accountBalance) == 0 {
		log.Printf("%s doSubBalance got empty list: %v", e.GetName(), jsonBalanceReturn)
		return nil
	}
	for _, balance := range accountBalance[0].List {
		totalAmount, err := strconv.ParseFloat(balance.Balance, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doSubBalance parse err: %+v %v", e.GetName(), balance, err)
			return operation.Error
		}

		c := e.GetCoinBySymbol(balance.Currency)
		if c == nil {
			continue
		}
		b := exchange.AssetBalance{
			Coin:             c,
			BalanceAvailable: totalAmount,
			BalanceFrozen:    0,
		}
		operation.BalanceList = append(operation.BalanceList, b)
	}

	return nil
}

func (e *Huobi) doSubAllBalance(operation *exchange.AccountOperation) error { //TODO, test with sub account
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	accountBalance := SubAllAccountBalances{}
	strRequest := "/v1/subuser/aggregate-balance"

	mapParams := make(map[string]string)

	jsonSubAllBalanceReturn := e.ApiKeyRequest("GET", mapParams, strRequest)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonSubAllBalanceReturn
	}

	if err := json.Unmarshal([]byte(jsonSubAllBalanceReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s doSubAllBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonSubAllBalanceReturn)
		return operation.Error
	} else if jsonResponse.Status != "ok" {
		operation.Error = fmt.Errorf("%s doSubAllBalance failed: %v", e.GetName(), jsonSubAllBalanceReturn)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		operation.Error = fmt.Errorf("%s doSubAllBalance Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	operation.BalanceList = []exchange.AssetBalance{}
	if len(accountBalance) == 0 {
		log.Printf("%s doSubAllBalance got empty list: %v", e.GetName(), jsonSubAllBalanceReturn)
		return nil
	}
	for _, balance := range accountBalance {
		totalAmount, err := strconv.ParseFloat(balance.Balance, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doSubBalance parse err: %+v %v", e.GetName(), balance, err)
			return operation.Error
		}

		c := e.GetCoinBySymbol(balance.Currency)
		if c == nil {
			continue
		}
		b := exchange.AssetBalance{
			Coin:             c,
			BalanceAvailable: totalAmount,
			BalanceFrozen:    0,
		}
		operation.BalanceList = append(operation.BalanceList, b)
	}

	return nil
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
			Pair:      op.Pair,
			OrderID:   fmt.Sprintf("%d", data.ID),
			Side:      data.Type,
			Timestamp: data.CreatedAt,
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
			Pair:      op.Pair,
			OrderID:   fmt.Sprintf("%d", data.ID),
			Side:      data.Type,
			Timestamp: data.FinishedAt,
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

	op.DepositAddresses = make(map[exchange.ChainType]*exchange.DepositAddr)
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
	} else if jsonResponse.Code != 200 {
		op.Error = fmt.Errorf("%s Get Deposit Address Failed: %v", e.GetName(), jsonDepositAddress)
		return op.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &address); err != nil {
		op.Error = fmt.Errorf("%s Get Deposit Address Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return op.Error
	}

	for _, data := range address {
		if data.Currency == data.Chain {
			addr := &exchange.DepositAddr{
				Coin:    op.Coin,
				Address: data.Address,
				Tag:     data.AddressTag,
				Chain:   exchange.MAINNET,
			}
			op.DepositAddresses[addr.Chain] = addr
		}
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
			ID:        fmt.Sprintf("%d", data.ID),
			Coin:      e.GetCoinBySymbol(data.Currency),
			Quantity:  data.Amount,
			Tag:       data.AddressTag,
			Address:   data.Address,
			TxHash:    data.TxHash,
			ChainType: exchange.MAINNET,
			Status:    data.State,
			TimeStamp: data.UpdatedAt,
		}
		result = append(result, history)
	}
	op.DepositHistory = result

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
			ID:        fmt.Sprintf("%d", data.ID),
			Coin:      e.GetCoinBySymbol(data.Currency),
			Quantity:  data.Amount,
			Tag:       data.AddressTag,
			Address:   data.Address,
			TxHash:    data.TxHash,
			ChainType: exchange.MAINNET,
			Status:    data.State,
			TimeStamp: data.UpdatedAt,
		}
		result = append(result, history)
	}
	op.WithdrawalHistory = result

	return nil
}
