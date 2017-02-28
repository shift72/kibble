package models

// ServiceConfig -
type ServiceConfig map[string]string

// FeatureToggles - store feature toggles
type FeatureToggles map[string]bool

// Site -
type Site struct {
	Config  ServiceConfig
	Toggles FeatureToggles
}
