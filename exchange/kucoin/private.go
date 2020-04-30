package kucoin

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/bitontop/gored/exchange"
)

func (e *Kucoin) DoAccountOperation(operation *exchange.AccountOperation) error {
	switch operation.Type {
	case exchange.Transfer:
		return e.transfer(operation)
	case exchange.BalanceList:
		return e.getAllBalance(operation)
	case exchange.Balance:
		return e.getBalance(operation)
	case exchange.Withdraw:
		return e.doWithdraw(operation)

	case exchange.SubBalanceList:
		if operation.Wallet == exchange.SpotWallet {
			return e.doSubBalance(operation) // spot trading sub account
		}
	case exchange.GetSubAccountList:
		// if operation.Wallet == exchange.SpotWallet {
		return e.doSubAccountList(operation) // all type sub account
	// }
	case exchange.SubAllBalanceList:
		if operation.Wallet == exchange.SpotWallet {
			return e.doSubAllBalance(operation) // All spot trading and main sub account
		}
	}
	return fmt.Errorf("Operation type invalid: %v", operation.Type)
}

func (e *Kucoin) doSubAllBalance(operation *exchange.AccountOperation) error { //TODO, test with sub account
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := JsonResponse{}
	accountBalance := SubAllAccountBalances{}
	strRequest := "/api/v1/sub-accounts"

	mapParams := make(map[string]string)

	jsonBalanceReturn := e.ApiKeyRequest("GET", strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonBalanceReturn
	}

	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s doSubAllBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonBalanceReturn)
		return operation.Error
	} else if jsonResponse.Code != "200000" {
		operation.Error = fmt.Errorf("%s doSubAllBalance Failed: %v", e.GetName(), jsonBalanceReturn)
		return operation.Error
	}

	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		operation.Error = fmt.Errorf("%s doSubAllBalance Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	balanceMap := make(map[string]exchange.AssetBalance)
	operation.BalanceList = []exchange.AssetBalance{}
	for _, account := range accountBalance {
		// trading accounts
		for _, balance := range account.TradeAccounts {
			freeAmount, err := strconv.ParseFloat(balance.Available, 64)
			if err != nil {
				operation.Error = fmt.Errorf("%s UpdateSubBalances parse err: %+v %v", e.GetName(), balance, err)
				return operation.Error
			}
			total, err := strconv.ParseFloat(balance.Balance, 64)
			if err != nil {
				operation.Error = fmt.Errorf("%s UpdateSubBalances parse err: %+v %v", e.GetName(), balance, err)
				return operation.Error
			}

			c := e.GetCoinBySymbol(balance.Currency)
			if c == nil {
				continue
			}
			b := exchange.AssetBalance{
				Coin:             c,
				BalanceAvailable: freeAmount,
				BalanceFrozen:    total - freeAmount,
			}

			// update balance for coin c
			oldBalance, ok := balanceMap[c.Code]
			if ok {
				b.BalanceAvailable += oldBalance.BalanceAvailable
				b.BalanceFrozen += oldBalance.BalanceFrozen
			}
			balanceMap[c.Code] = b

		}

		// main accounts
		for _, balance := range account.MainAccounts {
			freeAmount, err := strconv.ParseFloat(balance.Available, 64)
			if err != nil {
				operation.Error = fmt.Errorf("%s UpdateSubBalances parse err: %+v %v", e.GetName(), balance, err)
				return operation.Error
			}
			total, err := strconv.ParseFloat(balance.Balance, 64)
			if err != nil {
				operation.Error = fmt.Errorf("%s UpdateSubBalances parse err: %+v %v", e.GetName(), balance, err)
				return operation.Error
			}

			c := e.GetCoinBySymbol(balance.Currency)
			if c == nil {
				continue
			}
			b := exchange.AssetBalance{
				Coin:             c,
				BalanceAvailable: freeAmount,
				BalanceFrozen:    total - freeAmount,
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

func (e *Kucoin) doSubAccountList(operation *exchange.AccountOperation) error { //TODO, test with sub account
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := JsonResponse{}
	accountList := SubAccountList{}
	strRequest := "/api/v1/sub/user"

	mapParams := make(map[string]string)

	jsonSubAccountReturn := e.ApiKeyRequest("GET", strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonSubAccountReturn
	}

	if err := json.Unmarshal([]byte(jsonSubAccountReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s doSubAccountList Json Unmarshal Err: %v, %s", e.GetName(), err, jsonSubAccountReturn)
		return operation.Error
	} else if jsonResponse.Code != "200000" {
		operation.Error = fmt.Errorf("%s doSubAccountList Failed: %v", e.GetName(), jsonSubAccountReturn)
		return operation.Error
	}

	if err := json.Unmarshal(jsonResponse.Data, &accountList); err != nil {
		operation.Error = fmt.Errorf("%s doSubAccountList Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	operation.SubAccountList = []*exchange.SubAccountInfo{}
	for _, account := range accountList {

		a := &exchange.SubAccountInfo{
			ID: account.UserID,
			// Status:    account.Status,
			Activated: true,
			// AccountType: exchange.SpotWallet,
			// TimeStamp: account.CreateTime,
		}
		operation.SubAccountList = append(operation.SubAccountList, a)
	}

	return nil
}

// for spot tradeAccounts balance, no mainAccounts
func (e *Kucoin) doSubBalance(operation *exchange.AccountOperation) error { //TODO, test with sub account
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := JsonResponse{}
	accountBalance := SubAccountBalances{}
	strRequest := fmt.Sprintf("/api/v1/sub-accounts/%v", operation.SubAccountID)

	mapParams := make(map[string]string)
	mapParams["subUserId"] = operation.SubAccountID

	jsonBalanceReturn := e.ApiKeyRequest("GET", strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonBalanceReturn
	}

	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s doSubBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonBalanceReturn)
		return operation.Error
	} else if jsonResponse.Code != "200000" {
		operation.Error = fmt.Errorf("%s doSubBalance Failed: %v", e.GetName(), jsonBalanceReturn)
		return operation.Error
	}

	if err := json.Unmarshal(jsonResponse.Data, &accountBalance); err != nil {
		operation.Error = fmt.Errorf("%s doSubBalance Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	operation.BalanceList = []exchange.AssetBalance{}
	for _, balance := range accountBalance.TradeAccounts {
		freeAmount, err := strconv.ParseFloat(balance.Available, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s UpdateSubBalances parse err: %+v %v", e.GetName(), balance, err)
			return operation.Error
		}
		total, err := strconv.ParseFloat(balance.Balance, 64)
		if err != nil {
			operation.Error = fmt.Errorf("%s UpdateSubBalances parse err: %+v %v", e.GetName(), balance, err)
			return operation.Error
		}

		c := e.GetCoinBySymbol(balance.Currency)
		if c == nil {
			continue
		}
		b := exchange.AssetBalance{
			Coin:             c,
			BalanceAvailable: freeAmount,
			BalanceFrozen:    total - freeAmount,
		}
		operation.BalanceList = append(operation.BalanceList, b)
	}

	return nil
}

func (e *Kucoin) doWithdraw(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("Kucoin API Key or Secret Key or passphrase are nil.")
	}

	if operation.WithdrawTag != "" {
		operation.Error = fmt.Errorf("%s Withdraw Failed, got tag: %v, for coin: %v", e.GetName(), operation.WithdrawTag, operation.Coin.Code)
		return operation.Error
	}

	// need to use inner transfer before withdraw

	jsonResponse := JsonResponse{}
	withdraw := Withdraw{}
	strRequestUrl := "/api/v1/withdrawals"

	mapParams := make(map[string]string)
	mapParams["currency"] = e.GetSymbolByCoin(operation.Coin)
	mapParams["address"] = operation.WithdrawAddress
	mapParams["amount"] = operation.WithdrawAmount

	jsonCreateWithdraw := e.ApiKeyRequest("POST", strRequestUrl, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequestUrl
		operation.CallResponce = jsonCreateWithdraw
	}

	if err := json.Unmarshal([]byte(jsonCreateWithdraw), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Json Unmarshal Err: %v, %s", e.GetName(), err, jsonCreateWithdraw)
		return operation.Error
	} else if jsonResponse.Code != "200000" {
		operation.Error = fmt.Errorf("%s Withdraw Failed: %v", e.GetName(), jsonCreateWithdraw)
		return operation.Error
	}

	if err := json.Unmarshal(jsonResponse.Data, &withdraw); err != nil {
		operation.Error = fmt.Errorf("%s Withdraw Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	operation.WithdrawID = withdraw.WithdrawalID

	return nil
}

func (e *Kucoin) transfer(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
	}

	innerTrans := InnerTrans{}
	strRequestUrl := "/api/v2/accounts/inner-transfer"

	mapParams := make(map[string]string)
	mapParams["clientOid"] = "12345"
	mapParams["currency"] = e.GetSymbolByCoin(operation.Coin)
	mapParams["amount"] = operation.TransferAmount
	switch operation.TransferFrom {
	case exchange.AssetWallet:
		mapParams["from"] = "main"
	case exchange.SpotWallet:
		mapParams["from"] = "trade"
	}
	switch operation.TransferDestination {
	case exchange.AssetWallet:
		mapParams["to"] = "main"
	case exchange.SpotWallet:
		mapParams["to"] = "trade"
	}

	jsonTransferReturn := e.ApiKeyRequest("POST", strRequestUrl, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequestUrl
		operation.CallResponce = jsonTransferReturn
	}

	// log.Printf("return: %v", jsonTransferReturn)
	if err := json.Unmarshal([]byte(jsonTransferReturn), &innerTrans); err != nil {
		operation.Error = fmt.Errorf("%s InnerTrans Json Unmarshal Err: %v, %s", e.GetName(), err, jsonTransferReturn)
		return operation.Error
	} else if innerTrans.Code != "200000" {
		operation.Error = fmt.Errorf("%s InnerTrans Failed: %v", e.GetName(), jsonTransferReturn)
		return operation.Error
	} else if innerTrans.Msg != "" {
		operation.Error = fmt.Errorf("%s InnerTrans Failed: %v", e.GetName(), jsonTransferReturn)
		return operation.Error
	}

	log.Printf("InnerTrans response %v", jsonTransferReturn)

	return nil
}

func (e *Kucoin) getAllBalance(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	accountID := AccountID{}
	strRequest := "/api/v1/accounts"
	accountType := ""
	// balanceList := []exchange.AssetBalance{}

	mapParams := make(map[string]string)
	if operation.Wallet == exchange.AssetWallet {
		mapParams["type"] = "main" // "trade"
		accountType = "main"
	} else if operation.Wallet == exchange.SpotWallet {
		mapParams["type"] = "trade"
		accountType = "trade"
	}

	jsonAllBalanceReturn := e.ApiKeyRequest("GET", strRequest, nil)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonAllBalanceReturn
	}

	// log.Printf("jsonAllBalanceReturn: %v", jsonAllBalanceReturn)
	if err := json.Unmarshal([]byte(jsonAllBalanceReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s getAllBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonAllBalanceReturn)
		return operation.Error
	} else if jsonResponse.Code != "200000" {
		operation.Error = fmt.Errorf("%s getAllBalance Failed: %v", e.GetName(), jsonAllBalanceReturn)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountID); err != nil {
		operation.Error = fmt.Errorf("%s getAllBalance Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	for _, account := range accountID {
		if account.Balance == "0" {
			continue
		}
		if account.Type == accountType {
			frozen, err := strconv.ParseFloat(account.Holds, 64)
			avaliable, err := strconv.ParseFloat(account.Available, 64)
			if err != nil {
				return fmt.Errorf("%s balance parse fail: %v %+v", e.GetName(), err, account)
			}

			balance := exchange.AssetBalance{
				Coin:             e.GetCoinBySymbol(account.Currency),
				BalanceAvailable: avaliable,
				BalanceFrozen:    frozen,
			}
			operation.BalanceList = append(operation.BalanceList, balance)
		}
	}

	return nil
	// return fmt.Errorf("%s getBalance fail: %v", e.GetName(), jsonBalanceReturn)
}

func (e *Kucoin) getBalance(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	accountID := AccountID{}
	strRequest := "/api/v1/accounts"
	accountType := ""

	mapParams := make(map[string]string)
	mapParams["currency"] = e.GetSymbolByCoin(operation.Coin)
	if operation.Wallet == exchange.AssetWallet {
		mapParams["type"] = "main" // "trade"
		accountType = "main"
	} else if operation.Wallet == exchange.SpotWallet {
		mapParams["type"] = "trade"
		accountType = "trade"
	}

	jsonBalanceReturn := e.ApiKeyRequest("GET", strRequest, nil)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonBalanceReturn
	}

	// log.Printf("jsonBalanceReturn: %v", jsonBalanceReturn)
	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s getBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonBalanceReturn)
		return operation.Error
	} else if jsonResponse.Code != "200000" {
		operation.Error = fmt.Errorf("%s getBalance Failed: %v", e.GetName(), jsonBalanceReturn)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &accountID); err != nil {
		operation.Error = fmt.Errorf("%s getBalance Result Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	}

	for _, account := range accountID {
		if account.Type == accountType && account.Currency == mapParams["currency"] {
			frozen, err := strconv.ParseFloat(account.Holds, 64)
			avaliable, err := strconv.ParseFloat(account.Available, 64)
			if err != nil {
				return fmt.Errorf("%s balance parse fail: %v %+v", e.GetName(), err, account)
			}
			operation.BalanceFrozen = frozen
			operation.BalanceAvailable = avaliable
			return nil
		}
	}

	return fmt.Errorf("%s getBalance fail: %v", e.GetName(), jsonBalanceReturn)
}
