package gnucash

import (
	"github.com/shopspring/decimal"
)

// RationalToDecimal converts GnuCash rational number (numerator/denominator) to decimal
func RationalToDecimal(numerator, denominator int64) decimal.Decimal {
	if denominator == 0 {
		return decimal.Zero
	}
	return decimal.NewFromInt(numerator).Div(decimal.NewFromInt(denominator))
}

// DecimalToRational converts a decimal to rational number representation
func DecimalToRational(d decimal.Decimal, denominator int64) (int64, int64) {
	numerator := d.Mul(decimal.NewFromInt(denominator)).IntPart()
	return numerator, denominator
}

// FormatAmount formats a rational amount as a string
func FormatAmount(numerator, denominator int64) string {
	d := RationalToDecimal(numerator, denominator)
	return d.StringFixed(2)
}

// NormalizeSign adjusts the sign for GnuCash's accounting conventions
// Assets and Expenses are positive when debited
// Liabilities, Equity, and Income are positive when credited
func NormalizeSign(value int64, isDebitAccount bool) int64 {
	if isDebitAccount {
		return value
	}
	return -value
}
