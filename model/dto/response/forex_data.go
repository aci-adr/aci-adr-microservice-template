package response

import "time"

type ForexDataResponse struct {
	Id string `json:"id"`

	TenantId int `json:"tenantId"`

	BankId int `json:"bankId"`

	BaseCurrency string `json:"baseCurrency"`

	TargetCurrency string `json:"targetCurrency"`

	Tier string `json:"tier"`

	DirectIndirectFlag string `json:"directIndirectFlag"`

	Multiplier float64 `json:"multiplier"`

	BuyRate float64 `json:"buyRate"`

	SellRate float64 `json:"sellRate"`

	TolerancePercentage int `json:"tolerancePercentage"`

	EffectiveDate *time.Time `json:"effectiveDate"`

	ExpirationDate *time.Time `json:"expirationDate"`

	ContractRequirementThreshold string `json:"contractRequirementThreshold"`
}
