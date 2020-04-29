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

	// case exchange.Transfer:
	// 	return e.transfer(operation)
	// case exchange.BalanceList:
	// 	return e.getAllBalance(operation)
	// case exchange.Balance:
	// 	return e.getBalance(operation)

	case exchange.Withdraw:
		return e.doWithdraw(operation)

	case exchange.SubBalanceList:
		if operation.Wallet == exchange.SpotWallet {
			return e.doSubBalance(operation)
		}

	}
	return fmt.Errorf("Operation type invalid: %v", operation.Type)
}

// could also get all sub
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
