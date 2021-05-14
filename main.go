package main

import (
	"context"
	"log"
	"time"

	"github.com/google/gops/agent"
	"github.com/tsuru/kube-dbaas/web"

	mongov1 "github.com/mongodb/mongodb-kubernetes-operator/api/v1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatalf("could not initialize gops agent: %v", err)
	}

	scheme := runtime.NewScheme()
	utilruntime.Must(corev1.AddToScheme(scheme))
	utilruntime.Must(mongov1.AddToScheme(scheme))

	syncPeriod := time.Minute
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: ":8080",
		Port:               9443,
		SyncPeriod:         &syncPeriod,
	})
	if err != nil {
		log.Fatalf("unable to create manager: %v", err)
	}

	go func() {
		err = mgr.Start(context.Background())
		if err != nil {
			log.Fatalf("unable to start manager: %v", err)
		}
	}()

	api, err := web.New(mgr.GetClient())
	if err != nil {
		log.Fatalf("could not create Kube-DBAAS API: %v", err)
	}
	if err := api.Start(); err != nil {
		log.Fatalf("could not start the Kube-DBAAS API: %v", err)
	}
}
