package controller

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/cache"

	"github.com/lburgazzoli/camel-go/pkg/controller/logger"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

var (
	Scheme = runtime.NewScheme()
	Log    = ctrl.Log.WithName("controller")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(Scheme))
}

func Start(options Options, setup func(manager.Manager, Options) error) error {
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&logger.Options)))

	ctx := ctrl.SetupSignalHandler()

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                        Scheme,
		HealthProbeBindAddress:        options.ProbeAddr,
		LeaderElection:                options.EnableLeaderElection,
		LeaderElectionID:              options.LeaderElectionID,
		LeaderElectionReleaseOnCancel: options.ReleaseLeaderElectionOnCancel,
		LeaderElectionNamespace:       options.LeaderElectionNamespace,

		Metrics: metricsserver.Options{
			BindAddress: options.MetricsAddr,
		},
		Cache: cache.Options{
			ByObject: options.WatchSelectors,
		},
		PprofBindAddress: options.PprofAddr,
	})

	if err != nil {
		return fmt.Errorf("unable to create manager: %w", err)
	}

	if err := setup(mgr, options); err != nil {
		return fmt.Errorf("unable to set up controllers: %w", err)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		return fmt.Errorf("unable to set up health check: %w", err)
	}

	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		return fmt.Errorf("unable to set up readiness check: %w", err)
	}

	Log.Info("starting manager")

	if err := mgr.Start(ctx); err != nil {
		return fmt.Errorf("problem running manager: %w", err)
	}

	return nil
}
