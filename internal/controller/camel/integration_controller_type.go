package camel

import (
	"context"

	camelApi "github.com/lburgazzoli/camel-go/api/camel/v2alpha1"
	"github.com/lburgazzoli/camel-go/pkg/controller"
	"github.com/lburgazzoli/camel-go/pkg/controller/client"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/builder"
)

const (
	IntegrationGeneration = "camel.apache.org/integration.generation"
	IntegrationName       = "camel.apache.org/integration.name"
	IntegrationNamespace  = "camel.apache.org/integration.namespace"

	FinalizerName = "camel.apache.org/finalizer"
	FieldManager  = "camel-control-plane"

	ConditionReconciled = "Reconcile"
	ConditionReady      = "Ready"
	PhaseError          = "Error"
	PhaseReady          = "Ready"
)

type Options struct {
}

type ReconciliationRequest struct {
	*client.Client
	types.NamespacedName

	Reconciler  *Reconciler
	ClusterType controller.ClusterType
	Resource    *camelApi.Integration
}

type Action interface {
	Configure(context.Context, *client.Client, *builder.Builder) (*builder.Builder, error)
	Run(context.Context, *ReconciliationRequest) error
	Cleanup(context.Context, *ReconciliationRequest) error
}
