package usecases

import (
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/AliceDiNunno/rack-controller/src/core/domain"
)

type Logger interface {
	Error(args ...interface{})
	Fatal(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
}

type UserRepo interface {
	IsEmpty() bool
	CreateUser(user *domain.User) *e.Error
	FindByMail(mail string) (*domain.User, *e.Error)
}

type UserTokenRepo interface {
	CreateToken(token *domain.AccessToken) *e.Error
	FindByToken(token string) (*domain.AccessToken, *e.Error)
}

type JwtSignatureRepo interface {
	SaveSignature(signature *domain.JwtSignature) *e.Error
	CheckIfSignatureExists(signature string) bool
}

type ProjectRepo interface {
}

type EnvironmentRepo interface {
}

type ServiceRepo interface {
}

type Kubernetes interface {
	GetNodes() ([]domain.Node, *e.Error)
	GetNode(string) (*domain.Node, *e.Error)

	ListDeployments(namespace string) ([]domain.Deployment, *e.Error)
	GetDeployment(namespace string, deploymentName string) (*domain.Deployment, *e.Error)
	GetPodsOfADeployment(namespace string, deploymentName string) ([]domain.Pod, *e.Error)
	GetEnvironmentOfADeployment(namespace string, deploymentName string) ([]domain.Environment, *e.Error)
	GetPortsOfADeployment(namespace string, deploymentName string) ([]domain.Port, *e.Error)
	GetConfigMapsOfADeployment(namespace string, name string) ([]string, *e.Error)
	GetSecretsOfADeployment(namespace string, name string) ([]string, *e.Error)
	DeleteDeployment(namespace string, name string) *e.Error
	RestartDeployment(namespace string, name string) *e.Error
	CreateDeployment(namespace string, request domain.DeploymentCreationRequest) *e.Error
	ExposePorts(namespace string, name string, ports []domain.Port) *e.Error

	GetPods(namespace string) ([]domain.Pod, *e.Error)
	GetPod(namespace string, podName string) (*domain.Pod, *e.Error)
	GetPodLogs(namespace string, podName string) (string, *e.Error)
	DeletePod(namespace string, podName string) *e.Error

	GetNamespaces() ([]string, *e.Error)

	GetConfigMapList(namespace string) ([]string, *e.Error)
	GetConfigMap(namespace string, name string) (domain.ConfigMap, *e.Error)
	CreateConfigMap(namespace string, request domain.ConfigMapCreationRequest) *e.Error
	UpdateConfigMap(namespace string, name string, request domain.ConfigMapUpdateRequest) *e.Error
	DeleteConfigMap(namespace string, name string) *e.Error

	GetSecretsList(namespace string) ([]string, *e.Error)
	GetSecret(namespace string, name string) (domain.Secret, *e.Error)
	CreateSecret(namespace string, request domain.SecretCreationRequest) *e.Error
	UpdateSecret(namespace string, name string, request domain.SecretUpdateRequest) *e.Error
	DeleteSecret(namespace string, name string) *e.Error
}

type interactor struct {
	userRepo        UserRepo
	userTokenRepo   UserTokenRepo
	jwtSignature    JwtSignatureRepo
	projectRepo     ProjectRepo
	environmentRepo EnvironmentRepo
	serviceRepo     ServiceRepo
	kubernetes      Kubernetes
}

func NewInteractor(u UserRepo, ut UserTokenRepo, js JwtSignatureRepo,
	repo ProjectRepo, env EnvironmentRepo, s ServiceRepo,
	k8s Kubernetes) interactor {
	return interactor{
		userRepo:        u,
		userTokenRepo:   ut,
		jwtSignature:    js,
		projectRepo:     repo,
		environmentRepo: env,
		serviceRepo:     s,
		kubernetes:      k8s,
	}
}
