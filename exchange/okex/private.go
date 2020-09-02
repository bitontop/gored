package okex

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/bitontop/gored/exchange"
)

func (e *Okex) DoAccountOperation(operation *exchange.AccountOperation) error {
	switch operation.Type {
	case exchange.SubAccountTransfer:
		if operation.Wallet == exchange.SpotWallet {
			return e.subTransfer(operation)
		}
	case exchange.Transfer:
		if operation.Wallet == exchange.SpotWallet {
			return e.transfer(operation)
		}
	case exchange.BalanceList:
		if operation.Wallet == exchange.AssetWallet || operation.Wallet == exchange.SpotWallet {
			return e.getAllBalance(operation)
		}
	case exchange.Balance:
		if operation.Wallet == exchange.SpotWallet {
			return e.getBalance(operation)
		}
	case exchange.Withdraw:
		if operation.Wallet == exchange.SpotWallet {
			return e.doWithdraw(operation)
		}
	case exchange.GetTransferHistory:
		if operation.Wallet == exchange.SpotWallet {
			return e.doGetTransferHistory(operation)
		}
	case exchange.GetOpenOrder:
		if operation.Wallet == exchange.SpotWallet {
			return e.getOpenOrder(operation)
		}

	}

	return fmt.Errorf("%s Operation type invalid: %s %v", operation.Ex, operation.Wallet, operation.Type)
}

// all types of transfer in same api, all types of wallets. Only support spot now.
// put subAccount login name into 'SubTransferFrom' or 'SubTransferTo'
func (e *Okex) subTransfer(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
	}

	trans := Transfer{}
	strRequest := "/api/account/v3/transfer"

	mapParams := make(map[string]interface{})
	mapParams["currency"] = e.GetSymbolByCoin(operation.Coin)
	mapParams["amount"] = operation.SubTransferAmount
	mapParams["from"] = "1" // 1 spot, 6 asset
	mapParams["to"] = "1"
	if operation.SubTransferFrom != "" {
		mapParams["type"] = "2"
		mapParams["sub_account"] = operation.SubTransferFrom
	} else if operation.SubTransferTo != "" {
		mapParams["type"] = "1"
		mapParams["sub_account"] = operation.SubTransferTo
	} else {
		return fmt.Errorf("%s doSubTransfer failed, missing subAccount param", e.GetName())
	}

	jsonTransferReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonTransferReturn
	}

	if err := json.Unmarshal([]byte(jsonTransferReturn), &trans); err != nil {
		errorJson := ErrorMsg{}
		if err := json.Unmarshal([]byte(jsonTransferReturn), &errorJson); err != nil {
			operation.Error = fmt.Errorf("%s doSubTransfer Err: %v", e.GetName(), jsonTransferReturn)
			return operation.Error
		} else {
			operation.Error = fmt.Errorf("%s doSubTransfer Json Unmarshal Err: %v, %s", e.GetName(), err, jsonTransferReturn)
			return operation.Error
		}
	} else if !trans.Result {
		operation.Error = fmt.Errorf("%s doSubTransfer failed: %v", e.GetName(), jsonTransferReturn)
		return operation.Error
	}

	// log.Printf("SubTransfer response %v", jsonTransferReturn)

	return nil
}

func (e *Okex) getOpenOrder(op *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	openOrders := OpenOrders{}
	strRequest := "/api/spot/v3/orders_pending"

	mapParams := make(map[string]interface{})
	if op.Pair != nil {
		mapParams["instrument_id"] = e.GetSymbolByPair(op.Pair)
	}

	jsonOrders := e.ApiKeyRequest("GET", mapParams, strRequest)
	if op.DebugMode {
		op.RequestURI = strRequest
		op.CallResponce = jsonOrders
	}

	if err := json.Unmarshal([]byte(jsonOrders), &openOrders); err != nil {
		op.Error = fmt.Errorf("%s Get OpenOrders Json Unmarshal Err: %v, %s", e.GetName(), err, jsonOrders)
		return op.Error
	}

	result := []*exchange.Order{}
	for _, data := range openOrders {
		order := &exchange.Order{
			Pair:      e.GetPairBySymbol(data.InstrumentID),
			OrderID:   data.OrderID,
			Timestamp: data.Timestamp.UnixNano(),
		}

		switch data.Side {
		case "buy":
			order.Direction = exchange.Buy
		case "sell":
			order.Direction = exchange.Sell
		}

		order.Quantity, _ = strconv.ParseFloat(data.Size, 64)
		order.Rate, _ = strconv.ParseFloat(data.Price, 64)

		order.DealRate = order.Rate
		order.DealQuantity = 0.0

		dealQ, _ := strconv.ParseFloat(data.FilledSize, 64)
		dealTotal, _ := strconv.ParseFloat(data.FilledNotional, 64)
		if dealQ > 0 {
			order.DealQuantity = dealQ
		}
		if dealTotal > 0 && dealQ > 0 {
			order.DealRate = dealTotal / dealQ
		}

		if data.State == "0" && order.DealQuantity == 0 {
			order.Status = exchange.New
		} else if data.State == "0" && order.DealQuantity < order.Quantity {
			order.Status = exchange.Partial
		} else {
			order.Status = exchange.Other
		}

		result = append(result, order)
	}
	op.OpenOrders = result

	return nil
}

// only 1 month data
func (e *Okex) doGetTransferHistory(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key or passphrase are nil.", e.GetName())
	}

	transfer := TransferHistory{}
	strRequest := "/api/account/v3/ledger"

	// 	//============
	// strRequest = "/api/account/v3/deposit/history"
	// 	// ==============

	mapParams := make(map[string]interface{})
	if operation.TransferStartTime != 0 {
		mapParams["after"] = operation.TransferStartTime
	}
	if operation.TransferEndTime != 0 {
		mapParams["before"] = operation.TransferEndTime
	}

	jsonTransferOutHistory := e.ApiKeyRequest("GET", mapParams, strRequest)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonTransferOutHistory
	}

	if err := json.Unmarshal([]byte(jsonTransferOutHistory), &transfer); err != nil {
		operation.Error = fmt.Errorf("%s doGetTransferHistory Json Unmarshal Err: %v, %s", e.GetName(), err, jsonTransferOutHistory)
		return operation.Error
	}

	// store info into orders
	operation.TransferOutHistory = []*exchange.TransferHistory{}
	operation.TransferInHistory = []*exchange.TransferHistory{}
	for _, tx := range transfer {
		c := e.GetCoinBySymbol(tx.Currency)
		quantity, err := strconv.ParseFloat(tx.Amount, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s doGetTransferHistory parse quantity Err: %v, %v", e.GetName(), err, tx.Amount)
			return operation.Error
		}

		record := &exchange.TransferHistory{
			ID:        tx.LedgerID,
			Coin:      c,
			Quantity:  quantity,
			TimeStamp: tx.Timestamp.UnixNano(),
		}

		switch tx.Typename {
		case "To: subaccount":
			record.Type = exchange.TransferIn
			operation.TransferInHistory = append(operation.TransferInHistory, record)
		case "From: subaccount":
			record.Type = exchange.TransferOut
			record.Quantity *= -1
			operation.TransferOutHistory = append(operation.TransferOutHistory, record)
		default:
			continue
		}
	}

	return nil
}

func (e *Okex) doWithdraw(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
	}

	if operation.WithdrawTag != "" {
		operation.Error = fmt.Errorf("%s Withdraw Failed, got tag: %v, for coin: %v", e.GetName(), operation.WithdrawTag, operation.Coin.Code)
		return operation.Error
	}

	withdrawResponse := WithdrawResponse{}
	strRequest := "/api/account/v3/withdrawal"

	mapParams := make(map[string]interface{})
	mapParams["currency"] = e.GetSymbolByCoin(operation.Coin)
	mapParams["amount"] = operation.WithdrawAmount
	mapParams["destination"] = "4"
	mapParams["to_address"] = operation.WithdrawAddress
	mapParams["trade_pwd"] = e.TradePassword
	mapParams["fee"] = e.GetTxFee(operation.Coin)

	jsonSubmitWithdraw := e.ApiKeyRequest("POST", mapParams, strRequest)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonSubmitWithdraw
	}

	if err := json.Unmarshal([]byte(jsonSubmitWithdraw), &withdrawResponse); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Json Unmarshal Err: %v, %s", e.GetName(), err, jsonSubmitWithdraw)
		return operation.Error
	} else if !withdrawResponse.Result {
		operation.Error = fmt.Errorf("%s Withdraw Failed: %v", e.GetName(), jsonSubmitWithdraw)
		return operation.Error
	} else if withdrawResponse.WithdrawalID == "" {
		operation.Error = fmt.Errorf("%s Withdraw Failed: %v", e.GetName(), jsonSubmitWithdraw)
		return operation.Error
	}

	operation.WithdrawID = withdrawResponse.WithdrawalID

	return nil
}

func (e *Okex) transfer(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
	}

	trans := Transfer{}
	strRequest := "/api/account/v3/transfer"

	mapParams := make(map[string]interface{})
	mapParams["currency"] = e.GetSymbolByCoin(operation.Coin)
	mapParams["amount"] = operation.TransferAmount
	switch operation.TransferFrom {
	case exchange.AssetWallet:
		mapParams["from"] = "6"
	case exchange.SpotWallet:
		mapParams["from"] = "1"
	}
	switch operation.TransferDestination {
	case exchange.AssetWallet:
		mapParams["to"] = "6"
	case exchange.SpotWallet:
		mapParams["to"] = "1"
	}

	jsonTransferReturn := e.ApiKeyRequest("POST", mapParams, strRequest)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonTransferReturn
	}

	if err := json.Unmarshal([]byte(jsonTransferReturn), &trans); err != nil {
		errorJson := ErrorMsg{}
		if err := json.Unmarshal([]byte(jsonTransferReturn), &errorJson); err != nil {
			operation.Error = fmt.Errorf("%s Transfer Err: %v", e.GetName(), jsonTransferReturn)
			return operation.Error
		} else {
			operation.Error = fmt.Errorf("%s Transfer Json Unmarshal Err: %v, %s", e.GetName(), err, jsonTransferReturn)
			return operation.Error
		}
	} else if !trans.Result {
		operation.Error = fmt.Errorf("%s Transfer failed: %v", e.GetName(), jsonTransferReturn)
		return operation.Error
	}
	log.Printf("%s Transfer return: %+v", e.GetName(), jsonTransferReturn)

	return nil
}

func (e *Okex) getAllBalance(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	balance := AssetBalance{}
	operation.BalanceList = []exchange.AssetBalance{}

	strRequest := ""
	switch operation.Wallet {
	case exchange.AssetWallet:
		strRequest = "/api/account/v3/wallet" // asset api
	case exchange.SpotWallet:
		strRequest = "/api/spot/v3/accounts" // coin api
	}

	jsonAllBalanceReturn := e.ApiKeyRequest("GET", nil, strRequest)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonAllBalanceReturn
	}

	// log.Printf("jsonAllBalanceReturn: %v", jsonAllBalanceReturn) //====================
	if jsonAllBalanceReturn == "[]" {
		// log.Printf("getAllBalance empty return: %v", jsonAllBalanceReturn)
		for _, c := range e.GetCoins() { // set all coin balance to 0
			b := exchange.AssetBalance{
				Coin: c,
			}
			operation.BalanceList = append(operation.BalanceList, b)
		}
		return nil
	} else if err := json.Unmarshal([]byte(jsonAllBalanceReturn), &balance); err != nil {
		errorJson := ErrorMsg{}
		if err := json.Unmarshal([]byte(jsonAllBalanceReturn), &errorJson); err != nil {
			operation.Error = fmt.Errorf("%s getAllBalance Err: %v", e.GetName(), jsonAllBalanceReturn)
			return operation.Error
		} else {
			operation.Error = fmt.Errorf("%s getAllBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonAllBalanceReturn)
			return operation.Error
		}
	}

	for _, account := range balance {
		// if account.Balance == "0" {
		// 	continue
		// }
		frozen, err := strconv.ParseFloat(account.Hold, 64)
		available, err := strconv.ParseFloat(account.Available, 64)
		if err != nil {
			return fmt.Errorf("%s balance parse fail: %v %+v", e.GetName(), err, balance)
		}

		b := exchange.AssetBalance{
			Coin:             e.GetCoinBySymbol(account.Currency),
			BalanceAvailable: available,
			BalanceFrozen:    frozen,
		}
		operation.BalanceList = append(operation.BalanceList, b)
	}

	return nil
	// return fmt.Errorf("%s getBalance get %v account balance fail: %v", e.GetName(), symbol, jsonBalanceReturn)
}

func (e *Okex) getBalance(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	symbol := e.GetSymbolByCoin(operation.Coin)
	balance := AssetBalance{}
	strRequest := fmt.Sprintf("/api/account/v3/wallet/%s", e.GetSymbolByCoin(operation.Coin))

	jsonBalanceReturn := e.ApiKeyRequest("GET", nil, strRequest)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonBalanceReturn
	}

	if jsonBalanceReturn == "[]" {
		log.Printf("getBalance empty return: %v", jsonBalanceReturn)
		operation.BalanceFrozen = 0
		operation.BalanceAvailable = 0
		return nil
	} else if err := json.Unmarshal([]byte(jsonBalanceReturn), &balance); err != nil {
		errorJson := ErrorMsg{}
		if err := json.Unmarshal([]byte(jsonBalanceReturn), &errorJson); err != nil {
			operation.Error = fmt.Errorf("%s getBalance Err: %v", e.GetName(), jsonBalanceReturn)
			return operation.Error
		} else {
			operation.Error = fmt.Errorf("%s getBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonBalanceReturn)
			return operation.Error
		}
	}

	for _, account := range balance {
		if account.Currency == symbol {
			frozen, err := strconv.ParseFloat(account.Hold, 64)
			available, err := strconv.ParseFloat(account.Available, 64)
			if err != nil {
				return fmt.Errorf("%s balance parse fail: %v %+v", e.GetName(), err, balance)
			}
			operation.BalanceFrozen = frozen
			operation.BalanceAvailable = available
			return nil
		}
	}

	operation.Error = fmt.Errorf("%s getBalance get %v account balance fail: %v", e.GetName(), symbol, jsonBalanceReturn)
	return operation.Error
}
