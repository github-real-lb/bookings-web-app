package models

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]any
	CSRFToken string // Security Token to prevent Cross Site Request Forgery (CSRF)
	Flash     string // Flash message
	Warning   string // Warning message
	Error     string // Error message
}
