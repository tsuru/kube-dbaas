package mongo

import (
	"context"

	mongov1 "github.com/mongodb/mongodb-kubernetes-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
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
	err := e.cli.Create(ctx, &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      create.Name + "-secret",
			Namespace: "kube-dbaas",
			Labels: map[string]string{
				"tsuru.io/team": create.Team,
			},
		},
		StringData: map[string]string{
			"password": "admin",
		},
	})
	if err != nil {
		return err
	}

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
					Modes: []mongov1.AuthMode{"SCRAM"},
				},
			},
			Users: []mongov1.MongoDBUser{
				{
					Name: "tsuru",
					DB:   "admin",
					PasswordSecretRef: mongov1.SecretKeyReference{
						Name: create.Name + "-secret",
					},
					Roles: []mongov1.Role{
						{Name: "clusterAdmin", DB: "admin"},
						{Name: "userAdminAnyDatabase", DB: "admin"},
					},
					ScramCredentialsSecretName: "my-scram",
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
