package request

import "time"

type CreateForexDataRequest struct {
	TenantId int `json:"tenantId" binding:"required"`

	BankId int `json:"bankId" binding:"required"`

	BaseCurrency string `json:"baseCurrency" binding:"required"`

	TargetCurrency string `json:"targetCurrency" binding:"required"`

	Tier string `json:"tier" binding:"required"`

	DirectIndirectFlag string `json:"directIndirectFlag,omitempty"`

	Multiplier float64 `json:"multiplier,omitempty"`

	BuyRate float64 `json:"buyRate" binding:"required"`

	SellRate float64 `json:"sellRate" binding:"required"`

	TolerancePercentage int `json:"tolerancePercentage"`

	EffectiveDate *time.Time `json:"effectiveDate"`

	ExpirationDate *time.Time `json:"expirationDate"`

	ContractRequirementThreshold string `json:"contractRequirementThreshold,omitempty"`
}
