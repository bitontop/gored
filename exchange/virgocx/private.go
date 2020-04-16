package virgocx

import (
	"encoding/json"
	"fmt"

	"github.com/bitontop/gored/exchange"
)

func (e *Virgocx) DoAccountOperation(operation *exchange.AccountOperation) error {
	switch operation.Type {
	case exchange.BalanceList:
		return e.getAllBalance(operation)
	case exchange.Balance:
		return e.getBalance(operation)
	case exchange.PlaceOrder:
		return e.doPlaceOrder(operation)
	}
	return fmt.Errorf("Operation type invalid: %v", operation.Type)
}

func (e *Virgocx) doPlaceOrder(operation *exchange.AccountOperation) error {

	return nil
}

func (e *Virgocx) getAllBalance(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	jsonResponse := JsonResponse{}
	balance := AccountBalances{}
	strRequest := "/member/accounts"

	mapParams := make(map[string]string)

	jsonAllBalanceReturn := e.ApiKeyRequest("GET", strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonAllBalanceReturn
	}

	if err := json.Unmarshal([]byte(jsonAllBalanceReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s getAllBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonAllBalanceReturn)
		return operation.Error
	} else if jsonResponse.Code != 0 {
		operation.Error = fmt.Errorf("%s getAllBalance Failed: %v", e.GetName(), jsonAllBalanceReturn)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &balance); err != nil {
		operation.Error = fmt.Errorf("%s getAllBalance Data Unmarshal Err: %v, %s", e.GetName(), err, jsonAllBalanceReturn)
		return operation.Error
	}

	for _, account := range balance {
		if account.Total == 0 {
			continue
		}

		b := exchange.AssetBalance{
			Coin:             e.GetCoinBySymbol(account.CoinName),
			BalanceAvailable: account.Balance,
			BalanceFrozen:    account.FreezingBalance,
		}
		operation.BalanceList = append(operation.BalanceList, b)
	}

	return nil
}

func (e *Virgocx) getBalance(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil.", e.GetName())
	}

	symbol := e.GetSymbolByCoin(operation.Coin)
	jsonResponse := JsonResponse{}
	balance := AccountBalances{}
	strRequest := "/member/accounts"

	mapParams := make(map[string]string)

	jsonBalanceReturn := e.ApiKeyRequest("GET", strRequest, mapParams)
	if operation.DebugMode {
		operation.RequestURI = strRequest
		operation.CallResponce = jsonBalanceReturn
	}

	if err := json.Unmarshal([]byte(jsonBalanceReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s getBalance Json Unmarshal Err: %v, %s", e.GetName(), err, jsonBalanceReturn)
		return operation.Error
	} else if jsonResponse.Code != 0 {
		operation.Error = fmt.Errorf("%s getBalance Failed: %v", e.GetName(), jsonBalanceReturn)
		return operation.Error
	}
	if err := json.Unmarshal(jsonResponse.Data, &balance); err != nil {
		operation.Error = fmt.Errorf("%s getBalance Data Unmarshal Err: %v, %s", e.GetName(), err, jsonBalanceReturn)
		return operation.Error
	}

	for _, account := range balance {
		if account.CoinName == symbol {
			operation.BalanceFrozen = account.FreezingBalance
			operation.BalanceAvailable = account.Balance
			return nil
		}
	}

	operation.Error = fmt.Errorf("%s getBalance get %v account balance fail: %v", e.GetName(), symbol, jsonBalanceReturn)
	return operation.Error
}
