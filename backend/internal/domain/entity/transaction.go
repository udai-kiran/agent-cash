package entity

import "time"

// Transaction represents a GnuCash transaction
type Transaction struct {
	GUID              string
	CurrencyGUID      string
	CurrencyMnemonic  string
	Num               *string
	PostDate          time.Time
	EnterDate         time.Time
	Description       *string
	Splits            []*Split
}
