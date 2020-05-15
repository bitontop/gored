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
	return fmt.Errorf("%s Operation type invalid: %s %v", operation.Ex, operation.Wallet, operation.Type)
}

func (e *Virgocx) doPlaceOrder(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" {
		return fmt.Errorf("%s API Key or Secret Key are nil", e.GetName())
	}

	mapParams := make(map[string]string)
	mapParams["symbol"] = e.GetSymbolByPair(operation.Pair)
	mapParams["qty"] = fmt.Sprintf("%v", operation.Quantity)
	mapParams["price"] = fmt.Sprintf("%v", operation.Rate)

	if operation.TradeType == exchange.TRADE_LIMIT {
		mapParams["category"] = "1" // limit
	} else if operation.TradeType == exchange.TRADE_MARKET {
		mapParams["category"] = "3" // quick trade
	}

	if operation.OrderDirection == exchange.Buy {
		mapParams["type"] = "1" // buy: 1, sell: 2
	} else if operation.OrderDirection == exchange.Sell {
		mapParams["type"] = "2" // buy: 1, sell: 2
	}

	jsonResponse := &JsonResponse{}
	placeOrder := PlaceOrder{}
	strRequest := "/member/addOrder"

	jsonPlaceReturn := e.ApiKeyRequest("POST", strRequest, mapParams)
	if err := json.Unmarshal([]byte(jsonPlaceReturn), &jsonResponse); err != nil {
		return fmt.Errorf("%s PlaceOrder Json Unmarshal Err: %v %v", e.GetName(), err, jsonPlaceReturn)
	} else if jsonResponse.Code != 0 {
		return fmt.Errorf("%s PlaceOrder Failed: %v", e.GetName(), jsonPlaceReturn)
	}
	if err := json.Unmarshal(jsonResponse.Data, &placeOrder); err != nil {
		return fmt.Errorf("%s PlaceOrder Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
	}

	order := &exchange.Order{
		Pair:         operation.Pair,
		OrderID:      placeOrder.OrderID,
		Rate:         operation.Rate,
		Quantity:     operation.Quantity,
		Direction:    operation.OrderDirection,
		Status:       exchange.New,
		JsonResponse: jsonPlaceReturn,
	}

	operation.Order = order

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
