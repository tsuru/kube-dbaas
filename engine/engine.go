package engine

import (
	"context"

	"github.com/tsuru/kube-dbaas/types"
)

type Engine interface {
	CreateInstance(ctx context.Context, create *types.CreateArgs) error
	AppEnvVars(ctx context.Context, instanceName string) (map[string]string, error)
	Status(ctx context.Context, instanceName string) (address string, ready bool, err error)
	UpdateInstance()
	DeleteInstance()
}
