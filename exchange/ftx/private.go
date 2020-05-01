package ftx

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/bitontop/gored/exchange"
)

func (e *Ftx) DoAccountOperation(operation *exchange.AccountOperation) error {
	switch operation.Type {
	case exchange.BalanceList:
		if operation.Wallet == exchange.SpotWallet {
			return e.getAllBalance(operation)
		}
	case exchange.Balance:
		if operation.Wallet == exchange.SpotWallet {
			return e.getBalance(operation)
		}

	// Private operation
	case exchange.GetOpenOrder:
		if operation.Wallet == exchange.SpotWallet {
			return e.doGetOpenOrder(operation)
		}
	case exchange.GetOrderHistory:
		if operation.Wallet == exchange.SpotWallet {
			return e.doGetOrderHistory(operation)
		}
	case exchange.GetWithdrawalHistory:
		if operation.Wallet == exchange.SpotWallet {
			return e.doGetWithdrawalHistory(operation)
		}
	case exchange.GetDepositHistory:
		if operation.Wallet == exchange.SpotWallet {
			return e.doGetDepositHistory(operation)
		}
	case exchange.GetDepositAddress:
		if operation.Wallet == exchange.SpotWallet {
			return e.doGetDepositAddress(operation)
		}
	}
	return fmt.Errorf("Operation type invalid: %v", operation.Type)
}

func (e *Ftx) doGetOpenOrder(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key or passphrase are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	openOrders := OpenOrders{}
	symbol := e.GetSymbolByPair(operation.Pair)
	strRequest := "/orders" // /orders?market={market}

	mapParams := make(map[string]string)
	mapParams["market"] = symbol

	jsonGetOpenOrder := e.ApiKeyRequest("GET", strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonGetOpenOrder
	}

	if err := json.Unmarshal([]byte(jsonGetOpenOrder), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s doGetOpenOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonGetOpenOrder)
		return operation.Error
	} else if !jsonResponse.Success {
		operation.Error = fmt.Errorf("%s doGetOpenOrder Failed: %v", e.GetName(), jsonGetOpenOrder)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Result, &openOrders); err != nil {
		operation.Error = fmt.Errorf("%s doGetOpenOrder Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
		return operation.Error
	}

	// store info into orders
	operation.OpenOrders = []*exchange.Order{}
	for _, o := range openOrders {

		order := &exchange.Order{
			Pair:         operation.Pair,
			OrderID:      fmt.Sprintf("%v", o.ID),
			Rate:         o.Price,
			Quantity:     o.Size,
			DealRate:     o.AvgFillPrice,
			DealQuantity: o.Size - o.RemainingSize,
			Timestamp:    o.CreatedAt.UnixNano(),
			// JsonResponse: jsonGetOpenOrder,
		}

		switch o.Side {
		case "buy":
			order.Direction = exchange.Buy
		case "sell":
			order.Direction = exchange.Sell
		}

		if o.Status == "new" {
			order.Status = exchange.New
		} else if o.Status == "closed" && o.RemainingSize == 0 {
			order.Status = exchange.Filled
		} else if o.Status == "closed" && o.RemainingSize != 0 {
			order.Status = exchange.Cancelled
		} else if o.Status == "open" && o.Size == o.RemainingSize {
			order.Status = exchange.New
		} else if o.Status == "open" && o.Size > o.RemainingSize {
			order.Status = exchange.Partial
		} else {
			order.Status = exchange.Other
		}

		operation.OpenOrders = append(operation.OpenOrders, order)
	}

	return nil
}

func (e *Ftx) doGetOrderHistory(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key or passphrase are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	closeOrders := CloseOrders{}
	symbol := e.GetSymbolByPair(operation.Pair)
	strRequest := "/orders/history" // /orders/history?market={market}

	mapParams := make(map[string]string)
	mapParams["market"] = symbol

	jsonGetOrderHistory := e.ApiKeyRequest("GET", strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonGetOrderHistory
	}

	if err := json.Unmarshal([]byte(jsonGetOrderHistory), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s doGetOrderHistory Json Unmarshal Err: %v %v", e.GetName(), err, jsonGetOrderHistory)
		return operation.Error
	} else if !jsonResponse.Success {
		operation.Error = fmt.Errorf("%s doGetOrderHistory Failed: %v", e.GetName(), jsonGetOrderHistory)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Result, &closeOrders); err != nil {
		operation.Error = fmt.Errorf("%s doGetOrderHistory Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
		return operation.Error
	}

	// store info into orders
	operation.OrderHistory = []*exchange.Order{}
	for _, o := range closeOrders {

		order := &exchange.Order{
			Pair:         operation.Pair,
			OrderID:      fmt.Sprintf("%v", o.ID),
			Rate:         o.Price,
			Quantity:     o.Size,
			DealRate:     o.AvgFillPrice,
			DealQuantity: o.Size - o.RemainingSize,
			Timestamp:    o.CreatedAt.UnixNano(),
			// JsonResponse: jsonGetOrderHistory,
		}

		switch o.Side {
		case "buy":
			order.Direction = exchange.Buy
		case "sell":
			order.Direction = exchange.Sell
		}

		if o.Status == "new" {
			order.Status = exchange.New
		} else if o.Status == "closed" && o.RemainingSize == 0 {
			order.Status = exchange.Filled
		} else if o.Status == "closed" && o.RemainingSize != 0 {
			order.Status = exchange.Cancelled
		} else if o.Status == "open" && o.Size == o.RemainingSize {
			order.Status = exchange.New
		} else if o.Status == "open" && o.Size > o.RemainingSize {
			order.Status = exchange.Partial
		} else {
			order.Status = exchange.Other
		}

		operation.OrderHistory = append(operation.OrderHistory, order)
	}

	return nil
}

// FTX doesn't provide chainType information, use default MAINNET
func (e *Ftx) doGetWithdrawalHistory(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key or passphrase are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	withdrawHistory := WithdrawHistory{}
	strRequest := "/wallet/withdrawals"

	mapParams := make(map[string]string)

	jsonGetWithdrawalHistory := e.ApiKeyRequest("GET", strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonGetWithdrawalHistory
	}

	if err := json.Unmarshal([]byte(jsonGetWithdrawalHistory), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s doGetWithdrawalHistory Json Unmarshal Err: %v %v", e.GetName(), err, jsonGetWithdrawalHistory)
		return operation.Error
	} else if !jsonResponse.Success {
		operation.Error = fmt.Errorf("%s doGetWithdrawalHistory Failed: %v", e.GetName(), jsonGetWithdrawalHistory)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Result, &withdrawHistory); err != nil {
		operation.Error = fmt.Errorf("%s doGetWithdrawalHistory Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
		return operation.Error
	}

	// store info into orders
	operation.WithdrawalHistory = []*exchange.WDHistory{}
	for _, withdrawRecord := range withdrawHistory {
		c := e.GetCoinBySymbol(withdrawRecord.Coin)
		quantity, err := strconv.ParseFloat(withdrawRecord.Size, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doGetWithdrawalHistory parse quantity Err: %v, %v", e.GetName(), err, withdrawRecord.Size)
			return operation.Error
		}

		var chainType exchange.ChainType
		chainType = exchange.MAINNET
		statusMsg := withdrawRecord.Status

		record := &exchange.WDHistory{
			ID:        fmt.Sprintf("%v", withdrawRecord.ID),
			Coin:      c,
			Quantity:  quantity,
			Tag:       withdrawRecord.Tag,
			Address:   withdrawRecord.Address,
			TxHash:    withdrawRecord.Txid,
			ChainType: chainType,
			Status:    statusMsg,
			TimeStamp: withdrawRecord.Time.UnixNano(),
		}

		operation.WithdrawalHistory = append(operation.WithdrawalHistory, record)
	}

	return nil
}

// FTX doesn't provide chainType information, use default MAINNET
func (e *Ftx) doGetDepositHistory(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key or passphrase are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	depositHistory := DepositHistory{}
	strRequest := "/wallet/deposits"

	mapParams := make(map[string]string)

	jsonGetDepositHistory := e.ApiKeyRequest("GET", strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonGetDepositHistory
	}

	if err := json.Unmarshal([]byte(jsonGetDepositHistory), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s doGetDepositHistory Json Unmarshal Err: %v %v", e.GetName(), err, jsonGetDepositHistory)
		return operation.Error
	} else if !jsonResponse.Success {
		operation.Error = fmt.Errorf("%s doGetDepositHistory Failed: %v", e.GetName(), jsonGetDepositHistory)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Result, &depositHistory); err != nil {
		operation.Error = fmt.Errorf("%s doGetDepositHistory Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
		return operation.Error
	}

	// store info into orders
	operation.DepositHistory = []*exchange.WDHistory{}
	for _, depositRecord := range depositHistory {
		c := e.GetCoinBySymbol(depositRecord.Coin)
		quantity, err := strconv.ParseFloat(depositRecord.Size, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doGetDepositHistory parse quantity Err: %v, %v", e.GetName(), err, depositRecord.Size)
			return operation.Error
		}

		var chainType exchange.ChainType
		chainType = exchange.MAINNET
		statusMsg := depositRecord.Status

		record := &exchange.WDHistory{
			ID:        fmt.Sprintf("%v", depositRecord.ID),
			Coin:      c,
			Quantity:  quantity,
			TxHash:    depositRecord.Txid,
			ChainType: chainType,
			Status:    statusMsg,
			TimeStamp: depositRecord.Time.UnixNano(),
		}

		operation.DepositHistory = append(operation.DepositHistory, record)
	}

	return nil
}

// FTX doesn't provide chainType information, use default MAINNET
func (e *Ftx) doGetDepositAddress(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key or passphrase are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	depositAddress := DepositAddress{}
	strRequest := fmt.Sprintf("/wallet/deposit_address/%v", e.GetSymbolByCoin(operation.Coin))

	mapParams := make(map[string]string)

	jsonGetDepositAddress := e.ApiKeyRequest("GET", strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonGetDepositAddress
	}

	if err := json.Unmarshal([]byte(jsonGetDepositAddress), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s doGetDepositAddress Json Unmarshal Err: %v %v", e.GetName(), err, jsonGetDepositAddress)
		return operation.Error
	} else if !jsonResponse.Success {
		operation.Error = fmt.Errorf("%s doGetDepositAddress Failed: %v", e.GetName(), jsonGetDepositAddress)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Result, &depositAddress); err != nil {
		operation.Error = fmt.Errorf("%s doGetDepositAddress Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
		return operation.Error
	}

	operation.DepositAddresses = make(map[exchange.ChainType]*exchange.DepositAddr)

	var chain exchange.ChainType
	chain = exchange.MAINNET

	depoAddr := &exchange.DepositAddr{
		Coin:    operation.Coin,
		Address: depositAddress.Address,
		Tag:     depositAddress.Tag,
		Chain:   chain,
	}

	operation.DepositAddresses[chain] = depoAddr

	return nil
}

func (e *Ftx) getAllBalance(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalances{}
	strRequest := "/wallet/balances"

	jsonBalanceReturn := e.ApiKeyRequest("GET", strRequest, make(map[string]string))
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonBalanceReturn
	}

	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s GetAllBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return operation.Error
	} else if !jsonResponse.Success {
		operation.Error = fmt.Errorf("%s GetAllBalances Failed: %v", e.GetName(), jsonBalanceReturn)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Result, &accountBalance); err != nil {
		operation.Error = fmt.Errorf("%s GetAllBalances Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
		return operation.Error
	}

	operation.BalanceList = []exchange.AssetBalance{}
	for _, account := range accountBalance {
		if account.Total == 0 {
			continue
		}

		balance := exchange.AssetBalance{
			Coin:             e.GetCoinBySymbol(account.Coin),
			BalanceAvailable: account.Free,
			BalanceFrozen:    account.Total - account.Free,
		}
		operation.BalanceList = append(operation.BalanceList, balance)

	}

	return nil
}

func (e *Ftx) getBalance(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	accountBalance := AccountBalances{}
	strRequest := "/wallet/balances"

	jsonBalanceReturn := e.ApiKeyRequest("GET", strRequest, make(map[string]string))
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonBalanceReturn
	}

	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s GetBalances Json Unmarshal Err: %v %v", e.GetName(), err, jsonBalanceReturn)
		return operation.Error
	} else if !jsonResponse.Success {
		operation.Error = fmt.Errorf("%s GetBalances Failed: %v", e.GetName(), jsonBalanceReturn)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Result, &accountBalance); err != nil {
		operation.Error = fmt.Errorf("%s GetBalances Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Result)
		return operation.Error
	}

	operation.BalanceFrozen = 0
	operation.BalanceAvailable = 0
	for _, account := range accountBalance {
		if account.Coin == e.GetSymbolByCoin(operation.Coin) {
			operation.BalanceFrozen = account.Total - account.Free
			operation.BalanceAvailable = account.Free
			return nil
		}
	}

	return nil
}
