package oksim

import (
	"encoding/json"
	"fmt"

	"github.com/bitontop/gored/exchange"
)

func (e *Oksim) DoAccountOperation(operation *exchange.AccountOperation) error {
	switch operation.Type {
	case exchange.Transfer:
		return e.transfer(operation)
	}
	return fmt.Errorf("%s Operation type invalid: %s %v", operation.Ex, operation.Wallet, operation.Type)
}

func (e *Oksim) transfer(operation *exchange.AccountOperation) error {
	if e.API_KEY == "" || e.API_SECRET == "" || e.Passphrase == "" {
		return fmt.Errorf("%s API Key, Secret Key or Passphrase are nil", e.GetName())
	}

	jsonResponse := &JsonResponse{}
	strRequest := "/api/v5/asset/transfer"

	mapParams := make(map[string]interface{})
	mapParams["ccy"] = e.GetSymbolByCoin(operation.Coin)
	mapParams["amt"] = operation.TransferAmount
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

	if err := json.Unmarshal([]byte(jsonTransferReturn), &jsonResponse); err != nil {
		operation.Error = fmt.Errorf("%s Transfer Json Unmarshal Err: %v %v", e.GetName(), err, jsonTransferReturn)
		return operation.Error
	} else if jsonResponse.Code != "0" {
		operation.Error = fmt.Errorf("%s Transfer Err: Code: %v Msg: %v", e.GetName(), jsonResponse.Code, jsonResponse.Msg)
		return operation.Error
	}

	transfer := []*Transfer{}
	if err := json.Unmarshal(jsonResponse.Data, &transfer); err != nil {
		operation.Error = fmt.Errorf("%s Transfer Data Unmarshal Err: %v %s", e.GetName(), err, jsonResponse.Data)
		return operation.Error
	} else if len(transfer) == 0 {
		operation.Error = fmt.Errorf("%s Transfer Failed: %v", e.GetName(), jsonTransferReturn)
		return operation.Error
	}

	return nil
}
