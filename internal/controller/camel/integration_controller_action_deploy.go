package camel

import (
	"context"
	"fmt"

	"github.com/lburgazzoli/camel-go/pkg/components/dapr"
	"github.com/lburgazzoli/camel-go/pkg/health"

	"slices"

	"github.com/go-logr/logr"
	"github.com/lburgazzoli/camel-go/pkg/controller/client"
	"github.com/lburgazzoli/camel-go/pkg/controller/gc"
	"github.com/lburgazzoli/camel-go/pkg/util/dsl"
	"github.com/lburgazzoli/camel-go/pkg/util/resources"
	"github.com/lburgazzoli/camel-go/pkg/util/resources/apply"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1ac "k8s.io/client-go/applyconfigurations/apps/v1"
	corev1ac "k8s.io/client-go/applyconfigurations/core/v1"
	metav1ac "k8s.io/client-go/applyconfigurations/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
)

func NewDeployAction() Action {
	return &DeployAction{
		l:             ctrl.Log.WithName("action").WithName("deploy"),
		subscriptions: make(map[string]struct{}),
		gc:            gc.New(),
	}
}

type DeployAction struct {
	l             logr.Logger
	subscriptions map[string]struct{}
	gc            *gc.GC
}

func (a *DeployAction) Configure(_ context.Context, _ *client.Client, b *builder.Builder) (*builder.Builder, error) {
	return b, nil
}

func (a *DeployAction) Run(ctx context.Context, rc *ReconciliationRequest) error {
	deploymentCondition := metav1.Condition{
		Type:               "Deployment",
		Status:             metav1.ConditionTrue,
		Reason:             "Deployed",
		Message:            "Deployed",
		ObservedGeneration: rc.Resource.Generation,
	}

	err := a.deploy(ctx, rc)
	if err != nil {
		deploymentCondition.Status = metav1.ConditionFalse
		deploymentCondition.Reason = "Failure"
		deploymentCondition.Message = err.Error()
	}

	meta.SetStatusCondition(&rc.Resource.Status.Conditions, deploymentCondition)

	return err
}

func (a *DeployAction) Cleanup(_ context.Context, _ *ReconciliationRequest) error {
	return nil
}

func (a *DeployAction) deploy(ctx context.Context, rc *ReconciliationRequest) error {

	//
	// ConfigMap
	//

	cm, err := a.configmap(ctx, rc)
	if err != nil {
		return err
	}

	_, err = rc.Client.CoreV1().ConfigMaps(rc.Resource.Namespace).Apply(
		ctx,
		cm,
		metav1.ApplyOptions{
			FieldManager: FieldManager,
			Force:        true,
		},
	)

	if err != nil {
		return err
	}

	//
	// Deployment
	//

	d, err := a.deployment(ctx, rc)
	if err != nil {
		return err
	}

	_, err = rc.Client.AppsV1().Deployments(rc.Resource.Namespace).Apply(
		ctx,
		d,
		metav1.ApplyOptions{
			FieldManager: FieldManager,
			Force:        true,
		},
	)

	if err != nil {
		return err
	}

	return nil
}

const (
	RuntimeContainerName  = "integration"
	RuntimeContainerImage = "quay.io/lburgazzoli/camel-go:latest"
	RuntimeRoutesPath     = "/etc/camel/sources.d/routes.yaml"
	RuntimeRoutesSubPath  = "routes.yaml"
)

func (a *DeployAction) configmap(_ context.Context, rc *ReconciliationRequest) (*corev1ac.ConfigMapApplyConfiguration, error) {
	labels := LabelsForIntegration(rc)
	annotations := AnnotationForIntegration(rc)

	data, err := dsl.ToYamlDSL(rc.Resource.Spec.Flows)
	if err != nil {
		return nil, err
	}

	resource := corev1ac.ConfigMap(rc.Resource.Name, rc.Resource.Namespace).
		WithOwnerReferences(apply.WithOwnerReference(rc.Resource)).
		WithLabels(labels).
		WithAnnotations(annotations).
		WithData(map[string]string{
			RuntimeRoutesSubPath: string(data),
		})

	return resource, nil
}

func (a *DeployAction) deployment(_ context.Context, rc *ReconciliationRequest) (*appsv1ac.DeploymentApplyConfiguration, error) {
	labels := LabelsForIntegration(rc)
	lsec := LabelsForIntegrationSelector(rc)
	annotations := AnnotationForIntegration(rc)
	podannotations := AnnotationForIntegration(rc)

	m, err := dsl.NewInspector().Extract(rc.Resource.Spec.Flows)
	if err != nil {
		return nil, err
	}

	envs := make([]*corev1ac.EnvVarApplyConfiguration, 0)
	envs = append(envs, apply.WithEnvFromField(resources.EnvVarNamespace, "metadata.namespace"))
	envs = append(envs, apply.WithEnvFromField(resources.EnvVarPodName, "metadata.name"))
	envs = append(envs, apply.WithEnv(resources.EnvVarIntegrationChecksum, rc.Checksum))

	ports := make([]*corev1ac.ContainerPortApplyConfiguration, 0)
	ports = append(ports, apply.WithPort(HttpPortName, HttpPort))
	ports = append(ports, apply.WithPort(health.DefaultPortName, health.DefaultPort))

	if slices.Contains(m.Capabilities(), dsl.Capability_DAPR) {
		envs = append(envs, apply.WithEnv(dapr.EnvVarAddress, fmt.Sprintf(":%d", dapr.DefaultPort)))
		ports = append(ports, apply.WithPort(dapr.DefaultPortName, dapr.DefaultPort))

		podannotations[dapr.AnnotationAppID] = rc.Resource.Namespace + "-" + rc.Resource.Name
		podannotations[dapr.AnnotationAppPort] = fmt.Sprintf("%d", dapr.DefaultPort)
		podannotations[dapr.AnnotationAppProtocol] = dapr.DefaultProtocol
	}

	resource := appsv1ac.Deployment(rc.Resource.Name, rc.Resource.Namespace).
		WithOwnerReferences(apply.WithOwnerReference(rc.Resource)).
		WithLabels(labels).
		WithAnnotations(annotations).
		WithSpec(appsv1ac.DeploymentSpec().
			WithReplicas(1).
			WithSelector(metav1ac.LabelSelector().WithMatchLabels(lsec)).
			WithTemplate(corev1ac.PodTemplateSpec().
				WithLabels(labels).
				WithAnnotations(podannotations).
				WithSpec(corev1ac.PodSpec().
					// WithServiceAccountName(rc.Resource.Name).
					WithVolumes(corev1ac.Volume().
						WithName("routes").
						WithConfigMap(
							corev1ac.ConfigMapVolumeSource().
								WithName(rc.Resource.Name).
								WithItems(corev1ac.KeyToPath().
									WithKey(RuntimeRoutesSubPath).
									WithPath(RuntimeRoutesSubPath),
								),
						),
					).
					WithContainers(corev1ac.Container().
						WithImage(RuntimeContainerImage).
						WithImagePullPolicy(corev1.PullAlways).
						WithName(RuntimeContainerName).
						WithArgs(
							"run",
							"--health-check-enabled", "true",
							"--health-check-address", health.DefaultAddress,
							"--dev",
							"--route", RuntimeRoutesPath,
						).
						WithEnv(envs...).
						WithPorts(ports...).
						WithReadinessProbe(apply.WithHTTPProbe(ReadinessProbePath, health.DefaultPort)).
						WithLivenessProbe(apply.WithHTTPProbe(LivenessProbePath, health.DefaultPort)).
						WithResources(corev1ac.ResourceRequirements().
							WithRequests(corev1.ResourceList{
								corev1.ResourceMemory: DefaultMemory,
								corev1.ResourceCPU:    DefaultCPU,
							}),
						).
						WithVolumeMounts(corev1ac.VolumeMount().
							WithName("routes").
							WithMountPath(RuntimeRoutesPath).
							WithSubPath(RuntimeRoutesSubPath),
						),
					),
				),
			),
		)

	return resource, nil
}
