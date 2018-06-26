package camel

// ==========================
//
// Model
//
// ==========================

// Step --
type Step map[string]interface{}

// StepHandler --
type StepHandler func(step Step, route *RouteDefinition) (*RouteDefinition, error)

// Flow --
type Flow struct {
	ID    string `yaml:"id"`
	Name  string `yaml:"name"`
	Steps []Step `yaml:"steps"`
}

// Integration --
type Integration struct {
	Flows []Flow `yaml:"flows"`
}

// TypedStep --
type TypedStep struct {
	Type string `yaml:"type"`
}
