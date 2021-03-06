package mongo

import (
	"context"
	"net/url"

	mongov1 "github.com/mongodb/mongodb-kubernetes-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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
						{Name: "readWrite", DB: create.Name},
					},
					ScramCredentialsSecretName: "my-scram",
				},
			},
			StatefulSetConfiguration: mongov1.StatefulSetConfiguration{
				SpecWrapper: mongov1.StatefulSetSpecWrapper{
					Spec: appsv1.StatefulSetSpec{
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									{
										Name: "mongod",
										Resources: corev1.ResourceRequirements{
											Limits: corev1.ResourceList{
												"cpu":    resource.MustParse("100m"),
												"memory": resource.MustParse("256Mi"),
											},
											Requests: corev1.ResourceList{
												"cpu":    resource.MustParse("100m"),
												"memory": resource.MustParse("256Mi"),
											},
										},
									},

									{
										Name: "mongodb-agent",
										Resources: corev1.ResourceRequirements{
											Limits: corev1.ResourceList{
												"cpu":    resource.MustParse("50m"),
												"memory": resource.MustParse("128Mi"),
											},
											Requests: corev1.ResourceList{
												"cpu":    resource.MustParse("50m"),
												"memory": resource.MustParse("128Mi"),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})
}

func (e *Engine) UpdateInstance() {

}

func (e *Engine) DeleteInstance() {

}

func (e *Engine) Status(ctx context.Context, instanceName string) (address string, ready bool, err error) {
	instance := &mongov1.MongoDBCommunity{}

	err = e.cli.Get(ctx, client.ObjectKey{
		Namespace: "kube-dbaas",
		Name:      instanceName,
	}, instance)

	if err != nil {
		return "", false, err
	}

	return instance.Status.MongoURI, instance.Status.Phase == mongov1.Running, nil
}

func (e *Engine) AppEnvVars(ctx context.Context, instanceName string) (map[string]string, error) {
	instance := &mongov1.MongoDBCommunity{}

	err := e.cli.Get(ctx, client.ObjectKey{
		Namespace: "kube-dbaas",
		Name:      instanceName,
	}, instance)

	if err != nil {
		return nil, err
	}

	u, err := url.Parse(instance.Status.MongoURI)

	if err != nil {
		return nil, err
	}

	u.User = url.UserPassword("tsuru", "admin")
	u.Path = "/" + instanceName

	return map[string]string{
		"DBAAS_MONGODB_HOSTS":    u.Host,
		"DBAAS_MONGODB_PASSWORD": "admin",
		"DBAAS_MONGODB_ENDPOINT": u.String(),
		"DBAAS_MONGODB_USER":     "tsuru",
		"DBAAS_MONGODB_PORT":     "27017",
	}, nil
}
