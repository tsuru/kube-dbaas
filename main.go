package main

import (
	"log"

	"github.com/google/gops/agent"
	"github.com/tsuru/kube-dbaas/web"
)

func main() {
	if err := agent.Listen(agent.Options{}); err != nil {
		log.Fatalf("could not initialize gops agent: %v", err)
	}

	api, err := web.New()
	if err != nil {
		log.Fatalf("could not create Kube-DBAAS API: %v", err)
	}
	if err := api.Start(); err != nil {
		log.Fatalf("could not start the Kube-DBAAS API: %v", err)
	}
}
