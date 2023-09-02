package dsl

import (
	"slices"
	"strings"

	daprPubSub "github.com/lburgazzoli/camel-go/pkg/components/dapr/pubsub"
)

type Capability int

const (
	Capability_DAPR Capability = iota
	Capability_KNATIVE
)

func NewMetadata() *Metadata {
	return &Metadata{
		FromURIs: make([]string, 0),
		ToURIs:   make([]string, 0),
	}
}

type Metadata struct {
	// All starting URIs of defined routes
	FromURIs []string
	// All end URIs of defined routes
	ToURIs []string
}

func (m *Metadata) Capabilities() []Capability {
	capabilities := make([]Capability, 0)

	for i := range m.FromURIs {

		if strings.HasPrefix(m.FromURIs[i], daprPubSub.Scheme+":") {
			if !slices.Contains(capabilities, Capability_DAPR) {
				capabilities = append(capabilities, Capability_DAPR)
			}
		}
	}

	for i := range m.ToURIs {
		if strings.HasPrefix(m.ToURIs[i], daprPubSub.Scheme+":") {
			if !slices.Contains(capabilities, Capability_DAPR) {
				capabilities = append(capabilities, Capability_DAPR)
			}
		}
	}

	return capabilities
}
