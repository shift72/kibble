package models

// Config -
type Config struct {
	DefaultLanguage string            `json:"defaultLanguage"`
	Languages       map[string]string `json:"languages"`
	Routes          []Route           `json:"routes"`
}
