package operator

import (
	camelCtrl "github.com/lburgazzoli/camel-go/internal/controller/camel"
	"github.com/lburgazzoli/camel-go/pkg/controller"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	rtcache "sigs.k8s.io/controller-runtime/pkg/cache"
	rtclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func NewOperatorCmd() *cobra.Command {
	controllerOpts := controller.Options{
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

			controllerOpts.WatchSelectors = map[rtclient.Object]rtcache.ByObject{
				&corev1.Secret{}:     {Label: selector},
				&corev1.ConfigMap{}:  {Label: selector},
				&appsv1.Deployment{}: {Label: selector},
			}

			return controller.Start(controllerOpts, func(manager manager.Manager, opts controller.Options) error {
				_, err := camelCtrl.NewReconciler(cmd.Context(), manager, camelCtrl.Options{})
				if err != nil {
					return errors.Wrap(err, "unable to set-up reconciler")
				}

				return nil
			})
		},
	}

	cmd.Flags().StringVar(&controllerOpts.LeaderElectionID, "leader-election-id", controllerOpts.LeaderElectionID, "The leader election ID of the operator.")
	cmd.Flags().StringVar(&controllerOpts.LeaderElectionNamespace, "leader-election-namespace", controllerOpts.LeaderElectionNamespace, "The leader election namespace.")
	cmd.Flags().BoolVar(&controllerOpts.EnableLeaderElection, "leader-election", controllerOpts.EnableLeaderElection, "Enable leader election for controller manager.")
	cmd.Flags().BoolVar(&controllerOpts.ReleaseLeaderElectionOnCancel, "leader-election-release", controllerOpts.ReleaseLeaderElectionOnCancel, "If the leader should step down voluntarily.")

	cmd.Flags().StringVar(&controllerOpts.MetricsAddr, "metrics-bind-address", controllerOpts.MetricsAddr, "The address the metric endpoint binds to.")
	cmd.Flags().StringVar(&controllerOpts.ProbeAddr, "health-probe-bind-address", controllerOpts.ProbeAddr, "The address the probe endpoint binds to.")
	cmd.Flags().StringVar(&controllerOpts.PprofAddr, "pprof-bind-address", controllerOpts.PprofAddr, "The address the pprof endpoint binds to.")

	return &cmd
}
