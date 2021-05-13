package engine

import (
	"context"

	"github.com/tsuru/kube-dbaas/types"
)

type Engine interface {
	CreateInstance(ctx context.Context, create *types.CreateArgs) error
	UpdateInstance()
	DeleteInstance()
	Status()
}
