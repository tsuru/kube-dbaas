package engine

import (
	"context"

	"github.com/tsuru/kube-dbaas/types"
)

type Engine interface {
	CreateInstance(ctx context.Context, create *types.CreateArgs) error
	AppEnvVars(ctx context.Context, instanceName string) (map[string]string, error)
	UpdateInstance()
	DeleteInstance()
	Status()
}
