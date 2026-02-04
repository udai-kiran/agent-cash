package dto

// CommodityResponse represents a commodity in API responses
type CommodityResponse struct {
	GUID      string `json:"guid"`
	Namespace string `json:"namespace"`
	Mnemonic  string `json:"mnemonic"`
	Fullname  string `json:"fullname"`
	Fraction  int    `json:"fraction"`
}
