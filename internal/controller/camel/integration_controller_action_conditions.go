package camel

import (
	"context"
	"fmt"
	gc "github.com/lburgazzoli/camel-go/pkg/controller/gc"
	"github.com/lburgazzoli/camel-go/pkg/util/conditions"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"

	"github.com/go-logr/logr"
	"github.com/lburgazzoli/camel-go/pkg/controller/client"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
)

func NewConditionsAction() Action {
	return &ConditionsAction{
		l:             ctrl.Log.WithName("action").WithName("conditions"),
		subscriptions: make(map[string]struct{}),
		gc:            gc.New(),
	}
}

type ConditionsAction struct {
	gc            *gc.GC
	l             logr.Logger
	subscriptions map[string]struct{}
}

func (a *ConditionsAction) Configure(_ context.Context, _ *client.Client, b *builder.Builder) (*builder.Builder, error) {
	return b, nil
}

func (a *ConditionsAction) Run(ctx context.Context, rc *ReconciliationRequest) error {
	crs, err := WithCurrentGenerationSelector(rc)
	if err != nil {
		return errors.Wrap(err, "cannot compute current resources selector")
	}

	deployments, err := rc.Client.AppsV1().Deployments(rc.Resource.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: crs.String(),
	})

	if err != nil {
		return errors.Wrap(err, "cannot list deployments")
	}

	ready := 0
	for i := range deployments.Items {
		if conditions.ConditionStatus(deployments.Items[i], appsv1.DeploymentAvailable) == corev1.ConditionTrue {
			ready++
		}
	}

	var readyCondition metav1.Condition

	if len(deployments.Items) > 0 {
		if ready == len(deployments.Items) {
			readyCondition = metav1.Condition{
				Type:               ConditionReady,
				Status:             metav1.ConditionTrue,
				Reason:             "Ready",
				Message:            fmt.Sprintf("%d/%d deployments ready", ready, len(deployments.Items)),
				ObservedGeneration: rc.Resource.Generation,
			}
		} else {
			readyCondition = metav1.Condition{
				Type:               ConditionReady,
				Status:             metav1.ConditionFalse,
				Reason:             "InProgress",
				Message:            fmt.Sprintf("%d/%d deployments ready", ready, len(deployments.Items)),
				ObservedGeneration: rc.Resource.Generation,
			}
		}
	} else {
		readyCondition = metav1.Condition{
			Type:               ConditionReady,
			Status:             metav1.ConditionFalse,
			Reason:             "InProgress",
			Message:            "no deployments",
			ObservedGeneration: rc.Resource.Generation,
		}
	}

	meta.SetStatusCondition(&rc.Resource.Status.Conditions, readyCondition)

	return nil
}

func (a *ConditionsAction) Cleanup(_ context.Context, _ *ReconciliationRequest) error {
	return nil
}
