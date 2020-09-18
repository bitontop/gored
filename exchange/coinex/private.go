package coinex

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/bitontop/gored/exchange"
)

func (e *Coinex) DoAccountOperation(operation *exchange.AccountOperation) error {
	switch operation.Type {
	case exchange.SubAccountTransfer:
		if operation.Wallet == exchange.SpotWallet {
			return e.subTransfer(operation)
		}
	// case exchange.Transfer:
	// 	return e.transfer(operation)
	// case exchange.Balance:
	// 	return e.getBalance(operation)
	case exchange.BalanceList:
		if operation.Wallet == exchange.SpotWallet {
			return e.doAllBalance(operation)
		}

	case exchange.Withdraw:
		return e.doWithdraw(operation)

	case exchange.SubBalanceList:
		if operation.Wallet == exchange.SpotWallet {
			return e.doSubBalance(operation)
		}

	case exchange.SubAllBalanceList:
		if operation.Wallet == exchange.SpotWallet {
			return e.doSubAllBalance(operation) // All spot trading and main sub account
		}

	case exchange.GetOpenOrder:
		if operation.Wallet == exchange.SpotWallet {
			return e.doGetOpenOrder(operation)
		}
	case exchange.GetWithdrawalHistory:
		if operation.Wallet == exchange.SpotWallet {
			return e.doGetWithdrawalHistory(operation)
		}
	case exchange.GetDepositHistory:
		if operation.Wallet == exchange.SpotWallet {
			return e.doGetDepositHistory(operation)
		}
	case exchange.GetTransferHistory:
		if operation.Wallet == exchange.SpotWallet {
			return e.doGetTransferHistory(operation)
		}

	}
	return fmt.Errorf("%s Operation type invalid: %s %v", operation.Ex, operation.Wallet, operation.Type)
}

func (e *Coinex) subTransfer(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
	}

	jsonResponse := JsonResponse{}
	strRequest := "/v1/sub_account/transfer"

	mapParams := make(map[string]string)
	mapParams["access_id"] = e.API_KEY
	mapParams["coin_type"] = e.GetSymbolByCoin(operation.Coin)
	mapParams["amount"] = operation.SubTransferAmount

	if operation.SubTransferFrom != "" {
		mapParams["transfer_side"] = "in" // to post "in" or "out", in for deposit, out for withdrawal
		mapParams["transfer_account"] = operation.SubTransferFrom
	} else if operation.SubTransferTo != "" {
		mapParams["transfer_side"] = "out" // to post "in" or "out", in for deposit, out for withdrawal
		mapParams["transfer_account"] = operation.SubTransferTo
	} else {
		return fmt.Errorf("%s doSubTransfer failed, missing subAccount name", e.GetName())
	}

	jsonTransferReturn := e.ApiKeyRequest("POST", strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonTransferReturn
	}

	if err := json.Unmarshal([]byte(jsonTransferReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s doSubTransfer Json Unmarshal Err: %v, %s", e.GetName(), err, jsonTransferReturn)
		return operation.Error
	} else if jsonResponse.Code != 0 {
		operation.Error = fmt.Errorf("%s doSubTransfer failed: %v", e.GetName(), jsonTransferReturn)
		return operation.Error
	}

	// log.Printf("SubTransfer response %v", jsonTransferReturn)

	return nil
}

func (e *Coinex) doGetTransferHistory(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key or passphrase are nil.", e.GetName())
	}

	jsonResponse := JsonResponse{}
	transfer := TransferHistory{}
	strRequest := "/v1/sub_account/transfer/history"

	subUserName := operation.SubUserName

	mapParams := make(map[string]string)
	mapParams["access_id"] = e.API_KEY
	mapParams["sub_user_name"] = subUserName

	jsonTransferHistory := e.ApiKeyRequest("GET", strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonTransferHistory
	}

	if err := json.Unmarshal([]byte(jsonTransferHistory), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s doGetTransferHistory Json Unmarshal Err: %v, %s", e.GetName(), err, jsonTransferHistory)
		return operation.Error
	} else if jsonResponse.Code != 0 {
		operation.Error = fmt.Errorf("%s doGetTransferHistory failed: %v", e.GetName(), jsonTransferHistory)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &transfer); err != nil {
		operation.Error = fmt.Errorf("%s doGetTransferHistory Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	// store info into orders
	operation.TransferOutHistory = []*exchange.TransferHistory{}
	operation.TransferInHistory = []*exchange.TransferHistory{}
	for _, tx := range transfer.Data {
		c := e.GetCoinBySymbol(tx.CoinType)
		quantity, err := strconv.ParseFloat(tx.Amount, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doGetTransferHistory parse quantity Err: %v, %v", e.GetName(), err, tx.Amount)
			return operation.Error
		}

		record := &exchange.TransferHistory{
			Coin:      c,
			Quantity:  quantity,
			TimeStamp: tx.Time,
			StatusMsg: tx.Status,
		}

		if tx.TransferTo == subUserName {
			record.Type = exchange.TransferIn
			operation.TransferInHistory = append(operation.TransferInHistory, record)
		} else if tx.TransferFrom == subUserName {
			record.Type = exchange.TransferOut
			operation.TransferOutHistory = append(operation.TransferOutHistory, record)
		}
	}

	return nil
}

func (e *Coinex) doGetWithdrawalHistory(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key or passphrase are nil.", e.GetName())
	}

	jsonResponse := JsonResponse{}
	withdrawHistory := WithdrawHistory{}
	strRequest := "/v1/balance/coin/withdraw"

	mapParams := make(map[string]string)
	mapParams["access_id"] = e.API_KEY

	jsonGetWithdrawalHistory := e.ApiKeyRequest("GET", strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		// operation.MapParams = fmt.Sprintf("%+v", mapParams)
		operation.CallResponce = jsonGetWithdrawalHistory
	}

	if err := json.Unmarshal([]byte(jsonGetWithdrawalHistory), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s doGetWithdrawalHistory Json Unmarshal Err: %v, %s", e.GetName(), err, jsonGetWithdrawalHistory)
		return operation.Error
	} else if jsonResponse.Code != 0 {
		operation.Error = fmt.Errorf("%s doGetWithdrawalHistory failed: %v", e.GetName(), jsonGetWithdrawalHistory)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &withdrawHistory); err != nil {
		operation.Error = fmt.Errorf("%s doGetWithdrawalHistory Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	// store info into orders
	operation.WithdrawalHistory = []*exchange.WDHistory{}
	for _, withdrawRecord := range withdrawHistory {
		c := e.GetCoinBySymbol(withdrawRecord.CoinType)
		quantity, err := strconv.ParseFloat(withdrawRecord.Amount, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doGetWithdrawalHistory parse quantity Err: %v, %v", e.GetName(), err, withdrawRecord.Amount)
			return operation.Error
		}
		// var chainType exchange.ChainType
		// if withdrawRecord.Network == "BTC" {
		// 	chainType = exchange.MAINNET
		// } else if withdrawRecord.Network == "ETH" {
		// 	chainType = exchange.ERC20
		// } else {
		// 	chainType = exchange.OTHER
		// }

		statusMsg := withdrawRecord.Status

		record := &exchange.WDHistory{
			ID:       fmt.Sprintf("%v", withdrawRecord.CoinWithdrawID),
			Coin:     c,
			Quantity: quantity,
			Tag:      "",
			Address:  withdrawRecord.CoinAddress,
			TxHash:   withdrawRecord.TxID,
			// ChainType: chainType, // not provided
			Status:    statusMsg,
			TimeStamp: withdrawRecord.CreateTime,
		}

		operation.WithdrawalHistory = append(operation.WithdrawalHistory, record)
	}

	return nil
}

func (e *Coinex) doGetDepositHistory(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key or passphrase are nil.", e.GetName())
	}

	jsonResponse := JsonResponse{}
	depositHistory := DepositHistory{}
	strRequest := "/v1/balance/coin/deposit"

	mapParams := make(map[string]string)
	mapParams["access_id"] = e.API_KEY

	jsonGetDepositHistory := e.ApiKeyRequest("GET", strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		// operation.MapParams = fmt.Sprintf("%+v", mapParams)
		operation.CallResponce = jsonGetDepositHistory
	}

	if err := json.Unmarshal([]byte(jsonGetDepositHistory), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s doGetDepositHistory Json Unmarshal Err: %v, %s", e.GetName(), err, jsonGetDepositHistory)
		return operation.Error
	} else if jsonResponse.Code != 0 {
		operation.Error = fmt.Errorf("%s doGetDepositHistory failed: %v", e.GetName(), jsonGetDepositHistory)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &depositHistory); err != nil {
		operation.Error = fmt.Errorf("%s doGetDepositHistory Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	// store info into orders
	operation.DepositHistory = []*exchange.WDHistory{}
	for _, depositRecord := range depositHistory {
		c := e.GetCoinBySymbol(depositRecord.CoinType)
		quantity, err := strconv.ParseFloat(depositRecord.Amount, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doGetDepositHistory parse quantity Err: %v, %v", e.GetName(), err, depositRecord.Amount)
			return operation.Error
		}
		// var chainType exchange.ChainType
		// if depositRecord.Network == "BTC" {
		// 	chainType = exchange.MAINNET
		// } else if depositRecord.Network == "ETH" {
		// 	chainType = exchange.ERC20
		// } else {
		// 	chainType = exchange.OTHER
		// }

		statusMsg := depositRecord.Status

		record := &exchange.WDHistory{
			ID:       fmt.Sprintf("%v", depositRecord.CoinDepositID),
			Coin:     c,
			Quantity: quantity,
			Tag:      "",
			Address:  depositRecord.CoinAddress,
			TxHash:   depositRecord.TxID,
			// ChainType: chainType, // not provided
			Status:    statusMsg,
			TimeStamp: depositRecord.CreateTime,
		}

		operation.DepositHistory = append(operation.DepositHistory, record)
	}

	return nil
}

func (e *Coinex) doGetOpenOrder(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key or passphrase are nil.", e.GetName())
	}

	jsonResponse := JsonResponse{}
	openOrders := OpenOrders{}
	strRequest := "/v1/order/pending"

	mapParams := make(map[string]string)
	mapParams["access_id"] = e.API_KEY
	if operation.Pair != nil {
		mapParams["market"] = e.GetSymbolByPair(operation.Pair)
	}
	mapParams["page"] = "1"
	mapParams["limit"] = "100"

	jsonGetOpenOrder := e.ApiKeyRequest("GET", strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		// operation.MapParams = fmt.Sprintf("%+v", mapParams)
		operation.CallResponce = jsonGetOpenOrder
	}

	if err := json.Unmarshal([]byte(jsonGetOpenOrder), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s doGetOpenOrder Json Unmarshal Err: %v, %s", e.GetName(), err, jsonGetOpenOrder)
		return operation.Error
	} else if jsonResponse.Code != 0 {
		operation.Error = fmt.Errorf("%s doGetOpenOrder failed: %v", e.GetName(), jsonGetOpenOrder)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &openOrders); err != nil {
		operation.Error = fmt.Errorf("%s doGetOpenOrder Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	// store info into orders
	operation.OpenOrders = []*exchange.Order{}
	for _, o := range openOrders.Data {
		rate, err := strconv.ParseFloat(o.Price, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doGetOpenOrder parse rate Err: %v, %v", e.GetName(), err, o.Price)
			return operation.Error
		}
		quantity, err := strconv.ParseFloat(o.Amount, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doGetOpenOrder parse quantity Err: %v, %v", e.GetName(), err, o.Amount)
			return operation.Error
		}
		dealQuantity, err := strconv.ParseFloat(o.DealAmount, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doGetOpenOrder parse dealQuantity Err: %v, %v", e.GetName(), err, o.DealAmount)
			return operation.Error
		}

		order := &exchange.Order{
			Pair:         e.GetPairBySymbol(o.Market),
			OrderID:      fmt.Sprintf("%v", o.ID),
			Rate:         rate,
			Quantity:     quantity,
			DealRate:     rate,
			DealQuantity: dealQuantity,
			Timestamp:    o.CreateTime,
			// JsonResponse: jsonGetOpenOrder,
		}

		switch o.Type {
		case "buy":
			order.Direction = exchange.Buy
		case "sell":
			order.Direction = exchange.Sell
		}

		if o.Status == "done" {
			order.Status = exchange.Filled
		} else if o.Status == "part_deal" {
			order.Status = exchange.Partial
		} else if o.Status == "not_deal" {
			order.Status = exchange.New
		} else {
			order.Status = exchange.Other
		}

		operation.OpenOrders = append(operation.OpenOrders, order)
	}

	return nil
}

func (e *Coinex) doAllBalance(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := JsonResponse{}
	accountBalance := make(map[string]*AccountBalances)
	strRequest := "/v1/balance/info"
	operation.BalanceList = []exchange.AssetBalance{}

	mapParams := make(map[string]string)
	mapParams["access_id"] = e.API_KEY

	jsonAllBalanceReturn := e.ApiKeyRequest("GET", strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonAllBalanceReturn
	}

	if err := json.Unmarshal([]byte(jsonAllBalanceReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s doAllBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonAllBalanceReturn)
		return operation.Error
	} else if jsonResponse.Code != 0 {
		operation.Error = fmt.Errorf("%s doAllBalance failed: %v", e.GetName(), jsonAllBalanceReturn)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		operation.Error = fmt.Errorf("%s doAllBalance Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	for symbol, balance := range accountBalance {
		freeamount, err := strconv.ParseFloat(balance.Available, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doAllBalance err: %+v %v", e.GetName(), balance, err)
			return operation.Error
		}
		locked, err := strconv.ParseFloat(balance.Frozen, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doAllBalance err: %+v %v", e.GetName(), balance, err)
			return operation.Error
		}

		c := e.GetCoinBySymbol(symbol)
		if c == nil {
			continue
		}
		b := exchange.AssetBalance{
			Coin:             c,
			BalanceAvailable: freeamount,
			BalanceFrozen:    locked,
		}
		operation.BalanceList = append(operation.BalanceList, b)
	}

	return nil
}

func (e *Coinex) doSubAllBalance(operation *exchange.AccountOperation) error { // tested
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := JsonResponse{}
	accountBalance := SubAccountBalances{}
	strRequest := "/v1/sub_account/balance"

	mapParams := make(map[string]string)
	mapParams["access_id"] = e.API_KEY
	// mapParams["sub_user_name"] = url.QueryEscape(operation.SubAccountID) //

	jsonBalanceReturn := e.ApiKeyRequest("GET", strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonBalanceReturn
	}

	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s doSubBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonBalanceReturn)
		return operation.Error
	} else if jsonResponse.Code != 0 {
		operation.Error = fmt.Errorf("%s doSubBalance failed: %v", e.GetName(), jsonBalanceReturn)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		operation.Error = fmt.Errorf("%s doSubBalance Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	balanceMap := make(map[string]exchange.AssetBalance)
	operation.BalanceList = []exchange.AssetBalance{}
	for accountName, account := range accountBalance {
		if mapParams["sub_user_name"] != "" && mapParams["sub_user_name"] != accountName {
			continue
		}
		for symbol, balance := range account {
			freeamount, err := strconv.ParseFloat(balance.Available, 64)
			if err != nil {
				operation.Error = fmt.Errorf("%s UpdateSubBalances parse err: %+v %v", e.GetName(), balance, err)
				return operation.Error
			}
			locked, err := strconv.ParseFloat(balance.Frozen, 64)
			if err != nil {
				operation.Error = fmt.Errorf("%s UpdateSubBalances parse err: %+v %v", e.GetName(), balance, err)
				return operation.Error
			}

			c := e.GetCoinBySymbol(symbol)
			if c == nil {
				continue
			}
			b := exchange.AssetBalance{
				Coin:             c,
				BalanceAvailable: freeamount,
				BalanceFrozen:    locked,
			}

			// update balance for coin c
			oldBalance, ok := balanceMap[c.Code]
			if ok {
				b.BalanceAvailable += oldBalance.BalanceAvailable
				b.BalanceFrozen += oldBalance.BalanceFrozen
			}
			balanceMap[c.Code] = b
		}
	}

	// store aggregated balance into list
	for _, balance := range balanceMap {
		operation.BalanceList = append(operation.BalanceList, balance)
	}

	return nil
}

func (e *Coinex) doSubBalance(operation *exchange.AccountOperation) error { // tested
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := JsonResponse{}
	accountBalance := SubAccountBalances{}
	strRequest := "/v1/sub_account/balance"

	mapParams := make(map[string]string)
	mapParams["access_id"] = e.API_KEY
	mapParams["sub_user_name"] = url.QueryEscape(operation.SubAccountID) //

	jsonBalanceReturn := e.ApiKeyRequest("GET", strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonBalanceReturn
	}

	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s doSubBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonBalanceReturn)
		return operation.Error
	} else if jsonResponse.Code != 0 {
		operation.Error = fmt.Errorf("%s doSubBalance failed: %v", e.GetName(), jsonBalanceReturn)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		operation.Error = fmt.Errorf("%s doSubBalance Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	operation.BalanceList = []exchange.AssetBalance{}
	for accountName, account := range accountBalance {
		if mapParams["sub_user_name"] != "" && mapParams["sub_user_name"] != accountName {
			continue
		}
		for symbol, balance := range account {
			freeamount, err := strconv.ParseFloat(balance.Available, 64)
			if err != nil {
				operation.Error = fmt.Errorf("%s UpdateSubBalances parse err: %+v %v", e.GetName(), balance, err)
				return operation.Error
			}
			locked, err := strconv.ParseFloat(balance.Frozen, 64)
			if err != nil {
				operation.Error = fmt.Errorf("%s UpdateSubBalances parse err: %+v %v", e.GetName(), balance, err)
				return operation.Error
			}

			c := e.GetCoinBySymbol(symbol)
			if c == nil {
				continue
			}
			b := exchange.AssetBalance{
				Coin:             c,
				BalanceAvailable: freeamount,
				BalanceFrozen:    locked,
			}
			operation.BalanceList = append(operation.BalanceList, b)
		}

	}

	return nil
}

func (e *Coinex) doWithdraw(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("coinex API Key or Secret Key are nil.")
	}

	jsonResponse := JsonResponse{}
	withdraw := Withdraw{}
	strRequestUrl := "/v1/balance/coin/withdraw"

	mapParams := make(map[string]string)
	mapParams["access_id"] = e.API_KEY
	mapParams["coin_type"] = e.GetSymbolByCoin(operation.Coin)
	mapParams["transfer_method"] = "onchain"
	mapParams["actual_amount"] = operation.WithdrawAmount

	if operation.WithdrawTag != "" {
		mapParams["coin_address"] = fmt.Sprintf("%s:%s", operation.WithdrawAddress, operation.WithdrawTag)
	} else {
		mapParams["coin_address"] = operation.WithdrawAddress
	}

	jsonWithdraw := e.ApiKeyPost(strRequestUrl, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequestUrl
		operation.CallResponce = jsonWithdraw
	}

	if err := json.Unmarshal([]byte(jsonWithdraw), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Json Unmarshal Err: %v, %s", e.GetName(), err, jsonWithdraw)
		return operation.Error
	} else if jsonResponse.Code != 0 {
		operation.Error = fmt.Errorf("%s Withdraw Failed: %v", e.GetName(), jsonWithdraw)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &withdraw); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	operation.WithdrawID = fmt.Sprintf("%v", withdraw.CoinWithdrawID)

	return nil
}
