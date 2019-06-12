package bitbay

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

import (
	"encoding/json"
)

type JsonResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

type PairsData struct {
	Status string               `json:"status"`
	Pairs  map[string]*PairInfo `json:"items"`
}

type PairInfo struct {
	Market struct {
		Code  string `json:"code"`
		First struct {
			Currency string `json:"currency"`
			MinOffer string `json:"minOffer"`
			Scale    int    `json:"scale"`
		} `json:"first"`
		Second struct {
			Currency string `json:"currency"`
			MinOffer string `json:"minOffer"`
			Scale    int    `json:"scale"`
		} `json:"second"`
	} `json:"market"`
	Time         string `json:"time"`
	HighestBid   string `json:"highestBid"`
	LowestAsk    string `json:"lowestAsk"`
	Rate         string `json:"rate"`
	PreviousRate string `json:"previousRate"`
}

type OrderBook struct {
	Status string `json:"status"`
	Sell   []struct {
		Ra string `json:"ra"`
		Ca string `json:"ca"`
		Sa string `json:"sa"`
		Pa string `json:"pa"`
		Co int    `json:"co"`
	} `json:"sell"`
	Buy []struct {
		Ra string `json:"ra"`
		Ca string `json:"ca"`
		Sa string `json:"sa"`
		Pa string `json:"pa"`
		Co int    `json:"co"`
	} `json:"buy"`
	Timestamp string `json:"timestamp"`
}

type AccountBalances struct {
	Status   string `json:"status"`
	Balances []struct {
		ID             string  `json:"id"`
		UserID         string  `json:"userId"`
		AvailableFunds float64 `json:"availableFunds"`
		TotalFunds     float64 `json:"totalFunds"`
		LockedFunds    int     `json:"lockedFunds"`
		Currency       string  `json:"currency"`
		Type           string  `json:"type"`
		Name           string  `json:"name"`
		BalanceEngine  string  `json:"balanceEngine"`
	} `json:"balances"`
	Errors interface{} `json:"errors"`
}

type PlaceOrder struct {
	Status       string `json:"status"`
	Completed    bool   `json:"completed"`
	OfferID      string `json:"offerId"`
	Transactions []struct {
		Amount string `json:"amount"`
		Rate   string `json:"rate"`
	} `json:"transactions"`
}

type OrderStatus struct {
	Status string `json:"status"`
	Items  []struct {
		Market          string `json:"market"`
		OfferType       string `json:"offerType"`
		ID              string `json:"id"`
		CurrentAmount   string `json:"currentAmount"`
		LockedAmount    string `json:"lockedAmount"`
		Rate            string `json:"rate"`
		StartAmount     string `json:"startAmount"`
		Time            string `json:"time"`
		PostOnly        bool   `json:"postOnly"`
		Mode            string `json:"mode"`
		ReceivedAmount  string `json:"receivedAmount"`
		FirstBalanceID  string `json:"firstBalanceId"`
		SecondBalanceID string `json:"secondBalanceId"`
	} `json:"items"`
}

type CancelOrder struct {
	Status string        `json:"status"`
	Errors []interface{} `json:"errors"`
}
