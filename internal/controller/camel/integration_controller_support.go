package camel

import (
	"context"
	"github.com/lburgazzoli/camel-go/pkg/util/resources"
	"strconv"
	"strings"

	"github.com/lburgazzoli/camel-go/pkg/controller/predicates"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	ctrlCli "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

//nolint:unused
func gcSelector(rc *ReconciliationRequest) (labels.Selector, error) {

	namespace, err := labels.NewRequirement(
		IntegrationNamespace,
		selection.Equals,
		[]string{rc.Resource.Namespace})

	if err != nil {
		return nil, errors.Wrap(err, "cannot determine release namespace requirement")
	}

	name, err := labels.NewRequirement(
		IntegrationName,
		selection.Equals,
		[]string{rc.Resource.Name})

	if err != nil {
		return nil, errors.Wrap(err, "cannot determine release name requirement")
	}

	generation, err := labels.NewRequirement(
		IntegrationGeneration,
		selection.LessThan,
		[]string{strconv.FormatInt(rc.Resource.Generation, 10)})

	if err != nil {
		return nil, errors.Wrap(err, "cannot determine generation requirement")
	}

	selector := labels.NewSelector().
		Add(*namespace).
		Add(*name).
		Add(*generation)

	return selector, nil
}

//nolint:unused
func labelsToRequest(_ context.Context, object ctrlCli.Object) []reconcile.Request {
	allLabels := object.GetLabels()
	if allLabels == nil {
		return nil
	}

	namespace := allLabels[IntegrationNamespace]
	if namespace == "" {
		return nil
	}
	name := allLabels[IntegrationName]
	if name == "" {
		return nil
	}

	return []reconcile.Request{{
		NamespacedName: types.NamespacedName{
			Namespace: namespace,
			Name:      name,
		},
	}}
}

//nolint:unused
func dependantWithLabels(watchUpdate bool, watchDelete bool, watchStatus bool) predicate.Predicate {
	return predicate.And(
		&predicates.HasLabel{
			Name: IntegrationNamespace,
		},
		&predicates.HasLabel{
			Name: IntegrationName,
		},
		&predicates.DependentPredicate{
			WatchUpdate: watchUpdate,
			WatchDelete: watchDelete,
			WatchStatus: watchStatus,
		},
	)
}

func WithIntegrationLabelsSelector() (labels.Selector, error) {

	hasReleaseNamespaceLabel, err := labels.NewRequirement(IntegrationNamespace, selection.Exists, []string{})
	if err != nil {
		return nil, err
	}
	hasReleaseNameLabel, err := labels.NewRequirement(IntegrationName, selection.Exists, []string{})
	if err != nil {
		return nil, err
	}

	selector := labels.NewSelector().
		Add(*hasReleaseNamespaceLabel).
		Add(*hasReleaseNameLabel)

	return selector, nil
}

func WithCurrentGenerationSelector(rc *ReconciliationRequest) (labels.Selector, error) {
	namespace, err := labels.NewRequirement(
		IntegrationNamespace,
		selection.Equals,
		[]string{rc.Resource.Namespace})

	if err != nil {
		return nil, errors.Wrap(err, "cannot determine release namespace requirement")
	}

	name, err := labels.NewRequirement(
		IntegrationName,
		selection.Equals,
		[]string{rc.Resource.Name})

	if err != nil {
		return nil, errors.Wrap(err, "cannot determine release name requirement")
	}

	generation, err := labels.NewRequirement(
		IntegrationGeneration,
		selection.Equals,
		[]string{strconv.FormatInt(rc.Resource.Generation, 10)})

	if err != nil {
		return nil, errors.Wrap(err, "cannot determine generation requirement")
	}

	selector := labels.NewSelector().
		Add(*namespace).
		Add(*name).
		Add(*generation)

	return selector, nil
}

func AnnotationForIntegration(rc *ReconciliationRequest) map[string]string {
	return map[string]string{
		IntegrationGeneration: strconv.FormatInt(rc.Resource.Generation, 10),
		IntegrationNamespace:  rc.Resource.Namespace,
		IntegrationName:       rc.Resource.Name,
		IntegrationChecksum:   rc.Checksum,
	}
}

func LabelsForIntegration(rc *ReconciliationRequest) map[string]string {
	return map[string]string{
		resources.KubernetesLabelAppName:      strings.ToLower(rc.Resource.Kind),
		resources.KubernetesLabelAppInstance:  rc.Resource.GetName(),
		resources.KubernetesLabelAppComponent: "runtime",
		resources.KubernetesLabelAppPartOf:    "camel",
		resources.KubernetesLabelAppManagedBy: FieldManager,
	}
}

func LabelsForIntegrationSelector(rc *ReconciliationRequest) map[string]string {
	return map[string]string{
		resources.KubernetesLabelAppName:     strings.ToLower(rc.Resource.Kind),
		resources.KubernetesLabelAppInstance: rc.Resource.GetName(),
	}
}
