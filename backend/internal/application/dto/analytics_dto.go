package dto

import "time"

// IncomeExpenseData represents income vs expense data for a period
type IncomeExpenseData struct {
	Period  string  `json:"period"`
	Income  string  `json:"income"`
	Expense string  `json:"expense"`
	Net     string  `json:"net"`
}

// IncomeExpenseResponse represents the income/expense analytics response
type IncomeExpenseResponse struct {
	Data             []IncomeExpenseData `json:"data"`
	TotalIncome      string              `json:"total_income"`
	TotalExpense     string              `json:"total_expense"`
	NetTotal         string              `json:"net_total"`
	CurrencyMnemonic string              `json:"currency_mnemonic,omitempty"`
}

// CategoryBreakdownItem represents spending by category
type CategoryBreakdownItem struct {
	Category string `json:"category"`
	Amount   string `json:"amount"`
	Count    int    `json:"count"`
}

// CategoryBreakdownResponse represents category breakdown analytics
type CategoryBreakdownResponse struct {
	Income           []CategoryBreakdownItem `json:"income"`
	Expense          []CategoryBreakdownItem `json:"expense"`
	CurrencyMnemonic string                  `json:"currency_mnemonic,omitempty"`
}

// TrendDataPoint represents a data point in a trend
type TrendDataPoint struct {
	Date   time.Time `json:"date"`
	Amount string    `json:"amount"`
}

// TrendResponse represents trend analytics
type TrendResponse struct {
	AccountName string           `json:"account_name"`
	Data        []TrendDataPoint `json:"data"`
}

// NetWorthItem represents an asset or liability balance
type NetWorthItem struct {
	AccountName string `json:"account_name"`
	AccountType string `json:"account_type"`
	Balance     string `json:"balance"`
}

// NetWorthResponse represents net worth analytics
type NetWorthResponse struct {
	Assets           []NetWorthItem `json:"assets"`
	Liabilities      []NetWorthItem `json:"liabilities"`
	TotalAssets      string         `json:"total_assets"`
	TotalLiabilities string         `json:"total_liabilities"`
	NetWorth         string         `json:"net_worth"`
	CurrencyMnemonic string         `json:"currency_mnemonic,omitempty"`
}
