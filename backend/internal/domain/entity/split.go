package entity

import "github.com/shopspring/decimal"

// Split represents a GnuCash split (part of a transaction)
type Split struct {
	GUID           string
	TxGUID         string
	AccountGUID    string
	Memo           *string
	Action         *string
	ReconcileState string
	ValueNum       int64
	ValueDenom     int64
	QuantityNum    int64
	QuantityDenom  int64
	Value          decimal.Decimal
	Quantity       decimal.Decimal
	Account        *Account
}
