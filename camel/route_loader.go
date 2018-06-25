package camel

// Step --
type Step struct {
	URI string `yaml:"uri"`
}

// Flow --
type Flow struct {
	Steps []Step `yaml:"steps"`
}
