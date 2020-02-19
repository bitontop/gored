package binance

// Contributor 2015-2020 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"fmt"
	"log"
	utils "github.com/bitontop/gored/utils"
	exchange "github.com/bitontop/gored/exchange"
)



/*************** PUBLIC  API ***************/
func (e *Binance) LoadPublicData(operation *exchange.PublicOperation) error {
	switch operation.Type {

	case exchange.TradeHistory:
		return e.doTradeHistory(operation)

	}
	return fmt.Errorf("LoadPublicData :: Operation type invalid: %+v", operation.Type)
}

func (e *Binance) doTradeHistory(operation *exchange.PublicOperation) error {


	
	get := &utils.HttpGet{
		URI:         fmt.Sprintf("https://api.binance.com/api/v3/trades?symbol=%s&limit=%d", 
		e.GetSymbolByPair(operation.Pair),
		5,//TRADE_HISTORY_MAX_LIMIT,
	),
		
	}

	err := utils.HttpGetRequest(get)

	if err != nil {
		log.Printf("%+v", err)
		return err

	} else {
		log.Printf("%+v  ERR:%+v", string(get.ResponseBody),err)
		// if err := json.Unmarshal(post.ResponseBody, &e.User); err != nil {
		// 	return err
		// }
	}

	return nil
}

