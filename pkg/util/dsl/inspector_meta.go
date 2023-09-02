package dsl

import (
	"slices"
	"strings"

	daprPubSub "github.com/lburgazzoli/camel-go/pkg/components/dapr/pubsub"
)

type Capability int

const (
	CapabilityDAPR Capability = iota
	CapabilityKNATIVE
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
			if !slices.Contains(capabilities, CapabilityDAPR) {
				capabilities = append(capabilities, CapabilityDAPR)
			}
		}
	}

	for i := range m.ToURIs {
		if strings.HasPrefix(m.ToURIs[i], daprPubSub.Scheme+":") {
			if !slices.Contains(capabilities, CapabilityKNATIVE) {
				capabilities = append(capabilities, CapabilityKNATIVE)
			}
		}
	}

	return capabilities
}
