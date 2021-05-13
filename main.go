package main

import (
	"log"
	"time"

	"github.com/google/gops/agent"
	"github.com/tsuru/kube-dbaas/web"

	extensionsruntime "github.com/tsuru/rpaas-operator/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatalf("could not initialize gops agent: %v", err)
	}

	syncPeriod := time.Minute
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             extensionsruntime.NewScheme(),
		MetricsBindAddress: ":8080",
		Port:               9443,
		SyncPeriod:         &syncPeriod,
	})
	if err != nil {
		log.Fatalf("unable to start manager: %v", err)
	}

	api, err := web.New(mgr.GetClient())
	if err != nil {
		log.Fatalf("could not create Kube-DBAAS API: %v", err)
	}
	if err := api.Start(); err != nil {
		log.Fatalf("could not start the Kube-DBAAS API: %v", err)
	}
}
