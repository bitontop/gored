package binance

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bitontop/gored/exchange"
)

func (e *Binance) DoAccountOperation(operation *exchange.AccountOperation) error {
	switch operation.Type {
	case exchange.Withdraw:
		return e.doWithdraw(operation)
	// case exchange.Transfer:
	// 	return e.transfer(operation)
	// case exchange.BalanceList:
	// 	return e.getAllBalance(operation)
	// case exchange.Balance:
	// 	return e.getBalance(operation)

	// Contract operation
	case exchange.PlaceOrder:
		if operation.Wallet == exchange.ContractWallet {
			return e.doContractPlaceOrder(operation)
		}
	// case exchange.GetOrderStatus: // operation model changed
	// 	if operation.Wallet == exchange.ContractWallet {
	// 		return e.doContractOrderStatus(operation)
	// 	}
	case exchange.CancelOrder:
		if operation.Wallet == exchange.ContractWallet {
			return e.doContractCancelOrder(operation)
		}
	case exchange.BalanceList:
		if operation.Wallet == exchange.ContractWallet {
			return e.doContractAllBalance(operation)
		}
	// case exchange.Balance:
	// 	if operation.Wallet == exchange.ContractWallet {
	// 		return e.doContractBalance(operation)
	// 	}

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

func (e *Binance) doGetOpenOrder(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key or passphrase are nil.", e.GetName())
	}

	openOrders := OpenOrders{}
	strRequest := "/api/v3/openOrders"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(operation.Pair)

	jsonGetOpenOrder := e.ApiKeyGet(mapParams, strRequest)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		// operation.MapParams = fmt.Sprintf("%+v", mapParams)
		operation.CallResponce = jsonGetOpenOrder
	}

	if err := json.Unmarshal([]byte(jsonGetOpenOrder), &openOrders); err != nil {
		operation.Error = fmt.Errorf("%s doGetOpenOrder Json Unmarshal Err: %v, %s", e.GetName(), err, jsonGetOpenOrder)
		return operation.Error
	}

	// store info into orders
	operation.OpenOrders = []*exchange.Order{}
	for _, o := range openOrders {
		rate, err := strconv.ParseFloat(o.Price, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doGetOpenOrder parse rate Err: %v, %v", e.GetName(), err, o.Price)
			return operation.Error
		}
		quantity, err := strconv.ParseFloat(o.OrigQty, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doGetOpenOrder parse quantity Err: %v, %v", e.GetName(), err, o.OrigQty)
			return operation.Error
		}
		dealQuantity, err := strconv.ParseFloat(o.ExecutedQty, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doGetOpenOrder parse dealQuantity Err: %v, %v", e.GetName(), err, o.ExecutedQty)
			return operation.Error
		}

		order := &exchange.Order{
			Pair:         operation.Pair,
			OrderID:      fmt.Sprintf("%v", o.OrderID),
			Rate:         rate,
			Quantity:     quantity,
			Side:         o.Side,
			DealRate:     rate,
			DealQuantity: dealQuantity,
			Timestamp:    o.UpdateTime,
			// JsonResponse: jsonGetOpenOrder,
		}

		if o.Status == "CANCELED" {
			order.Status = exchange.Cancelled
		} else if o.Status == "FILLED" {
			order.Status = exchange.Filled
		} else if o.Status == "PARTIALLY_FILLED" {
			order.Status = exchange.Partial
		} else if o.Status == "REJECTED" {
			order.Status = exchange.Rejected
		} else if o.Status == "Expired" {
			order.Status = exchange.Expired
		} else if o.Status == "NEW" {
			order.Status = exchange.New
		} else {
			order.Status = exchange.Other
		}

		operation.OpenOrders = append(operation.OpenOrders, order)
	}

	return nil
}

func (e *Binance) doGetOrderHistory(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key or passphrase are nil.", e.GetName())
	}

	closeOrders := CloseOrders{}
	strRequest := "/api/v3/myTrades"

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(operation.Pair)

	jsonGetOrderHistory := e.ApiKeyGet(mapParams, strRequest)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		// operation.MapParams = fmt.Sprintf("%+v", mapParams)
		operation.CallResponce = jsonGetOrderHistory
	}

	if err := json.Unmarshal([]byte(jsonGetOrderHistory), &closeOrders); err != nil {
		operation.Error = fmt.Errorf("%s doGetOrderHistory Json Unmarshal Err: %v, %s", e.GetName(), err, jsonGetOrderHistory)
		return operation.Error
	}

	// store info into orders
	operation.OrderHistory = []*exchange.Order{}
	for _, o := range closeOrders {
		rate, err := strconv.ParseFloat(o.Price, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doGetOrderHistory parse rate Err: %v, %v", e.GetName(), err, o.Price)
			return operation.Error
		}
		quantity, err := strconv.ParseFloat(o.Qty, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doGetOrderHistory parse quantity Err: %v, %v", e.GetName(), err, o.Qty)
			return operation.Error
		}

		side := ""
		if o.IsBuyer {
			side = "Buy"
		} else {
			side = "Sell"
		}

		order := &exchange.Order{
			Pair:         operation.Pair,
			OrderID:      fmt.Sprintf("%v", o.OrderID),
			Rate:         rate,
			Quantity:     quantity,
			Side:         side,
			DealRate:     rate,
			DealQuantity: quantity,
			Timestamp:    o.Time,
			// JsonResponse: jsonGetOrderHistory,
		}

		order.Status = exchange.Filled

		operation.OrderHistory = append(operation.OrderHistory, order)
	}

	return nil
}

func (e *Binance) doGetWithdrawalHistory(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key or passphrase are nil.", e.GetName())
	}

	withdrawHistory := WithdrawHistory{}
	strRequest := "/sapi/v1/capital/withdraw/history"

	mapParams := make(map[string]string)

	jsonGetWithdrawalHistory := e.ApiKeyGet(mapParams, strRequest)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		// operation.MapParams = fmt.Sprintf("%+v", mapParams)
		operation.CallResponce = jsonGetWithdrawalHistory
	}

	if err := json.Unmarshal([]byte(jsonGetWithdrawalHistory), &withdrawHistory); err != nil {
		operation.Error = fmt.Errorf("%s doGetWithdrawalHistory Json Unmarshal Err: %v, %s", e.GetName(), err, jsonGetWithdrawalHistory)
		return operation.Error
	}

	// store info into orders
	operation.WithdrawalHistory = []*exchange.WDHistory{}
	for _, withdrawRecord := range withdrawHistory {
		c := e.GetCoinBySymbol(withdrawRecord.Coin)
		quantity, err := strconv.ParseFloat(withdrawRecord.Amount, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doGetWithdrawalHistory parse quantity Err: %v, %v", e.GetName(), err, withdrawRecord.Amount)
			return operation.Error
		}
		var chainType exchange.ChainType
		if withdrawRecord.Network == "BTC" {
			chainType = exchange.MAINNET
		} else if withdrawRecord.Network == "ETH" {
			chainType = exchange.ERC20
		} else {
			chainType = exchange.OTHER
		}

		statusMsg := ""
		if withdrawRecord.Status == 0 {
			statusMsg = "Confirm email sent"
		} else if withdrawRecord.Status == 1 {
			statusMsg = "Canceled by user"
		} else if withdrawRecord.Status == 2 {
			statusMsg = "Waiting for Confirmation"
		} else if withdrawRecord.Status == 3 {
			statusMsg = "Rejected"
		} else if withdrawRecord.Status == 4 {
			statusMsg = "Processing"
		} else if withdrawRecord.Status == 5 {
			statusMsg = "Failed"
		} else if withdrawRecord.Status == 6 {
			statusMsg = "Completed"
		}

		record := &exchange.WDHistory{
			ID:        withdrawRecord.ID,
			Coin:      c,
			Quantity:  quantity,
			Tag:       "",
			Address:   withdrawRecord.Address,
			TxHash:    withdrawRecord.TxID,
			ChainType: chainType,
			Status:    statusMsg,
			TimeStamp: withdrawRecord.ApplyTime.UnixNano(),
		}

		operation.WithdrawalHistory = append(operation.WithdrawalHistory, record)
	}

	return nil
}

func (e *Binance) doGetDepositHistory(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key or passphrase are nil.", e.GetName())
	}

	depositHistory := DepositHistory{}
	strRequest := "/sapi/v1/capital/deposit/hisrec"

	mapParams := make(map[string]string)

	jsonGetDepositHistory := e.ApiKeyGet(mapParams, strRequest)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		// operation.MapParams = fmt.Sprintf("%+v", mapParams)
		operation.CallResponce = jsonGetDepositHistory
	}

	if err := json.Unmarshal([]byte(jsonGetDepositHistory), &depositHistory); err != nil {
		operation.Error = fmt.Errorf("%s doGetDepositHistory Json Unmarshal Err: %v, %s", e.GetName(), err, jsonGetDepositHistory)
		return operation.Error
	}

	// store info into orders
	operation.DepositHistory = []*exchange.WDHistory{}
	for _, depositRecord := range depositHistory {
		c := e.GetCoinBySymbol(depositRecord.Coin)
		quantity, err := strconv.ParseFloat(depositRecord.Amount, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doGetDepositHistory parse quantity Err: %v, %v", e.GetName(), err, depositRecord.Amount)
			return operation.Error
		}
		var chainType exchange.ChainType
		if depositRecord.Network == "BTC" {
			chainType = exchange.MAINNET
		} else if depositRecord.Network == "ETH" {
			chainType = exchange.ERC20
		} else {
			chainType = exchange.OTHER
		}

		statusMsg := ""
		if depositRecord.Status == 0 {
			statusMsg = "Confirm email sent"
		} else if depositRecord.Status == 1 {
			statusMsg = "Canceled by user"
		} else if depositRecord.Status == 2 {
			statusMsg = "Waiting for Confirmation"
		} else if depositRecord.Status == 3 {
			statusMsg = "Rejected"
		} else if depositRecord.Status == 4 {
			statusMsg = "Processing"
		} else if depositRecord.Status == 5 {
			statusMsg = "Failed"
		} else if depositRecord.Status == 6 {
			statusMsg = "Completed"
		}

		record := &exchange.WDHistory{
			// ID:        depositRecord.ID,
			Coin:      c,
			Quantity:  quantity,
			Tag:       depositRecord.AddressTag,
			Address:   depositRecord.Address,
			TxHash:    depositRecord.TxID,
			ChainType: chainType,
			Status:    statusMsg,
			TimeStamp: depositRecord.InsertTime,
		}

		operation.DepositHistory = append(operation.DepositHistory, record)
	}

	return nil
}

func (e *Binance) doGeTransferHistory(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key or passphrase are nil.", e.GetName())
	}

	transfer := TransferHistory{}
	strRequest := "/sapi/v1/sub-account/transfer/subUserHistory"

	mapParams := make(map[string]string)

	jsonTransferOutHistory := e.ApiKeyGet(mapParams, strRequest)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonTransferOutHistory
	}

	if err := json.Unmarshal([]byte(jsonTransferOutHistory), &transfer); err != nil {
		operation.Error = fmt.Errorf("%s doTransferOutHistory Json Unmarshal Err: %v, %s", e.GetName(), err, jsonTransferOutHistory)
		return operation.Error
	}

	// store info into orders
	operation.TransferOutHistory = []*exchange.TransferHistory{}
	operation.TransferInHistory = []*exchange.TransferHistory{}
	for _, tx := range transfer {
		c := e.GetCoinBySymbol(tx.Asset)
		quantity, err := strconv.ParseFloat(tx.Qty, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doGetTransferHistory parse quantity Err: %v, %v", e.GetName(), err, tx.Qty)
			return operation.Error
		}

		record := &exchange.TransferHistory{
			Coin:      c,
			Quantity:  quantity,
			TimeStamp: tx.Time,
		}

		switch tx.Type {
		case 1:
			record.Type = exchange.TransferIn
			operation.TransferInHistory = append(operation.TransferInHistory, record)
		case 2:
			record.Type = exchange.TransferOut
			operation.TransferOutHistory = append(operation.TransferOutHistory, record)
		}
	}

	return nil
}

func (e *Binance) doGetDepositAddress(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key or passphrase are nil.", e.GetName())
	}

	strRequest := "/sapi/v1/capital/deposit/address"

	operation.DepositAddresses = make(map[exchange.ChainType]*exchange.DepositAddr)
	for _, network := range []string{"BTC", "ETH"} {
		depositAddress := DepositAddress{}
		mapParams := make(map[string]string)
		mapParams["coin"] = e.GetSymbolByCoin(operation.Coin)
		mapParams["network"] = network

		jsonGetDepositAddress := e.ApiKeyGet(mapParams, strRequest)
		if operation.DebugMode {
			operation.RequestURI = strRequest
			// operation.MapParams = fmt.Sprintf("%+v", mapParams)
			operation.CallResponce = jsonGetDepositAddress
		}

		if err := json.Unmarshal([]byte(jsonGetDepositAddress), &depositAddress); err != nil {
			operation.Error = fmt.Errorf("%s doGetDepositAddress Json Unmarshal Err: %v, %s", e.GetName(), err, jsonGetDepositAddress)
			return operation.Error
		} else if depositAddress.Code == -9000 { // no deposit addr
			log.Printf("%v Coin %v No deposit addr for network: %v", e.GetName(), mapParams["coin"], network)
			continue
		} else if depositAddress.Code != 0 {
			operation.Error = fmt.Errorf("%s doGetDepositAddress fail: %s", e.GetName(), jsonGetDepositAddress)
			return operation.Error
		}

		var chain exchange.ChainType
		if mapParams["network"] == "BTC" {
			chain = exchange.MAINNET
		} else if mapParams["network"] == "ETH" {
			chain = exchange.ERC20
		} else {
			chain = exchange.OTHER
		}

		// store info into orders
		depoAddr := &exchange.DepositAddr{
			Coin:    operation.Coin,
			Address: depositAddress.Address,
			Tag:     depositAddress.Tag,
			Chain:   chain,
		}

		operation.DepositAddresses[chain] = depoAddr
	}

	return nil
}

func (e *Binance) doWithdraw(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	withdraw := WithdrawResponse{}
	strRequest := "/wapi/v3/withdraw.html"

	mapParams := make(map[string]string)
	mapParams["asset"] = e.GetSymbolByCoin(operation.Coin)
	mapParams["address"] = operation.WithdrawAddress
	if operation.WithdrawTag != "" { //this part is not working yet
		mapParams["addressTag"] = operation.WithdrawTag
	}
	mapParams["amount"] = operation.WithdrawAmount
	mapParams["timestamp"] = fmt.Sprintf("%d", time.Now().UnixNano()/1e6)

	jsonSubmitWithdraw := e.WApiKeyRequest("POST", mapParams, strRequest)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		// operation.MapParams = fmt.Sprintf("%+v", mapParams)
		operation.CallResponce = jsonSubmitWithdraw
	}

	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &withdraw); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Json Unmarshal Err: %v, %s", e.GetName(), err, jsonSubmitWithdraw)
		return operation.Error
	}
	if !withdraw.Success {
		operation.Error = fmt.Errorf("%s Withdraw Failed: %v", e.GetName(), jsonSubmitWithdraw)
		return operation.Error
	}

	operation.WithdrawID = withdraw.ID

	return nil
}
