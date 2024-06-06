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
	"crypto/sha256"
	"encoding/base64"
	"sort"

	"github.com/lburgazzoli/camel-go/pkg/controller/reconciler"

	camelApi "github.com/lburgazzoli/camel-go/api/camel/v2alpha1"
	"go.uber.org/multierr"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/log"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)
	l.Info("Reconciling", "resource", req.NamespacedName.String())

	rr := ReconciliationRequest{
		Client: r.Client(),
		NamespacedName: types.NamespacedName{
			Name:      req.Name,
			Namespace: req.Namespace,
		},
		ClusterType: r.ClusterType,
		Reconciler:  r,
		Resource:    &camelApi.Integration{},
	}

	err := r.Client().Get(ctx, req.NamespacedName, rr.Resource)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// no CR found anymore, maybe deleted
			return ctrl.Result{}, nil
		}
	}

	sum := sha256.New()
	for i := range rr.Resource.Spec.Flows {
		sum.Write(rr.Resource.Spec.Flows[i].RawMessage)
	}

	rr.Checksum = base64.RawURLEncoding.EncodeToString(sum.Sum(nil))

	if rr.Resource.ObjectMeta.DeletionTimestamp.IsZero() {
		err := reconciler.AddFinalizer(ctx, r.Client(), rr.Resource, FinalizerName)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else {
		for i := len(r.actions) - 1; i >= 0; i-- {
			if err := r.actions[i].Cleanup(ctx, &rr); err != nil {
				return ctrl.Result{}, err
			}
		}

		err := reconciler.RemoveFinalizer(ctx, r.Client(), rr.Resource, FinalizerName)
		if err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	//
	// Reconcile
	//

	reconcileCondition := metav1.Condition{
		Type:               ConditionReconciled,
		Status:             metav1.ConditionTrue,
		Reason:             "Reconciled",
		Message:            "Reconciled",
		ObservedGeneration: rr.Resource.Generation,
	}

	var allErrors error

	for i := range r.actions {
		if err := r.actions[i].Run(ctx, &rr); err != nil {
			allErrors = multierr.Append(allErrors, err)
		}
	}

	if allErrors != nil {
		reconcileCondition.Status = metav1.ConditionFalse
		reconcileCondition.Reason = FailureReason
		reconcileCondition.Message = FailureReason

		rr.Resource.Status.Phase = PhaseError
	} else {
		rr.Resource.Status.ObservedGeneration = rr.Resource.Generation
		rr.Resource.Status.Phase = PhaseReady
	}

	meta.SetStatusCondition(&rr.Resource.Status.Conditions, reconcileCondition)

	sort.SliceStable(rr.Resource.Status.Conditions, func(i, j int) bool {
		return rr.Resource.Status.Conditions[i].Type < rr.Resource.Status.Conditions[j].Type
	})

	//
	// Update status
	//

	err = r.Client().Status().Update(ctx, rr.Resource)
	if err != nil && k8serrors.IsConflict(err) {
		l.Info(err.Error())
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		allErrors = multierr.Append(allErrors, err)
	}

	return ctrl.Result{}, allErrors
}
