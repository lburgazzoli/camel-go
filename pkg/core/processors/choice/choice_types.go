package choice

import camel "github.com/lburgazzoli/camel-go/pkg/api"

// branchDone is a marker type to help when/otherwise to break out
type branchDone struct {
	M camel.Message
}
