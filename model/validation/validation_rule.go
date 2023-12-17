package validation

type ValidationRule struct {
	ParamName string
	Required  bool
	// Add a ParamType field (e.g., "int" or "string")
	ParamType string
}
