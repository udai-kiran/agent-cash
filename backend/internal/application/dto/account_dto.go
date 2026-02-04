package dto

// AccountResponse represents an account in API responses
type AccountResponse struct {
	GUID              string            `json:"guid"`
	Name              string            `json:"name"`
	Type              string            `json:"type"`
	Code              *string           `json:"code,omitempty"`
	Description       *string           `json:"description,omitempty"`
	Hidden            bool              `json:"hidden"`
	Placeholder       bool              `json:"placeholder"`
	ParentGUID        *string           `json:"parent_guid,omitempty"`
	Balance           string            `json:"balance"`
	BalanceNum        int64             `json:"balance_num"`
	BalanceDenom      int64             `json:"balance_denom"`
	CommodityMnemonic string            `json:"commodity_mnemonic,omitempty"`
	Children          []AccountResponse `json:"children,omitempty"`
}

// AccountBalanceResponse represents an account balance
type AccountBalanceResponse struct {
	GUID         string `json:"guid"`
	Name         string `json:"name"`
	Balance      string `json:"balance"`
	BalanceNum   int64  `json:"balance_num"`
	BalanceDenom int64  `json:"balance_denom"`
}
