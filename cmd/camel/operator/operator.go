package operator

import (
	"flag"

	camelApi "github.com/lburgazzoli/camel-go/api/camel/v2alpha1"
	camelCtrl "github.com/lburgazzoli/camel-go/internal/controller/camel"
	"github.com/lburgazzoli/camel-go/pkg/controller"
	"github.com/lburgazzoli/camel-go/pkg/controller/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog/v2"
	rtcache "sigs.k8s.io/controller-runtime/pkg/cache"
	rtclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func init() {
	utilruntime.Must(camelApi.AddToScheme(controller.Scheme))
}

func NewOperatorCmd() *cobra.Command {
	co := controller.Options{
		MetricsAddr:                   ":8080",
		ProbeAddr:                     ":8081",
		PprofAddr:                     "",
		LeaderElectionID:              "9bsv0s25god0031unma0.camel.apache.org",
		EnableLeaderElection:          true,
		ReleaseLeaderElectionOnCancel: true,
		LeaderElectionNamespace:       "",
	}

	cmd := cobra.Command{
		Use:   "operator",
		Short: "operator",
		RunE: func(cmd *cobra.Command, args []string) error {
			selector, err := camelCtrl.WithIntegrationLabelsSelector()
			if err != nil {
				return errors.Wrap(err, "unable to compute cache's watch selector")
			}

			co.WatchSelectors = map[rtclient.Object]rtcache.ByObject{
				&corev1.Secret{}:     {Label: selector},
				&corev1.ConfigMap{}:  {Label: selector},
				&appsv1.Deployment{}: {Label: selector},
			}

			return controller.Start(co, func(manager manager.Manager, opts controller.Options) error {
				_, err := camelCtrl.NewReconciler(cmd.Context(), manager, camelCtrl.Options{})
				if err != nil {
					return errors.Wrap(err, "unable to set-up reconciler")
				}

				return nil
			})
		},
	}

	cmd.Flags().StringVar(
		&co.LeaderElectionID, "leader-election-id", co.LeaderElectionID, "The leader election ID of the operator.")
	cmd.Flags().StringVar(
		&co.LeaderElectionNamespace, "leader-election-namespace", co.LeaderElectionNamespace, "The leader election namespace.")
	cmd.Flags().BoolVar(
		&co.EnableLeaderElection, "leader-election", co.EnableLeaderElection, "Enable leader election for controller manager.")
	cmd.Flags().BoolVar(
		&co.ReleaseLeaderElectionOnCancel, "leader-election-release", co.ReleaseLeaderElectionOnCancel, "If the leader should step down voluntarily.")
	cmd.Flags().StringVar(
		&co.MetricsAddr, "metrics-bind-address", co.MetricsAddr, "The address the metric endpoint binds to.")
	cmd.Flags().StringVar(
		&co.ProbeAddr, "health-probe-bind-address", co.ProbeAddr, "The address the probe endpoint binds to.")
	cmd.Flags().StringVar(
		&co.PprofAddr, "pprof-bind-address", co.PprofAddr, "The address the pprof endpoint binds to.")

	fs := flag.NewFlagSet("", flag.PanicOnError)

	klog.InitFlags(fs)
	logger.Options.BindFlags(fs)

	cmd.Flags().AddGoFlagSet(fs)

	return &cmd
}
