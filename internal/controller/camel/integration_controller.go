/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package camel

import (
	"context"

	camlApi "github.com/lburgazzoli/camel-go/api/camel/v2alpha1"

	"k8s.io/client-go/tools/record"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/go-logr/logr"
	"github.com/lburgazzoli/camel-go/pkg/controller"
	"github.com/lburgazzoli/camel-go/pkg/controller/client"

	ctrlRt "sigs.k8s.io/controller-runtime"
	ctrlCli "sigs.k8s.io/controller-runtime/pkg/client"
	ctrl "sigs.k8s.io/controller-runtime/pkg/controller"
)

func NewReconciler(ctx context.Context, manager ctrlRt.Manager, _ Options) (*Reconciler, error) {
	c, err := client.NewClient(manager.GetConfig(), manager.GetScheme(), manager.GetClient())
	if err != nil {
		return nil, err
	}

	rec := Reconciler{}
	rec.l = ctrlRt.Log.WithName("controller")
	rec.Client = c
	rec.Scheme = manager.GetScheme()
	rec.ClusterType = controller.ClusterTypeVanilla
	rec.manager = manager
	rec.recorder = manager.GetEventRecorderFor(FieldManager)

	isOpenshift, err := c.IsOpenShift()
	if err != nil {
		return nil, err
	}
	if isOpenshift {
		rec.ClusterType = controller.ClusterTypeOpenShift
	}

	rec.actions = append(rec.actions, NewDeployAction())
	rec.actions = append(rec.actions, NewConditionsAction())

	err = rec.init(ctx)
	if err != nil {
		return nil, err
	}

	return &rec, nil
}

// +kubebuilder:rbac:groups=camel.apache.org,resources=integrations,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=camel.apache.org,resources=integrations/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=camel.apache.org,resources=integrations/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=events,verbs=*
// +kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=*
// +kubebuilder:rbac:groups="",resources=secrets,verbs=*
// +kubebuilder:rbac:groups="",resources=services,verbs=*
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=*
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=*

type Reconciler struct {
	*client.Client

	Scheme      *runtime.Scheme
	ClusterType controller.ClusterType
	actions     []Action
	l           logr.Logger
	manager     ctrlRt.Manager
	controller  ctrl.Controller
	recorder    record.EventRecorder
}

func (r *Reconciler) init(ctx context.Context) error {
	c := ctrlRt.NewControllerManagedBy(r.manager)

	c = c.For(&camlApi.Integration{}, builder.WithPredicates(
		predicate.GenerationChangedPredicate{},
	))

	for i := range r.actions {
		b, err := r.actions[i].Configure(ctx, r.Client, c)
		if err != nil {
			return err
		}

		c = b
	}

	ct, err := c.Build(r)
	if err != nil {
		return err
	}

	r.controller = ct

	return nil
}

func (r *Reconciler) Watch(obj ctrlCli.Object, eh handler.EventHandler, predicates ...predicate.Predicate) error {
	return r.controller.Watch(
		source.Kind(
			r.manager.GetCache(),
			obj,
			eh,
			predicates...))
}

func (r *Reconciler) EnqueueRequestForOwner(owner ctrlCli.Object, opts ...handler.OwnerOption) handler.EventHandler {
	return handler.EnqueueRequestForOwner(
		r.manager.GetScheme(),
		r.manager.GetRESTMapper(),
		owner,
		opts...,
	)
}

func (r *Reconciler) EnqueueRequestsFromMapFunc(fn func(context.Context, ctrlCli.Object) []reconcile.Request) handler.EventHandler {
	return handler.EnqueueRequestsFromMapFunc(fn)
}

func (r *Reconciler) Event(object runtime.Object, eventType string, reason string, message string) {
	r.recorder.Event(
		object,
		eventType,
		reason,
		message,
	)
}
