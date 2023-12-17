package request

type NatConvertRequest struct {
	TenantID       int     `json:"tenantId"`
	BankID         int     `json:"bankId"`
	BaseCurrency   string  `json:"baseCurrency"`
	TargetCurrency string  `json:"targetCurrency"`
	Tier           string  `json:"tier"`
	Amount         float64 `json:"amount"`
	InitiatedOn    int64   `json:"initiatedOn"`
}
