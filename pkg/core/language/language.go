package language

type Language struct {
	Jq       *Jq       `yaml:"jq,omitempty"`
	Mustache *Mustache `yaml:"mustache,omitempty"`
	Wasm     *Wasm     `yaml:"wasm,omitempty"`
}
