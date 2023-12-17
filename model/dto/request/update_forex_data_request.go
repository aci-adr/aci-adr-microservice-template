package request

import "time"

type UpdateForexDataRequest struct {
	DirectIndirectFlag string `json:"directIndirectFlag,omitempty"`

	Multiplier int32 `json:"multiplier,omitempty"`

	BuyRate float64 `json:"buyRate" binding:"required"`

	SellRate float64 `json:"sellRate" binding:"required"`

	TolerancePercentage int32 `json:"tolerancePercentage,omitempty"`

	EffectiveDate *time.Time `json:"effectiveDate,omitempty"`

	ExpirationDate *time.Time `json:"expirationDate,omitempty"`

	ContractRequirementThreshold string `json:"contractRequirementThreshold,omitempty"`
}
