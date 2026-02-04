package entity

// Commodity represents a GnuCash commodity (currency, stock, etc.)
type Commodity struct {
	GUID      string
	Namespace string
	Mnemonic  string
	Fullname  string
	Fraction  int
}
