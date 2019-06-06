package gemini

// Copyright (c) 2015-2019 Bitontop Technologies Inc.
// Distributed under the MIT software license, see the accompanying
// file COPYING or http://www.opensource.org/licenses/mit-license.php.

type ErrorResponse struct {
	// error return
	Result  string `json:"result"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type PairsFile []struct {
	Symbol        string  `json:"symbol"`
	Base          string  `json:"base"`
	Quote         string  `json:"quote"`
	MinOrderSize  float64 `json:"min_order_size"`
	MinOrderIncre float64 `json:"min_order_incre"`
	MinPriceIncre float64 `json:"min_price_incre"`
}

type OrderBook struct {
	Bids []struct {
		Price     string `json:"price"`
		Amount    string `json:"amount"`
		Timestamp string `json:"timestamp"`
	} `json:"bids"`
	Asks []struct {
		Price     string `json:"price"`
		Amount    string `json:"amount"`
		Timestamp string `json:"timestamp"`
	} `json:"asks"`
}

type AccountBalances []struct {
	Type                   string `json:"type"`
	Currency               string `json:"currency"`
	Amount                 string `json:"amount"`
	Available              string `json:"available"`
	AvailableForWithdrawal string `json:"availableForWithdrawal"`
}

type Withdrawal struct {
	Address      string `json:"address"`
	Amount       string `json:"amount"`
	TxHash       string `json:"txHash"`
	WithdrawalID string `json:"withdrawalId"`
	Message      string `json:"message"`
}

type PlaceOrder struct {
	OrderID           string        `json:"order_id"`
	ID                string        `json:"id"`
	Symbol            string        `json:"symbol"`
	Exchange          string        `json:"exchange"`
	AvgExecutionPrice string        `json:"avg_execution_price"`
	Side              string        `json:"side"`
	Type              string        `json:"type"`
	Timestamp         string        `json:"timestamp"`
	Timestampms       int64         `json:"timestampms"`
	IsLive            bool          `json:"is_live"`
	IsCancelled       bool          `json:"is_cancelled"`
	IsHidden          bool          `json:"is_hidden"`
	WasForced         bool          `json:"was_forced"`
	ExecutedAmount    string        `json:"executed_amount"`
	RemainingAmount   string        `json:"remaining_amount"`
	Options           []interface{} `json:"options"`
	Price             string        `json:"price"`
	OriginalAmount    string        `json:"original_amount"`
}
