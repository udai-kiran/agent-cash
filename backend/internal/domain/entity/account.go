package entity

import (
	"github.com/shopspring/decimal"
)

// AccountType represents the type of account in GnuCash
type AccountType string

const (
	AccountTypeRoot       AccountType = "ROOT"
	AccountTypeBank       AccountType = "BANK"
	AccountTypeCash       AccountType = "CASH"
	AccountTypeCredit     AccountType = "CREDIT"
	AccountTypeAsset      AccountType = "ASSET"
	AccountTypeLiability  AccountType = "LIABILITY"
	AccountTypeStock      AccountType = "STOCK"
	AccountTypeMutual     AccountType = "MUTUAL"
	AccountTypeCurrency   AccountType = "CURRENCY"
	AccountTypeIncome     AccountType = "INCOME"
	AccountTypeExpense    AccountType = "EXPENSE"
	AccountTypeEquity     AccountType = "EQUITY"
	AccountTypeReceivable AccountType = "RECEIVABLE"
	AccountTypePayable    AccountType = "PAYABLE"
)

// Account represents a GnuCash account
type Account struct {
	GUID            string
	Name            string
	AccountType     AccountType
	CommodityGUID   *string
	CommoditySCU    int
	ParentGUID      *string
	Code            *string
	Description     *string
	Hidden          bool
	Placeholder     bool
	Children           []*Account
	Balance            decimal.Decimal
	BalanceNum         int64
	BalanceDenom       int64
	CommodityMnemonic  string
}

// IsDebitAccount returns true if the account increases with debits
func (a *Account) IsDebitAccount() bool {
	switch a.AccountType {
	case AccountTypeAsset, AccountTypeBank, AccountTypeCash,
		AccountTypeStock, AccountTypeMutual, AccountTypeReceivable, AccountTypeExpense:
		return true
	default:
		return false
	}
}

// IsCreditAccount returns true if the account increases with credits
func (a *Account) IsCreditAccount() bool {
	return !a.IsDebitAccount()
}
