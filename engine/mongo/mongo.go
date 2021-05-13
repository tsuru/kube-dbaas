package mongo

import (
	"context"

	mongov1 "github.com/mongodb/mongodb-kubernetes-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/tsuru/kube-dbaas/engine"
	"github.com/tsuru/kube-dbaas/types"
)

type Engine struct {
	cli client.Client
}

func New(cli client.Client) engine.Engine {
	return &Engine{
		cli: cli,
	}
}

func (e *Engine) CreateInstance(ctx context.Context, create *types.CreateArgs) error {
	return e.cli.Create(ctx, &mongov1.MongoDBCommunity{
		ObjectMeta: metav1.ObjectMeta{
			Name:      create.Name,
			Namespace: "kube-dbaas",
			Labels: map[string]string{
				"tsuru.io/team": create.Team,
			},
		},
		Spec: mongov1.MongoDBCommunitySpec{
			Type:    "ReplicaSet",
			Version: "4.2.6",
			Members: 1,
			Security: mongov1.Security{
				Authentication: mongov1.Authentication{
					Modes: []mongov1.AuthMode{"CRAM"},
				},
			},
		},
	})

}

func (e *Engine) UpdateInstance() {

}

func (e *Engine) DeleteInstance() {

}

func (e *Engine) Status() {

}
