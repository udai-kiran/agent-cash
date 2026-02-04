package dto

import "time"

// SplitResponse represents a split in API responses
type SplitResponse struct {
	GUID           string         `json:"guid"`
	TxGUID         string         `json:"tx_guid"`
	AccountGUID    string         `json:"account_guid"`
	Memo           *string        `json:"memo,omitempty"`
	Action         *string        `json:"action,omitempty"`
	ReconcileState string         `json:"reconcile_state"`
	ValueNum       int64          `json:"value_num"`
	ValueDenom     int64          `json:"value_denom"`
	QuantityNum    int64          `json:"quantity_num"`
	QuantityDenom  int64          `json:"quantity_denom"`
	Value          string         `json:"value"`
	Quantity       string         `json:"quantity"`
	Account        *AccountSummary `json:"account,omitempty"`
}

// AccountSummary represents basic account info in transaction responses
type AccountSummary struct {
	GUID string `json:"guid"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// TransactionResponse represents a transaction in API responses
type TransactionResponse struct {
	GUID             string          `json:"guid"`
	CurrencyGUID     string          `json:"currency_guid"`
	CurrencyMnemonic string          `json:"currency_mnemonic,omitempty"`
	Num              *string         `json:"num,omitempty"`
	PostDate         time.Time       `json:"post_date"`
	EnterDate        time.Time       `json:"enter_date"`
	Description      *string         `json:"description,omitempty"`
	Splits           []SplitResponse `json:"splits"`
}

// TransactionListResponse represents a paginated list of transactions
type TransactionListResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	Total        int64                 `json:"total"`
	Limit        int                   `json:"limit"`
	Offset       int                   `json:"offset"`
}
